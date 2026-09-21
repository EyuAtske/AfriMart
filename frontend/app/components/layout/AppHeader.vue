<script setup lang="ts">
const isMenuOpen = ref(false)
const searchQuery = ref('')
const router = useRouter()
const { cartProducts } = useMarketplace()
const { isCartBouncing } = useFlyToCart()

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
  if (search) {
    await router.push({
      path: '/ai',
      query: { prompt: search }
    })
  } else {
    await router.push('/ai')
  }
  isMenuOpen.value = false
}

const submitSearch = async () => {
  const search = searchQuery.value.trim()
  if (!search) return

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

  isMenuOpen.value = false
}

const openCategories = ref<Record<string, boolean>>({
  MEN: true,
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
  <header class="fixed inset-x-0 top-0 z-50 w-full bg-[#f5f1e9]">
    <!-- Navbar -->
    <nav
      class="mx-auto flex h-11 items-center justify-between px-3 sm:px-5 lg:h-14 lg:px-12 relative"
    >
      <!-- Left side: Hamburger & Mobile Logo -->
      <div class="flex items-center gap-2 sm:gap-3">
        <!-- Mobile hamburger -->
        <button
          type="button"
          aria-label="Toggle menu"
          class="flex h-8 w-8 items-center justify-center text-[#302d29] lg:hidden"
          @click="isMenuOpen = !isMenuOpen"
        >
          <!-- Hamburger -->
          <svg
            v-if="!isMenuOpen"
            xmlns="http://www.w3.org/2000/svg"
            width="22"
            height="22"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
          >
            <path d="M4 7h16" />
            <path d="M4 12h16" />
            <path d="M4 17h16" />
          </svg>

          <!-- Close -->
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            width="22"
            height="22"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
          >
            <path d="M5 5l14 14" />
            <path d="M19 5L5 19" />
          </svg>
        </button>

        <!-- Logo (Mobile: Shifts to the left) -->
        <NuxtLink
          to="/"
          class="whitespace-nowrap font-serif text-base tracking-[0.15em] text-[#24211e] sm:text-lg lg:hidden"
        >
          AFRIMART
        </NuxtLink>
      </div>

      <!-- Desktop navigation -->
      <div class="hidden items-center gap-5 lg:flex xl:gap-7">
        <!-- New In -->
        <NuxtLink
          to="/"
          class="py-5 text-xs tracking-[0.08em] text-[#302d29] transition-opacity hover:opacity-60 xl:text-sm"
        >
          NEW IN
        </NuxtLink>

        <!-- Categories -->
        <div
          v-for="category in categories"
          :key="category.name"
          class="group static"
        >
          <!-- Category name -->
          <NuxtLink
            :to="category.to"
            class="block py-5 text-xs tracking-[0.08em] text-[#302d29] transition-opacity hover:opacity-60 xl:text-sm"
          >
            {{ category.name }}
          </NuxtLink>

          <!-- Mega menu -->
          <div
            class="invisible absolute left-0 right-0 top-full border-t border-[#ded8ce] bg-[#f5f1e9] opacity-0 shadow-[0_12px_30px_rgba(48,45,41,0.08)] transition-all duration-200 group-hover:visible group-hover:opacity-100"
          >
            <div class="mx-auto max-w-6xl px-12 py-8">
              <!-- Header -->
              <div
                class="mb-7 flex items-center justify-between border-b border-[#ded8ce] pb-5"
              >
                <div>
                  <p
                    class="text-xs tracking-[0.18em] text-[#806344]"
                  >
                    SHOP
                  </p>

                  <h3
                    class="mt-1 font-serif text-2xl text-[#211f1d]"
                  >
                    {{ category.name }}
                  </h3>
                </div>

                <NuxtLink
                  :to="category.to"
                  class="group/all flex items-center gap-2 text-xs tracking-[0.12em] text-[#806344]"
                >
                  VIEW ALL
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="17"
                    height="17"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    class="transition-transform group-hover/all:translate-x-1"
                  >
                    <path d="M5 12h14" />
                    <path d="m13 6 6 6-6 6" />
                  </svg>
                </NuxtLink>
              </div>

              <!-- Category columns -->
              <div class="grid grid-cols-3 gap-12">
                <div
                  v-for="section in category.sections"
                  :key="section.title"
                >
                  <!-- Section title -->
                  <h4
                    class="mb-4 text-xs font-medium tracking-[0.14em] text-[#806344]"
                  >
                    {{ section.title }}
                  </h4>

                  <!-- Links -->
                  <div class="flex flex-col gap-2.5">
                    <NuxtLink
                      v-for="item in section.items"
                      :key="item.name"
                      :to="item.to"
                      class="w-fit text-sm text-[#302d29] transition-colors hover:text-[#806344]"
                    >
                      {{ item.name }}
                    </NuxtLink>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Logo (Desktop: Absolutely Centered) -->
      <NuxtLink
        to="/"
        class="absolute left-1/2 -translate-x-1/2 hidden whitespace-nowrap font-serif lg:text-2xl lg:tracking-[0.3em] text-[#24211e] lg:block"
      >
        AFRIMART
      </NuxtLink>

      <!-- Right actions -->
      <div class="ml-auto flex items-center gap-2.5 sm:gap-4 lg:gap-5">
        <!-- Search Input with AI support -->
        <form
          class="relative flex items-center"
          @submit.prevent="submitSearch"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            class="absolute left-3 text-[#302d29] pointer-events-none"
          >
            <circle cx="11" cy="11" r="7" />
            <path d="m20 20-4-4" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search or ask AI..."
            class="w-36 sm:w-60 lg:w-72 h-9 sm:h-10 pl-9 pr-9 text-xs sm:text-sm bg-transparent rounded-full border border-[#d9d0c4] text-[#302d29] placeholder-[#92877b] focus:outline-none focus:border-[#806344] transition-all"
          />
          <button
            type="button"
            title="Ask AI Assistant"
            class="absolute right-1.5 flex h-6 w-6 sm:h-7 sm:w-7 items-center justify-center rounded-full bg-[#806344]/15 hover:bg-[#806344] hover:text-white text-[#806344] text-[10px] sm:text-xs font-bold transition-all"
            @click="handleAiSearch"
          >
            AI
          </button>
        </form>

        <!-- Account -->
        <NuxtLink
          to="/account"
          aria-label="Account"
          class="text-[#302d29] transition-opacity hover:opacity-60"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="21"
            height="21"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
          >
            <circle cx="12" cy="7" r="3.5" />
            <path d="M5 21c.6-4 3.1-6 7-6s6.4 2 7 6" />
          </svg>
        </NuxtLink>

        <!-- Cart -->
        <NuxtLink
          id="header-cart-icon"
          to="/cart"
          aria-label="Shopping cart"
          class="relative text-[#302d29] transition-all duration-200 hover:opacity-60"
          :class="{ 'scale-125 text-[#806344]': isCartBouncing }"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="21"
            height="21"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            class="transition-transform duration-200"
            :class="{ 'scale-110 rotate-6': isCartBouncing }"
          >
            <path d="M5 8h14l-1 12H6L5 8Z" />
            <path d="M9 8a3 3 0 0 1 6 0" />
          </svg>

          <span
            v-if="cartCount"
            class="absolute -right-2 -top-2 flex h-5 min-w-5 items-center justify-center rounded-full bg-[#806344] px-1 text-[10px] font-semibold leading-none text-white transition-transform duration-200"
            :class="{ 'scale-130 bg-[#211f1d]': isCartBouncing }"
          >
            {{ cartCount }}
          </span>
        </NuxtLink>
      </div>
    </nav>

    <!-- Mobile menu -->
    <div
      v-if="isMenuOpen"
      class="border-t border-[#ded8ce] bg-[#f5f1e9] shadow-xl max-h-[calc(100vh-3.5rem)] overflow-y-auto lg:hidden"
    >
      <div class="px-5 py-3 divide-y divide-[#ded8ce]">
        <!-- New In -->
        <NuxtLink
          to="/"
          class="flex items-center justify-between py-3 text-sm font-medium tracking-[0.14em] text-[#302d29] hover:text-[#806344]"
          @click="isMenuOpen = false"
        >
          <span>NEW IN</span>
          <span class="text-xs text-[#806344]">EXPLORE →</span>
        </NuxtLink>

        <!-- Mobile categories accordion -->
        <div
          v-for="category in categories"
          :key="category.name"
          class="py-3"
        >
          <!-- Category header toggle -->
          <div class="flex items-center justify-between">
            <NuxtLink
              :to="category.to"
              class="text-sm font-medium tracking-[0.14em] text-[#302d29] hover:text-[#806344]"
              @click="isMenuOpen = false"
            >
              {{ category.name }}
            </NuxtLink>

            <button
              type="button"
              class="flex h-8 w-8 items-center justify-center rounded-full text-[#756a60] hover:bg-[#ded8ce]/60"
              :aria-label="`Toggle ${category.name} menu`"
              @click="toggleCategory(category.name)"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                class="transition-transform duration-200"
                :class="{ 'rotate-180': openCategories[category.name] }"
              >
                <path d="m6 9 6 6 6-6" />
              </svg>
            </button>
          </div>

          <!-- Category items: responsive 2-column grid -->
          <div
            v-if="openCategories[category.name]"
            class="mt-3 space-y-3.5 rounded-lg bg-[#eee8df]/60 p-3.5"
          >
            <div
              v-for="section in category.sections"
              :key="section.title"
            >
              <p class="mb-1.5 text-[10px] font-semibold tracking-[0.14em] text-[#806344] uppercase">
                {{ section.title }}
              </p>

              <div class="grid grid-cols-2 gap-x-2 gap-y-1.5">
                <NuxtLink
                  v-for="item in section.items"
                  :key="item.name"
                  :to="item.to"
                  class="rounded px-1.5 py-1 text-xs text-[#5d574f] hover:bg-[#ded8ce]/60 hover:text-[#211f1d] truncate"
                  @click="isMenuOpen = false"
                >
                  {{ item.name }}
                </NuxtLink>
              </div>
            </div>

            <NuxtLink
              :to="category.to"
              class="inline-flex items-center gap-1.5 pt-1 text-xs font-semibold tracking-[0.12em] text-[#806344] hover:underline"
              @click="isMenuOpen = false"
            >
              VIEW ALL {{ category.name }} →
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
