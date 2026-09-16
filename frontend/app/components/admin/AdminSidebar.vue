<script setup lang="ts">
import { computed } from 'vue'
import { useAdminDashboard } from '~/composables/useAdminDashboard'

const props = defineProps<{
  activeTab: string
  isOpen?: boolean
}>()

const emit = defineEmits<{
  (e: 'select', tabId: string): void
  (e: 'close'): void
}>()

const { pendingReportsCount, systemAlertsCount } = useAdminDashboard()

interface NavItem {
  id: string
  label: string
  icon: string
  badge?: number
}

interface NavSection {
  title: string
  items: NavItem[]
}

const navigationSections = computed<NavSection[]>(() => [
  {
    title: 'OVERVIEW',
    items: [
      { id: 'dashboard', label: 'Dashboard', icon: 'grid' },
      { id: 'analytics', label: 'Analytics', icon: 'analytics' }
    ]
  },
  {
    title: 'MARKETPLACE',
    items: [
      { id: 'users', label: 'Users', icon: 'users' },
      { id: 'shops', label: 'Shops', icon: 'shops' },
      { id: 'products', label: 'Products', icon: 'products' },
      { id: 'categories', label: 'Categories', icon: 'categories' },
      { id: 'orders', label: 'Orders', icon: 'orders' }
    ]
  },
  {
    title: 'MODERATION',
    items: [
      {
        id: 'reports',
        label: 'Reports',
        icon: 'reports',
        badge: pendingReportsCount.value
      }
    ]
  },
  {
    title: 'INFRASTRUCTURE',
    items: [
      {
        id: 'system',
        label: 'System',
        icon: 'system',
        badge: systemAlertsCount.value
      }
    ]
  }
])

const handleSelect = (id: string) => {
  emit('select', id)
  emit('close')
}
</script>

