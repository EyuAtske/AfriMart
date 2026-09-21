<script setup lang="ts">
import type { Product } from '~/types/product'

const router = useRouter()
const route = useRoute()
const { filterProducts } = useMarketplace()

const prompt = ref('')
const suggestions = ref<string[]>([])
const matchingProducts = ref<Product[]>([])
const isThinking = ref(false)

const submitPrompt = async () => {
  const value = prompt.value.trim()
  if (!value) return

  isThinking.value = true

  await new Promise(resolve => setTimeout(resolve, 250))

  // Find products matching keywords in prompt
  const words = value.toLowerCase().split(/\s+/).filter(w => w.length > 2)
  const results = filterProducts({ search: value })
  
  if (results.length) {
    matchingProducts.value = results
  } else if (words.length) {
    // Try matching individual key words
    const partials = filterProducts({ search: words[0] })
    matchingProducts.value = partials
  } else {
    matchingProducts.value = []
  }

  suggestions.value = [
    `Browse all results for "${value}"`,
    'Compare prices under 3,000 ETB',
    'Check seller stock before checkout'
  ]

  isThinking.value = false
}

const searchSuggestion = async (suggestion: string) => {
  const value = prompt.value.trim()
  await router.push({
    path: '/products',
    query: { search: value || undefined }
  })
}

onMounted(() => {
  if (route.query.prompt) {
    prompt.value = String(route.query.prompt)
    submitPrompt()
  }
})

watch(() => route.query.prompt, (newPrompt) => {
  if (newPrompt && String(newPrompt) !== prompt.value) {
    prompt.value = String(newPrompt)
    submitPrompt()
  }
})
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <UiAppCard
      as="section"
      padding="large"
      class="mx-auto max-w-4xl"
    >
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
          class="h-12 min-w-0 flex-1 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-5 text-sm text-[#211f1d] outline-none transition placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
        />

        <UiAppButton
          type="submit"
          variant="primary"
          :disabled="isThinking"
        >
          {{ isThinking ? 'Thinking...' : 'Ask' }}
        </UiAppButton>
      </form>

      <div
        v-if="matchingProducts.length"
        class="mt-10 border-t border-[#ded6cc] pt-8"
      >
        <h3 class="font-serif text-2xl text-[#211f1d] mb-6">
          Suggested Products ({{ matchingProducts.length }})
        </h3>
        <MarketplaceProductGrid :products="matchingProducts" />
      </div>
    </UiAppCard>
  </main>
</template>
