<script setup lang="ts">
const props = defineProps<{
  id: number
  shop: string
  name: string
  price: string
  rating: string | number
  image: string
  stock?: number
}>()

const emit = defineEmits<{
  addToCart: [id: number]
}>()
</script>

<template>
  <article class="group overflow-hidden rounded-xl border border-[#e2d8c8] bg-[#eee5d5] transition duration-300 hover:-translate-y-1 hover:shadow-[0_18px_45px_rgba(33,31,29,0.10)]">
    <NuxtLink
      :to="`/products/${id}`"
      class="block"
    >
      <div class="relative overflow-hidden bg-white">
        <img
          :src="image"
          :alt="name"
          class="block h-auto w-full transition duration-700 group-hover:scale-105"
        />

        <span
          v-if="stock !== undefined"
          class="absolute left-3 top-3 rounded-full bg-[#f8f4ec]/95 px-3 py-1 text-[10px] font-medium uppercase tracking-[0.12em] text-[#5d4b37]"
        >
          {{ stock }} left
        </span>
      </div>
    </NuxtLink>

    <div class="px-4 py-4 sm:px-5">
      <NuxtLink
        :to="`/shops/${shop.toLowerCase().replace(/\s+/g, '-')}`"
        class="inline-block text-[10px] uppercase tracking-[0.14em] text-[#8b6b46] transition hover:text-[#211f1d] hover:underline sm:text-xs"
      >
        {{ shop }}
      </NuxtLink>

      <NuxtLink
        :to="`/products/${id}`"
        class="mt-1 block text-base text-[#211f1d] transition hover:text-[#806344] sm:text-lg"
      >
        {{ name }}
      </NuxtLink>

      <div class="mt-3 flex items-center justify-between gap-2">
        <p class="text-base font-medium text-[#806344] sm:text-lg">
          {{ price }}
        </p>

        <div class="flex items-center gap-1 text-xs text-[#81715f] sm:text-sm">
          <span>*</span>
          <span>{{ rating }}</span>
        </div>
      </div>

      <button
        type="button"
        class="mt-4 h-11 w-full rounded-full border border-[#806344] text-xs font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2"
        @click="emit('addToCart', props.id)"
      >
        Add to cart
      </button>
    </div>
  </article>
</template>
