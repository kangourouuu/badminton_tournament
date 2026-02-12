<script setup>
import { computed, ref, onMounted, onUnmounted, nextTick, watch } from "vue";
import MatchCard from "./MatchCard.vue";
import { useBracketConnectors } from "../composables/useBracketConnectors";

const props = defineProps({
  group: {
    type: Object,
    required: true,
  },
  isAdmin: {
    type: Boolean,
    default: false,
  },
  side: {
    type: String,
    default: "left", // 'left' or 'right'
  },
});

const emit = defineEmits(["match-click"]);

// Filter matches by label
const matches = computed(() => props.group.matches || []);
const m1 = computed(() => matches.value.find((m) => m.label === "M1"));
const m2 = computed(() => matches.value.find((m) => m.label === "M2"));
const winners = computed(() =>
  matches.value.find((m) => m.label === "Winners"),
);
const losers = computed(() => matches.value.find((m) => m.label === "Losers"));
const decider = computed(() =>
  matches.value.find((m) => m.label === "Decider"),
);

const onMatchClick = (match) => {
  emit("match-click", match);
};

// --- SVG CONNECTORS ---
const containerRef = ref(null);
const { paths, updateConnectors } = useBracketConnectors(containerRef);

const updatePaths = async () => {
  const gid = props.group.id;
  const connections = [
    {
      startId: `gsl-${gid}-m1`,
      endId: `gsl-${gid}-winners`,
      type: "orthogonal",
    },
    {
      startId: `gsl-${gid}-m1`,
      endId: `gsl-${gid}-losers`,
      type: "orthogonal",
    },
    {
      startId: `gsl-${gid}-m2`,
      endId: `gsl-${gid}-winners`,
      type: "orthogonal",
    },
    {
      startId: `gsl-${gid}-m2`,
      endId: `gsl-${gid}-losers`,
      type: "orthogonal",
    },
    {
      startId: `gsl-${gid}-winners`,
      endId: `gsl-${gid}-decider`,
      type: "orthogonal",
    },
    {
      startId: `gsl-${gid}-losers`,
      endId: `gsl-${gid}-decider`,
      type: "orthogonal",
    },

    // Elimination paths
    {
      startId: `gsl-${gid}-losers`,
      endId: `elim-m4-${gid}`,
      type: "straight",
      stroke: "#ef4444",
      opacity: 0.6,
      marker: "url(#crosshead)",
    },
    {
      startId: `gsl-${gid}-decider`,
      endId: `elim-m5-${gid}`,
      type: "straight",
      stroke: "#ef4444",
      opacity: 0.6,
      marker: "url(#crosshead)",
    },
  ];
  await updateConnectors(connections);
};

watch(() => props.group.matches, updatePaths, { deep: true });
onMounted(() => {
  window.addEventListener("resize", updatePaths);
  setTimeout(updatePaths, 150);
});
onUnmounted(() => window.removeEventListener("resize", updatePaths));
</script>

<template>
  <div class="relative group-card p-4 transition-all duration-300">
    <!-- Group Header (Floating) -->
    <div
      class="mb-4 flex flex-col"
      :class="side === 'right' ? 'items-end' : 'items-start'"
    >
      <span
        class="text-3xl font-black text-white/90 uppercase tracking-[0.2em] italic font-outfit"
      >
        {{ group.name }}
      </span>
      <div
        class="h-1 w-24 rounded-full bg-gradient-to-r transition-all duration-500 group-hover:w-48"
        :class="
          group.pool === 'Mesoneer'
            ? 'from-blue-600 to-indigo-600'
            : 'from-emerald-600 to-teal-600'
        "
      ></div>
    </div>

    <div ref="containerRef" class="relative w-[700px] h-[320px]">
      <!-- SVG Connectors -->
      <svg class="absolute inset-0 w-full h-full pointer-events-none z-0">
        <defs>
          <marker
            id="crosshead"
            markerWidth="10"
            markerHeight="10"
            refX="5"
            refY="5"
            orient="auto"
          >
            <path
              d="M 2 2 L 8 8 M 8 2 L 2 8"
              stroke="#ef4444"
              stroke-width="2"
              stroke-linecap="round"
            />
          </marker>
        </defs>
        <path
          v-for="(path, i) in paths"
          :key="i"
          :d="path.d"
          :stroke="path.stroke || '#cbd5e1'"
          :stroke-width="path.stroke_width || 2"
          fill="none"
          :stroke-opacity="path.opacity || 0.2"
          :marker-end="path.marker"
        />
      </svg>

      <!-- Matches Grid -->
      <div class="absolute inset-0 z-10 grid grid-cols-3 gap-4">
        <!-- Col 1: Opening -->
        <div class="flex flex-col justify-between py-2">
          <MatchCard
            v-if="m1"
            :id="`gsl-${group.id}-m1`"
            :match="m1"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="scale-90"
          />
          <MatchCard
            v-if="m2"
            :id="`gsl-${group.id}-m2`"
            :match="m2"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="scale-90"
          />
        </div>

        <!-- Col 2: Winners/Losers -->
        <div class="flex flex-col justify-between py-2">
          <MatchCard
            v-if="winners"
            :id="`gsl-${group.id}-winners`"
            :match="winners"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="scale-90"
          />
          <MatchCard
            v-if="losers"
            :id="`gsl-${group.id}-losers`"
            :match="losers"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="scale-90"
          />
        </div>

        <!-- Col 3: Decider -->
        <div class="flex flex-col justify-center py-2">
          <MatchCard
            v-if="decider"
            :id="`gsl-${group.id}-decider`"
            :match="decider"
            :is-admin="isAdmin"
            @click="onMatchClick"
            class="scale-90"
          />
        </div>
      </div>

      <!-- ELIMINATION ENDPOINTS (Hidden anchors) -->
      <div
        :id="`elim-m4-${group.id}`"
        class="absolute left-1/2 top-[85%] w-1 h-1"
      ></div>
      <div
        :id="`elim-m5-${group.id}`"
        class="absolute left-full top-1/2 w-1 h-1"
      ></div>
    </div>
  </div>
</template>

<style scoped>
.group-card {
  filter: drop-shadow(0 0 10px rgba(0, 0, 0, 0.1));
}
</style>
