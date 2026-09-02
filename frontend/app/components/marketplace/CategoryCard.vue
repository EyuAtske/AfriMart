<script setup lang="ts">
defineProps<{
  title: string
  image: string
  to: string
}>()
</script>

<template>
  <NuxtLink
    :to="to"
    class="group block overflow-hidden rounded-md"
  >
    <!-- Image -->
    <div class="aspect-[4/3] overflow-hidden">
      <img
        :src="image"
        :alt="title"
        class="h-full w-full object-cover object-top transition-transform duration-700 group-hover:scale-105"
      />
    </div>

    <!-- Card info -->
    <div class="bg-[#e8dcc9] px-2 py-3 sm:px-4 sm:py-5 lg:px-7 lg:py-7">
      <h2
        class="font-serif text-sm text-[#211f1d] sm:text-lg lg:text-2xl"
      >
        {{ title }}
      </h2>

      <div
        class="mt-1 flex items-center gap-1.5 text-[7px] tracking-[0.08em] text-[#806344] sm:mt-2 sm:gap-2 sm:text-[10px] lg:gap-3 lg:text-sm"
      >
        <span>EXPLORE NOW</span>

        <span
          class="block h-px w-3 bg-[#806344] sm:w-5 lg:w-7"
        />
      </div>
    </div>
  </NuxtLink>
</template>