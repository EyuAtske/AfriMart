<script setup lang="ts">
import ProductCard from './ProductCard.vue'

export interface Product {
  id?: string | number
  shop: string
  name: string
  price: string
  rating?: string | number
  image: string
}

defineProps<{
  products: Product[]
}>()
</script>

<template>
  <div
    class="columns-2 gap-3 sm:columns-2 sm:gap-4 md:columns-3 lg:columns-4 lg:gap-5"
  >
    <div
      v-for="product in products"
      :key="product.id ?? product.name"
      class="mb-3 break-inside-avoid sm:mb-4 lg:mb-5"
    >
      <ProductCard
        :shop="product.shop"
        :name="product.name"
        :price="product.price"
        :rating="product.rating ?? '0'"
        :image="product.image"
      />
    </div>
  </div>
</template>