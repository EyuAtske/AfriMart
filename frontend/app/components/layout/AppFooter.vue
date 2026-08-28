<template>
  <footer class="bg-[#1d1a18] text-white">
    <div class="mx-auto max-w-[1800px] px-6 py-16 lg:px-12 lg:py-24">

      <!-- Main footer -->
      <div class="grid gap-12 sm:grid-cols-2 lg:grid-cols-4 lg:gap-16">

        <!-- Brand -->
        <div>
          <NuxtLink
            to="/"
            class="font-serif text-3xl tracking-[0.3em] text-[#f5f1e9]"
          >
            AFRIMART
          </NuxtLink>

          <p
            class="mt-7 max-w-xs text-base leading-8 text-[#998d82]"
          >
            Discover clothing from independent sellers.
            Buy, sell, and find your style in one marketplace.
          </p>
        </div>

        <!-- Shop -->
        <div>
          <h3
            class="text-sm font-semibold tracking-[0.16em] text-[#f5f1e9]"
          >
            SHOP
          </h3>

          <nav class="mt-7 flex flex-col gap-5">
            <NuxtLink
              to="/products"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              New Arrivals
            </NuxtLink>

            <NuxtLink
              to="/products?category=Men"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              Men
            </NuxtLink>

            <NuxtLink
              to="/products?category=Women"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              Women
            </NuxtLink>

            <NuxtLink
              to="/products?category=Kids"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              Kids
            </NuxtLink>

            <NuxtLink
              to="/cart"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              Shopping Cart
            </NuxtLink>
          </nav>
        </div>

        <!-- Sell -->
        <div>
          <h3
            class="text-sm font-semibold tracking-[0.16em] text-[#f5f1e9]"
          >
            SELL
          </h3>

          <nav class="mt-7 flex flex-col gap-5">
            <NuxtLink
              to="/sell"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              Start Selling
            </NuxtLink>

            <NuxtLink
              to="/shop"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              My Shop
            </NuxtLink>

            <NuxtLink
              to="/seller/products"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              My Products
            </NuxtLink>

            <NuxtLink
              to="/seller/orders"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              Seller Orders
            </NuxtLink>
          </nav>
        </div>

        <!-- Account -->
        <div>
          <h3
            class="text-sm font-semibold tracking-[0.16em] text-[#f5f1e9]"
          >
            ACCOUNT
          </h3>

          <nav class="mt-7 flex flex-col gap-5">
            <NuxtLink
              to="/profile"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              My Profile
            </NuxtLink>

            <NuxtLink
              to="/orders"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              My Orders
            </NuxtLink>

            <NuxtLink
              to="/ai"
              class="text-base text-[#998d82] transition-colors hover:text-white"
            >
              AI Shopping Assistant
            </NuxtLink>
          </nav>
        </div>

      </div>

    </div>
  </footer>
</template>