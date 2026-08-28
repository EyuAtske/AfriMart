<script setup lang="ts">
import LoadMoreButton from '~/components/LoadMoreButton.vue'
import CategoryCard from '~/components/marketplace/CategoryCard.vue'

useSeoMeta({
  title: 'Afrimart — Curated Fashion & Independent Sellers Marketplace',
  description: 'Discover clothing and unique items from independent sellers on Afrimart. Quality pieces, ethical shopping, one unified marketplace.',
  ogTitle: 'Afrimart — Curated Fashion Marketplace',
  ogDescription: 'Discover clothing and unique items from independent sellers on Afrimart.',
  ogImage: '/images/herotemp.png'
})

const productsToShow = ref(10)

const loadMore = () => {
  productsToShow.value += 5
}
</script>

<template>
  <main>
    <!-- Hero section -->
    <section class="relative min-h-screen overflow-hidden pt-11 lg:pt-14">
      <!-- Hero pic -->
      <img
        src="/images/herotemp.png"
        alt="Afrimart fashion collection"
        class="absolute inset-0 h-full w-full object-cover object-center"
      />

      <!-- Soft overlay only on the text side -->
      <div
        class="absolute inset-y-0 left-0 w-full bg-gradient-to-t from-[#f5f1e9]/90 via-[#f5f1e9]/35 to-transparent sm:w-3/4 sm:bg-gradient-to-r sm:from-[#f5f1e9]/90 sm:via-[#f5f1e9]/50 sm:to-transparent lg:w-1/2 lg:from-[#f5f1e9]/95 lg:via-[#f5f1e9]/55 lg:to-transparent"
      ></div>

      <!-- Hero content texts -->
      <div
        class="relative z-10 flex min-h-screen items-end px-5 pb-12 pt-24 sm:px-8 sm:pb-16 sm:pt-32 lg:items-center lg:px-16 lg:pb-10"
      >
        <div
          class="max-w-xl rounded-md bg-[#f5f1e9]/75 p-5 backdrop-blur-[2px] sm:bg-transparent sm:p-0 sm:backdrop-blur-0 lg:max-w-xl"
        >
          <!-- Heading -->
          <h1
            class="font-serif text-5xl leading-[0.95] tracking-[-0.03em] text-[#211f1d] sm:text-6xl lg:text-7xl xl:text-[6.2rem]"
          >
            Your style.<br />
            Your<br />
            marketplace.
          </h1>

          <!-- Description -->
          <p
            class="mt-6 max-w-md text-sm leading-6 text-[#4c4945] sm:mt-8 sm:text-base sm:leading-7 lg:mt-10 lg:text-lg lg:leading-8"
          >
            Discover clothing from independent sellers.
            <br class="hidden sm:block" />
            Quality pieces. Great prices. One marketplace.
          </p>

          <!-- Actions -->
          <div
            class="mt-7 flex flex-wrap items-center gap-5 sm:mt-8 sm:gap-7 lg:mt-10"
          >
            <!-- Shop Now -->
           <NuxtLink to="/products" class="group inline-flex items-center gap-4 rounded-full border border-[#806344] px-10 py-4 text-sm font-medium tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 focus:ring-offset-[#f5f1e9]">
             <span>SHOP NOW</span>
          </NuxtLink>

            <!-- Explore Items -->
            <NuxtLink to="/shop" class="group inline-flex items-center gap-4 rounded-full border border-[#806344] px-10 py-4 text-sm font-medium tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 focus:ring-offset-[#f5f1e9]">
              <span>SELL ITEMS</span>
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <!-- Categories section -->
    <section
      class="bg-[#f5f1e9] px-3 py-12 sm:px-4 sm:py-16 md:px-6 md:py-20 lg:px-12 lg:py-24"
    >
      <div class="grid grid-cols-3 gap-2 sm:gap-4 lg:gap-6">
        <CategoryCard
          title="Men"
          image="/images/mencat.jpg"
          to="/products?category=Men"
        />

        <CategoryCard
          title="Women"
          image="/images/womencat.jpg"
          to="/products?category=Women"
        />

        <CategoryCard
          title="Kids"
          image="/images/kidscat.jpg"
          to="/products?category=Kids"
        />
      </div>
    </section>

 <MarketplaceFeaturedProducts
  :limit="productsToShow"
/>

<LoadMoreButton
  @load-more="loadMore"
/>

<br>
    <MarketplaceAiSearch />
  </main>
</template>
