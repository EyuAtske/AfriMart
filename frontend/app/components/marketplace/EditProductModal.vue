<script setup lang="ts">
import type { SellerProduct } from '~/composables/useSellerShop'
import type { ProductCategory } from '~/types/product'
import AuthInput from '~/components/auth/AuthInput.vue'

const props = defineProps<{
  isOpen: boolean
  product: SellerProduct | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const { updateSellerProduct } = useSellerShop()
const { showToast } = useToast()

const categories: ProductCategory[] = ['Men', 'Women', 'Kids', 'Shoes', 'Accessories']

const form = reactive({
  name: '',
  description: '',
  category: (categories[0] || 'Men') as ProductCategory,
  price: 1000,
  stock: 1,
  image: '',
  status: 'Active' as 'Active' | 'Draft'
})

const error = ref('')

watch(
  () => props.product,
  (newProduct) => {
    if (newProduct) {
      form.name = newProduct.name
      form.description = newProduct.description
      form.category = newProduct.category
      form.price = newProduct.price
      form.stock = newProduct.stock
      form.image = newProduct.image
      form.status = newProduct.status
      error.value = ''
    }
  },
  { immediate: true }
)

const handleImageUpload = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) {
    form.image = URL.createObjectURL(file)
  }
}

const submit = () => {
  if (!props.product) return
  error.value = ''

  if (
    !form.name.trim() ||
    !form.description.trim() ||
    form.price < 1 ||
    form.stock < 0 ||
    !form.image
  ) {
    error.value = 'Please provide valid values for all required fields.'
    return
  }

  updateSellerProduct(props.product.id, {
    name: form.name.trim(),
    description: form.description.trim(),
    category: form.category as ProductCategory,
    price: form.price,
    stock: form.stock,
    image: form.image,
    status: form.status
  })

  showToast('Product updated successfully!')
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isOpen && product"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-xs"
      @click.self="emit('close')"
    >
      <div class="w-full max-w-2xl max-h-[90vh] overflow-y-auto rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-2xl sm:p-8">
        <div class="flex items-center justify-between border-b border-[#ded6cc] pb-4">
          <h3 class="font-serif text-2xl text-[#211f1d]">
            Edit Product
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
          <div class="grid gap-6 sm:grid-cols-[160px_1fr]">
            <label class="flex flex-col items-center justify-center overflow-hidden rounded-[8px] border border-dashed border-[#b9aa98] bg-[#f5f1e9] h-36 cursor-pointer hover:border-[#806344]">
              <img
                v-if="form.image"
                :src="form.image"
                alt="Product preview"
                class="h-full w-full object-cover"
              />
              <span v-else class="text-xs text-[#756a60] text-center px-2">Change Image</span>
              <input type="file" accept="image/*" class="sr-only" @change="handleImageUpload" />
            </label>

            <div class="space-y-4">
              <AuthInput
                v-model="form.name"
                label="Product name"
                placeholder="Product name"
                name="edit-product-name"
              />

              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-2">
                  <label class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                    Category
                  </label>
                  <select
                    v-model="form.category"
                    class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none hover:border-[#9e8b77] focus:border-[#806344]"
                  >
                    <option v-for="cat in categories" :key="cat" :value="cat">
                      {{ cat }}
                    </option>
                  </select>
                </div>

                <div class="space-y-2">
                  <label class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                    Status
                  </label>
                  <select
                    v-model="form.status"
                    class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none hover:border-[#9e8b77] focus:border-[#806344]"
                  >
                    <option value="Active">Active</option>
                    <option value="Draft">Draft</option>
                  </select>
                </div>
              </div>
            </div>
          </div>

          <div class="space-y-2">
            <label class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
              Description
            </label>
            <textarea
              v-model="form.description"
              rows="3"
              class="w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none hover:border-[#9e8b77] focus:border-[#806344]"
            />
          </div>

          <div class="grid gap-5 sm:grid-cols-2">
            <label class="space-y-2">
              <span class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                Price (ETB)
              </span>
              <input
                v-model.number="form.price"
                type="number"
                min="1"
                step="50"
                class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none hover:border-[#9e8b77] focus:border-[#806344]"
              />
            </label>

            <label class="space-y-2">
              <span class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                Stock Count
              </span>
              <input
                v-model.number="form.stock"
                type="number"
                min="0"
                step="1"
                class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none hover:border-[#9e8b77] focus:border-[#806344]"
              />
            </label>
          </div>

          <p v-if="error" class="text-sm font-medium text-red-600">
            {{ error }}
          </p>

          <div class="flex justify-end gap-3 pt-4 border-t border-[#ded6cc]">
            <button
              type="button"
              class="h-12 rounded-full border border-[#cfc4b5] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#756a60] transition hover:bg-[#eee8df]"
              @click="emit('close')"
            >
              Cancel
            </button>

            <button
              type="submit"
              class="h-12 rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
            >
              Update Product
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>
