<script setup>
import { ref, watch } from "vue";

const props = defineProps(["pool", "participants"]);
const emit = defineEmits(["generated"]);

const isSpinning = ref(false);
const generatedTeam = ref(null);

const spin = () => {
  if (isSpinning.value || props.participants.length < 2) return;
  isSpinning.value = true;

  // Simulate spin time
  setTimeout(() => {
    isSpinning.value = false;
    // Basic random logic for UI effect (actual logic in backend)
    emit("generated");
  }, 2000);
};
</script>

<template>
  <div
    class="flex flex-col items-center gap-6 p-8 border border-border rounded-lg bg-white"
  >
    <h3 class="text-xl font-bold text-primary">Team Generator: {{ pool }}</h3>

    <div
      class="relative w-64 h-64 rounded-full border-4 border-primary flex items-center justify-center overflow-hidden bg-purple-50"
    >
      <div
        :class="[
          'w-full h-full absolute transition-transform duration-[2000ms] ease-out',
          isSpinning ? 'rotate-[1080deg]' : 'rotate-0',
        ]"
        style="
          background: conic-gradient(
            from 0deg,
            #e9d5ff 0deg 90deg,
            #fdfbff 90deg 180deg,
            #e9d5ff 180deg 270deg,
            #fdfbff 270deg 360deg
          );
        "
      ></div>
      <div
        class="z-10 bg-white p-4 rounded-full shadow-sm font-bold text-primary"
      >
        {{ isSpinning ? "SPINNING..." : "READY" }}
      </div>
    </div>

    <button
      @click="spin"
      :disabled="isSpinning || participants.length < 2"
      class="px-6 py-2 bg-primary text-white rounded-sm font-medium hover:bg-purple-700 disabled:opacity-50 transition-colors"
    >
      {{ isSpinning ? "Spinning..." : "Spin for New Team" }}
    </button>

    <div v-if="participants.length < 2" class="text-red-500 text-sm">
      Not enough participants to form a team.
    </div>
  </div>
</template>
