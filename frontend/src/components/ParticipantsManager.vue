<script setup>
import { ref, computed, onMounted } from "vue";
import api from "../services/api";

const props = defineProps({
  participants: { type: Array, required: true },
  teams: { type: Array, required: true },
});

onMounted(() => {
  console.log("ParticipantsManager Mounted");
  console.log("Participants:", props.participants);
  console.log("Teams:", props.teams);
});

const emit = defineEmits(["refresh"]);

// -- STATE --
const filters = ref({
  search: "",
  pool: "ALL",
  status: "ALL", // ALL, FREE, BUSY
});

const showTeamModal = ref(false);
const teamForm = ref({
  player1_id: "",
  player2_id: "",
  pool: "", // Auto-detected from P1
});

// -- COMPUTED --
const isAssigned = (id) =>
  props.teams.some((t) => t.player1_id === id || t.player2_id === id);

const getTeamName = (id) => {
  const t = props.teams.find((t) => t.player1_id === id || t.player2_id === id);
  return t ? t.name : null;
};

const filteredParticipants = computed(() => {
  return props.participants.filter((p) => {
    // 1. Search
    const matchSearch =
      !filters.value.search ||
      p.name.toLowerCase().includes(filters.value.search.toLowerCase()) ||
      p.email.toLowerCase().includes(filters.value.search.toLowerCase());

    // 2. Pool
    const matchPool =
      filters.value.pool === "ALL" || p.pool === filters.value.pool;

    // 3. Status
    const pAssigned = isAssigned(p.id);
    const matchStatus =
      filters.value.status === "ALL" ||
      (filters.value.status === "FREE" && !pAssigned) ||
      (filters.value.status === "BUSY" && pAssigned);

    return matchSearch && matchPool && matchStatus;
  });
});

// For Team Builder Modal
const freeParticipants = computed(() => {
  return props.participants.filter((p) => !isAssigned(p.id));
});

const p1Options = computed(() => freeParticipants.value);

const p2Options = computed(() => {
  if (!teamForm.value.player1_id) return [];
  const p1 = props.participants.find((p) => p.id === teamForm.value.player1_id);
  if (!p1) return [];
  // Must be same pool as P1
  return freeParticipants.value.filter(
    (p) => p.id !== p1.id && p.pool === p1.pool,
  );
});

// -- ACTIONS --
const openTeamModal = () => {
  teamForm.value = { player1_id: "", player2_id: "", pool: "" };
  showTeamModal.value = true;
};

const createTeam = async () => {
  try {
    const p1 = props.participants.find(
      (p) => p.id === teamForm.value.player1_id,
    );
    if (!p1) return;

    // Safety check (should cover backend logic too)
    const p2 = props.participants.find(
      (p) => p.id === teamForm.value.player2_id,
    );
    if (p1.pool !== p2.pool) {
      alert("Error: Players must be in the same pool.");
      return;
    }

    await api.post("/teams", {
      player1_id: teamForm.value.player1_id,
      player2_id: teamForm.value.player2_id,
    });

    showTeamModal.value = false;
    emit("refresh");
  } catch (err) {
    alert(
      "Failed to create team: " + (err.response?.data?.error || err.message),
    );
  }
};

const autoPairTeams = async () => {
  if (
    !confirm(
      "Auto-pair all remaining free participants into teams? Pairs will be random within same pool.",
    )
  )
    return;

  try {
    const res = await api.post("/teams/auto-pair", {
      tournament_id: "00000000-0000-0000-0000-000000000000",
    });
    alert(res.data.message);
    emit("refresh");
  } catch (err) {
    alert("Failed to auto-pair: " + (err.response?.data?.error || err.message));
  }
};
</script>

