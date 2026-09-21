<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'
import type { OrderStatus } from '~/types/order'

definePageMeta({
  middleware: 'auth'
})

const { orders, getOrderProducts, updateOrderStatus, fetchSellerOrders } = useMarketplace()
const { showToast } = useToast()
const statuses: OrderStatus[] = ['Pending', 'Confirmed', 'Processing', 'Shipped', 'Delivered', 'Cancelled']

onMounted(() => {
  fetchSellerOrders()
})

const handleStatusChange = async (orderId: number | string, newStatus: OrderStatus) => {
  try {
    await updateOrderStatus(orderId, newStatus)
    showToast(`Order status updated to "${newStatus}"`)
  } catch (err: any) {
    showToast(err?.message || 'Failed to update order status', 'error')
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 lg:flex-row">
      <AccountSidebar active="seller-orders" />

      <section class="min-w-0 flex-1">
        <div class="mb-8">
          <h1 class="font-serif text-4xl tracking-[-0.025em] text-[#211f1d]">
            Seller Orders
          </h1>

          <p class="mt-2 text-base text-[#756a60]">
            Manage buyer orders and update delivery statuses.
          </p>
        </div>

        <div class="grid gap-5">
          <UiAppCard
            v-for="order in orders"
            :key="order.id"
          >
            <div class="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
              <div>
                <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
                  Order #{{ order.id }}
                </p>

                <h2 class="mt-2 font-serif text-3xl text-[#211f1d]">
                  {{ order.buyerName }}
                </h2>

                <p class="mt-2 text-sm leading-6 text-[#756a60]">
                  {{ order.deliveryAddress }} / {{ order.phone }}
                </p>
              </div>

              <label class="w-full space-y-2 sm:w-60">
                <span class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                  Order status
                </span>

                <select
                  :value="order.status"
                  class="h-12 w-full rounded-md border border-[#cfc4b5] bg-[#f5f1e9] px-4 text-sm text-[#211f1d] outline-none transition hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                  @change="handleStatusChange(order.id, ($event.target as HTMLSelectElement).value as OrderStatus)"
                >
                  <option
                    v-for="status in statuses"
                    :key="status"
                    :value="status"
                  >
                    {{ status }}
                  </option>
                </select>
              </label>
            </div>

            <div class="mt-6 overflow-x-auto rounded-lg border border-[#ded6cc]">
              <div class="min-w-[320px]">
                <div class="grid grid-cols-[1fr_90px_110px] bg-[#eee8df] px-4 py-3 text-xs font-medium uppercase tracking-[0.12em] text-[#665c53]">
                  <span>Product</span>
                  <span>Qty</span>
                  <span class="text-right">Total</span>
                </div>

                <div
                  v-for="item in getOrderProducts(order)"
                  :key="item?.productId"
                  class="grid grid-cols-[1fr_90px_110px] items-center border-t border-[#ded6cc] px-4 py-3 text-sm"
                >
                  <span
                    v-if="item"
                    class="font-medium text-[#211f1d] truncate pr-2"
                  >
                    {{ item.product.name }}
                  </span>

                  <span
                    v-if="item"
                    class="text-[#756a60]"
                  >
                    {{ item.quantity }}
                  </span>

                  <span
                    v-if="item"
                    class="text-right font-semibold text-[#211f1d]"
                  >
                    {{ formatPrice(item.lineTotal) }}
                  </span>
                </div>
              </div>
            </div>

            <div class="mt-5 flex flex-wrap items-center justify-between gap-3 border-b border-[#ded6cc] pb-5">
              <div class="flex flex-wrap gap-2">
                <UiAppBadge variant="dark">
                  {{ formatPrice(order.total) }}
                </UiAppBadge>

                <UiAppBadge variant="success">
                  {{ order.paymentMethod }}
                </UiAppBadge>
              </div>
            </div>

            <div class="mt-5">
              <OrdersOrderTracker :status="order.status" />
            </div>
          </UiAppCard>
        </div>
      </section>
    </div>
  </main>
</template>
