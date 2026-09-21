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
      <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-amber-100 text-3xl text-amber-800">
        ℹ️
      </div>

      <p class="mt-6 text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
        Payment Service Status
      </p>

      <h1 class="mt-2 font-serif text-4xl text-[#211f1d]">
        Payment Processing Unavailable
      </h1>

      <p class="mt-4 text-base leading-7 text-[#756a60]">
        Online payment endpoints are not deployed on the backend server.
        No live transactions or order simulations have been processed.
      </p>

      <div class="mt-8 flex flex-col gap-3 sm:flex-row sm:justify-center">
        <UiAppButton
          to="/cart"
          variant="primary"
        >
          View cart
        </UiAppButton>

        <UiAppButton
          to="/products"
          variant="secondary"
        >
          Browse catalog
        </UiAppButton>
      </div>
    </UiAppCard>
  </main>
</template>
