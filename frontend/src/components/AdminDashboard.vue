<script setup>
import { ref, onMounted, computed, watch } from "vue";
import Wheel from "./Wheel.vue";
import BracketNode from "./BracketNode.vue";
import api from "../services/api";

const activePool = ref("Mesoneer");
const participants = ref([]);
const matches = ref([]);
const teams = ref([]);
const selectedTeamIds = ref([]);

const showResultModal = ref(false);
const selectedMatch = ref(null);

// Score State
const sets = ref([
  { a: 0, b: 0 },
  { a: 0, b: 0 },
  { a: 0, b: 0 },
]);

const loadData = async () => {
  try {
    const res = await api.getBracket(activePool.value);
    matches.value = res.data;

    // Also load teams for bracket gen
    const teamRes = await api.getTeams(activePool.value);
    teams.value = teamRes.data;
    selectedTeamIds.value = [];
  } catch (e) {
    console.error(e);
  }
};

const switchPool = (pool) => {
  activePool.value = pool;
  loadData();
};

const onTeamsGenerated = async () => {
  await api.generateTeams(activePool.value);
  loadData();
};

const toggleTeamSelection = (id) => {
  if (selectedTeamIds.value.includes(id)) {
    selectedTeamIds.value = selectedTeamIds.value.filter((x) => x !== id);
  } else {
    if (selectedTeamIds.value.length < 4) {
      selectedTeamIds.value.push(id);
    } else {
      alert("Select exactly 4 teams");
    }
  }
};

const generateBracket = async () => {
  if (selectedTeamIds.value.length !== 4) {
    alert("Please select exactly 4 teams.");
    return;
  }
  try {
    await api.generateBracket(selectedTeamIds.value, activePool.value);
    loadData();
  } catch (e) {
    alert("Failed to generate bracket: " + e.response?.data?.error);
  }
};

const onUpdateResult = (match) => {
  selectedMatch.value = { ...match, video_url: match.video_url || "" };
  // Reset sets
  sets.value = [
    { a: 0, b: 0 },
    { a: 0, b: 0 },
    { a: 0, b: 0 },
  ];
  // Parse existing details if needed? For now start fresh or simple parse
  // If match has score_details, ideally parse it.
  showResultModal.value = true;
};

const saveResult = async () => {
  // Calculate Winner
  let winsA = 0;
  let winsB = 0;
  let details = [];

  for (const s of sets.value) {
    if (s.a === 0 && s.b === 0) continue; // Skip empty
    details.push(`${s.a}-${s.b}`);
    if (s.a > s.b) winsA++;
    else if (s.b > s.a) winsB++;
  }

  const scoreDetails = details.join(", ");

  try {
    await api.updateMatch(selectedMatch.value.id, {
      score_a: winsA,
      score_b: winsB,
      score_details: scoreDetails,
      video_url: selectedMatch.value.video_url,
    });
    showResultModal.value = false;
    loadData();
  } catch (e) {
    alert("Failed to save");
  }
};

