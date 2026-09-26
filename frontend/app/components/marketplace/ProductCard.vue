<script setup lang="ts">
import type { ProductMedia } from '~/types/product'

const props = defineProps<{
  id: number | string
  shop: string
  name: string
  price: string
  rating: string | number
  image: string
  stock?: number
  media?: ProductMedia[]
}>()

const emit = defineEmits<{
  addToCart: [id: number | string]
}>()

const { flyToCart } = useFlyToCart()
const { isOwnProduct } = useMarketplace()
const imageEl = ref<HTMLImageElement | null>(null)

const isSelfProduct = computed(() => isOwnProduct({ id: props.id, shop: props.shop } as any))
const isSoldOut = computed(() => props.stock !== undefined && props.stock <= 0)

const displayImage = computed(() => {
  if (props.media?.length) {
    const primary = props.media.find(m => m.isPrimary)
    return primary?.url || props.media[0]?.url || props.image
  }
  return props.image
})

const onImageError = (event: Event) => {
  const target = event.target as HTMLImageElement
  if (target && !target.src.endsWith('/images/shop.jpg')) {
    target.src = '/images/shop.jpg'
  }
}

const handleAddToCart = () => {
  if (isSoldOut.value || isSelfProduct.value) return
  flyToCart(imageEl.value, displayImage.value)
  emit('addToCart', props.id)
}
</script>

<template>
  <article class="group overflow-hidden rounded-xl border border-[#d9d0c4] bg-[#eee8df] transition duration-300 hover:-translate-y-1 hover:shadow-[0_18px_45px_rgba(33,31,29,0.10)]">
    <NuxtLink
      :to="`/products/${id}`"
      class="block"
    >
      <div class="relative overflow-hidden bg-white">
        <img
          ref="imageEl"
          :src="displayImage"
          :alt="name"
          @error="onImageError"
          class="block h-auto w-full transition duration-700 group-hover:scale-105"
        />

        <!-- SOLD OUT Badge -->
        <span
          v-if="isSoldOut"
          class="absolute left-3 top-3 rounded-full bg-red-700 text-white px-3 py-1 text-xs font-bold uppercase tracking-wider shadow"
        >
          Sold Out
        </span>

        <!-- YOUR PRODUCT Badge -->
        <span
          v-else-if="isSelfProduct"
          class="absolute left-3 top-3 rounded-full bg-[#806344] text-white px-3 py-1 text-xs font-semibold uppercase tracking-wider shadow"
        >
          Your Product
        </span>

        <!-- Stock Badge -->
        <span
          v-else-if="stock !== undefined"
          class="absolute left-3 top-3 rounded-full bg-[#faf8f4]/95 px-3 py-1 text-xs font-medium uppercase tracking-[0.12em] text-[#5d4b37]"
        >
          {{ stock }} left
        </span>
      </div>
    </NuxtLink>

    <div class="p-3 sm:p-5">
      <NuxtLink
        :to="`/shops/${shop.toLowerCase().replace(/\s+/g, '-')}`"
        class="inline-block text-[10px] sm:text-xs uppercase tracking-[0.14em] text-[#806344] transition hover:text-[#211f1d] hover:underline truncate max-w-full"
      >
        {{ shop }}
      </NuxtLink>

      <NuxtLink
        :to="`/products/${id}`"
        class="mt-1 block font-serif text-sm sm:text-lg text-[#211f1d] transition hover:text-[#806344] line-clamp-2"
      >
        {{ name }}
      </NuxtLink>

      <div class="mt-2 sm:mt-3 flex items-center justify-between gap-1">
        <p class="text-sm sm:text-lg font-medium text-[#806344]">
          {{ price }}
        </p>

        <div class="flex items-center gap-1 text-[11px] sm:text-sm text-[#756a60]">
          <span>★</span>
          <span>{{ rating }}</span>
        </div>
      </div>

      <UiAppButton
        class="mt-3 sm:mt-4 w-full text-xs"
        size="small"
        :disabled="isSoldOut || isSelfProduct"
        @click="handleAddToCart"
      >
        <template v-if="isSoldOut">Sold Out</template>
        <template v-else-if="isSelfProduct">Your Product</template>
        <template v-else>Add to cart</template>
      </UiAppButton>
    </div>
  </article>
</template>
