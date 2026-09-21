<script setup lang="ts">
defineProps<{
  title: string
  value: string | number
  subtitle: string
  trend?: string
  trendPositive?: boolean
}>()
</script>

<template>
  <div class="flex flex-col justify-between rounded-xl border border-[#ebe5dd] bg-white p-5 shadow-xs transition hover:border-[#dfd7cc] hover:shadow-sm">
    <div>
      <p class="text-[11px] font-semibold tracking-[0.14em] text-[#8a7a6c] uppercase">
        {{ title }}
      </p>

      <p class="mt-2.5 font-serif text-3xl font-medium tracking-tight text-[#211f1d]">
        {{ value }}
      </p>

      <p class="mt-0.5 text-xs text-[#756a60]">
        {{ subtitle }}
      </p>
    </div>

    <div v-if="trend" class="mt-3.5 flex items-center text-xs font-medium text-[#2d7d46]">
      <span>{{ trend }}</span>
    </div>
  </div>
</template>