<template>
  <div class="space-y-6">
    <!-- Toolbar -->
    <div
      class="bg-white p-4 rounded-lg border border-gray-200 flex flex-wrap gap-4 items-center justify-between"
    >
      <div class="flex gap-4 flex-1">
        <input
          v-model="filters.search"
          type="text"
          placeholder="Search name/email..."
          class="input-material min-w-[200px]"
        />
        <select v-model="filters.pool" class="input-material w-32">
          <option value="ALL">All Pools</option>
          <option value="Mesoneer">Mesoneer</option>
          <option value="Lab">Lab</option>
        </select>
        <select v-model="filters.status" class="input-material w-32">
          <option value="ALL">All Status</option>
          <option value="FREE">Free</option>
          <option value="BUSY">In Team</option>
        </select>
        </select>
      </div>

      <div class="flex gap-2">
         <button
          @click="autoPairTeams"
          class="px-4 py-2 border border-violet-200 text-violet-600 font-medium rounded-sm hover:bg-violet-50 flex items-center gap-2"
        >
          <span>🎲 Auto-Pair</span>
        </button>
        <button
          @click="openTeamModal"
          class="px-4 py-2 bg-violet-600 text-white font-medium rounded-sm btn-animated flex items-center gap-2"
        >
          <span>+ Form Team</span>
        </button>
      </div>
    </div>

    <!-- Table -->
    <div
      class="bg-white rounded-lg shadow-sm border border-gray-100 overflow-hidden overflow-x-auto"
    >
      <table class="w-full text-sm text-left">
        <thead class="bg-gray-50 text-gray-500 font-medium uppercase text-xs">
          <tr>
            <th class="px-6 py-3">Name</th>
            <th class="px-6 py-3">Email</th>
            <th class="px-6 py-3">Pool</th>
            <th class="px-6 py-3">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100">
          <tr
            v-for="p in filteredParticipants"
            :key="p.id"
            class="hover:bg-gray-50 transition-colors"
          >
            <td class="px-6 py-3 font-medium text-gray-900">{{ p.name }}</td>
            <td class="px-6 py-3 text-gray-500">{{ p.email }}</td>
            <td class="px-6 py-3">
              <span
                :class="{
                  'bg-violet-100 text-violet-700': p.pool === 'Mesoneer',
                  'bg-blue-100 text-blue-700': p.pool === 'Lab',
                }"
                class="px-2 py-0.5 rounded-full text-xs font-medium"
              >
                {{ p.pool }}
              </span>
            </td>
            <td class="px-6 py-3">
              <span
                v-if="getTeamName(p.id)"
                class="bg-green-100 text-green-700 px-2 py-0.5 rounded-md text-xs font-medium border border-green-200"
              >
                Team: {{ getTeamName(p.id) }}
              </span>
              <span
                v-else
                class="bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full text-xs font-medium"
              >
                Free
              </span>
            </td>
          </tr>
          <tr v-if="filteredParticipants.length === 0">
            <td colspan="4" class="px-6 py-8 text-center text-gray-400">
              No participants found.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Logic: Team Builder Modal -->
    <div
      v-if="showTeamModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
    >
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6">
        <h3 class="text-lg font-bold mb-4">Form New Team</h3>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Player 1 (Free)</label
            >
            <select
              v-model="teamForm.player1_id"
              @change="teamForm.player2_id = ''"
              class="input-material"
            >
              <option value="" disabled>Select Player 1</option>
              <option v-for="p in p1Options" :key="p.id" :value="p.id">
                {{ p.name }} ({{ p.pool }})
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Player 2 (Same Pool)</label
            >
            <select
              v-model="teamForm.player2_id"
              :disabled="!teamForm.player1_id"
              class="input-material disabled:opacity-50"
            >
              <option value="" disabled>Select Player 2</option>
              <option v-for="p in p2Options" :key="p.id" :value="p.id">
                {{ p.name }}
              </option>
            </select>
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <button
            @click="showTeamModal = false"
            class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-sm transition-colors"
          >
            Cancel
          </button>
          <button
            @click="createTeam"
            :disabled="!teamForm.player1_id || !teamForm.player2_id"
            class="px-4 py-2 bg-violet-600 text-white rounded-sm hover:bg-violet-700 disabled:opacity-50 btn-animated"
          >
            Create Team
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
