<script setup lang="ts">
import ProductGrid from '~/components/marketplace/ProductGrid.vue'
import LoadMoreButton from '~/components/LoadMoreButton.vue'

useSeoMeta({
  title: 'Browse Products — Afrimart Marketplace',
  description: 'Explore unique clothing, dresses, hoodies, sneakers, and accessories from independent sellers across Africa.',
  ogTitle: 'Browse Products — Afrimart Marketplace',
  ogDescription: 'Explore unique clothing, dresses, hoodies, sneakers, and accessories from independent sellers.'
})

const route = useRoute()
const { categories, filterProducts } = useMarketplace()
const { productRepo } = useRepositories()

const search = ref(typeof route.query.search === 'string' ? route.query.search : '')
const selectedCategory = ref(
  typeof route.query.category === 'string' ? route.query.category : 'All'
)
const itemsPerPage = 6
const visibleCount = ref(itemsPerPage)

const { data: asyncProducts } = await useAsyncData(
  'products-catalog-list',
  async () => {
    const res = await productRepo.getProducts({ page: 1, pageSize: 50 })
    return res.data
  }
)

watch(
  () => route.query.search,
  (value) => {
    search.value = typeof value === 'string' ? value : ''
    visibleCount.value = itemsPerPage
  }
)

watch(
  () => route.query.category,
  (value) => {
    selectedCategory.value = typeof value === 'string' ? value : 'All'
    visibleCount.value = itemsPerPage
  }
)

watch([search, selectedCategory], () => {
  visibleCount.value = itemsPerPage
})

const filteredProducts = computed(() =>
  filterProducts({
    search: search.value,
    category: selectedCategory.value
  })
)

const displayedProducts = computed(() =>
  filteredProducts.value.slice(0, visibleCount.value)
)

const hasMore = computed(() =>
  visibleCount.value < filteredProducts.value.length
)

const handleLoadMore = () => {
  visibleCount.value += itemsPerPage
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <section class="mx-auto max-w-[1800px]">
      <div class="mb-8 flex flex-col gap-6 lg:mb-12 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
            Marketplace
          </p>

          <h1 class="mt-2 font-serif text-4xl tracking-[-0.025em] text-[#211f1d] sm:text-5xl">
            Browse products
          </h1>

          <p class="mt-3 max-w-2xl text-base leading-7 text-[#756a60]">
            Search independent sellers, filter by category, and find pieces ready for checkout.
          </p>
        </div>

        <div class="rounded-full bg-[#211f1d] px-5 py-2 text-xs font-medium uppercase tracking-[0.14em] text-white">
          Showing {{ displayedProducts.length }} of {{ filteredProducts.length }} items
        </div>
      </div>

      <!-- Search & Category Filter Controls -->
      <div class="mb-8 grid gap-4 rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-4 shadow-[0_20px_70px_rgba(33,31,29,0.05)] sm:grid-cols-[1fr_240px]">
        <input
          v-model="search"
          type="search"
          placeholder="Search products or shops"
          class="h-12 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-5 text-sm text-[#211f1d] outline-none transition focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
        />

        <select
          v-model="selectedCategory"
          class="h-12 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-5 text-sm text-[#211f1d] outline-none transition focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
        >
          <option
            v-for="category in categories"
            :key="category"
            :value="category"
          >
            {{ category }}
          </option>
        </select>
      </div>

      <template v-if="filteredProducts.length">
        <ProductGrid :products="displayedProducts" />

        <LoadMoreButton
          v-if="hasMore"
          @load-more="handleLoadMore"
        />
      </template>

      <div
        v-else
        class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-10 text-center"
      >
        <h2 class="font-serif text-3xl text-[#211f1d]">
          No products found
        </h2>

        <p class="mt-2 text-sm text-[#756a60]">
          Try a different search term or category.
        </p>
      </div>
    </section>
  </main>
</template>
