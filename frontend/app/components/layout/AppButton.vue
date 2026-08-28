<script setup lang="ts">
withDefaults(
  defineProps<{
    type?: 'button' | 'submit'
    variant?: 'primary' | 'secondary' | 'danger'
    disabled?: boolean
  }>(),
  {
    type: 'button',
    variant: 'primary',
    disabled: false
  }
)
</script>

<template>
  <button
    :type="type"
    :disabled="disabled"
    class="h-10 rounded-full px-5 text-xs font-medium uppercase tracking-[0.12em] transition focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
    :class="{
      'bg-[#211f1d] text-white hover:bg-[#3b3733]':
        variant === 'primary',

      'border border-[#d9d0c4] bg-[#faf8f4] text-[#211f1d] hover:bg-[#f0ebe4]':
        variant === 'secondary',

      'bg-red-600 text-white hover:bg-red-700':
        variant === 'danger'
    }"
  >
    <slot />
  </button>
</template>