<script setup lang="ts">
const route = useRoute()

const orderId = computed(() =>
  typeof route.query.order === 'string' ? route.query.order : ''
)
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <UiAppCard
      as="section"
      padding="large"
      class="mx-auto max-w-2xl text-center"
    >
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
        <UiAppButton
          to="/orders"
          variant="primary"
        >
          View orders
        </UiAppButton>

        <UiAppButton
          to="/products"
          variant="secondary"
        >
          Keep shopping
        </UiAppButton>
      </div>
    </UiAppCard>
  </main>
</template>
