<script setup>
import { ref, onMounted } from "vue";
import api from "../services/api";
import GSLGrid from "../components/GSLGrid.vue";
import TheNavbar from "../components/TheNavbar.vue";

const groups = ref([]);
const loading = ref(true);

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

onMounted(() => {
  fetchData();
  // Optional: Polling
  setInterval(fetchData, 30000);
});
</script>

<template>
  <div class="min-h-screen pb-12 bg-gray-50">
    <TheNavbar />

    <main class="container mx-auto px-4 py-8 space-y-12">
      <div v-if="loading" class="text-center py-12 text-gray-500 animate-pulse">
        Loading tournament data...
      </div>

      <div
        v-else-if="groups.length === 0"
        class="text-center py-12 text-gray-500"
      >
        No groups active yet.
      </div>

      <div v-else v-for="group in groups" :key="group.id" class="space-y-4">
        <h2
          class="text-xl font-bold text-gray-800 border-l-4 border-violet-500 pl-3"
        >
          {{ group.name }}
        </h2>
        <GSLGrid :matches="group.matches || []" />
      </div>
    </main>
  </div>
</template>