onMounted(loadData);
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto">
    <header class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold text-primary">Tournament Admin</h1>
      <div class="flex gap-2">
        <button
          @click="switchPool('Mesoneer')"
          :class="[
            'px-4 py-2 rounded-sm font-medium transition-colors',
            activePool === 'Mesoneer'
              ? 'bg-primary text-white'
              : 'bg-white border border-border text-gray-600',
          ]"
        >
          Mesoneer
        </button>
        <button
          @click="switchPool('Lab')"
          :class="[
            'px-4 py-2 rounded-sm font-medium transition-colors',
            activePool === 'Lab'
              ? 'bg-primary text-white'
              : 'bg-white border border-border text-gray-600',
          ]"
        >
          Lab
        </button>
      </div>
    </header>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Sidebar: Actions -->
      <div class="space-y-8">
        <Wheel
          :pool="activePool"
          :participants="[]"
          @generated="onTeamsGenerated"
        />

        <div class="bg-white p-6 rounded-lg border border-border">
          <h3 class="font-bold text-primary mb-4">Bracket Actions</h3>
          <p class="text-xs text-slate-500 mb-2">
            Select 4 Teams to start a GSL Group:
          </p>

          <div
            class="max-h-48 overflow-y-auto space-y-2 mb-4 border border-slate-100 p-2 rounded-sm"
          >
            <div v-for="t in teams" :key="t.id" class="flex items-center gap-2">
              <input
                type="checkbox"
                :id="'team-' + t.id"
                :checked="selectedTeamIds.includes(t.id)"
                @change="toggleTeamSelection(t.id)"
                class="rounded-sm border-gray-300 text-primary focus:ring-primary"
              />
              <label
                :for="'team-' + t.id"
                class="text-sm text-slate-700 truncate cursor-pointer"
                >{{ t.name }}</label
              >
            </div>
            <div v-if="teams.length === 0" class="text-xs text-gray-400 italic">
              No teams available. Spin the wheel!
            </div>
          </div>

          <button
            @click="generateBracket"
            :disabled="selectedTeamIds.length !== 4"
            class="w-full py-3 bg-gray-100 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed rounded-sm font-medium text-gray-700 transition-colors"
          >
            Auto-Generate GSL Bracket
          </button>
        </div>
      </div>

      <!-- Main: Bracket -->
      <div class="lg:col-span-2 space-y-8">
        <h2 class="text-xl font-bold text-gray-800 border-b border-border pb-2">
          Bracket: {{ activePool }}
        </h2>

        <div
          v-if="matches.length === 0"
          class="text-center py-20 text-gray-400"
        >
          No matches generated yet. Select 4 teams to start.
        </div>

        <div v-else class="flex flex-wrap gap-8 justify-center relative">
          <div class="flex flex-col gap-8">
            <BracketNode
              v-for="m in matches.slice(0, 2)"
              :key="m.id"
              :match="m"
              :admin="true"
              @update-result="onUpdateResult"
            />
          </div>
          <div class="flex flex-col gap-8 justify-center">
            <BracketNode
              v-for="m in matches.slice(2, 4)"
              :key="m.id"
              :match="m"
              :admin="true"
              @update-result="onUpdateResult"
            />
          </div>
          <div class="flex flex-col justify-center">
            <BracketNode
              v-if="matches[4]"
              :match="matches[4]"
              :admin="true"
              @update-result="onUpdateResult"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div
      v-if="showResultModal"
      class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
    >
      <div class="bg-white p-6 rounded-lg w-full max-w-md space-y-4 shadow-xl">
        <h3 class="font-bold text-lg text-primary">Update Match Result</h3>
        <p class="text-sm text-slate-500">
          {{ selectedMatch?.team_a?.name || "TBD" }} vs
          {{ selectedMatch?.team_b?.name || "TBD" }}
        </p>

        <!-- Set Scores -->
        <div class="space-y-3 bg-slate-50 p-4 rounded-sm">
          <div
            class="flex justify-between text-xs font-bold text-slate-500 uppercase"
          >
            <span>Set</span>
            <span>Team A</span>
            <span>Team B</span>
          </div>
          <div
            v-for="(set, idx) in sets"
            :key="idx"
            class="flex items-center gap-4"
          >
            <span class="text-xs font-bold w-6 text-slate-400"
              >#{{ idx + 1 }}</span
            >
            <input
              v-model.number="set.a"
              type="number"
              class="w-full p-2 border border-border rounded-sm text-center"
            />
            <span class="text-slate-300">-</span>
            <input
              v-model.number="set.b"
              type="number"
              class="w-full p-2 border border-border rounded-sm text-center"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-bold text-gray-500 mb-1"
            >Video URL (YouTube)</label
          >
          <input
            v-model="selectedMatch.video_url"
            type="text"
            class="w-full p-2 border border-border rounded-sm"
            placeholder="https://youtube.com/..."
          />
        </div>

        <div class="flex justify-end gap-2 pt-4">
          <button
            @click="showResultModal = false"
            class="px-4 py-2 hover:bg-gray-100 rounded-sm font-medium"
          >
            Cancel
          </button>
          <button
            @click="saveResult"
            class="px-4 py-2 bg-primary text-white rounded-sm hover:bg-violet-700 font-medium"
          >
            Save Final Score
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
