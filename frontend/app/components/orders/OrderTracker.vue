<script setup lang="ts">
import type { OrderStatus } from '~/types/order'

const props = defineProps<{
  status: OrderStatus
}>()

const steps: { name: OrderStatus; label: string; desc: string }[] = [
  { name: 'Pending', label: 'Order Placed', desc: 'Waiting for seller confirmation' },
  { name: 'Confirmed', label: 'Confirmed', desc: 'Seller confirmed the order' },
  { name: 'Processing', label: 'Processing', desc: 'Preparing shipment' },
  { name: 'Shipped', label: 'Shipped / In Transit', desc: 'In transit with courier' },
  { name: 'Delivered', label: 'Delivered', desc: 'Package delivered safely' }
]

const getStepIndex = (status: OrderStatus) => {
  switch (status) {
    case 'Confirmed': return 1
    case 'Processing': return 2
    case 'Shipped': return 3
    case 'Delivered': return 4
    default: return 0
  }
}

const currentStepIndex = computed(() => getStepIndex(props.status))
const isCancelled = computed(() => props.status === 'Cancelled')
</script>

<template>
  <div class="rounded-[10px] border border-[#d9d0c4] bg-[#faf8f4] p-5 sm:p-6">
    <p class="mb-4 text-xs font-medium uppercase tracking-[0.16em] text-[#806344]">
      Delivery Progress Timeline
    </p>

    <div
      v-if="isCancelled"
      class="rounded-lg border border-red-200 bg-red-50 p-4 text-sm font-medium text-red-800"
    >
      Order cancelled
    </div>

    <div v-else class="relative flex flex-col justify-between gap-6 sm:flex-row sm:items-center sm:gap-0">
      <div class="absolute left-8 right-8 top-4 -z-0 hidden h-[2px] bg-[#e0d6c9] sm:block">
        <div
          class="h-full bg-[#806344] transition-all duration-500"
          :style="{ width: `${(currentStepIndex / (steps.length - 1)) * 100}%` }"
        />
      </div>

      <div
        v-for="(step, index) in steps"
        :key="step.name"
        class="relative z-10 flex items-center gap-4 text-left sm:w-1/5 sm:flex-col sm:gap-2 sm:text-center"
      >
        <div
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full border-2 text-xs font-bold transition-all duration-300"
          :class="[
            index <= currentStepIndex
              ? 'border-[#806344] bg-[#806344] text-white shadow-md'
              : 'border-[#cfc4b5] bg-[#f5f1e9] text-[#756a60]'
          ]"
        >
          <span v-if="index < currentStepIndex">✓</span>
          <span v-else>{{ index + 1 }}</span>
        </div>

        <div>
          <p
            class="text-sm font-semibold tracking-wide"
            :class="index <= currentStepIndex ? 'text-[#211f1d]' : 'text-[#92877b]'"
          >
            {{ step.label }}
          </p>

          <p class="mt-0.5 max-w-[140px] text-xs text-[#756a60] sm:mx-auto">
            {{ step.desc }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
