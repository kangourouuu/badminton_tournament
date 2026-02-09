<script setup>
import { ref, onMounted, computed } from "vue";
import api from "../services/api";

const teams = ref([]);
const participants = ref([]);
const loading = ref(true);

// Edit Modal
const showEditModal = ref(false);
const editingTeam = ref(null);
const editForm = ref({
  player1_id: "",
  player2_id: "",
});

const fetchTeams = async () => {
  try {
    const res = await api.get("/teams");
    teams.value = res.data || [];
  } catch (err) {
    console.error("Failed to fetch teams", err);
  }
};

const fetchParticipants = async () => {
  try {
    const res = await api.get("/participants");
    participants.value = res.data || [];
  } catch (err) {
    console.error("Failed to fetch participants", err);
  }
};

const refresh = async () => {
  loading.value = true;
  await Promise.all([fetchTeams(), fetchParticipants()]);
  loading.value = false;
};

// Actions
const disbandTeam = async (id) => {
  if (!confirm("Are you sure? This will free the players.")) return;
  try {
    await api.delete(`/teams/${id}`);
    teams.value = teams.value.filter((t) => t.id !== id);
  } catch (err) {
    alert("Failed to disband: " + (err.response?.data?.error || err.message));
  }
};

const openEdit = (team) => {
  editingTeam.value = team;
  editForm.value = {
    player1_id: team.player1_id,
    player2_id: team.player2_id,
  };
  showEditModal.value = true;
};

const saveEdit = async () => {
  try {
    await api.put(`/teams/${editingTeam.value.id}`, editForm.value);
    showEditModal.value = false;
    refresh(); // Refresh to see name changes etc
  } catch (err) {
    alert("Failed to update: " + (err.response?.data?.error || err.message));
  }
};

// Helpers
const getPoolParticipants = (pool) => {
  return participants.value.filter((p) => p.pool === pool);
};

onMounted(refresh);
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h3 class="text-xl font-bold text-gray-800">Team Management</h3>
      <button @click="refresh" class="text-sm text-violet-600 hover:underline">
        Refresh
      </button>
    </div>

    <!-- Team List -->
    <div v-if="loading" class="text-center py-8 text-gray-400">
      Loading Teams...
    </div>
    <div v-else class="overflow-hidden bg-white shadow sm:rounded-md">
      <ul role="list" class="divide-y divide-gray-200">
        <li
          v-for="team in teams"
          :key="team.id"
          class="px-4 py-4 sm:px-6 flex items-center justify-between hover:bg-gray-50"
        >
          <div class="flex flex-col">
            <span class="text-sm font-medium text-violet-600 truncate">{{
              team.name
            }}</span>
            <span class="text-xs text-gray-500">
              {{ team.pool }} •
              <span class="font-mono text-gray-400">{{
                team.id.slice(0, 8)
              }}</span>
            </span>
          </div>
          <div class="flex gap-2">
            <button
              @click="openEdit(team)"
              class="text-indigo-600 hover:text-indigo-900 border border-indigo-200 px-3 py-1 rounded text-xs font-semibold hover:bg-indigo-50 transition"
            >
              Edit
            </button>
            <button
              @click="disbandTeam(team.id)"
              class="text-red-600 hover:text-red-900 border border-red-200 px-3 py-1 rounded text-xs font-semibold hover:bg-red-50 transition"
            >
              Disband
            </button>
          </div>
        </li>
      </ul>
    </div>

    <!-- Edit Modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
    >
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6">
        <h3 class="text-lg font-bold mb-4">Edit Team</h3>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700"
              >Player 1</label
            >
            <select
              v-model="editForm.player1_id"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-violet-500 focus:ring-violet-500 sm:text-sm"
            >
              <option
                v-for="p in getPoolParticipants(editingTeam.pool)"
                :key="p.id"
                :value="p.id"
              >
                {{ p.name }} ({{ p.email }})
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700"
              >Player 2</label
            >
            <select
              v-model="editForm.player2_id"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-violet-500 focus:ring-violet-500 sm:text-sm"
            >
              <option
                v-for="p in getPoolParticipants(editingTeam.pool)"
                :key="p.id"
                :value="p.id"
              >
                {{ p.name }} ({{ p.email }})
              </option>
            </select>
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <button
            @click="showEditModal = false"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            @click="saveEdit"
            class="px-4 py-2 text-sm font-medium text-white bg-violet-600 rounded-md hover:bg-violet-700"
          >
            Save Changes
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
