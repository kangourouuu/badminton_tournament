<script setup>
import { ref, onMounted, computed } from "vue";
import { marked } from "marked";
import api from "../services/api";
import TheNavbar from "../components/TheNavbar.vue";

const rules = ref("");
const loading = ref(true);

const compiledRules = computed(() => {
  return marked(rules.value);
});

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
  <div class="min-h-screen bg-gray-50 relative overflow-hidden">
    <!-- Navbar -->
    <div class="relative z-20">
      <TheNavbar />
    </div>

    <!-- Main Content -->
    <main class="relative z-10 container mx-auto px-4 py-12">
      <header class="text-center mb-12">
        <h1
          class="text-4xl md:text-5xl font-extrabold text-gray-900 tracking-tight mb-2 uppercase"
        >
          Tournament Rules
        </h1>
        <p
          class="text-violet-600 font-medium tracking-widest uppercase text-sm"
        >
          Official Playbook & Guidelines
        </p>
      </header>

      <div class="max-w-4xl mx-auto">
        <div v-if="loading" class="text-center text-gray-400 py-12">
          <div
            class="animate-spin h-8 w-8 border-4 border-violet-500 border-t-transparent rounded-full mx-auto mb-4"
          ></div>
          <p class="uppercase tracking-widest text-xs font-semibold">
            Loading Playbook...
          </p>
        </div>

        <div
          v-else
          class="bg-white border border-gray-200 rounded-lg shadow-sm p-8 md:p-12 prose prose-violet max-w-none"
          v-html="compiledRules"
        ></div>
      </div>
    </main>
  </div>
</template>
