<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'
import type { OrderStatus } from '~/types/order'

definePageMeta({
  middleware: 'auth'
})

const { orders, getOrderProducts, updateOrderStatus } = useMarketplace()
const { showToast } = useToast()
const statuses: OrderStatus[] = ['Ordered', 'Shipped', 'Delivered']

const handleStatusChange = (orderId: number, newStatus: OrderStatus) => {
  updateOrderStatus(orderId, newStatus)
  showToast(`Order #${orderId} updated to status "${newStatus}"`)
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-8">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 lg:flex-row">
      <AccountSidebar active="seller-orders" />

      <section class="min-w-0 flex-1">
        <div class="mb-8">
          <h1 class="font-serif text-4xl tracking-[-0.025em] text-[#211f1d]">
            Seller Orders
          </h1>

          <p class="mt-2 text-base text-[#756a60]">
            Process mock buyer orders and update delivery status.
          </p>
        </div>

        <div class="grid gap-5">
          <article
            v-for="order in orders"
            :key="order.id"
            class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-5 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-6"
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
                <span class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                  Order status
                </span>

                <select
                  :value="order.status"
                  class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#f5f1e9] px-4 text-sm text-[#211f1d] outline-none transition focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
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

            <div class="mt-6 overflow-hidden rounded-[8px] border border-[#ded6cc]">
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
                  class="font-medium text-[#211f1d]"
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

            <div class="mt-5 flex flex-wrap items-center justify-between gap-3 border-b border-[#ded6cc] pb-5">
              <div class="flex flex-wrap gap-2">
                <span class="rounded-full bg-[#211f1d] px-3 py-1 text-[11px] font-medium uppercase tracking-[0.12em] text-white">
                  {{ formatPrice(order.total) }}
                </span>

                <span class="rounded-full bg-[#e6eee5] px-3 py-1 text-[11px] font-medium uppercase tracking-[0.12em] text-[#536653]">
                  {{ order.paymentMethod }}
                </span>
              </div>
            </div>

            <div class="mt-5">
              <OrdersOrderTracker :status="order.status" />
            </div>
          </article>
        </div>
      </section>
    </div>
  </main>
</template>
