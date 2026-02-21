<script setup>
import { ref, onMounted, nextTick, onUnmounted, computed, watch } from "vue";
import api from "../services/api";
import MatchNode from "../components/MatchNode.vue";
import TheNavbar from "../components/TheNavbar.vue";
import ScoreModal from "../components/ScoreModal.vue";
import MatchDetailsModal from "../components/MatchDetailsModal.vue";

const props = defineProps({
  isAdmin: { type: Boolean, default: false },
});

const groups = ref([]);
const loading = ref(true);
const zoom = ref(0.7);
const selectedMatch = ref(null);
const isModalOpen = ref(false);

const mesoneerGroups = computed(() =>
  groups.value.filter((g) => g.pool === "Mesoneer" && g.name !== "KNOCKOUT"),
);
const labGroups = computed(() =>
  groups.value.filter((g) => g.pool === "Lab" && g.name !== "KNOCKOUT"),
);

const knockoutMatches = computed(() => {
  const kg = groups.value.find((g) => g.name === "KNOCKOUT");
  return kg && kg.matches ? kg.matches : [];
});

const sf1 = computed(
  () =>
    knockoutMatches.value.find((m) => m.label === "SF1") || {
      id: "ghost-sf1",
      label: "SF1",
      isGhost: true,
    },
);
const sf2 = computed(
  () =>
    knockoutMatches.value.find((m) => m.label === "SF2") || {
      id: "ghost-sf2",
      label: "SF2",
      isGhost: true,
    },
);
const final = computed(
  () =>
    knockoutMatches.value.find((m) => m.label === "Final") || {
      id: "ghost-final",
      label: "Final",
      isGhost: true,
    },
);
const bronze = computed(
  () =>
    knockoutMatches.value.find((m) => m.label === "Bronze") || {
      id: "ghost-bronze",
      label: "Bronze",
      isGhost: true,
    },
);

const fetchData = async () => {
  try {
    const response = await api.get("/groups");
    groups.value = response.data || [];
  } catch (err) {
    console.error("Failed to fetch data", err);
  } finally {
    loading.value = false;
  }
};

const isMatchDetailsOpen = ref(false);
const selectedMatchDetails = ref(null);

const getMatch = (matches, label) => matches?.find((m) => m.label === label);

const handleMatchClick = (match) => {
  if (props.isAdmin) {
    selectedMatch.value = match;
    isModalOpen.value = true;
  } else {
    // Open read-only modal
    selectedMatchDetails.value = match;
    isMatchDetailsOpen.value = true;
  }
};

