<script setup>
import { computed } from "vue";
import MatchCard from "./MatchCard.vue";

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
</script>

<template>
  <div
    class="w-full overflow-x-auto pb-8 pt-4 px-4 bg-gray-50/50 rounded-xl border border-gray-100"
  >
    <div class="relative min-w-[1000px] h-[400px]">
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
        </defs>

        <!-- M1 (Top-Left) to M3 (Top-Mid) [Winner] -->
        <path
          d="M 280 60 C 330 60, 330 60, 380 60"
          stroke="#cbd5e1"
          stroke-width="2"
          fill="none"
          marker-end="url(#arrowhead)"
        />
        <!-- M1 (Top-Left) to M4 (Bot-Mid) [Loser - Cross] -->
        <path
          d="M 280 70 C 330 70, 330 310, 380 310"
          stroke="#f1f5f9"
          stroke-width="2"
          stroke-dasharray="4"
          fill="none"
        />

        <!-- M2 (Bot-Left) to M3 (Top-Mid) [Winner - Cross] -->
        <path
          d="M 280 300 C 330 300, 330 70, 380 70"
          stroke="#cbd5e1"
          stroke-width="2"
          fill="none"
          marker-end="url(#arrowhead)"
        />
        <!-- M2 (Bot-Left) to M4 (Bot-Mid) [Loser] -->
        <path
          d="M 280 310 C 330 310, 330 310, 380 310"
          stroke="#f1f5f9"
          stroke-width="2"
          stroke-dasharray="4"
          fill="none"
        />

        <!-- M3 (Top-Mid) to M5 (Center-Right) [Loser] -->
        <path
          d="M 660 70 C 710 70, 710 185, 760 185"
          stroke="#f1f5f9"
          stroke-width="2"
          stroke-dasharray="4"
          fill="none"
        />

        <!-- M4 (Bot-Mid) to M5 (Center-Right) [Winner] -->
        <path
          d="M 660 300 C 710 300, 710 195, 760 195"
          stroke="#cbd5e1"
          stroke-width="2"
          fill="none"
          marker-end="url(#arrowhead)"
        />
      </svg>

      <!-- Grid Columns -->
      <div class="grid grid-cols-3 h-full absolute inset-0 z-10">
        <!-- COL 1: Opening (M1, M2) -->
        <div class="flex flex-col justify-between py-4 px-8">
          <!-- M1 -->
          <div v-if="m1" class="relative w-64">
            <div
              class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
            >
              Opening A
            </div>
            <MatchCard :match="m1" :is-admin="isAdmin" @click="onMatchClick" />
          </div>
          <div
            v-else
            class="h-24 w-64 bg-gray-100 rounded flex items-center justify-center text-gray-400 text-xs"
          >
            Waiting M1
          </div>

          <!-- M2 -->
          <div v-if="m2" class="relative w-64">
            <div
              class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
            >
              Opening B
            </div>
            <MatchCard :match="m2" :is-admin="isAdmin" @click="onMatchClick" />
          </div>
          <div
            v-else
            class="h-24 w-64 bg-gray-100 rounded flex items-center justify-center text-gray-400 text-xs"
          >
            Waiting M2
          </div>
        </div>

        <!-- COL 2: Winners (M3) & Losers (M4) -->
        <div class="flex flex-col justify-between py-4 px-8">
          <!-- M3 -->
          <div v-if="winners" class="relative w-64">
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
            class="h-24 w-64 bg-gray-100 rounded flex items-center justify-center text-gray-400 text-xs"
          >
            Waiting Winners
          </div>

          <!-- M4 -->
          <div v-if="losers" class="relative w-64">
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
            class="h-24 w-64 bg-gray-100 rounded flex items-center justify-center text-gray-400 text-xs"
          >
            Waiting Elimination
          </div>
        </div>

        <!-- COL 3: Decider (M5) -->
        <div class="flex flex-col justify-center px-8">
          <!-- M5 -->
          <div v-if="decider" class="relative w-64">
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
            class="h-24 w-64 bg-gray-100 rounded flex items-center justify-center text-gray-400 text-xs"
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
