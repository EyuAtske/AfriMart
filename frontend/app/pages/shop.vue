<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'
import EditShopModal from '~/components/marketplace/EditShopModal.vue'
import EditProductModal from '~/components/marketplace/EditProductModal.vue'
import type { SellerProduct } from '~/composables/useSellerShop'
import type { ProductCategory } from '~/types/product'

definePageMeta({
  middleware: 'auth'
})

const { shop, hasShop, createShop, addProduct, deleteSellerProduct, toggleProductStatus } = useSellerShop()
const { showToast } = useToast()

const categories: ProductCategory[] = ['Men', 'Women', 'Kids', 'Shoes', 'Accessories']

const shopForm = reactive({
  name: '',
  description: ''
})

const productForm = reactive({
  name: '',
  description: '',
  category: categories[0] as ProductCategory,
  price: 1500,
  stock: 1,
  image: ''
})

const showShopForm = ref(false)
const isEditShopOpen = ref(false)
const isEditProductOpen = ref(false)
const editingProduct = ref<SellerProduct | null>(null)

const shopError = ref('')
const productError = ref('')
const uploadedFileName = ref('')

const buttonClass = 'h-12 rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50'

const submitShop = () => {
  shopError.value = ''

  if (!shopForm.name.trim() || !shopForm.description.trim()) {
    shopError.value = 'Please add your shop name and description.'
    return
  }

  createShop(shopForm)
  showShopForm.value = false
  showToast('Shop created successfully!')
}

const handleImageUpload = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) return

  uploadedFileName.value = file.name
  productForm.image = URL.createObjectURL(file)
}

const submitProduct = () => {
  productError.value = ''

  if (
    !productForm.name.trim() ||
    !productForm.description.trim() ||
    !productForm.category ||
    productForm.price < 1 ||
    productForm.stock < 1 ||
    !productForm.image
  ) {
    productError.value = 'Please upload a picture and complete every product field.'
    return
  }

  addProduct({
    name: productForm.name,
    description: productForm.description,
    category: productForm.category as ProductCategory,
    price: productForm.price,
    stock: productForm.stock,
    image: productForm.image
  })
  showToast(`Added product "${productForm.name}" to marketplace!`)

  productForm.name = ''
  productForm.description = ''
  productForm.category = categories[0] || 'Men'
  productForm.price = 1500
  productForm.stock = 1
  productForm.image = ''
  uploadedFileName.value = ''
}

const openProductEdit = (product: SellerProduct) => {
  editingProduct.value = product
  isEditProductOpen.value = true
}

const handleDeleteProduct = (id: number, name: string) => {
  if (confirm(`Are you sure you want to delete "${name}"?`)) {
    deleteSellerProduct(id)
    showToast(`Product "${name}" deleted.`)
  }
}

