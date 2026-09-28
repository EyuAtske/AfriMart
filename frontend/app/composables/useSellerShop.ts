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

  const fetchSellerProducts = async (): Promise<Product[]> => {
    if (!isLoggedIn.value) return []

    isLoadingProducts.value = true
    try {
      let currentShop = shop.value
      if (!currentShop) {
        currentShop = await shopRepo.getMyShop()
        if (currentShop) {
          shop.value = currentShop
          user.value.role = 'seller'
        }
      }

      if (!currentShop) return []

      const shopId = currentShop.backendId || currentShop.id
      if (!shopId) return []

      const items = await productRepo.getProductsByShop(String(shopId), 1000, 0)
      if (shop.value) {
        shop.value.products = items
      }
      return items
    } catch (err: any) {
      console.warn('Seller products fetch notice:', err?.message || err)
      return shop.value?.products || []
    } finally {
      isLoadingProducts.value = false
    }
  }

  // Hydrate shop from backend when logged in and no shop is loaded yet
  const shopHydrated = useState<boolean>('seller-shop-hydrated', () => false)
  if (import.meta.client && !shopHydrated.value && isLoggedIn.value) {
    shopHydrated.value = true
    fetchSellerProducts().catch((err) => {
      console.warn('Seller shop hydration notice:', err?.message || err)
    })
  }

  const createShop = (details: CreateShopDTO) => {
    const ownerEmail = user.value.email || 'seller@afrimart.com'
    return shopRepo.createShop(ownerEmail, details)
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
    await fetchSellerProducts().catch(() => {
      if (created && shop.value && !shop.value.products.some(p => String(p.id) === String(created.id))) {
        shop.value.products.unshift(created)
      }
    })

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
