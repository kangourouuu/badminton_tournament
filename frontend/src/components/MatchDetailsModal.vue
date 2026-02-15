<script setup>
import { computed } from "vue";

const props = defineProps({
  isOpen: Boolean,
  match: Object,
});

const emit = defineEmits(["close"]);

// Helper to parse sets detail string or object
const sets = computed(() => {
  if (!props.match?.sets_detail?.sets) return [];
  return props.match.sets_detail.sets;
});

const finished = computed(() => !!props.match?.winner_id);

const teamAClass = computed(() => {
  if (!finished.value) return "text-gray-900";
  return props.match.winner_id === props.match.team_a_id
    ? "font-black text-violet-700"
    : "text-gray-500 opacity-60";
});

const teamBClass = computed(() => {
  if (!finished.value) return "text-gray-900";
  return props.match.winner_id === props.match.team_b_id
    ? "font-black text-violet-700"
    : "text-gray-500 opacity-60";
});
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-[100] flex items-center justify-center bg-gray-900/20 backdrop-blur-sm p-4"
    @click.self="$emit('close')"
  >
    <div
      class="bg-white rounded-lg border border-gray-100 shadow-2xl w-full max-w-md overflow-hidden font-outfit"
    >
      <!-- Header -->
      <div
        class="bg-gray-50 px-6 py-4 flex justify-between items-center border-b border-gray-100"
      >
        <span
          class="text-[10px] font-black uppercase tracking-[0.2em] text-gray-400"
        >
          Match Details
        </span>
        <button
          @click="$emit('close')"
          class="text-gray-400 hover:text-gray-900 transition-colors"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
              clip-rule="evenodd"
            />
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="p-8">
        <!-- Teams & Score Header -->
        <div class="flex justify-between items-center mb-10">
          <!-- Team A -->
          <div class="flex-1 text-center">
            <div class="text-[10px] font-bold text-gray-300 uppercase mb-2">
              Team A
            </div>
            <div class="text-lg leading-tight" :class="teamAClass">
              {{ match.team_a?.name || "TBD" }}
            </div>
          </div>

          <!-- Score -->
          <div class="px-6 flex flex-col items-center">
            <div class="text-3xl font-black text-gray-900 tracking-tighter">
              {{ match.score || "vs" }}
            </div>
          </div>

          <!-- Team B -->
          <div class="flex-1 text-center">
            <div class="text-[10px] font-bold text-gray-300 uppercase mb-2">
              Team B
            </div>
            <div class="text-lg leading-tight" :class="teamBClass">
              {{ match.team_b?.name || "TBD" }}
            </div>
          </div>
        </div>

        <!-- Sets Detail -->
        <div v-if="sets.length > 0" class="space-y-4 mb-8">
          <div
            class="text-[10px] font-bold text-gray-300 uppercase tracking-widest text-center mb-4"
          >
            Set History
          </div>
          <div class="flex justify-center gap-2">
            <div
              v-for="(s, i) in sets"
              :key="i"
              class="flex flex-col items-center bg-gray-50 rounded px-3 py-2 border border-gray-100 min-w-[60px]"
            >
              <span class="text-[8px] font-bold text-gray-400 mb-1"
                >SET {{ s.set || i + 1 }}</span
              >
              <div class="text-sm font-bold text-gray-700">
                {{ s.a }}-{{ s.b }}
              </div>
            </div>
          </div>
        </div>

        <!-- Video Button -->
        <div
          v-if="match.video_url"
          class="text-center pt-4 border-t border-gray-50"
        >
          <a
            :href="match.video_url"
            target="_blank"
            class="inline-flex items-center gap-2 text-violet-600 hover:text-violet-800 transition-colors font-bold text-xs uppercase tracking-widest group"
          >
            <span>Watch Recording</span>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 transform group-hover:translate-x-1 transition-transform"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
              />
            </svg>
          </a>
        </div>
        <div v-else class="text-center pt-4 border-t border-gray-50 opacity-50">
          <span
            class="text-[10px] font-bold text-gray-300 uppercase tracking-widest"
            >No Recording Available</span
          >
        </div>
      </div>
    </div>
  </div>
</template>
