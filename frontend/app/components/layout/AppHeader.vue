<script setup lang="ts">
const isMenuOpen = ref(false)
const searchQuery = ref('')

const categories = [
  {
    name: 'MEN',
    to: '/men',
    sections: [
      {
        title: 'CLOTHING',
        items: [
          { name: 'T-Shirts', to: '/men/clothing/t-shirts' },
          { name: 'Shirts', to: '/men/clothing/shirts' },
          { name: 'Trousers', to: '/men/clothing/trousers' },
          { name: 'Jeans', to: '/men/clothing/jeans' },
          { name: 'Jackets', to: '/men/clothing/jackets' },
          { name: 'Hoodies', to: '/men/clothing/hoodies' },
        ],
      },
      {
        title: 'SHOES',
        items: [
          { name: 'Sneakers', to: '/men/shoes/sneakers' },
          { name: 'Formal Shoes', to: '/men/shoes/formal' },
          { name: 'Boots', to: '/men/shoes/boots' },
          { name: 'Sandals', to: '/men/shoes/sandals' },
        ],
      },
      {
        title: 'ACCESSORIES',
        items: [
          { name: 'Bags', to: '/men/accessories/bags' },
          { name: 'Watches', to: '/men/accessories/watches' },
          { name: 'Belts', to: '/men/accessories/belts' },
          { name: 'Hats', to: '/men/accessories/hats' },
        ],
      },
    ],
  },

  {
    name: 'WOMEN',
    to: '/women',
    sections: [
      {
        title: 'CLOTHING',
        items: [
          { name: 'Dresses', to: '/women/clothing/dresses' },
          { name: 'Tops', to: '/women/clothing/tops' },
          { name: 'Trousers', to: '/women/clothing/trousers' },
          { name: 'Jeans', to: '/women/clothing/jeans' },
          { name: 'Skirts', to: '/women/clothing/skirts' },
          { name: 'Jackets', to: '/women/clothing/jackets' },
          { name: 'Sweaters', to: '/women/clothing/sweaters' },
        ],
      },
      {
        title: 'SHOES',
        items: [
          { name: 'Sneakers', to: '/women/shoes/sneakers' },
          { name: 'Heels', to: '/women/shoes/heels' },
          { name: 'Boots', to: '/women/shoes/boots' },
          { name: 'Sandals', to: '/women/shoes/sandals' },
          { name: 'Flats', to: '/women/shoes/flats' },
        ],
      },
      {
        title: 'ACCESSORIES',
        items: [
          { name: 'Bags', to: '/women/accessories/bags' },
          { name: 'Jewelry', to: '/women/accessories/jewelry' },
          { name: 'Belts', to: '/women/accessories/belts' },
          { name: 'Hats', to: '/women/accessories/hats' },
          { name: 'Scarves', to: '/women/accessories/scarves' },
        ],
      },
    ],
  },

  {
    name: 'KIDS',
    to: '/kids',
    sections: [
      {
        title: 'CLOTHING',
        items: [
          { name: 'T-Shirts', to: '/kids/clothing/t-shirts' },
          { name: 'Dresses', to: '/kids/clothing/dresses' },
          { name: 'Trousers', to: '/kids/clothing/trousers' },
          { name: 'Jeans', to: '/kids/clothing/jeans' },
          { name: 'Jackets', to: '/kids/clothing/jackets' },
          { name: 'Hoodies', to: '/kids/clothing/hoodies' },
        ],
      },
      {
        title: 'SHOES',
        items: [
          { name: 'Sneakers', to: '/kids/shoes/sneakers' },
          { name: 'Sandals', to: '/kids/shoes/sandals' },
          { name: 'Boots', to: '/kids/shoes/boots' },
        ],
      },
      {
        title: 'ACCESSORIES',
        items: [
          { name: 'Bags', to: '/kids/accessories/bags' },
          { name: 'Hats', to: '/kids/accessories/hats' },
          { name: 'Belts', to: '/kids/accessories/belts' },
        ],
      },
    ],
  },
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
        <!--z Search Input -->
        <div class="relative flex items-center">
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
            placeholder="Search products or ask AI ✦..."
            class="w-36 sm:w-56 lg:w-64 h-9 sm:h-10 pl-9 pr-3.5 text-xs sm:text-sm bg-transparent rounded-full border border-[#d7d0c4] text-[#302d29] placeholder-[#8c857b] focus:outline-none focus:border-[#806344] transition-all"
          />
        </div>

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
          to="/cart"
          aria-label="Shopping cart"
          class="relative text-[#302d29] transition-opacity hover:opacity-60"
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
            <path d="M5 8h14l-1 12H6L5 8Z" />
            <path d="M9 8a3 3 0 0 1 6 0" />
          </svg>

          <span
            class="absolute -right-2 -top-2 flex h-5 w-5 items-center justify-center rounded-full bg-[#806344] text-[10px] font-medium text-white"
          >
            2
          </span>
        </NuxtLink>
      </div>
    </nav>

    <!-- Mobile menu -->
    <div
      v-if="isMenuOpen"
      class="border-t border-[#ded8ce] bg-[#f5f1e9] lg:hidden"
    >
      <div class="px-6 py-4">
        <!-- New In -->
        <NuxtLink
          to="/"
          class="block border-b border-[#ded8ce] py-4 text-sm tracking-[0.12em] text-[#302d29]"
          @click="isMenuOpen = false"
        >
          NEW IN
        </NuxtLink>

        <!-- Mobile categories -->
        <div
          v-for="category in categories"
          :key="category.name"
          class="border-b border-[#ded8ce] py-4"
        >
          <!-- Main category -->
          <NuxtLink
            :to="category.to"
            class="block text-sm tracking-[0.12em] text-[#302d29]"
            @click="isMenuOpen = false"
          >
            {{ category.name }}
          </NuxtLink>

          <!-- Sections -->
          <div class="mt-3 grid grid-cols-3 gap-4">
            <div
              v-for="section in category.sections"
              :key="section.title"
            >
              <p
                class="mb-2 text-[10px] tracking-[0.12em] text-[#806344]"
              >
                {{ section.title }}
              </p>

              <div class="flex flex-col gap-1.5">
                <NuxtLink
                  v-for="item in section.items"
                  :key="item.name"
                  :to="item.to"
                  class="text-xs text-[#5d574f]"
                  @click="isMenuOpen = false"
                >
                  {{ item.name }}
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>