<script setup>
import { ref } from "vue";
import { RouterLink } from "vue-router";

const isOpen = ref(false);

const toggleMenu = () => {
  isOpen.value = !isOpen.value;
};
</script>

<template>
  <nav class="bg-white border-b border-purple-100 sticky top-0 z-50">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <RouterLink to="/" class="flex items-center gap-2">
          <span class="text-2xl">🏸</span>
          <span class="text-xl font-bold text-violet-800 tracking-tight">
            Badminton Tourney
          </span>
        </RouterLink>

        <!-- Desktop Menu -->
        <div class="hidden md:flex items-center space-x-8">
          <RouterLink
            to="/"
            class="text-sm font-medium text-gray-600 hover:text-violet-600 transition-colors"
            active-class="text-violet-600 font-bold"
          >
            Bracket
          </RouterLink>
          <RouterLink
            to="/rules"
            class="text-sm font-medium text-gray-600 hover:text-violet-600 transition-colors"
            active-class="text-violet-600 font-bold"
          >
            Rules
          </RouterLink>
          <RouterLink
            to="/contact"
            class="text-sm font-medium text-gray-600 hover:text-violet-600 transition-colors"
            active-class="text-violet-600 font-bold"
          >
            Contact
          </RouterLink>
        </div>

        <!-- Mobile Menu Button -->
        <button
          @click="toggleMenu"
          class="md:hidden p-2 rounded-md text-gray-600 hover:text-violet-600 hover:bg-purple-50 focus:outline-none"
        >
          <svg
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              v-if="!isOpen"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            />
            <path
              v-else
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile Menu Dropdown -->
    <div v-if="isOpen" class="md:hidden border-t border-purple-50 bg-white">
      <div class="px-4 pt-2 pb-4 space-y-1">
        <RouterLink
          to="/"
          class="block px-3 py-3 rounded-md text-base font-medium text-gray-700 hover:text-violet-600 hover:bg-purple-50"
          active-class="bg-purple-50 text-violet-700"
          @click="isOpen = false"
        >
          Bracket
        </RouterLink>
        <RouterLink
          to="/rules"
          class="block px-3 py-3 rounded-md text-base font-medium text-gray-700 hover:text-violet-600 hover:bg-purple-50"
          active-class="bg-purple-50 text-violet-700"
          @click="isOpen = false"
        >
          Rules
        </RouterLink>
        <RouterLink
          to="/contact"
          class="block px-3 py-3 rounded-md text-base font-medium text-gray-700 hover:text-violet-600 hover:bg-purple-50"
          active-class="bg-purple-50 text-violet-700"
          @click="isOpen = false"
        >
          Contact
        </RouterLink>
      </div>
    </div>
  </nav>
</template>
