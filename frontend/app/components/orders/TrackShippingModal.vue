<script setup lang="ts">
import type { MarketplaceOrder, OrderStatus } from '~/types/order'
import { formatPrice } from '~/composables/useMarketplace'
import { formatPaymentSummary } from '~/utils/orderDisplay'

const props = defineProps<{
  isOpen: boolean
  order: MarketplaceOrder | null
}>()

const emit = defineEmits<{
  close: []
}>()

const trackingSteps: { key: OrderStatus; title: string; desc: string }[] = [
  { key: 'Pending', title: 'Order Placed', desc: 'Order received and waiting for seller confirmation' },
  { key: 'Confirmed', title: 'Confirmed', desc: 'Seller confirmed the order' },
  { key: 'Processing', title: 'Processing', desc: 'Seller is preparing the shipment' },
  { key: 'Shipped', title: 'Shipped / In Transit', desc: 'Package picked up and on its way' },
  { key: 'Delivered', title: 'Delivered', desc: 'Order delivered to buyer address' }
]

const getStepIndex = (status: OrderStatus) => {
  const idx = trackingSteps.findIndex(step => step.key === status)
  return idx >= 0 ? idx : 0
}

const currentStepIndex = computed(() =>
  props.order ? getStepIndex(props.order.status) : 0
)

const isCancelled = computed(() => props.order?.status === 'Cancelled')
const orderItems = computed(() => props.order?.items || [])
const orderNumber = computed(() => props.order ? props.order.backendId || String(props.order.id) : '')
const deliveryAddress = computed(() => {
  if (!props.order) return ''
  return [props.order.deliveryAddress, props.order.deliveryCity].filter(Boolean).join(', ')
})
const paymentLabel = computed(() => {
  if (!props.order) return 'Payment details unavailable'
  return formatPaymentSummary(props.order)
})
const amountPaid = computed(() => props.order?.paymentAmount)
</script>

<template>
  <UiAppModal
    :is-open="isOpen && !!order"
    title="Shipping & Order Tracking"
    :kicker="order ? `Order #${orderNumber}` : ''"
    max-width="lg"
    @close="emit('close')"
  >
    <template v-if="order" #header>
      <p class="mt-2 text-sm text-[#756a60]">
        <span v-if="order.date">Placed on {{ order.date }}</span>
        <span v-if="order.date && deliveryAddress"> - </span>
        <span v-if="deliveryAddress">
          Delivery to: <strong class="text-[#211f1d]">{{ deliveryAddress }}</strong>
        </span>
      </p>
    </template>

    <template v-if="order">
      <div
        v-if="isCancelled"
        class="rounded-lg border border-red-200 bg-red-50 p-5 text-red-800"
      >
        <h3 class="text-sm font-semibold">Order Cancelled</h3>
        <p class="mt-1 text-sm">
          This order was cancelled. The status is coming from the backend order record.
        </p>
      </div>

      <div v-else class="rounded-lg border border-[#ded6cc] bg-[#f5f1e9] p-5">
        <h3 class="mb-4 text-xs font-semibold uppercase tracking-[0.14em] text-[#806344]">
          Current Delivery Progress
        </h3>

        <div class="relative grid gap-5 md:grid-cols-5">
          <div
            v-for="(step, index) in trackingSteps"
            :key="step.key"
            class="flex items-start gap-3 md:flex-col"
          >
            <div
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-xs font-bold transition-all"
              :class="
                index <= currentStepIndex
                  ? 'bg-[#806344] text-white shadow-md'
                  : 'border-2 border-[#cfc4b5] bg-white text-[#92877b]'
              "
            >
              {{ index < currentStepIndex ? '✓' : index + 1 }}
            </div>

            <div>
              <p
                class="text-sm font-semibold tracking-wide"
                :class="index <= currentStepIndex ? 'text-[#211f1d]' : 'text-[#92877b]'"
              >
                {{ step.title }}
              </p>
              <p class="mt-0.5 text-xs leading-snug text-[#756a60]">
                {{ step.desc }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-6 space-y-4">
        <h3 class="text-sm font-semibold uppercase tracking-[0.14em] text-[#211f1d]">
          Items in this Shipment ({{ orderItems.length }})
        </h3>

        <div class="max-h-60 space-y-3 overflow-y-auto pr-1">
          <div
            v-for="item in orderItems"
            :key="item.backendItemId || item.productId"
            class="flex items-center gap-4 rounded-lg border border-[#ded6cc] bg-white p-3.5 shadow-sm"
          >
            <img
              v-if="item.productImage"
              :src="item.productImage"
              :alt="item.productName || 'Order item'"
              class="h-16 w-14 rounded-md object-cover object-top"
            />

            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-[#211f1d]">
                {{ item.productName || 'Product details unavailable' }}
              </p>
              <p class="mt-1 text-xs text-[#756a60]">
                <span v-if="item.productShop">Seller: {{ item.productShop }} | </span>Qty: {{ item.quantity }}
              </p>
              <p v-if="item.price !== undefined" class="mt-0.5 text-xs font-medium text-[#806344]">
                Unit Price: {{ formatPrice(item.price) }}
              </p>
            </div>

            <div class="text-right">
              <p class="text-sm font-semibold text-[#211f1d]">
                {{ item.lineTotal !== undefined ? formatPrice(item.lineTotal) : 'Amount unavailable' }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-6 grid gap-4 border-t border-[#ded6cc] pt-4 sm:grid-cols-3">
        <div>
          <span class="text-xs font-medium uppercase tracking-[0.14em] text-[#756a60]">Payment</span>
          <p class="text-sm font-semibold text-[#211f1d]">{{ paymentLabel }}</p>
          <p v-if="order.transactionId" class="mt-1 text-xs text-[#756a60]">
            Transaction {{ order.transactionId }}
          </p>
        </div>

        <div>
          <span class="text-xs font-medium uppercase tracking-[0.14em] text-[#756a60]">Amount Paid</span>
          <p class="text-sm font-semibold text-[#211f1d]">
            {{ amountPaid !== undefined ? formatPrice(amountPaid) : 'Unavailable' }}
          </p>
        </div>

        <div class="sm:text-right">
          <span class="text-xs font-medium uppercase tracking-[0.14em] text-[#756a60]">Total Amount</span>
          <p class="font-serif text-2xl text-[#211f1d]">{{ formatPrice(order.total) }}</p>
        </div>
      </div>
    </template>
  </UiAppModal>
</template>
