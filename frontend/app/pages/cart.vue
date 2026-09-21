<script setup lang="ts">
import type { CartProductItem } from '~/types/order'

definePageMeta({
  middleware: 'auth'
})

const {
  cartProducts,
  cartSubtotal,
  updateCartQuantity,
  removeFromCart
} = useMarketplace()
const { gtag } = useGtag()
const { track } = useAnalytics()
const deliveryFee = computed(() => (cartSubtotal.value > 0 ? 250 : 0))
const orderTotal = computed(() => cartSubtotal.value + deliveryFee.value)
if (import.meta.client && cartProducts.value.length) {
  gtag('event', 'view_cart', {
    currency: 'ETB',
    value: cartSubtotal.value,
    items: cartProducts.value.map(item => ({
      item_id: String(item.product.id),
      item_name: item.product.name,
      item_category: item.product.category,
      price: Number(item.product.price),
      quantity: item.quantity
    }))
  })
}

const cartError = ref<string | null>(null)

const handleDecreaseQuantity = async (item: CartProductItem) => {
  cartError.value = null
  try {
    const newQty = item.quantity - 1
    await updateCartQuantity(item.productId, newQty)
    track('cart_quantity_updated', {
      product_id: String(item.product.id),
      product_name: item.product.name,
      quantity: newQty
    })
  } catch (err: any) {
    cartError.value = err?.message || 'Failed to update cart'
  }
}

const handleIncreaseQuantity = async (item: CartProductItem) => {
  cartError.value = null
  try {
    const newQty = item.quantity + 1
    await updateCartQuantity(item.productId, newQty)
    track('cart_quantity_updated', {
      product_id: String(item.product.id),
      product_name: item.product.name,
      quantity: newQty
    })
  } catch (err: any) {
    cartError.value = err?.message || 'Failed to update cart'
  }
}

const handleRemoveFromCart = async (item: CartProductItem) => {
  cartError.value = null
  try {
    await removeFromCart(item.productId)
    track('cart_item_removed', {
      product_id: String(item.product.id),
      product_name: item.product.name,
      quantity: item.quantity
    })
  } catch (err: any) {
    cartError.value = err?.message || 'Failed to remove item'
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <section class="mx-auto max-w-7xl">
      <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
            Checkout
          </p>

          <h1 class="mt-2 font-serif text-4xl tracking-[-0.025em] text-[#211f1d] sm:text-5xl">
            Shopping cart
          </h1>
        </div>

        <NuxtLink
          to="/products"
          class="text-sm font-medium text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline"
        >
          Continue shopping
        </NuxtLink>
      </div>

      <!-- Error Message Banner -->
      <div
        v-if="cartError"
        class="mb-6 rounded-lg border border-red-300 bg-red-50 p-4 text-sm text-red-800"
      >
        {{ cartError }}
      </div>

      <div
        v-if="cartProducts.length"
        class="grid gap-8 lg:grid-cols-[1fr_380px]"
      >
        <section class="space-y-4">
          <UiAppCard
            v-for="item in cartProducts"
            :key="item?.productId"
            class="p-4 sm:p-5"
          >
            <div
              v-if="item"
              class="flex gap-3.5 sm:gap-5 items-start"
            >
              <NuxtLink
                :to="`/products/${item.product.id}`"
                class="h-28 w-24 sm:h-36 sm:w-28 shrink-0 overflow-hidden rounded-lg bg-[#eee8df]"
              >
                <img
                  :src="item.product.image"
                  :alt="item.product.name"
                  class="h-full w-full object-cover object-top"
                />
              </NuxtLink>

              <div class="min-w-0 flex-1">
                <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-1 sm:gap-4">
                  <div class="min-w-0">
                    <p class="text-[10px] sm:text-xs font-medium uppercase tracking-[0.14em] text-[#806344]">
                      {{ item.product.shop }}
                    </p>

                    <NuxtLink
                      :to="`/products/${item.product.id}`"
                      class="mt-0.5 sm:mt-1 block font-serif text-lg sm:text-2xl text-[#211f1d] transition hover:text-[#806344] truncate"
                    >
                      {{ item.product.name }}
                    </NuxtLink>

                    <p class="mt-1 text-xs sm:text-sm text-[#756a60]">
                      {{ item.product.category }}<span v-if="item.product.subCategory"> · {{ item.product.subCategory }}</span> / {{ item.product.stock }} available
                    </p>
                  </div>

                  <p class="text-base sm:text-lg font-semibold text-[#211f1d]">
                    {{ formatPrice(item.lineTotal) }}
                  </p>
                </div>

                <div class="mt-4 flex flex-wrap items-center justify-between gap-3">
                  <div class="flex items-center rounded-full border border-[#cfc4b5] bg-[#f5f1e9] p-1">
                    <button
                      type="button"
                      aria-label="Decrease quantity"
                      class="flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-full text-base sm:text-lg text-[#5d4b37] transition hover:bg-[#eee8df]"
                      @click="handleDecreaseQuantity(item)"
                    >
                      -
                    </button>

                    <span class="flex h-8 min-w-10 sm:h-9 sm:min-w-12 items-center justify-center text-xs sm:text-sm font-medium text-[#211f1d]">
                      {{ item.quantity }}
                    </span>

                    <button
                      type="button"
                      aria-label="Increase quantity"
                      class="flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-full text-base sm:text-lg text-[#5d4b37] transition hover:bg-[#eee8df]"
                      @click="handleIncreaseQuantity(item)"
                    >
                      +
                    </button>
                  </div>

                  <button
                    type="button"
                    class="text-xs sm:text-sm font-medium text-red-600 underline-offset-4 transition hover:underline"
                    @click="handleRemoveFromCart(item)"
                  >
                    Remove
                  </button>
                </div>
              </div>
            </div>
          </UiAppCard>
        </section>

        <UiAppCard as="aside" class="h-fit">
          <h2 class="font-serif text-3xl text-[#211f1d]">
            Summary
          </h2>

          <div class="mt-6 space-y-4 border-b border-[#ded6cc] pb-6 text-sm">
            <div class="flex justify-between gap-4 text-[#665c53]">
              <span>Subtotal</span>
              <span>{{ formatPrice(cartSubtotal) }}</span>
            </div>

            <div class="flex justify-between gap-4 text-[#665c53]">
              <span>Delivery</span>
              <span>{{ formatPrice(deliveryFee) }}</span>
            </div>
          </div>

          <div class="mt-5 flex justify-between gap-4 text-lg font-semibold text-[#211f1d]">
            <span>Total</span>
            <span>{{ formatPrice(orderTotal) }}</span>
          </div>

          <div class="mt-7">
            <UiAppButton
              to="/checkout"
              variant="secondary"
              class="w-full"
            >
              Checkout
            </UiAppButton>
          </div>
        </UiAppCard>
      </div>

      <UiAppEmptyState
        v-else
        title="Your cart is empty"
        description="Add a few products and your checkout summary will appear here."
        action-label="Browse products"
        action-to="/products"
      />
    </section>
  </main>
</template>
