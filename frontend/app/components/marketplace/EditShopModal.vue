<script setup lang="ts">
const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const { shop, updateShop } = useSellerShop()
const { showToast } = useToast()

const form = reactive({
  name: '',
  description: ''
})

const error = ref('')

watch(() => props.isOpen, (open) => {
  if (open && shop.value) {
    form.name = shop.value.name
    form.description = shop.value.description
    error.value = ''
  }
})

const submit = () => {
  error.value = ''

  if (!form.name.trim() || !form.description.trim()) {
    error.value = 'Please provide both shop name and description.'
    return
  }

  updateShop({
    name: form.name,
    description: form.description
  })

  showToast('Shop details updated successfully!')
  emit('close')
}
</script>

<template>
  <UiAppModal
    :is-open="isOpen"
    title="Edit Shop Details"
    @close="emit('close')"
  >
    <form class="space-y-5" @submit.prevent="submit">
      <AuthInput
        v-model="form.name"
        label="Shop name"
        placeholder="Shop name"
        name="edit-shop-name"
      />

      <div class="space-y-2">
        <label
          for="edit-shop-description"
          class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]"
        >
          Description
        </label>

        <textarea
          id="edit-shop-description"
          v-model="form.description"
          rows="4"
          placeholder="Describe your shop"
          class="w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
        />
      </div>

      <UiAppAlert v-if="error">
        {{ error }}
      </UiAppAlert>

      <div class="flex justify-end gap-3 pt-2">
        <UiAppButton
          variant="ghost"
          @click="emit('close')"
        >
          Cancel
        </UiAppButton>

        <UiAppButton type="submit">
          Save Changes
        </UiAppButton>
      </div>
    </form>
  </UiAppModal>
</template>
