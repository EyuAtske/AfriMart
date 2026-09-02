<script setup lang="ts">
import SkeletonCard from './SkeletonCard.vue'

withDefaults(
  defineProps<{
    count?: number
  }>(),
  {
    count: 6
  }
)
</script>

<template>
  <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
    <SkeletonCard v-for="i in count" :key="i" />
  </div>
</template>
