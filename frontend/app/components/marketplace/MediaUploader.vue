<script setup lang="ts">
import type { ProductMedia } from '~/types/product'
import { ALL_ACCEPTED_TYPES, MAX_FILES, formatFileSize } from '~/repositories/mock/MockMediaRepository'

const props = withDefaults(
  defineProps<{
    modelValue?: ProductMedia[]
    existing?: ProductMedia[]
  }>(),
  {
    modelValue: () => [],
    existing: () => []
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: ProductMedia[]]
}>()

const { mediaRepo } = useRepositories()



const mediaItems = ref<ProductMedia[]>([...props.existing])
const errors = ref<string[]>([])
const isDraggingOver = ref(false)
const isUploading = ref(false)

onBeforeUnmount(() => {
  if (import.meta.client) {
    for (const item of mediaItems.value) {
      if (item.url?.startsWith('blob:')) {
        URL.revokeObjectURL(item.url)
      }
    }
  }
})

// Sync existing prop on mount (edit mode)
watch(
  () => props.existing,
  (val) => {
    if (val.length && !mediaItems.value.length) {
      mediaItems.value = [...val]
      emitUpdate()
    }
  },
  { immediate: true }
)

function emitUpdate() {
  // Recalculate positions before emitting
  mediaItems.value.forEach((item, idx) => {
    item.position = idx
  })
  emit('update:modelValue', [...mediaItems.value])
}

async function handleFiles(files: FileList | File[]) {
  const fileArray = Array.from(files)
  if (!fileArray.length) return

  // Check total count before uploading
  const remainingSlots = MAX_FILES - mediaItems.value.length
  if (remainingSlots <= 0) {
    errors.value = [`Maximum ${MAX_FILES} files already added.`]
    return
  }

  const filesToUpload = fileArray.slice(0, remainingSlots)
  if (fileArray.length > remainingSlots) {
    errors.value = [`Only ${remainingSlots} more file(s) can be added (max ${MAX_FILES}).`]
  } else {
    errors.value = []
  }

  isUploading.value = true

  const result = await mediaRepo.uploadMedia(filesToUpload)

  isUploading.value = false

  if (result.errors.length) {
    errors.value = [...errors.value, ...result.errors]
  }

  if (result.uploaded.length) {
    // If no primary exists yet, set the first image as primary
    const hasPrimary = mediaItems.value.some(m => m.isPrimary) ||
      result.uploaded.some(m => m.isPrimary)

    if (!hasPrimary) {
      const firstImage = result.uploaded.find(m => m.type === 'image')
      if (firstImage) {
        firstImage.isPrimary = true
      }
    } else {
      // Don't override existing primary
      result.uploaded.forEach(m => { m.isPrimary = false })
    }

    mediaItems.value = [...mediaItems.value, ...result.uploaded]
    emitUpdate()
  }
}

function onFileInput(event: Event) {
  const input = event.target as HTMLInputElement
  if (input.files) {
    handleFiles(input.files)
    input.value = '' // Reset so same file can be re-selected
  }
}

function onDrop(event: DragEvent) {
  isDraggingOver.value = false
  if (event.dataTransfer?.files) {
    handleFiles(event.dataTransfer.files)
  }
}

function onDragOver(event: DragEvent) {
  event.preventDefault()
  isDraggingOver.value = true
}

function onDragLeave() {
  isDraggingOver.value = false
}

function removeItem(index: number) {
  const item = mediaItems.value[index]
  if (!item) return

  // Revoke object URL on client
  if (import.meta.client && item.url.startsWith('blob:')) {
    URL.revokeObjectURL(item.url)
  }

  const wasPrimary = item.isPrimary
  mediaItems.value.splice(index, 1)

  // If removed item was primary, assign primary to first remaining image
  if (wasPrimary) {
    const firstImage = mediaItems.value.find(m => m.type === 'image')
    if (firstImage) {
      firstImage.isPrimary = true
    }
  }

  emitUpdate()
}

function moveItem(index: number, direction: -1 | 1) {
  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= mediaItems.value.length) return

  const items = [...mediaItems.value]
  const temp = items[index]!
  items[index] = items[newIndex]!
  items[newIndex] = temp
  mediaItems.value = items
  emitUpdate()
}

function setPrimary(index: number) {
  const item = mediaItems.value[index]
  if (!item || item.type !== 'image') return

  mediaItems.value.forEach(m => { m.isPrimary = false })
  item.isPrimary = true
  emitUpdate()
}

function updateAlt(index: number, alt: string) {
  const item = mediaItems.value[index]
  if (item) {
    item.alt = alt
    emitUpdate()
  }
}

const acceptString = ALL_ACCEPTED_TYPES.join(',')
</script>

