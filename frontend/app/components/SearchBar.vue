<script setup lang="ts">
const router = useRouter()

const searchQuery = ref('')
const isFocused = ref(false)

const recentSearches = ref<string[]>([])

const MAX_RECENT_SEARCHES = 5

// Load recent searches from localStorage
onMounted(() => {
  const savedSearches = localStorage.getItem('afrimart-recent-searches')

  if (savedSearches) {
    try {
      recentSearches.value = JSON.parse(savedSearches)
    } catch {
      recentSearches.value = []
    }
  }
})

// Save recent searches
const saveRecentSearch = (query: string) => {
  const cleanedQuery = query.trim()

  if (!cleanedQuery) return

  // Remove duplicate
  const filtered = recentSearches.value.filter(
    search => search.toLowerCase() !== cleanedQuery.toLowerCase()
  )

  // Put newest search first
  recentSearches.value = [
    cleanedQuery,
    ...filtered
  ].slice(0, MAX_RECENT_SEARCHES)

  localStorage.setItem(
    'afrimart-recent-searches',
    JSON.stringify(recentSearches.value)
  )
}

// Perform normal product search
const performSearch = (query = searchQuery.value) => {
  const cleanedQuery = query.trim()

  if (!cleanedQuery) return

  saveRecentSearch(cleanedQuery)

  searchQuery.value = cleanedQuery
  isFocused.value = false

  router.push({
    path: '/products',
    query: {
      search: cleanedQuery
    }
  })
}

// Click recent search
const selectRecentSearch = (search: string) => {
  searchQuery.value = search

  performSearch(search)
}

// Remove one recent search
const removeRecentSearch = (search: string) => {
  recentSearches.value = recentSearches.value.filter(
    item => item !== search
  )

  localStorage.setItem(
    'afrimart-recent-searches',
    JSON.stringify(recentSearches.value)
  )
}

// Clear all recent searches
const clearRecentSearches = () => {
  recentSearches.value = []

  localStorage.removeItem('afrimart-recent-searches')
}

// Future AI search
const performAiSearch = () => {
  const cleanedQuery = searchQuery.value.trim()

  if (!cleanedQuery) return

  saveRecentSearch(cleanedQuery)

  // AI SERVICE WILL BE CONNECTED HERE LATER
  console.log('AI search:', cleanedQuery)
}

// Decide whether the query looks like a normal search
// or a natural-language AI request.
const isNaturalLanguageQuery = computed(() => {
  const query = searchQuery.value.trim()

  if (!query) return false

  const aiIndicators = [
    'i want',
    'i need',
    'find me',
    'looking for',
    'show me',
    'something',
    'outfit',
    'occasion',
    'under',
    'around',
    'cheap',
    'affordable',
    'recommend',
    'recommendation'
  ]

  const lowerQuery = query.toLowerCase()

  return aiIndicators.some(indicator =>
    lowerQuery.includes(indicator)
  )
})
</script>

