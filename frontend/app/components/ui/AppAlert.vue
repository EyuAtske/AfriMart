<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'error' | 'success' | 'warning' | 'info'
  }>(),
  {
    variant: 'error'
  }
)
</script>

<template>
  <div
    role="alert"
    class="rounded-lg border px-4 py-3 text-sm font-medium"
    :class="{
      'border-red-200 bg-red-50 text-red-700': variant === 'error',
      'border-[#cfe0cc] bg-[#eef7ec] text-[#3f6a3f]': variant === 'success',
      'border-amber-200 bg-amber-50 text-amber-800': variant === 'warning',
      'border-[#d9d0c4] bg-[#eee8df] text-[#211f1d]': variant === 'info'
    }"
  >
    <slot />
  </div>
</template>
