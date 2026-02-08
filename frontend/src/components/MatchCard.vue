<script setup>
import { computed } from "vue";

const props = defineProps({
  match: {
    type: Object,
    required: true,
  },
  isAdmin: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["click", "play-video"]);

// Determine status
const status = computed(() => {
  if (props.match.winner_id) return "finished";
  if (props.match.team_a_id && props.match.team_b_id) return "ready"; // or live if we had that state
  return "scheduled";
});

const cardClasses = computed(() => {
  const base =
    "bg-white rounded-sm p-3 border-2 transition-all duration-200 relative overflow-hidden";
  if (status.value === "finished") return `${base} border-gray-200 opacity-90`;
  if (status.value === "ready") return `${base} border-purple-500 shadow-sm`; // "Live" look
  return `${base} border-gray-100 text-gray-400 border-dashed`; // Empty/Scheduled
});

const getTeamName = (teamId) => {
  if (!teamId) return "TBD";
  if (props.match.team_a?.id === teamId) return props.match.team_a.name;
  if (props.match.team_b?.id === teamId) return props.match.team_b.name;
  // Fallback if relation not loaded but ID exists (shouldn't happen with eager load)
  return "Team " + teamId.substring(0, 4);
};

const isWinner = (teamId) => {
  return status.value === "finished" && props.match.winner_id === teamId;
};

const handleClick = () => {
  if (props.isAdmin) {
    emit("click", props.match);
  } else if (props.match.video_url) {
    window.open(props.match.video_url, "_blank"); // Simple video handler for now
  }
};
</script>

<template>
  <div
    :class="cardClasses"
    @click="handleClick"
    class="cursor-pointer hover:shadow-md group"
  >
    <!-- Label -->
    <div
      class="absolute top-0 right-0 bg-gray-100 text-xs text-gray-500 px-2 py-0.5 rounded-bl-sm font-mono"
    >
      {{ match.label }}
    </div>

    <div class="space-y-2 mt-2">
      <!-- Team A -->
      <div class="flex justify-between items-center">
        <span
          :class="{
            'font-bold text-gray-900': true,
            'text-green-600': isWinner(match.team_a_id),
          }"
        >
          {{ match.team_a?.name || "Waiting..." }}
        </span>
        <span
          v-if="match.score && isWinner(match.team_a_id)"
          class="text-xs font-bold bg-green-100 text-green-700 px-1.5 py-0.5 rounded"
          >WIN</span
        >
      </div>

      <!-- VS / Score -->
      <div
        class="text-xs text-center text-gray-400 font-medium py-1 border-t border-b border-gray-50 my-1"
      >
        <span v-if="match.score" class="text-purple-600 font-bold text-sm">{{
          match.score
        }}</span>
        <span v-else>vs</span>
      </div>

      <!-- Team B -->
      <div class="flex justify-between items-center">
        <span
          :class="{
            'font-bold text-gray-900': true,
            'text-green-600': isWinner(match.team_b_id),
          }"
        >
          {{ match.team_b?.name || "Waiting..." }}
        </span>
        <span
          v-if="match.score && isWinner(match.team_b_id)"
          class="text-xs font-bold bg-green-100 text-green-700 px-1.5 py-0.5 rounded"
          >WIN</span
        >
      </div>
    </div>

    <!-- Video Indicator -->
    <div
      v-if="match.video_url && !isAdmin"
      class="absolute inset-0 bg-black/5 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
    >
      <div class="bg-white rounded-full p-2 shadow-sm text-red-600">▶</div>
    </div>
  </div>
</template>
