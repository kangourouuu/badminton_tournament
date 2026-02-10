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

    // Connect Team A (Source Group A)
    // We assume SF1 Team A comes from Group A (Mesoneer) or P1
    // Logic: Connect Group Card Center to SF1 Input

    // Heuristic: Connect Logic
    // Group A matches -> SF1
    // Group B matches -> SF2
    // We iterate ALL groups to draw lines to SF
    // If Group Pool = Mesoneer -> SF1? No, logic is complex.
    // Let's draw generic "Qualification Lines" from all groups to the center.
  });

  // NEW LOGIC: Draw lines from Group Containers to Knockout Bracket
  // Mesoneer Groups -> SF1 / SF2 (Flowing Right)
  // Lab Groups -> SF1 / SF2 (Flowing Right)

  mesoneerGroups.value.forEach((g) => {
    const startId = `group-${g.id}`;
    const endId = `knockout-stage`; // Target the whole stage or specific SF
    const p1 = getElementCenter(startId);
    // Target SF1 for Mesoneer (Visual approximation)
    const p2 = getElementCenter("knockout-sf1");

    if (p1 && p2) {
      newPaths.push(drawFlowConnector(p1, p2, "Mesoneer"));
    }
  });

  labGroups.value.forEach((g) => {
    const startId = `group-${g.id}`;
    // Target SF2 for Lab (Visual approximation)
    const p2 = getElementCenter("knockout-sf2");
    if (getElementCenter(startId) && p2) {
      newPaths.push(drawFlowConnector(getElementCenter(startId), p2, "Lab"));
    }
  });

  paths.value = newPaths;
};

const drawFlowConnector = (p1, p2, type) => {
  // Sigmoid / Bezier Curve
  const cx = (p1.x + p2.x) / 2;
  return {
    d: `M ${p1.x} ${p1.y} C ${cx} ${p1.y}, ${cx} ${p2.y}, ${p2.x} ${p2.y}`,
    color: type === "Mesoneer" ? "#60a5fa" : "#34d399", // Blue / Emerald
    dash: "5,5",
    width: 2,
  };
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
            :stroke="path.color"
            :stroke-width="path.width"
            fill="none"
            :stroke-dasharray="path.dash"
            opacity="0.4"
            marker-end="url(#arrowhead-flow)"
          />
          <defs>
            <marker
              id="arrowhead-flow"
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
                class="rounded-2xl p-4 transition-transform duration-300 hover:scale-[1.01]"
              >
                <div
                  class="px-2 py-2 mb-2 border-b border-indigo-500/30 flex justify-between items-center"
                >
                  <span
                    class="font-black text-2xl text-indigo-100 uppercase tracking-widest font-outfit"
                    >{{ group.name }}</span
                  >
                  <span
                    class="text-[10px] font-bold tracking-widest text-blue-300 bg-blue-900/40 px-2 py-1 rounded border border-blue-500/20"
                    >MESONEER</span
                  >
                </div>
                <div class="p-4">
                  <GSLGrid :matches="group.matches" :group-id="group.id" />
                </div>
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
                class="rounded-2xl p-4 transition-transform duration-300 hover:scale-[1.01]"
              >
                <div
                  class="px-2 py-2 mb-2 border-b border-emerald-500/30 flex justify-between items-center"
                >
                  <span
                    class="font-black text-2xl text-emerald-100 uppercase tracking-widest font-outfit"
                    >{{ group.name }}</span
                  >
                  <span
                    class="text-[10px] font-bold tracking-widest text-emerald-300 bg-emerald-900/40 px-2 py-1 rounded border border-emerald-500/20"
                    >LAB</span
                  >
                </div>
                <div class="p-4">
                  <GSLGrid :matches="group.matches" :group-id="group.id" />
                </div>
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
              id="knockout-stage"
              class="transform scale-110 origin-top-center"
            >
              <KnockoutGrid
                :matches="knockoutGroup ? knockoutGroup.matches : []"
              />
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
