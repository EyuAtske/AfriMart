<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'default' | 'brand' | 'dark' | 'success' | 'warning' | 'muted'
  }>(),
  {
    variant: 'default'
  }
)
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium uppercase tracking-[0.12em]"
    :class="{
      'bg-[#eee8df] text-[#211f1d]': variant === 'default',
      'bg-[#806344] text-white': variant === 'brand',
      'bg-[#211f1d] text-white': variant === 'dark',
      'bg-green-100 text-green-800 border border-green-300': variant === 'success',
      'bg-amber-100 text-amber-800 border border-amber-300': variant === 'warning',
      'bg-[#92877b] text-white': variant === 'muted'
    }"
  >
    <slot />
  </span>
</template>
