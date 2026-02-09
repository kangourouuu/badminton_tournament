<script setup>
import { ref, onMounted } from "vue";
import api from "../services/api";
import TheNavbar from "../components/TheNavbar.vue";

const rules = ref("");
const loading = ref(true);

onMounted(async () => {
  try {
    const res = await api.get("/public/rules");
    rules.value = res.data.content;
  } catch (err) {
    console.error("Failed to fetch rules", err);
    rules.value = "Failed to load rules.";
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="min-h-screen bg-gray-50 pb-12">
    <TheNavbar />

    <main class="container mx-auto px-4 py-12">
      <div class="max-w-3xl mx-auto space-y-8">
        <div class="text-center">
          <h1 class="text-3xl font-bold text-violet-800">Tournament Rules</h1>
          <p class="text-gray-500 mt-2">Official Playbook & Guidelines</p>
        </div>

        <div
          class="bg-white rounded-lg border border-purple-200 p-8 min-h-[400px]"
        >
          <div v-if="loading" class="text-center text-gray-400 py-12">
            Loading rules...
          </div>
          <div
            v-else
            class="prose prose-purple max-w-none whitespace-pre-wrap text-gray-700"
          >
            {{ rules }}
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
