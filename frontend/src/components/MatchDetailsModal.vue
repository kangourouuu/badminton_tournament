<script setup>
import { computed, watch, ref } from "vue";
import api from "../services/api";

const props = defineProps({
  isOpen: Boolean,
  match: Object, // Can be partial, used for ID
});

const emit = defineEmits(["close"]);

const matchData = ref(null);
const loading = ref(false);

const fetchData = async () => {
  if (!props.match?.id) return;
  loading.value = true;
  try {
    // GET /matches/:id
    const res = await api.get(`/matches/${props.match.id}`);
    matchData.value = res.data;
  } catch (err) {
    console.error("Failed to fetch match details", err);
  } finally {
    loading.value = false;
  }
};

watch(
  () => props.isOpen,
  (newVal) => {
    if (newVal) {
      matchData.value = null; // Reset
      fetchData();
    }
  },
);

// Helper to parse sets detail string or object
const sets = computed(() => {
  if (!matchData.value) return [];
  const details = matchData.value.sets_detail;
  if (!details) return [];

  // If it's the JSONB object from DB: { sets: [...] }
  if (typeof details === "object" && details.sets) {
    return details.sets;
  }

  // If string (legacy/fallback) "21-19, 15-21"
  if (typeof details === "string") {
    return details.split(",").map((s, i) => {
      const parts = s.trim().split("-");
      return { set: i + 1, a: parts[0], b: parts[1] };
    });
  }
  return [];
});

const finished = computed(() => !!matchData.value?.winner_id);

const teamAClass = computed(() => {
  if (!matchData.value) return "text-gray-900";
  if (!finished.value) return "text-gray-900";
  return matchData.value.winner_id === matchData.value.team_a_id
    ? "font-black text-violet-700"
    : "text-gray-500 opacity-60";
});

const teamBClass = computed(() => {
  if (!matchData.value) return "text-gray-900";
  if (!finished.value) return "text-gray-900";
  return matchData.value.winner_id === matchData.value.team_b_id
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
      class="bg-white rounded-lg border border-gray-100 shadow-2xl w-full max-w-md overflow-hidden font-outfit min-h-[300px] flex flex-col"
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

      <!-- Loading State -->
      <div
        v-if="loading || !matchData"
        class="flex-1 flex items-center justify-center p-10"
      >
        <div
          class="animate-spin rounded-full h-8 w-8 border-b-2 border-violet-600"
        ></div>
      </div>

      <!-- Content -->
      <div v-else class="p-8">
        <!-- Teams & Score Header -->
        <div class="flex justify-between items-center mb-10">
          <!-- Team A -->
          <div class="flex-1 text-center">
            <div class="text-[10px] font-bold text-gray-300 uppercase mb-2">
              Team A
            </div>
            <div class="text-lg leading-tight" :class="teamAClass">
              {{ matchData.team_a?.name || "TBD" }}
            </div>
          </div>

          <!-- Score -->
          <div class="px-6 flex flex-col items-center">
            <div class="text-3xl font-black text-gray-900 tracking-tighter">
              {{ matchData.score || "vs" }}
            </div>
          </div>

          <!-- Team B -->
          <div class="flex-1 text-center">
            <div class="text-[10px] font-bold text-gray-300 uppercase mb-2">
              Team B
            </div>
            <div class="text-lg leading-tight" :class="teamBClass">
              {{ matchData.team_b?.name || "TBD" }}
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
          <div class="flex flex-col gap-2 items-center">
            <div
              v-for="(s, i) in sets"
              :key="i"
              class="flex justify-between items-center bg-gray-50 rounded px-4 py-2 border border-gray-100 w-full max-w-[200px]"
            >
              <span
                class="text-[10px] font-black text-gray-400 uppercase tracking-wider"
                >SET {{ s.set || i + 1 }}</span
              >
              <div class="text-sm font-bold text-gray-700 font-mono">
                {{ s.a }} - {{ s.b }}
              </div>
            </div>
          </div>
        </div>

        <!-- Video Button -->
        <div
          v-if="matchData.video_url"
          class="text-center pt-4 border-t border-gray-50"
        >
          <a
            :href="matchData.video_url"
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
