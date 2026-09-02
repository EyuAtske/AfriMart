<script setup lang="ts">
const emit = defineEmits<{
  loadMore: []
}>()
</script>

<template>
  <div class="flex justify-center mt-10">
    <button
      type="button"
      @click="emit('loadMore')"
      class="group inline-flex items-center gap-4 rounded-full border border-[#806344] px-10 py-4 text-sm tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 focus:ring-offset-[#f5f1e9]"
    >
      <span>LOAD MORE</span>

      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        class="transition-transform duration-300 group-hover:translate-y-1"
      >
        <path d="M12 5v14" />
        <path d="m6 13 6 6 6-6" />
      </svg>
    </button>
  </div>
</template>