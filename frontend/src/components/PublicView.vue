<script setup>
import { ref, onMounted } from "vue";
import BracketNode from "./BracketNode.vue";
import VideoModal from "./VideoModal.vue";
import api from "../services/api";

const activePool = ref("Mesoneer");
const matches = ref([]);
const videoUrl = ref(null);

const loadData = async () => {
  const res = await api.getBracket(activePool.value);
  matches.value = res.data;
};

const switchPool = (pool) => {
  activePool.value = pool;
  loadData();
};

onMounted(loadData);
</script>

<template>
  <div>
    <!-- Pool Toggles (Segmented Control) -->
    <div class="mb-8 flex gap-2">
      <button
        @click="switchPool('Mesoneer')"
        :class="[
          'rounded-sm px-6 py-2 font-medium transition-colors shadow-none text-sm',
          activePool === 'Mesoneer'
            ? 'bg-violet-600 text-white'
            : 'bg-white border border-purple-200 text-purple-700 hover:bg-purple-50',
        ]"
      >
        Mesoneer Pool
      </button>
      <button
        @click="switchPool('Lab')"
        :class="[
          'rounded-sm px-6 py-2 font-medium transition-colors shadow-none text-sm',
          activePool === 'Lab'
            ? 'bg-violet-600 text-white'
            : 'bg-white border border-purple-200 text-purple-700 hover:bg-purple-50',
        ]"
      >
        Lab Pool
      </button>
    </div>

    <!-- Empty State -->
    <div
      v-if="matches.length === 0"
      class="flex h-64 flex-col items-center justify-center rounded-sm border border-purple-200 bg-white p-10 text-center"
    >
      <div class="mb-4 text-purple-300">
        <!-- SVG Icon: Calendar/Loading -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="48"
          height="48"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M21 12a9 9 0 1 1-6.219-8.56" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-purple-900">Bracket is Generating</h3>
      <p class="text-slate-500 font-light mt-1">
        Matches are being arranged by the Admin. Stay tuned!
      </p>
    </div>

    <!-- Bracket Display -->
    <div v-else class="flex flex-wrap gap-12 justify-center items-center">
      <!-- Round 1 -->
      <div class="flex flex-col gap-12">
        <BracketNode
          v-for="m in matches.slice(0, 2)"
          :key="m.id"
          :match="m"
          @open-video="(url) => (videoUrl = url)"
        />
      </div>
      <!-- Round 2 (Winners) -->
      <div class="flex flex-col gap-12">
        <BracketNode
          v-for="m in matches.slice(2, 4)"
          :key="m.id"
          :match="m"
          @open-video="(url) => (videoUrl = url)"
        />
      </div>
      <!-- Final -->
      <div>
        <BracketNode
          v-if="matches[4]"
          :match="matches[4]"
          @open-video="(url) => (videoUrl = url)"
        />
      </div>
    </div>

    <VideoModal v-if="videoUrl" :url="videoUrl" @close="videoUrl = null" />
  </div>
</template>
