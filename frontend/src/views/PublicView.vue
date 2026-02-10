<script setup>
import { ref, onMounted, nextTick, onUnmounted, computed, watch } from "vue";
import api from "../services/api";
import GSLGrid from "../components/GSLGrid.vue";
import KnockoutGrid from "../components/KnockoutGrid.vue";
import TheNavbar from "../components/TheNavbar.vue";

const groups = ref([]);
const loading = ref(true);
const zoom = ref(0.8); // Default zoom out for overview

const mesoneerGroups = computed(() =>
  groups.value.filter((g) => g.pool === "Mesoneer" && g.name !== "KNOCKOUT"),
);
const labGroups = computed(() =>
  groups.value.filter((g) => g.pool === "Lab" && g.name !== "KNOCKOUT"),
);
const knockoutGroup = computed(() =>
  groups.value.find((g) => g.name === "KNOCKOUT"),
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

const getElementCenter = (id) => {
  const el = document.getElementById(id);
  if (!el || !containerRef.value) return null;
  const rect = el.getBoundingClientRect();
  const cRect = containerRef.value.getBoundingClientRect();
  return {
    x: rect.left - cRect.left + rect.width / 2,
    y: rect.top - cRect.top + rect.height / 2,
  };
};

const updatePaths = async () => {
  await nextTick();
  if (!knockoutGroup.value) {
    paths.value = [];
    return;
  }

  const newPaths = [];

  // Helper to find which group a team came from
  const findSourceGroup = (teamId) => {
    if (!teamId) return null;
    for (const g of groups.value) {
      if (g.name === "KNOCKOUT") continue;
      // Simple heuristic: check if team is in any match of this group
      for (const m of g.matches || []) {
        if (m.team_a_id === teamId || m.team_b_id === teamId) return g.id;
      }
    }
    return null;
  };

  // Parse Knockout Matches to find input teams
  const kMatches = knockoutGroup.value.matches || [];
  const sf1 = kMatches.find((m) => m.label === "SF1");
  const sf2 = kMatches.find((m) => m.label === "SF2");

  const targets = [
    { match: sf1, targetId: "knockout-sf1" }, // Need IDs on KnockoutGrid
    { match: sf2, targetId: "knockout-sf2" },
  ];

  targets.forEach(({ match, targetId }) => {
    if (!match) return;

    // Connect Team A
    if (match.team_a_id) {
      const groupA = findSourceGroup(match.team_a_id);
      if (groupA) {
        // In GSLGrid, we need to ID the "winner" box or similar.
        // For now, let's just target the whole group card center?
        // Or we add IDs to GSLGrid.
        // Let's target the group container.
        const startId = `group-${groupA}`;
        const p1 = getElementCenter(startId);
        const p2 = getElementCenter(targetId);
        if (p1 && p2) {
          newPaths.push(drawConnector(p1, p2));
        }
      }
    }

    // Connect Team B
    if (match.team_b_id) {
      const groupB = findSourceGroup(match.team_b_id);
      if (groupB) {
        const startId = `group-${groupB}`;
        const p1 = getElementCenter(startId);
        const p2 = getElementCenter(targetId);
        if (p1 && p2) {
          newPaths.push(drawConnector(p1, p2));
        }
      }
    }
  });

  paths.value = newPaths;
};

const drawConnector = (p1, p2) => {
  // Bezier from P1 to P2
  const cx = (p1.x + p2.x) / 2;
  return {
    d: `M ${p1.x} ${p1.y} C ${cx} ${p1.y}, ${cx} ${p2.y}, ${p2.x} ${p2.y}`,
    color: "#cbd5e1", // slate-300
  };
};

watch(groups, updatePaths, { deep: true });

onMounted(() => {
  fetchData();
  window.addEventListener("resize", updatePaths);
  setInterval(fetchData, 30000);
});

onUnmounted(() => {
  window.removeEventListener("resize", updatePaths);
});
</script>

<template>
  <div
    class="min-h-screen bg-tech-pattern overflow-hidden font-sans text-slate-300"
  >
    <TheNavbar />

    <!-- Zoom Controls -->
    <div class="fixed bottom-8 left-8 z-50 flex gap-2">
      <button
        @click="zoom -= 0.1"
        class="bg-slate-800 text-white p-2 rounded-full shadow-lg border border-slate-700 hover:bg-slate-700"
      >
        -
      </button>
      <div
        class="bg-slate-900 px-3 py-2 rounded-lg border border-slate-700 font-mono text-xs flex items-center"
      >
        {{ Math.round(zoom * 100) }}%
      </div>
      <button
        @click="zoom += 0.1"
        class="bg-slate-800 text-white p-2 rounded-full shadow-lg border border-slate-700 hover:bg-slate-700"
      >
        +
      </button>
    </div>

    <main
      class="p-10 transition-transform origin-top-left ease-out duration-300"
      :style="{
        transform: `scale(${zoom})`,
        width: `${100 / zoom}%`,
        height: `${100 / zoom}%`,
      }"
    >
      <div
        v-if="loading"
        class="text-center py-20 text-white text-xl animate-pulse"
      >
        Initializing Tournament System...
      </div>

      <div
        v-else
        ref="containerRef"
        class="relative flex gap-20 min-w-[1800px]"
      >
        <!-- SVG Layer (Global Connectors) -->
        <svg class="absolute inset-0 w-full h-full pointer-events-none z-0">
          <path
            v-for="(path, i) in paths"
            :key="i"
            :d="path.d"
            stroke="#cbd5e1"
            stroke-width="3"
            fill="none"
            stroke-opacity="0.2"
            stroke-dasharray="10"
          />
        </svg>

        <!-- Left Panel: Groups (Mesoneer Top, Lab Bottom) -->
        <div class="flex flex-col gap-24 w-3/4 z-10">
          <!-- Top Row: Mesoneer -->
          <div class="border-b border-slate-800/50 pb-12">
            <h3
              class="text-2xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-400 uppercase tracking-widest mb-8 pl-4 border-l-4 border-blue-500"
            >
              Session A: Mesoneer
            </h3>
            <div class="flex flex-wrap gap-12">
              <div
                v-for="group in mesoneerGroups"
                :key="group.id"
                :id="`group-${group.id}`"
                class="bg-white/95 backdrop-blur rounded-2xl shadow-xl overflow-hidden border border-slate-200 transform hover:scale-105 transition-transform duration-300"
              >
                <div
                  class="bg-slate-50 px-6 py-3 border-b border-slate-100 flex justify-between items-center"
                >
                  <span
                    class="font-bold text-slate-700 uppercase tracking-wider"
                    >{{ group.name }}</span
                  >
                  <span
                    class="text-xs font-mono text-blue-500 bg-blue-50 px-2 py-1 rounded"
                    >MESONEER</span
                  >
                </div>
                <GSLGrid :matches="group.matches" :group-id="group.id" />
              </div>
            </div>
          </div>

          <!-- Bottom Row: Lab -->
          <div>
            <h3
              class="text-2xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-teal-400 uppercase tracking-widest mb-8 pl-4 border-l-4 border-emerald-500"
            >
              Session B: Lab
            </h3>
            <div class="flex flex-wrap gap-12">
              <div
                v-for="group in labGroups"
                :key="group.id"
                :id="`group-${group.id}`"
                class="bg-white/95 backdrop-blur rounded-2xl shadow-xl overflow-hidden border border-slate-200 transform hover:scale-105 transition-transform duration-300"
              >
                <div
                  class="bg-slate-50 px-6 py-3 border-b border-slate-100 flex justify-between items-center"
                >
                  <span
                    class="font-bold text-slate-700 uppercase tracking-wider"
                    >{{ group.name }}</span
                  >
                  <span
                    class="text-xs font-mono text-emerald-500 bg-emerald-50 px-2 py-1 rounded"
                    >LAB</span
                  >
                </div>
                <GSLGrid :matches="group.matches" :group-id="group.id" />
              </div>
            </div>
          </div>
        </div>

        <!-- Right Panel: Knockout (Sticky) -->
        <div class="w-1/4 z-10 flex flex-col justify-center">
          <div class="sticky top-10">
            <h3
              class="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-amber-300 to-orange-500 uppercase tracking-widest mb-12 text-center drop-shadow-sm"
            >
              Champion's Stage
            </h3>

            <div
              v-if="knockoutGroup"
              id="knockout-stage"
              class="transform scale-110 origin-top-center"
            >
              <KnockoutGrid :matches="knockoutGroup.matches" />
            </div>
            <div
              v-else
              class="h-96 border-4 border-dashed border-slate-700 rounded-3xl flex items-center justify-center text-slate-600 font-mono text-lg"
            >
              Waiting for Group Stages...
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
