<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import api from "../services/api";

const password = ref("");
const error = ref("");
const loading = ref(false);
const router = useRouter();

const handleLogin = async () => {
  loading.value = true;
  error.value = "";
  try {
    const response = await api.post("/auth/login", {
      password: password.value,
    });
    localStorage.setItem("token", response.data.token);
    router.push("/admin");
  } catch (err) {
    error.value = "Invalid password";
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-50">
    <div
      class="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-sm border border-purple-100"
    >
      <h1 class="text-2xl font-bold text-center text-purple-900">
        Admin Login
      </h1>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block mb-2 text-sm font-medium text-gray-700"
            >Password</label
          >
          <input
            v-model="password"
            type="password"
            class="w-full px-4 py-2 border border-purple-200 rounded-sm focus:outline-none focus:border-violet-600 focus:ring-1 focus:ring-violet-600 transition-colors"
            placeholder="Enter admin password"
            required
          />
        </div>
        <div v-if="error" class="text-sm text-red-600 text-center">
          {{ error }}
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="w-full px-4 py-2 text-white bg-violet-600 rounded-sm hover:bg-violet-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-violet-500 disabled:opacity-50 transition-colors"
        >
          {{ loading ? "Logging in..." : "Login" }}
        </button>
      </form>
    </div>
  </div>
</template>
