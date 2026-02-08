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
  <div class="min-h-screen bg-background">
    <div
      class="bg-primary text-white py-12 px-6 text-center mb-8 relative overflow-hidden"
    >
      <h1 class="text-4xl font-bold font-sans tracking-tight relative z-10">
        Badminton Tournament
      </h1>
      <div class="flex justify-center gap-4 mt-8 relative z-10">
        <button
          @click="switchPool('Mesoneer')"
          :class="[
            'px-6 py-2 rounded-full font-bold transition-all',
            activePool === 'Mesoneer'
              ? 'bg-white text-primary'
              : 'bg-primary/50 text-white hover:bg-white/20',
          ]"
        >
          Mesoneer Pool
        </button>
        <button
          @click="switchPool('Lab')"
          :class="[
            'px-6 py-2 rounded-full font-bold transition-all',
            activePool === 'Lab'
              ? 'bg-white text-primary'
              : 'bg-primary/50 text-white hover:bg-white/20',
          ]"
        >
          Lab Pool
        </button>
      </div>

      <!-- Decoration -->
      <div
        class="absolute top-0 opacity-10 w-full h-full left-0 pointer-events-none"
      >
        <div
          class="absolute -right-20 -top-20 w-96 h-96 bg-white rounded-full"
        ></div>
        <div
          class="absolute -left-20 bottom-0 w-64 h-64 bg-white rounded-full"
        ></div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-6 pb-20">
      <div v-if="matches.length === 0" class="text-center py-20 text-gray-400">
        Matches are being arranged. Stay tuned!
      </div>
      <div v-else class="flex flex-wrap gap-12 justify-center items-center">
        <div class="flex flex-col gap-12">
          <BracketNode
            v-for="m in matches.slice(0, 2)"
            :key="m.id"
            :match="m"
            @open-video="(url) => (videoUrl = url)"
          />
        </div>
        <div class="flex flex-col gap-12">
          <BracketNode
            v-for="m in matches.slice(2, 4)"
            :key="m.id"
            :match="m"
            @open-video="(url) => (videoUrl = url)"
          />
        </div>
        <div>
          <BracketNode
            v-if="matches[4]"
            :match="matches[4]"
            @open-video="(url) => (videoUrl = url)"
          />
        </div>
      </div>
    </div>

    <VideoModal v-if="videoUrl" :url="videoUrl" @close="videoUrl = null" />
  </div>
</template>
