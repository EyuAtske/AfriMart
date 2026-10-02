<script setup lang="ts">
import type { Product } from '~/types/product'

const route = useRoute()
const { products } = useMarketplace()
const { shopRepo, productRepo } = useRepositories()
const { gtag } = useGtag()
const selectedCategory = ref('All')

const slugify = (value: string) =>
  value.toLowerCase().replace(/\s+/g, '-')

const currentSlug = computed(() => String(route.params.slug))
const shopId = computed(() => typeof route.query.shopId === 'string' ? route.query.shopId : '')
const requestedShopName = computed(() =>
  typeof route.query.shopName === 'string' && route.query.shopName.trim()
    ? route.query.shopName
    : products.value.find(product => String(product.shopId || '') === shopId.value)?.shop || ''
)
const shopProducts = ref<Product[] | null>(null)
const isLoadingShopProducts = ref(Boolean(shopId.value))
const shopProductsError = ref('')

const { data: shopData, error } = await useAsyncData(
  `shop-profile-${currentSlug.value}-${shopId.value}`,
  async () => {
    const found = await shopRepo.getShopBySlug(currentSlug.value)
    if (found) return found

    if (shopId.value) {
      return {
        id: shopId.value,
        backendId: shopId.value,
        name: requestedShopName.value,
        slug: currentSlug.value,
        description: '',
        ownerEmail: '',
        products: [],
        paymentMethods: []
      }
    }

    throw createError({ statusCode: 404, statusMessage: 'Shop not found', fatal: true })
  }
)

  if (error.value || !shopData.value) {
   throw createError({ statusCode: 404, statusMessage: 'Shop not found', fatal: true })
  }
  onMounted(() => {
    if (shopData.value) {
       gtag('event', 'view_shop', {
        shop_name: shopData.value.name
      })
   }
   loadShopProducts()
   })
  useSeoMeta({
    title: computed(() => shopData.value ? `${shopData.value.name} — Seller Storefront` : 'Shop Profile'),
    description: computed(() => shopData.value?.description || 'Browse storefront products on Afrimart.'),
    ogTitle: computed(() => shopData.value ? `${shopData.value.name} Storefront` : 'Shop Profile'),
    ogDescription: computed(() => shopData.value?.description || 'Browse storefront products on Afrimart.')
})

const shopName = computed(() => shopData.value?.name || requestedShopName.value)
const shopDescription = computed(() => shopData.value?.description || `Storefront by ${shopName.value}`)

const loadShopProducts = async () => {
  if (!shopId.value) return

  isLoadingShopProducts.value = true
  shopProductsError.value = ''
  try {
    const { data } = await productRepo.getProducts({ page: 1, pageSize: 1000 })
    const matchingProducts = data.filter(product =>
      product.shopId
        ? String(product.shopId) === shopId.value
        : slugify(product.shop) === currentSlug.value
    )
    shopProducts.value = matchingProducts.map(product => ({ ...product, shop: shopName.value }))
  } catch (err: any) {
    shopProductsError.value = err?.message || "Unable to load this shop's products. Please try again."
  } finally {
    isLoadingShopProducts.value = false
  }
}

const allShopProducts = computed(() => {
  const matchingProducts = shopProducts.value ?? products.value.filter(product =>
    slugify(product.shop) === currentSlug.value
  )
  return matchingProducts.filter(product => product.status === 'Active')
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
      <UiAppCard padding="large" class="mb-8 overflow-hidden">
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
          <div class="flex flex-wrap items-center gap-3 rounded-lg border border-[#ded6cc] bg-[#f5f1e9] p-4 text-xs">
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
            class="rounded-full px-4 py-2 text-xs font-medium uppercase tracking-[0.12em] transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-[#806344]"
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
      </UiAppCard>

      <!-- Shop Products Grid -->
      <div v-if="isLoadingShopProducts" class="py-16 text-center">
        <div class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-[#806344] border-r-transparent align-[-0.125em]"></div>
        <p class="mt-4 text-sm text-[#756a60]">Loading shop products...</p>
      </div>

      <div v-else-if="shopProductsError" class="py-12 text-center">
        <UiAppAlert variant="error" class="mb-6 max-w-xl mx-auto">
          {{ shopProductsError }}
        </UiAppAlert>
        <UiAppButton variant="secondary" @click="loadShopProducts">Retry</UiAppButton>
      </div>

      <MarketplaceProductGrid
        v-else-if="filteredShopProducts.length"
        :products="filteredShopProducts"
      />

      <UiAppEmptyState
        v-else
        title="No active products"
        :description="`No products match this category filter in ${shopName}.`"
      />
    </section>
  </main>
</template>

