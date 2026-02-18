<script setup>
import { ref, computed } from "vue";
import api from "../services/api";

const props = defineProps({
  teams: { type: Array, required: true },
});

const emit = defineEmits(["refresh"]);

// -- STATE --
const selectedPool = ref("Mesoneer");
const selectedTeamIds = ref([]);
const isProcessing = ref(false);

// -- COMPUTED --
const poolTeams = computed(() => {
  return props.teams.filter((t) => t.pool === selectedPool.value);
});

// -- MODAL STATE --
const showConfigModal = ref(false);
const configModal = ref({
  title: "",
  description: "",
  prefix: "", // Displayed before input
  value: "", // User input
  placeholder: "",
  mode: "", // 'MANUAL' or 'AUTO'
});

// -- ACTIONS --
const openManualCreationModal = () => {
  if (selectedTeamIds.value.length !== 4) return;

  configModal.value = {
    title: "Create Single Group",
    description: `Create a group for ${selectedPool.value} with the 4 selected teams.`,
    prefix: "",
    value: "Group A",
    placeholder: "Enter full group name",
    mode: "MANUAL",
  };
  showConfigModal.value = true;
};

const openAutoGenerateModal = () => {
  if (poolTeams.value.length === 0) {
    alert("No teams in this pool.");
    return;
  }
  if (poolTeams.value.length % 4 !== 0) {
    alert(
      `Cannot auto-generate: ${poolTeams.value.length} teams is not divisible by 4.`,
    );
    return;
  }

  configModal.value = {
    title: "Auto-Generate All Groups",
    description: `Randomly assign ${poolTeams.value.length} teams into ${poolTeams.value.length / 4} groups.`,
    prefix: "Group ",
    value: "A",
    placeholder: "Enter suffix (e.g. A)",
    mode: "AUTO",
  };
  showConfigModal.value = true;
};

const handleModalConfirm = async () => {
  if (!configModal.value.value) return;

  isProcessing.value = true;
  try {
    if (configModal.value.mode === "MANUAL") {
      const name = configModal.value.value;
      await api.post("/groups", {
        name: name,
        pool: selectedPool.value,
        tournament_id: "00000000-0000-0000-0000-000000000000",
        team_ids: selectedTeamIds.value,
      });
      alert("Group created successfully!");
      selectedTeamIds.value = [];
    } else {
      const prefix = configModal.value.prefix + configModal.value.value.trim();
      await api.post("/groups/auto-generate", {
        pool: selectedPool.value,
        tournament_id: "00000000-0000-0000-0000-000000000000",
        name_prefix: prefix,
      });
      alert("Groups auto-generated successfully!");
    }
    emit("refresh");
    showConfigModal.value = false;
  } catch (err) {
    console.error("Operation failed:", err);
    alert("Failed: " + (err.response?.data?.error || err.message));
  } finally {
    isProcessing.value = false;
  }
};
</script>