<template>
  <div class="relative w-[min(70vw,520px)]">
    <!-- Search input -->
    <div
      class="flex h-10 items-center border-b border-[#806344] bg-[#f5f1e9]"
    >
      <!-- Search icon -->
      <button
        type="button"
        aria-label="Search"
        class="mr-3 shrink-0 text-[#806344]"
        @click="performSearch()"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="19"
          height="19"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.6"
        >
          <circle cx="11" cy="11" r="7" />
          <path d="m20 20-4-4" />
        </svg>
      </button>

      <!-- Input -->
      <input
        v-model="searchQuery"
        type="text"
        autocomplete="off"
        :placeholder="
          isNaturalLanguageQuery
            ? 'Ask AfriMart...'
            : 'Search products...'
        "
        class="min-w-0 flex-1 bg-transparent text-sm text-[#302d29] outline-none placeholder:text-[#9b9388]"
        @focus="isFocused = true"
        @keydown.enter="performSearch()"
      />

      <!-- AI indicator -->
      <span
        v-if="isNaturalLanguageQuery"
        class="mr-2 hidden text-[9px] tracking-[0.12em] text-[#806344] sm:block"
      >
        AI
      </span>

      <!-- Clear input -->
      <button
        v-if="searchQuery"
        type="button"
        aria-label="Clear search"
        class="ml-2 shrink-0 text-[#806344] transition-opacity hover:opacity-60"
        @click="searchQuery = ''"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="17"
          height="17"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.6"
        >
          <path d="M6 6l12 12" />
          <path d="M18 6L6 18" />
        </svg>
      </button>
    </div>

    <!-- Search dropdown -->
    <div
      v-if="isFocused"
      class="absolute left-0 right-0 top-full z-50 mt-2 border border-[#ded8ce] bg-[#f5f1e9] shadow-[0_12px_30px_rgba(48,45,41,0.10)]"
    >
      <!-- Search suggestions -->
      <div class="p-4 sm:p-5">

        <!-- AI suggestion -->
        <button
          v-if="searchQuery && isNaturalLanguageQuery"
          type="button"
          class="mb-5 flex w-full items-start gap-3 border-b border-[#ded8ce] pb-4 text-left"
          @click="performAiSearch"
        >
          <span
            class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-[#806344] text-white"
          >
            ★
          </span>

          <span>
            <span
              class="block text-xs tracking-[0.08em] text-[#302d29]"
            >
              ASK AFRIMART
            </span>

            <span
              class="mt-1 block text-xs leading-5 text-[#806344]"
            >
              Find products based on your description
            </span>
          </span>
        </button>

        <!-- Recent searches -->
        <div v-if="recentSearches.length">
          <div class="mb-3 flex items-center justify-between">
            <p
              class="text-[10px] font-medium tracking-[0.14em] text-[#806344]"
            >
              RECENT SEARCHES
            </p>

            <button
              type="button"
              class="text-[10px] tracking-[0.08em] text-[#9b9388] transition-colors hover:text-[#806344]"
              @click="clearRecentSearches"
            >
              CLEAR
            </button>
          </div>

          <div class="flex flex-col">
            <div
              v-for="search in recentSearches"
              :key="search"
              class="group flex items-center"
            >
              <!-- Search -->
              <button
                type="button"
                class="flex min-w-0 flex-1 items-center gap-3 py-2 text-left text-sm text-[#4c4945] transition-colors hover:text-[#806344]"
                @click="selectRecentSearch(search)"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="15"
                  height="15"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.5"
                  class="shrink-0 text-[#9b9388]"
                >
                  <circle cx="12" cy="12" r="8" />
                  <path d="M12 8v4l2.5 2" />
                </svg>

                <span class="truncate">
                  {{ search }}
                </span>
              </button>

              <!-- Remove -->
              <button
                type="button"
                aria-label="Remove recent search"
                class="ml-2 hidden shrink-0 text-[#b0a89e] hover:text-[#806344] group-hover:block"
                @click.stop="removeRecentSearch(search)"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <path d="M6 6l12 12" />
                  <path d="M18 6L6 18" />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div
          v-else-if="!searchQuery"
          class="py-3"
        >
          <p
            class="text-[10px] tracking-[0.14em] text-[#806344]"
          >
            SEARCH AFRIMART
          </p>

          <p
            class="mt-2 text-xs leading-5 text-[#8d857b]"
          >
            Search for products or describe what you're looking for.
          </p>
        </div>

        <!-- Normal search hint -->
        <div
          v-if="searchQuery && !isNaturalLanguageQuery"
          class="pt-2"
        >
         <button
          type="button"
            class="flex w-full items-center gap-3 text-left text-sm text-[#4c4945] transition-colors hover:text-[#806344]"
            @click="performSearch(searchQuery)"
>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <circle cx="11" cy="11" r="7" />
              <path d="m20 20-4-4" />
            </svg>

            <span>
              Search for
              <strong class="font-medium text-[#302d29]">
                "{{ searchQuery }}"
              </strong>
            </span>
          </button>
        </div>

      </div>
    </div>
  </div>
</template>