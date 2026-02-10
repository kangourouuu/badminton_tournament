<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from "vue";
import MatchCard from "./MatchCard.vue";
import { useBracketConnectors } from "../composables/useBracketConnectors";

const props = defineProps({
  matches: {
    type: Array,
    default: () => [],
  },
  isAdmin: {
    type: Boolean,
    default: false,
  },
  groupId: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["match-click"]);

// Filter matches by label
const m1 = computed(() => props.matches.find((m) => m.label === "M1"));
const m2 = computed(() => props.matches.find((m) => m.label === "M2"));
const winners = computed(() =>
  props.matches.find((m) => m.label === "Winners"),
); // M3
const losers = computed(() => props.matches.find((m) => m.label === "Losers")); // M4
const decider = computed(() =>
  props.matches.find((m) => m.label === "Decider"),
); // M5

const onMatchClick = (match) => {
  emit("match-click", match);
};

// --- DYNAMIC SVG CONNECTORS ---
const containerRef = ref(null);
const { paths, updateConnectors } = useBracketConnectors(containerRef);

const updatePaths = async () => {
  const gid = props.groupId;

  // Helper for path styles
  const getStyle = (match, type) => {
    // type: 'win' or 'lose'
    if (!match || !match.winner_id) {
      // Pending: Faint Guidelines
      return {
        stroke: "#cbd5e1",
        opacity: 0.3,
        dash: type === "lose" ? "4" : "0",
        marker: type === "win" ? "url(#arrowhead)" : "",
      };
    }

    // Finished Match
    if (type === "win") {
      return {
        stroke: "#a855f7", // Purple-500
        opacity: 1.0,
        dash: "0",
        marker: "url(#arrowhead-active)",
      };
    } else {
      return {
        stroke: "#94a3b8", // Slate-400
        opacity: 0.4,
        dash: "4",
        marker: "",
      };
    }
  };

  const sM1 = getStyle(m1.value, "win");
  const sM1_L = getStyle(m1.value, "lose");

  const sM2 = getStyle(m2.value, "win");
  const sM2_L = getStyle(m2.value, "lose");

  const sM3 = getStyle(winners.value, "lose"); // Winner qualifies (no line), Loser to M5
  // Actually Winner of M3 goes to Knockout (handled by PublicView lines).
  // Only Loser of M3 goes to Decider (M5).

  const sM4 = getStyle(losers.value, "win"); // Winner to Decider (M5), Loser Out

  const connections = [
    // M1 -> Winners (Win)
    {
      startId: `gsl-${gid}-m1`,
      endId: `gsl-${gid}-winners`,
      type: "orthogonal",
      ...sM1,
    },
    // M1 -> Losers (Lose)
    {
      startId: `gsl-${gid}-m1`,
      endId: `gsl-${gid}-losers`,
      type: "orthogonal",
      ...sM1_L,
    },

    // M2 -> Winners (Win)
    {
      startId: `gsl-${gid}-m2`,
      endId: `gsl-${gid}-winners`,
      type: "orthogonal",
      ...sM2,
    },
    // M2 -> Losers (Lose)
    {
      startId: `gsl-${gid}-m2`,
      endId: `gsl-${gid}-losers`,
      type: "orthogonal",
      ...sM2_L,
    },

    // Winners -> Decider (Lose)
    {
      startId: `gsl-${gid}-winners`,
      endId: `gsl-${gid}-decider`,
      type: "orthogonal",
      ...sM3,
    },
    // Losers -> Decider (Win)
    {
      startId: `gsl-${gid}-losers`,
      endId: `gsl-${gid}-decider`,
      type: "orthogonal",
      ...sM4,
    },
  ];
  await updateConnectors(connections);
};

watch(() => props.matches, updatePaths, { deep: true });

onMounted(() => {
  window.addEventListener("resize", updatePaths);
  setTimeout(updatePaths, 100);
});

onUnmounted(() => {
  window.removeEventListener("resize", updatePaths);
});
</script>

<template>
  <div class="w-full overflow-x-auto pb-8 pt-4 px-4">
    <div ref="containerRef" class="relative min-w-[1000px] h-[400px]">
      <!-- SVG Connector Layer -->
      <svg class="absolute inset-0 w-full h-full pointer-events-none z-0">
        <defs>
          <marker
            id="arrowhead"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
          >
            <polygon points="0 0, 10 3.5, 0 7" fill="#cbd5e1" />
          </marker>
          <marker
            id="arrowhead-active"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
          >
            <polygon points="0 0, 10 3.5, 0 7" fill="#a855f7" />
          </marker>
        </defs>

        <path
          v-for="(path, i) in paths"
          :key="i"
          :d="path.d"
          :stroke="path.stroke || '#cbd5e1'"
          stroke-width="2"
          fill="none"
          :stroke-opacity="path.opacity || 1"
          :stroke-dasharray="path.dash || '0'"
          :marker-end="path.marker"
        />
      </svg>

      <!-- Grid Columns -->
      <div class="grid grid-cols-3 h-full absolute inset-0 z-10">
        <!-- COL 1: Opening (M1, M2) -->
        <div class="flex flex-col justify-between py-4 px-8">
          <!-- M1 -->
          <div :id="`gsl-${groupId}-m1`" v-if="m1" class="relative w-64">
            <div
              class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
            >
              Opening A
            </div>
            <MatchCard :match="m1" :is-admin="isAdmin" @click="onMatchClick" />
          </div>
          <div
            v-else
            :id="`gsl-${groupId}-m1`"
            class="h-24 w-64 bg-white/50 border-2 border-dashed border-gray-200 rounded flex items-center justify-center text-gray-400 text-xs font-mono"
          >
            Waiting M1
          </div>

          <!-- M2 -->
          <div :id="`gsl-${groupId}-m2`" v-if="m2" class="relative w-64">
            <div
              class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
            >
              Opening B
            </div>
            <MatchCard :match="m2" :is-admin="isAdmin" @click="onMatchClick" />
          </div>
          <div
            v-else
            :id="`gsl-${groupId}-m2`"
            class="h-24 w-64 bg-white/50 border-2 border-dashed border-gray-200 rounded flex items-center justify-center text-gray-400 text-xs font-mono"
          >
            Waiting M2
          </div>
        </div>

        <!-- COL 2: Winners (M3) & Losers (M4) -->
        <div class="flex flex-col justify-between py-4 px-8">
          <!-- M3 -->
          <div
            :id="`gsl-${groupId}-winners`"
            v-if="winners"
            class="relative w-64"
          >
            <div
              class="text-xs text-purple-600 mb-1 uppercase tracking-wider font-bold"
            >
              Winners Match
            </div>
            <MatchCard
              :match="winners"
              :is-admin="isAdmin"
              @click="onMatchClick"
              class="border-purple-200 shadow-purple-50"
            />
          </div>
          <div
            v-else
            :id="`gsl-${groupId}-winners`"
            class="h-24 w-64 bg-purple-50/30 border-2 border-dashed border-purple-100 rounded flex items-center justify-center text-purple-300 text-xs font-mono"
          >
            Waiting Winners
          </div>

          <!-- M4 -->
          <div
            :id="`gsl-${groupId}-losers`"
            v-if="losers"
            class="relative w-64"
          >
            <div
              class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
            >
              Elimination
            </div>
            <MatchCard
              :match="losers"
              :is-admin="isAdmin"
              @click="onMatchClick"
            />
          </div>
          <div
            v-else
            :id="`gsl-${groupId}-losers`"
            class="h-24 w-64 bg-white/50 border-2 border-dashed border-gray-200 rounded flex items-center justify-center text-gray-400 text-xs font-mono"
          >
            Waiting Elimination
          </div>
        </div>

        <!-- COL 3: Decider (M5) -->
        <div class="flex flex-col justify-center px-8">
          <!-- M5 -->
          <div
            :id="`gsl-${groupId}-decider`"
            v-if="decider"
            class="relative w-64"
          >
            <div
              class="text-xs text-orange-500 mb-1 uppercase tracking-wider font-bold"
            >
              Decider Match
            </div>
            <MatchCard
              :match="decider"
              :is-admin="isAdmin"
              @click="onMatchClick"
              class="border-orange-200 shadow-orange-50"
            />
          </div>
          <div
            v-else
            :id="`gsl-${groupId}-decider`"
            class="h-24 w-64 bg-orange-50/30 border-2 border-dashed border-orange-100 rounded flex items-center justify-center text-orange-300 text-xs font-mono"
          >
            Waiting Decider
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Scoped styles mainly for specific overrides, layout handled by Tailwind */
</style>
