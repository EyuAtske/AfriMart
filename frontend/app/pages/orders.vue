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
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
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
          <UiAppCard
            v-for="order in orders"
            :key="order.id"
          >
            <div class="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
              <div>
                <div class="flex items-center gap-3">
                  <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
                    Order #{{ order.id }}
                  </p>
                  <UiAppBadge variant="default">
                    {{ order.status }}
                  </UiAppBadge>
                </div>

                <p class="mt-1.5 text-sm text-[#756a60]">
                  Placed on {{ order.date }}
                </p>
              </div>

              <div class="flex flex-wrap items-center gap-2">
                <UiAppBadge variant="success">
                  {{ order.paymentStatus }}
                </UiAppBadge>

                <UiAppBadge variant="dark">
                  {{ formatPrice(order.total) }}
                </UiAppBadge>

                <UiAppButton
                  size="small"
                  class="ml-2"
                  @click="openTrackModal(order)"
                >
                  Track Shipping
                </UiAppButton>
              </div>
            </div>

            <!-- Compact Order Item Summary List -->
            <div class="mt-5 space-y-2 border-t border-[#ded6cc] pt-4">
              <div
                v-for="item in getOrderProducts(order)"
                :key="item?.productId"
                class="flex items-center justify-between gap-4 rounded-lg border border-[#ded6cc] bg-[#f5f1e9] p-3"
              >
                <div class="flex items-center gap-3 min-w-0">
                  <img
                    v-if="item"
                    :src="item.product.image"
                    :alt="item.product.name"
                    class="h-12 w-10 rounded-md object-cover object-top shrink-0"
                  />
                  <div v-if="item" class="min-w-0">
                    <p class="truncate text-sm font-medium text-[#211f1d]">{{ item.product.name }}</p>
                    <p class="text-xs text-[#756a60]">Qty {{ item.quantity }} · {{ item.product.shop }}</p>
                  </div>
                </div>

                <!-- Review / Under Review Badge for Delivered Items -->
                <div v-if="order.status === 'Delivered' && item" class="shrink-0">
                  <template v-if="getUserReviewForProduct(item.product.id, order.id)">
                    <UiAppBadge
                      v-if="getUserReviewForProduct(item.product.id, order.id)?.status === 'pending'"
                      variant="warning"
                    >
                      <span>⏳</span> Under review
                    </UiAppBadge>
                    <UiAppBadge
                      v-else
                      variant="success"
                    >
                      <span>✓</span> Reviewed
                    </UiAppBadge>
                  </template>
                  <UiAppButton
                    v-else
                    variant="primary"
                    size="small"
                    class="h-8 px-3.5 text-xs"
                    @click="openReviewModal(item.product, order.id)"
                  >
                    Add Review
                  </UiAppButton>
                </div>
              </div>
            </div>
          </UiAppCard>
        </div>

        <UiAppEmptyState
          v-else
          title="No orders yet"
          description="Orders created from checkout will show up here."
          action-label="Start shopping"
          action-to="/products"
        />
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
