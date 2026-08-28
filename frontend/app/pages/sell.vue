<script setup lang="ts">
const { hasShop } = useSellerShop()
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <section class="mx-auto grid max-w-7xl gap-8 lg:grid-cols-[1fr_0.9fr] lg:items-center">
      <div>
        <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
          Seller marketplace
        </p>

        <h1 class="mt-3 font-serif text-5xl leading-tight tracking-[-0.025em] text-[#211f1d] sm:text-6xl">
          Open your shop and start listing today.
        </h1>

        <p class="mt-6 max-w-2xl text-base leading-8 text-[#756a60]">
          Afrimart is mocked on the frontend for now, so you can test the shop creation, product listing, payments, and order dashboards before the backend services arrive.
        </p>

        <div class="mt-8 flex flex-col gap-3 sm:flex-row">
          <NuxtLink
            to="/shop"
            class="inline-flex h-12 items-center justify-center rounded-full bg-[#211f1d] px-6 text-sm font-medium uppercase tracking-[0.14em] text-white transition hover:bg-[#3b3733]"
          >
            {{ hasShop ? 'Go to dashboard' : 'Become a seller' }}
          </NuxtLink>

          <NuxtLink
            v-if="hasShop"
            to="/seller/orders"
            class="inline-flex h-12 items-center justify-center rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
          >
            Seller orders
          </NuxtLink>
        </div>
      </div>

      <div class="overflow-hidden rounded-[8px] border border-[#d9d0c4] bg-[#faf8f4] shadow-[0_20px_70px_rgba(33,31,29,0.06)]">
        <img
          src="/images/product8.jpg"
          alt="Seller product display"
          class="h-[520px] w-full object-cover object-top"
        />
      </div>
    </section>
  </main>
</template>
