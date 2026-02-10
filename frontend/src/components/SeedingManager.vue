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

// -- ACTIONS --
const randomizeAndCreate = async () => {
  if (selectedTeamIds.value.length !== 4) return;
  if (
    !confirm(
      `Confirm create group for ${selectedPool.value}? This will shuffle the 4 selected teams and generate matches.`,
    )
  )
    return;

  isProcessing.value = true;
  try {
    // Determine group name (e.g. "Group X") - For now let's just use a prompt or auto-gen
    const name = prompt("Enter Group Name (e.g. 'Group A'):", "Group A");
    if (!name) return;

    await api.post("/groups", {
      name: name,
      pool: selectedPool.value,
      tournament_id: "00000000-0000-0000-0000-000000000000",
      team_ids: selectedTeamIds.value,
    });

    alert("Group created and seeded successfully!");
    selectedTeamIds.value = [];
    emit("refresh");
  } catch (err) {
    console.error("Group creation failed:", err);
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

      <div class="text-sm text-gray-500">
        Selected:
        <span class="font-bold text-gray-900">{{
          selectedTeamIds.length
        }}</span>
        / 4
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
          @click="randomizeAndCreate"
          :disabled="isProcessing"
          class="bg-gradient-to-r from-violet-600 to-indigo-600 text-white font-bold py-3 px-6 rounded-full shadow-lg hover:shadow-xl hover:scale-105 transition-transform flex items-center gap-2"
        >
          <span v-if="isProcessing" class="animate-spin">⏳</span>
          <span v-else>🎲</span>
          Randomize & Create Group
        </button>
      </div>
    </transition>
  </div>
</template>