<template>
  <div class="space-y-4">
    <!-- Drop zone -->
    <label
      class="flex min-h-44 cursor-pointer flex-col items-center justify-center rounded-[10px] border-2 border-dashed px-6 py-8 text-center transition-colors"
      :class="[
        isDraggingOver
          ? 'border-[#806344] bg-[#806344]/5'
          : 'border-[#b9aa98] bg-[#f5f1e9] hover:border-[#806344]',
        isUploading ? 'pointer-events-none opacity-60' : ''
      ]"
      @drop.prevent="onDrop"
      @dragover="onDragOver"
      @dragleave="onDragLeave"
    >
      <!-- Upload icon -->
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="32"
        height="32"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.4"
        class="mb-3 text-[#806344]"
      >
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
        <polyline points="17 8 12 3 7 8" />
        <line x1="12" y1="3" x2="12" y2="15" />
      </svg>

      <span v-if="isUploading" class="text-sm font-medium text-[#806344]">
        Uploading…
      </span>
      <template v-else>
        <span class="block font-serif text-lg text-[#211f1d]">
          Drag files here or click to browse
        </span>
        <span class="mt-1.5 block text-xs leading-5 text-[#756a60]">
          Images: JPEG, PNG, WebP (max 5 MB) · Up to {{ MAX_FILES }} files
        </span>
        <span v-if="mediaItems.length" class="mt-1 block text-xs font-semibold text-[#806344]">
          {{ mediaItems.length }} of {{ MAX_FILES }} images selected
        </span>
      </template>

      <input
        type="file"
        :accept="acceptString"
        multiple
        class="sr-only"
        aria-label="Upload product media files"
        @change="onFileInput"
      />
    </label>

    <!-- Upload progress bar -->
    <div
      v-if="isUploading"
      class="h-1 overflow-hidden rounded-full bg-[#e8dcc9]"
    >
      <div class="h-full w-2/3 animate-pulse rounded-full bg-[#806344]" />
    </div>

    <!-- Errors -->
    <div
      v-if="errors.length"
      class="rounded-[8px] border border-red-200 bg-red-50 px-4 py-3"
    >
      <p
        v-for="(err, i) in errors"
        :key="i"
        class="text-sm text-red-700"
      >
        {{ err }}
      </p>
    </div>

    <!-- Media grid -->
    <div
      v-if="mediaItems.length"
      class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4"
    >
      <div
        v-for="(item, index) in mediaItems"
        :key="item.id"
        class="group relative overflow-hidden rounded-[8px] border bg-white"
        :class="item.isPrimary ? 'border-[#806344] ring-2 ring-[#806344]/30' : 'border-[#d9d0c4]'"
      >
        <!-- Preview -->
        <div class="relative aspect-square overflow-hidden bg-[#f5f1e9]">
          <img
            v-if="item.type === 'image'"
            :src="item.url"
            :alt="item.alt || item.fileName"
            class="h-full w-full object-cover"
          />
          <div v-else class="relative h-full w-full">
            <video
              :src="item.url"
              controls
              preload="metadata"
              class="h-full w-full object-cover"
            />
            <!-- Video type badge -->
            <span class="absolute left-2 top-2 rounded bg-black/70 px-1.5 py-0.5 text-[10px] font-semibold uppercase text-white">
              Video
            </span>
          </div>

          <!-- Primary badge -->
          <span
            v-if="item.isPrimary"
            class="absolute right-2 top-2 rounded-full bg-[#806344] px-2.5 py-1 text-[10px] font-semibold uppercase tracking-wider text-white shadow"
          >
            Cover
          </span>

          <!-- Remove button -->
          <button
            type="button"
            :aria-label="`Remove ${item.fileName}`"
            class="absolute left-2 top-2 flex h-6 w-6 items-center justify-center rounded-full bg-black/60 text-xs text-white opacity-0 transition-opacity group-hover:opacity-100 hover:bg-red-600"
            :class="{ 'left-auto right-2': item.isPrimary, 'top-8': item.isPrimary }"
            @click="removeItem(index)"
          >
            ✕
          </button>
        </div>

        <!-- Info & controls -->
        <div class="px-2.5 py-2">
          <p class="truncate text-xs font-medium text-[#211f1d]">
            {{ item.fileName }}
          </p>
          <p class="text-[10px] text-[#756a60]">
            {{ formatFileSize(item.fileSize) }}
          </p>

          <!-- Alt text input (images only) -->
          <input
            v-if="item.type === 'image'"
            :value="item.alt"
            type="text"
            placeholder="Alt text…"
            :aria-label="`Alt text for ${item.fileName}`"
            class="mt-1.5 h-7 w-full rounded border border-[#d9d0c4] bg-[#faf8f4] px-2 text-[11px] text-[#211f1d] outline-none transition placeholder:text-[#b9aa98] focus:border-[#806344]"
            @input="updateAlt(index, ($event.target as HTMLInputElement).value)"
          />

          <!-- Action buttons -->
          <div class="mt-2 flex items-center gap-1">
            <button
              type="button"
              :disabled="index === 0"
              :aria-label="`Move ${item.fileName} left`"
              class="flex h-6 w-6 items-center justify-center rounded bg-[#f5f1e9] text-[10px] text-[#756a60] transition hover:bg-[#ded6cc] disabled:opacity-30 disabled:cursor-not-allowed"
              @click="moveItem(index, -1)"
            >
              ←
            </button>

            <button
              type="button"
              :disabled="index === mediaItems.length - 1"
              :aria-label="`Move ${item.fileName} right`"
              class="flex h-6 w-6 items-center justify-center rounded bg-[#f5f1e9] text-[10px] text-[#756a60] transition hover:bg-[#ded6cc] disabled:opacity-30 disabled:cursor-not-allowed"
              @click="moveItem(index, 1)"
            >
              →
            </button>

            <button
              v-if="item.type === 'image' && !item.isPrimary"
              type="button"
              :aria-label="`Set ${item.fileName} as cover image`"
              class="ml-auto rounded bg-[#f5f1e9] px-2 py-1 text-[10px] font-medium text-[#806344] transition hover:bg-[#806344] hover:text-white"
              @click="setPrimary(index)"
            >
              Set cover
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty state (when no items and not uploading) -->
    <p
      v-if="!mediaItems.length && !isUploading"
      class="text-center text-sm text-[#756a60]"
    >
      No media files selected yet.
    </p>
  </div>
</template>
