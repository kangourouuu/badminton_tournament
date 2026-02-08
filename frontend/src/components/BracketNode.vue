<script setup>
const props = defineProps({
  match: Object, // Match object
  admin: Boolean, // Is admin mode?
});
const emit = defineEmits(["open-video", "update-result"]);

// Helper to get winner class
const getWinnerClass = (teamId) => {
  if (!props.match.winner_id) return "";
  return props.match.winner_id === teamId
    ? "font-bold text-primary"
    : "text-gray-400 opacity-70";
};
</script>

<template>
  <div
    class="flex flex-col gap-1 w-64 border border-border bg-white rounded-sm p-2 relative"
  >
    <!-- Header -->
    <div class="flex justify-between text-xs text-gray-500 mb-1">
      <span>{{ match.description || "Match " + match.id }}</span>
      <button
        v-if="match.video_url"
        @click="$emit('open-video', match.video_url)"
        class="hover:text-primary"
      >
        ▶ Watch Replay
      </button>
    </div>

    <!-- Team A -->
    <div
      class="flex justify-between items-center p-1 rounded-sm bg-purple-50"
      :class="getWinnerClass(match.team_a_id)"
    >
      <span class="truncate">{{ match.team_a?.name || "TBD" }}</span>
      <span class="font-mono">{{ match.score_a }}</span>
    </div>

    <!-- Team B -->
    <div
      class="flex justify-between items-center p-1 rounded-sm bg-purple-50 mt-1"
      :class="getWinnerClass(match.team_b_id)"
    >
      <span class="truncate">{{ match.team_b?.name || "TBD" }}</span>
      <span class="font-mono">{{ match.score_b }}</span>
    </div>

    <!-- Score Details -->
    <div
      v-if="match.score_details"
      class="px-1 mt-1 text-[10px] text-gray-400 text-center"
    >
      {{ match.score_details }}
    </div>

    <!-- Admin Actions -->
    <button
      v-if="admin && match.team_a_id && match.team_b_id && !match.winner_id"
      @click="$emit('update-result', match)"
      class="absolute -right-2 -top-2 bg-primary text-white w-5 h-5 rounded-full flex items-center justify-center text-xs hover:bg-purple-700"
    >
      ✎
    </button>
  </div>
</template>
