<script setup lang="ts">
import AuthInput from '~/components/auth/AuthInput.vue'

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

const buttonClass = 'h-12 rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none'

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
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-xs"
      @click.self="emit('close')"
    >
      <div class="w-full max-w-lg overflow-hidden rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-2xl sm:p-8">
        <div class="flex items-center justify-between border-b border-[#ded6cc] pb-4">
          <h3 class="font-serif text-2xl text-[#211f1d]">
            Edit Shop Details
          </h3>

          <button
            type="button"
            class="text-[#756a60] hover:text-[#211f1d]"
            @click="emit('close')"
          >
            ✕
          </button>
        </div>

        <form class="mt-6 space-y-5" @submit.prevent="submit">
          <AuthInput
            v-model="form.name"
            label="Shop name"
            placeholder="Shop name"
            name="edit-shop-name"
          />

          <div class="space-y-2">
            <label
              for="edit-shop-description"
              class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]"
            >
              Description
            </label>

            <textarea
              id="edit-shop-description"
              v-model="form.description"
              rows="4"
              placeholder="Describe your shop"
              class="w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
            />
          </div>

          <p v-if="error" class="text-sm font-medium text-red-600">
            {{ error }}
          </p>

          <div class="flex justify-end gap-3 pt-2">
            <button
              type="button"
              class="h-12 rounded-full border border-[#cfc4b5] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#756a60] transition hover:bg-[#eee8df]"
              @click="emit('close')"
            >
              Cancel
            </button>

            <button
              type="submit"
              :class="buttonClass"
            >
              Save Changes
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>
