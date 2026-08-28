<script setup lang="ts">
const route = useRoute()

const orderId = computed(() =>
  typeof route.query.order === 'string' ? route.query.order : ''
)
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <section class="mx-auto max-w-2xl rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-8 text-center shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-10">
      <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-[#e6eee5] text-3xl text-[#536653]">
        OK
      </div>

      <p class="mt-6 text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
        Payment status
      </p>

      <h1 class="mt-2 font-serif text-4xl text-[#211f1d]">
        Order confirmed
      </h1>

      <p class="mt-4 text-base leading-7 text-[#756a60]">
        Your mock order has been created successfully{{ orderId ? ` as #${orderId}` : '' }}.
      </p>

      <div class="mt-8 flex flex-col gap-3 sm:flex-row sm:justify-center">
        <NuxtLink
          to="/orders"
          class="inline-flex h-12 items-center justify-center rounded-full bg-[#211f1d] px-6 text-sm font-medium uppercase tracking-[0.14em] text-white transition hover:bg-[#3b3733]"
        >
          View orders
        </NuxtLink>

        <NuxtLink
          to="/products"
          class="inline-flex h-12 items-center justify-center rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
        >
          Keep shopping
        </NuxtLink>
      </div>
    </section>
  </main>
</template>