<template>
  <div class="space-y-6">
    <div
      class="flex justify-between items-center bg-white p-4 rounded-lg border border-gray-200"
    >
      <div class="space-x-4">
        <label class="font-medium text-gray-700">Select Pool:</label>
        <button
          @click="
            selectedPool = 'Mesoneer';
            selectedTeamIds = [];
          "
          :class="
            selectedPool === 'Mesoneer'
              ? 'bg-violet-100 text-violet-700 border-violet-200'
              : 'bg-white text-gray-500 border-gray-200 hover:bg-gray-50'
          "
          class="px-4 py-1 rounded-full text-sm font-bold border transition-colors"
        >
          Mesoneer
        </button>
        <button
          @click="
            selectedPool = 'Lab';
            selectedTeamIds = [];
          "
          :class="
            selectedPool === 'Lab'
              ? 'bg-blue-100 text-blue-700 border-blue-200'
              : 'bg-white text-gray-500 border-gray-200 hover:bg-gray-50'
          "
          class="px-4 py-1 rounded-full text-sm font-bold border transition-colors"
        >
          Lab
        </button>
      </div>

      <div class="text-sm text-gray-500 flex items-center gap-4">
        <div>
          Selected:
          <span class="font-bold text-gray-900">{{
            selectedTeamIds.length
          }}</span>
          / 4
        </div>
        <div class="h-4 w-px bg-gray-200"></div>

        <div class="flex items-center gap-2">
          <span class="text-xs font-medium text-gray-500">Group</span>
          <input
            v-model="groupNameSuffix"
            type="text"
            placeholder="A"
            class="w-12 text-center border border-gray-300 rounded-sm text-sm py-1 focus:ring-violet-500 focus:border-violet-500 uppercase"
            maxlength="5"
          />
        </div>

        <button
          @click="autoGenerate"
          :disabled="isProcessing || poolTeams.length === 0"
          class="text-violet-600 hover:text-violet-800 font-bold flex items-center gap-1 disabled:opacity-50"
        >
          <span>🎲</span>
          Auto-Generate All
        </button>
      </div>
    </div>

    <!-- Team Selection Grid -->
    <div
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
    >
      <div
        v-for="team in poolTeams"
        :key="team.id"
        @click="
          selectedTeamIds.includes(team.id)
            ? (selectedTeamIds = selectedTeamIds.filter((id) => id !== team.id))
            : selectedTeamIds.length < 4 && selectedTeamIds.push(team.id)
        "
        :class="
          selectedTeamIds.includes(team.id)
            ? 'ring-2 ring-violet-500 bg-violet-50'
            : 'hover:bg-gray-50 bg-white border-gray-200'
        "
        class="border rounded-lg p-4 cursor-pointer transition-all shadow-sm select-none"
      >
        <div class="font-bold text-gray-800 mb-1">{{ team.name }}</div>
        <div class="text-xs text-gray-500 flex gap-1">
          <span class="bg-gray-100 px-1 rounded">{{ team.pool }}</span>
        </div>
      </div>
    </div>

    <div v-if="poolTeams.length === 0" class="text-center py-12 text-gray-400">
      No teams in this pool. Go create some in "Participants & Matchmaking"!
    </div>

    <!-- Floating Action Button -->
    <transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="transform translate-y-20 opacity-0"
      enter-to-class="transform translate-y-0 opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="transform translate-y-0 opacity-100"
      leave-to-class="transform translate-y-20 opacity-0"
    >
      <div
        v-if="selectedTeamIds.length === 4"
        class="fixed bottom-8 right-8 z-30"
      >
        <button
          @click="openManualCreationModal"
          :disabled="isProcessing"
          class="bg-gradient-to-r from-violet-600 to-indigo-600 text-white font-bold py-3 px-6 rounded-full shadow-lg hover:shadow-xl hover:scale-105 transition-transform flex items-center gap-2"
        >
          <span v-if="isProcessing" class="animate-spin">⏳</span>
          <span v-else>✨</span>
          Create Group
        </button>
      </div>
    </transition>

    <!-- Custom Config Modal -->
    <div
      v-if="showConfigModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
    >
      <div class="bg-white rounded-lg shadow-xl w-full max-w-sm overflow-hidden font-outfit">
        <div class="bg-violet-600 px-6 py-4 flex justify-between items-center">
            <h3 class="font-bold text-white">{{ configModal.title }}</h3>
            <button @click="showConfigModal = false" class="text-violet-200 hover:text-white">✕</button>
        </div>
        
        <div class="p-6 space-y-4">
            <p class="text-sm text-gray-600">{{ configModal.description }}</p>
            
            <div class="space-y-1">
                <label class="block text-xs font-bold text-gray-500 uppercase">Group Name</label>
                <div class="flex items-center border border-gray-300 rounded-sm focus-within:ring-2 focus-within:ring-violet-500 overflow-hidden">
                    <span v-if="configModal.prefix" class="bg-gray-50 px-3 py-2 text-gray-500 font-medium border-r border-gray-200 select-none">
                        {{ configModal.prefix }}
                    </span>
                    <input 
                        v-model="configModal.value"
                        type="text" 
                        :placeholder="configModal.placeholder"
                        class="flex-1 px-3 py-2 outline-none text-gray-900 font-medium w-full"
                        @keyup.enter="handleModalConfirm"
                        autoFocus
                    />
                </div>
            </div>
            
            <div class="flex justify-end gap-3 pt-2">
                <button 
                  @click="showConfigModal = false"
                  class="px-4 py-2 text-gray-500 hover:bg-gray-50 rounded-sm text-sm font-medium"
                >Cancel</button>
                <button 
                  @click="handleModalConfirm"
                  :disabled="isProcessing || !configModal.value"
                  class="px-6 py-2 bg-violet-600 text-white font-bold rounded-sm hover:bg-violet-700 disabled:opacity-50 shadow-md transition-all active:scale-95 flex items-center gap-2"
                >
                    <span v-if="isProcessing" class="animate-spin">⏳</span>
                    <span>Confirm</span>
                </button>
            </div>
        </div>
      </div>
    </div>
  </div>
</template>
