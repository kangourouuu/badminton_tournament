<script setup>
import { computed } from "vue";

const props = defineProps(["url"]);
const emit = defineEmits(["close"]);

const embedUrl = computed(() => {
  if (!props.url) return "";
  // Convert standard YouTube URL to embed
  // Very basic parsing
  let videoId = "";
  if (props.url.includes("youtube.com/watch?v=")) {
    videoId = props.url.split("v=")[1].split("&")[0];
  } else if (props.url.includes("youtu.be/")) {
    videoId = props.url.split("youtu.be/")[1];
  }

  if (videoId) return `https://www.youtube.com/embed/${videoId}?autoplay=1`;
  return props.url; // Fallback
});
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm"
    @click.self="$emit('close')"
  >
    <div
      class="bg-black w-full max-w-4xl aspect-video rounded-sm overflow-hidden relative shadow-2xl border border-gray-800"
    >
      <button
        @click="$emit('close')"
        class="absolute top-4 right-4 text-white text-2xl hover:text-primary z-10"
      >
        &times;
      </button>
      <iframe
        v-if="embedUrl"
        :src="embedUrl"
        class="w-full h-full"
        frameborder="0"
        allow="
          accelerometer;
          autoplay;
          encrypted-media;
          gyroscope;
          picture-in-picture;
        "
        allowfullscreen
      >
      </iframe>
      <div v-else class="flex items-center justify-center h-full text-white">
        Invalid Video URL
      </div>
    </div>
  </div>
</template>
