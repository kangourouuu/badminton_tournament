<script setup>
import { ref, watch, computed } from "vue";

const props = defineProps({
  isOpen: Boolean,
  match: Object,
});

const emit = defineEmits(["close", "save"]);

const winnerId = ref("");
const score = ref("");
const videoUrl = ref("");
const knockoutSets = ref([""]); // Array for Best of 3 sets

// Determine if this is a Group Stage match (M1-M5) or Knockout (SF/Final)
const isGroupStage = computed(() => {
  if (!props.match) return true;
  // Matches with strictly M1-M5 labels are Group Stage
  return [
    "M1",
    "M2",
    "M3",
    "Winners",
    "M4",
    "Losers",
    "M5",
    "Decider",
  ].includes(props.match.label);
});

// Sync data when match changes
watch(
  () => props.match,
  (newMatch) => {
    if (newMatch) {
      winnerId.value = newMatch.winner_id || "";
      videoUrl.value = newMatch.video_url || "";

      // Parse Score
      if (isGroupStage.value) {
        score.value = newMatch.score || "1-0";
      } else {
        // Flattened string "21-19, 21-18" -> Array
        knockoutSets.value = newMatch.score ? newMatch.score.split(", ") : [""];
      }
    }
  },
  { immediate: true },
);

const addSet = () => {
  if (knockoutSets.value.length < 3) {
    knockoutSets.value.push("");
  }
};

const save = () => {
  let finalScore = "";

  if (isGroupStage.value) {
    finalScore = score.value || "1-0";
  } else {
    // Filter empty sets and join
    finalScore = knockoutSets.value.filter((s) => s.trim() !== "").join(", ");
  }

  emit("save", {
    id: props.match.id,
    winner_id: winnerId.value,
    score: finalScore,
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
            >Winner <span class="text-red-500">*</span></label
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

        <!-- Group Stage Logic (1 Set) -->
        <div
          v-if="isGroupStage"
          class="p-4 bg-gray-50 rounded-sm border border-gray-200"
        >
          <label class="block text-xs font-bold text-gray-500 uppercase mb-2"
            >Group Stage Scoring (1 Set)</label
          >
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-gray-700">Score:</span>
            <input
              v-model="score"
              type="text"
              placeholder="1-0"
              class="flex-1 border border-gray-300 rounded-sm p-2 text-sm focus:ring-2 focus:ring-violet-500"
            />
          </div>
          <p class="text-xs text-gray-400 mt-1">
            Default is "1-0". You can enter points (e.g. "31-29") if needed.
          </p>
        </div>

        <!-- Knockout Stage Logic (Best of 3) -->
        <div
          v-else
          class="p-4 bg-purple-50 rounded-sm border border-purple-100"
        >
          <label class="block text-xs font-bold text-purple-600 uppercase mb-2"
            >Knockout Scoring (Best of 3)</label
          >
          <div class="space-y-2">
            <div
              class="flex gap-2"
              v-for="(set, index) in knockoutSets"
              :key="index"
            >
              <span class="text-xs text-gray-500 w-12 pt-2"
                >Set {{ index + 1 }}</span
              >
              <input
                v-model="knockoutSets[index]"
                type="text"
                placeholder="21-19"
                class="flex-1 border border-gray-300 rounded-sm p-1 text-sm text-center"
              />
            </div>
            <button
              v-if="knockoutSets.length < 3"
              @click="addSet"
              class="text-xs text-violet-600 hover:text-violet-800 font-medium underline"
            >
              + Add Set
            </button>
          </div>
        </div>

        <!-- Video URL -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Match's Video URL</label
          >
          <input
            v-model="videoUrl"
            type="text"
            placeholder="https://youtu.be/..."
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
