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
  <section class="bg-[#eee5d5] px-6 py-16 lg:px-12 lg:py-20">
    <div class="mx-auto flex max-w-[1800px] flex-col gap-10 lg:flex-row lg:items-center lg:justify-between lg:gap-16">
      <div class="flex max-w-xl items-start gap-5">
        <div class="mt-1 flex h-10 w-10 shrink-0 items-center justify-center rounded-full border border-[#806344] text-sm font-semibold text-[#806344]">
          AI
        </div>

        <div>
          <h2 class="font-serif text-3xl leading-tight tracking-[-0.02em] text-[#211f1d] sm:text-4xl">
            Looking for something specific?
          </h2>

          <p class="mt-3 max-w-lg text-base leading-7 text-[#806f5e] sm:text-lg">
            Try a natural-language search and browse matching products from the mocked catalog.
          </p>
        </div>
      </div>

      <form
        class="flex w-full max-w-4xl items-center rounded-full border border-[#d7ccbb] bg-[#f8f4ec] p-2"
        @submit.prevent="submitSearch"
      >
        <input
          v-model="query"
          type="text"
          placeholder="What are you looking for?"
          class="min-w-0 flex-1 bg-transparent px-5 py-3 text-base text-[#211f1d] outline-none placeholder:text-[#b9aa98] sm:px-7 sm:py-4 sm:text-lg"
        />

        <button
          type="submit"
          aria-label="Search with AI"
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full border border-[#806344] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white hover:scale-105 focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 sm:h-14 sm:w-14"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="22"
            height="22"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path d="M5 12h14" />
            <path d="m13 6 6 6-6 6" />
          </svg>
        </button>
      </form>
    </div>
  </section>
</template>
