<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const router = useRouter()
const { user } = useAuth()
const { cartProducts, cartSubtotal, createOrder, retryPostCheckoutRefresh, isOwnProduct } = useMarketplace()
const { gtag } = useGtag()
const { $posthog } = useNuxtApp()
const checkoutForm = reactive({
  fullName: '',
  phone: '',
  address: '',
  deliveryCity: 'Addis Ababa',
  deliveryNotes: ''
})

const isSubmitting = ref(false)
const errorMessage = ref('')
const orderSuccess = ref(false)
const refreshWarning = ref('')
const createdOrderId = ref('')
const isRetrying = ref(false)

const deliveryFee = computed(() => (cartSubtotal.value > 0 ? 250 : 0))
const orderTotal = computed(() => cartSubtotal.value + deliveryFee.value)

if (import.meta.client && cartProducts.value.length) {
  gtag('event', 'begin_checkout', {
    currency: 'ETB',
    value: orderTotal.value,
    items: cartProducts.value.map(item => ({
      item_id: String(item.product.id),
      item_name: item.product.name,
      item_category: item.product.category,
      price: Number(item.product.price),
      quantity: item.quantity
    }))
  })
}

watch(
  user,
  (currentUser) => {
    if (currentUser.name || currentUser.username) {
      checkoutForm.fullName = currentUser.name || currentUser.username
    }
  },
  { immediate: true }
)

const retryRefresh = async () => {
  isRetrying.value = true
  try {
    await retryPostCheckoutRefresh()
    refreshWarning.value = ''
    await navigateTo(`/orders?highlight=${createdOrderId.value}`)
  } catch (err: any) {
    refreshWarning.value = err?.message || 'Still unable to refresh data. Your order has been placed — check your order history.'
  } finally {
    isRetrying.value = false
  }
}

