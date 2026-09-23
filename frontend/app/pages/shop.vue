<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'
import EditShopModal from '~/components/marketplace/EditShopModal.vue'
import EditProductModal from '~/components/marketplace/EditProductModal.vue'
import MediaUploader from '~/components/marketplace/MediaUploader.vue'
import type { SellerProduct } from '~/composables/useSellerShop'
import { MAIN_CATEGORIES, GENDERS, HIERARCHICAL_ITEMS, type ProductCategory, type ProductGender, type ProductSubCategory, type ProductMedia } from '~/types/product'
import { extractError } from '~/repositories/api/apiHelpers'

definePageMeta({
  middleware: 'auth'
})
const { gtag } = useGtag()
const { shop, hasShop, createShop, addProduct, deleteSellerProduct, toggleProductStatus } = useSellerShop()
const { showToast } = useToast()

const isSubmitting = ref(false)

const categories = MAIN_CATEGORIES
const genders = GENDERS

const shopForm = reactive({
  name: '',
  description: ''
})

const productForm = reactive<{
  name: string
  description: string
  category: ProductCategory
  gender: ProductGender
  subCategory: ProductSubCategory
  price: number
  stock: number
  image: string
}>({
  name: '',
  description: '',
  category: 'Clothing',
  gender: 'Men',
  subCategory: 'T-Shirts',
  price: 1500,
  stock: 1,
  image: ''
})

const availableSubCategories = computed(() => {
  const cat = (productForm.category === 'Accessories' ? 'Accessories' : 'Clothing') as 'Clothing' | 'Accessories'
  const gen = (productForm.gender || 'Men') as ProductGender
  return HIERARCHICAL_ITEMS[cat]?.[gen] || []
})

watch([() => productForm.category, () => productForm.gender], () => {
  const subs = availableSubCategories.value
  if (!subs.includes(productForm.subCategory)) {
    productForm.subCategory = (subs[0] || 'Other') as ProductSubCategory
  }
})

const productMedia = ref<ProductMedia[]>([])

const showShopForm = ref(false)
const isEditShopOpen = ref(false)
const isEditProductOpen = ref(false)
const editingProduct = ref<SellerProduct | null>(null)

const shopError = ref('')
const productError = ref('')

const fieldErrors = reactive({
  name: '',
  description: '',
  category: '',
  subCategory: '',
  price: '',
  stock: '',
  media: ''
})

const clearFieldErrors = () => {
  fieldErrors.name = ''
  fieldErrors.description = ''
  fieldErrors.category = ''
  fieldErrors.subCategory = ''
  fieldErrors.price = ''
  fieldErrors.stock = ''
  fieldErrors.media = ''
}

const submitShop = () => {
  shopError.value = ''

  if (!shopForm.name.trim() || !shopForm.description.trim()) {
    shopError.value = 'Please add your shop name and description.'
    return
  }

  createShop(shopForm)

  gtag('event', 'shop_created', {
    shop_name: shopForm.name
  })

  showShopForm.value = false
  showToast('Shop created successfully!')
}

