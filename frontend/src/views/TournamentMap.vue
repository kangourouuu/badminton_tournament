<script setup>
import { ref, onMounted, nextTick, onUnmounted, computed, watch } from "vue";
import api from "../services/api";
import GroupCard from "../components/GroupCard.vue";
import MatchCard from "../components/MatchCard.vue";
import TheNavbar from "../components/TheNavbar.vue";
import ScoreModal from "../components/ScoreModal.vue";

const groups = ref([]);
const loading = ref(true);
const zoom = ref(0.6); // Wide overview zoom
const isAdmin = ref(false); // Public by default
const selectedMatch = ref(null);
const isModalOpen = ref(false);

const mesoneerGroups = computed(() =>
  groups.value.filter((g) => g.pool === "Mesoneer" && g.name !== "KNOCKOUT"),
);
const labGroups = computed(() =>
  groups.value.filter((g) => g.pool === "Lab" && g.name !== "KNOCKOUT"),
);

const knockoutMatches = computed(() => {
  const kg = groups.value.find((g) => g.name === "KNOCKOUT");
  return kg ? kg.matches : [];
});

const sf1 = computed(() =>
  knockoutMatches.value.find((m) => m.label === "SF1"),
);
const sf2 = computed(() =>
  knockoutMatches.value.find((m) => m.label === "SF2"),
);
const final = computed(() =>
  knockoutMatches.value.find((m) => m.label === "Final"),
);
const bronze = computed(() =>
  knockoutMatches.value.find((m) => m.label === "Bronze"),
);

const fetchData = async () => {
  try {
    const response = await api.get("/groups");
    groups.value = response.data || [];
  } catch (err) {
    console.error("Failed to fetch data", err);
  } finally {
    loading.value = false;
  }
};

// SVG Logic
const paths = ref([]);
const containerRef = ref(null);

const getElementPoint = (id, side = "center") => {
  const el = document.getElementById(id);
  if (!el || !containerRef.value) return null;
  const rect = el.getBoundingClientRect();
  const cRect = containerRef.value.getBoundingClientRect();

  const x =
    side === "right"
      ? rect.right - cRect.left
      : side === "left"
        ? rect.left - cRect.left
        : rect.left - cRect.left + rect.width / 2;
  const y = rect.top - cRect.top + rect.height / 2;
  return { x, y };
};

const drawGlobalPath = (p1, p2, type) => {
  const midX = p1.x + (p2.x - p1.x) * 0.5;
  return {
    d: `M ${p1.x} ${p1.y} L ${midX} ${p1.y} L ${midX} ${p2.y} L ${p2.x} ${p2.y}`,
    color: type === "Mesoneer" ? "#60a5fa" : "#34d399",
  };
};

const updateGlobalPaths = async () => {
  await nextTick();
  const newPaths = [];

  // 1. Connect Mesoneer Groups to SF1/SF2
  mesoneerGroups.value.forEach((g) => {
    // Top 1 -> SF1
    const pWinners = getElementPoint(`gsl-${g.id}-winners`, "right");
    const pSF1 = getElementPoint("ko-sf1", "left");
    if (pWinners && pSF1)
      newPaths.push(drawGlobalPath(pWinners, pSF1, "Mesoneer"));

    // Top 2 -> SF2
    const pDecider = getElementPoint(`gsl-${g.id}-decider`, "right");
    const pSF2 = getElementPoint("ko-sf2", "left");
    if (pDecider && pSF2)
      newPaths.push(drawGlobalPath(pDecider, pSF2, "Mesoneer"));
  });

  // 2. Connect Lab Groups to SF1/SF2
  labGroups.value.forEach((g) => {
    // Top 1 -> SF2
    const pWinners = getElementPoint(`gsl-${g.id}-winners`, "left");
    const pSF2 = getElementPoint("ko-sf2", "right");
    if (pWinners && pSF2) newPaths.push(drawGlobalPath(pWinners, pSF2, "Lab"));

    // Top 2 -> SF1
    const pDecider = getElementPoint(`gsl-${g.id}-decider`, "left");
    const pSF1 = getElementPoint("ko-sf1", "right");
    if (pDecider && pSF1) newPaths.push(drawGlobalPath(pDecider, pSF1, "Lab"));
  });

  // 3. SF -> Final (Winner flow)
  const pSF1 = getElementPoint("ko-sf1", "right");
  const pSF2 = getElementPoint("ko-sf2", "right");
  const pFinalL = getElementPoint("ko-final", "left");
  if (pSF1 && pFinalL) newPaths.push(drawGlobalPath(pSF1, pFinalL, "Mesoneer"));
  if (pSF2 && pFinalL) newPaths.push(drawGlobalPath(pSF2, pFinalL, "Lab"));

  // 4. SF -> Bronze (Loser flow - using different color)
  const pBronzeL = getElementPoint("ko-bronze", "left");
  if (pSF1 && pBronzeL) {
    const p = drawGlobalPath(pSF1, pBronzeL, "Other");
    p.color = "#94a3b8";
    newPaths.push(p);
  }
  if (pSF2 && pBronzeL) {
    const p = drawGlobalPath(pSF2, pBronzeL, "Other");
    p.color = "#94a3b8";
    newPaths.push(p);
  }

  paths.value = newPaths;
};

const handleMatchClick = (match) => {
  // Only admin or video link?
  // Let's allow public to open video, admin to edit
  if (localStorage.getItem("token")) {
    selectedMatch.value = match;
    isModalOpen.value = true;
  }
};

const saveMatch = async (data) => {
  try {
    await api.put(`/matches/${data.id}`, data);
    isModalOpen.value = false;
    fetchData();
  } catch (err) {
    alert("Failed to save: " + err.message);
  }
};

