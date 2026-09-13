<script setup lang="ts">
import type { ProductMedia } from '~/types/product'

const props = defineProps<{
  media?: ProductMedia[]
  fallbackImage: string
}>()

const sortedMedia = computed(() => {
  if (!props.media?.length) return []
  return [...props.media].sort((a, b) => a.position - b.position)
})

const hasMedia = computed(() => sortedMedia.value.length > 0)

// Default to the primary image, or first item
const selectedIndex = ref(0)

const selectedItem = computed(() => {
  if (!hasMedia.value) return null
  return sortedMedia.value[selectedIndex.value] || sortedMedia.value[0] || null
})

function selectItem(index: number) {
  selectedIndex.value = index
}

// Reset selection when media changes
watch(
  () => props.media,
  () => {
    selectedIndex.value = sortedMedia.value.findIndex(m => m.isPrimary)
    if (selectedIndex.value < 0) selectedIndex.value = 0
  },
  { immediate: true }
)
</script>

<template>
  <!-- Multi-media gallery -->
  <div v-if="hasMedia" class="space-y-3">
    <!-- Main view -->
    <div class="overflow-hidden rounded-lg border border-[#d9d0c4] bg-[#eee8df]">
      <!-- Image main view -->
      <img
        v-if="selectedItem?.type === 'image'"
        :src="selectedItem.url"
        :alt="selectedItem.alt || selectedItem.fileName"
        class="aspect-[4/5] sm:aspect-auto sm:min-h-[380px] sm:max-h-[700px] w-full object-cover object-top"
      />

      <!-- Video main view -->
      <video
        v-else-if="selectedItem?.type === 'video'"
        :key="selectedItem.id"
        :src="selectedItem.url"
        controls
        preload="metadata"
        class="aspect-[4/5] sm:aspect-auto sm:min-h-[380px] sm:max-h-[700px] w-full bg-black object-contain"
      />
    </div>

    <!-- Thumbnail strip -->
    <div
      v-if="sortedMedia.length > 1"
      class="flex gap-2 overflow-x-auto pb-1"
    >
      <button
        v-for="(item, index) in sortedMedia"
        :key="item.id"
        type="button"
        :aria-label="`View ${item.type === 'video' ? 'video' : 'image'}: ${item.alt || item.fileName}`"
        class="relative h-10 w-10 sm:h-16 sm:w-16 shrink-0 overflow-hidden rounded-[6px] border-2 transition-all"
        :class="
          selectedIndex === index
            ? 'border-[#806344] ring-1 ring-[#806344]/30'
            : 'border-[#d9d0c4] hover:border-[#b9aa98]'
        "
        @click="selectItem(index)"
      >
        <img
          v-if="item.type === 'image'"
          :src="item.url"
          :alt="item.alt || item.fileName"
          class="h-full w-full object-cover"
        />

        <div
          v-else
          class="flex h-full w-full items-center justify-center bg-[#211f1d]"
        >
          <!-- Play triangle icon -->
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="white"
            class="opacity-80"
          >
            <polygon points="6,4 20,12 6,20" />
          </svg>
        </div>

        <!-- Primary badge (small) -->
        <span
          v-if="item.isPrimary"
          class="absolute bottom-0.5 right-0.5 rounded bg-[#806344] px-1 py-px text-[8px] font-bold uppercase text-white"
        >
          Cover
        </span>
      </button>
    </div>
  </div>

  <!-- Fallback: single image (legacy products without media array) -->
  <div
    v-else
    class="overflow-hidden rounded-lg border border-[#d9d0c4] bg-[#eee8df]"
  >
    <img
      :src="fallbackImage"
      alt="Product image"
      class="aspect-[4/5] sm:aspect-auto sm:min-h-[380px] sm:max-h-[700px] w-full object-cover object-top"
    />
  </div>
</template>