const submitProduct = async () => {
  productError.value = ''
  clearFieldErrors()

  let isValid = true

  if (!productForm.name.trim()) {
    fieldErrors.name = 'Product name is required.'
    isValid = false
  }

  if (!productForm.description.trim()) {
    fieldErrors.description = 'Product description is required.'
    isValid = false
  }

  if (!productForm.category) {
    fieldErrors.category = 'Category is required.'
    isValid = false
  }

  if (!productForm.subCategory) {
    fieldErrors.subCategory = 'Item type is required.'
    isValid = false
  }

  if (productForm.price === null || productForm.price === undefined || isNaN(productForm.price) || productForm.price < 0) {
    fieldErrors.price = 'Price must be a valid non-negative number.'
    isValid = false
  }

  if (productForm.stock === null || productForm.stock === undefined || isNaN(productForm.stock) || productForm.stock < 1) {
    fieldErrors.stock = 'Stock must be at least 1.'
    isValid = false
  }

  const hasMedia = productMedia.value.length > 0
  const primaryImage = productMedia.value.find(m => m.isPrimary)
  const coverUrl = primaryImage?.url || productMedia.value[0]?.url || productForm.image

  if (!hasMedia && !productForm.image) {
    fieldErrors.media = 'At least one product image is required.'
    isValid = false
  } else if (productMedia.value.length > 10) {
    fieldErrors.media = 'A maximum of 10 images can be uploaded.'
    isValid = false
  }

  if (!isValid) {
    productError.value = 'Please fix the highlighted field errors below.'
    return
  }

  isSubmitting.value = true

  try {
    const created = await addProduct({
      name: productForm.name.trim(),
      description: productForm.description.trim(),
      category: productForm.category,
      gender: productForm.gender,
      subCategory: productForm.subCategory,
      price: productForm.price,
      stock: productForm.stock,
      image: coverUrl || '',
      media: hasMedia ? productMedia.value : undefined
    })

    if (!created) {
      throw new Error('Failed to create product.')
    }

    gtag('event', 'product_created', {
      product_name: productForm.name,
      category: productForm.category,
      price: productForm.price
    })

    showToast(`Added product "${productForm.name}" to marketplace!`)

    productForm.name = ''
    productForm.description = ''
    productForm.category = 'Clothing'
    productForm.gender = 'Men'
    productForm.subCategory = 'T-Shirts'
    productForm.price = 1500
    productForm.stock = 1
    productForm.image = ''
    productMedia.value = []
    clearFieldErrors()
  } catch (err: any) {
    productError.value = extractError(err, 'Failed to publish product. Please try again.')
    showToast(productError.value, 'error')
  } finally {
    isSubmitting.value = false
  }
}

const openProductEdit = (product: SellerProduct) => {
  editingProduct.value = product
  isEditProductOpen.value = true
}

const handleDeleteProduct = (id: number | string, name: string) => {
  if (confirm(`Are you sure you want to delete "${name}"?`)) {
    deleteSellerProduct(id)
    showToast(`Product "${name}" deleted.`)
  }
}

const handleToggleStatus = (id: number | string, currentStatus: string) => {
  toggleProductStatus(id)
  const newStatus = currentStatus === 'Active' ? 'Draft' : 'Active'
  showToast(`Product status updated to ${newStatus}.`)
}

const totalStockCount = computed(() => {
  if (!shop.value) return 0
  return shop.value.products.reduce((acc, p) => acc + p.stock, 0)
})

