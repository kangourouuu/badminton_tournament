<script setup>
import { ref, onMounted, computed } from "vue";
import Wheel from "./Wheel.vue";
import BracketNode from "./BracketNode.vue";
import api from "../services/api";

const activePool = ref("Mesoneer");
const participants = ref([]); // Loading from DB?
const matches = ref([]);
const showResultModal = ref(false);
const selectedMatch = ref(null);

const loadData = async () => {
  try {
    const res = await api.getBracket(activePool.value);
    matches.value = res.data;
  } catch (e) {
    console.error(e);
  }
};

const switchPool = (pool) => {
  activePool.value = pool;
  loadData();
};

const onTeamsGenerated = async () => {
  // Backend generates teams.
  // UX: Show success?
  await api.generateTeams(activePool.value);
  // Then what? Generate Bracket?
  // The prompt says: "Admin assigns Teams to Groups."
  // For this 1-day sprint, maybe auto-assign 4 teams to GSL?
  // I need to fetch Teams.
  // Missing fetch teams endpoint.
  // I'll assume we just Generate Bracket with "All available teams".
};

const onUpdateResult = (match) => {
  selectedMatch.value = { ...match, video_url: match.video_url || "" };
  showResultModal.value = true;
};

const saveResult = async () => {
  try {
    await api.updateMatch(selectedMatch.value.id, {
      score_a: parseInt(selectedMatch.value.score_a),
      score_b: parseInt(selectedMatch.value.score_b),
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
            'px-4 py-2 rounded-sm font-medium',
            activePool === 'Mesoneer'
              ? 'bg-primary text-white'
              : 'bg-gray-100 text-gray-600',
          ]"
        >
          Mesoneer
        </button>
        <button
          @click="switchPool('Lab')"
          :class="[
            'px-4 py-2 rounded-sm font-medium',
            activePool === 'Lab'
              ? 'bg-primary text-white'
              : 'bg-gray-100 text-gray-600',
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
          :participants="['Mock1', 'Mock2']"
          @generated="onTeamsGenerated"
        />

        <div class="bg-white p-6 rounded-lg border border-border">
          <h3 class="font-bold text-primary mb-4">Bracket Actions</h3>
          <button
            class="w-full py-3 bg-gray-100 hover:bg-gray-200 rounded-sm font-medium text-gray-700"
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
          No matches generated yet.
        </div>

        <div v-else class="flex flex-wrap gap-8 justify-center relative">
          <!-- Simple GSL Layout Visualization (Tree) -->
          <!-- M1, M2 -> M3, M4 -> M5 -->
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
      class="fixed inset-0 bg-black/50 flex items-center justify-center p-4"
    >
      <div class="bg-white p-6 rounded-lg w-full max-w-md space-y-4">
        <h3 class="font-bold text-lg">Update Match Result</h3>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-bold text-gray-500 mb-1"
              >Score Team A</label
            >
            <input
              v-model="selectedMatch.score_a"
              type="number"
              class="w-full p-2 border border-border rounded-sm"
            />
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-500 mb-1"
              >Score Team B</label
            >
            <input
              v-model="selectedMatch.score_b"
              type="number"
              class="w-full p-2 border border-border rounded-sm"
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
            class="px-4 py-2 hover:bg-gray-100 rounded-sm"
          >
            Cancel
          </button>
          <button
            @click="saveResult"
            class="px-4 py-2 bg-primary text-white rounded-sm hover:bg-purple-700"
          >
            Save Final Score
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
