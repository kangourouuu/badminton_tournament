<script setup>
import { ref, computed } from "vue";
import AdminDashboard from "./components/AdminDashboard.vue";
import PublicView from "./components/PublicView.vue";

const currentPath = ref(window.location.hash);

window.addEventListener("hashchange", () => {
  currentPath.value = window.location.hash;
});

const currentView = computed(() => {
  if (currentPath.value === "#/admin") return AdminDashboard;
  return PublicView;
});
</script>

<template>
  <div class="min-h-screen bg-[#FDFBFF] font-sans text-slate-900">
    <header class="border-b border-purple-200 bg-white px-6 py-4">
      <div class="mx-auto max-w-7xl flex items-center gap-2">
        <!-- Optional Logo Icon could go here -->
        <h1 class="text-xl font-bold tracking-tight text-purple-900">
          Badminton Tournament
        </h1>
      </div>
    </header>

    <main class="mx-auto max-w-7xl p-6">
      <component :is="currentView" />
    </main>
  </div>
</template>
