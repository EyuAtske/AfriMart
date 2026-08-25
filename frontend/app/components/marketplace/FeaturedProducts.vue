<script setup lang="ts">
import ProductGrid from './ProductGrid.vue'

const products = [
  {
    id: 1,
    shop: 'Atelier North',
    name: 'Shirt for men',
    price: '1,400 ETB',
    rating: '4.8',
    image: '/images/product1.jpg'
  },
  {
    id: 2,
    shop: 'Beyond Score',
    name: 'Tank Tops',
    price: '2,500 ETB',
    rating: '4.9',
    image: '/images/product2.jpg'
  },
  {
    id: 3,
    shop: 'True Form',
    name: 'Casual Outfit Set',
    price: '3,500 ETB',
    rating: '4.6',
    image: '/images/product3.jpg'
  },
  {
    id: 4,
    shop: 'Minimal Studio',
    name: 'Everyday fit for kids',
    price: '2,900 ETB',
    rating: '4.8',
    image: '/images/product4.jpg'
  },
  {
    id: 5,
    shop: 'Atelier North',
    name: 'Cute dress for kids',
    price: '3,800 ETB',
    rating: '4.7',
    image: '/images/product5.jpg'
  },
  {
    id: 6,
    shop: 'Urban Thread',
    name: 'Hoodie',
    price: '4,700 ETB',
    rating: '4.8',
    image: '/images/product6.jpg'
  },
  {
    id: 7,
    shop: 'Mara Studio',
    name: 'Shirt for men',
    price: '1,550 ETB',
    rating: '4.9',
    image: '/images/product7.jpg'
  },
  {
    id: 8,
    shop: 'ute dress',
    name: 'Relaxed Denim',
    price: '5,200 ETB',
    rating: '4.7',
    image: '/images/product8.jpg'
  },
  {
    id: 9,
    shop: 'Minimal Studio',
    name: 'Watch for women',
    price: '2,500 ETB',
    rating: '4.9',
    image: '/images/product9.jpg'
  },
  {
    id: 10,
    shop: 'Mara Studio',
    name: 'Hat',
    price: '1,100 ETB',
    rating: '4.8',
    image: '/images/product10.jpg'
  }
]
</script>

<template>
  <section
    class="bg-[#f5f1e9] px-3 py-12 sm:px-6 sm:py-16 lg:px-12 lg:py-24"
  >
    <!-- Header -->
    <div class="mb-8 flex items-end justify-between sm:mb-10 lg:mb-12">
      <div>
        <h2
          class="font-serif text-2xl tracking-[-0.02em] text-[#211f1d] sm:text-4xl lg:text-5xl"
        >
          Featured Products
        </h2>

        <div
          class="mt-3 h-px w-10 bg-[#806344] sm:mt-5 sm:w-14"
        ></div>
      </div>

      <NuxtLink to="/products"
        class="group flex items-center gap-1.5 text-[9px] tracking-[0.1em] text-[#806344] sm:gap-3 sm:text-sm"
      >
        VIEW ALL

        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          class="transition-transform duration-300 group-hover:translate-x-1 sm:h-5 sm:w-5"
        >
          <path d="M5 12h14" />
          <path d="m13 6 6 6-6 6" />
        </svg>
      </NuxtLink>
    </div>

    <!-- Products -->
    <ProductGrid :products="products" />
    
  </section>
</template>