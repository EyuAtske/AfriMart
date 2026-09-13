<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    type?: 'button' | 'submit'
    variant?: 'primary' | 'secondary' | 'ghost'
    size?: 'default' | 'small'
    disabled?: boolean
    to?: string
  }>(),
  {
    type: 'button',
    variant: 'secondary',
    size: 'default',
    disabled: false
  }
)

const baseClasses =
  'group inline-flex items-center justify-center rounded-full font-medium uppercase tracking-[0.14em] transition-[colors,transform,box-shadow] duration-200 ease-out focus:outline-none focus-visible:ring-2 focus-visible:ring-[#806344] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 disabled:pointer-events-none motion-reduce:transition-none hover:-translate-y-[1px] hover:shadow-[0_2px_8px_rgba(0,0,0,0.08)] active:translate-y-0 active:scale-[0.98] active:shadow-none active:duration-100'

const sizeClasses = computed(() =>
  props.size === 'small' ? 'h-10 px-5 text-xs' : 'h-12 px-6 text-sm'
)

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'primary':
      return 'bg-[#211f1d] text-white hover:bg-[#3b3733]'
    case 'ghost':
      return 'border border-[#cfc4b5] text-[#756a60] hover:bg-[#eee8df]'
    case 'secondary':
    default:
      return 'border border-[#806344] text-[#5d4b37] hover:bg-[#806344] hover:text-white'
  }
})

const buttonClasses = computed(() =>
  [baseClasses, sizeClasses.value, variantClasses.value].join(' ')
)
</script>

<template>
  <NuxtLink v-if="to" :to="to" :class="buttonClasses">
    <slot />
  </NuxtLink>
  <button v-else :type="type" :disabled="disabled" :class="buttonClasses">
    <slot />
  </button>
</template>