const totalInventoryValue = computed(() => {
  if (!shop.value) return 0
  return shop.value.products.reduce((acc, p) => acc + (p.price * p.stock), 0)
})
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 lg:flex-row">
      <AccountSidebar active="shop" />

      <section class="min-w-0 flex-1">
        <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="font-serif text-4xl tracking-[-0.025em] text-[#211f1d]">
              My Shop
            </h1>

            <p class="mt-2 text-base text-[#756a60]">
              Create your seller profile, manage listings, and update store settings.
            </p>
          </div>

          <div v-if="hasShop && shop" class="flex gap-2">
            <UiAppButton
              :to="`/shops/${shop.name.toLowerCase().replace(/\s+/g, '-')}`"
              variant="ghost"
              size="small"
            >
              View Public Shop
            </UiAppButton>

            <UiAppButton
              variant="secondary"
              size="small"
              @click="isEditShopOpen = true"
            >
              Edit Shop
            </UiAppButton>
          </div>
        </div>

        <UiAppCard
          v-if="!hasShop"
          padding="none"
          class="overflow-hidden"
        >
          <div class="grid lg:grid-cols-[1fr_0.85fr]">
            <div class="p-6 sm:p-8">
              <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
                Seller tools
              </p>

              <h2 class="mt-3 font-serif text-3xl leading-tight text-[#211f1d] sm:text-4xl">
                Start selling your pieces on Afrimart.
              </h2>

              <p class="mt-4 max-w-xl text-base leading-7 text-[#756a60]">
                Give your shop a name, add a short description, and your product listing tools will open here.
              </p>

              <div v-if="!showShopForm" class="mt-8">
                <UiAppButton
                  variant="secondary"
                  @click="showShopForm = true"
                >
                  Become a seller
                </UiAppButton>
              </div>

              <form
                v-else
                class="mt-8 space-y-5"
                @submit.prevent="submitShop"
              >
                <AuthInput
                  v-model="shopForm.name"
                  label="Shop name"
                  placeholder="Give your shop a name"
                  name="shop-name"
                />

                <div class="space-y-2">
                  <label
                    for="shop-description"
                    class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                  >
                    Description
                  </label>

                  <textarea
                    id="shop-description"
                    v-model="shopForm.description"
                    rows="4"
                    placeholder="Describe what your shop sells"
                    class="w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                  />
                </div>

                <UiAppAlert v-if="shopError">
                  {{ shopError }}
                </UiAppAlert>

                <div class="flex flex-wrap gap-3">
                  <UiAppButton
                    type="submit"
                    variant="secondary"
                  >
                    Create shop
                  </UiAppButton>

                  <UiAppButton
                    variant="ghost"
                    @click="showShopForm = false"
                  >
                    Cancel
                  </UiAppButton>
                </div>
              </form>
            </div>

            <div class="relative hidden overflow-hidden bg-[#211f1d] p-8 text-[#f5f1e9] lg:flex lg:flex-col lg:justify-end">
              <img
                src="/images/product1.jpg"
                alt="Sell with style on Afrimart"
                class="absolute inset-0 h-full w-full object-cover object-center"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-[#211f1d]/95 via-[#211f1d]/60 to-[#211f1d]/30" />

              <div class="relative z-10">
                <p class="font-serif text-5xl leading-none text-white">
                  Sell with style.
                </p>

                <p class="mt-5 text-sm leading-6 text-[#ded6cc]">
                  Your shop, products, and payment setup stay together in this account area.
                </p>
              </div>
            </div>
          </div>
        </UiAppCard>

        <div
          v-else-if="shop"
          class="space-y-6"
        >
          <!-- Shop Overview Header & Metrics -->
          <UiAppCard padding="large">
            <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
              <div>
                <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
                  Seller shop
                </p>

                <h2 class="mt-2 font-serif text-3xl text-[#211f1d]">
                  {{ shop.name }}
                </h2>

                <p class="mt-3 max-w-2xl text-base leading-7 text-[#756a60]">
                  {{ shop.description }}
                </p>
              </div>

              <UiAppButton
                to="/seller/products"
                variant="secondary"
                size="small"
                class="shrink-0"
              >
                Manage Inventory
              </UiAppButton>
            </div>

            <!-- Quick Metrics Grid -->
            <div class="mt-6 grid grid-cols-2 gap-4 border-t border-[#ded6cc] pt-6 sm:grid-cols-3">
              <div class="rounded-lg border border-[#e8e0d5] bg-[#f5f1e9] p-4">
                <p class="text-xs font-medium uppercase tracking-[0.14em] text-[#756a60]">Listed Products</p>
                <p class="mt-1 font-serif text-2xl text-[#211f1d]">{{ shop.products.length }}</p>
              </div>

              <div class="rounded-lg border border-[#e8e0d5] bg-[#f5f1e9] p-4">
                <p class="text-xs font-medium uppercase tracking-[0.14em] text-[#756a60]">Total Items in Stock</p>
                <p class="mt-1 font-serif text-2xl text-[#211f1d]">{{ totalStockCount }}</p>
              </div>

              <div class="col-span-2 rounded-lg border border-[#e8e0d5] bg-[#f5f1e9] p-4 sm:col-span-1">
                <p class="text-xs font-medium uppercase tracking-[0.14em] text-[#756a60]">Stock Value</p>
                <p class="mt-1 font-serif text-2xl text-[#211f1d]">{{ formatPrice(totalInventoryValue) }}</p>
              </div>
            </div>
          </UiAppCard>

          <!-- Product Creation Form -->
          <UiAppCard padding="large">
            <div class="mb-6">
              <h2 class="text-xl font-medium text-[#211f1d]">
                List a product
              </h2>

              <p class="mt-1 text-sm text-[#756a60]">
                Upload a product picture, describe it, and choose a category.
              </p>
            </div>

            <form
              class="grid gap-6 lg:grid-cols-[0.8fr_1fr]"
              @submit.prevent="submitProduct"
            >
              <div>
                <MediaUploader v-model="productMedia" />
                <p v-if="fieldErrors.media" class="mt-2 text-xs font-medium text-red-600">
                  {{ fieldErrors.media }}
                </p>
              </div>

              <div class="space-y-5">
                <div>
                  <AuthInput
                    v-model="productForm.name"
                    label="Product name"
                    placeholder="Name your product"
                    name="product-name"
                  />
                  <p v-if="fieldErrors.name" class="mt-1 text-xs font-medium text-red-600">
                    {{ fieldErrors.name }}
                  </p>
                </div>

                <div class="space-y-2">
                  <label
                    for="product-description"
                    class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                  >
                    Description
                  </label>

                  <textarea
                    id="product-description"
                    v-model="productForm.description"
                    rows="4"
                    placeholder="Describe size, condition, fabric, and anything buyers should know"
                    class="w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                  />
                  <p v-if="fieldErrors.description" class="mt-1 text-xs font-medium text-red-600">
                    {{ fieldErrors.description }}
                  </p>
                </div>

                <div class="grid gap-5 sm:grid-cols-3">
                  <div class="space-y-2">
                    <label
                      for="product-category"
                      class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                    >
                      Category
                    </label>

                    <select
                      id="product-category"
                      v-model="productForm.category"
                      class="h-12 w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                    >
                      <option
                        v-for="category in categories"
                        :key="category"
                        :value="category"
                      >
                        {{ category }}
                      </option>
                    </select>
                    <p v-if="fieldErrors.category" class="mt-1 text-xs font-medium text-red-600">
                      {{ fieldErrors.category }}
                    </p>
                  </div>

                  <div class="space-y-2">
                    <label
                      for="product-gender"
                      class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                    >
                      Gender
                    </label>

                    <select
                      id="product-gender"
                      v-model="productForm.gender"
                      class="h-12 w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                    >
                      <option
                        v-for="gender in genders"
                        :key="gender"
                        :value="gender"
                      >
                        {{ gender }}
                      </option>
                    </select>
                  </div>

                  <div class="space-y-2">
                    <label
                      for="product-subcategory"
                      class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                    >
                      Item Type
                    </label>

                    <select
                      id="product-subcategory"
                      v-model="productForm.subCategory"
                      class="h-12 w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                    >
                      <option
                        v-for="subCat in availableSubCategories"
                        :key="subCat"
                        :value="subCat"
                      >
                        {{ subCat }}
                      </option>
                    </select>
                    <p v-if="fieldErrors.subCategory" class="mt-1 text-xs font-medium text-red-600">
                      {{ fieldErrors.subCategory }}
                    </p>
                  </div>
                </div>

                <div class="grid gap-5 sm:grid-cols-2">
                  <div class="space-y-2">
                    <label
                      for="product-price"
                      class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                    >
                      Price (ETB)
                    </label>

                    <input
                      id="product-price"
                      v-model.number="productForm.price"
                      type="number"
                      min="0"
                      step="any"
                      class="h-12 w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                    />
                    <p v-if="fieldErrors.price" class="mt-1 text-xs font-medium text-red-600">
                      {{ fieldErrors.price }}
                    </p>
                  </div>

                  <div class="space-y-2">
                    <label
                      for="product-stock"
                      class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                    >
                      Stock
                    </label>

                    <input
                      id="product-stock"
                      v-model.number="productForm.stock"
                      type="number"
                      min="1"
                      step="1"
                      class="h-12 w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                    />
                    <p v-if="fieldErrors.stock" class="mt-1 text-xs font-medium text-red-600">
                      {{ fieldErrors.stock }}
                    </p>
                  </div>
                </div>

                <UiAppAlert v-if="productError">
                  {{ productError }}
                </UiAppAlert>

                <UiAppButton
                  type="submit"
                  variant="secondary"
                  :loading="isSubmitting"
                >
                  {{ isSubmitting ? 'Publishing product…' : 'List product' }}
                </UiAppButton>
              </div>
            </form>
          </UiAppCard>

          <!-- Listed Products Management List -->
          <UiAppCard
            v-if="shop.products.length"
            padding="large"
          >
            <div class="mb-6 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
              <h2 class="text-xl font-medium text-[#211f1d]">
                Listed Products ({{ shop.products.length }})
              </h2>

              <NuxtLink
                to="/seller/products"
                class="text-xs font-medium uppercase tracking-[0.14em] text-[#806344] hover:underline"
              >
                Full Inventory Management →
              </NuxtLink>
            </div>

            <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
              <article
                v-for="product in shop.products"
                :key="product.id"
                class="overflow-hidden rounded-lg border border-[#ded6cc] bg-[#f5f1e9] flex flex-col justify-between"
              >
                <div>
                  <div class="relative">
                    <img
                      :src="product.image"
                      :alt="product.name"
                      class="h-48 w-full object-cover"
                    />
                    <div class="absolute top-3 right-3">
                      <UiAppBadge
                        :variant="product.status === 'Active' ? 'brand' : 'muted'"
                      >
                        {{ product.status }}
                      </UiAppBadge>
                    </div>
                  </div>

                  <div class="p-4">
                    <p class="text-xs font-medium uppercase tracking-[0.14em] text-[#806344]">
                      {{ product.category }}<span v-if="product.subCategory"> · {{ product.subCategory }}</span>
                    </p>

                    <h3 class="mt-1 font-serif text-xl text-[#211f1d]">
                      {{ product.name }}
                    </h3>

                    <p class="mt-2 line-clamp-2 text-sm leading-6 text-[#756a60]">
                      {{ product.description }}
                    </p>

                    <div class="mt-4 flex items-center justify-between gap-3 text-sm">
                      <span class="font-semibold text-[#211f1d]">
                        {{ formatPrice(product.price) }}
                      </span>

                      <span class="text-[#756a60]">
                        {{ product.stock }} in stock
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Action Controls -->
                <div class="flex items-center justify-between border-t border-[#ded6cc] px-4 py-3 bg-[#eee8df]/60 text-xs">
                  <div class="flex gap-2">
                    <button
                      type="button"
                      class="font-medium text-[#806344] hover:underline"
                      @click="openProductEdit(product)"
                    >
                      Edit
                    </button>

                    <button
                      type="button"
                      class="font-medium text-[#756a60] hover:underline"
                      @click="handleToggleStatus(product.id, product.status)"
                    >
                      {{ product.status === 'Active' ? 'Set Draft' : 'Set Active' }}
                    </button>
                  </div>

                  <button
                    type="button"
                    class="font-medium text-red-600 hover:underline"
                    @click="handleDeleteProduct(product.id, product.name)"
                  >
                    Delete
                  </button>
                </div>
              </article>
            </div>
          </UiAppCard>
        </div>

        <EditShopModal
          :is-open="isEditShopOpen"
          @close="isEditShopOpen = false"
        />

        <EditProductModal
          :is-open="isEditProductOpen"
          :product="editingProduct"
          @close="isEditProductOpen = false"
        />
      </section>
    </div>
  </main>
</template>
