<script setup>
import { ref, onMounted } from "vue";
import api from "../services/api";
import GSLGrid from "../components/GSLGrid.vue";
import KnockoutGrid from "../components/KnockoutGrid.vue";
import TheNavbar from "../components/TheNavbar.vue";

const groups = ref([]);
const loading = ref(true);
const activeTab = ref(null);

import { computed } from "vue"; // Import computed locally

const activeGroup = computed(() => {
  return groups.value.find((g) => g.id === activeTab.value);
});

const fetchData = async () => {
  try {
    const response = await api.get("/groups");
    groups.value = response.data || [];
    // Set default tab if not set
    if (groups.value.length > 0 && !activeTab.value) {
      activeTab.value = groups.value[0].id;
    }
  } catch (err) {
    console.error("Failed to fetch data", err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchData();
  // Optional: Polling
  setInterval(fetchData, 30000);
});
</script>

<template>
  <div class="min-h-screen pb-12 bg-tech-pattern">
    <TheNavbar />

    <main class="container mx-auto px-4 py-8 space-y-8">
      <!-- Loading State -->
      <div v-if="loading" class="text-center py-12">
        <div
          class="animate-spin h-8 w-8 border-4 border-violet-500 border-t-transparent rounded-full mx-auto mb-4"
        ></div>
        <p
          class="text-gray-500 uppercase tracking-widest text-xs font-semibold"
        >
          Loading Arena...
        </p>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="groups.length === 0"
        class="text-center py-12 text-gray-500 bg-white border border-dashed border-gray-300 rounded-sm"
      >
        No active tournaments found.
      </div>

      <!-- Content -->
      <div v-else>
        <!-- Tabs Navigation -->
        <div
          class="flex overflow-x-auto gap-2 pb-4 mb-6 border-b border-gray-100 no-scrollbar"
        >
          <button
            v-for="group in groups"
            :key="group.id"
            @click="activeTab = group.id"
            class="px-6 py-2 rounded-full text-sm font-bold uppercase tracking-wider whitespace-nowrap transition-all duration-200 border"
            :class="
              activeTab === group.id
                ? 'bg-violet-600 text-white border-violet-600 shadow-md transform scale-105'
                : 'bg-white text-gray-500 border-gray-200 hover:border-violet-300 hover:text-violet-600'
            "
          >
            {{ group.name }}
          </button>

          <!-- Future Knockout Tab Placeholder -->
          <!-- 
          <button class="px-6 py-2 rounded-full text-sm font-bold uppercase tracking-wider whitespace-nowrap bg-amber-400 text-amber-900 border border-amber-500 shadow-sm relative overflow-hidden">
             <span class="relative z-10">🏆 Knockout Stage</span>
          </button> 
          -->
        </div>

        <!-- Active Group Display -->
        <transition
          mode="out-in"
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0 translate-y-2"
        >
          <div v-if="activeGroup" :key="activeGroup.id" class="space-y-6">
            <div class="flex items-center justify-between">
              <h2
                class="text-2xl font-black text-violet-900 uppercase tracking-tighter"
              >
                {{ activeGroup.name }} Arena
              </h2>
              <div
                class="text-xs font-mono text-gray-400 bg-white px-2 py-1 rounded border border-gray-100"
              >
                Group ID: {{ activeGroup.id.substring(0, 8) }}...
              </div>
            </div>

            <!-- The Bracket -->
            <KnockoutGrid
              v-if="activeGroup.name === 'KNOCKOUT'"
              :matches="activeGroup.matches || []"
            />
            <GSLGrid v-else :matches="activeGroup.matches || []" />
          </div>
        </transition>
      </div>
    </main>
  </div>
</template>
