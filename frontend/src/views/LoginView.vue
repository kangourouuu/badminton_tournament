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
    const response = await api.login(password.value);
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
  <div class="flex min-h-[80vh] items-center justify-center">
    <div
      class="w-full max-w-md rounded-sm border border-purple-200 bg-white p-8"
    >
      <h2 class="mb-6 text-center text-2xl font-bold text-purple-900">
        Admin Login
      </h2>

      <form @submit.prevent="handleLogin" class="flex flex-col gap-4">
        <div>
          <label class="mb-1 block text-sm font-medium text-slate-700"
            >Password</label
          >
          <input
            v-model="password"
            type="password"
            class="w-full rounded-sm border border-purple-200 p-2 focus:border-violet-600 focus:outline-none"
            placeholder="Enter admin password"
            required
          />
        </div>

        <div v-if="error" class="text-sm text-red-600">
          {{ error }}
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="rounded-sm bg-violet-600 px-4 py-2 font-bold text-white transition-colors hover:bg-violet-700 disabled:opacity-50"
        >
          {{ loading ? "Logging in..." : "Login" }}
        </button>
      </form>
    </div>
  </div>
</template>
