<script setup>
import { ref, watch } from "vue";

const props = defineProps({
  isOpen: Boolean,
  match: Object,
});

const emit = defineEmits(["close", "save"]);

const winnerId = ref("");
const score = ref("");
const videoUrl = ref("");

// Sync data when match changes
watch(
  () => props.match,
  (newMatch) => {
    if (newMatch) {
      winnerId.value = newMatch.winner_id || "";
      score.value = newMatch.score || "";
      videoUrl.value = newMatch.video_url || "";
    }
  },
  { immediate: true },
);

const save = () => {
  emit("save", {
    id: props.match.id,
    winner_id: winnerId.value,
    score: score.value,
    video_url: videoUrl.value,
  });
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
  >
    <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6 space-y-6">
      <h3 class="text-xl font-bold text-gray-800">Update Match Result</h3>

      <div class="space-y-4">
        <!-- Winner Selection -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Winner</label
          >
          <select
            v-model="winnerId"
            class="w-full border border-gray-300 rounded-sm p-2 focus:ring-2 focus:ring-violet-500 focus:border-violet-500"
          >
            <option value="">Select Winner...</option>
            <option v-if="match.team_a" :value="match.team_a.id">
              {{ match.team_a.name }}
            </option>
            <option v-if="match.team_b" :value="match.team_b.id">
              {{ match.team_b.name }}
            </option>
          </select>
        </div>

        <!-- Score -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Score (e.g. 21-19, 21-18)</label
          >
          <input
            v-model="score"
            type="text"
            class="w-full border border-gray-300 rounded-sm p-2 focus:ring-2 focus:ring-violet-500 focus:border-violet-500"
          />
        </div>

        <!-- Video URL -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >YouTube URL</label
          >
          <input
            v-model="videoUrl"
            type="text"
            class="w-full border border-gray-300 rounded-sm p-2 focus:ring-2 focus:ring-violet-500 focus:border-violet-500"
          />
        </div>
      </div>

      <div class="flex justify-end space-x-3 pt-4 border-t border-gray-100">
        <button
          @click="$emit('close')"
          class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-sm transition-colors"
        >
          Cancel
        </button>
        <button
          @click="save"
          class="px-4 py-2 bg-violet-600 text-white rounded-sm hover:bg-violet-700 transition-colors"
        >
          Save Result
        </button>
      </div>
    </div>
  </div>
</template>
