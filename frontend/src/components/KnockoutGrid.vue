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

// Filter matches
const sf1 = computed(() => props.matches.find((m) => m.label === "SF1"));
const sf2 = computed(() => props.matches.find((m) => m.label === "SF2"));
const bronze = computed(() => props.matches.find((m) => m.label === "Bronze"));
const final = computed(() => props.matches.find((m) => m.label === "Final"));

const onMatchClick = (match) => {
  emit("match-click", match);
};
</script>

<template>
  <div class="w-full overflow-x-auto pb-4">
    <div
      class="grid grid-cols-2 gap-16 items-center bg-white p-12 rounded-lg border border-purple-50 min-w-[800px] relative"
    >
      <!-- Semi Finals Column -->
      <div class="flex flex-col gap-24 relative z-10">
        <div v-if="sf1" class="relative">
          <div
            class="text-xs text-purple-600 mb-2 uppercase tracking-wider font-bold"
          >
            Semi-Final 1
          </div>
          <MatchCard :match="sf1" :is-admin="isAdmin" @click="onMatchClick" />
          <!-- Connector to Final -->
          <div class="absolute top-1/2 -right-8 w-8 h-0.5 bg-gray-200"></div>
          <div
            class="absolute top-1/2 -right-8 w-0.5 h-[150%] bg-gray-200 origin-top"
          ></div>
        </div>

        <div v-if="sf2" class="relative">
          <div
            class="text-xs text-purple-600 mb-2 uppercase tracking-wider font-bold"
          >
            Semi-Final 2
          </div>
          <MatchCard :match="sf2" :is-admin="isAdmin" @click="onMatchClick" />
          <!-- Connector to Final -->
          <div class="absolute top-1/2 -right-8 w-8 h-0.5 bg-gray-200"></div>
          <div
            class="absolute top-1/2 -right-8 w-0.5 h-[150%] bg-gray-200 origin-bottom transform -scale-y-100"
          ></div>
        </div>
      </div>

      <!-- Finals Column -->
      <div class="flex flex-col gap-12 relative z-10">
        <div v-if="final" class="relative">
          <!-- Incoming Connector Stub -->
          <div class="absolute top-1/2 -left-8 w-8 h-0.5 bg-gray-200"></div>

          <div
            class="text-xs text-amber-500 mb-2 uppercase tracking-wider font-black flex items-center gap-2"
          >
            <span>🏆 Grand Final</span>
          </div>
          <MatchCard
            :match="final"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="border-amber-400 ring-4 ring-amber-50"
          />
        </div>

        <div v-if="bronze" class="mt-8 opacity-80">
          <div
            class="text-xs text-gray-400 mb-2 uppercase tracking-wider font-bold"
          >
            🥉 Bronze Match
          </div>
          <MatchCard
            :match="bronze"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="border-dashed"
          />
        </div>
      </div>
    </div>
  </div>
</template>
