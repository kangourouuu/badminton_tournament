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

// simple scoring condition (Only Opening Match or Winners/Decider in Group Stage)
const isGroupStage = computed(() => props.stage === "GROUP");

// Score inputs for Group Stage (Simplified)
const groupScoreA = ref("");
const groupScoreB = ref("");

// Auto-calculate score string (e.g. "2-1")
const calculatedScore = computed(() => {
  if (isGroupStage.value) {
    if (!groupScoreA.value || !groupScoreB.value) return "";
    return parseInt(groupScoreA.value) > parseInt(groupScoreB.value)
      ? "1-0"
      : "0-1";
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

// Auto-detect winner for Group Stage
watch([groupScoreA, groupScoreB], () => {
  if (isGroupStage.value && groupScoreA.value && groupScoreB.value) {
    const sA = parseInt(groupScoreA.value);
    const sB = parseInt(groupScoreB.value);
    if (sA > sB) winnerId.value = props.match.team_a?.id;
    else if (sB > sA) winnerId.value = props.match.team_b?.id;
  }
});

// Sync data when match changes
watch(
  () => props.match,
  (newMatch) => {
    if (newMatch) {
      winnerId.value = newMatch.winner_id || "";
      videoUrl.value = newMatch.video_url || "";

      // Reset Group scores
      groupScoreA.value = "";
      groupScoreB.value = "";

      // Parse score if exists (e.g. "1-0")
      if (newMatch.score && isGroupStage.value) {
        // Just reset inputs to clean slate for new entry
      }

      // Parse detailed sets if available
      if (newMatch.sets_detail && newMatch.sets_detail.sets) {
        const loadedSets = newMatch.sets_detail.sets;
        sets.value = [
          { a: loadedSets[0]?.a ?? null, b: loadedSets[0]?.b ?? null },
          { a: loadedSets[1]?.a ?? null, b: loadedSets[1]?.b ?? null },
          { a: loadedSets[2]?.a ?? null, b: loadedSets[2]?.b ?? null },
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
  let validSets = [];
  let score = calculatedScore.value || "";

  if (isGroupStage.value) {
    // Backend strictly expects some notation, but Masterplan says always '1-0' effectively
    if (!winnerId.value)
      return alert("Please enter scores to determine winner");
    score = calculatedScore.value;
    validSets = [];
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

        <!-- FEATURE A: SIMPLIFIED GROUP SCORING (Winner detection) -->
        <div v-if="isGroupStage" class="space-y-6">
          <div
            class="flex justify-center items-center gap-8 px-4 py-8 bg-slate-50 rounded-xl border border-slate-100"
          >
            <div class="flex flex-col items-center gap-3">
              <input
                v-model="groupScoreA"
                type="number"
                placeholder="21"
                class="w-24 h-24 bg-white border-4 border-slate-200 rounded-2xl text-center text-4xl font-black text-slate-800 focus:border-violet-500 focus:ring-4 focus:ring-violet-500/20 outline-none transition-all shadow-sm"
              />
              <span class="text-[10px] font-bold text-slate-400 uppercase">{{
                match.team_a?.name
              }}</span>
              <span
                v-if="winnerId === match.team_a?.id"
                class="text-xs font-black text-violet-600 bg-violet-100 px-2 py-0.5 rounded animate-bounce"
                >WINNER</span
              >
            </div>

            <div class="text-2xl font-black text-slate-300">VS</div>

            <div class="flex flex-col items-center gap-3">
              <input
                v-model="groupScoreB"
                type="number"
                placeholder="19"
                class="w-24 h-24 bg-white border-4 border-slate-200 rounded-2xl text-center text-4xl font-black text-slate-800 focus:border-violet-500 focus:ring-4 focus:ring-violet-500/20 outline-none transition-all shadow-sm"
              />
              <span class="text-[10px] font-bold text-slate-400 uppercase">{{
                match.team_b?.name
              }}</span>
              <span
                v-if="winnerId === match.team_b?.id"
                class="text-xs font-black text-violet-600 bg-violet-100 px-2 py-0.5 rounded animate-bounce"
                >WINNER</span
              >
            </div>
          </div>

          <p
            class="text-center text-[10px] text-slate-400 uppercase tracking-widest font-bold"
          >
            System will automatically detect winner based on score
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
