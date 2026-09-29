import type { Product, ProductCategory, ProductGender, ProductSubCategory, CreateProductDTO, UpdateProductDTO, ProductMedia } from '~/types/product'
import type { Shop, CreateShopDTO, UpdateShopDTO, PaymentMethod } from '~/types/shop'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import { useRepositories } from '~/composables/useRepositories'

export type SellerProduct = Product
export type SellerShop = Shop

export const useSellerShop = () => {
  const { shop, user, isLoggedIn } = useMockDataStore()
  const { shopRepo, productRepo } = useRepositories()

  const hasShop = computed(() => Boolean(shop.value))
  const isLoadingProducts = useState<boolean>('seller-products-loading', () => false)
  const isLoadingShop = useState<boolean>('seller-shop-loading', () => isLoggedIn.value)
  const shopError = useState<string>('seller-shop-error', () => '')
  const shopHydratedFor = useState<string>('seller-shop-hydrated-for', () => '')
  const shopRequestId = useState<number>('seller-shop-request-id', () => 0)
  const currentAccountKey = () => String(user.value.id || user.value.email || '')

  const fetchSellerProducts = async (): Promise<Product[]> => {
    if (!isLoggedIn.value || !currentAccountKey()) {
      shop.value = null
      shopError.value = ''
      isLoadingShop.value = false
      isLoadingProducts.value = false
      return []
    }

    const accountKey = currentAccountKey()
    const requestId = ++shopRequestId.value
    let resolvedShop: Shop | null = null
    shopError.value = ''
    isLoadingShop.value = true
    isLoadingProducts.value = true
    try {
      const currentShop = await shopRepo.getMyShop()
      if (shopRequestId.value !== requestId || !isLoggedIn.value || currentAccountKey() !== accountKey) return []

      if (!currentShop) {
        shop.value = null
        user.value.role = 'buyer'
        return []
      }

      resolvedShop = currentShop
      currentShop.products = []
      shop.value = currentShop
      user.value.role = 'seller'

      const shopId = currentShop.backendId || currentShop.id
      if (!shopId) return []

      const items = await productRepo.getProductsByShop(String(shopId), 1000, 0)
      if (shopRequestId.value !== requestId || !isLoggedIn.value || currentAccountKey() !== accountKey) return []

      if (shop.value && String(shop.value.backendId || shop.value.id) === String(shopId)) {
        shop.value.products = items
      }
      return items
    } catch (err: any) {
      console.warn('Seller products fetch notice:', err?.message || err)
      if (shopRequestId.value === requestId && isLoggedIn.value && currentAccountKey() === accountKey) {
        shopError.value = err?.message || 'Unable to load your shop. Please try again.'
        shop.value = resolvedShop
        if (resolvedShop) resolvedShop.products = []
      }
      return []
    } finally {
      if (shopRequestId.value === requestId) {
        isLoadingShop.value = false
        isLoadingProducts.value = false
      }
    }
  }

  // Hydrate shop from backend when logged in and no shop is loaded yet
  watch(
    [isLoggedIn, currentAccountKey],
    ([loggedIn, accountKey]) => {
      if (!loggedIn) {
        shopHydratedFor.value = ''
          shopRequestId.value++
          shop.value = null
          shopError.value = ''
          isLoadingShop.value = false
          isLoadingProducts.value = false
        return
      }
      if (!import.meta.client || !accountKey || shopHydratedFor.value === accountKey) return

      shopHydratedFor.value = accountKey
      fetchSellerProducts().catch((err) => {
        console.warn('Seller shop hydration notice:', err?.message || err)
      })
    },
    { immediate: true }
  )

  const createShop = async (details: CreateShopDTO) => {
    const ownerEmail = user.value.email || 'seller@afrimart.com'
    const accountKey = currentAccountKey()
    const requestId = ++shopRequestId.value
    shopError.value = ''
    isLoadingShop.value = true
    try {
      const created = await shopRepo.createShop(ownerEmail, details)
      if (
        shopRequestId.value !== requestId ||
        !isLoggedIn.value ||
        currentAccountKey() !== accountKey
      ) return created

      created.products = []
      shop.value = created
      user.value.role = 'seller'
      shopHydratedFor.value = accountKey
      return created
    } catch (err: any) {
      if (shopRequestId.value === requestId && isLoggedIn.value && currentAccountKey() === accountKey) {
        shopError.value = err?.message || 'Unable to create your shop. Please try again.'
      }
      throw err
    } finally {
      if (shopRequestId.value === requestId) isLoadingShop.value = false
    }
  }

  const updateShop = (details: UpdateShopDTO) => {
    if (!shop.value) return null
    return shopRepo.updateShop(shop.value.slug, details)
  }

  const addProduct = async (product: {
    name: string
    description: string
    category: ProductCategory
    categoryId?: string
    gender?: ProductGender
    subCategory?: ProductSubCategory
    subcategoryId?: string
    brand?: string
    color?: string
    size?: string
    price: number
    stock: number
    image?: string
    media?: ProductMedia[]
    files?: File[]
  }) => {
    if (!shop.value) {
      const myShop = await shopRepo.getMyShop()
      if (myShop) {
        shop.value = myShop
        user.value.role = 'seller'
      } else {
        throw new Error('You must create a shop before creating products.')
      }
    }

    // Extract File objects from media items if files not directly passed
    const files = product.files || product.media?.map(m => m.file).filter((f): f is File => f instanceof File) || []

    // Derive image from primary media for backward compat
    const primaryMedia = product.media?.find(m => m.isPrimary)
    const image = primaryMedia?.url || product.image || ''

    const dto: CreateProductDTO = {
      name: product.name,
      description: product.description,
      category: product.category,
      categoryId: product.categoryId,
      gender: product.gender,
      subCategory: product.subCategory,
      subcategoryId: product.subcategoryId,
      brand: product.brand,
      color: product.color,
      size: product.size,
      price: product.price,
      stock: product.stock,
      image,
      status: 'Active',
      media: product.media,
      files
    }
    const created = await productRepo.createProduct(shop.value.name, dto)

    // Refresh seller products from backend to ensure consistent state
    const refreshedProducts = await fetchSellerProducts()
    if (
      shop.value &&
      !refreshedProducts.some(p => String(p.id) === String(created.id)) &&
      (!created.shopId || String(created.shopId) === String(shop.value.backendId || shop.value.id)) &&
      !shop.value.products.some(p => String(p.id) === String(created.id))
    ) {
      shop.value.products.unshift(created)
    }

    return created
  }

  const updateSellerProduct = (id: number | string, updates: UpdateProductDTO) => {
    // Derive image from primary media if media is being updated
    if (updates.media?.length) {
      const primaryMedia = updates.media.find(m => m.isPrimary)
      if (primaryMedia) {
        updates.image = primaryMedia.url
      }
    }
    return productRepo.updateProduct(id, updates)
  }

  const deleteSellerProduct = (id: number | string) => {
    return productRepo.deleteProduct(id)
  }

  const toggleProductStatus = (id: number | string) => {
    return productRepo.toggleProductStatus(id)
  }

  const updateStock = (id: number | string, stock: number) => {
    const validStock = Math.max(0, stock)
    return productRepo.updateProduct(id, { stock: validStock })
  }

  const addPaymentMethod = (method: {
    type: 'Telebirr' | 'CBE'
    accountName: string
    accountNumber: string
  }) => {
    if (!shop.value) return

    shop.value.paymentMethods.unshift({
      id: Date.now(),
      type: method.type,
      accountName: method.accountName.trim(),
      accountNumber: method.accountNumber.trim()
    })
  }

  const deactivateShop = async () => {
    if (!shop.value) return null
    const shopId = shop.value.backendId || shop.value.id
    return shopRepo.deactivateShop(shopId)
  }

  const activateShop = async () => {
    if (!shop.value) return null
    const shopId = shop.value.backendId || shop.value.id
    return shopRepo.activateShop(shopId)
  }

  return {
    shop,
    hasShop,
    isLoadingShop,
    shopError,
    isLoadingProducts,
    fetchSellerProducts,
    createShop,
    updateShop,
    addProduct,
    updateSellerProduct,
    deleteSellerProduct,
    toggleProductStatus,
    updateStock,
    addPaymentMethod,
    deactivateShop,
    activateShop
  }
}