const handleToggleStatus = (id: number, currentStatus: string) => {
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
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-8">
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
            <NuxtLink
              :to="`/shops/${shop.name.toLowerCase().replace(/\s+/g, '-')}`"
              class="inline-flex h-10 items-center justify-center rounded-full border border-[#cfc4b5] px-4 text-xs font-medium uppercase tracking-[0.14em] text-[#756a60] transition hover:bg-[#eee8df]"
            >
              View Public Shop
            </NuxtLink>

            <button
              type="button"
              class="inline-flex h-10 items-center justify-center rounded-full border border-[#806344] px-4 text-xs font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
              @click="isEditShopOpen = true"
            >
              Edit Shop
            </button>
          </div>
        </div>

        <div
          v-if="!hasShop"
          class="overflow-hidden rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] shadow-[0_20px_70px_rgba(33,31,29,0.06)]"
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

              <button
                v-if="!showShopForm"
                type="button"
                :class="`${buttonClass} mt-8`"
                @click="showShopForm = true"
              >
                Become a seller
              </button>

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
                    class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                  >
                    Description
                  </label>

                  <textarea
                    id="shop-description"
                    v-model="shopForm.description"
                    rows="4"
                    placeholder="Describe what your shop sells"
                    class="w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                  />
                </div>

                <p
                  v-if="shopError"
                  class="text-sm font-medium text-red-600"
                >
                  {{ shopError }}
                </p>

                <div class="flex flex-wrap gap-3">
                  <button
                    type="submit"
                    :class="buttonClass"
                  >
                    Create shop
                  </button>

                  <button
                    type="button"
                    :class="buttonClass"
                    @click="showShopForm = false"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>

            <div class="hidden bg-[#211f1d] p-8 text-[#f5f1e9] lg:flex lg:flex-col lg:justify-end">
              <p class="font-serif text-5xl leading-none">
                Sell with style.
              </p>

              <p class="mt-5 text-sm leading-6 text-[#d8cfc2]">
                Your shop, products, and payment setup stay together in this account area.
              </p>
            </div>
          </div>
        </div>

        <div
          v-else-if="shop"
          class="space-y-6"
        >
          <!-- Shop Overview Header & Metrics -->
          <section class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
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

              <NuxtLink
                to="/seller/products"
                class="inline-flex shrink-0 h-10 items-center justify-center rounded-full border border-[#806344] px-5 text-xs font-medium uppercase tracking-[0.14em] text-[#5d4b37] hover:bg-[#806344] hover:text-white transition"
              >
                Manage Inventory
              </NuxtLink>
            </div>

            <!-- Quick Metrics Grid -->
            <div class="mt-6 grid grid-cols-2 gap-4 border-t border-[#ded6cc] pt-6 sm:grid-cols-3">
              <div class="rounded-[8px] border border-[#e8e0d5] bg-[#f5f1e9] p-4">
                <p class="text-[11px] font-medium uppercase tracking-[0.14em] text-[#756a60]">Listed Products</p>
                <p class="mt-1 font-serif text-2xl text-[#211f1d]">{{ shop.products.length }}</p>
              </div>

              <div class="rounded-[8px] border border-[#e8e0d5] bg-[#f5f1e9] p-4">
                <p class="text-[11px] font-medium uppercase tracking-[0.14em] text-[#756a60]">Total Items in Stock</p>
                <p class="mt-1 font-serif text-2xl text-[#211f1d]">{{ totalStockCount }}</p>
              </div>

              <div class="col-span-2 rounded-[8px] border border-[#e8e0d5] bg-[#f5f1e9] p-4 sm:col-span-1">
                <p class="text-[11px] font-medium uppercase tracking-[0.14em] text-[#756a60]">Stock Value</p>
                <p class="mt-1 font-serif text-2xl text-[#211f1d]">{{ formatPrice(totalInventoryValue) }}</p>
              </div>
            </div>
          </section>

          <!-- Product Creation Form -->
          <section class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
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
              <label class="flex min-h-72 cursor-pointer flex-col items-center justify-center overflow-hidden rounded-[10px] border border-dashed border-[#b9aa98] bg-[#f5f1e9] text-center transition hover:border-[#806344]">
                <img
                  v-if="productForm.image"
                  :src="productForm.image"
                  alt="Product preview"
                  class="h-full max-h-80 w-full object-cover"
                />

                <span
                  v-else
                  class="px-6"
                >
                  <span class="block font-serif text-2xl text-[#211f1d]">
                    Upload picture
                  </span>

                  <span class="mt-2 block text-sm leading-6 text-[#756a60]">
                    Choose a clear photo of the product.
                  </span>
                </span>

                <input
                  type="file"
                  accept="image/*"
                  class="sr-only"
                  @change="handleImageUpload"
                />
              </label>

              <div class="space-y-5">
                <AuthInput
                  v-model="productForm.name"
                  label="Product name"
                  placeholder="Name your product"
                  name="product-name"
                />

                <div class="space-y-2">
                  <label
                    for="product-description"
                    class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                  >
                    Description
                  </label>

                  <textarea
                    id="product-description"
                    v-model="productForm.description"
                    rows="4"
                    placeholder="Describe size, condition, fabric, and anything buyers should know"
                    class="w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                  />
                </div>

                <div class="space-y-2">
                  <label
                    for="product-category"
                    class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]"
                  >
                    Category
                  </label>

                  <select
                    id="product-category"
                    v-model="productForm.category"
                    class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                  >
                    <option
                      v-for="category in categories"
                      :key="category"
                      :value="category"
                    >
                      {{ category }}
                    </option>
                  </select>
                </div>

                <div class="grid gap-5 sm:grid-cols-2">
                  <label class="space-y-2">
                    <span class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                      Price
                    </span>

                    <input
                      v-model.number="productForm.price"
                      type="number"
                      min="1"
                      step="50"
                      class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                    />
                  </label>

                  <label class="space-y-2">
                    <span class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
                      Stock
                    </span>

                    <input
                      v-model.number="productForm.stock"
                      type="number"
                      min="1"
                      step="1"
                      class="h-12 w-full rounded-[5px] border border-[#cfc4b5] bg-[#faf8f4] px-4 text-sm text-[#211f1d] outline-none transition-all hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
                    />
                  </label>
                </div>

                <p
                  v-if="uploadedFileName"
                  class="text-xs text-[#756a60]"
                >
                  Selected: {{ uploadedFileName }}
                </p>

                <p
                  v-if="productError"
                  class="text-sm font-medium text-red-600"
                >
                  {{ productError }}
                </p>

                <button
                  type="submit"
                  :class="buttonClass"
                >
                  List product
                </button>
              </div>
            </form>
          </section>

          <!-- Listed Products Management List -->
          <section
            v-if="shop.products.length"
            class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8"
          >
            <div class="mb-6 flex items-center justify-between">
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
                class="overflow-hidden rounded-[10px] border border-[#ded6cc] bg-[#f5f1e9] flex flex-col justify-between"
              >
                <div>
                  <div class="relative">
                    <img
                      :src="product.image"
                      :alt="product.name"
                      class="h-48 w-full object-cover"
                    />
                    <span
                      class="absolute top-3 right-3 rounded-full px-3 py-1 text-[11px] font-semibold tracking-wider text-white uppercase shadow-md"
                      :class="product.status === 'Active' ? 'bg-[#806344]' : 'bg-gray-500'"
                    >
                      {{ product.status }}
                    </span>
                  </div>

                  <div class="p-4">
                    <p class="text-xs font-medium uppercase tracking-[0.14em] text-[#806344]">
                      {{ product.category }}
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
          </section>
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
