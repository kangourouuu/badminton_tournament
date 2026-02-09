<script setup>
import { computed, ref } from "vue";

const props = defineProps({
  match: {
    type: Object,
    required: true,
  },
  isAdmin: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["click", "play-video"]);

const showCopiedToast = ref(false);

// Determine status
const status = computed(() => {
  if (props.match.winner_id) return "finished";
  if (props.match.team_a_id && props.match.team_b_id) return "ready"; // or live if we had that state
  return "scheduled";
});

const cardClasses = computed(() => {
  const base = "border-2";
  if (status.value === "finished") return `${base} border-gray-200 opacity-95`;
  if (status.value === "ready")
    return `${base} border-violet-500 ring-2 ring-violet-100 shadow-md transform scale-[1.02]`;
  return `${base} border-dashed border-gray-200 bg-gray-50/50`;
});

const headerClasses = computed(() => {
  if (status.value === "finished")
    return "bg-gray-100 text-gray-500 border-gray-200";
  if (status.value === "ready")
    return "bg-violet-600 text-white border-violet-600";
  return "bg-gray-50 text-gray-400 border-gray-200";
});

const getTeamName = (teamId) => {
  if (!teamId) return "TBD";
  if (props.match.team_a?.id === teamId) return props.match.team_a.name;
  if (props.match.team_b?.id === teamId) return props.match.team_b.name;
  // Fallback if relation not loaded but ID exists (shouldn't happen with eager load)
  return "Team " + teamId.substring(0, 4);
};

const isWinner = (teamId) => {
  return status.value === "finished" && props.match.winner_id === teamId;
};

const handleClick = async () => {
  if (props.isAdmin) {
    emit("click", props.match);
  } else if (props.match.video_url) {
    // Copy to clipboard
    try {
      await navigator.clipboard.writeText(props.match.video_url);
      showCopiedToast.value = true;
      setTimeout(() => (showCopiedToast.value = false), 2000);
      // Also open behavior? Requirement says "Copy... AND show Toast".
      // "Alternatively... Play Icon... opens link".
      // Let's do both: default click copies.
      // We already have a play icon overlay that opens link on click?
      // The existing overlay does NOTHING but show an icon.
    } catch (err) {
      console.error("Failed to copy", err);
      // Fallback or just ignore
    }
  }
};

const openVideo = (e) => {
  e.stopPropagation(); // Prevent card click
  if (props.match.video_url) {
    window.open(props.match.video_url, "_blank");
  }
};
</script>

<template>
  <div
    :class="cardClasses"
    @click="handleClick"
    class="cursor-pointer hover:shadow-md group transition-all duration-200 bg-white rounded-sm border-2 overflow-hidden relative"
  >
    <!-- Header Strip -->
    <div
      class="flex justify-between items-center px-3 py-1.5 border-b"
      :class="headerClasses"
    >
      <span class="text-[10px] font-bold uppercase tracking-wider opacity-80">{{
        match.label
      }}</span>
      <span
        v-if="match.score"
        class="text-[10px] font-bold uppercase bg-white/20 px-1.5 rounded"
        >{{ status === "finished" ? "Finished" : "Live" }}</span
      >
    </div>

    <div class="p-3 space-y-3">
      <!-- Team A -->
      <div
        class="flex justify-between items-center p-2 rounded-sm transition-colors"
        :class="isWinner(match.team_a_id) ? 'bg-violet-50' : ''"
      >
        <div class="flex items-center gap-2">
          <!-- Placeholder Avatar -->
          <div
            class="w-6 h-6 rounded-full bg-gray-100 flex items-center justify-center text-xs text-gray-400 font-bold"
          >
            {{ match.team_a?.name?.charAt(0) || "?" }}
          </div>
          <span
            :class="{
              'font-medium text-sm text-gray-700': true,
              'font-bold text-violet-900': isWinner(match.team_a_id),
            }"
          >
            {{ match.team_a?.name || "TBD" }}
          </span>
        </div>
        <span
          v-if="isWinner(match.team_a_id)"
          class="text-xs font-bold text-violet-600"
          >WIN</span
        >
      </div>

      <!-- Score display in middle if needed, or simple divider -->
      <!-- For 'Tech-Flat Pro', let's format the score distinctly if it exists -->
      <div v-if="match.score" class="text-center py-1">
        <span class="text-2xl font-black text-gray-800 tracking-tighter">{{
          match.score
        }}</span>
      </div>
      <div v-else class="text-center text-xs text-gray-300 py-1 font-mono">
        VS
      </div>

      <!-- Team B -->
      <div
        class="flex justify-between items-center p-2 rounded-sm transition-colors"
        :class="isWinner(match.team_b_id) ? 'bg-violet-50' : ''"
      >
        <div class="flex items-center gap-2">
          <div
            class="w-6 h-6 rounded-full bg-gray-100 flex items-center justify-center text-xs text-gray-400 font-bold"
          >
            {{ match.team_b?.name?.charAt(0) || "?" }}
          </div>
          <span
            :class="{
              'font-medium text-sm text-gray-700': true,
              'font-bold text-violet-900': isWinner(match.team_b_id),
            }"
          >
            {{ match.team_b?.name || "TBD" }}
          </span>
        </div>
        <span
          v-if="isWinner(match.team_b_id)"
          class="text-xs font-bold text-violet-600"
          >WIN</span
        >
      </div>
    </div>

    <!-- Video Overlay (Play Button) -->
    <div
      v-if="match.video_url && !isAdmin"
      @click="openVideo"
      class="absolute inset-0 bg-black/5 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity z-10"
    >
      <div
        class="bg-white rounded-full p-3 shadow-lg text-violet-600 transform scale-110 hover:scale-125 transition-transform"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="w-6 h-6"
        >
          <path
            fill-rule="evenodd"
            d="M4.5 5.653c0-1.426 1.529-2.33 2.779-1.643l11.54 6.348c1.295.712 1.295 2.573 0 3.285L7.28 19.991c-1.25.687-2.779-.217-2.779-1.643V5.653z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
    </div>

    <!-- Copied Toast -->
    <div
      v-if="showCopiedToast"
      class="absolute inset-0 flex items-end justify-center pb-2 pointer-events-none z-20"
    >
      <div
        class="bg-gray-800 text-white text-xs px-2 py-1 rounded shadow-md opacity-90"
      >
        Link copied!
      </div>
    </div>
  </div>
</template>
