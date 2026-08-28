<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'
import TrackShippingModal from '~/components/orders/TrackShippingModal.vue'
import AddReviewModal from '~/components/orders/AddReviewModal.vue'
import type { MarketplaceOrder } from '~/types/order'
import type { Product } from '~/types/product'

definePageMeta({
  middleware: 'auth'
})

const { orders, getOrderProducts, getUserReviewForProduct, formatPrice } = useMarketplace()

const selectedTrackOrder = ref<MarketplaceOrder | null>(null)
const isTrackModalOpen = ref(false)

const selectedReviewProduct = ref<Product | null>(null)
const selectedReviewOrderId = ref<number | undefined>(undefined)
const isReviewModalOpen = ref(false)

const openTrackModal = (order: MarketplaceOrder) => {
  selectedTrackOrder.value = order
  isTrackModalOpen.value = true
}

const openReviewModal = (product: Product, orderId?: number) => {
  selectedReviewProduct.value = product
  selectedReviewOrderId.value = orderId
  isReviewModalOpen.value = true
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-8">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 lg:flex-row">
      <AccountSidebar active="orders" />

      <section class="min-w-0 flex-1">
        <div class="mb-8">
          <h1 class="font-serif text-4xl tracking-[-0.025em] text-[#211f1d]">
            My Orders
          </h1>

          <p class="mt-2 text-base text-[#756a60]">
            View purchases, payment status, and delivery progress.
          </p>
        </div>

        <div
          v-if="orders.length"
          class="space-y-5"
        >
          <article
            v-for="order in orders"
            :key="order.id"
            class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-5 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-6"
          >
            <div class="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
              <div>
                <div class="flex items-center gap-3">
                  <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
                    Order #{{ order.id }}
                  </p>
                  <span class="rounded-full bg-[#eee8df] px-3 py-0.5 text-[11px] font-semibold text-[#211f1d]">
                    {{ order.status }}
                  </span>
                </div>

                <p class="mt-1.5 text-sm text-[#756a60]">
                  Placed on {{ order.date }}
                </p>
              </div>

              <div class="flex flex-wrap items-center gap-2">
                <span class="rounded-full bg-[#e6eee5] px-3 py-1 text-[11px] font-medium uppercase tracking-[0.12em] text-[#536653]">
                  {{ order.paymentStatus }}
                </span>

                <span class="rounded-full bg-[#211f1d] px-3 py-1 text-[11px] font-medium uppercase tracking-[0.12em] text-white">
                  {{ formatPrice(order.total) }}
                </span>

                <button
                  type="button"
                  class="ml-2 inline-flex h-9 items-center justify-center rounded-full border border-[#806344] px-4 text-xs font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
                  @click="openTrackModal(order)"
                >
                  Track Shipping
                </button>
              </div>
            </div>

            <!-- Compact Order Item Summary List -->
            <div class="mt-5 space-y-2 border-t border-[#ded6cc] pt-4">
              <div
                v-for="item in getOrderProducts(order)"
                :key="item?.productId"
                class="flex items-center justify-between gap-4 rounded-[8px] border border-[#ded6cc] bg-[#f5f1e9] p-3"
              >
                <div class="flex items-center gap-3 min-w-0">
                  <img
                    v-if="item"
                    :src="item.product.image"
                    :alt="item.product.name"
                    class="h-12 w-10 rounded object-cover object-top shrink-0"
                  />
                  <div v-if="item" class="min-w-0">
                    <p class="truncate text-sm font-medium text-[#211f1d]">{{ item.product.name }}</p>
                    <p class="text-xs text-[#756a60]">Qty {{ item.quantity }} · {{ item.product.shop }}</p>
                  </div>
                </div>

                <!-- Review / Under Review Badge for Delivered Items -->
                <div v-if="order.status === 'Delivered' && item" class="shrink-0">
                  <template v-if="getUserReviewForProduct(item.product.id, order.id)">
                    <span
                      v-if="getUserReviewForProduct(item.product.id, order.id)?.status === 'pending'"
                      class="inline-flex items-center gap-1.5 rounded-full bg-amber-100 px-3.5 py-1 text-xs font-medium text-amber-800 border border-amber-300"
                    >
                      <span>⏳</span> Under review
                    </span>
                    <span
                      v-else
                      class="inline-flex items-center gap-1.5 rounded-full bg-green-100 px-3.5 py-1 text-xs font-medium text-green-800 border border-green-300"
                    >
                      <span>✓</span> Reviewed
                    </span>
                  </template>
                  <button
                    v-else
                    type="button"
                    class="inline-flex h-8 items-center justify-center rounded-full bg-[#806344] px-3.5 text-xs font-medium uppercase tracking-[0.12em] text-white transition hover:bg-[#5d4b37]"
                    @click="openReviewModal(item.product, order.id)"
                  >
                    Add Review
                  </button>
                </div>
              </div>
            </div>
          </article>
        </div>

        <section
          v-else
          class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-10 text-center"
        >
          <h2 class="font-serif text-3xl text-[#211f1d]">
            No orders yet
          </h2>

          <p class="mt-3 text-sm leading-6 text-[#756a60]">
            Orders created from checkout will show up here.
          </p>
        </section>
      </section>

      <!-- Track Shipping Modal -->
      <TrackShippingModal
        :is-open="isTrackModalOpen"
        :order="selectedTrackOrder"
        @close="isTrackModalOpen = false"
        @open-review="(product) => { isTrackModalOpen = false; openReviewModal(product, selectedTrackOrder?.id); }"
      />

      <!-- Add Review Modal -->
      <AddReviewModal
        :is-open="isReviewModalOpen"
        :product="selectedReviewProduct"
        :order-id="selectedReviewOrderId"
        @close="isReviewModalOpen = false"
      />
    </div>
  </main>
</template>
