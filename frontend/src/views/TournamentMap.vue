<script setup>
import { ref, onMounted, nextTick, onUnmounted, computed, watch } from "vue";
import api from "../services/api";
import MatchNode from "../components/MatchNode.vue";
import TheNavbar from "../components/TheNavbar.vue";
import ScoreModal from "../components/ScoreModal.vue";
import MatchDetailsModal from "../components/MatchDetailsModal.vue";

const props = defineProps({
  isAdmin: { type: Boolean, default: false },
});

const groups = ref([]);
const loading = ref(true);
const zoom = ref(0.7);
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

const isMatchDetailsOpen = ref(false);
const selectedMatchDetails = ref(null);

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

const drawGlobalPath = (p1, p2, color = "#e2e8f0") => {
  // Cubic Bezier Curve
  // M startX startY C cp1X cp1Y, cp2X cp2Y, endX endY
  // cp1 = startX + 50, startY
  // cp2 = endX - 50, endY
  const cp1X = p1.x + 50;
  const cp1Y = p1.y;
  const cp2X = p2.x - 50;
  const cp2Y = p2.y;

  return {
    d: `M ${p1.x} ${p1.y} C ${cp1X} ${cp1Y}, ${cp2X} ${cp2Y}, ${p2.x} ${p2.y}`,
    color: color,
  };
};

const updateGlobalPaths = async () => {
  await nextTick();
  const newPaths = [];

  // Group Connectors
  groups.value.forEach((g) => {
    if (g.name === "KNOCKOUT") return;

    // Layout-aware sides
    const isPoolA = g.pool === "Mesoneer";
    const outSide = isPoolA ? "right" : "left";
    const inSide = isPoolA ? "left" : "right";

    // Internal GSL connections
    const m1 = g.matches.find((m) => m.label === "M1");
    const m2 = g.matches.find((m) => m.label === "M2");
    const winners = g.matches.find((m) => m.label === "Winners");
    const losers = g.matches.find((m) => m.label === "Losers");
    const decider = g.matches.find((m) => m.label === "Decider");

    if (m1 && winners)
      newPaths.push(
        drawGlobalPath(
          getElementPoint(`match-${m1.id}`, outSide),
          getElementPoint(`match-${winners.id}`, inSide),
        ),
      );
    if (m1 && losers)
      newPaths.push(
        drawGlobalPath(
          getElementPoint(`match-${m1.id}`, outSide),
          getElementPoint(`match-${losers.id}`, inSide),
        ),
      );
    if (m2 && winners)
      newPaths.push(
        drawGlobalPath(
          getElementPoint(`match-${m2.id}`, outSide),
          getElementPoint(`match-${winners.id}`, inSide),
        ),
      );
    if (m2 && losers)
      newPaths.push(
        drawGlobalPath(
          getElementPoint(`match-${m2.id}`, outSide),
          getElementPoint(`match-${losers.id}`, inSide),
        ),
      );
    if (winners && decider)
      newPaths.push(
        drawGlobalPath(
          getElementPoint(`match-${winners.id}`, outSide),
          getElementPoint(`match-${decider.id}`, inSide),
        ),
      );
    if (losers && decider)
      newPaths.push(
        drawGlobalPath(
          getElementPoint(`match-${losers.id}`, outSide),
          getElementPoint(`match-${decider.id}`, inSide),
        ),
      );

    // Promotion paths to SF
    if (winners) {
      const pW = getElementPoint(`match-${winners.id}`, outSide);
      const targetSF = isPoolA ? sf1.value : sf2.value;
      const pSF = getElementPoint(`match-${targetSF?.id}`, inSide);
      if (pW && pSF) newPaths.push(drawGlobalPath(pW, pSF, "#8b5cf6"));
    }
    if (decider) {
      const pD = getElementPoint(`match-${decider.id}`, outSide);
      const targetSF = isPoolA ? sf2.value : sf1.value;
      const pSF = getElementPoint(`match-${targetSF?.id}`, inSide);
      if (pD && pSF) newPaths.push(drawGlobalPath(pD, pSF, "#8b5cf6"));
    }
  });

  // SF -> Final
  if (sf1.value && final.value)
    newPaths.push(
      drawGlobalPath(
        getElementPoint(`match-${sf1.value.id}`, "right"),
        getElementPoint(`match-${final.value.id}`, "left"),
        "#8b5cf6",
      ),
    );
  if (sf2.value && final.value)
    newPaths.push(
      drawGlobalPath(
        getElementPoint(`match-${sf2.value.id}`, "right"),
        getElementPoint(`match-${final.value.id}`, "left"),
        "#8b5cf6",
      ),
    );

  paths.value = newPaths;
};

