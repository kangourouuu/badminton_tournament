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
  <div
    class="w-full overflow-x-auto pb-8 pt-4 px-4 bg-white/50 rounded-xl border border-gray-100"
  >
    <div class="relative min-w-[800px] h-[500px]">
      <!-- SVG Connector Layer -->
      <svg class="absolute inset-0 w-full h-full pointer-events-none z-0">
        <defs>
          <marker
            id="arrowhead-ko"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
          >
            <polygon points="0 0, 10 3.5, 0 7" fill="#cbd5e1" />
          </marker>
        </defs>

        <!-- SF1 (Top-Left) to Final (Mid-Right) [Winner] -->
        <!-- Start: 280,80 -> Go Right 100 -> Go Down to 200 -> Go Right to 450 -->
        <path
          d="M 280 80 L 380 80 L 380 200 L 450 200"
          stroke="#cbd5e1"
          stroke-width="2"
          fill="none"
          marker-end="url(#arrowhead-ko)"
        />

        <!-- SF2 (Bot-Left) to Final (Mid-Right) [Winner] -->
        <!-- Start: 280,320 -> Go Right 100 -> Go Up to 200 -> Go Right to 450 -->
        <path
          d="M 280 320 L 380 320 L 380 200 L 450 200"
          stroke="#cbd5e1"
          stroke-width="2"
          fill="none"
          marker-end="url(#arrowhead-ko)"
        />

        <!-- SF1 (Top-Left) to Bronze (Bot-Right) [Loser - Dashed] -->
        <!-- 280,90 -> Right -> Down -> Right -->
        <path
          d="M 280 90 L 360 90 L 360 420 L 450 420"
          stroke="#f1f5f9"
          stroke-width="2"
          stroke-dasharray="4"
          fill="none"
        />

        <!-- SF2 (Bot-Left) to Bronze (Bot-Right) [Loser - Dashed] -->
        <!-- 280,330 -> Right -> Down -> Right -->
        <path
          d="M 280 330 L 360 330 L 360 420 L 450 420"
          stroke="#f1f5f9"
          stroke-width="2"
          stroke-dasharray="4"
          fill="none"
        />
      </svg>

      <!-- Matches Layer -->
      <div class="absolute inset-0 z-10">
        <!-- SF1 -->
        <div class="absolute top-10 left-8 w-64">
          <div
            class="text-xs text-purple-600 mb-1 uppercase tracking-wider font-bold"
          >
            Semi-Final 1
          </div>
          <MatchCard
            v-if="sf1"
            :match="sf1"
            :is-admin="isAdmin"
            @click="onMatchClick"
          />
          <div
            v-else
            class="h-24 bg-white/50 border-2 border-dashed border-gray-200 rounded flex items-center justify-center text-gray-400 text-xs font-mono"
          >
            Waiting for Qualifiers
          </div>
        </div>

        <!-- SF2 -->
        <div class="absolute top-70 left-8 w-64" style="top: 280px">
          <div
            class="text-xs text-purple-600 mb-1 uppercase tracking-wider font-bold"
          >
            Semi-Final 2
          </div>
          <MatchCard
            v-if="sf2"
            :match="sf2"
            :is-admin="isAdmin"
            @click="onMatchClick"
          />
          <div
            v-else
            class="h-24 bg-white/50 border-2 border-dashed border-gray-200 rounded flex items-center justify-center text-gray-400 text-xs font-mono"
          >
            Waiting for Qualifiers
          </div>
        </div>

        <!-- Final -->
        <div
          class="absolute top-40 left-120 w-64"
          style="left: 480px; top: 160px"
        >
          <div
            class="text-xs text-amber-500 mb-1 uppercase tracking-wider font-black flex items-center gap-2"
          >
            <span>🏆 Grand Final</span>
          </div>
          <MatchCard
            v-if="final"
            :match="final"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="border-amber-400 ring-4 ring-amber-50 shadow-xl"
          />
          <div
            v-else
            class="h-24 bg-amber-50/50 border-2 border-amber-200 rounded flex items-center justify-center text-amber-400 text-xs font-bold tracking-widest uppercase"
          >
            TBD
          </div>
        </div>

        <!-- Bronze -->
        <div
          v-if="bronze"
          class="absolute w-64"
          style="left: 480px; top: 380px"
        >
          <div
            class="text-xs text-gray-400 mb-1 uppercase tracking-wider font-bold text-center"
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
        <div
          v-else
          class="absolute w-64 h-24 bg-gray-50 rounded flex items-center justify-center text-gray-300 text-xs border border-dashed border-gray-200"
          style="left: 480px; top: 380px"
        >
          Waiting Bronze
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
