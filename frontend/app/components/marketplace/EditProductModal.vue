<script setup lang="ts">
import type { SellerProduct } from '~/composables/useSellerShop'
import { CATEGORY_SUBCATEGORIES, type ProductCategory, type ProductSubCategory, type ProductMedia } from '~/types/product'

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

const form = reactive<{
  name: string
  description: string
  category: ProductCategory
  subCategory: ProductSubCategory
  price: number
  stock: number
  image: string
  status: 'Active' | 'Draft'
}>({
  name: '',
  description: '',
  category: (categories[0] || 'Men') as ProductCategory,
  subCategory: (CATEGORY_SUBCATEGORIES[categories[0] || 'Men']?.[0] || 'T-Shirts') as ProductSubCategory,
  price: 1000,
  stock: 1,
  image: '',
  status: 'Active' as 'Active' | 'Draft'
})

const availableSubCategories = computed(() => {
  if (!form.category) return []
  return CATEGORY_SUBCATEGORIES[form.category] || []
})

watch(() => form.category, (newCat) => {
  const subs = CATEGORY_SUBCATEGORIES[newCat] || []
  if (!subs.includes(form.subCategory)) {
    form.subCategory = (subs[0] || 'Other') as ProductSubCategory
  }
})

const formMedia = ref<ProductMedia[]>([])

const error = ref('')

watch(
  () => props.product,
  (newProduct) => {
    if (newProduct) {
      form.name = newProduct.name
      form.description = newProduct.description
      form.category = newProduct.category
      form.subCategory = (newProduct.subCategory || CATEGORY_SUBCATEGORIES[newProduct.category]?.[0] || 'Other') as ProductSubCategory
      form.price = newProduct.price
      form.stock = newProduct.stock
      form.image = newProduct.image
      form.status = newProduct.status
      formMedia.value = newProduct.media ? [...newProduct.media] : []
      error.value = ''
    }
  },
  { immediate: true }
)

const selectClasses = 'h-12 w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15'

const submit = () => {
  if (!props.product) return
  error.value = ''

  if (
    !form.name.trim() ||
    !form.description.trim() ||
    form.price < 1 ||
    form.stock < 0 ||
    (!formMedia.value.length && !form.image)
  ) {
    error.value = 'Please provide valid values for all required fields.'
    return
  }

  const hasMedia = formMedia.value.length > 0
  const primaryMedia = formMedia.value.find(m => m.isPrimary)
  const coverUrl = primaryMedia?.url || formMedia.value[0]?.url || form.image

  updateSellerProduct(props.product.id, {
    name: form.name.trim(),
    description: form.description.trim(),
    category: form.category as ProductCategory,
    subCategory: form.subCategory,
    price: form.price,
    stock: form.stock,
    image: coverUrl || form.image,
    status: form.status,
    media: hasMedia ? formMedia.value : undefined
  })

  showToast('Product updated successfully!')
  emit('close')
}
</script>

<template>
  <UiAppModal
    :is-open="isOpen && !!product"
    title="Edit Product"
    max-width="lg"
    @close="emit('close')"
  >
    <form class="space-y-5" @submit.prevent="submit">
      <MarketplaceMediaUploader v-model="formMedia" :existing="props.product?.media || []" />

      <div class="space-y-4">
        <AuthInput
          v-model="form.name"
          label="Product name"
          placeholder="Product name"
          name="edit-product-name"
        />

        <div class="grid gap-4 sm:grid-cols-3">
          <div class="space-y-2">
            <label class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
              Category
            </label>
            <select
              v-model="form.category"
              :class="selectClasses"
            >
              <option v-for="cat in categories" :key="cat" :value="cat">
                {{ cat }}
              </option>
            </select>
          </div>

          <div class="space-y-2">
            <label class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
              Item Type
            </label>
            <select
              v-model="form.subCategory"
              :class="selectClasses"
            >
              <option v-for="subCat in availableSubCategories" :key="subCat" :value="subCat">
                {{ subCat }}
              </option>
            </select>
          </div>

          <div class="space-y-2">
            <label class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
              Status
            </label>
            <select
              v-model="form.status"
              :class="selectClasses"
            >
              <option value="Active">Active</option>
              <option value="Draft">Draft</option>
            </select>
          </div>
        </div>
      </div>

      <div class="space-y-2">
        <label class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
          Description
        </label>
        <textarea
          v-model="form.description"
          rows="3"
          class="w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
        />
      </div>

      <div class="grid gap-5 sm:grid-cols-2">
        <label class="space-y-2">
          <span class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
            Price (ETB)
          </span>
          <input
            v-model.number="form.price"
            type="number"
            min="1"
            step="50"
            class="h-12 w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
          />
        </label>

        <label class="space-y-2">
          <span class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
            Stock Count
          </span>
          <input
            v-model.number="form.stock"
            type="number"
            min="0"
            step="1"
            class="h-12 w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
          />
        </label>
      </div>

      <UiAppAlert v-if="error">
        {{ error }}
      </UiAppAlert>

      <div class="flex justify-end gap-3 border-t border-[#ded6cc] pt-4">
        <UiAppButton
          variant="ghost"
          @click="emit('close')"
        >
          Cancel
        </UiAppButton>

        <UiAppButton type="submit">
          Update Product
        </UiAppButton>
      </div>
    </form>
  </UiAppModal>
</template>
