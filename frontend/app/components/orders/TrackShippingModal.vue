<script setup lang="ts">
import type { MarketplaceOrder } from '~/types/order'
import type { Product } from '~/types/product'
import { formatPrice, useMarketplace } from '~/composables/useMarketplace'

const props = defineProps<{
  isOpen: boolean
  order: MarketplaceOrder | null
}>()

const emit = defineEmits<{
  close: []
  openReview: [product: Product]
}>()

const { getOrderProducts, getUserReviewForProduct } = useMarketplace()

const trackingSteps = [
  { key: 'Ordered', title: 'Order Placed', desc: 'Order received & confirmed by seller' },
  { key: 'Shipped', title: 'Shipped / In Transit', desc: 'Package picked up & on its way' },
  { key: 'Delivered', title: 'Delivered', desc: 'Order delivered to buyer address' }
] as const

const getStepIndex = (status: string) => {
  const idx = trackingSteps.findIndex(s => s.key === status)
  return idx >= 0 ? idx : 0
}

const currentStepIndex = computed(() =>
  props.order ? getStepIndex(props.order.status) : 0
)

const orderItems = computed(() =>
  props.order ? getOrderProducts(props.order) : []
)
</script>

<template>
  <UiAppModal
    :is-open="isOpen && !!order"
    title="Shipping & Order Tracking"
    :kicker="order ? `Order #${order.id}` : ''"
    max-width="lg"
    @close="emit('close')"
  >
    <template v-if="order" #header>
      <p class="mt-2 text-sm text-[#756a60]">
        Placed on {{ order.date }} — Delivery to: <strong class="text-[#211f1d]">{{ order.deliveryAddress }}</strong>
      </p>
    </template>

    <template v-if="order">
      <!-- Tracking Pipeline Steps -->
      <div class="rounded-lg border border-[#ded6cc] bg-[#f5f1e9] p-5">
        <h3 class="text-xs font-semibold uppercase tracking-[0.14em] text-[#806344] mb-4">
          Current Delivery Progress
        </h3>

        <div class="relative flex flex-col sm:flex-row justify-between items-start gap-6">
          <div
            v-for="(step, index) in trackingSteps"
            :key="step.key"
            class="flex flex-1 items-start gap-3 relative z-10"
          >
            <div
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-xs font-bold transition-all"
              :class="
                index <= currentStepIndex
                  ? 'bg-[#806344] text-white shadow-md'
                  : 'border-2 border-[#cfc4b5] bg-white text-[#92877b]'
              "
            >
              {{ index <= currentStepIndex ? '✓' : index + 1 }}
            </div>

            <div>
              <p
                class="text-sm font-semibold tracking-wide"
                :class="index <= currentStepIndex ? 'text-[#211f1d]' : 'text-[#92877b]'"
              >
                {{ step.title }}
              </p>
              <p class="mt-0.5 text-xs text-[#756a60] leading-snug">
                {{ step.desc }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Detailed Product Items -->
      <div class="mt-6 space-y-4">
        <h3 class="text-sm font-semibold uppercase tracking-[0.14em] text-[#211f1d]">
          Items in this Shipment ({{ orderItems.length }})
        </h3>

        <div class="space-y-3 max-h-60 overflow-y-auto pr-1">
          <div
            v-for="item in orderItems"
            :key="item?.productId"
            class="flex items-center gap-4 rounded-lg border border-[#ded6cc] bg-white p-3.5 shadow-sm"
          >
            <img
              v-if="item"
              :src="item.product.image"
              :alt="item.product.name"
              class="h-16 w-14 rounded-md object-cover object-top"
            />

            <div v-if="item" class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-[#211f1d]">
                {{ item.product.name }}
              </p>
              <p class="mt-1 text-xs text-[#756a60]">
                Seller: {{ item.product.shop }} | Qty: {{ item.quantity }}
              </p>
              <p class="text-xs text-[#806344] font-medium mt-0.5">
                Unit Price: {{ formatPrice(item.product.price) }}
              </p>
            </div>

            <div v-if="item" class="text-right flex flex-col items-end gap-1.5">
              <p class="text-sm font-semibold text-[#211f1d]">
                {{ formatPrice(item.lineTotal) }}
              </p>

              <!-- Add Review or Review Status Badge if Delivered -->
              <div v-if="order.status === 'Delivered'">
                <template v-if="getUserReviewForProduct(item.product.id, order.id)">
                  <UiAppBadge
                    v-if="getUserReviewForProduct(item.product.id, order.id)?.status === 'pending'"
                    variant="warning"
                  >
                    <span>⏳</span> Under review
                  </UiAppBadge>
                  <UiAppBadge v-else variant="success">
                    <span>✓</span> Reviewed
                  </UiAppBadge>
                </template>
                <UiAppButton
                  v-else
                  variant="primary"
                  size="small"
                  class="h-8 px-3.5 text-xs"
                  @click="emit('openReview', item.product)"
                >
                  Add Review
                </UiAppButton>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer Total -->
      <div class="mt-6 flex items-center justify-between border-t border-[#ded6cc] pt-4">
        <div>
          <span class="text-xs font-medium uppercase tracking-[0.14em] text-[#756a60]">Payment Method</span>
          <p class="text-sm font-semibold text-[#211f1d]">{{ order.paymentMethod }} ({{ order.paymentStatus }})</p>
        </div>

        <div class="text-right">
          <span class="text-xs font-medium uppercase tracking-[0.14em] text-[#756a60]">Total Amount</span>
          <p class="font-serif text-2xl text-[#211f1d]">{{ formatPrice(order.total) }}</p>
        </div>
      </div>
    </template>
  </UiAppModal>
</template>
