<script setup lang="ts">
defineProps<{
  active: 'profile' | 'orders' | 'shop' | 'payments' | 'seller-orders'
}>()

const { logout } = useAuth()
const { hasShop } = useSellerShop()

const allItems = [
  { key: 'profile', label: 'Profile', to: '/profile' },
  { key: 'orders', label: 'My Orders', to: '/orders' },
  { key: 'shop', label: 'My Shop', to: '/shop' },
  { key: 'seller-orders', label: 'Seller Orders', to: '/seller/orders', requiresShop: true },
  { key: 'payments', label: 'Payments', to: '/payments', requiresShop: true }
] as const

const items = computed(() =>
  allItems.filter(item => !('requiresShop' in item && item.requiresShop) || hasShop.value)
)
</script>

<template>
  <aside class="h-fit w-full lg:w-60">
    <h2 class="px-3 py-2 font-serif text-2xl text-[#211f1d]">
      My Account
    </h2>

    <nav class="mt-3 space-y-1">
      <NuxtLink
        v-for="item in items"
        :key="item.key"
        :to="item.to"
        class="block rounded-lg px-3 py-3 text-base text-[#665c53] transition hover:bg-[#eee8df] hover:text-[#211f1d]"
        :class="{
          'bg-[#eee8df] font-medium text-[#211f1d]': active === item.key
        }"
      >
        {{ item.label }}
      </NuxtLink>

      <button
        type="button"
        class="w-full rounded-lg px-3 py-3 text-left text-base text-red-600 transition hover:bg-red-50"
        @click="logout"
      >
        Logout
      </button>
    </nav>
  </aside>
</template>
