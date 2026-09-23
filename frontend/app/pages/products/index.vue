<script setup lang="ts">
import ProductGrid from '~/components/marketplace/ProductGrid.vue'
import LoadMoreButton from '~/components/LoadMoreButton.vue'
const { gtag } = useGtag()
const { track } = useAnalytics()

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

const { data: asyncProducts, pending, error, refresh } = await useAsyncData(
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
  filterProducts(
    {
      search: search.value,
      category: selectedCategory.value
    },
    asyncProducts.value || []
  )
)

const displayedProducts = computed(() =>
  filteredProducts.value.slice(0, visibleCount.value)
)

const hasMore = computed(() =>
  visibleCount.value < filteredProducts.value.length
)

const handleLoadMore = () => {
  visibleCount.value += itemsPerPage

  track('load_more_products', {
    category: selectedCategory.value,
    search_term: search.value.trim() || undefined,
    products_visible: visibleCount.value
  })
}
const handleSearch = () => {
  const term = search.value.trim()

  if (!term) return

  gtag('event', 'search', {
    search_term: term
  })
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

        <UiAppBadge v-if="!pending && !error" variant="dark">
          Showing {{ displayedProducts.length }} of {{ filteredProducts.length }} items
        </UiAppBadge>
      </div>

      <!-- Search & Category Filter Controls -->
      <UiAppCard class="mb-8 p-4 shadow-[0_20px_70px_rgba(33,31,29,0.05)]">
        <div class="grid gap-4 sm:grid-cols-[1fr_240px]">
          <input v-model="search" type="search" placeholder="Search products or shops" @keyup.enter="handleSearch"
            class="h-12 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-5 text-sm text-[#211f1d] outline-none transition placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
          />

          <select
            v-model="selectedCategory"
            @change="track('category_filter_used', {
              category: selectedCategory
            })"
            class="h-12 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-5 text-sm text-[#211f1d] outline-none transition hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
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
      </UiAppCard>

      <!-- Loading State -->
      <div v-if="pending" class="py-16 text-center">
        <div class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-[#806344] border-r-transparent align-[-0.125em]"></div>
        <p class="mt-4 text-sm text-[#756a60]">Loading catalog products...</p>
      </div>

      <!-- Error State with Retry -->
      <div v-else-if="error" class="py-12 text-center">
        <UiAppAlert variant="error" class="mb-6 max-w-xl mx-auto">
          {{ error.message || 'Unable to retrieve products from the server. Please try again.' }}
        </UiAppAlert>
        <UiAppButton variant="secondary" @click="() => refresh()">Retry</UiAppButton>
      </div>

      <!-- Empty Catalog State (0 items in entire catalog) -->
      <UiAppEmptyState
        v-else-if="!asyncProducts || asyncProducts.length === 0"
        title="Catalog is empty"
        description="No products are available in the marketplace yet. Check back soon or list your own products!"
        action-label="Open your shop"
        action-to="/sell"
      />

      <!-- Filtered Products Grid -->
      <template v-else-if="filteredProducts.length">
        <ProductGrid :products="displayedProducts" />

        <LoadMoreButton
          v-if="hasMore"
          @load-more="handleLoadMore"
        />
      </template>

      <!-- Filtered Empty State (search/filter matched 0) -->
      <UiAppEmptyState
        v-else
        title="No products found"
        description="Try a different search term or category."
      />
    </section>
  </main>
</template>
