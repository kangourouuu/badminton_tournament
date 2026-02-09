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
  <div class="min-h-screen bg-playbook relative overflow-hidden">
    <!-- Navbar -->
    <div class="relative z-20">
      <TheNavbar />
    </div>

    <!-- Watermark (Decorative) -->
    <div
      class="fixed bottom-[-5%] right-[-5%] w-96 h-96 opacity-5 pointer-events-none z-0 rotate-12 text-purple-900"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="currentColor"
        class="w-full h-full"
      >
        <path
          fill-rule="evenodd"
          d="M5.166 2.621v.56c0 .556.405.88.57 1.013.336.27.42.337.564.912l.006.026c.266 1.066.975 3.901 3.975 5.56 1.487.822 3.141 1.054 4.719 1.054 1.577 0 3.232-.232 4.719-1.054 3-1.66 3.71-4.494 3.975-5.56l.006-.026c.144-.575.228-.642.564-.912.165-.133.57-.457.57-1.013v-.56c0-1.162-1.077-2.023-2.18-1.573C21.493 1.527 18.067 2.75 12 2.75c-6.067 0-9.493-1.223-10.655-1.693-1.103-.45-2.179.411-2.179 1.573zM12 13.25c-3.149 0-5.875-.796-7.83-1.896.792 1.346 1.954 2.457 3.33 3.197 2.46.993 4.5 2.7 4.5 5.2V22l1.6-.8c2.4-.6 3.6-2 3.6-4 0-1.4 1.6-3 3.5-4h.246c.224 0 .44-.065.656-.145-1.954 1.1-4.68 1.896-7.83 1.896H12z"
          clip-rule="evenodd"
        />
      </svg>
    </div>

    <!-- Main Content -->
    <main class="relative z-10 container mx-auto px-4 py-12">
      <header class="text-center mb-12">
        <h1
          class="text-4xl md:text-5xl font-extrabold text-purple-900 tracking-tight mb-2 uppercase"
        >
          Tournament Rules
        </h1>
        <p
          class="text-violet-600 font-medium tracking-widest uppercase text-sm"
        >
          Official Playbook & Guidelines
        </p>
      </header>

      <div class="max-w-3xl mx-auto space-y-8">
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
          class="bg-white border border-purple-200 border-t-4 border-t-violet-600 rounded-sm p-6 md:p-12 relative overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-500"
        >
          <!-- Internal Decorative numbering for the main card -->
          <div
            class="absolute top-0 right-0 p-4 opacity-10 pointer-events-none"
          >
            <span class="text-9xl font-bold text-purple-900">01</span>
          </div>

          <div
            class="prose prose-purple max-w-none whitespace-pre-wrap text-slate-600 leading-relaxed"
            v-html="compiledRules"
          ></div>
        </div>
      </div>
    </main>
  </div>
</template>
