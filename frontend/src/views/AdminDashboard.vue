<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import api from "../services/api";
import TheWheel from "../components/TheWheel.vue";
import GSLGrid from "../components/GSLGrid.vue";
import ScoreModal from "../components/ScoreModal.vue";

const router = useRouter();
const activeTab = ref("brackets"); // 'data', 'teams', 'brackets'

// Data State
const participants = ref([]);
const teams = ref([]);
const groups = ref([]);
const loading = ref(false);

// Modal State
const showCreateGroupModal = ref(false);
const showScoreModal = ref(false);
const selectedMatch = ref(null);

// Create Group Form
const newGroupName = ref("");
const selectedTeamIds = ref([]);

// -- INIT --
const fetchData = async () => {
  loading.value = true;
  try {
    const [pRes, tRes, gRes] = await Promise.all([
      api.get("/participants"),
      api.get("/teams"),
      api.get("/groups"),
    ]);
    participants.value = pRes.data || [];
    teams.value = tRes.data || [];
    groups.value = gRes.data || [];
  } catch (err) {
    console.error("Fetch error", err);
    if (err.response?.status === 401) {
      router.push("/login");
    }
  } finally {
    loading.value = false;
  }
};

onMounted(fetchData);

const logout = () => {
  localStorage.removeItem("token");
  router.push("/login");
};

// -- TEAMS --
const onTeamsGenerated = () => {
  fetchData(); // Refresh teams list
};

// -- GROUPS --
const availableTeams = computed(() => {
  // In a real app, filter out teams already in a group.
  // For now, just show all teams sorted by creation
  return teams.value;
});

const createGroup = async () => {
  if (selectedTeamIds.value.length !== 4 || !newGroupName.value) return;

  try {
    await api.post("/groups", {
      name: newGroupName.value,
      tournament_id: "00000000-0000-0000-0000-000000000000", // Placeholder or fetch actual
      team_ids: selectedTeamIds.value,
    });
    showCreateGroupModal.value = false;
    newGroupName.value = "";
    selectedTeamIds.value = [];
    fetchData();
  } catch (err) {
    alert("Failed to create group: " + err.message);
  }
};

// -- MATCHES --
const openScoreModal = (match) => {
  selectedMatch.value = match;
  showScoreModal.value = true;
};

