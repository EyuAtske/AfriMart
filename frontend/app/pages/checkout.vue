<script setup lang="ts">
const router = useRouter()
const { user, isLoggedIn } = useAuth()
const { cartProducts, cartSubtotal, createOrder } = useMarketplace()
const { mockRequest } = useApiClient()
const { gtag } = useGtag()
const { track, isFeatureEnabled } = useAnalytics()
const newCheckoutEnabled = computed(() =>
  isFeatureEnabled('new-checkout')
)

const checkoutForm = reactive({
  fullName: '',
  phone: '',
  address: '',
  city: 'Addis Ababa',
  paymentMethod: 'Cash on delivery'
})

const errorMessage = ref('')
const isSubmitting = ref(false)

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

const submitCheckout = async () => {
  errorMessage.value = ''

  if (!cartProducts.value.length) {
    errorMessage.value = 'Your cart is empty.'
    return
  }

  if (
    !checkoutForm.fullName.trim() ||
    !checkoutForm.phone.trim() ||
    !checkoutForm.address.trim() ||
    !checkoutForm.city.trim()
  ) {
    console.log('[DEBUG] submitCheckout aborted: incomplete details', { ...checkoutForm })
    errorMessage.value = 'Please complete the delivery details.'
    return
  }

  const rawPhone = checkoutForm.phone.trim()
  const cleanPhone = rawPhone.replace(/[\s\-()]/g, '')
  const ethiopianPhoneRegex = /^\+251[79]\d{8}$/

  if (!ethiopianPhoneRegex.test(cleanPhone)) {
    errorMessage.value = 'Please enter a valid Ethiopian phone number: +251 followed by 7 or 9 and 8 digits (e.g. +251911234567 or +251711234567).'
    return
  }

  isSubmitting.value = true

  try {
    await mockRequest(
      {
        cart: cartProducts.value,
        delivery: checkoutForm
      },
      { delay: 0 }
    )

    const order = await createOrder({
      buyerName: checkoutForm.fullName.trim(),
      deliveryAddress: `${checkoutForm.address.trim()}, ${checkoutForm.city.trim()}`,
      phone: cleanPhone,
      paymentMethod: checkoutForm.paymentMethod,
      deliveryFee: deliveryFee.value
    })
if (!order) {
  errorMessage.value = 'Could not create the mock order.'
  return
}

gtag('event', 'purchase', {
  transaction_id: String(order.id),
  value: orderTotal.value,
  currency: 'ETB',
  items: cartProducts.value.map(item => ({
    item_id: String(item.product.id),
    item_name: item.product.name,
    price: item.product.price,
    quantity: item.quantity
  }))
})

isLoggedIn.value = true
router.push(`/payment/success?order=${order.id}`)
  } catch (error) {
    router.push('/payment/failure')
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

        <p
          v-if="!isLoggedIn"
          class="mt-3 max-w-2xl text-sm leading-6 text-[#756a60]"
        >
          You can complete this mock checkout as a guest. When backend auth is connected, this page can become protected.
        </p>
      </div>

      <div
        v-if="cartProducts.length"
        class="grid gap-8 lg:grid-cols-[1fr_380px]"
      >
        <form
          action="javascript:void(0)"
          class="rounded-xl border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8"
          @submit.prevent="submitCheckout"
        >
          <h2 class="font-serif text-3xl text-[#211f1d]">
            Delivery details
          </h2>

          <div class="mt-6 grid gap-5 sm:grid-cols-2">
            <AuthInput
              v-model="checkoutForm.fullName"
              label="Full name"
              placeholder="Name for delivery"
              name="checkout-name"
            />

            <AuthInput
              v-model="checkoutForm.phone"
              label="Phone"
              type="tel"
              placeholder="+251911234567"
              name="checkout-phone"
            />
          </div>

          <div class="mt-5 grid gap-5 sm:grid-cols-[1fr_220px]">
            <AuthInput
              v-model="checkoutForm.address"
              label="Address"
              placeholder="Street, building, area"
              name="checkout-address"
            />

            <AuthInput
              v-model="checkoutForm.city"
              label="City"
              placeholder="City"
              name="checkout-city"
            />
          </div>

          <section class="mt-8">
            <h2 class="font-serif text-3xl text-[#211f1d]">
              Payment method
            </h2>

            <p
              v-if="newCheckoutEnabled"
              class="mt-2 text-xs font-medium uppercase tracking-[0.16em] text-[#806344]"
            >
              New checkout experience
            </p>
            
            <div class="mt-5 grid gap-3 sm:grid-cols-3">
              <label
                v-for="method in ['Cash on delivery', 'Telebirr', 'CBE']"
                :key="method"
                class="cursor-pointer rounded-lg border p-4 transition"
                :class="
                  checkoutForm.paymentMethod === method
                    ? 'border-[#806344] bg-[#eee8df]'
                    : 'border-[#d9d0c4] bg-[#f5f1e9] hover:border-[#b9aa98]'
                "
              >
                <input
                  v-model="checkoutForm.paymentMethod"
                  type="radio"
                  name="payment-method"
                  :value="method"
                  class="sr-only"
                  @change="track('payment_method_selected', {
    payment_method: method
  })"
/>

                <span class="block text-sm font-medium text-[#211f1d]">
                  {{ method }}
                </span>

                <span class="mt-2 block text-xs leading-5 text-[#756a60]">
                  {{
                    method === 'Cash on delivery'
                      ? 'Pay when the order arrives.'
                      : 'Mock online payment approval.'
                  }}
                </span>
              </label>
            </div>
          </section>

          <UiAppAlert
            v-if="errorMessage"
            class="mt-5"
          >
            {{ errorMessage }}
          </UiAppAlert>

          <div class="mt-8">
            <UiAppButton
              type="submit"
              variant="secondary"
              :disabled="isSubmitting"
              class="w-full sm:w-auto"
            >
              {{ isSubmitting ? 'Creating order...' : 'Place order' }}
            </UiAppButton>
          </div>
        </form>

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
