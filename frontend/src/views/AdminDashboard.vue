<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import api from "../services/api";
import GSLGrid from "../components/GSLGrid.vue";
import KnockoutGrid from "../components/KnockoutGrid.vue";
import ScoreModal from "../components/ScoreModal.vue";
import ParticipantsManager from "../components/ParticipantsManager.vue";
import SeedingManager from "../components/SeedingManager.vue";

const router = useRouter();
const activeTab = ref("brackets"); // 'brackets', 'seeding', 'participants', 'rules'

// Data State
const participants = ref([]);
const teams = ref([]);
const groups = ref([]);
const loading = ref(false);

// Modal State
const showScoreModal = ref(false);
const selectedMatch = ref(null);

// -- INIT --
const fetchData = async () => {
  loading.value = true;
  try {
    const results = await Promise.allSettled([
      api.get("/participants"),
      api.get("/teams"),
      api.get("/groups"),
      api.get("/public/rules"),
    ]);

    // 0: Participants
    if (results[0].status === "fulfilled") {
      participants.value = results[0].value.data || [];
    } else {
      console.error("Failed to load participants", results[0].reason);
    }

    // 1: Teams
    if (results[1].status === "fulfilled") {
      teams.value = results[1].value.data || [];
    } else {
      console.error("Failed to load teams", results[1].reason);
    }

    // 2: Groups
    if (results[2].status === "fulfilled") {
      groups.value = results[2].value.data || [];
    } else {
      console.error("Failed to load groups", results[2].reason);
    }

    // 3: Rules
    if (results[3].status === "fulfilled") {
      rulesContent.value = results[3].value.data.content || "";
    } else {
      console.error("Failed to load rules", results[3].reason);
    }
  } catch (err) {
    console.error("Critical Fetch error", err);
    if (err.response?.status === 401) {
      router.push("/login"); // Handle global auth error if needed
    }
  } finally {
    loading.value = false;
  }
};

// -- COMPUTED --
const availableTeams = computed(() => {
  if (!teams.value.length) return [];
  const assignedTeamIds = new Set();

  groups.value.forEach((g) => {
    if (g.matches) {
      g.matches.forEach((m) => {
        if (m.team_a_id) assignedTeamIds.add(m.team_a_id);
        if (m.team_b_id) assignedTeamIds.add(m.team_b_id);
      });
    }
  });

  return teams.value.filter((t) => !assignedTeamIds.has(t.id));
});

onMounted(fetchData);

const logout = () => {
  localStorage.removeItem("token");
  router.push("/login");
};

// -- MATCHES --
const selectedMatchStage = ref("KNOCKOUT"); // Default to detailed

const openScoreModal = (match, stage = "KNOCKOUT") => {
  selectedMatch.value = match;
  selectedMatchStage.value = stage;
  showScoreModal.value = true;
};

const saveMatchResult = async (data) => {
  try {
    await api.post(`/matches/${data.id}`, {
      winner_id: data.winner_id,
      score: data.score,
      sets_detail: data.sets_detail,
      video_url: data.video_url,
    });
    showScoreModal.value = false;
    fetchData(); // Refresh to update brackets
  } catch (err) {
    alert("Failed to update match: " + err.message);
  }
};
// ... (rest of simple script logic if any remains before template starts?)

// ... (rest of simple script logic if any remains before template starts?)
// Actually lines 121-144 were garbage.
// The real code continues at line 145 (generateKnockout) in previous version but here I am creating a replacement block that bridges the gap properly.
// Wait, looking at file content:
// Line 118: };
// Line 119: // ... (rest of script)
// Line 121: // ... (in template)
// ... garbage ...
// Line 144: // -- BRACKETS AUTOMATION --
// So I need to replace from line 119 to 144 with *nothing* (or just the commented out section if needed, but better clear it).

// AND update template below.
// Since replace_file_content works on contiguous blocks or multi-chunks.
// I will use multi-chunks.

// Chunk 1: Remove garbage from script.
// Chunk 2: Update Grids in Template.
// Chunk 3: Update ScoreModal in Template.

// ... (rest of simple script logic if any remains before template starts?)
// Actually lines 121-144 were garbage.
// The real code continues at line 145 (generateKnockout) in previous version but here I am creating a replacement block that bridges the gap properly.
// Wait, looking at file content:
// Line 118: };
// Line 119: // ... (rest of script)
// Line 121: // ... (in template)
// ... garbage ...
// Line 144: // -- BRACKETS AUTOMATION --
// So I need to replace from line 119 to 144 with *nothing* (or just the commented out section if needed, but better clear it).

// AND update template below.
// Since replace_file_content works on contiguous blocks or multi-chunks.
// I will use multi-chunks.

// Chunk 1: Remove garbage from script.
// Chunk 2: Update Grids in Template.
// Chunk 3: Update ScoreModal in Template.

// -- RULES --
const rulesContent = ref("");
const savingRules = ref(false);

