<script setup>
import { ref, onMounted, nextTick, onUnmounted, computed, watch } from "vue";
import api from "../services/api";
import GSLGrid from "../components/GSLGrid.vue";
import KnockoutGrid from "../components/KnockoutGrid.vue";
import TheNavbar from "../components/TheNavbar.vue";

const groups = ref([]);
const loading = ref(true);

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
  <div class="min-h-screen pb-12 bg-tech-pattern overflow-x-hidden">
    <TheNavbar />

    <main class="container-fluid px-4 py-8 relative">
      <div v-if="loading" class="text-center py-20 text-white">
        Loading Map...
      </div>

      <div
        v-else
        ref="containerRef"
        class="relative min-h-[800px] flex justify-between gap-8"
      >
        <!-- SVG Layer -->
        <svg class="absolute inset-0 w-full h-full pointer-events-none z-0">
          <path
            v-for="(path, i) in paths"
            :key="i"
            :d="path.d"
            stroke="#64748b"
            stroke-width="2"
            fill="none"
            stroke-opacity="0.4"
          />
        </svg>

        <!-- Left Column: Mesoneer -->
        <div class="flex flex-col gap-16 w-1/3 items-center z-10">
          <h3
            class="text-xl font-bold text-blue-400 uppercase tracking-widest mb-4"
          >
            Mesoneer Groups
          </h3>
          <div
            v-for="group in mesoneerGroups"
            :key="group.id"
            :id="`group-${group.id}`"
            class="scale-90 origin-top"
          >
            <div class="text-center text-white mb-2 font-bold">
              {{ group.name }}
            </div>
            <GSLGrid :matches="group.matches" :group-id="group.id" />
          </div>
        </div>

        <!-- Center Column: Knockout -->
        <div class="flex flex-col justify-center w-1/3 items-center z-10">
          <h3
            class="text-xl font-bold text-amber-400 uppercase tracking-widest mb-4"
          >
            The Finals
          </h3>
          <div v-if="knockoutGroup" class="scale-100" id="knockout-stage">
            <KnockoutGrid :matches="knockoutGroup.matches" />
            <!-- Note: KnockoutGrid needs to expose IDs for SF1/SF2 if we want precise connecting -->
            <!-- We will assume KnockoutGrid has slots or ids we can target, or we wrap it -->
          </div>
          <div v-else class="text-gray-500 italic">
            Knockout Stage not yet generated
          </div>
        </div>

        <!-- Right Column: Lab -->
        <div class="flex flex-col gap-16 w-1/3 items-center z-10">
          <h3
            class="text-xl font-bold text-emerald-400 uppercase tracking-widest mb-4"
          >
            Lab Groups
          </h3>
          <div
            v-for="group in labGroups"
            :key="group.id"
            :id="`group-${group.id}`"
            class="scale-90 origin-top"
          >
            <div class="text-center text-white mb-2 font-bold">
              {{ group.name }}
            </div>
            <GSLGrid :matches="group.matches" :group-id="group.id" />
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
