<script setup lang="ts">
withDefaults(
  defineProps<{
    type?: 'button' | 'submit'
    variant?: 'primary' | 'secondary' | 'ghost'
    size?: 'default' | 'small'
    disabled?: boolean
    to?: string
  }>(),
  {
    type: 'button',
    variant: 'primary',
    size: 'default',
    disabled: false
  }
)
</script>

<template>
  <UiAppButton
    :type="type"
    :variant="variant"
    :size="size"
    :disabled="disabled"
    :to="to"
  >
    <slot />
  </UiAppButton>
</template>