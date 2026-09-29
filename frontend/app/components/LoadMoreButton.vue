<script setup lang="ts">
const emit = defineEmits<{
  loadMore: []
}>()
</script>

<template>
  <div class="flex justify-center mt-10">
    <UiAppButton
      variant="secondary"
      class="h-[52px] gap-4 px-10"
      @click="emit('loadMore')"
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
      >
        <path d="M12 5v14" />
        <path d="m6 13 6 6 6-6" />
      </svg>
    </UiAppButton>
  </div>
</template>