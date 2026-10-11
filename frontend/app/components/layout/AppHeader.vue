<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const isMenuOpen = ref(false)
const isSearchOpen = ref(false)
const searchQuery = ref('')
const isMounted = ref(false)

// Expose alias if referenced
const isOpen = isMenuOpen

const { cartProducts } = useMarketplace()
const { isCartBouncing } = useFlyToCart()

const closeOverlays = () => {
  isMenuOpen.value = false
  isSearchOpen.value = false
}

// Ensure overlays are reset and closed on any route transition
watch(
  () => route.fullPath,
  () => {
    closeOverlays()
  }
)

// Body scroll lock management when drawer is open
watch(isMenuOpen, (open) => {
  if (import.meta.client) {
    document.body.style.overflow = open ? 'hidden' : ''
  }
})

let mediaQuery: MediaQueryList | null = null
let handleMediaChange: ((e: MediaQueryListEvent | MediaQueryList) => void) | null = null

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    closeOverlays()
  }
}

onMounted(() => {
  isMounted.value = true
  closeOverlays()

  if (import.meta.client) {
    window.addEventListener('keydown', handleKeydown)
    mediaQuery = window.matchMedia('(min-width: 1024px)')
    handleMediaChange = (e: MediaQueryListEvent | MediaQueryList) => {
      if (e.matches) {
        closeOverlays()
      }
    }
    mediaQuery.addEventListener?.('change', handleMediaChange)
  }
})

onUnmounted(() => {
  if (import.meta.client) {
    window.removeEventListener('keydown', handleKeydown)
    if (mediaQuery && handleMediaChange) {
      mediaQuery.removeEventListener?.('change', handleMediaChange)
    }
    document.body.style.overflow = ''
  }
})

const cartCount = computed(() =>
  cartProducts.value.reduce((total, item) => total + (item?.quantity || 0), 0)
)

const isAiQuery = (text: string) => {
  const query = text.toLowerCase()
  const aiKeywords = [
    'i want', 'i need', 'find me', 'looking for', 'show me',
    'recommend', 'suggest', 'outfit', 'under', 'cheap', 'what', 'style'
  ]
  return aiKeywords.some(kw => query.includes(kw))
}

const handleAiSearch = async () => {
  const search = searchQuery.value.trim()
  closeOverlays()
  if (search) {
    await router.push({
      path: '/ai',
      query: { prompt: search }
    })
  } else {
    await router.push('/ai')
  }
}

const submitSearch = async () => {
  const search = searchQuery.value.trim()
  if (!search) return

  closeOverlays()
  if (isAiQuery(search)) {
    await router.push({
      path: '/ai',
      query: { prompt: search }
    })
  } else {
    await router.push({
      path: '/products',
      query: { search }
    })
  }
}

const openCategories = ref<Record<string, boolean>>({
  MEN: false,
  WOMEN: false,
  KIDS: false
})

const toggleCategory = (name: string) => {
  openCategories.value[name] = !openCategories.value[name]
}

