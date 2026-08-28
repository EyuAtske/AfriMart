<script setup lang="ts">
const {
  cartProducts,
  cartSubtotal,
  updateCartQuantity,
  removeFromCart
} = useMarketplace()

const deliveryFee = computed(() => (cartSubtotal.value > 0 ? 250 : 0))
const orderTotal = computed(() => cartSubtotal.value + deliveryFee.value)
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

      <div
        v-if="cartProducts.length"
        class="grid gap-8 lg:grid-cols-[1fr_380px]"
      >
        <section class="space-y-4">
          <article
            v-for="item in cartProducts"
            :key="item?.productId"
            class="rounded-[10px] border border-[#d9d0c4] bg-[#faf8f4] p-4 shadow-[0_20px_70px_rgba(33,31,29,0.05)] sm:p-5"
          >
            <div
              v-if="item"
              class="grid gap-5 sm:grid-cols-[140px_1fr]"
            >
              <NuxtLink
                :to="`/products/${item.product.id}`"
                class="h-44 overflow-hidden rounded-[8px] bg-[#eee8df] sm:h-36"
              >
                <img
                  :src="item.product.image"
                  :alt="item.product.name"
                  class="h-full w-full object-cover object-top"
                />
              </NuxtLink>

              <div class="min-w-0">
                <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                  <div>
                    <p class="text-xs font-medium uppercase tracking-[0.14em] text-[#806344]">
                      {{ item.product.shop }}
                    </p>

                    <NuxtLink
                      :to="`/products/${item.product.id}`"
                      class="mt-1 block font-serif text-2xl text-[#211f1d] transition hover:text-[#806344]"
                    >
                      {{ item.product.name }}
                    </NuxtLink>

                    <p class="mt-2 text-sm leading-6 text-[#756a60]">
                      {{ item.product.category }} / {{ item.product.stock }} available
                    </p>
                  </div>

                  <p class="text-lg font-semibold text-[#211f1d]">
                    {{ formatPrice(item.lineTotal) }}
                  </p>
                </div>

                <div class="mt-5 flex flex-wrap items-center justify-between gap-4">
                  <div class="flex items-center rounded-full border border-[#cfc4b5] bg-[#f5f1e9] p-1">
                    <button
                      type="button"
                      class="flex h-9 w-9 items-center justify-center rounded-full text-lg text-[#5d4b37] transition hover:bg-[#eee8df]"
                      @click="updateCartQuantity(item.productId, item.quantity - 1)"
                    >
                      -
                    </button>

                    <span class="flex h-9 min-w-12 items-center justify-center text-sm font-medium text-[#211f1d]">
                      {{ item.quantity }}
                    </span>

                    <button
                      type="button"
                      class="flex h-9 w-9 items-center justify-center rounded-full text-lg text-[#5d4b37] transition hover:bg-[#eee8df]"
                      @click="updateCartQuantity(item.productId, item.quantity + 1)"
                    >
                      +
                    </button>
                  </div>

                  <button
                    type="button"
                    class="text-sm font-medium text-red-600 underline-offset-4 transition hover:underline"
                    @click="removeFromCart(item.productId)"
                  >
                    Remove
                  </button>
                </div>
              </div>
            </div>
          </article>
        </section>

        <aside class="h-fit rounded-[10px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)]">
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

          <NuxtLink
            to="/checkout"
            class="mt-7 inline-flex h-12 w-full items-center justify-center rounded-full bg-[#211f1d] px-6 text-sm font-medium uppercase tracking-[0.14em] text-white transition hover:bg-[#3b3733]"
          >
            Checkout
          </NuxtLink>
        </aside>
      </div>

      <section
        v-else
        class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-10 text-center shadow-[0_20px_70px_rgba(33,31,29,0.06)]"
      >
        <h2 class="font-serif text-3xl text-[#211f1d]">
          Your cart is empty
        </h2>

        <p class="mt-3 text-sm leading-6 text-[#756a60]">
          Add a few products and your checkout summary will appear here.
        </p>

        <NuxtLink
          to="/products"
          class="mt-6 inline-flex h-12 items-center justify-center rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
        >
          Browse products
        </NuxtLink>
      </section>
    </section>
  </main>
</template>
