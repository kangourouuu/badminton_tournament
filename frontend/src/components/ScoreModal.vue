<script setup>
import { ref, watch, computed } from "vue";

const props = defineProps({
  isOpen: Boolean,
  match: Object,
  stage: { type: String, default: "KNOCKOUT" }, // 'GROUP' or 'KNOCKOUT'
  isAdmin: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["close", "save"]);

const winnerId = ref("");
const videoUrl = ref("");

// Stage detection
const isGroupStage = computed(() => {
  return props.stage === "GROUP";
});

// Detailed Sets: [ {a: null, b: null}, ... ] - Used for Knockout (BO3)
const knockoutSets = ref([
  { a: null, b: null },
  { a: null, b: null },
  { a: null, b: null },
]);

// Group Stage Inputs: 1 Set
const groupSet = ref({ a: null, b: null });

// Auto-calculate score string (e.g. "2-1")
const calculatedSummaryScore = computed(() => {
  if (isGroupStage.value) {
    if (groupSet.value.a === null || groupSet.value.b === null) return "";
    return parseInt(groupSet.value.a) > parseInt(groupSet.value.b)
      ? "1-0"
      : "0-1";
  }

  let aWins = 0;
  let bWins = 0;
  knockoutSets.value.forEach((s) => {
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

// Auto-detect winner based on score
watch(
  [groupSet, knockoutSets, calculatedSummaryScore],
  () => {
    if (isGroupStage.value) {
      if (groupSet.value.a !== null && groupSet.value.b !== null) {
        if (parseInt(groupSet.value.a) > parseInt(groupSet.value.b))
          winnerId.value = props.match.team_a_id;
        else if (parseInt(groupSet.value.b) > parseInt(groupSet.value.a))
          winnerId.value = props.match.team_b_id;
        else winnerId.value = ""; // Tie
      }
    } else {
      // BO3 Logic: Who reached 2 sets first?
      let teamASetsWon = 0;
      let teamBSetsWon = 0;

      const s1 = knockoutSets.value[0];
      const s2 = knockoutSets.value[1];
      const s3 = knockoutSets.value[2];

      // Set 1
      if (s1.a !== null && s1.b !== null) {
        if (parseInt(s1.a) > parseInt(s1.b)) teamASetsWon++;
        else if (parseInt(s1.b) > parseInt(s1.a)) teamBSetsWon++;
      }

      // Set 2
      if (s2.a !== null && s2.b !== null) {
        if (parseInt(s2.a) > parseInt(s2.b)) teamASetsWon++;
        else if (parseInt(s2.b) > parseInt(s2.a)) teamBSetsWon++;
      }

      // Set 3 (Only count if scores are entered)
      if (
        s3.a !== null &&
        s3.b !== null &&
        (parseInt(s3.a) > 0 || parseInt(s3.b) > 0)
      ) {
        if (parseInt(s3.a) > parseInt(s3.b)) teamASetsWon++;
        else if (parseInt(s3.b) > parseInt(s3.a)) teamBSetsWon++;
      }

      if (teamASetsWon === 2) {
        winnerId.value = props.match.team_a_id;
      } else if (teamBSetsWon === 2) {
        winnerId.value = props.match.team_b_id;
      } else {
        winnerId.value = ""; // Incomplete or tie
      }
    }
  },
  { deep: true },
);

// Sync data when match changes
watch(
  () => props.match,
  (newMatch) => {
    if (newMatch) {
      winnerId.value = newMatch.winner_id || "";
      videoUrl.value = newMatch.video_url || "";

      // Reset
      groupSet.value = { a: null, b: null };
      knockoutSets.value = [
        { a: null, b: null },
        { a: null, b: null },
        { a: null, b: null },
      ];

      // Parse existing sets if they exist
      if (newMatch.sets_detail && newMatch.sets_detail.sets) {
        const loaded = newMatch.sets_detail.sets;
        if (isGroupStage.value) {
          groupSet.value.a = loaded[0]?.a ?? null;
          groupSet.value.b = loaded[0]?.b ?? null;
        } else {
          knockoutSets.value = [
            { a: loaded[0]?.a ?? null, b: loaded[0]?.b ?? null },
            { a: loaded[1]?.a ?? null, b: loaded[1]?.b ?? null },
            { a: loaded[2]?.a ?? null, b: loaded[2]?.b ?? null },
          ];
        }
      }
    }
  },
  { immediate: true },
);

const isValid = computed(() => {
  if (isGroupStage.value) {
    if (groupSet.value.a === null || groupSet.value.b === null) return false;
    if (parseInt(groupSet.value.a) === parseInt(groupSet.value.b)) return false; // No ties in 1 set
    return true;
  } else {
    // Knockout: Must have winner (2 sets won)
    if (!winnerId.value) return false;

    // Validate set 1 and 2 must be complete
    const s1 = knockoutSets.value[0];
    const s2 = knockoutSets.value[1];
    if (s1.a === null || s1.b === null || s2.a === null || s2.b === null)
      return false;
    if (parseInt(s1.a) === parseInt(s1.b) || parseInt(s2.a) === parseInt(s2.b))
      return false; // No individual ties

    return true;
  }
});

const save = () => {
  if (!isValid.value)
    return alert("Invalid score. Ensure winner is determined (2 sets won).");

  let finalSets = [];
  if (isGroupStage.value) {
    finalSets = [
      { set: 1, a: parseInt(groupSet.value.a), b: parseInt(groupSet.value.b) },
    ];
  } else {
    finalSets = knockoutSets.value
      .map((s, i) => ({ set: i + 1, a: parseInt(s.a), b: parseInt(s.b) }))
      .filter(
        (s) => s.a !== null && s.b !== null && !isNaN(s.a) && !isNaN(s.b),
      );
  }

  emit("save", {
    id: props.match.id,
    winner_id: winnerId.value,
    score: calculatedSummaryScore.value,
    sets_detail: { sets: finalSets },
    video_url: videoUrl.value,
    status: "finished",
  });
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-gray-900/10 backdrop-blur-[2px]"
  >
    <div
      class="bg-white rounded-[4px] border border-gray-100 shadow-2xl w-full max-w-md overflow-hidden font-outfit"
    >
      <!-- Header -->
      <div
        class="px-6 py-4 border-b border-gray-50 flex justify-between items-center"
      >
        <h3 class="text-xs font-black text-gray-400 uppercase tracking-[0.2em]">
          Update Result
        </h3>
        <button
          @click="$emit('close')"
          class="text-gray-300 hover:text-gray-900 transition-colors text-lg italic tracking-tighter"
        >
          Exit
        </button>
      </div>

      <div class="p-8 space-y-10">
        <!-- Teams Display -->
        <div class="flex justify-between items-center px-2">
          <div class="w-2/5 text-center">
            <div class="text-[10px] font-bold text-gray-300 uppercase mb-2">
              Team A
            </div>
            <div
              class="text-sm font-black text-gray-900 leading-tight h-10 flex items-center justify-center"
            >
              {{ match.team_a?.name || "TBD" }}
            </div>
          </div>
          <div class="text-[10px] font-black text-gray-200">VS</div>
          <div class="w-2/5 text-center">
            <div class="text-[10px] font-bold text-gray-300 uppercase mb-2">
              Team B
            </div>
            <div
              class="text-sm font-black text-gray-900 leading-tight h-10 flex items-center justify-center"
            >
              {{ match.team_b?.name || "TBD" }}
            </div>
          </div>
        </div>

        <!-- Inputs Group -->
        <div v-if="isGroupStage" class="space-y-4">
          <div class="flex items-center gap-6 justify-center">
            <input
              v-model="groupSet.a"
              type="number"
              class="input-material text-center text-4xl font-black w-20"
              placeholder="0"
            />
            <span class="text-gray-200">-</span>
            <input
              v-model="groupSet.b"
              type="number"
              class="input-material text-center text-4xl font-black w-20"
              placeholder="0"
            />
          </div>
          <p
            class="text-[10px] font-bold text-gray-400 text-center uppercase tracking-widest"
          >
            1 Set - Max 21 Points
          </p>
        </div>

        <!-- Inputs Knockout -->
        <div v-else class="space-y-6">
          <div
            v-for="(set, i) in knockoutSets"
            :key="i"
            class="flex items-center gap-4 justify-center"
          >
            <span class="text-[10px] font-black text-gray-300 w-8"
              >SET {{ i + 1 }}</span
            >
            <input
              v-model="set.a"
              type="number"
              class="input-material text-center text-xl font-bold w-12"
              placeholder="0"
            />
            <input
              v-model="set.b"
              type="number"
              class="input-material text-center text-xl font-bold w-12"
              placeholder="0"
            />
          </div>
          <p
            class="text-[10px] font-bold text-gray-400 text-center uppercase tracking-widest"
          >
            BO3 - Max 21 Points
          </p>
        </div>

        <!-- Details -->
        <div class="space-y-1 pt-4">
          <label
            class="text-[10px] font-bold text-gray-400 uppercase tracking-widest"
            >Video Recording</label
          >
          <input
            v-model="videoUrl"
            type="text"
            class="input-material text-xs placeholder:text-gray-200"
            placeholder="https://youtube.com/..."
          />
        </div>
      </div>

      <!-- Footer -->
      <div v-if="isAdmin" class="p-6 bg-gray-50 flex justify-end gap-3">
        <button
          @click="save"
          :disabled="!isValid"
          class="btn-primary w-full max-w-[140px] uppercase tracking-widest text-[10px] font-black disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Save Score
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
