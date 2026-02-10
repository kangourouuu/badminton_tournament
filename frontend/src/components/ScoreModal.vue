<script setup>
import { ref, watch, computed } from "vue";

const props = defineProps({
  isOpen: Boolean,
  match: Object,
  stage: {
    type: String,
    default: "KNOCKOUT", // 'GROUP' or 'KNOCKOUT'
  },
});

const emit = defineEmits(["close", "save"]);

const winnerId = ref("");
const videoUrl = ref("");

// Detailed Sets: [ {a: null, b: null}, ... ]
const sets = ref([
  { a: null, b: null },
  { a: null, b: null },
  { a: null, b: null },
]);

// simple scoring condition
const isSimpleScoring = computed(() => {
  if (props.stage !== "GROUP") return false;
  const label = props.match?.label || "";
  return label === "M1" || label === "M2";
});

// Auto-calculate score string (e.g. "2-1")
const calculatedScore = computed(() => {
  if (isSimpleScoring.value) {
    return winnerId.value ? "1-0" : "";
  }
  let aWins = 0;
  let bWins = 0;
  sets.value.forEach((s) => {
    if (s.a !== null && s.b !== null && s.a !== "" && s.b !== "") {
      const pA = parseInt(s.a);
      const pB = parseInt(s.b);
      if (pA > pB) aWins++;
      else if (pB > pA) bWins++;
    }
  });
  if (aWins === 0 && bWins === 0) return "";
  return `${aWins}-${bWins}`;
});

// Sync data when match changes
watch(
  () => props.match,
  (newMatch) => {
    if (newMatch) {
      winnerId.value = newMatch.winner_id || "";
      videoUrl.value = newMatch.video_url || "";

      // Parse detailed sets if available
      if (newMatch.sets_detail && newMatch.sets_detail.sets) {
        // Map back to our structure
        const loadedSets = newMatch.sets_detail.sets;
        sets.value = [
          { a: loadedSets[0]?.a ?? null, b: loadedSets[0]?.b ?? null },
          { a: loadedSets[1]?.a ?? null, b: loadedSets[1]?.b ?? null },
          { a: loadedSets[2]?.a ?? null, b: loadedSets[2]?.b ?? null },
        ];
      } else if (newMatch.score) {
        // Legacy/Simple score support could go here, but let's reset to clean slate
        // or try to parse if format was "21-19, 15-21"
        // For now, reset to empty to encourage re-entry
        sets.value = [
          { a: null, b: null },
          { a: null, b: null },
          { a: null, b: null },
        ];
      } else {
        sets.value = [
          { a: null, b: null },
          { a: null, b: null },
          { a: null, b: null },
        ];
      }
    }
  },
  { immediate: true },
);