const handlePlaceOrder = async () => {
  errorMessage.value = ''
  refreshWarning.value = ''
  orderSuccess.value = false

  const selfItem = cartProducts.value.find(item => item && isOwnProduct(item.product))
  if (selfItem) {
    errorMessage.value = `You cannot purchase "${selfItem.product.name}" because it is listed by your own shop.`
    return
  }

  const soldOutItem = cartProducts.value.find(item => item && item.product.stock < item.quantity)
  if (soldOutItem) {
    errorMessage.value = `"${soldOutItem.product.name}" is sold out or has insufficient stock (${soldOutItem.product.stock} available).`
    return
  }

  if (!checkoutForm.fullName.trim()) {
    errorMessage.value = 'Please enter recipient name for delivery.'
    return
  }
  if (!checkoutForm.phone.trim()) {
    errorMessage.value = 'Please enter a contact phone number.'
    return
  }
  if (!checkoutForm.address.trim()) {
    errorMessage.value = 'Please enter delivery address.'
    return
  }
  if (!checkoutForm.deliveryCity.trim()) {
    errorMessage.value = 'Please enter delivery city.'
    return
  }

  isSubmitting.value = true
  try {
    const result = await createOrder({
      buyerName: checkoutForm.fullName.trim(),
      phone: checkoutForm.phone.trim(),
      deliveryAddress: checkoutForm.address.trim(),
      deliveryCity: checkoutForm.deliveryCity.trim(),
      deliveryNotes: checkoutForm.deliveryNotes.trim()
    })

    if (result?.order) {
      createdOrderId.value = result.order.backendId || String(result.order.id)
      if (import.meta.client) {
      $posthog?.capture('purchase_completed', {
      order_id: createdOrderId.value,
      currency: 'ETB',
      value: orderTotal.value,
      items: cartProducts.value.map(item => ({
      product_id: String(item.product.id),
      product_name: item.product.name,
      category: item.product.category,
      price: Number(item.product.price),
      quantity: item.quantity
    }))
  })
}

      if (result.refreshError) {
        // Order confirmed, but data refresh failed — stay on page with warning
        orderSuccess.value = true
        refreshWarning.value = result.refreshError
      } else {
        await navigateTo(`/orders?highlight=${createdOrderId.value}`)
      }
    } else {
      throw new Error('Order creation failed. Please try again.')
    }
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to place order. Please try again.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <section class="mx-auto max-w-7xl">
      <div class="mb-8">
        <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
          Order creation
        </p>

        <h1 class="mt-2 font-serif text-4xl tracking-[-0.025em] text-[#211f1d] sm:text-5xl">
          Checkout
        </h1>
      </div>

      <!-- Error Message Banner -->
      <div
        v-if="errorMessage"
        class="mb-6 rounded-lg border border-red-300 bg-red-50 p-4 text-red-800"
      >
        <div class="flex items-center gap-3">
          <span class="text-xl">⚠️</span>
          <p class="text-sm font-medium">{{ errorMessage }}</p>
        </div>
      </div>

      <!-- Order placed with data-refresh warning -->
      <div
        v-if="orderSuccess"
        class="mb-6 space-y-3"
      >
        <div class="rounded-lg border border-green-300 bg-green-50 p-4 text-green-800">
          <div class="flex items-center gap-3">
            <span class="text-xl">✅</span>
            <p class="text-sm font-medium">Your order has been placed successfully!</p>
          </div>
        </div>

        <div
          v-if="refreshWarning"
          class="rounded-lg border border-amber-300 bg-amber-50 p-4 text-amber-800"
        >
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div class="flex items-center gap-3 min-w-0">
              <span class="shrink-0 text-xl">⚠️</span>
              <p class="text-sm font-medium">{{ refreshWarning }}</p>
            </div>
            <div class="flex shrink-0 items-center gap-2">
              <button
                :disabled="isRetrying"
                class="inline-flex h-9 items-center justify-center rounded-lg border border-amber-400 bg-white px-4 text-xs font-medium text-amber-800 transition hover:bg-amber-100 disabled:opacity-50"
                @click="retryRefresh"
              >
                {{ isRetrying ? 'Retrying…' : 'Retry refresh' }}
              </button>
              <NuxtLink
                :to="`/orders?highlight=${createdOrderId}`"
                class="inline-flex h-9 items-center justify-center rounded-lg bg-[#806344] px-4 text-xs font-medium text-white transition hover:bg-[#6b5438]"
              >
                View orders
              </NuxtLink>
            </div>
          </div>
        </div>

        <div v-else class="flex justify-center pt-2">
          <NuxtLink
            :to="`/orders?highlight=${createdOrderId}`"
            class="inline-flex h-12 items-center justify-center rounded-lg bg-[#806344] px-8 text-sm font-medium text-white transition hover:bg-[#6b5438]"
          >
            View your orders →
          </NuxtLink>
        </div>
      </div>

      <div
        v-else-if="cartProducts.length"
        class="grid gap-8 lg:grid-cols-[1fr_380px]"
      >
        <div class="rounded-xl border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
          <h2 class="font-serif text-3xl text-[#211f1d]">
            Delivery Details
          </h2>

          <div class="mt-6 grid gap-5 sm:grid-cols-2">
            <AuthInput
              v-model="checkoutForm.fullName"
              label="Recipient Name"
              placeholder="Full name for delivery"
              name="checkout-name"
            />

            <AuthInput
              v-model="checkoutForm.phone"
              label="Phone Number"
              type="tel"
              placeholder="+251911234567"
              name="checkout-phone"
            />
          </div>

          <div class="mt-5 grid gap-5 sm:grid-cols-[1fr_220px]">
            <AuthInput
              v-model="checkoutForm.address"
              label="Street Address"
              placeholder="Street, building, house no."
              name="checkout-address"
            />

            <AuthInput
              v-model="checkoutForm.deliveryCity"
              label="City"
              placeholder="Addis Ababa"
              name="checkout-city"
            />
          </div>

          <div class="mt-5">
            <AuthInput
              v-model="checkoutForm.deliveryNotes"
              label="Delivery Notes (Optional)"
              placeholder="Special instructions for driver..."
              name="checkout-notes"
            />
          </div>

          <section class="mt-8">
            <h2 class="font-serif text-3xl text-[#211f1d]">
              Payment Method
            </h2>

            <div class="mt-4 space-y-3">
              <!-- Cash on Delivery (Enabled & Default) -->
              <label class="flex items-center justify-between rounded-lg border-2 border-[#806344] bg-[#f5f1e9] p-4 cursor-pointer">
                <div class="flex items-center gap-3">
                  <input
                    type="radio"
                    name="payment"
                    value="cod"
                    checked
                    class="h-4 w-4 text-[#806344] focus:ring-[#806344]"
                  />
                  <div>
                    <p class="font-medium text-[#211f1d]">Cash on Delivery</p>
                    <p class="text-xs text-[#756a60] mt-0.5">Payment is collected in cash when your package is delivered to your doorstep.</p>
                  </div>
                </div>
                <span class="rounded bg-[#806344]/10 px-2.5 py-1 text-xs font-semibold text-[#806344]">Selected</span>
              </label>

              <!-- Online / Mobile Payment Options (Disabled & Coming Soon) -->
              <div class="flex items-center justify-between rounded-lg border border-[#ded6cc] bg-[#f0ede8] p-4 opacity-60 cursor-not-allowed">
                <div class="flex items-center gap-3">
                  <input type="radio" name="payment" disabled class="h-4 w-4" />
                  <div>
                    <p class="font-medium text-[#756a60]">Telebirr / Mobile Money</p>
                    <p class="text-xs text-[#a0958b] mt-0.5">Digital mobile payment integration.</p>
                  </div>
                </div>
                <span class="rounded bg-gray-200 px-2.5 py-1 text-xs font-medium text-gray-600">Coming soon</span>
              </div>

              <div class="flex items-center justify-between rounded-lg border border-[#ded6cc] bg-[#f0ede8] p-4 opacity-60 cursor-not-allowed">
                <div class="flex items-center gap-3">
                  <input type="radio" name="payment" disabled class="h-4 w-4" />
                  <div>
                    <p class="font-medium text-[#756a60]">Credit / Debit Card</p>
                    <p class="text-xs text-[#a0958b] mt-0.5">Visa, Mastercard, or local card.</p>
                  </div>
                </div>
                <span class="rounded bg-gray-200 px-2.5 py-1 text-xs font-medium text-gray-600">Coming soon</span>
              </div>
            </div>
          </section>

          <div class="mt-8 flex flex-wrap items-center gap-4">
            <UiAppButton
              variant="primary"
              size="default"
              :disabled="isSubmitting"
              class="min-w-[200px]"
              @click="handlePlaceOrder"
            >
              <span v-if="isSubmitting">Placing order...</span>
              <span v-else>Place order</span>
            </UiAppButton>

            <NuxtLink
              to="/cart"
              class="inline-flex h-12 items-center justify-center rounded-lg border border-[#b8ab9b] px-6 text-sm font-medium text-[#211f1d] transition hover:bg-[#eee8df]"
            >
              Return to cart
            </NuxtLink>
          </div>
        </div>

        <UiAppCard as="aside" class="h-fit">
          <h2 class="font-serif text-3xl text-[#211f1d]">
            Order summary
          </h2>

          <div class="mt-6 space-y-4">
            <article
              v-for="item in cartProducts"
              :key="item?.productId"
              class="flex gap-3"
            >
              <img
                v-if="item"
                :src="item.product.image"
                :alt="item.product.name"
                class="h-16 w-14 rounded-md object-cover object-top"
              />

              <div
                v-if="item"
                class="min-w-0 flex-1"
              >
                <p class="truncate text-sm font-medium text-[#211f1d]">
                  {{ item.product.name }}
                </p>

                <p class="mt-1 text-xs text-[#756a60]">
                  Qty {{ item.quantity }}
                </p>
              </div>

              <p
                v-if="item"
                class="text-sm font-semibold text-[#211f1d]"
              >
                {{ formatPrice(item.lineTotal) }}
              </p>
            </article>
          </div>

          <div class="mt-6 space-y-4 border-t border-[#ded6cc] pt-6 text-sm">
            <div class="flex justify-between gap-4 text-[#665c53]">
              <span>Subtotal</span>
              <span>{{ formatPrice(cartSubtotal) }}</span>
            </div>

            <div class="flex justify-between gap-4 text-[#665c53]">
              <span>Delivery</span>
              <span>{{ formatPrice(deliveryFee) }}</span>
            </div>

            <div class="flex justify-between gap-4 text-lg font-semibold text-[#211f1d]">
              <span>Total</span>
              <span>{{ formatPrice(orderTotal) }}</span>
            </div>
          </div>
        </UiAppCard>
      </div>

      <UiAppEmptyState
        v-else
        title="Nothing to checkout"
        description="Add some items to your cart before proceeding to checkout."
        action-label="Browse products"
        action-to="/products"
      />
    </section>
  </main>
</template>