watch(groups, updateGlobalPaths, { deep: true });
onMounted(() => {
  fetchData();
  window.addEventListener("resize", updateGlobalPaths);
});
onUnmounted(() => window.removeEventListener("resize", updateGlobalPaths));
</script>

<template>
  <div
    class="min-h-screen bg-tech-pattern overflow-hidden font-sans text-slate-300"
  >
    <TheNavbar />

    <!-- Zoom UI -->
    <div class="fixed bottom-8 left-8 z-50 flex gap-2">
      <button
        @click="zoom -= 0.1"
        class="bg-slate-800 p-2 rounded-full border border-slate-700 hover:bg-slate-700"
      >
        -
      </button>
      <div
        class="bg-slate-900 px-3 py-2 rounded-lg border border-slate-700 font-mono text-xs"
      >
        {{ Math.round(zoom * 100) }}%
      </div>
      <button
        @click="zoom += 0.1"
        class="bg-slate-800 p-2 rounded-full border border-slate-700 hover:bg-slate-700"
      >
        +
      </button>
    </div>

    <main
      class="p-20 transition-transform origin-top-left"
      :style="{ transform: `scale(${zoom})`, width: `${100 / zoom}%` }"
    >
      <div
        v-if="loading"
        class="text-center py-40 animate-pulse text-2xl font-black"
      >
        SYNCHRONIZING GRAND MAP...
      </div>

      <div
        v-else
        ref="containerRef"
        class="relative flex justify-between gap-16"
      >
        <!-- SVG Global Layer -->
        <svg class="absolute inset-0 w-full h-full pointer-events-none z-0">
          <path
            v-for="(p, i) in paths"
            :key="i"
            :d="p.d"
            :stroke="p.color"
            stroke-width="3"
            fill="none"
            opacity="0.3"
            marker-end="url(#arrowhead-grand)"
          />
          <defs>
            <marker
              id="arrowhead-grand"
              markerWidth="10"
              markerHeight="7"
              refX="9"
              refY="3.5"
              orient="auto"
            >
              <polygon points="0 0, 10 3.5, 0 7" fill="#cbd5e1" />
            </marker>
          </defs>
        </svg>

        <!-- Column 1: Mesoneer -->
        <div class="flex flex-col gap-12 z-10">
          <h2
            class="text-4xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-500 uppercase italic pl-4 border-l-8 border-blue-600 mb-10"
          >
            Session A
          </h2>
          <GroupCard
            v-for="g in mesoneerGroups"
            :key="g.id"
            :group="g"
            side="left"
            @match-click="handleMatchClick"
          />
        </div>

        <!-- Column 2: Champions Stage (SF / FINAL) -->
        <div class="flex flex-col justify-center gap-32 z-10 pt-48">
          <div class="text-center mb-10">
            <h2
              class="text-6xl font-black text-transparent bg-clip-text bg-gradient-to-b from-amber-300 to-orange-600 uppercase italic tracking-tighter"
            >
              Champions Stage
            </h2>
          </div>

          <div class="flex flex-col gap-24 items-center">
            <!-- Semi Final 1 -->
            <div class="relative">
              <div
                class="text-[10px] font-bold text-amber-500 uppercase tracking-widest mb-2 text-center"
              >
                Semi Final 1
              </div>
              <MatchCard
                id="ko-sf1"
                v-if="sf1"
                :match="sf1"
                @click="handleMatchClick"
                class="scale-110"
              />
            </div>

            <!-- Grand Final -->
            <div class="relative py-12">
              <div
                class="text-xs font-black text-amber-500 uppercase tracking-[0.5em] mb-4 text-center"
              >
                🏆 Grand Final 🏆
              </div>
              <MatchCard
                id="ko-final"
                v-if="final"
                :match="final"
                @click="handleMatchClick"
                class="scale-150 shadow-2xl ring-4 ring-amber-500/20"
              />
            </div>

            <!-- Semi Final 2 -->
            <div class="relative">
              <div
                class="text-[10px] font-bold text-amber-500 uppercase tracking-widest mb-2 text-center"
              >
                Semi Final 2
              </div>
              <MatchCard
                id="ko-sf2"
                v-if="sf2"
                :match="sf2"
                @click="handleMatchClick"
                class="scale-110"
              />
            </div>

            <!-- Bronze -->
            <div class="relative mt-12 opacity-80">
              <div
                class="text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-2 text-center"
              >
                Bronze Match
              </div>
              <MatchCard
                id="ko-bronze"
                v-if="bronze"
                :match="bronze"
                @click="handleMatchClick"
                class="scale-90"
              />
            </div>
          </div>
        </div>

        <!-- Column 3: Lab -->
        <div class="flex flex-col gap-12 z-10">
          <h2
            class="text-4xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-teal-500 uppercase italic pr-4 border-r-8 border-emerald-600 mb-10 text-right"
          >
            Session B
          </h2>
          <GroupCard
            v-for="g in labGroups"
            :key="g.id"
            :group="g"
            side="right"
            @match-click="handleMatchClick"
          />
        </div>
      </div>
    </main>

    <ScoreModal
      :is-open="isModalOpen"
      :match="selectedMatch"
      @close="isModalOpen = false"
      @save="saveMatch"
    />
  </div>
</template>

<style scoped>
.bg-tech-pattern {
  background-image: radial-gradient(
    circle at 2px 2px,
    rgba(255, 255, 255, 0.05) 1px,
    transparent 0
  );
  background-size: 32px 32px;
  background-color: #0f172a;
}
</style>
