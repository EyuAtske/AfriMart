<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  shop: string
  name: string
  price: string
  rating: string | number
  images?: string[]
  image?: string
}>()

const currentSlide = ref(0)

const productImages = computed(() => {
  if (props.images?.length) {
    return props.images
  }

  if (props.image) {
    return [props.image]
  }

  return []
})

const totalSlides = computed(() => productImages.value.length)

const nextSlide = () => {
  if (totalSlides.value <= 1) return

  currentSlide.value =
    (currentSlide.value + 1) % totalSlides.value
}

const previousSlide = () => {
  if (totalSlides.value <= 1) return

  currentSlide.value =
    (currentSlide.value - 1 + totalSlides.value) % totalSlides.value
}
</script>

<template>
  <article
    class="group overflow-hidden rounded-xl border border-[#e2d8c8] bg-[#eee5d5]"
  >
    <!-- Product image -->
    <div class="relative bg-white">
      <img
        v-if="productImages.length"
        :src="productImages[currentSlide]"
        :alt="`${name} - image ${currentSlide + 1}`"
        class="block h-auto w-full transition-opacity duration-300"
      />

      <!-- Previous arrow -->
      <button
        v-if="totalSlides > 1"
        type="button"
        aria-label="Previous image"
        class="absolute left-2 top-1/2 flex h-8 w-8 -translate-y-1/2 items-center justify-center rounded-full bg-[#f8f4ec]/90 text-[#5d4b37] shadow-sm transition hover:bg-white lg:opacity-0 lg:group-hover:opacity-100"
        @click.stop="previousSlide"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>

      <!-- Next arrow -->
      <button
        v-if="totalSlides > 1"
        type="button"
        aria-label="Next image"
        class="absolute right-2 top-1/2 flex h-8 w-8 -translate-y-1/2 items-center justify-center rounded-full bg-[#f8f4ec]/90 text-[#5d4b37] shadow-sm transition hover:bg-white lg:opacity-0 lg:group-hover:opacity-100"
        @click.stop="nextSlide"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path d="m9 18 6-6-6-6" />
        </svg>
      </button>

      <!-- Slide indicators -->
      <div
        v-if="totalSlides > 1"
        class="absolute bottom-3 left-1/2 flex -translate-x-1/2 items-center gap-1.5 rounded-full bg-[#211f1d]/35 px-2.5 py-1.5 backdrop-blur-sm"
      >
        <button
          v-for="(_, index) in productImages"
          :key="index"
          type="button"
          :aria-label="`Go to image ${index + 1}`"
          class="h-1.5 rounded-full transition-all duration-300"
          :class="
            currentSlide === index
              ? 'w-4 bg-white'
              : 'w-1.5 bg-white/60'
          "
          @click.stop="currentSlide = index"
        />
      </div>

      <!-- Add to cart -->
      <button
        type="button"
        aria-label="Add to cart"
        class="absolute bottom-3 right-3 flex h-10 w-10 items-center justify-center rounded-full border border-[#ddd3c4] bg-[#f8f4ec] text-[#5d4b37] transition hover:bg-white"
        @click.stop
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="17"
          height="17"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path d="M6 8h12l1 12H5L6 8Z" />
          <path d="M9 8a3 3 0 0 1 6 0" />
        </svg>
      </button>
    </div>

    <!-- Product information -->
    <div class="px-4 py-4 sm:px-5">
      <!-- Shop -->
      <NuxtLink
        :to="`/shops/${shop.toLowerCase().replace(/\s+/g, '-')}`"
        class="inline-block text-[10px] uppercase tracking-[0.14em] text-[#8b6b46] transition hover:text-[#211f1d] hover:underline sm:text-xs"
      >
        {{ shop }}
      </NuxtLink>

      <!-- Product name -->
      <h3 class="mt-1 text-base text-[#211f1d] sm:text-lg">
        {{ name }}
      </h3>

      <!-- Price / rating -->
      <div class="mt-2 flex items-center justify-between gap-2 sm:mt-3">
        <p class="text-base font-medium text-[#806344] sm:text-lg">
          {{ price }}
        </p>

        <div class="flex items-center gap-1 text-xs text-[#81715f] sm:text-sm">
          <span>★</span>
          <span>{{ rating }}</span>
        </div>
      </div>
    </div>
  </article>
</template>