<script setup lang="ts">
const { toasts, removeToast } = useToast()
</script>

<template>
  <div class="fixed bottom-6 right-6 z-50 flex flex-col gap-3 max-w-sm w-full pointer-events-none">
    <TransitionGroup
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 translate-y-4 scale-95"
      enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 translate-y-2 scale-95"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="pointer-events-auto flex items-center justify-between rounded-[8px] border border-[#d9d0c4] bg-[#211f1d] px-4 py-3 text-white shadow-xl"
      >
        <div class="flex items-center gap-3">
          <span
            v-if="toast.type === 'success'"
            class="flex h-6 w-6 items-center justify-center rounded-full bg-[#806344] text-xs font-bold text-white"
          >
            ✓
          </span>
          <span
            v-else-if="toast.type === 'error'"
            class="flex h-6 w-6 items-center justify-center rounded-full bg-red-600 text-xs font-bold text-white"
          >
            ✕
          </span>
          <span
            v-else
            class="flex h-6 w-6 items-center justify-center rounded-full bg-[#806344]/40 text-xs font-bold text-white"
          >
            i
          </span>

          <p class="text-sm font-medium leading-tight text-[#f5f1e9]">
            {{ toast.message }}
          </p>
        </div>

        <button
          type="button"
          class="ml-4 text-xs text-[#b0a79d] hover:text-white"
          @click="removeToast(toast.id)"
        >
          ✕
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
