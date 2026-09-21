<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  actionLabel?: string
  actionTo?: string
}>()
</script>

<template>
  <section class="rounded-xl border border-[#d9d0c4] bg-[#faf8f4] p-10 text-center shadow-[0_20px_70px_rgba(33,31,29,0.06)]">
    <h2 class="font-serif text-3xl text-[#211f1d]">
      {{ title }}
    </h2>

    <p v-if="description" class="mt-3 text-sm leading-6 text-[#756a60]">
      {{ description }}
    </p>

    <UiAppButton
      v-if="actionLabel && actionTo"
      :to="actionTo"
      variant="secondary"
      class="mt-6"
    >
      {{ actionLabel }}
    </UiAppButton>
  </section>
</template>
