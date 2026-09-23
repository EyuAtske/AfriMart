<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const { addToCart, getProductReviews, isOwnProduct } = useMarketplace()
const { productRepo } = useRepositories()
const { gtag } = useGtag()

const productId = computed(() => String(route.params.id || ''))
const quantity = ref(1)
const added = ref(false)
const actionError = ref('')

const isSelfProduct = computed(() => product.value ? isOwnProduct(product.value) : false)
const isSoldOut = computed(() => product.value ? product.value.stock <= 0 : false)

const productReviews = computed(() =>
  product.value ? getProductReviews(product.value.id) : []
)

const approvedReviews = computed(() =>
  productReviews.value.filter(r => r.status === 'approved')
)

const pendingReviews = computed(() =>
  productReviews.value.filter(r => r.status === 'pending')
)

const calculatedRating = computed(() => {
  if (!approvedReviews.value.length) return product.value?.rating || '5.0'
  const sum = approvedReviews.value.reduce((acc, r) => acc + r.rating, 0)
  return (sum / approvedReviews.value.length).toFixed(1)
})

const { data: product, pending, error, refresh } = await useAsyncData(
  `product-detail-${route.params.id}`,
  async () => {
    const id = String(route.params.id || '').replace(/[()[\]<>{}`'"]/g, '').trim()
    if (!id) {
      throw createError({ statusCode: 404, statusMessage: 'Invalid Product ID', fatal: false })
    }
    const found = await productRepo.getProductById(id)
    if (!found) {
      throw createError({ statusCode: 404, statusMessage: 'Product not found', fatal: false })
    }
    return found
  }
)

if (import.meta.client && product.value) {
  gtag('event', 'view_item', {
    currency: 'ETB',
    value: Number(product.value.price),
    items: [
      {
        item_id: String(product.value.id),
        item_name: product.value.name,
        item_category: product.value.category,
        price: Number(product.value.price),
        quantity: 1
      }
    ]
  })
}
useSeoMeta({
  title: computed(() => product.value ? `${product.value.name} — ${product.value.shop}` : 'Product Details'),
  description: computed(() => product.value?.description || 'Product details on Afrimart.'),
  ogTitle: computed(() => product.value ? `${product.value.name} (${formatPrice(product.value.price)})` : 'Product Details'),
  ogDescription: computed(() => product.value?.description || 'Product details on Afrimart.'),
  ogImage: computed(() => product.value?.image || '/images/herotemp.png')
})

const { data: relatedProducts } = await useAsyncData(
  `product-related-${route.params.id}`,
  async () => {
    if (!product.value) return []
    try {
      const res = await productRepo.getProducts({ category: product.value.category, pageSize: 5 })
      return res.data.filter(item => String(item.id) !== String(product.value?.id)).slice(0, 4)
    } catch {
      return []
    }
  }
)

const { flyToCart } = useFlyToCart()
const galleryContainerEl = ref<HTMLElement | null>(null)

const addSelectedQuantity = async () => {
  if (!product.value) return
  actionError.value = ''
  added.value = false

  if (isSelfProduct.value) {
    actionError.value = 'You cannot purchase products listed by your own shop.'
    return
  }

  if (isSoldOut.value) {
    actionError.value = 'This product is currently sold out.'
    return
  }

  try {
    flyToCart(galleryContainerEl.value, product.value.image)

    for (let index = 0; index < quantity.value; index += 1) {
      await addToCart(product.value.id)
    }

    added.value = true
  } catch (err: any) {
    actionError.value = err?.message || 'Failed to add item to cart'
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-3 py-8 sm:px-6 sm:py-16 lg:px-12 lg:py-20">
    <!-- Loading State -->
    <div v-if="pending" class="mx-auto max-w-7xl py-12 text-center">
      <div class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-[#806344] border-r-transparent align-[-0.125em]"></div>
      <p class="mt-4 text-sm text-[#756a60]">Loading product details...</p>
    </div>

    <!-- Error State with Retry -->
    <div v-else-if="error && !product && error.statusCode !== 404" class="mx-auto max-w-xl py-12 text-center">
      <UiAppAlert variant="error" class="mb-6">
        {{ error.statusMessage || 'Unable to load product details. Please try again.' }}
      </UiAppAlert>
      <div class="flex justify-center gap-4">
        <UiAppButton variant="secondary" @click="() => refresh()">Retry</UiAppButton>
        <UiAppButton to="/products" variant="ghost">Back to products</UiAppButton>
      </div>
    </div>

    <!-- Product Content -->
    <section
      v-else-if="product"
      class="mx-auto max-w-7xl"
    >
      <button
        type="button"
        class="mb-4 sm:mb-6 text-xs sm:text-sm font-medium text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline"
        @click="router.back()"
      >
        ← Back to browsing
      </button>

      <!-- Side-by-side product view on mobile & desktop -->
      <div class="grid grid-cols-2 gap-3 sm:gap-8 lg:grid-cols-[1fr_0.85fr] lg:gap-12 items-start">
        <div ref="galleryContainerEl" class="min-w-0">
          <MarketplaceMediaGallery
            :media="product.media"
            :fallback-image="product.image"
          />
        </div>

        <div class="flex flex-col justify-start min-w-0">
          <p class="text-[10px] sm:text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
            {{ product.shop }}
          </p>

          <h1 class="mt-1 sm:mt-3 font-serif text-sm sm:text-3xl lg:text-4xl leading-tight tracking-[-0.025em] text-[#211f1d]">
            {{ product.name }}
          </h1>

          <div class="mt-2 sm:mt-5 flex flex-wrap items-center gap-1.5 sm:gap-3">
            <UiAppBadge variant="dark" class="text-[11px] sm:text-xs">
              {{ formatPrice(product.price) }}
            </UiAppBadge>

            <span class="rounded-full border border-[#d9d0c4] px-2 py-0.5 sm:px-3 sm:py-1 text-[9px] sm:text-xs text-[#665c53]">
              ★ {{ calculatedRating }} Rating
            </span>

            <span
              v-if="isSoldOut"
              class="rounded-full bg-red-100 text-red-800 border border-red-300 px-2 py-0.5 sm:px-3 sm:py-1 text-[9px] sm:text-xs font-bold uppercase tracking-wider"
            >
              Sold Out
            </span>
            <span
              v-else-if="isSelfProduct"
              class="rounded-full bg-amber-100 text-amber-900 border border-amber-300 px-2 py-0.5 sm:px-3 sm:py-1 text-[9px] sm:text-xs font-semibold"
            >
              Your Product
            </span>
            <span
              v-else
              class="rounded-full border border-[#d9d0c4] px-2 py-0.5 sm:px-3 sm:py-1 text-[9px] sm:text-xs text-[#665c53]"
            >
              {{ product.stock }} in stock
            </span>
          </div>

          <p class="mt-2 sm:mt-6 max-w-2xl text-[11px] sm:text-base leading-relaxed sm:leading-8 text-[#665c53] line-clamp-3 sm:line-clamp-none">
            {{ product.description }}
          </p>

          <UiAppCard class="mt-3 sm:mt-8 p-2.5 sm:p-5">
            <UiAppAlert v-if="isSelfProduct" variant="warning" class="mb-3 text-xs sm:text-sm">
              This item is listed by your shop. Sellers cannot purchase their own products.
            </UiAppAlert>
            <UiAppAlert v-else-if="isSoldOut" variant="error" class="mb-3 text-xs sm:text-sm">
              This product is currently sold out.
            </UiAppAlert>

            <div class="flex flex-col gap-2.5 sm:flex-row sm:items-end">
              <label class="space-y-1 sm:space-y-2 sm:w-28">
                <span class="block text-[10px] sm:text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                  Quantity
                </span>

                <input
                  v-model="quantity"
                  type="number"
                  min="1"
                  :max="product.stock"
                  :disabled="isSoldOut || isSelfProduct"
                  class="h-8 sm:h-12 w-full rounded-md border border-[#cfc4b5] bg-[#f5f1e9] px-2 text-xs sm:text-sm text-[#211f1d] outline-none transition hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15 disabled:opacity-50"
                />
              </label>

              <div class="flex flex-1 items-center gap-2">
                <UiAppButton
                  class="flex-1 text-xs sm:text-sm h-8 sm:h-12"
                  variant="secondary"
                  :disabled="isSoldOut || isSelfProduct"
                  @click="addSelectedQuantity"
                >
                  <template v-if="isSoldOut">Sold Out</template>
                  <template v-else-if="isSelfProduct">Your Product</template>
                  <template v-else>Add to cart</template>
                </UiAppButton>

                <UiAppButton
                  to="/cart"
                  variant="ghost"
                  class="text-xs sm:text-sm h-8 sm:h-12"
                >
                  Cart
                </UiAppButton>
              </div>
            </div>
          </UiAppCard>

          <UiAppAlert
            v-if="actionError"
            variant="error"
            class="mt-3 text-xs sm:text-sm"
          >
            {{ actionError }}
          </UiAppAlert>

          <UiAppAlert
            v-else-if="added"
            variant="success"
            class="mt-3 text-xs sm:text-sm"
          >
            Added to cart. You can continue browsing or view cart.
          </UiAppAlert>
        </div>
      </div>

      <!-- Customer Reviews & Ratings Section -->
      <UiAppCard padding="large" class="mt-10 sm:mt-16 p-4 sm:p-8">
        <div class="flex flex-row items-center justify-between border-b border-[#ded6cc] pb-4 sm:pb-6 gap-2">
          <div>
            <p class="text-[10px] sm:text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
              Feedback & Ratings
            </p>
            <h2 class="mt-0.5 sm:mt-1 font-serif text-lg sm:text-3xl text-[#211f1d]">
              Customer Reviews
            </h2>
          </div>

          <div class="flex items-center gap-2 sm:gap-3 rounded-full bg-[#eee8df] px-3 py-1.5 sm:px-5 sm:py-2.5 shrink-0">
            <span class="font-serif text-lg sm:text-3xl font-bold text-[#211f1d]">{{ calculatedRating }}</span>
            <div>
              <div class="flex text-amber-500 text-xs sm:text-sm">
                <span v-for="star in 5" :key="star">{{ star <= Math.round(Number(calculatedRating)) ? '★' : '☆' }}</span>
              </div>
              <p class="text-[10px] sm:text-xs text-[#756a60] mt-0.5">{{ approvedReviews.length }} reviews</p>
            </div>
          </div>
        </div>

        <!-- Pending Admin Approval Notice -->
        <UiAppAlert
          v-if="pendingReviews.length"
          variant="warning"
          class="mt-6 flex items-center justify-between"
        >
          <div class="flex items-center gap-2">
            <span>⏳</span>
            <span>You have <strong>{{ pendingReviews.length }} review(s)</strong> pending admin approval.</span>
          </div>
          <span class="text-xs font-semibold uppercase tracking-wider text-amber-700">Under Review</span>
        </UiAppAlert>

        <!-- Approved Reviews List -->
        <div v-if="approvedReviews.length" class="mt-6 divide-y divide-[#ded6cc]">
          <article
            v-for="review in approvedReviews"
            :key="review.id"
            class="py-5 first:pt-0 last:pb-0"
          >
            <div class="flex items-center justify-between gap-4">
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-full bg-[#806344] text-sm font-bold text-white uppercase">
                  {{ review.author.charAt(0) }}
                </div>
                <div>
                  <p class="text-sm font-semibold text-[#211f1d]">{{ review.author }}</p>
                  <p class="text-xs text-[#756a60]">{{ review.date }}</p>
                </div>
              </div>

              <div class="flex text-amber-500 text-sm">
                <span v-for="star in 5" :key="star">{{ star <= review.rating ? '★' : '☆' }}</span>
              </div>
            </div>

            <p class="mt-3 text-sm leading-6 text-[#665c53]">
              {{ review.comment }}
            </p>
          </article>
        </div>

        <div v-else-if="!pendingReviews.length" class="mt-6 text-center py-6 text-sm text-[#756a60]">
          No customer reviews for this product yet. Be the first to leave a review after your delivery!
        </div>
      </UiAppCard>

      <section
        v-if="relatedProducts && relatedProducts.length"
        class="mt-16"
      >
        <div class="mb-6 flex items-end justify-between gap-4">
          <h2 class="font-serif text-3xl text-[#211f1d]">
            More from {{ product.category }}
          </h2>

          <NuxtLink
            :to="`/products?category=${product.category}`"
            class="text-xs font-medium uppercase tracking-[0.14em] text-[#806344] hover:underline"
          >
            View all
          </NuxtLink>
        </div>

        <MarketplaceProductGrid :products="relatedProducts" />
      </section>
    </section>

    <UiAppEmptyState
      v-else
      title="Product not found"
      description="The requested product could not be found in our catalog."
      action-label="Browse products"
      action-to="/products"
      class="mx-auto max-w-xl"
    />
  </main>
</template>
