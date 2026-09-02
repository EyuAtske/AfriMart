<script setup lang="ts">
const router = useRouter()
const prompt = ref('')
const suggestions = ref<string[]>([])
const isThinking = ref(false)

const submitPrompt = async () => {
  const value = prompt.value.trim()
  if (!value) return

  isThinking.value = true

  await new Promise(resolve => setTimeout(resolve, 250))

  suggestions.value = [
    `Search for "${value}"`,
    'Compare prices under 3,000 ETB',
    'Check seller stock before checkout'
  ]

  isThinking.value = false
}

const searchSuggestion = async () => {
  const value = prompt.value.trim()
  if (!value) return

  await router.push({
    path: '/products',
    query: { search: value }
  })
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <section class="mx-auto max-w-4xl rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
      <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
        AI shopping assistant
      </p>

      <h1 class="mt-2 font-serif text-4xl tracking-[-0.025em] text-[#211f1d] sm:text-5xl">
        Find the right item faster.
      </h1>

      <form
        class="mt-8 flex flex-col gap-3 sm:flex-row"
        @submit.prevent="submitPrompt"
      >
        <input
          v-model="prompt"
          type="text"
          placeholder="Example: a relaxed denim outfit"
          class="h-12 min-w-0 flex-1 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-5 text-sm text-[#211f1d] outline-none transition focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
        />

        <button
          type="submit"
          :disabled="isThinking"
          class="h-12 rounded-full bg-[#211f1d] px-6 text-sm font-medium uppercase tracking-[0.14em] text-white transition hover:bg-[#3b3733] disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ isThinking ? 'Thinking...' : 'Ask' }}
        </button>
      </form>

      <div
        v-if="suggestions.length"
        class="mt-8 grid gap-3"
      >
        <button
          v-for="suggestion in suggestions"
          :key="suggestion"
          type="button"
          class="rounded-[8px] border border-[#ded6cc] bg-[#f5f1e9] px-4 py-4 text-left text-sm text-[#211f1d] transition hover:border-[#806344]"
          @click="searchSuggestion"
        >
          {{ suggestion }}
        </button>
      </div>
    </section>
  </main>
</template>
