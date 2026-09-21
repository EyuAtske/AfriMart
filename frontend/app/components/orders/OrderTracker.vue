<script setup lang="ts">
import type { OrderStatus } from '~/types/order'

const props = defineProps<{
  status: OrderStatus
}>()

const steps: { name: OrderStatus; label: string; desc: string }[] = [
  { name: 'Ordered', label: 'Order Placed', desc: 'Order received & confirmed' },
  { name: 'Shipped', label: 'Shipped', desc: 'In transit with courier' },
  { name: 'Delivered', label: 'Delivered', desc: 'Package delivered safely' }
]

const getStepIndex = (status: OrderStatus) => {
  switch (status) {
    case 'Ordered': return 0
    case 'Shipped': return 1
    case 'Delivered': return 2
    default: return 0
  }
}

const currentStepIndex = computed(() => getStepIndex(props.status))
</script>

<template>
  <div class="rounded-[10px] border border-[#d9d0c4] bg-[#faf8f4] p-5 sm:p-6">
    <p class="text-xs font-medium uppercase tracking-[0.16em] text-[#806344] mb-4">
      Delivery Progress Timeline
    </p>

    <div class="relative flex flex-col sm:flex-row sm:items-center justify-between gap-6 sm:gap-0">
      <!-- Line behind steps (Desktop) -->
      <div class="hidden sm:block absolute left-8 right-8 top-4 h-[2px] bg-[#e0d6c9] -z-0">
        <div
          class="h-full bg-[#806344] transition-all duration-500"
          :style="{ width: `${(currentStepIndex / 2) * 100}%` }"
        />
      </div>

      <div
        v-for="(step, index) in steps"
        :key="step.name"
        class="relative z-10 flex sm:flex-col items-center gap-4 sm:gap-2 text-left sm:text-center sm:w-1/3"
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

          <p class="text-xs text-[#756a60] mt-0.5 max-w-[140px] sm:mx-auto">
            {{ step.desc }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