const categories = [
  {
    name: 'MEN',
    to: '/products?category=Men',
    sections: [
      {
        title: 'CLOTHING',
        items: [
          { name: 'T-Shirts', to: '/products?category=Men&search=t-shirt' },
          { name: 'Shirts', to: '/products?category=Men&search=shirt' },
          { name: 'Trousers', to: '/products?category=Men&search=trousers' },
          { name: 'Jeans', to: '/products?category=Men&search=jeans' },
          { name: 'Jackets', to: '/products?category=Men&search=jacket' },
          { name: 'Hoodies', to: '/products?category=Men&search=hoodie' }
        ]
      },
      {
        title: 'SHOES',
        items: [
          { name: 'Sneakers', to: '/products?category=Shoes&search=sneakers' },
          { name: 'Formal Shoes', to: '/products?category=Shoes&search=formal' },
          { name: 'Boots', to: '/products?category=Shoes&search=boots' },
          { name: 'Sandals', to: '/products?category=Shoes&search=sandals' }
        ]
      },
      {
        title: 'ACCESSORIES',
        items: [
          { name: 'Bags', to: '/products?category=Accessories&search=bag' },
          { name: 'Watches', to: '/products?category=Accessories&search=watch' },
          { name: 'Belts', to: '/products?category=Accessories&search=belt' },
          { name: 'Hats', to: '/products?category=Accessories&search=hat' }
        ]
      }
    ]
  },
  {
    name: 'WOMEN',
    to: '/products?category=Women',
    sections: [
      {
        title: 'CLOTHING',
        items: [
          { name: 'Dresses', to: '/products?category=Women&search=dress' },
          { name: 'Tops', to: '/products?category=Women&search=top' },
          { name: 'Trousers', to: '/products?category=Women&search=trousers' },
          { name: 'Jeans', to: '/products?category=Women&search=jeans' },
          { name: 'Skirts', to: '/products?category=Women&search=skirt' },
          { name: 'Jackets', to: '/products?category=Women&search=jacket' },
          { name: 'Sweaters', to: '/products?category=Women&search=sweater' }
        ]
      },
      {
        title: 'SHOES',
        items: [
          { name: 'Sneakers', to: '/products?category=Shoes&search=sneakers' },
          { name: 'Heels', to: '/products?category=Shoes&search=heels' },
          { name: 'Boots', to: '/products?category=Shoes&search=boots' },
          { name: 'Sandals', to: '/products?category=Shoes&search=sandals' },
          { name: 'Flats', to: '/products?category=Shoes&search=flats' }
        ]
      },
      {
        title: 'ACCESSORIES',
        items: [
          { name: 'Bags', to: '/products?category=Accessories&search=bag' },
          { name: 'Jewelry', to: '/products?category=Accessories&search=jewelry' },
          { name: 'Belts', to: '/products?category=Accessories&search=belt' },
          { name: 'Hats', to: '/products?category=Accessories&search=hat' },
          { name: 'Scarves', to: '/products?category=Accessories&search=scarf' }
        ]
      }
    ]
  },
  {
    name: 'KIDS',
    to: '/products?category=Kids',
    sections: [
      {
        title: 'CLOTHING',
        items: [
          { name: 'T-Shirts', to: '/products?category=Kids&search=t-shirt' },
          { name: 'Dresses', to: '/products?category=Kids&search=dress' },
          { name: 'Trousers', to: '/products?category=Kids&search=trousers' },
          { name: 'Jeans', to: '/products?category=Kids&search=jeans' },
          { name: 'Jackets', to: '/products?category=Kids&search=jacket' },
          { name: 'Hoodies', to: '/products?category=Kids&search=hoodie' }
        ]
      },
      {
        title: 'SHOES',
        items: [
          { name: 'Sneakers', to: '/products?category=Shoes&search=sneakers' },
          { name: 'Sandals', to: '/products?category=Shoes&search=sandals' },
          { name: 'Boots', to: '/products?category=Shoes&search=boots' }
        ]
      },
      {
        title: 'ACCESSORIES',
        items: [
          { name: 'Bags', to: '/products?category=Accessories&search=bag' },
          { name: 'Hats', to: '/products?category=Accessories&search=hat' },
          { name: 'Belts', to: '/products?category=Accessories&search=belt' }
        ]
      }
    ]
  }
]
</script>