<template>
  <div>
    <!-- Mobile Backdrop -->
    <div
      v-if="isOpen"
      class="fixed inset-0 z-40 bg-black/30 backdrop-blur-xs transition-opacity lg:hidden"
      @click="emit('close')"
    />

    <!-- Sidebar Container -->
    <aside
      class="fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r border-[#ebe5dd] bg-[#faf7f2] px-4 py-6 transition-transform duration-300 ease-in-out lg:static lg:inset-auto lg:z-auto lg:w-60 lg:translate-x-0 lg:border-r lg:border-[#ebe5dd] lg:bg-transparent lg:px-0 lg:py-2"
      :class="[
        isOpen ? 'translate-x-0 shadow-2xl' : '-translate-x-full lg:translate-x-0'
      ]"
    >
      <!-- Mobile Header with Close Button -->
      <div class="mb-4 flex items-center justify-between px-3 lg:hidden">
        <span class="font-serif text-xl font-medium text-[#211f1d]">Admin Portal</span>
        <button
          type="button"
          aria-label="Close sidebar"
          class="rounded-md p-1.5 text-[#6e6358] hover:bg-[#ede7de]"
          @click="emit('close')"
        >
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Navigation Sections -->
      <nav class="flex-1 space-y-6 overflow-y-auto px-1">
        <div v-for="section in navigationSections" :key="section.title">
          <p class="px-3 pb-2 text-[11px] font-semibold tracking-[0.14em] text-[#8a7a6c] uppercase">
            {{ section.title }}
          </p>

          <div class="space-y-1">
            <button
              v-for="item in section.items"
              :key="item.id"
              type="button"
              class="group flex w-full items-center justify-between rounded-lg px-3 py-2 text-sm font-medium transition-all duration-150"
              :class="[
                activeTab === item.id
                  ? 'bg-[#eae3d8] text-[#211f1d] shadow-2xs font-semibold'
                  : 'text-[#61574d] hover:bg-[#efe9e0] hover:text-[#211f1d]'
              ]"
              @click="handleSelect(item.id)"
            >
              <div class="flex items-center gap-3">
                <!-- Grid / Dashboard Icon -->
                <svg
                  v-if="item.icon === 'grid'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.6"
                >
                  <rect x="2" y="2" width="5" height="5" rx="1" />
                  <rect x="9" y="2" width="5" height="5" rx="1" />
                  <rect x="2" y="9" width="5" height="5" rx="1" />
                  <rect x="9" y="9" width="5" height="5" rx="1" />
                </svg>

                <!-- Analytics Icon -->
                <svg
                  v-else-if="item.icon === 'analytics'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.6"
                  stroke-linecap="round"
                >
                  <path d="M2.5 13.5L13.5 2.5" />
                  <path d="M8.5 2.5H13.5V7.5" />
                </svg>

                <!-- Users Icon -->
                <svg
                  v-else-if="item.icon === 'users'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.6"
                >
                  <circle cx="8" cy="5" r="3" />
                  <path d="M2.5 14C2.5 11.5 5 9.5 8 9.5C11 9.5 13.5 11.5 13.5 14" stroke-linecap="round" />
                </svg>

                <!-- Shops Icon -->
                <svg
                  v-else-if="item.icon === 'shops'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.6"
                  stroke-linecap="round"
                >
                  <path d="M2 5L3 2H13L14 5" />
                  <path d="M2.5 5C2.5 6 3.5 7 4.5 7C5.5 7 6.5 6 6.5 5C6.5 6 7.5 7 8.5 7C9.5 7 10.5 6 10.5 5C10.5 6 11.5 7 12.5 7C13.5 7 14.5 6 14.5 5" />
                  <path d="M3 7V13.5H13V7" />
                </svg>

                <!-- Products Icon -->
                <svg
                  v-else-if="item.icon === 'products'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.6"
                >
                  <rect x="2.5" y="2.5" width="11" height="11" rx="1.5" />
                  <path d="M2.5 6.5H13.5" />
                  <path d="M8 2.5V6.5" />
                </svg>

                <!-- Categories Icon -->
                <svg
                  v-else-if="item.icon === 'categories'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.6"
                  stroke-linecap="round"
                >
                  <line x1="2.5" y1="4" x2="13.5" y2="4" />
                  <line x1="2.5" y1="8" x2="13.5" y2="8" />
                  <line x1="2.5" y1="12" x2="9.5" y2="12" />
                </svg>

                <!-- Orders Icon -->
                <svg
                  v-else-if="item.icon === 'orders'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.6"
                >
                  <rect x="3.5" y="3.5" width="9" height="9" transform="rotate(45 8 8)" rx="1" />
                </svg>

                <!-- Reports Icon -->
                <svg
                  v-else-if="item.icon === 'reports'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="currentColor"
                >
                  <path d="M3 2V14H4.5V9H11.5L13.5 5.5L11.5 2H3Z" />
                </svg>

                <!-- System Icon -->
                <svg
                  v-else-if="item.icon === 'system'"
                  class="h-4 w-4 shrink-0 transition-colors"
                  :class="activeTab === item.id ? 'text-[#211f1d]' : 'text-[#8a7a6c] group-hover:text-[#211f1d]'"
                  viewBox="0 0 16 16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.6"
                >
                  <circle cx="8" cy="8" r="5" />
                  <circle cx="8" cy="8" r="1.5" fill="currentColor" />
                </svg>

                <span>{{ item.label }}</span>
              </div>

              <!-- Badge (e.g. Reports '3', System '1') -->
              <span
                v-if="item.badge && item.badge > 0"
                class="flex h-5 min-w-5 items-center justify-center rounded-full bg-[#faeae7] px-1.5 text-[11px] font-semibold text-[#c83928]"
              >
                {{ item.badge }}
              </span>
            </button>
          </div>
        </div>
      </nav>

      <!-- Bottom link to marketplace -->
      <div class="mt-6 border-t border-[#ebe5dd] pt-4">
        <NuxtLink
          to="/"
          class="flex items-center gap-2 rounded-lg px-3 py-2 text-xs font-medium text-[#7a6f63] transition hover:bg-[#efe9e0] hover:text-[#211f1d]"
        >
          <svg class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          <span>Return to Storefront</span>
        </NuxtLink>
      </div>
    </aside>
  </div>
</template>
