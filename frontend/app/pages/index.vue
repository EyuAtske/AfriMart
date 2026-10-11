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
    <!-- Hero section: Immersive Full-Width Centered Banner -->
    <section class="relative w-full overflow-hidden min-h-[75vh] sm:min-h-[80vh] lg:min-h-[85vh] xl:min-h-[88vh] flex items-center justify-center">
      <!-- Full-Width Background Imagery -->
      <img
        src="/images/herotemp.png"
        alt="Afrimart curated fashion collection"
        class="absolute inset-0 h-full w-full object-cover object-center"
      />

      <!-- Soft, Elegant Gradient Overlay (keeps photography bright, warm & natural) -->
      <div class="absolute inset-0 bg-gradient-to-b from-black/30 via-black/10 to-black/30" />

      <!-- Centered Content Container -->
      <div class="relative z-10 mx-auto w-full max-w-4xl px-5 sm:px-8 lg:px-12 py-20 sm:py-24 lg:py-32 text-center">
        <!-- Massive Serif Headline -->
        <h1 class="font-serif text-4xl sm:text-6xl md:text-7xl lg:text-8xl xl:text-[6.5rem] font-normal tracking-[-0.025em] text-white leading-[1.02] sm:leading-[0.96] drop-shadow-[0_4px_24px_rgba(0,0,0,0.65)]">
          Your style.<br />
          Your marketplace.
        </h1>

        <!-- Subtitle -->
        <p class="mt-6 sm:mt-8 text-base sm:text-lg md:text-xl text-white/90 leading-relaxed max-w-2xl mx-auto font-light drop-shadow-[0_2px_12px_rgba(0,0,0,0.7)]">
          Discover clothing from independent sellers. Quality pieces, ethical shopping, and authentic style unified in one curated marketplace.
        </p>

        <!-- Centered High-End CTA Buttons -->
        <div class="mt-8 sm:mt-10 flex flex-wrap items-center justify-center gap-4 sm:gap-6">
          <NuxtLink
            to="/products"
            class="bg-[#eee8df] text-[#4a3b2c] border border-[#806344]/80 px-8 sm:px-10 py-3.5 sm:py-4 text-xs uppercase tracking-[0.18em] font-semibold rounded-full shadow-lg hover:bg-[#806344] hover:text-white hover:border-[#806344] hover:scale-[1.03] active:scale-[0.98] transition-all duration-200 inline-flex items-center justify-center cursor-pointer"
          >
            SHOP NOW
          </NuxtLink>

          <NuxtLink
            to="/shop"
            class="bg-[#eee8df] text-[#4a3b2c] border border-[#806344]/80 px-8 sm:px-10 py-3.5 sm:py-4 text-xs uppercase tracking-[0.18em] font-semibold rounded-full shadow-lg hover:bg-[#806344] hover:text-white hover:border-[#806344] hover:scale-[1.03] active:scale-[0.98] transition-all duration-200 inline-flex items-center justify-center cursor-pointer"
          >
            SELL ITEMS
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- Categories section -->
    <section
      class="bg-[#f5f1e9] px-4 py-12 sm:px-6 sm:py-16 md:py-20 lg:px-12 lg:py-24"
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

    <div class="pb-12 sm:pb-16">
      <LoadMoreButton
        @load-more="loadMore"
      />
    </div>

    <MarketplaceAiSearch />
  </main>
</template>
