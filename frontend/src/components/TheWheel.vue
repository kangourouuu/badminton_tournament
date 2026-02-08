<script setup>
import { ref } from "vue";
import api from "../services/api";

const props = defineProps({
  pool: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["teams-generated"]);

const isSpinning = ref(false);
const generatedTeams = ref([]);
const error = ref("");

const generate = async () => {
  isSpinning.value = true;
  error.value = "";
  generatedTeams.value = [];

  // Simulate spin time
  await new Promise((resolve) => setTimeout(resolve, 2000));

  try {
    const response = await api.post("/teams/generate", { pool: props.pool });
    generatedTeams.value = response.data.teams;
    emit("teams-generated");
  } catch (err) {
    error.value =
      "Failed to generate teams: " + (err.response?.data?.error || err.message);
  } finally {
    isSpinning.value = false;
  }
};
</script>

<template>
  <div
    class="flex flex-col items-center justify-center space-y-8 p-8 bg-white rounded-lg border border-purple-50 shadow-sm"
  >
    <div class="relative">
      <!-- Wheel Visual -->
      <div
        :class="{ 'animate-spin': isSpinning }"
        class="w-48 h-48 rounded-full border-8 border-violet-200 border-t-violet-600 flex items-center justify-center transition-all duration-1000 ease-in-out"
      >
        <span v-if="!isSpinning" class="text-4xl">🏸</span>
        <span v-else class="text-4xl">🎲</span>
      </div>
    </div>

    <div class="space-y-4 text-center w-full">
      <button
        @click="generate"
        :disabled="isSpinning"
        class="px-8 py-3 bg-violet-600 text-white font-bold rounded-sm hover:bg-violet-700 disabled:opacity-50 transition-colors shadow-sm"
      >
        {{ isSpinning ? "Randomizing..." : "Generate " + pool + " Teams" }}
      </button>

      <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>

      <!-- Result List -->
      <div
        v-if="generatedTeams.length > 0"
        class="mt-8 text-left w-full max-w-md mx-auto"
      >
        <h3 class="text-lg font-bold text-gray-800 mb-4 border-b pb-2">
          Generated Teams
        </h3>
        <ul class="space-y-2">
          <li
            v-for="team in generatedTeams"
            :key="team.id"
            class="flex justify-between p-3 bg-gray-50 rounded-sm border border-gray-100"
          >
            <span class="font-medium text-purple-900">{{ team.name }}</span>
            <span class="text-xs text-gray-400 font-mono">{{ team.pool }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(1800deg);
  }
}
.animate-spin {
  animation: spin 2s cubic-bezier(0.25, 1, 0.5, 1) infinite;
}
</style>
