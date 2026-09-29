<script setup lang="ts">
import ProductGrid from '~/components/marketplace/ProductGrid.vue'
import LoadMoreButton from '~/components/LoadMoreButton.vue'
import { ensureCategoryCatalog, STATIC_CATALOG } from '~/utils/categoryCatalog'
const { gtag } = useGtag()
const { track } = useAnalytics()

useSeoMeta({
  title: 'Browse Products — Afrimart Marketplace',
  description: 'Explore unique clothing, dresses, hoodies, sneakers, and accessories from independent sellers across Africa.',
  ogTitle: 'Browse Products — Afrimart Marketplace',
  ogDescription: 'Explore unique clothing, dresses, hoodies, sneakers, and accessories from independent sellers.'
})

const route = useRoute()
const { categories, filterProducts, products: marketplaceProducts } = useMarketplace()
const { productRepo } = useRepositories()

await ensureCategoryCatalog()

const search = ref(typeof route.query.search === 'string' ? route.query.search : '')
const selectedCategory = ref(
  typeof route.query.category === 'string' ? route.query.category : 'All'
)
const selectedSubcategory = ref(
  typeof route.query.subcategory === 'string' ? route.query.subcategory : 'All'
)
const minPrice = ref<number | ''>('')
const maxPrice = ref<number | ''>('')
const availability = ref<'all' | 'available' | 'sold-out'>('all')
const sortBy = ref<'default' | 'price-asc' | 'price-desc'>('default')
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

watch(
  () => route.query.subcategory,
  (value) => {
    selectedSubcategory.value = typeof value === 'string' ? value : 'All'
    visibleCount.value = itemsPerPage
  }
)

watch([search, selectedCategory, selectedSubcategory, minPrice, maxPrice, availability, sortBy], () => {
  visibleCount.value = itemsPerPage
})

const subcategories = computed(() => {
  const category = STATIC_CATALOG.find(item => item.name === selectedCategory.value)
  const available = selectedCategory.value === 'All'
    ? STATIC_CATALOG.flatMap(item => item.subcategories)
    : category?.subcategories || []
  return ['All', ...new Set(available.map(subcategory => subcategory.name))]
})

const catalogProducts = computed(() => {
  return marketplaceProducts.value && marketplaceProducts.value.length > 0
    ? marketplaceProducts.value
    : (asyncProducts.value || [])
})

const filteredProducts = computed(() =>
  filterProducts(
    {
      search: search.value,
      category: selectedCategory.value,
      subCategory: selectedSubcategory.value,
      minPrice: minPrice.value === '' ? undefined : minPrice.value,
      maxPrice: maxPrice.value === '' ? undefined : maxPrice.value,
      availability: availability.value,
      sortBy: sortBy.value
    },
    catalogProducts.value
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

const clearMarketplaceFilters = () => {
  selectedCategory.value = 'All'
  selectedSubcategory.value = 'All'
  minPrice.value = ''
  maxPrice.value = ''
  availability.value = 'all'
  sortBy.value = 'default'
}

const handleCategoryChange = () => {
  selectedSubcategory.value = 'All'
  track('category_filter_used', {
    category: selectedCategory.value
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

      <!-- Search and marketplace filters -->
      <UiAppCard class="mb-8 p-4 shadow-[0_20px_70px_rgba(33,31,29,0.05)]">
        <div class="grid gap-4 sm:grid-cols-[minmax(0,1fr)_220px]">
          <input v-model="search" type="search" placeholder="Search products or shops" @keyup.enter="handleSearch"
            class="h-12 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-5 text-sm text-[#211f1d] outline-none transition placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
          />

          <select
            v-model="selectedCategory"
            @change="handleCategoryChange"
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

        <div class="mt-4 grid grid-cols-2 gap-3 md:grid-cols-3 xl:grid-cols-6">
          <label class="min-w-0 space-y-1.5 text-xs font-medium text-[#4d4035]">
            <span class="block uppercase tracking-[0.1em]">Subcategory</span>
            <select v-model="selectedSubcategory" class="h-11 w-full min-w-0 rounded-md border border-[#cfc4b5] bg-[#f5f1e9] px-3 text-sm text-[#211f1d] outline-none focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15">
              <option v-for="subcategory in subcategories" :key="subcategory" :value="subcategory">{{ subcategory }}</option>
            </select>
          </label>

          <label class="min-w-0 space-y-1.5 text-xs font-medium text-[#4d4035]">
            <span class="block uppercase tracking-[0.1em]">Min price (ETB)</span>
            <input v-model.number="minPrice" type="number" min="0" step="any" placeholder="Any" class="h-11 w-full min-w-0 rounded-md border border-[#cfc4b5] bg-[#f5f1e9] px-3 text-sm text-[#211f1d] outline-none placeholder:text-[#92877b] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15" />
          </label>

          <label class="min-w-0 space-y-1.5 text-xs font-medium text-[#4d4035]">
            <span class="block uppercase tracking-[0.1em]">Max price (ETB)</span>
            <input v-model.number="maxPrice" type="number" min="0" step="any" placeholder="Any" class="h-11 w-full min-w-0 rounded-md border border-[#cfc4b5] bg-[#f5f1e9] px-3 text-sm text-[#211f1d] outline-none placeholder:text-[#92877b] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15" />
          </label>

          <label class="min-w-0 space-y-1.5 text-xs font-medium text-[#4d4035]">
            <span class="block uppercase tracking-[0.1em]">Availability</span>
            <select v-model="availability" class="h-11 w-full min-w-0 rounded-md border border-[#cfc4b5] bg-[#f5f1e9] px-3 text-sm text-[#211f1d] outline-none focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15">
              <option value="all">All</option>
              <option value="available">In stock</option>
              <option value="sold-out">Sold out</option>
            </select>
          </label>

          <label class="min-w-0 space-y-1.5 text-xs font-medium text-[#4d4035]">
            <span class="block uppercase tracking-[0.1em]">Sort</span>
            <select v-model="sortBy" class="h-11 w-full min-w-0 rounded-md border border-[#cfc4b5] bg-[#f5f1e9] px-3 text-sm text-[#211f1d] outline-none focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15">
              <option value="default">Default order</option>
              <option value="price-asc">Price: low to high</option>
              <option value="price-desc">Price: high to low</option>
            </select>
          </label>

          <button type="button" class="self-end justify-self-start pb-2 text-xs font-semibold uppercase tracking-[0.1em] text-[#806344] hover:underline" @click="clearMarketplaceFilters">
            Clear filters
          </button>
        </div>
      </UiAppCard>

      <nav v-if="selectedCategory !== 'All'" aria-label="Breadcrumb" class="mb-5 flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1 text-sm text-[#756a60]">
        <NuxtLink to="/" class="hover:text-[#806344] hover:underline">Home</NuxtLink>
        <span aria-hidden="true">/</span>
        <NuxtLink :to="{ path: '/products', query: { category: selectedCategory } }" class="hover:text-[#806344] hover:underline">{{ selectedCategory }}</NuxtLink>
        <template v-if="selectedSubcategory !== 'All'">
          <span aria-hidden="true">/</span>
          <span class="min-w-0 truncate text-[#211f1d]">{{ selectedSubcategory }}</span>
        </template>
      </nav>

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
        v-else-if="!catalogProducts || catalogProducts.length === 0"
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
