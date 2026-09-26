<script setup lang="ts">
import ProductGrid from './ProductGrid.vue'

const props = withDefaults(
  defineProps<{
    limit?: number
  }>(),
  {
    limit: 8
  }
)

const { productRepo } = useRepositories()
const { products: marketplaceProducts } = useMarketplace()

const { data: apiProducts, pending, error } = await useAsyncData(
  'home-featured-products',
  async () => {
    const res = await productRepo.getProducts({ page: 1, pageSize: 50 })
    return res.data
  }
)

const featuredProducts = computed(() => {
  const source = marketplaceProducts.value && marketplaceProducts.value.length > 0
    ? marketplaceProducts.value
    : (apiProducts.value || [])
  return source
    .filter(product => product.status === 'Active')
    .slice(0, props.limit)
})
</script>

<template>
  <section
    class="bg-[#f5f1e9] px-4 py-12 sm:px-6 sm:py-16 lg:px-12 lg:py-24"
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
        class="group flex items-center gap-1.5 text-xs tracking-[0.1em] text-[#806344] sm:gap-3 sm:text-sm"
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

    <!-- Loading State -->
    <div v-if="pending" class="py-12 text-center">
      <div class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-[#806344] border-r-transparent align-[-0.125em]"></div>
      <p class="mt-4 text-sm text-[#756a60]">Loading products...</p>
    </div>

    <!-- Products -->
    <ProductGrid v-else-if="featuredProducts.length" :products="featuredProducts" />

    <!-- Empty state -->
    <p v-else class="py-12 text-center text-sm text-[#756a60]">
      No products available yet. Check back soon!
    </p>
  </section>
</template>
