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
  <div
    class="grid grid-cols-1 md:grid-cols-3 gap-8 items-center bg-white p-6 rounded-lg border border-purple-50"
  >
    <!-- Column 1: Opening Matches -->
    <div class="flex flex-col gap-12">
      <div v-if="m1">
        <div
          class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
        >
          Opening A
        </div>
        <MatchCard :match="m1" :is-admin="isAdmin" @click="onMatchClick" />
      </div>
      <div v-if="m2">
        <div
          class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
        >
          Opening B
        </div>
        <MatchCard :match="m2" :is-admin="isAdmin" @click="onMatchClick" />
      </div>
    </div>

    <!-- Column 2: Winners & Losers -->
    <div class="flex flex-col gap-24">
      <!-- Larger gap to align with M1/M2 center? or just space out -->
      <div v-if="winners">
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
      <div v-if="losers">
        <div
          class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-semibold"
        >
          Elimination Match
        </div>
        <MatchCard :match="losers" :is-admin="isAdmin" @click="onMatchClick" />
      </div>
    </div>

    <!-- Column 3: Decider -->
    <div class="flex flex-col justify-center">
      <div v-if="decider">
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
</template>