const save = () => {
  // Construct sets_detail payload
  let validSets = [];
  let score = calculatedScore.value || "0-0";

  if (isSimpleScoring.value) {
    // For Group stage, we enforce a simple "1-0" structure invisibly if needed,
    // or just send empty sets. Backend should handle simplified score string.
    // We can mock a single set 1-0 for consistency?
    // Master plan says just "immediately sets score='1-0'".
    score = "1-0";
    validSets = []; // No detailed sets for simplified group mode
  } else {
    validSets = sets.value
      .map((s, i) => ({
        set: i + 1,
        a: s.a && parseInt(s.a),
        b: s.b && parseInt(s.b),
      }))
      .filter((s) => s.a !== null && !isNaN(s.a));
  }

  emit("save", {
    id: props.match.id,
    winner_id: winnerId.value,
    score: score,
    sets_detail: { sets: validSets },
    video_url: videoUrl.value,
  });
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm transition-opacity"
  >
    <div
      class="bg-white rounded-xl shadow-2xl w-full max-w-lg overflow-hidden transform transition-all font-outfit"
    >
      <!-- Header -->
      <div class="bg-slate-900 px-6 py-4 flex justify-between items-center">
        <h3 class="text-lg font-bold text-white uppercase tracking-wider">
          Update Match Result
        </h3>
        <button
          @click="$emit('close')"
          class="text-slate-400 hover:text-white transition-colors"
        >
          ✕
        </button>
      </div>

      <div class="p-8 space-y-8">
        <!-- Teams Header -->
        <div class="flex justify-between items-center px-4">
          <div class="text-center w-1/3">
            <div class="text-sm text-slate-400 uppercase font-bold mb-1">
              Team A
            </div>
            <div class="font-black text-slate-800 text-lg leading-tight">
              {{ match.team_a?.name || "TBD" }}
            </div>
          </div>
          <div class="text-2xl font-black text-slate-300">VS</div>
          <div class="text-center w-1/3">
            <div class="text-sm text-slate-400 uppercase font-bold mb-1">
              Team B
            </div>
            <div class="font-black text-slate-800 text-lg leading-tight">
              {{ match.team_b?.name || "TBD" }}
            </div>
          </div>
        </div>

        <!-- FEATURE A: SIMPLIFIED GROUP SCORING (Only M1/M2) -->
        <div v-if="isSimpleScoring">
          <label
            class="block text-xs font-bold text-slate-400 uppercase mb-4 text-center"
          >
            Select Winner (Opening Match)
          </label>
          <div class="grid grid-cols-2 gap-4">
            <button
              v-if="match.team_a"
              class="p-6 border-2 rounded-lg font-outfit font-bold text-lg transition-all flex flex-col items-center gap-2 group relative overflow-hidden"
              :class="
                winnerId === match.team_a.id
                  ? 'border-violet-600 bg-purple-50 text-violet-900 shadow-md ring-2 ring-violet-200'
                  : 'border-gray-200 text-gray-600 hover:border-gray-300 hover:bg-gray-50'
              "
              @click="winnerId = match.team_a.id"
            >
              <span class="text-2xl transition-transform group-hover:scale-110"
                >🏆</span
              >
              <span>{{ match.team_a.name }}</span>
              <div
                v-if="winnerId === match.team_a.id"
                class="absolute top-2 right-2 text-violet-600"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fill-rule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                    clip-rule="evenodd"
                  />
                </svg>
              </div>
            </button>

            <button
              v-if="match.team_b"
              class="p-6 border-2 rounded-lg font-outfit font-bold text-lg transition-all flex flex-col items-center gap-2 group relative overflow-hidden"
              :class="
                winnerId === match.team_b.id
                  ? 'border-violet-600 bg-purple-50 text-violet-900 shadow-md ring-2 ring-violet-200'
                  : 'border-gray-200 text-gray-600 hover:border-gray-300 hover:bg-gray-50'
              "
              @click="winnerId = match.team_b.id"
            >
              <span class="text-2xl transition-transform group-hover:scale-110"
                >🏆</span
              >
              <span>{{ match.team_b.name }}</span>
              <div
                v-if="winnerId === match.team_b.id"
                class="absolute top-2 right-2 text-violet-600"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fill-rule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                    clip-rule="evenodd"
                  />
                </svg>
              </div>
            </button>
          </div>
          <p class="text-center text-xs text-gray-400 mt-4 italic">
            Group match outcomes are effectively "Win/Loss". Detailed scores are
            optional for internal records but simplified here for speed.
          </p>
        </div>

        <!-- FEATURE A: KNOCKOUT DETAILED SCORING -->
        <div v-else>
          <!-- Winner Selection Dropdown (Legacy/Backup for Knockout) -->
          <div class="relative mb-8">
            <label class="block text-xs font-bold text-slate-400 uppercase mb-2"
              >Match Winner</label
            >
            <select
              v-model="winnerId"
              class="input-material font-bold text-slate-800"
            >
              <option value="" disabled>Select Winner...</option>
              <option v-if="match.team_a" :value="match.team_a.id">
                {{ match.team_a.name }}
              </option>
              <option v-if="match.team_b" :value="match.team_b.id">
                {{ match.team_b.name }}
              </option>
            </select>
          </div>

          <!-- Detailed Scoring -->
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase mb-4"
              >Set Scores (Best of 3)</label
            >
            <div class="space-y-4 font-outfit">
              <div
                v-for="(set, i) in sets"
                :key="i"
                class="flex items-center gap-4"
              >
                <div
                  class="w-16 text-xs font-bold text-slate-400 uppercase pt-2"
                >
                  Set {{ i + 1 }}
                </div>
                <input
                  v-model="set.a"
                  type="number"
                  placeholder="0"
                  class="input-material text-center font-bold text-lg text-slate-700 focus:text-violet-700"
                />
                <span class="text-slate-300 font-bold">-</span>
                <input
                  v-model="set.b"
                  type="number"
                  placeholder="0"
                  class="input-material text-center font-bold text-lg text-slate-700 focus:text-violet-700"
                />
              </div>
            </div>

            <!-- Calculated Score Display -->
            <div class="mt-4 text-center">
              <span class="text-xs font-bold text-slate-400 uppercase mr-2"
                >Calculated Score:</span
              >
              <span class="text-xl font-black text-violet-600">{{
                calculatedScore || "0-0"
              }}</span>
            </div>
          </div>
        </div>

        <!-- Video URL (Common) -->
        <div>
          <label class="block text-xs font-bold text-slate-400 uppercase mb-2"
            >Video Recording URL</label
          >
          <input
            v-model="videoUrl"
            type="text"
            placeholder="https://youtu.be/..."
            class="input-material text-sm"
          />
        </div>
      </div>

      <!-- Footer -->
      <div
        class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end space-x-4"
      >
        <button
          @click="$emit('close')"
          class="px-6 py-2 text-slate-500 font-bold hover:text-slate-700 transition-colors"
        >
          Cancel
        </button>
        <button
          @click="save"
          class="px-8 py-2 bg-gradient-to-r from-violet-600 to-indigo-600 text-white font-bold rounded-full shadow-lg hover:shadow-xl hover:scale-105 transition-all"
        >
          Save Details
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Scoped overrides if needed, mostly Tailwind */
input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
</style>