const saveMatch = async (data) => {
  try {
    await api.post(`/matches/${data.id}`, data);
    isModalOpen.value = false;
    fetchData();
  } catch (err) {
    alert("Failed to save: " + err.message);
  }
};

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="min-h-screen bg-dot-pattern font-outfit text-gray-900 pb-20">
    <TheNavbar />

    <!-- View Controls -->
    <div
      class="fixed bottom-10 left-10 z-50 flex items-center gap-4 bg-white/80 backdrop-blur-md border border-gray-100 px-4 py-2 rounded-full shadow-sm"
    >
      <button
        @click="zoom = Math.max(0.3, zoom - 0.1)"
        class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 transition-colors text-lg font-light"
      >
        -
      </button>
      <span
        class="text-[10px] font-black uppercase tracking-widest text-gray-400 w-12 text-center"
        >{{ Math.round(zoom * 100) }}%</span
      >
      <button
        @click="zoom = Math.min(1.5, zoom + 0.1)"
        class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 transition-colors text-lg font-light"
      >
        +
      </button>
    </div>

    <main
      class="origin-top-left transition-transform duration-300 px-20 pt-32"
      :style="{ transform: `scale(${zoom})`, width: `${100 / zoom}%` }"
    >
      <div
        v-if="loading"
        class="text-center py-40 animate-pulse text-xs font-black tracking-[0.5em] text-gray-300 uppercase"
      >
        Initializing Tournament Grid
      </div>

      <div v-else class="relative flex justify-between gap-40 min-w-[2000px]">
        <!-- Column: Session A -->
        <div class="flex flex-col gap-32 z-10 justify-center">
          <div
            v-for="g in mesoneerGroups"
            :key="g.id"
            class="space-y-12 pb-20 border-b border-gray-50 last:border-0"
          >
            <div
              class="text-[16px] font-bold text-black uppercase tracking-widest"
            >
              {{ g.name }}
            </div>
            <div class="flex gap-16">
              <!-- Round 1 -->
              <div class="flex flex-col justify-between py-4">
                <MatchNode
                  v-if="getMatch(g.matches, 'M1')"
                  :match="getMatch(g.matches, 'M1')"
                  @click="handleMatchClick"
                />
                <MatchNode
                  v-if="getMatch(g.matches, 'M2')"
                  :match="getMatch(g.matches, 'M2')"
                  @click="handleMatchClick"
                />
              </div>
              <!-- Round 2 -->
              <div class="flex flex-col justify-between py-4">
                <MatchNode
                  v-if="getMatch(g.matches, 'Winners')"
                  :match="getMatch(g.matches, 'Winners')"
                  @click="handleMatchClick"
                />
                <MatchNode
                  v-if="getMatch(g.matches, 'Losers')"
                  :match="getMatch(g.matches, 'Losers')"
                  @click="handleMatchClick"
                />
              </div>
              <!-- Decider -->
              <div class="flex flex-col justify-center">
                <MatchNode
                  v-if="getMatch(g.matches, 'Decider')"
                  :match="getMatch(g.matches, 'Decider')"
                  @click="handleMatchClick"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Column: Champions (Center) -->
        <div
          class="flex flex-col items-center justify-center gap-8 z-10 pt-20 w-full"
        >
          <div class="text-center mb-10">
            <h1
              class="text-6xl font-black text-gray-900 tracking-tighter italic uppercase"
            >
              Champions Stage
            </h1>
            <div class="h-1 w-20 bg-violet-600 mx-auto mt-4"></div>
          </div>

          <div class="flex flex-col items-center gap-12 w-full">
            <!-- SF1 -->
            <div class="flex flex-col items-center gap-2">
              <div
                class="text-[20px] font-black text-black uppercase tracking-[0.5em] text-center"
              >
                Semi Final 1
              </div>
              <MatchNode
                :match="sf1"
                @click="!sf1.isGhost && handleMatchClick($event)"
                class="scale-110"
                :class="{ 'opacity-75': sf1.isGhost }"
              />
            </div>

            <!-- Final -->
            <div class="flex flex-col items-center gap-2">
              <div
                class="text-[20px] font-black text-gray-900 uppercase tracking-[1em] text-center mb-1"
              >
                Grand Final
              </div>
              <MatchNode
                :match="final"
                @click="!final.isGhost && handleMatchClick($event)"
                class="scale-[1.6] shadow-2xl border-violet-100"
                :class="{ 'opacity-75': final.isGhost }"
              />
            </div>

            <!-- SF2 -->
            <div class="flex flex-col items-center gap-2">
              <div
                class="text-[20px] font-black text-black uppercase tracking-[0.5em] text-center"
              >
                Semi Final 2
              </div>
              <MatchNode
                :match="sf2"
                @click="!sf2.isGhost && handleMatchClick($event)"
                class="scale-110"
                :class="{ 'opacity-75': sf2.isGhost }"
              />
            </div>

            <!-- Bronze -->
            <div class="mt-8 opacity-50 flex flex-col items-center gap-2">
              <div
                class="text-[20px] font-black text-black uppercase tracking-widest text-center"
              >
                Bronze Match
              </div>
              <MatchNode
                :match="bronze"
                @click="!bronze.isGhost && handleMatchClick($event)"
                class="scale-90"
                :class="{ 'opacity-75': bronze.isGhost }"
              />
            </div>
          </div>
        </div>

        <!-- Column: Session B -->
        <div class="flex flex-col gap-32 z-10 justify-center">
          <div
            v-for="g in labGroups"
            :key="g.id"
            class="space-y-12 pb-20 border-b border-gray-50 last:border-0 items-end flex flex-col"
          >
            <div
              class="text-[16px] font-bold text-black uppercase tracking-widest text-right"
            >
              {{ g.name }}
            </div>
            <div class="flex gap-16 flex-row-reverse">
              <!-- Round 1 -->
              <div class="flex flex-col justify-between py-4">
                <MatchNode
                  v-if="getMatch(g.matches, 'M1')"
                  :match="getMatch(g.matches, 'M1')"
                  @click="handleMatchClick"
                />
                <MatchNode
                  v-if="getMatch(g.matches, 'M2')"
                  :match="getMatch(g.matches, 'M2')"
                  @click="handleMatchClick"
                />
              </div>
              <!-- Round 2 -->
              <div class="flex flex-col justify-between py-4">
                <MatchNode
                  v-if="getMatch(g.matches, 'Winners')"
                  :match="getMatch(g.matches, 'Winners')"
                  @click="handleMatchClick"
                />
                <MatchNode
                  v-if="getMatch(g.matches, 'Losers')"
                  :match="getMatch(g.matches, 'Losers')"
                  @click="handleMatchClick"
                />
              </div>
              <!-- Decider -->
              <div class="flex flex-col justify-center">
                <MatchNode
                  v-if="getMatch(g.matches, 'Decider')"
                  :match="getMatch(g.matches, 'Decider')"
                  @click="handleMatchClick"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <ScoreModal
      :is-open="isModalOpen"
      :match="selectedMatch"
      :is-admin="true"
      @close="isModalOpen = false"
      @save="saveMatch"
    />

    <MatchDetailsModal
      v-if="isMatchDetailsOpen && selectedMatchDetails"
      :is-open="isMatchDetailsOpen"
      :match="selectedMatchDetails"
      @close="isMatchDetailsOpen = false"
    />
  </div>
</template>

<style scoped>
/* No additional styles needed, using Tailwind */
</style>
