<script setup lang="ts">
import { ref } from 'vue'
import { useAdminDashboard } from '~/composables/useAdminDashboard'
import AdminSalesChart from '~/components/admin/AdminSalesChart.vue'
import AdminOrdersChart from '~/components/admin/AdminOrdersChart.vue'

const {
  users,
  products,
  orders,
  totalUsers,
  usersGrowth,
  totalShops,
  shopsGrowth,
  totalProducts,
  totalOrders,
  ordersGrowth,
  formattedTotalSales,
  salesGrowth,
  monthlyMetrics,
  pendingReportsCount,
  systemAlertsCount,
  moderationReports,
  systemServices,
  resolveReport,
  dismissReport
} = useAdminDashboard()

const activeTab = ref('dashboard')
const isMobileSidebarOpen = ref(false)

const selectTab = (tabId: string) => {
  activeTab.value = tabId
  isMobileSidebarOpen.value = false
}
</script>

<template>
  <div class="min-h-screen bg-[#f7f4ed] text-[#211f1d] font-sans antialiased flex flex-col lg:flex-row">
    <!-- Mobile Drawer Backdrop -->
    <div
      v-if="isMobileSidebarOpen"
      class="fixed inset-0 z-40 bg-black/25 backdrop-blur-xs lg:hidden"
      @click="isMobileSidebarOpen = false"
    />

    <!-- Left Sidebar -->
    <aside
      class="fixed inset-y-0 left-0 z-50 flex w-56 flex-col justify-between border-r border-[#ede8e1] bg-[#f7f4ed] p-5 transition-transform duration-200 ease-in-out lg:static lg:inset-auto lg:z-auto lg:translate-x-0 shrink-0"
      :class="[
        isMobileSidebarOpen ? 'translate-x-0 shadow-xl' : '-translate-x-full lg:translate-x-0'
      ]"
    >
      <div>
        <!-- Mobile close button -->
        <div class="mb-4 flex items-center justify-between lg:hidden">
          <span class="font-serif text-lg font-medium text-[#211f1d]">Admin</span>
          <button
            type="button"
            class="p-1 text-[#85776a] hover:text-[#211f1d]"
            @click="isMobileSidebarOpen = false"
          >
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Section: OVERVIEW -->
        <div>
          <p class="px-2.5 text-[10px] font-bold tracking-[0.14em] text-[#85776a] uppercase">
            OVERVIEW
          </p>
          <div class="mt-2 space-y-0.5">
            <!-- Dashboard (Active) -->
            <button
              type="button"
              class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'dashboard'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('dashboard')"
            >
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4">
                <rect x="1" y="1" width="4.5" height="4.5" rx="0.5" />
                <rect x="8.5" y="1" width="4.5" height="4.5" rx="0.5" />
                <rect x="1" y="8.5" width="4.5" height="4.5" rx="0.5" />
                <rect x="8.5" y="8.5" width="4.5" height="4.5" rx="0.5" />
              </svg>
              <span>Dashboard</span>
            </button>

            <!-- Analytics -->
            <button
              type="button"
              class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'analytics'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('analytics')"
            >
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round">
                <line x1="2" y1="12" x2="12" y2="2" />
              </svg>
              <span>Analytics</span>
            </button>
          </div>
        </div>

        <!-- Section: MARKETPLACE -->
        <div class="mt-6">
          <p class="px-2.5 text-[10px] font-bold tracking-[0.14em] text-[#85776a] uppercase">
            MARKETPLACE
          </p>
          <div class="mt-2 space-y-0.5">
            <!-- Users -->
            <button
              type="button"
              class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'users'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('users')"
            >
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4">
                <circle cx="7" cy="7" r="5" />
                <circle cx="7" cy="7" r="1.5" fill="currentColor" />
              </svg>
              <span>Users</span>
            </button>

            <!-- Shops -->
            <button
              type="button"
              class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'shops'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('shops')"
            >
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4">
                <rect x="2" y="3" width="10" height="8" rx="0.75" />
                <line x1="2" y1="6" x2="12" y2="6" />
              </svg>
              <span>Shops</span>
            </button>

            <!-- Products -->
            <button
              type="button"
              class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'products'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('products')"
            >
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4">
                <rect x="2" y="2" width="10" height="10" rx="0.75" />
                <line x1="5.5" y1="2" x2="5.5" y2="12" />
              </svg>
              <span>Products</span>
            </button>

            <!-- Categories -->
            <button
              type="button"
              class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'categories'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('categories')"
            >
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round">
                <line x1="2" y1="3.5" x2="12" y2="3.5" />
                <line x1="2" y1="7" x2="12" y2="7" />
                <line x1="2" y1="10.5" x2="12" y2="10.5" />
              </svg>
              <span>Categories</span>
            </button>

            <!-- Orders -->
            <button
              type="button"
              class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'orders'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('orders')"
            >
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4">
                <rect x="3.5" y="3.5" width="7" height="7" transform="rotate(45 7 7)" rx="0.5" />
              </svg>
              <span>Orders</span>
            </button>
          </div>
        </div>

        <!-- Section: MODERATION -->
        <div class="mt-6">
          <p class="px-2.5 text-[10px] font-bold tracking-[0.14em] text-[#85776a] uppercase">
            MODERATION
          </p>
          <div class="mt-2 space-y-0.5">
            <!-- Reports -->
            <button
              type="button"
              class="flex w-full items-center justify-between rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'reports'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('reports')"
            >
              <div class="flex items-center gap-2.5">
                <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="currentColor">
                  <path d="M2.5 1.5v11h1.2V7.5h6l1.8-3-1.8-3h-7.2z" />
                </svg>
                <span>Reports</span>
              </div>
              <span class="flex h-4 min-w-4 items-center justify-center rounded-full bg-[#fdece9] px-1.5 text-[10px] font-bold text-[#cc3828]">
                {{ pendingReportsCount }}
              </span>
            </button>
          </div>
        </div>

        <!-- Section: INFRASTRUCTURE -->
        <div class="mt-6">
          <p class="px-2.5 text-[10px] font-bold tracking-[0.14em] text-[#85776a] uppercase">
            INFRASTRUCTURE
          </p>
          <div class="mt-2 space-y-0.5">
            <!-- System -->
            <button
              type="button"
              class="flex w-full items-center justify-between rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="[
                activeTab === 'system'
                  ? 'bg-[#eae3d7] text-[#211f1d]'
                  : 'text-[#61574d] hover:bg-[#ede6da] hover:text-[#211f1d]'
              ]"
              @click="selectTab('system')"
            >
              <div class="flex items-center gap-2.5">
                <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4">
                  <circle cx="7" cy="7" r="5" />
                  <circle cx="7" cy="7" r="1.5" fill="currentColor" />
                </svg>
                <span>System</span>
              </div>
              <span class="flex h-4 min-w-4 items-center justify-center rounded-full bg-[#fdece9] px-1.5 text-[10px] font-bold text-[#cc3828]">
                {{ systemAlertsCount }}
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- Footer back link -->
      <div class="pt-6 border-t border-[#ede8e1]">
        <NuxtLink
          to="/"
          class="flex items-center gap-2 text-xs font-medium text-[#85776a] hover:text-[#211f1d] transition"
        >
          <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          <span>Storefront</span>
        </NuxtLink>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 min-w-0 flex flex-col">
      <!-- Breadcrumb Bar matching screenshot -->
      <div class="border-b border-[#ede8e1] px-6 py-3.5 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <!-- Mobile toggle button -->
          <button
            type="button"
            class="p-1 -ml-1 text-[#85776a] hover:text-[#211f1d] lg:hidden"
            @click="isMobileSidebarOpen = true"
          >
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>

          <div class="text-xs font-normal">
            <span class="text-[#85776a]">Admin</span>
            <span class="mx-1.5 text-[#beb3a6]">·</span>
            <span class="text-[#211f1d] font-medium capitalize">{{ activeTab }}</span>
          </div>
        </div>
      </div>

      <!-- Main Body Canvas -->
      <div class="p-6 sm:p-8 flex-1">
        <!-- View: DASHBOARD -->
        <div v-if="activeTab === 'dashboard'" class="space-y-6">
          <!-- Header -->
          <div>
            <h1 class="font-serif text-[28px] sm:text-[32px] font-bold text-[#211f1d] tracking-tight leading-none">
              Dashboard
            </h1>
            <p class="mt-2 text-xs sm:text-[13px] text-[#786e64]">
              Platform overview as of today
            </p>
          </div>

          <!-- 5 Stats Cards in a row -->
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-5 gap-4">
            <!-- TOTAL USERS -->
            <div class="rounded-lg border border-[#ede8e1] bg-white p-5 flex flex-col justify-between">
              <div>
                <p class="text-[10px] sm:text-[11px] font-bold tracking-[0.12em] text-[#85776a] uppercase">
                  TOTAL USERS
                </p>
                <p class="mt-2.5 font-serif text-3xl sm:text-[32px] font-normal leading-tight text-[#211f1d]">
                  {{ totalUsers }}
                </p>
                <p class="mt-0.5 text-xs text-[#786e64]">
                  registered accounts
                </p>
              </div>
              <p class="mt-3.5 text-xs font-normal text-[#2d7d46]">
                {{ usersGrowth }}
              </p>
            </div>

            <!-- TOTAL SHOPS -->
            <div class="rounded-lg border border-[#ede8e1] bg-white p-5 flex flex-col justify-between">
              <div>
                <p class="text-[10px] sm:text-[11px] font-bold tracking-[0.12em] text-[#85776a] uppercase">
                  TOTAL SHOPS
                </p>
                <p class="mt-2.5 font-serif text-3xl sm:text-[32px] font-normal leading-tight text-[#211f1d]">
                  {{ totalShops }}
                </p>
                <p class="mt-0.5 text-xs text-[#786e64]">
                  active storefronts
                </p>
              </div>
              <p class="mt-3.5 text-xs font-normal text-[#2d7d46]">
                {{ shopsGrowth }}
              </p>
            </div>

            <!-- TOTAL PRODUCTS (No trend, matches screenshot) -->
            <div class="rounded-lg border border-[#ede8e1] bg-white p-5 flex flex-col justify-between">
              <div>
                <p class="text-[10px] sm:text-[11px] font-bold tracking-[0.12em] text-[#85776a] uppercase">
                  TOTAL PRODUCTS
                </p>
                <p class="mt-2.5 font-serif text-3xl sm:text-[32px] font-normal leading-tight text-[#211f1d]">
                  {{ totalProducts }}
                </p>
                <p class="mt-0.5 text-xs text-[#786e64]">
                  listed items
                </p>
              </div>
              <!-- Empty spacer to match height -->
              <div class="h-4 mt-3.5" />
            </div>

            <!-- TOTAL ORDERS -->
            <div class="rounded-lg border border-[#ede8e1] bg-white p-5 flex flex-col justify-between">
              <div>
                <p class="text-[10px] sm:text-[11px] font-bold tracking-[0.12em] text-[#85776a] uppercase">
                  TOTAL ORDERS
                </p>
                <p class="mt-2.5 font-serif text-3xl sm:text-[32px] font-normal leading-tight text-[#211f1d]">
                  {{ totalOrders }}
                </p>
                <p class="mt-0.5 text-xs text-[#786e64]">
                  all time
                </p>
              </div>
              <p class="mt-3.5 text-xs font-normal text-[#2d7d46]">
                {{ ordersGrowth }}
              </p>
            </div>

            <!-- TOTAL SALES -->
            <div class="rounded-lg border border-[#ede8e1] bg-white p-5 flex flex-col justify-between">
              <div>
                <p class="text-[10px] sm:text-[11px] font-bold tracking-[0.12em] text-[#85776a] uppercase">
                  TOTAL SALES
                </p>
                <p class="mt-2.5 font-serif text-3xl sm:text-[32px] font-normal leading-tight text-[#211f1d]">
                  {{ formattedTotalSales }}
                </p>
                <p class="mt-0.5 text-xs text-[#786e64]">
                  gross revenue
                </p>
              </div>
              <p class="mt-3.5 text-xs font-normal text-[#2d7d46]">
                {{ salesGrowth }}
              </p>
            </div>
          </div>

          <!-- 2 Charts side by side -->
          <div class="grid grid-cols-1 xl:grid-cols-2 gap-5">
            <AdminSalesChart :metrics="monthlyMetrics" />
            <AdminOrdersChart :metrics="monthlyMetrics" />
          </div>
        </div>

        <!-- View: USERS -->
        <div v-else-if="activeTab === 'users'" class="space-y-6">
          <div>
            <h1 class="font-serif text-[28px] font-bold text-[#211f1d]">Users</h1>
            <p class="mt-1 text-xs text-[#786e64]">{{ totalUsers }} registered platform accounts</p>
          </div>

          <div class="rounded-lg border border-[#ede8e1] bg-white overflow-hidden">
            <table class="w-full text-left text-xs">
              <thead class="bg-[#faf7f2] border-b border-[#ede8e1] text-[#85776a] uppercase tracking-wider font-semibold">
                <tr>
                  <th class="px-5 py-3">Name</th>
                  <th class="px-5 py-3">Email</th>
                  <th class="px-5 py-3">Role</th>
                  <th class="px-5 py-3">Joined</th>
                  <th class="px-5 py-3 text-right">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3ede4]">
                <tr v-for="u in users" :key="u.email" class="hover:bg-[#faf7f2]/50">
                  <td class="px-5 py-3.5 font-medium text-[#211f1d]">{{ u.name || u.username }}</td>
                  <td class="px-5 py-3.5 text-[#786e64]">{{ u.email }}</td>
                  <td class="px-5 py-3.5">
                    <span class="rounded px-2 py-0.5 text-[11px] font-semibold uppercase tracking-wider" :class="u.role === 'seller' ? 'bg-amber-50 text-amber-800' : 'bg-stone-100 text-stone-700'">
                      {{ u.role }}
                    </span>
                  </td>
                  <td class="px-5 py-3.5 text-[#786e64]">{{ u.created_at || 'Recent' }}</td>
                  <td class="px-5 py-3.5 text-right text-emerald-700 font-medium">Active</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- View: ORDERS -->
        <div v-else-if="activeTab === 'orders'" class="space-y-6">
          <div>
            <h1 class="font-serif text-[28px] font-bold text-[#211f1d]">Orders</h1>
            <p class="mt-1 text-xs text-[#786e64]">{{ totalOrders }} total orders recorded</p>
          </div>

          <div class="rounded-lg border border-[#ede8e1] bg-white overflow-hidden">
            <table class="w-full text-left text-xs">
              <thead class="bg-[#faf7f2] border-b border-[#ede8e1] text-[#85776a] uppercase tracking-wider font-semibold">
                <tr>
                  <th class="px-5 py-3">ID</th>
                  <th class="px-5 py-3">Buyer</th>
                  <th class="px-5 py-3">Date</th>
                  <th class="px-5 py-3">Payment</th>
                  <th class="px-5 py-3">Status</th>
                  <th class="px-5 py-3 text-right">Total</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3ede4]">
                <tr v-for="o in orders" :key="o.id" class="hover:bg-[#faf7f2]/50">
                  <td class="px-5 py-3.5 font-medium text-[#211f1d]">#{{ o.id }}</td>
                  <td class="px-5 py-3.5 text-[#786e64]">{{ o.buyerName }}</td>
                  <td class="px-5 py-3.5 text-[#786e64]">{{ o.date }}</td>
                  <td class="px-5 py-3.5 text-[#786e64]">{{ o.paymentMethod }}</td>
                  <td class="px-5 py-3.5">
                    <span class="rounded px-2 py-0.5 text-[11px] font-semibold" :class="o.status === 'Delivered' ? 'bg-emerald-50 text-emerald-700' : o.status === 'Shipped' ? 'bg-blue-50 text-blue-700' : 'bg-amber-50 text-amber-700'">
                      {{ o.status }}
                    </span>
                  </td>
                  <td class="px-5 py-3.5 text-right font-serif text-sm font-medium text-[#211f1d]">${{ o.total }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- View: SHOPS -->
        <div v-else-if="activeTab === 'shops'" class="space-y-6">
          <div>
            <h1 class="font-serif text-[28px] font-bold text-[#211f1d]">Shops</h1>
            <p class="mt-1 text-xs text-[#786e64]">{{ totalShops }} active storefronts</p>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <div v-for="shopName in ['Atelier North', 'Beyond Score', 'True Form', 'Minimal Studio', 'Urban Thread', 'Mara Studio']" :key="shopName" class="rounded-lg border border-[#ede8e1] bg-white p-5">
              <div class="flex items-center justify-between">
                <span class="font-serif text-lg font-medium text-[#211f1d]">{{ shopName }}</span>
                <span class="rounded bg-emerald-50 px-2 py-0.5 text-[10px] font-bold text-emerald-700 uppercase">Active</span>
              </div>
              <p class="mt-2 text-xs text-[#786e64]">Curated apparel and lifestyle products</p>
            </div>
          </div>
        </div>

        <!-- View: PRODUCTS -->
        <div v-else-if="activeTab === 'products'" class="space-y-6">
          <div>
            <h1 class="font-serif text-[28px] font-bold text-[#211f1d]">Products</h1>
            <p class="mt-1 text-xs text-[#786e64]">{{ totalProducts }} listed marketplace items</p>
          </div>

          <div class="rounded-lg border border-[#ede8e1] bg-white overflow-hidden">
            <table class="w-full text-left text-xs">
              <thead class="bg-[#faf7f2] border-b border-[#ede8e1] text-[#85776a] uppercase tracking-wider font-semibold">
                <tr>
                  <th class="px-5 py-3">Product</th>
                  <th class="px-5 py-3">Shop</th>
                  <th class="px-5 py-3">Category</th>
                  <th class="px-5 py-3">Stock</th>
                  <th class="px-5 py-3 text-right">Price</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3ede4]">
                <tr v-for="p in products" :key="p.id" class="hover:bg-[#faf7f2]/50">
                  <td class="px-5 py-3.5 font-medium text-[#211f1d]">{{ p.name }}</td>
                  <td class="px-5 py-3.5 text-[#786e64]">{{ p.shop }}</td>
                  <td class="px-5 py-3.5 text-[#786e64]">{{ p.category }}</td>
                  <td class="px-5 py-3.5 text-[#786e64]">{{ p.stock }} units</td>
                  <td class="px-5 py-3.5 text-right font-serif text-sm font-medium text-[#211f1d]">{{ p.price }} ETB</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- View: CATEGORIES -->
        <div v-else-if="activeTab === 'categories'" class="space-y-6">
          <div>
            <h1 class="font-serif text-[28px] font-bold text-[#211f1d]">Categories</h1>
            <p class="mt-1 text-xs text-[#786e64]">Catalog distribution by category</p>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <div v-for="cat in ['Men', 'Women', 'Kids', 'Accessories', 'Shoes']" :key="cat" class="rounded-lg border border-[#ede8e1] bg-white p-5">
              <span class="font-serif text-lg font-medium text-[#211f1d]">{{ cat }}</span>
              <p class="mt-1 text-xs text-[#786e64]">Active platform collection</p>
            </div>
          </div>
        </div>

        <!-- View: REPORTS -->
        <div v-else-if="activeTab === 'reports'" class="space-y-6">
          <div>
            <h1 class="font-serif text-[28px] font-bold text-[#211f1d]">Reports</h1>
            <p class="mt-1 text-xs text-[#786e64]">{{ pendingReportsCount }} moderation items requiring attention</p>
          </div>

          <div class="space-y-3">
            <div v-for="r in moderationReports" :key="r.id" class="rounded-lg border border-[#ede8e1] bg-white p-5">
              <div class="flex items-center justify-between">
                <span class="font-medium text-xs text-[#211f1d]">{{ r.targetTitle }}</span>
                <span class="text-[10px] font-bold uppercase rounded px-2 py-0.5 text-red-700 bg-red-50">{{ r.severity }}</span>
              </div>
              <p class="mt-2 text-xs text-[#786e64]">{{ r.reason }}</p>
              <div v-if="r.status === 'pending'" class="mt-3 flex gap-2">
                <button type="button" class="rounded bg-[#2d7d46] px-3 py-1 text-xs text-white font-medium hover:bg-[#236838]" @click="resolveReport(r.id)">Resolve</button>
                <button type="button" class="rounded border border-[#ede8e1] px-3 py-1 text-xs text-[#786e64] hover:bg-[#faf7f2]" @click="dismissReport(r.id)">Dismiss</button>
              </div>
              <div v-else class="mt-3 text-xs text-[#2d7d46] font-medium capitalize">Status: {{ r.status }}</div>
            </div>
          </div>
        </div>

        <!-- View: SYSTEM -->
        <div v-else-if="activeTab === 'system'" class="space-y-6">
          <div>
            <h1 class="font-serif text-[28px] font-bold text-[#211f1d]">System</h1>
            <p class="mt-1 text-xs text-[#786e64]">Core infrastructure telemetry and alerts</p>
          </div>

          <div class="space-y-3">
            <div v-for="s in systemServices" :key="s.name" class="rounded-lg border border-[#ede8e1] bg-white p-5">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="h-2 w-2 rounded-full" :class="s.status === 'healthy' ? 'bg-emerald-500' : 'bg-amber-500'" />
                  <span class="text-xs font-medium text-[#211f1d]">{{ s.name }}</span>
                </div>
                <span class="text-xs text-[#786e64]">Latency: {{ s.latency }}</span>
              </div>
              <p v-if="s.message" class="mt-2 text-xs text-amber-800 bg-amber-50 p-2 rounded">{{ s.message }}</p>
            </div>
          </div>
        </div>

        <!-- View: ANALYTICS -->
        <div v-else-if="activeTab === 'analytics'" class="space-y-6">
          <div>
            <h1 class="font-serif text-[28px] font-bold text-[#211f1d]">Analytics</h1>
            <p class="mt-1 text-xs text-[#786e64]">6-month sales and orders velocity</p>
          </div>

          <div class="grid grid-cols-1 xl:grid-cols-2 gap-5">
            <AdminSalesChart :metrics="monthlyMetrics" />
            <AdminOrdersChart :metrics="monthlyMetrics" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
