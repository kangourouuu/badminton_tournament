<script setup>
import TheNavbar from "../components/TheNavbar.vue";
</script>

<template>
  <div class="min-h-screen bg-gray-50 pb-12">
    <TheNavbar />

    <main class="container mx-auto px-4 py-12 flex justify-center">
      <div
        class="bg-white rounded-lg border border-purple-200 p-8 max-w-md w-full shadow-none text-center space-y-6"
      >
        <div
          class="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center mx-auto text-3xl"
        >
          📞
        </div>

        <div>
          <h2 class="text-2xl font-bold text-violet-800">Contact Us</h2>
          <p class="text-gray-500 mt-2">
            Get in touch with the tournament organizers.
          </p>
        </div>

        <div class="space-y-4 text-left">
          <div class="flex items-center gap-4 p-4 bg-gray-50 rounded-sm">
            <span class="text-xl">🏢</span>
            <div>
              <p class="text-xs text-gray-400 uppercase font-bold">Organizer</p>
              <p class="font-medium text-gray-800">Badminton Committee</p>
            </div>
          </div>

          <div class="flex items-center gap-4 p-4 bg-gray-50 rounded-sm">
            <span class="text-xl">📧</span>
            <div>
              <p class="text-xs text-gray-400 uppercase font-bold">Email</p>
              <p class="font-medium text-gray-800">contact@badminton.com</p>
            </div>
          </div>

          <div class="flex items-center gap-4 p-4 bg-gray-50 rounded-sm">
            <span class="text-xl">📱</span>
            <div>
              <p class="text-xs text-gray-400 uppercase font-bold">Hotline</p>
              <p class="font-medium text-gray-800">+84 90 123 4567</p>
            </div>
          </div>
        </div>

        <div class="pt-4 border-t border-purple-50">
          <p class="text-xs text-gray-400">
            Available 8:00 AM - 6:00 PM (Mon-Fri)
          </p>
        </div>
      </div>
    </main>
  </div>
</template>
