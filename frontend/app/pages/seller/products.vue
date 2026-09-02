<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'
import EditProductModal from '~/components/marketplace/EditProductModal.vue'
import type { SellerProduct } from '~/composables/useSellerShop'

definePageMeta({
  middleware: 'auth'
})

const { shop, hasShop, deleteSellerProduct, toggleProductStatus, updateStock } = useSellerShop()
const { showToast } = useToast()

const selectedCategory = ref('All')
const selectedStatus = ref('All')
const isEditProductOpen = ref(false)
const editingProduct = ref<SellerProduct | null>(null)

const categories = computed(() => {
  if (!shop.value) return ['All']
  return ['All', ...Array.from(new Set(shop.value.products.map(p => p.category)))]
})

const filteredProducts = computed(() => {
  if (!shop.value) return []
  return shop.value.products.filter((p) => {
    const matchesCategory = selectedCategory.value === 'All' || p.category === selectedCategory.value
    const matchesStatus = selectedStatus.value === 'All' || p.status === selectedStatus.value
    return matchesCategory && matchesStatus
  })
})

const openEditModal = (product: SellerProduct) => {
  editingProduct.value = product
  isEditProductOpen.value = true
}

const handleDelete = (id: number, name: string) => {
  if (confirm(`Delete "${name}" permanently?`)) {
    deleteSellerProduct(id)
    showToast(`Product "${name}" deleted.`)
  }
}

const handleToggleStatus = (id: number, currentStatus: string) => {
  toggleProductStatus(id)
  const nextStatus = currentStatus === 'Active' ? 'Draft' : 'Active'
  showToast(`Product status changed to ${nextStatus}`)
}

const adjustStock = (id: number, currentStock: number, delta: number) => {
  const newStock = Math.max(0, currentStock + delta)
  updateStock(id, newStock)
  showToast(`Stock updated to ${newStock}`)
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-8">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 lg:flex-row">
      <AccountSidebar active="shop" />

      <section class="min-w-0 flex-1">
        <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="font-serif text-4xl tracking-[-0.025em] text-[#211f1d]">
              Seller Products
            </h1>

            <p class="mt-2 text-base text-[#756a60]">
              Manage inventory, update stock, and control product visibility.
            </p>
          </div>

          <NuxtLink
            v-if="hasShop"
            to="/shop"
            class="inline-flex h-12 items-center justify-center rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
          >
            + Add Product
          </NuxtLink>
        </div>

        <section
          v-if="hasShop && shop"
          class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8"
        >
          <!-- Filters bar -->
          <div class="flex flex-col gap-4 border-b border-[#ded6cc] pb-6 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
                {{ shop.name }} Inventory
              </p>
              <p class="text-sm font-medium text-[#756a60] mt-1">
                Showing {{ filteredProducts.length }} of {{ shop.products.length }} products
              </p>
            </div>

            <div class="flex flex-wrap items-center gap-3">
              <select
                v-model="selectedCategory"
                class="h-10 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-4 text-xs text-[#211f1d] outline-none hover:border-[#806344]"
              >
                <option v-for="cat in categories" :key="cat" :value="cat">
                  Category: {{ cat }}
                </option>
              </select>

              <select
                v-model="selectedStatus"
                class="h-10 rounded-full border border-[#cfc4b5] bg-[#f5f1e9] px-4 text-xs text-[#211f1d] outline-none hover:border-[#806344]"
              >
                <option value="All">Status: All</option>
                <option value="Active">Status: Active</option>
                <option value="Draft">Status: Draft</option>
              </select>
            </div>
          </div>

          <!-- Product Grid -->
          <div
            v-if="filteredProducts.length"
            class="mt-6 grid gap-4 sm:grid-cols-2 xl:grid-cols-3"
          >
            <article
              v-for="product in filteredProducts"
              :key="product.id"
              class="overflow-hidden rounded-[10px] border border-[#ded6cc] bg-[#f5f1e9] flex flex-col justify-between"
            >
              <div>
                <div class="relative">
                  <img
                    :src="product.image"
                    :alt="product.name"
                    class="h-48 w-full object-cover object-top"
                  />
                  <span
                    class="absolute top-3 right-3 rounded-full px-3 py-1 text-[11px] font-semibold uppercase tracking-wider text-white shadow-md"
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

                  <p class="mt-2 text-sm leading-6 text-[#756a60] line-clamp-2">
                    {{ product.description }}
                  </p>

                  <div class="mt-4 flex items-center justify-between gap-3 text-sm">
                    <span class="font-semibold text-[#211f1d]">
                      {{ formatPrice(product.price) }}
                    </span>

                    <!-- Quick Stock Adjuster -->
                    <div class="flex items-center gap-2 rounded-full border border-[#cfc4b5] bg-[#faf8f4] px-2 py-1 text-xs">
                      <button
                        type="button"
                        class="h-5 w-5 rounded-full bg-[#ded6cc] text-[#211f1d] hover:bg-[#806344] hover:text-white"
                        @click="adjustStock(product.id, product.stock, -1)"
                      >
                        -
                      </button>
                      <span class="font-medium text-[#211f1d] min-w-[20px] text-center">{{ product.stock }}</span>
                      <button
                        type="button"
                        class="h-5 w-5 rounded-full bg-[#ded6cc] text-[#211f1d] hover:bg-[#806344] hover:text-white"
                        @click="adjustStock(product.id, product.stock, 1)"
                      >
                        +
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Action Toolbar -->
              <div class="flex items-center justify-between border-t border-[#ded6cc] px-4 py-3 bg-[#eee8df]/60 text-xs">
                <div class="flex gap-3">
                  <button
                    type="button"
                    class="font-medium text-[#806344] hover:underline"
                    @click="openEditModal(product)"
                  >
                    Edit
                  </button>

                  <button
                    type="button"
                    class="font-medium text-[#756a60] hover:underline"
                    @click="handleToggleStatus(product.id, product.status)"
                  >
                    {{ product.status === 'Active' ? 'Make Draft' : 'Make Active' }}
                  </button>
                </div>

                <button
                  type="button"
                  class="font-medium text-red-600 hover:underline"
                  @click="handleDelete(product.id, product.name)"
                >
                  Delete
                </button>
              </div>
            </article>
          </div>

          <p
            v-else
            class="mt-6 rounded-[8px] border border-[#ded6cc] bg-[#f5f1e9] p-8 text-center text-sm text-[#756a60]"
          >
            No products match the selected filters.
          </p>
        </section>

        <section
          v-else
          class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-8 shadow-[0_20px_70px_rgba(33,31,29,0.06)]"
        >
          <h2 class="font-serif text-3xl text-[#211f1d]">
            Create a shop first
          </h2>

          <p class="mt-3 text-sm leading-6 text-[#756a60]">
            Product management appears after your seller shop is created.
          </p>

          <NuxtLink
            to="/shop"
            class="mt-6 inline-flex h-12 items-center justify-center rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition hover:bg-[#806344] hover:text-white"
          >
            Start selling
          </NuxtLink>
        </section>

        <EditProductModal
          :is-open="isEditProductOpen"
          :product="editingProduct"
          @close="isEditProductOpen = false"
        />
      </section>
    </div>
  </main>
</template>

