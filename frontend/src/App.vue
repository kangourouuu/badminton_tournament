<script setup>
import { ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();
const isLoggedIn = ref(!!localStorage.getItem("token"));

watch(route, () => {
  isLoggedIn.value = !!localStorage.getItem("token");
});

const logout = () => {
  localStorage.removeItem("token");
  isLoggedIn.value = false;
  router.push("/");
};
</script>

<template>
  <div class="min-h-screen bg-[#FDFBFF] font-sans text-slate-900">
    <header class="border-b border-purple-200 bg-white px-6 py-4">
      <div class="mx-auto max-w-7xl flex items-center justify-between">
        <div
          class="flex items-center gap-2"
          @click="router.push('/')"
          role="button"
        >
          <!-- Optional Logo Icon could go here -->
          <h1 class="text-xl font-bold tracking-tight text-purple-900">
            Badminton Tournament
          </h1>
        </div>

        <nav>
          <div v-if="isLoggedIn" class="flex items-center gap-4">
            <span class="text-sm font-medium text-purple-900">Admin</span>
            <button
              @click="logout"
              class="text-sm font-medium text-slate-500 hover:text-red-600 transition-colors"
            >
              Logout
            </button>
          </div>
          <div v-else>
            <router-link
              to="/admin/login"
              class="text-sm font-medium text-violet-600 hover:text-violet-700 transition-colors"
            >
              Admin Login
            </router-link>
          </div>
        </nav>
      </div>
    </header>

    <main class="mx-auto max-w-7xl p-6">
      <router-view />
    </main>
  </div>
</template>
