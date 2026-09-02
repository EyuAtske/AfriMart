<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const { addToCart, getProductReviews } = useMarketplace()
const { productRepo } = useRepositories()

const productId = computed(() => Number(route.params.id))
const quantity = ref(1)
const added = ref(false)

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

const { data: product, error } = await useAsyncData(
  `product-detail-${route.params.id}`,
  async () => {
    const id = Number(route.params.id)
    if (isNaN(id)) {
      throw createError({ statusCode: 404, statusMessage: 'Invalid Product ID', fatal: true })
    }
    const found = await productRepo.getProductById(id)
    if (!found) {
      throw createError({ statusCode: 404, statusMessage: 'Product not found', fatal: true })
    }
    return found
  }
)

if (error.value || !product.value) {
  throw createError({ statusCode: 404, statusMessage: 'Product not found', fatal: true })
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
    const res = await productRepo.getProducts({ category: product.value.category, pageSize: 5 })
    return res.data.filter(item => item.id !== product.value?.id).slice(0, 4)
  }
)

const addSelectedQuantity = () => {
  if (!product.value) return

  for (let index = 0; index < quantity.value; index += 1) {
    addToCart(product.value.id)
  }

  added.value = true
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <section
      v-if="product"
      class="mx-auto max-w-7xl"
    >
      <button
        type="button"
        class="mb-6 text-sm font-medium text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline"
        @click="router.back()"
      >
        Back to browsing
      </button>

      <div class="grid gap-8 lg:grid-cols-[1fr_0.85fr] lg:gap-12">
        <div class="overflow-hidden rounded-[8px] border border-[#d9d0c4] bg-[#eee8df]">
          <img
            :src="product.image"
            :alt="product.name"
            class="h-full max-h-[760px] min-h-[420px] w-full object-cover object-top"
          />
        </div>

        <div class="flex flex-col justify-center">
          <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
            {{ product.shop }}
          </p>

          <h1 class="mt-3 font-serif text-4xl leading-tight tracking-[-0.025em] text-[#211f1d] sm:text-5xl">
            {{ product.name }}
          </h1>

          <div class="mt-5 flex flex-wrap items-center gap-3">
            <span class="rounded-full bg-[#211f1d] px-4 py-2 text-sm font-semibold text-white">
              {{ formatPrice(product.price) }}
            </span>

            <span class="rounded-full border border-[#d9d0c4] px-4 py-2 text-sm text-[#665c53]">
              ★ {{ calculatedRating }} Rating
            </span>

            <span class="rounded-full border border-[#d9d0c4] px-4 py-2 text-sm text-[#665c53]">
              {{ product.stock }} in stock
            </span>
          </div>

          <p class="mt-7 max-w-2xl text-base leading-8 text-[#665c53]">
            {{ product.description }}
          </p>

          <div class="mt-8 grid gap-4 rounded-[8px] border border-[#d9d0c4] bg-[#faf8f4] p-5 sm:grid-cols-[160px_1fr]">
            <label class="space-y-2">
              <span class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                Quantity
              </span>

              <input
                v-model="quantity"
                type="number"
                min="1"
                :max="product.stock"
                class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#f5f1e9] px-4 text-sm text-[#211f1d] outline-none transition focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
              />
            </label>

            <div class="flex items-end gap-3">
              <button
                type="button"
                class="h-12 flex-1 rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2"
                @click="addSelectedQuantity"
              >
                Add to cart
              </button>

              <NuxtLink
                to="/cart"
                class="inline-flex h-12 items-center justify-center rounded-full bg-[#211f1d] px-6 text-sm font-medium uppercase tracking-[0.14em] text-white transition hover:bg-[#3b3733]"
              >
                View cart
              </NuxtLink>
            </div>
          </div>

          <p
            v-if="added"
            class="mt-4 rounded-[8px] border border-[#cfe0cc] bg-[#eef7ec] px-4 py-3 text-sm font-medium text-[#3f6a3f]"
          >
            Added to cart. You can keep browsing or continue to checkout.
          </p>
        </div>
      </div>

      <!-- Customer Reviews & Ratings Section -->
      <section class="mt-16 rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between border-b border-[#ded6cc] pb-6">
          <div>
            <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
              Feedback & Ratings
            </p>
            <h2 class="mt-1 font-serif text-3xl text-[#211f1d]">
              Customer Reviews
            </h2>
          </div>

          <div class="flex items-center gap-3 rounded-full bg-[#eee8df] px-5 py-2.5">
            <span class="font-serif text-3xl font-bold text-[#211f1d]">{{ calculatedRating }}</span>
            <div>
              <div class="flex text-amber-500 text-sm">
                <span v-for="star in 5" :key="star">{{ star <= Math.round(Number(calculatedRating)) ? '★' : '☆' }}</span>
              </div>
              <p class="text-xs text-[#756a60] mt-0.5">Based on {{ approvedReviews.length }} reviews</p>
            </div>
          </div>
        </div>

        <!-- Pending Admin Approval Notice -->
        <div
          v-if="pendingReviews.length"
          class="mt-6 rounded-[8px] border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800 flex items-center justify-between"
        >
          <div class="flex items-center gap-2">
            <span>⏳</span>
            <span>You have <strong>{{ pendingReviews.length }} review(s)</strong> pending admin approval.</span>
          </div>
          <span class="text-xs font-semibold uppercase tracking-wider text-amber-700">Under Review</span>
        </div>

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
      </section>

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
            class="text-xs font-medium uppercase tracking-[0.14em] text-[#806344]"
          >
            View all
          </NuxtLink>
        </div>

        <MarketplaceProductGrid :products="relatedProducts" />
      </section>
    </section>

    <section
      v-else
      class="mx-auto max-w-xl rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-8 text-center"
    >
      <h1 class="font-serif text-3xl text-[#211f1d]">
        Product not found
      </h1>

      <p class="mt-3 text-sm leading-6 text-[#756a60]">
        This mocked product is not in the current marketplace data.
      </p>

      <NuxtLink
        to="/products"
        class="mt-6 inline-flex h-12 items-center justify-center rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
      >
        Browse products
      </NuxtLink>
    </section>
  </main>
</template>
