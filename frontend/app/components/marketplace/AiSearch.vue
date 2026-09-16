<script setup lang="ts">
const router = useRouter()
const query = ref('')

const submitSearch = async () => {
  const search = query.value.trim()
  if (!search) return

  await router.push({
    path: '/products',
    query: { search }
  })
}
</script>

<template>
  <section class="bg-[#eee8df] px-5 py-12 sm:px-8 sm:py-16 lg:px-12 lg:py-20">
    <div class="mx-auto flex w-full max-w-[1800px] flex-col gap-6 sm:gap-8 lg:flex-row lg:items-center lg:justify-between lg:gap-16">
      
      <!-- Text & Badge -->
      <div class="flex items-start gap-4 sm:gap-5 max-w-2xl">
        <div class="mt-1 flex h-11 w-11 sm:h-13 sm:w-13 shrink-0 items-center justify-center rounded-full border-2 border-[#806344] bg-[#806344]/10 text-sm sm:text-base font-bold text-[#806344]">
          AI
        </div>

        <div>
          <h2 class="font-serif text-2xl sm:text-3xl lg:text-4xl leading-tight tracking-[-0.02em] text-[#211f1d]">
            Looking for something specific?
          </h2>

          <p class="mt-2 text-sm sm:text-base leading-relaxed text-[#756a60]">
            Try a natural-language search and browse matching products from our catalog.
          </p>
        </div>
      </div>

      <!-- Search Form -->
      <form
        class="flex w-full lg:max-w-xl xl:max-w-2xl items-center rounded-full border-2 border-[#d9d0c4] bg-[#faf8f4] p-2 sm:p-2.5 shadow-sm transition-all focus-within:border-[#806344] focus-within:shadow-[0_4px_20px_rgba(128,99,68,0.15)]"
        @submit.prevent="submitSearch"
      >
        <input
          v-model="query"
          type="text"
          placeholder="e.g. casual linen shirt for warm weather"
          class="min-w-0 flex-1 bg-transparent px-4 py-3 sm:px-6 sm:py-3.5 text-sm sm:text-base text-[#211f1d] outline-none placeholder:text-[#92877b]"
        />

        <button
          type="submit"
          aria-label="Search with AI"
          class="flex h-11 w-11 sm:h-13 sm:w-13 shrink-0 items-center justify-center rounded-full bg-[#806344] text-white transition-all duration-300 hover:bg-[#684f34] hover:scale-105 focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            class="h-5 w-5 sm:h-6 sm:w-6"
          >
            <path d="M5 12h14" />
            <path d="m13 6 6 6-6 6" />
          </svg>
        </button>
      </form>
    </div>
  </section>
</template>
