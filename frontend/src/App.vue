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
  <component :is="currentView" />
</template>