const saveMatchResult = async (data) => {
  try {
    await api.post(`/matches/${data.id}`, {
      winner_id: data.winner_id,
      score: data.score,
      video_url: data.video_url,
    });
    showScoreModal.value = false;
    fetchData();
  } catch (err) {
    alert("Failed to update match: " + err.message);
  }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 pb-12">
    <!-- Navbar -->
    <nav
      class="bg-white border-b border-purple-100 px-6 py-4 flex justify-between items-center sticky top-0 z-20"
    >
      <h1 class="text-xl font-bold text-violet-800">🏸 Admin Dashboard</h1>
      <div class="flex items-center space-x-4">
        <span v-if="loading" class="text-xs text-gray-400 animate-pulse"
          >Syncing...</span
        >
        <button
          @click="logout"
          class="text-sm text-red-500 hover:text-red-700 font-medium"
        >
          Logout
        </button>
      </div>
    </nav>

    <!-- Tabs -->
    <div class="bg-white border-b border-gray-200 px-6 pt-2">
      <div class="flex space-x-8">
        <button
          @click="activeTab = 'brackets'"
          :class="{
            'border-violet-600 text-violet-700': activeTab === 'brackets',
            'border-transparent text-gray-500 hover:text-gray-700':
              activeTab !== 'brackets',
          }"
          class="pb-3 px-1 border-b-2 font-medium text-sm transition-colors"
        >
          Brackets & Matches
        </button>
        <button
          @click="activeTab = 'teams'"
          :class="{
            'border-violet-600 text-violet-700': activeTab === 'teams',
            'border-transparent text-gray-500 hover:text-gray-700':
              activeTab !== 'teams',
          }"
          class="pb-3 px-1 border-b-2 font-medium text-sm transition-colors"
        >
          Team Generation
        </button>
        <button
          @click="activeTab = 'data'"
          :class="{
            'border-violet-600 text-violet-700': activeTab === 'data',
            'border-transparent text-gray-500 hover:text-gray-700':
              activeTab !== 'data',
          }"
          class="pb-3 px-1 border-b-2 font-medium text-sm transition-colors"
        >
          Participants Data
        </button>
      </div>
    </div>

    <!-- Content -->
    <main class="container mx-auto px-6 py-8">
      <!-- Tab: BRACKETS -->
      <div v-show="activeTab === 'brackets'" class="space-y-8">
        <div class="flex justify-between items-center">
          <h2 class="text-lg font-bold text-gray-800">Active Groups</h2>
          <button
            @click="showCreateGroupModal = true"
            class="px-4 py-2 bg-violet-600 text-white rounded-sm hover:bg-violet-700 shadow-sm transition-colors flex items-center gap-2"
          >
            <span>+ Create Group</span>
          </button>
        </div>

        <div
          v-if="groups.length === 0"
          class="text-center py-12 bg-white rounded-lg border border-dashed border-gray-200 text-gray-400"
        >
          No groups created yet. Go create one!
        </div>

        <div v-for="group in groups" :key="group.id" class="space-y-4">
          <div class="flex items-center gap-3">
            <h3 class="text-xl font-bold text-gray-800">{{ group.name }}</h3>
            <span
              class="text-xs text-gray-400 bg-gray-100 px-2 py-0.5 rounded-full"
              >{{ group.id.split("-")[0] }}</span
            >
          </div>
          <GSLGrid
            :matches="group.matches || []"
            :is-admin="true"
            @match-click="openScoreModal"
          />
        </div>
      </div>

      <!-- Tab: TEAMS -->
      <div v-show="activeTab === 'teams'" class="space-y-8">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <!-- Wheel 1: Mesoneer -->
          <div class="space-y-4">
            <h3
              class="font-bold text-gray-700 border-l-4 border-violet-500 pl-2"
            >
              Mesoneer Pool
            </h3>
            <TheWheel pool="Mesoneer" @teams-generated="onTeamsGenerated" />
          </div>
          <!-- Wheel 2: Lab -->
          <div class="space-y-4">
            <h3 class="font-bold text-gray-700 border-l-4 border-blue-500 pl-2">
              Lab Pool
            </h3>
            <TheWheel pool="Lab" @teams-generated="onTeamsGenerated" />
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-sm border border-gray-100 p-6">
          <h3 class="font-bold text-gray-800 mb-4">All Generated Teams</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
            <div
              v-for="team in teams"
              :key="team.id"
              class="p-3 bg-gray-50 border border-gray-100 rounded-sm text-sm"
            >
              <div class="font-medium text-gray-900">{{ team.name }}</div>
              <div class="text-xs text-gray-500 mt-1">{{ team.pool }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Tab: DATA -->
      <div v-show="activeTab === 'data'" class="space-y-6">
        <div
          class="bg-white rounded-lg shadow-sm border border-gray-100 overflow-hidden"
        >
          <table class="w-full text-sm text-left">
            <thead
              class="bg-gray-50 text-gray-500 font-medium uppercase text-xs"
            >
              <tr>
                <th class="px-6 py-3">Name</th>
                <th class="px-6 py-3">Email</th>
                <th class="px-6 py-3">Pool</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr
                v-for="p in participants"
                :key="p.id"
                class="hover:bg-gray-50 transition-colors"
              >
                <td class="px-6 py-3 font-medium text-gray-900">
                  {{ p.name }}
                </td>
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
              </tr>
              <tr v-if="participants.length === 0">
                <td colspan="3" class="px-6 py-8 text-center text-gray-400">
                  No participants found.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>

    <!-- Modal: Create Group -->
    <div
      v-if="showCreateGroupModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
    >
      <div class="bg-white rounded-lg shadow-xl w-full max-w-lg p-6 space-y-6">
        <h3 class="text-xl font-bold text-gray-800">
          Create New Bracket Group
        </h3>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Group Name</label
            >
            <input
              v-model="newGroupName"
              type="text"
              placeholder="e.g. Group A"
              class="w-full border border-gray-300 rounded-sm p-2 focus:ring-2 focus:ring-violet-500"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >Select 4 Teams</label
            >
            <div
              class="h-48 overflow-y-auto border border-gray-200 rounded-sm p-2 space-y-1 bg-gray-50"
            >
              <label
                v-for="team in teams"
                :key="team.id"
                class="flex items-center space-x-2 p-2 hover:bg-white rounded-sm cursor-pointer transition-colors"
              >
                <input
                  type="checkbox"
                  :value="team.id"
                  v-model="selectedTeamIds"
                  :disabled="
                    selectedTeamIds.length >= 4 &&
                    !selectedTeamIds.includes(team.id)
                  "
                  class="text-violet-600 focus:ring-violet-500 rounded-sm"
                />
                <span class="text-sm text-gray-700"
                  >{{ team.name }}
                  <span class="text-xs text-gray-400"
                    >({{ team.pool }})</span
                  ></span
                >
              </label>
            </div>
            <p class="text-xs text-gray-500 mt-1 text-right">
              {{ selectedTeamIds.length }} / 4 selected
            </p>
          </div>
        </div>

        <div class="flex justify-end space-x-3 pt-4 border-t border-gray-100">
          <button
            @click="showCreateGroupModal = false"
            class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-sm transition-colors"
          >
            Cancel
          </button>
          <button
            @click="createGroup"
            :disabled="selectedTeamIds.length !== 4 || !newGroupName"
            class="px-4 py-2 bg-violet-600 text-white rounded-sm hover:bg-violet-700 disabled:opacity-50 transition-colors"
          >
            Create Bracket
          </button>
        </div>
      </div>
    </div>

    <!-- Modal: Score -->
    <ScoreModal
      :is-open="showScoreModal"
      :match="selectedMatch || {}"
      @close="showScoreModal = false"
      @save="saveMatchResult"
    />
  </div>
</template>
