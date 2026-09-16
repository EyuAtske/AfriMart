<script setup lang="ts">
defineProps<{
  text?: string
}>()
</script>

<template>
  <span class="relative inline-flex items-center justify-center overflow-hidden">
    <!-- Invisible label to reserve exact dimensions without shifting layout -->
    <span class="invisible select-none" aria-hidden="true">
      <slot>{{ text }}</slot>
    </span>

    <!-- Primary label: visible by default, slides upward & fades out on hover/keyboard focus -->
    <span
      class="absolute inset-0 flex items-center justify-center transition-all duration-200 ease-out group-hover:-translate-y-full group-hover:opacity-0 group-focus-visible:-translate-y-full group-focus-visible:opacity-0 motion-reduce:transition-none motion-reduce:group-hover:translate-y-0 motion-reduce:group-hover:opacity-100 motion-reduce:group-focus-visible:translate-y-0 motion-reduce:group-focus-visible:opacity-100"
    >
      <slot>{{ text }}</slot>
    </span>

    <!-- Duplicate label: enters from below & fades in on hover/keyboard focus -->
    <span
      class="absolute inset-0 flex items-center justify-center translate-y-full opacity-0 transition-all duration-200 ease-out group-hover:translate-y-0 group-hover:opacity-100 group-focus-visible:translate-y-0 group-focus-visible:opacity-100 motion-reduce:hidden"
      aria-hidden="true"
    >
      <slot>{{ text }}</slot>
    </span>
  </span>
</template>