const handleMatchClick = (match) => {
  if (props.isAdmin) {
    selectedMatch.value = match;
    isModalOpen.value = true;
  } else {
    // Open read-only modal
    selectedMatchDetails.value = match;
    isMatchDetailsOpen.value = true;
  }
};

const saveMatch = async (data) => {
  try {
    await api.post(`/matches/${data.id}`, data);
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
  setTimeout(updateGlobalPaths, 500);
});
onUnmounted(() => window.removeEventListener("resize", updateGlobalPaths));
</script>

<template>
  <div class="min-h-screen bg-dot-pattern font-outfit text-gray-900 pb-20">
    <TheNavbar />

    <!-- View Controls -->
    <div
      class="fixed bottom-10 left-10 z-50 flex items-center gap-4 bg-white/80 backdrop-blur-md border border-gray-100 px-4 py-2 rounded-full shadow-sm"
    >
      <button
        @click="zoom = Math.max(0.3, zoom - 0.1)"
        class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 transition-colors text-lg font-light"
      >
        -
      </button>
      <span
        class="text-[10px] font-black uppercase tracking-widest text-gray-400 w-12 text-center"
        >{{ Math.round(zoom * 100) }}%</span
      >
      <button
        @click="zoom = Math.min(1.5, zoom + 0.1)"
        class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 transition-colors text-lg font-light"
      >
        +
      </button>
    </div>

    <main
      class="origin-top-left transition-transform duration-300 px-20 pt-32"
      :style="{ transform: `scale(${zoom})`, width: `${100 / zoom}%` }"
    >
      <div
        v-if="loading"
        class="text-center py-40 animate-pulse text-xs font-black tracking-[0.5em] text-gray-300 uppercase"
      >
        Initializing Tournament Grid
      </div>

      <div
        v-else
        ref="containerRef"
        class="relative flex justify-between gap-40 min-w-[2000px]"
      >
        <!-- SVG Layer -->
        <svg class="absolute inset-0 w-full h-full pointer-events-none z-0">
          <path
            v-for="(p, i) in paths"
            :key="i"
            :d="p.d"
            :stroke="p.color"
            stroke-width="1.5"
            fill="none"
            class="transition-all duration-500"
          />
        </svg>

        <!-- Column: Session A -->
        <div class="flex flex-col gap-32 z-10">
          <div class="space-y-2 mb-10 pl-2">
            <h2
              class="text-xs font-black text-violet-600 uppercase tracking-[0.3em]"
            >
              Session A
            </h2>
            <div
              class="text-2xl font-black text-gray-900 italic uppercase tracking-tighter"
            >
              Mesoneer Base
            </div>
          </div>

          <div
            v-for="g in mesoneerGroups"
            :key="g.id"
            class="space-y-12 pb-20 border-b border-gray-50 last:border-0"
          >
            <div
              class="text-[10px] font-bold text-gray-300 uppercase tracking-widest"
            >
              {{ g.name }}
            </div>
            <div class="flex gap-16">
              <!-- Round 1 -->
              <div class="flex flex-col justify-between py-4">
                <MatchNode
                  v-if="g.matches[0]"
                  :match="g.matches.find((m) => m.label === 'M1')"
                  @click="handleMatchClick"
                />
                <MatchNode
                  v-if="g.matches[1]"
                  :match="g.matches.find((m) => m.label === 'M2')"
                  @click="handleMatchClick"
                />
              </div>
              <!-- Round 2 -->
              <div class="flex flex-col justify-between py-4">
                <MatchNode
                  v-if="g.matches[2]"
                  :match="g.matches.find((m) => m.label === 'Winners')"
                  @click="handleMatchClick"
                />
                <MatchNode
                  v-if="g.matches[3]"
                  :match="g.matches.find((m) => m.label === 'Losers')"
                  @click="handleMatchClick"
                />
              </div>
              <!-- Decider -->
              <div class="flex flex-col justify-center">
                <MatchNode
                  v-if="g.matches[4]"
                  :match="g.matches.find((m) => m.label === 'Decider')"
                  @click="handleMatchClick"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Column: Champions (Center) -->
        <div
          class="flex flex-col items-center justify-center gap-40 z-10 pt-20"
        >
          <div class="text-center mb-20">
            <h1
              class="text-6xl font-black text-gray-900 tracking-tighter italic uppercase"
            >
              Champions Stage
            </h1>
            <div class="h-1 w-20 bg-violet-600 mx-auto mt-4"></div>
          </div>

          <div class="flex flex-col items-center gap-24">
            <!-- SF1 -->
            <div class="space-y-4">
              <div
                class="text-[10px] font-black text-gray-300 uppercase tracking-[0.5em] text-center"
              >
                Semi Final 1
              </div>
              <MatchNode
                v-if="sf1"
                :match="sf1"
                @click="handleMatchClick"
                class="scale-110"
              />
            </div>

            <!-- Final -->
            <div class="space-y-6 pt-10">
              <div
                class="text-xs font-black text-violet-600 uppercase tracking-[1em] text-center"
              >
                Grand Final
              </div>
              <MatchNode
                v-if="final"
                :match="final"
                @click="handleMatchClick"
                class="scale-[1.6] shadow-2xl border-violet-100"
              />
            </div>

            <!-- SF2 -->
            <div class="space-y-4 pt-10">
              <div
                class="text-[10px] font-black text-gray-300 uppercase tracking-[0.5em] text-center"
              >
                Semi Final 2
              </div>
              <MatchNode
                v-if="sf2"
                :match="sf2"
                @click="handleMatchClick"
                class="scale-110"
              />
            </div>

            <!-- Bronze -->
            <div class="mt-20 opacity-50 space-y-2">
              <div
                class="text-[8px] font-black text-gray-400 uppercase tracking-widest text-center"
              >
                Bronze Match
              </div>
              <MatchNode
                v-if="bronze"
                :match="bronze"
                @click="handleMatchClick"
                class="scale-90"
              />
            </div>
          </div>
        </div>

        <!-- Column: Session B -->
        <div class="flex flex-col gap-32 z-10">
          <div class="space-y-2 mb-10 text-right pr-2">
            <h2
              class="text-xs font-black text-emerald-600 uppercase tracking-[0.3em]"
            >
              Session B
            </h2>
            <div
              class="text-2xl font-black text-gray-900 italic uppercase tracking-tighter"
            >
              Lab Colony
            </div>
          </div>

          <div
            v-for="g in labGroups"
            :key="g.id"
            class="space-y-12 pb-20 border-b border-gray-50 last:border-0 items-end flex flex-col"
          >
            <div
              class="text-[10px] font-bold text-gray-300 uppercase tracking-widest text-right"
            >
              {{ g.name }}
            </div>
            <div class="flex gap-16 flex-row-reverse">
              <!-- Round 1 -->
              <div class="flex flex-col justify-between py-4">
                <MatchNode
                  v-if="g.matches[0]"
                  :match="g.matches.find((m) => m.label === 'M1')"
                  @click="handleMatchClick"
                />
                <MatchNode
                  v-if="g.matches[1]"
                  :match="g.matches.find((m) => m.label === 'M2')"
                  @click="handleMatchClick"
                />
              </div>
              <!-- Round 2 -->
              <div class="flex flex-col justify-between py-4">
                <MatchNode
                  v-if="g.matches[2]"
                  :match="g.matches.find((m) => m.label === 'Winners')"
                  @click="handleMatchClick"
                />
                <MatchNode
                  v-if="g.matches[3]"
                  :match="g.matches.find((m) => m.label === 'Losers')"
                  @click="handleMatchClick"
                />
              </div>
              <!-- Decider -->
              <div class="flex flex-col justify-center">
                <MatchNode
                  v-if="g.matches[4]"
                  :match="g.matches.find((m) => m.label === 'Decider')"
                  @click="handleMatchClick"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <ScoreModal
      :is-open="isModalOpen"
      :match="selectedMatch"
      :is-admin="true"
      @close="isModalOpen = false"
      @save="saveMatch"
    />

    <MatchDetailsModal
      :is-open="isMatchDetailsOpen"
      :match="selectedMatchDetails"
      @close="isMatchDetailsOpen = false"
    />
  </div>
</template>

<style scoped>
/* No additional styles needed, using Tailwind */
</style>
