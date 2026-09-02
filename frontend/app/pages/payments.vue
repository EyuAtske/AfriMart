<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'

definePageMeta({
  middleware: 'auth'
})

const { shop, hasShop, addPaymentMethod } = useSellerShop()

const paymentForm = reactive({
  type: 'Telebirr' as 'Telebirr' | 'CBE',
  accountName: '',
  accountNumber: ''
})

const paymentError = ref('')

const buttonClass = 'h-12 rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50'

const submitPaymentMethod = () => {
  paymentError.value = ''

  if (!paymentForm.accountName.trim() || !paymentForm.accountNumber.trim()) {
    paymentError.value = 'Please add the account name and account number.'
    return
  }

  addPaymentMethod(paymentForm)

  paymentForm.type = 'Telebirr'
  paymentForm.accountName = ''
  paymentForm.accountNumber = ''
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-8">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 lg:flex-row">
      <AccountSidebar active="payments" />

      <section class="min-w-0 flex-1">
        <div class="mb-8">
          <h1 class="font-serif text-4xl tracking-[-0.025em] text-[#211f1d]">
            Payments
          </h1>

          <p class="mt-2 text-base text-[#756a60]">
            Add the ways your shop accepts money from buyers.
          </p>
        </div>

        <div
          v-if="!hasShop"
          class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8"
        >
          <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
            Seller account needed
          </p>

          <h2 class="mt-3 font-serif text-3xl text-[#211f1d]">
            Create a shop first.
          </h2>

          <p class="mt-3 max-w-xl text-base leading-7 text-[#756a60]">
            Payment methods are only available after you become a seller and create your shop.
          </p>

          <NuxtLink
            to="/shop"
            :class="`${buttonClass} mt-8 inline-flex items-center justify-center`"
          >
            Go to my shop
          </NuxtLink>
        </div>

        <div
          v-else-if="shop"
          class="space-y-6"
        >
          <section class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
            <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
              {{ shop.name }}
            </p>

            <h2 class="mt-2 font-serif text-3xl text-[#211f1d]">
              Add payment method
            </h2>

            <form
              class="mt-6 space-y-5"
              @submit.prevent="submitPaymentMethod"
            >
              <div class="grid grid-cols-2 gap-3 rounded-full border border-[#d9d0c4] bg-[#f5f1e9] p-1">
                <button
                  type="button"
                  class="h-11 rounded-full text-sm font-medium transition"
                  :class="paymentForm.type === 'Telebirr' ? 'bg-[#211f1d] text-white' : 'text-[#665c53] hover:bg-[#eee8df]'"
                  @click="paymentForm.type = 'Telebirr'"
                >
                  Telebirr
                </button>

                <button
                  type="button"
                  class="h-11 rounded-full text-sm font-medium transition"
                  :class="paymentForm.type === 'CBE' ? 'bg-[#211f1d] text-white' : 'text-[#665c53] hover:bg-[#eee8df]'"
                  @click="paymentForm.type = 'CBE'"
                >
                  CBE
                </button>
              </div>

              <div class="grid gap-5 sm:grid-cols-2">
                <AuthInput
                  v-model="paymentForm.accountName"
                  label="Account name"
                  placeholder="Name on the account"
                  name="payment-account-name"
                />

                <AuthInput
                  v-model="paymentForm.accountNumber"
                  label="Account number"
                  placeholder="Phone or bank account number"
                  name="payment-account-number"
                />
              </div>

              <p
                v-if="paymentError"
                class="text-sm font-medium text-red-600"
              >
                {{ paymentError }}
              </p>

              <button
                type="submit"
                :class="buttonClass"
              >
                Save method
              </button>
            </form>
          </section>

          <section class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
            <h2 class="text-xl font-medium text-[#211f1d]">
              Saved methods
            </h2>

            <div
              v-if="shop.paymentMethods.length"
              class="mt-5 grid gap-4 sm:grid-cols-2"
            >
              <article
                v-for="method in shop.paymentMethods"
                :key="method.id"
                class="rounded-[10px] border border-[#ded6cc] bg-[#f5f1e9] p-5"
              >
                <p class="text-xs font-medium uppercase tracking-[0.16em] text-[#806344]">
                  {{ method.type }}
                </p>

                <h3 class="mt-2 text-lg font-medium text-[#211f1d]">
                  {{ method.accountName }}
                </h3>

                <p class="mt-1 text-sm text-[#756a60]">
                  {{ method.accountNumber }}
                </p>
              </article>
            </div>

            <p
              v-else
              class="mt-4 text-sm leading-6 text-[#756a60]"
            >
              No payment methods added yet.
            </p>
          </section>
        </div>
      </section>
    </div>
  </main>
</template>