const saveRules = async () => {
  savingRules.value = true;
  try {
    await api.put("/admin/rules", { content: rulesContent.value });
    alert("Rules saved successfully!");
  } catch (err) {
    alert("Failed to save rules: " + err.message);
  } finally {
    savingRules.value = false;
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
    <div class="bg-white border-b border-gray-200 px-6 py-4">
      <div class="flex space-x-2 overflow-x-auto no-scrollbar">
        <button
          @click="activeTab = 'brackets'"
          :class="
            activeTab === 'brackets'
              ? 'bg-violet-100 text-violet-700 border-violet-200'
              : 'text-gray-500 hover:bg-gray-50 border-transparent'
          "
          class="px-4 py-2 rounded-full font-bold text-sm border transition-colors whitespace-nowrap"
        >
          Brackets & Matches
        </button>
        <button
          @click="activeTab = 'seeding'"
          :class="
            activeTab === 'seeding'
              ? 'bg-violet-100 text-violet-700 border-violet-200'
              : 'text-gray-500 hover:bg-gray-50 border-transparent'
          "
          class="px-4 py-2 rounded-full font-bold text-sm border transition-colors whitespace-nowrap"
        >
          Group Seeding
        </button>
        <button
          @click="activeTab = 'participants'"
          :class="
            activeTab === 'participants'
              ? 'bg-violet-100 text-violet-700 border-violet-200'
              : 'text-gray-500 hover:bg-gray-50 border-transparent'
          "
          class="px-4 py-2 rounded-full font-bold text-sm border transition-colors whitespace-nowrap"
        >
          Participants & Teams
        </button>
        <button
          @click="activeTab = 'rules'"
          :class="
            activeTab === 'rules'
              ? 'bg-violet-100 text-violet-700 border-violet-200'
              : 'text-gray-500 hover:bg-gray-50 border-transparent'
          "
          class="px-4 py-2 rounded-full font-bold text-sm border transition-colors whitespace-nowrap"
        >
          Tournament Rules
        </button>
      </div>
    </div>

    <!-- Content -->
    <main class="container mx-auto px-6 py-8">
      <!-- Tab: BRACKETS -->
      <div v-show="activeTab === 'brackets'" class="space-y-8">
        <div class="flex justify-between items-center">
          <h2 class="text-lg font-bold text-gray-800">Active Groups</h2>
        </div>

        <div
          v-if="groups.length === 0"
          class="text-center py-12 bg-white rounded-lg border border-dashed border-gray-200 text-gray-400"
        >
          No groups created yet. Go create one in "Group Seeding"!
        </div>

        <div v-for="group in groups" :key="group.id" class="space-y-4">
          <h3 class="text-xl font-bold text-gray-800">{{ group.name }}</h3>

          <div
            class="bg-white rounded-xl shadow-sm border border-gray-200 p-4 overflow-hidden"
          >
            <KnockoutGrid
              v-if="group.type === 'KNOCKOUT'"
              :matches="group.matches || []"
              :is-admin="true"
              @match-click="(m) => openScoreModal(m, 'KNOCKOUT')"
            />
            <GSLGrid
              v-else
              :matches="group.matches || []"
              :is-admin="true"
              @match-click="(m) => openScoreModal(m, 'GROUP')"
            />
          </div>
        </div>
      </div>

      <!-- Tab: SEEDING -->
      <div v-show="activeTab === 'seeding'" class="space-y-6">
        <SeedingManager :teams="availableTeams" @refresh="fetchData" />
      </div>

      <!-- Tab: PARTICIPANTS -->
      <div v-show="activeTab === 'participants'" class="space-y-6">
        <ParticipantsManager
          :participants="participants"
          :teams="teams"
          @refresh="fetchData"
        />
      </div>

      <!-- Tab: RULES -->
      <div v-show="activeTab === 'rules'" class="space-y-6">
        <div class="bg-white rounded-lg shadow-sm border border-gray-100 p-6">
          <h3 class="text-lg font-bold text-gray-800 mb-4">
            Edit Tournament Rules
          </h3>
          <p class="text-sm text-gray-500 mb-4">
            You can use simple text. It will be displayed as-is on the public
            page.
          </p>
          <textarea
            v-model="rulesContent"
            rows="15"
            class="w-full border border-gray-300 rounded-sm p-4 font-mono text-sm focus:ring-2 focus:ring-violet-500 focus:border-violet-500"
            placeholder="Enter rules here..."
          ></textarea>
          <div class="mt-4 flex justify-end">
            <button
              @click="saveRules"
              :disabled="savingRules"
              class="px-6 py-2 bg-violet-600 text-white font-medium rounded-sm hover:bg-violet-700 disabled:opacity-50 transition-colors"
            >
              {{ savingRules ? "Saving..." : "Save Rules" }}
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- Modal: Score -->
    <ScoreModal
      :is-open="showScoreModal"
      :match="selectedMatch || {}"
      :stage="selectedMatchStage"
      @close="showScoreModal = false"
      @save="saveMatchResult"
    />
  </div>
</template>
