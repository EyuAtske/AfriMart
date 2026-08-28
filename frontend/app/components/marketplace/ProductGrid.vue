<script setup lang="ts">
import ProductCard from './ProductCard.vue'
import type { Product } from '~/types/product'
import { formatPrice } from '~/composables/useMarketplace'

defineProps<{
  products: Product[]
}>()

const { addToCart } = useMarketplace()
</script>

<template>
  <div class="columns-1 gap-3 sm:columns-2 sm:gap-4 md:columns-3 lg:columns-4 lg:gap-5">
    <div
      v-for="product in products"
      :key="product.id"
      class="mb-3 break-inside-avoid sm:mb-4 lg:mb-5"
    >
      <ProductCard
        :id="product.id"
        :shop="product.shop"
        :name="product.name"
        :price="formatPrice(product.price)"
        :rating="product.rating"
        :image="product.image"
        :stock="product.stock"
        @add-to-cart="addToCart"
      />
    </div>
  </div>
</template>
