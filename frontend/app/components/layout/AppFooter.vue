<template>
  <footer class="w-full bg-[#1d1a18] text-white">
    <div class="mx-auto w-full max-w-[1800px] px-4 py-10 sm:px-6 sm:py-16 lg:px-12 lg:py-20">

      <!-- Main footer: 3 equal columns covering the full span -->
      <div class="grid w-full grid-cols-3 gap-4 sm:gap-10 lg:gap-16">

        <!-- Shop -->
        <div class="min-w-0">
          <h3
            class="text-[11px] sm:text-xs md:text-sm font-semibold tracking-[0.12em] sm:tracking-[0.16em] text-[#f5f1e9]"
          >
            SHOP
          </h3>

          <nav class="mt-2 sm:mt-4 lg:mt-6 flex flex-col gap-1.5 sm:gap-3 lg:gap-4">
            <NuxtLink
              to="/products"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              New Arrivals
            </NuxtLink>

            <NuxtLink
              to="/products?category=Men"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              Men
            </NuxtLink>

            <NuxtLink
              to="/products?category=Women"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              Women
            </NuxtLink>

            <NuxtLink
              to="/products?category=Kids"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              Kids
            </NuxtLink>

            <NuxtLink
              to="/cart"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              Shopping Cart
            </NuxtLink>
          </nav>
        </div>

        <!-- Sell -->
        <div class="min-w-0">
          <h3
            class="text-[11px] sm:text-xs md:text-sm font-semibold tracking-[0.12em] sm:tracking-[0.16em] text-[#f5f1e9]"
          >
            SELL
          </h3>

          <nav class="mt-2 sm:mt-4 lg:mt-6 flex flex-col gap-1.5 sm:gap-3 lg:gap-4">
            <NuxtLink
              to="/sell"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              Start Selling
            </NuxtLink>

            <NuxtLink
              to="/shop"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              My Shop
            </NuxtLink>

            <NuxtLink
              to="/seller/products"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              My Products
            </NuxtLink>

            <NuxtLink
              to="/seller/orders"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              Seller Orders
            </NuxtLink>
          </nav>
        </div>

        <!-- Account -->
        <div class="min-w-0">
          <h3
            class="text-[11px] sm:text-xs md:text-sm font-semibold tracking-[0.12em] sm:tracking-[0.16em] text-[#f5f1e9]"
          >
            ACCOUNT
          </h3>

          <nav class="mt-2 sm:mt-4 lg:mt-6 flex flex-col gap-1.5 sm:gap-3 lg:gap-4">
            <NuxtLink
              to="/profile"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              My Profile
            </NuxtLink>

            <NuxtLink
              to="/orders"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              My Orders
            </NuxtLink>

            <NuxtLink
              to="/ai"
              class="text-[10px] sm:text-xs md:text-sm lg:text-base text-[#998d82] transition-colors hover:text-white truncate"
            >
              AI Assistant
            </NuxtLink>
          </nav>
        </div>

      </div>

    </div>
  </footer>
</template>