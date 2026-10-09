<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'
import TrackShippingModal from '~/components/orders/TrackShippingModal.vue'
import AddReviewModal from '~/components/orders/AddReviewModal.vue'
import { formatPrice, useMarketplace } from '~/composables/useMarketplace'
import { useRepositories } from '~/composables/useRepositories'
import {
  formatPaymentSummary,
  getOrderItems,
  getOrderNumber,
  getOrderTitle,
  hasReviewableProduct,
  toReviewProduct
} from '~/utils/orderDisplay'
import type { CartItem, MarketplaceOrder } from '~/types/order'
import type { Product } from '~/types/product'

definePageMeta({
  middleware: 'auth'
})

const { orderRepo } = useRepositories()
const { getUserReviewForProduct } = useMarketplace()
const orders = ref<MarketplaceOrder[]>([])
const isLoadingOrders = ref(true)
const isOpeningTracking = ref(false)
const ordersError = ref('')
const trackingError = ref('')
const selectedTrackOrder = ref<MarketplaceOrder | null>(null)
const isTrackModalOpen = ref(false)

const selectedReviewProduct = ref<Product | null>(null)
const selectedReviewOrderId = ref<number | string | undefined>(undefined)
const isReviewModalOpen = ref(false)

const openReviewModal = (item: CartItem, orderId?: number | string) => {
  const reviewProduct = toReviewProduct(item)
  if (!reviewProduct) return
  selectedReviewProduct.value = reviewProduct
  selectedReviewOrderId.value = orderId
  isReviewModalOpen.value = true
}

const fetchOrders = async () => {
  isLoadingOrders.value = true
  ordersError.value = ''
  try {
    orders.value = await orderRepo.getOrders()
  } catch (err: any) {
    orders.value = []
    ordersError.value = err?.message || 'Could not load orders from the backend.'
  } finally {
    isLoadingOrders.value = false
  }
}

onMounted(fetchOrders)

const openTrackModal = async (order: MarketplaceOrder) => {
  trackingError.value = ''
  isOpeningTracking.value = true
  try {
    const latest = await orderRepo.getOrderById(getOrderNumber(order))
    if (!latest) {
      trackingError.value = 'This order could not be found in the backend.'
      return
    }

    const index = orders.value.findIndex(existing => getOrderNumber(existing) === getOrderNumber(latest))
    if (index !== -1) orders.value[index] = latest

    selectedTrackOrder.value = latest
    isTrackModalOpen.value = true
  } catch (err: any) {
    trackingError.value = err?.message || 'Could not load the latest tracking details.'
  } finally {
    isOpeningTracking.value = false
  }
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

        <UiAppAlert v-if="ordersError" class="mb-5">
          {{ ordersError }}
        </UiAppAlert>

        <UiAppAlert v-if="trackingError" class="mb-5">
          {{ trackingError }}
        </UiAppAlert>

        <div v-if="isLoadingOrders" class="space-y-4">
          <UiSkeletonCard v-for="index in 3" :key="index" />
        </div>

        <div v-else-if="orders.length" class="space-y-5">
          <UiAppCard
            v-for="order in orders"
            :key="getOrderNumber(order)"
          >
            <div class="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-3">
                  <p class="truncate text-lg font-semibold text-[#211f1d]">
                    {{ getOrderTitle(order) }}
                  </p>
                  <UiAppBadge variant="default">
                    {{ order.status }}
                  </UiAppBadge>
                </div>

                <p class="mt-1.5 text-sm text-[#756a60]">
                  Order #{{ getOrderNumber(order) }}
                  <span v-if="order.date"> - Placed on {{ order.date }}</span>
                </p>
              </div>

              <div class="flex flex-wrap items-center gap-2">
                <UiAppBadge variant="success">
                  {{ formatPaymentSummary(order) }}
                </UiAppBadge>

                <UiAppBadge variant="dark">
                  {{ formatPrice(order.total) }}
                </UiAppBadge>

                <UiAppButton
                  size="small"
                  class="ml-2"
                  :disabled="isOpeningTracking"
                  @click="openTrackModal(order)"
                >
                  {{ isOpeningTracking ? 'Loading...' : 'Track Order' }}
                </UiAppButton>
              </div>
            </div>

            <div class="mt-5 space-y-2 border-t border-[#ded6cc] pt-4">
              <p
                v-if="!getOrderItems(order).length"
                class="text-xs text-[#756a60]"
              >
                No item details available for this order.
              </p>
              <div
                v-for="item in getOrderItems(order)"
                :key="item.backendItemId || item.productId"
                class="flex items-center justify-between gap-4 rounded-lg border border-[#ded6cc] bg-[#f5f1e9] p-3"
              >
                <div class="flex min-w-0 items-center gap-3">
                  <img
                    v-if="item.productImage"
                    :src="item.productImage"
                    :alt="item.productName || 'Order item'"
                    class="h-12 w-10 shrink-0 rounded-md object-cover object-top"
                  />

                  <div class="min-w-0">
                    <p class="truncate text-sm font-medium text-[#211f1d]">
                      {{ item.productName || 'Product details unavailable' }}
                    </p>
                    <p class="text-xs text-[#756a60]">
                      Qty {{ item.quantity }}
                      <span v-if="item.productShop"> - {{ item.productShop }}</span>
                      <span v-if="item.price !== undefined"> - {{ formatPrice(item.price) }}</span>
                    </p>
                  </div>
                </div>

                <div class="flex items-center gap-3 shrink-0">
                  <p class="text-sm font-semibold text-[#211f1d]">
                    {{ item.lineTotal !== undefined ? formatPrice(item.lineTotal) : 'Amount unavailable' }}
                  </p>

                  <div v-if="order.status === 'Delivered'">
                    <template v-if="getUserReviewForProduct(item.productId, getOrderNumber(order))">
                      <UiAppBadge
                        v-if="getUserReviewForProduct(item.productId, getOrderNumber(order))?.status === 'pending'"
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
                      :disabled="!hasReviewableProduct(item)"
                      @click="openReviewModal(item, getOrderNumber(order))"
                    >
                      {{ hasReviewableProduct(item) ? 'Add Review' : 'Review unavailable' }}
                    </UiAppButton>
                  </div>
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

      <TrackShippingModal
        :is-open="isTrackModalOpen"
        :order="selectedTrackOrder"
        @close="isTrackModalOpen = false"
      />

      <AddReviewModal
        :is-open="isReviewModalOpen"
        :product="selectedReviewProduct"
        :order-id="selectedReviewOrderId"
        @close="isReviewModalOpen = false"
      />
    </div>
  </main>
</template>
