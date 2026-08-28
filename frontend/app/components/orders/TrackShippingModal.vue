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
  <Teleport to="body">
    <div
      v-if="isOpen && order"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm transition-opacity"
      @click.self="emit('close')"
    >
      <div class="relative w-full max-w-2xl max-h-[90vh] overflow-y-auto rounded-[16px] border border-[#d9d0c4] bg-[#faf8f4] p-6 sm:p-8 shadow-[0_25px_80px_rgba(0,0,0,0.25)]">
        <!-- Close Button -->
        <button
          type="button"
          class="absolute top-5 right-5 flex h-9 w-9 items-center justify-center rounded-full bg-[#eee8df] text-[#665c53] transition hover:bg-[#211f1d] hover:text-white"
          @click="emit('close')"
        >
          <span class="sr-only">Close modal</span>
          ✕
        </button>

        <!-- Header -->
        <div class="border-b border-[#ded6cc] pb-5">
          <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
            Order #{{ order.id }}
          </p>
          <h2 class="mt-1 font-serif text-3xl text-[#211f1d]">
            Shipping & Order Tracking
          </h2>
          <p class="mt-2 text-sm text-[#756a60]">
            Placed on {{ order.date }} — Delivery to: <strong class="text-[#211f1d]">{{ order.deliveryAddress }}</strong>
          </p>
        </div>

        <!-- Tracking Pipeline Steps -->
        <div class="my-6 rounded-[12px] border border-[#ded6cc] bg-[#f5f1e9] p-5">
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
                  class="text-sm font-semibold"
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
        <div class="space-y-4">
          <h3 class="text-sm font-semibold uppercase tracking-[0.14em] text-[#211f1d]">
            Items in this Shipment ({{ orderItems.length }})
          </h3>

          <div class="space-y-3 max-h-60 overflow-y-auto pr-1">
            <div
              v-for="item in orderItems"
              :key="item?.productId"
              class="flex items-center gap-4 rounded-[10px] border border-[#ded6cc] bg-white p-3.5 shadow-sm"
            >
              <img
                v-if="item"
                :src="item.product.image"
                :alt="item.product.name"
                class="h-16 w-14 rounded-[6px] object-cover object-top"
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
                    <span
                      v-if="getUserReviewForProduct(item.product.id, order.id)?.status === 'pending'"
                      class="inline-flex items-center gap-1 rounded-full bg-amber-100 px-3 py-1 text-[11px] font-medium text-amber-800"
                    >
                      <span>⏳</span> Under review
                    </span>
                    <span
                      v-else
                      class="inline-flex items-center gap-1 rounded-full bg-green-100 px-3 py-1 text-[11px] font-medium text-green-800"
                    >
                      <span>✓</span> Reviewed
                    </span>
                  </template>
                  <button
                    v-else
                    type="button"
                    class="inline-flex h-8 items-center justify-center rounded-full bg-[#806344] px-3.5 text-xs font-medium uppercase tracking-[0.12em] text-white transition hover:bg-[#5d4b37]"
                    @click="emit('openReview', item.product)"
                  >
                    Add Review
                  </button>
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
      </div>
    </div>
  </Teleport>
</template>
