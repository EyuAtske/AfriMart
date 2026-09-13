<script setup lang="ts">
withDefaults(
  defineProps<{
    isOpen: boolean
    title: string
    kicker?: string
    maxWidth?: 'sm' | 'md' | 'lg'
  }>(),
  {
    maxWidth: 'md'
  }
)

const emit = defineEmits<{
  close: []
}>()
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isOpen"
      role="dialog"
      aria-modal="true"
      :aria-label="title"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-xs"
      @click.self="emit('close')"
    >
      <div
        class="w-full max-h-[90vh] overflow-y-auto rounded-xl border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-xl sm:p-8"
        :class="{
          'max-w-sm': maxWidth === 'sm',
          'max-w-lg': maxWidth === 'md',
          'max-w-2xl': maxWidth === 'lg'
        }"
      >
        <div class="flex items-center justify-between border-b border-[#ded6cc] pb-4">
          <div>
            <p v-if="kicker" class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
              {{ kicker }}
            </p>

            <h3
              class="font-serif text-2xl text-[#211f1d]"
              :class="kicker ? 'mt-1' : ''"
            >
              {{ title }}
            </h3>
          </div>

          <button
            type="button"
            aria-label="Close modal"
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-[#eee8df] text-[#756a60] transition hover:bg-[#ded6cc] hover:text-[#211f1d]"
            @click="emit('close')"
          >
            ✕
          </button>
        </div>

        <!-- Optional header slot for subtitles/descriptions below the title bar -->
        <slot name="header" />

        <div class="mt-6">
          <slot />
        </div>
      </div>
    </div>
  </Teleport>
</template>