<template>
  <header class="sticky top-0 z-40 w-full backdrop-blur-md bg-[#faf8f4]/85 border-b border-stone-200/60 transition-all duration-300">
    <!-- Non-overlapping 3-Column Grid Layout -->
    <div class="mx-auto grid grid-cols-3 items-center h-16 lg:h-[72px] w-full px-4 sm:px-6 lg:px-12 max-w-7xl relative">
      <!-- 1. LEFT COLUMN: Category Links (Desktop) / Hamburger Trigger (Mobile/Tablet) -->
      <div class="flex items-center justify-start min-w-0">
        <!-- Mobile/Tablet Drawer Toggle -->
        <button
          type="button"
          aria-label="Open navigation menu"
          class="flex lg:hidden h-10 w-10 -ml-2 items-center justify-center text-stone-850 hover:text-stone-500 transition-colors"
          @click="isMenuOpen = true"
        >
          <svg
            class="h-5 w-5 stroke-[1.4]"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
          >
            <line x1="3" y1="7" x2="21" y2="7" />
            <line x1="3" y1="12" x2="21" y2="12" />
            <line x1="3" y1="17" x2="21" y2="17" />
          </svg>
        </button>

        <!-- Desktop Category Links -->
        <nav aria-label="Main categories" class="hidden lg:flex items-center gap-4 xl:gap-7 min-w-0">
          <NuxtLink
            to="/products?sort=newest"
            class="group/link relative py-5 text-xs font-medium uppercase tracking-[0.16em] xl:tracking-[0.2em] text-stone-800 transition-colors duration-200 hover:text-stone-500 shrink-0"
          >
            NEW IN
            <span class="absolute bottom-3.5 left-0 h-[1.5px] w-0 bg-stone-800 transition-all duration-250 ease-out group-hover/link:w-full" />
          </NuxtLink>

          <!-- Category items with Mega-menu -->
          <div
            v-for="category in categories"
            :key="category.name"
            class="group static"
          >
            <NuxtLink
              :to="category.to"
              class="relative block py-5 text-xs font-medium uppercase tracking-[0.16em] xl:tracking-[0.2em] text-stone-800 transition-colors duration-200 hover:text-stone-500 shrink-0"
            >
              {{ category.name }}
              <span class="absolute bottom-3.5 left-0 h-[1.5px] w-0 bg-stone-800 transition-all duration-250 ease-out group-hover:w-full" />
            </NuxtLink>

            <!-- Luxury Editorial Mega Menu Dropdown -->
            <div
              class="invisible absolute left-0 right-0 top-full border-t border-stone-200/60 bg-[#faf8f4]/95 backdrop-blur-md opacity-0 shadow-[0_20px_40px_rgba(33,31,29,0.06)] transition-all duration-200 ease-out group-hover:visible group-hover:opacity-100"
            >
              <div class="mx-auto max-w-6xl px-8 lg:px-12 py-8 lg:py-9">
                <div class="mb-7 flex items-center justify-between border-b border-stone-200/60 pb-4">
                  <div>
                    <span class="text-[10px] font-semibold tracking-[0.22em] text-[#806344] uppercase">
                      COLLECTION
                    </span>
                    <h3 class="mt-1 font-serif text-2xl tracking-[0.05em] text-stone-900">
                      {{ category.name }}
                    </h3>
                  </div>

                  <NuxtLink
                    :to="category.to"
                    class="group/viewall inline-flex items-center gap-1.5 text-xs font-medium tracking-[0.16em] text-[#806344] transition-colors hover:text-stone-900"
                  >
                    <span>VIEW ALL</span>
                    <svg
                      class="h-3.5 w-3.5 stroke-[1.5] transition-transform duration-200 group-hover/viewall:translate-x-1"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                    >
                      <path d="M5 12h14" />
                      <path d="m13 6 6 6-6 6" />
                    </svg>
                  </NuxtLink>
                </div>

                <div class="grid grid-cols-3 gap-8 lg:gap-12">
                  <div
                    v-for="section in category.sections"
                    :key="section.title"
                  >
                    <h4 class="mb-3 text-[11px] font-semibold uppercase tracking-[0.16em] text-[#806344]">
                      {{ section.title }}
                    </h4>

                    <div class="flex flex-col gap-2">
                      <NuxtLink
                        v-for="item in section.items"
                        :key="item.name"
                        :to="item.to"
                        class="w-fit text-xs text-stone-600 transition-colors duration-150 hover:text-stone-950"
                      >
                        {{ item.name }}
                      </NuxtLink>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </nav>
      </div>

      <!-- 2. CENTER COLUMN: Brand Logo (Naturally centered in grid flow, zero collision) -->
      <div class="flex items-center justify-center min-w-0 text-center">
        <NuxtLink
          to="/"
          class="group inline-flex items-center justify-center py-2 min-w-0 transition-opacity hover:opacity-80"
          aria-label="AFRIMART Home"
        >
          <span class="font-serif text-base sm:text-lg md:text-xl lg:text-2xl font-normal tracking-[0.2em] sm:tracking-[0.25em] text-stone-900 uppercase truncate select-none">
            AFRIMART
          </span>
        </NuxtLink>
      </div>

      <!-- 3. RIGHT COLUMN: Utility Actions (Search, AI Pill, Profile, Cart) -->
      <div class="flex items-center justify-end min-w-0 gap-2 sm:gap-3 lg:gap-3.5 xl:gap-5">
        <!-- Desktop Integrated Search & AI Pill (lg+) -->
        <form
          class="relative hidden lg:flex items-center min-w-0"
          @submit.prevent="submitSearch"
        >
          <div class="relative flex items-center rounded-full border border-stone-300/70 bg-white/60 pl-3 pr-1.5 py-1 backdrop-blur-xs transition-all duration-300 focus-within:border-stone-800 focus-within:bg-white focus-within:shadow-xs hover:border-stone-400">
            <svg
              class="h-3.5 w-3.5 text-stone-500 mr-2 shrink-0 stroke-[1.6]"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
            >
              <circle cx="11" cy="11" r="7" />
              <path d="m20 20-4-4" />
            </svg>

            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search or ask AI..."
              class="w-28 xl:w-44 2xl:w-52 bg-transparent text-xs text-stone-800 placeholder-stone-400 focus:outline-none tracking-wide"
            />

            <!-- Distinct Minimalist AI Pill Button -->
            <button
              type="button"
              title="Ask AI Stylist"
              class="ml-1 inline-flex items-center gap-1 rounded-full bg-stone-900 px-2 py-0.5 xl:px-2.5 xl:py-1 text-[10px] font-medium uppercase tracking-[0.12em] text-stone-100 transition-all duration-200 hover:bg-stone-700 hover:scale-[1.02] active:scale-95 shadow-2xs cursor-pointer shrink-0"
              @click="handleAiSearch"
            >
              <svg class="h-2.5 w-2.5 text-amber-300 fill-amber-300/40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="m12 3 1.912 5.885L20 10.8l-4.706 3.914L16.8 21 12 17.5 7.2 21l1.506-6.286L4 10.8l6.088-1.915L12 3Z"/>
              </svg>
              <span>AI</span>
            </button>
          </div>
        </form>

        <!-- Tablet/Mobile Search Icon Toggle (< lg) -->
        <button
          type="button"
          aria-label="Toggle search bar"
          class="flex lg:hidden h-9 w-9 shrink-0 items-center justify-center text-stone-800 hover:text-stone-500 transition-colors"
          @click="isSearchOpen = !isSearchOpen"
        >
          <svg class="h-4.5 w-4.5 stroke-[1.4]" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <circle cx="11" cy="11" r="7" />
            <path d="m20 20-4-4" />
          </svg>
        </button>

        <!-- User Account Profile Link -->
        <NuxtLink
          to="/account"
          aria-label="User Account"
          class="flex h-9 w-9 shrink-0 items-center justify-center text-stone-800 hover:text-stone-500 transition-colors"
        >
          <svg
            class="h-5 w-5 stroke-[1.4]"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
          >
            <circle cx="12" cy="7" r="3.5" />
            <path d="M5 21c.6-4 3.1-6 7-6s6.4 2 7 6" />
          </svg>
        </NuxtLink>

        <!-- Cart Link -->
        <NuxtLink
          id="header-cart-icon"
          to="/cart"
          aria-label="Shopping Cart"
          class="relative flex h-9 w-9 shrink-0 items-center justify-center text-stone-800 transition-all duration-200 hover:text-stone-500"
          :class="{ 'scale-115 text-stone-900': isCartBouncing }"
        >
          <svg
            class="h-5 w-5 stroke-[1.4] transition-transform duration-200"
            :class="{ 'scale-110 rotate-6': isCartBouncing }"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
          >
            <path d="M6 7h12l-1 13H7L6 7Z" />
            <path d="M9 7a3 3 0 0 1 6 0" />
          </svg>

          <!-- Minimalist Badge Counter -->
          <span
            v-if="cartCount"
            class="absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-stone-900 px-1 text-[9px] font-semibold tracking-tighter text-white transition-transform duration-200"
            :class="{ 'scale-125 bg-[#806344]': isCartBouncing }"
          >
            {{ cartCount }}
          </span>
        </NuxtLink>
      </div>
    </div>

    <!-- Tablet/Mobile Expandable Search Bar Popup -->
    <Transition
      enter-active-class="transition-all duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="isSearchOpen"
        class="border-t border-stone-200/60 bg-[#faf8f4]/95 backdrop-blur-md px-4 py-3 sm:px-6 shadow-sm lg:hidden"
      >
        <div class="mx-auto max-w-xl">
          <form
            class="flex items-center"
            @submit.prevent="submitSearch"
          >
            <div class="relative flex w-full items-center rounded-full border border-stone-300 bg-white/90 pl-3.5 pr-1.5 py-1.5 shadow-2xs focus-within:border-stone-800 transition-colors">
              <svg class="h-4 w-4 text-stone-500 mr-2 shrink-0 stroke-[1.5]" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <circle cx="11" cy="11" r="7" />
                <path d="m20 20-4-4" />
              </svg>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search or ask AI..."
                class="w-full bg-transparent text-xs sm:text-sm text-stone-800 placeholder-stone-400 focus:outline-none tracking-wide"
                autofocus
              />
              <button
                type="button"
                title="Ask AI Stylist"
                class="ml-1.5 inline-flex items-center gap-1 rounded-full bg-stone-900 px-2.5 py-1 text-[10px] font-medium uppercase tracking-[0.12em] text-stone-100 shrink-0"
                @click="handleAiSearch"
              >
                <svg class="h-2.5 w-2.5 text-amber-300 fill-amber-300/40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="m12 3 1.912 5.885L20 10.8l-4.706 3.914L16.8 21 12 17.5 7.2 21l1.506-6.286L4 10.8l6.088-1.915L12 3Z"/>
                </svg>
                <span>AI</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- Responsive Mobile/Tablet Drawer Navigation (Teleported to body for zero layout bleed) -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition-opacity duration-200 ease-out"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-opacity duration-150 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="isMounted && isMenuOpen"
          class="fixed inset-0 z-50 overflow-hidden lg:hidden"
          role="dialog"
          aria-modal="true"
          aria-label="Mobile navigation menu"
        >
          <!-- Backdrop -->
          <div
            class="fixed inset-0 bg-stone-900/40 backdrop-blur-xs transition-opacity duration-300"
            @click="closeOverlays"
          />

          <!-- Drawer Panel -->
          <div class="fixed inset-y-0 left-0 z-50 flex max-w-full">
            <div class="w-screen max-w-xs sm:max-w-sm bg-[#faf8f4] border-r border-stone-200/80 shadow-2xl flex flex-col justify-between">
              <!-- Drawer Header & Navigation -->
              <div>
                <div class="flex h-16 items-center justify-between border-b border-stone-200/60 px-5">
                  <NuxtLink
                    to="/"
                    class="font-serif text-lg tracking-[0.25em] text-stone-900 uppercase"
                    @click="closeOverlays"
                  >
                    AFRIMART
                  </NuxtLink>
                  <button
                    type="button"
                    aria-label="Close menu"
                    class="flex h-9 w-9 items-center justify-center text-stone-600 hover:text-stone-900 transition-colors"
                    @click="closeOverlays"
                  >
                    <svg class="h-5 w-5 stroke-[1.5]" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                      <path d="M18 6 6 18" />
                      <path d="m6 6 12 12" />
                    </svg>
                  </button>
                </div>

                <!-- Drawer Search Bar -->
                <div class="p-4 border-b border-stone-200/60">
                  <form
                    class="relative flex items-center"
                    @submit.prevent="submitSearch"
                  >
                    <div class="relative flex w-full items-center rounded-full border border-stone-300 bg-white/80 pl-3 pr-1.5 py-1.5">
                      <svg class="h-3.5 w-3.5 text-stone-500 mr-2 shrink-0 stroke-[1.5]" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                        <circle cx="11" cy="11" r="7" />
                        <path d="m20 20-4-4" />
                      </svg>
                      <input
                        v-model="searchQuery"
                        type="text"
                        placeholder="Search or ask AI..."
                        class="w-full bg-transparent text-xs text-stone-800 placeholder-stone-400 focus:outline-none"
                      />
                      <button
                        type="button"
                        class="ml-1 inline-flex items-center gap-1 rounded-full bg-stone-900 px-2.5 py-1 text-[10px] font-medium uppercase tracking-[0.12em] text-stone-100"
                        @click="handleAiSearch"
                      >
                        <span>AI</span>
                      </button>
                    </div>
                  </form>
                </div>

                <!-- Category Links & Accordions -->
                <div class="overflow-y-auto px-5 py-4 divide-y divide-stone-200/60 max-h-[calc(100dvh-14rem)]">
                  <!-- New In -->
                  <NuxtLink
                    to="/products?sort=newest"
                    class="flex items-center justify-between py-3.5 text-xs font-medium tracking-[0.18em] text-stone-900 hover:text-stone-500"
                    @click="closeOverlays"
                  >
                    <span>NEW IN</span>
                    <span class="text-[10px] text-[#806344]">EXPLORE →</span>
                  </NuxtLink>

                  <!-- Categories Accordion -->
                  <div
                    v-for="category in categories"
                    :key="category.name"
                    class="py-3.5"
                  >
                    <div class="flex items-center justify-between">
                      <NuxtLink
                        :to="category.to"
                        class="text-xs font-medium tracking-[0.18em] text-stone-900 hover:text-stone-500"
                        @click="closeOverlays"
                      >
                        {{ category.name }}
                      </NuxtLink>

                      <button
                        type="button"
                        class="flex h-7 w-7 items-center justify-center text-stone-500 hover:text-stone-900"
                        :aria-label="`Toggle ${category.name}`"
                        @click="toggleCategory(category.name)"
                      >
                        <svg
                          class="h-4 w-4 stroke-[1.5] transition-transform duration-200"
                          :class="{ 'rotate-180': openCategories[category.name] }"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                        >
                          <path d="m6 9 6 6 6-6" />
                        </svg>
                      </button>
                    </div>

                    <!-- Sub-items -->
                    <div
                      v-if="openCategories[category.name]"
                      class="mt-3 space-y-3 rounded-lg bg-[#eee8df]/50 p-3"
                    >
                      <div
                        v-for="section in category.sections"
                        :key="section.title"
                      >
                        <p class="mb-1.5 text-[9px] font-semibold tracking-[0.16em] text-[#806344] uppercase">
                          {{ section.title }}
                        </p>
                        <div class="grid grid-cols-2 gap-1.5">
                          <NuxtLink
                            v-for="item in section.items"
                            :key="item.name"
                            :to="item.to"
                            class="rounded py-1 text-xs text-stone-700 hover:text-stone-950 truncate"
                            @click="closeOverlays"
                          >
                            {{ item.name }}
                          </NuxtLink>
                        </div>
                      </div>

                      <NuxtLink
                        :to="category.to"
                        class="inline-block pt-1 text-[11px] font-semibold tracking-[0.14em] text-[#806344] hover:underline"
                        @click="closeOverlays"
                      >
                        VIEW ALL {{ category.name }} →
                      </NuxtLink>
                    </div>
                  </div>

                  <!-- AI Stylist Link -->
                  <NuxtLink
                    to="/ai"
                    class="flex items-center justify-between py-3.5 text-xs font-medium tracking-[0.18em] text-stone-900 hover:text-stone-500"
                    @click="closeOverlays"
                  >
                    <span class="flex items-center gap-1.5">
                      <span>AI STYLIST</span>
                      <span class="rounded bg-amber-200/60 px-1 py-0.5 text-[9px] font-bold text-amber-900">NEW</span>
                    </span>
                    <span class="text-[10px] text-[#806344]">ASK →</span>
                  </NuxtLink>
                </div>
              </div>

              <!-- Drawer Footer -->
              <div class="border-t border-stone-200/60 bg-[#f5f1e9] p-5 space-y-3">
                <div class="flex items-center justify-between text-xs tracking-wider text-stone-600">
                  <NuxtLink to="/account" class="hover:text-stone-950" @click="closeOverlays">
                    My Account
                  </NuxtLink>
                  <NuxtLink to="/orders" class="hover:text-stone-950" @click="closeOverlays">
                    Orders & Tracking
                  </NuxtLink>
                </div>
</div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </header>
</template>
