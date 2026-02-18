<script setup>
import { computed } from "vue";

const props = defineProps({
  match: {
    type: Object,
    required: true,
  },
  nextMatchLabel: {
    type: String,
    default: "",
  },
});

defineEmits(["click"]);

const isFinished = computed(() => props.match.winner_id);

const teamAClass = computed(() => {
  if (!isFinished.value) return "text-gray-700";
  return props.match.winner_id === props.match.team_a_id
    ? "font-black text-violet-700"
    : "text-gray-400";
});

const teamBClass = computed(() => {
  if (!isFinished.value) return "text-gray-700";
  return props.match.winner_id === props.match.team_b_id
    ? "font-black text-violet-700"
    : "text-gray-400";
});
</script>

<template>
  <div
    class="w-64 bg-white border border-gray-100 rounded-[4px] p-4 cursor-pointer hover:border-violet-200 transition-all active:scale-95 group font-outfit"
    :id="'match-' + match.id"
    @click="$emit('click', match)"
  >
    <!-- Teams -->
    <div class="space-y-3">
      <div class="flex justify-between items-center">
        <span
          class="text-sm tracking-tight truncate flex-1"
          :class="teamAClass"
        >
          {{ match.team_a?.name || "TBD" }}
        </span>
        <span
          v-if="isFinished && match.winner_id === match.team_a_id"
          class="text-[10px] bg-violet-100 text-violet-700 px-1.5 py-0.5 rounded font-bold uppercase ml-2"
          >Win</span
        >
      </div>

      <div
        class="h-[1px] bg-gray-50 group-hover:bg-violet-50 transition-colors"
      ></div>

      <div class="flex justify-between items-center">
        <span
          class="text-sm tracking-tight truncate flex-1"
          :class="teamBClass"
        >
          {{ match.team_b?.name || "TBD" }}
        </span>
        <span
          v-if="isFinished && match.winner_id === match.team_b_id"
          class="text-[10px] bg-violet-100 text-violet-700 px-1.5 py-0.5 rounded font-bold uppercase ml-2"
          >Win</span
        >
      </div>
    </div>

    <!-- Footer metadata (hidden if not finished?) -->
    <div
      class="mt-4 pt-3 border-t border-gray-50 flex justify-between items-center"
    >
      <div v-if="isFinished">
        <span
          class="text-[10px] font-bold text-violet-600 uppercase tracking-widest"
          >{{ match.score }}</span
        >
      </div>
      <div
        v-else-if="nextMatchLabel"
        class="text-[9px] font-bold text-gray-400 uppercase tracking-wider flex items-center gap-1"
      >
        <span>Winner to</span>
        <span class="bg-gray-100 text-gray-600 px-1 rounded">{{
          nextMatchLabel
        }}</span>
      </div>
      <div v-else></div>

      <div v-if="match.video_url" class="text-violet-400 hover:text-violet-600">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
          />
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
      </div>
    </div>
  </div>
</template>
