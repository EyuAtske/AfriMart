<script setup lang="ts">
const route = useRoute()
const { products } = useMarketplace()
const { shopRepo } = useRepositories()

const selectedCategory = ref('All')

const slugify = (value: string) =>
  value.toLowerCase().replace(/\s+/g, '-')

const currentSlug = computed(() => String(route.params.slug))

const { data: shopData, error } = await useAsyncData(
  `shop-profile-${currentSlug.value}`,
  async () => {
    const found = await shopRepo.getShopBySlug(currentSlug.value)
    if (!found) {
      throw createError({ statusCode: 404, statusMessage: 'Shop not found', fatal: true })
    }
    return found
  }
)

if (error.value || !shopData.value) {
  throw createError({ statusCode: 404, statusMessage: 'Shop not found', fatal: true })
}

useSeoMeta({
  title: computed(() => shopData.value ? `${shopData.value.name} — Seller Storefront` : 'Shop Profile'),
  description: computed(() => shopData.value?.description || 'Browse storefront products on Afrimart.'),
  ogTitle: computed(() => shopData.value ? `${shopData.value.name} Storefront` : 'Shop Profile'),
  ogDescription: computed(() => shopData.value?.description || 'Browse storefront products on Afrimart.')
})

const shopName = computed(() => shopData.value?.name || currentSlug.value.replace(/-/g, ' '))
const shopDescription = computed(() => shopData.value?.description || `Storefront by ${shopName.value}`)

const allShopProducts = computed(() => {
  return products.value.filter(product =>
    slugify(product.shop) === currentSlug.value && product.status === 'Active'
  )
})

const shopCategories = computed(() => [
  'All',
  ...Array.from(new Set(allShopProducts.value.map(p => p.category)))
])

const filteredShopProducts = computed(() => {
  if (selectedCategory.value === 'All') return allShopProducts.value
  return allShopProducts.value.filter(p => p.category === selectedCategory.value)
})

const totalStock = computed(() =>
  allShopProducts.value.reduce((sum, p) => sum + p.stock, 0)
)
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <section class="mx-auto max-w-7xl">
      <!-- Shop Header Banner -->
      <div class="mb-8 overflow-hidden rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
        <div class="flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-5">
            <div class="flex h-16 w-16 items-center justify-center rounded-full bg-[#806344] font-serif text-2xl uppercase text-white shadow-md">
              {{ shopName.charAt(0) }}
            </div>

            <div>
              <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
                Verified Seller Storefront
              </p>

              <h1 class="mt-1 font-serif text-3xl capitalize tracking-[-0.025em] text-[#211f1d] sm:text-4xl">
                {{ shopName }}
              </h1>

              <p class="mt-2 max-w-2xl text-sm leading-6 text-[#756a60]">
                {{ shopDescription }}
              </p>
            </div>
          </div>

          <!-- Shop Metadata Badge -->
          <div class="flex flex-wrap items-center gap-3 rounded-[10px] border border-[#ded6cc] bg-[#f5f1e9] p-4 text-xs">
            <div>
              <p class="font-medium text-[#756a60] uppercase tracking-wider">Listings</p>
              <p class="font-serif text-xl text-[#211f1d]">{{ allShopProducts.length }} items</p>
            </div>

            <div class="h-8 w-[1px] bg-[#ded6cc]" />

            <div>
              <p class="font-medium text-[#756a60] uppercase tracking-wider">Total Stock</p>
              <p class="font-serif text-xl text-[#211f1d]">{{ totalStock }} available</p>
            </div>
          </div>
        </div>

        <!-- Shop Category Filter Tabs -->
        <div v-if="shopCategories.length > 2" class="mt-8 flex flex-wrap gap-2 border-t border-[#ded6cc] pt-6">
          <button
            v-for="cat in shopCategories"
            :key="cat"
            type="button"
            class="rounded-full px-4 py-2 text-xs font-medium uppercase tracking-[0.12em] transition-all"
            :class="
              selectedCategory === cat
                ? 'bg-[#806344] text-white shadow-sm'
                : 'border border-[#cfc4b5] bg-[#f5f1e9] text-[#5d4b37] hover:bg-[#ded6cc]'
            "
            @click="selectedCategory = cat"
          >
            {{ cat }}
          </button>
        </div>
      </div>

      <!-- Shop Products Grid -->
      <MarketplaceProductGrid
        v-if="filteredShopProducts.length"
        :products="filteredShopProducts"
      />

      <section
        v-else
        class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-10 text-center"
      >
        <h2 class="font-serif text-3xl text-[#211f1d]">
          No active products
        </h2>

        <p class="mt-3 text-sm leading-6 text-[#756a60]">
          No products match this category filter in {{ shopName }}.
        </p>
      </section>
    </section>
  </main>
</template>

