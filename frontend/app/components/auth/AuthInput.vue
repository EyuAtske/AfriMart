<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue: string
    label: string
    placeholder?: string
    type?: string
    autocomplete?: string
    name?: string
  }>(),
  {
    placeholder: '',
    type: 'text',
    autocomplete: 'off',
    name: ''
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const showPassword = ref(false)

const inputType = computed(() =>
  props.type === 'password' && showPassword.value
    ? 'text'
    : props.type
)
</script>

<template>
  <div class="space-y-2">
    <label
      :for="name || label"
      class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]"
    >
      {{ label }}
    </label>

    <div class="relative">
      <input
        :id="name || label"
        :name="name"
        :type="inputType"
        :value="modelValue"
        :placeholder="placeholder"
        :autocomplete="autocomplete"
        class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
        @input="
          emit(
            'update:modelValue',
            ($event.target as HTMLInputElement).value
          )
        "
      />

      <button
        v-if="type === 'password'"
        type="button"
        :aria-label="showPassword ? 'Hide password' : 'Show password'"
        class="absolute right-3 top-1/2 -translate-y-1/2 rounded p-1 text-[#74685d] transition hover:text-[#211f1d] focus:outline-none focus:ring-2 focus:ring-[#806344]/30"
        @click="showPassword = !showPassword"
      >
        <svg
          v-if="!showPassword"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          class="h-4 w-4"
        >
          <path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6S2 12 2 12Z" />
          <circle cx="12" cy="12" r="2.5" />
        </svg>

        <svg
          v-else
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          class="h-4 w-4"
        >
          <path d="M3 3l18 18" />
          <path d="M10.6 6.2A9.8 9.8 0 0 1 12 6c6.5 0 10 6 10 6a17 17 0 0 1-3.1 3.7" />
          <path d="M6.2 6.2C3.5 8.1 2 12 2 12s3.5 6 10 6c1.8 0 3.3-.4 4.6-1" />
        </svg>
      </button>
    </div>
  </div>
</template>