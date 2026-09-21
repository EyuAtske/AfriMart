<script setup lang="ts">
withDefaults(
  defineProps<{
    padding?: 'default' | 'large' | 'none'
  }>(),
  {
    padding: 'default'
  }
)
</script>

<template>
  <div
    class="rounded-xl border border-[#d9d0c4] bg-[#faf8f4] shadow-[0_20px_70px_rgba(33,31,29,0.06)]"
    :class="{
      'p-5 sm:p-6': padding === 'default',
      'p-6 sm:p-8': padding === 'large'
    }"
  >
    <slot />
  </div>
</template>
