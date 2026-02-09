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
  <div class="w-full overflow-x-auto pb-4">
    <div
      class="grid grid-cols-3 gap-8 items-center bg-white p-6 rounded-lg border border-purple-50 min-w-[900px]"
    >
      <!-- Column 1: Opening Matches -->
      <div class="flex flex-col gap-12 relative">
        <div v-if="m1" class="match-node match-m1 relative">
          <div
            class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
          >
            Opening A
          </div>
          <MatchCard :match="m1" :is-admin="isAdmin" @click="onMatchClick" />
        </div>
        <div v-if="m2" class="match-node match-m2 relative">
          <div
            class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
          >
            Opening B
          </div>
          <MatchCard :match="m2" :is-admin="isAdmin" @click="onMatchClick" />
        </div>
      </div>

      <!-- Column 2: Winners & Losers -->
      <div class="flex flex-col gap-12 relative">
        <!-- Larger gap to align with M1/M2 center? or just space out -->
        <div v-if="winners" class="match-node match-m3 relative">
          <div
            class="text-xs text-purple-600 mb-1 uppercase tracking-wider font-bold"
          >
            Winners Match
          </div>
          <MatchCard
            :match="winners"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="border-purple-200"
          />
        </div>
        <div v-if="losers" class="match-node match-m4 relative">
          <div
            class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
          >
            Elimination Match
          </div>
          <MatchCard
            :match="losers"
            :is-admin="isAdmin"
            @click="onMatchClick"
          />
        </div>
      </div>

      <!-- Column 3: Decider -->
      <div class="flex flex-col justify-center relative">
        <div v-if="decider" class="match-node match-m5 relative">
          <div
            class="text-xs text-orange-500 mb-1 uppercase tracking-wider font-bold"
          >
            Decider Match
          </div>
          <MatchCard
            :match="decider"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="border-orange-200"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Connector Lines */
.match-node::after {
  content: "";
  position: absolute;
  background-color: #e5e7eb; /* gray-200 */
  z-index: 0;
}

/* M1 connects to M3 (Down-Right) */
.match-m1::after {
  width: 2rem;
  height: 2px;
  top: 60%;
  right: -2rem;
}
/* We ideally want lines to connect nodes specifically, but strict CSS grid lines are hard without fixed heights.
   Visual approximation: Just show outgoing lines indicating flow. */

/* Refined Approach: Use borders on pseodo elements to create bracket shapes */

/* M1 -> M3 */
.match-m1::before {
  content: "";
  position: absolute;
  top: 50%;
  right: -2rem; /* Extend into gap */
  width: 2rem;
  height: 50%; /* Go down? */
  border-top: 2px solid #e2e8f0;
  border-right: 2px solid #e2e8f0;
  border-top-right-radius: 8px;
  transform: translateY(0);
  z-index: 0;
}
/* But M1 goes to M3 (Win) AND M4 (Lose). This is too complex for simple CSS lines without SVG. 
   Let's stick to simple "Flow" indicators */

.match-node {
  position: relative;
}

/* Outgoing line */
.match-node::before {
  content: "";
  position: absolute;
  top: 50%;
  right: -32px; /* Pull into gap */
  width: 32px;
  height: 2px;
  background: #e2e8f0;
  z-index: 0;
}

/* Incoming line for Gen 2/3 */
.match-m3::after,
.match-m4::after,
.match-m5::after {
  content: "";
  position: absolute;
  top: 50%;
  left: -32px;
  width: 32px;
  height: 2px;
  background: #e2e8f0;
  z-index: 0;
}
</style>
