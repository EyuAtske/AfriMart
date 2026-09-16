import type { Product, ProductCategory, ProductSubCategory, CreateProductDTO, UpdateProductDTO, ProductMedia } from '~/types/product'
import type { Shop, CreateShopDTO, UpdateShopDTO, PaymentMethod } from '~/types/shop'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import { useRepositories } from '~/composables/useRepositories'

export type SellerProduct = Product
export type SellerShop = Shop

export const useSellerShop = () => {
  const { shop, user, isLoggedIn } = useMockDataStore()
  const { shopRepo, productRepo } = useRepositories()

  const hasShop = computed(() => Boolean(shop.value))

  // Hydrate shop from backend when logged in and no shop is loaded yet
  const shopHydrated = useState<boolean>('seller-shop-hydrated', () => false)
  if (import.meta.client && !shopHydrated.value && isLoggedIn.value && !shop.value) {
    shopHydrated.value = true
    shopRepo.getMyShop().then((myShop) => {
      if (myShop) {
        shop.value = myShop
        user.value.role = 'seller'
      }
    }).catch(() => {
      // Silently fail — user just doesn't have a shop
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

  const addProduct = (product: {
    name: string
    description: string
    category: ProductCategory
    subCategory?: ProductSubCategory
    price: number
    stock: number
    image: string
    media?: ProductMedia[]
  }) => {
    if (!shop.value) return null

    // Derive image from primary media for backward compat
    const primaryMedia = product.media?.find(m => m.isPrimary)
    const image = primaryMedia?.url || product.image

    const dto: CreateProductDTO = {
      name: product.name,
      description: product.description,
      category: product.category,
      subCategory: product.subCategory,
      price: product.price,
      stock: product.stock,
      image,
      status: 'Active',
      media: product.media
    }
    return productRepo.createProduct(shop.value.name, dto)
  }

  const updateSellerProduct = (id: number, updates: UpdateProductDTO) => {
    // Derive image from primary media if media is being updated
    if (updates.media?.length) {
      const primaryMedia = updates.media.find(m => m.isPrimary)
      if (primaryMedia) {
        updates.image = primaryMedia.url
      }
    }
    return productRepo.updateProduct(id, updates)
  }

  const deleteSellerProduct = (id: number) => {
    return productRepo.deleteProduct(id)
  }

  const toggleProductStatus = (id: number) => {
    return productRepo.toggleProductStatus(id)
  }

  const updateStock = (id: number, stock: number) => {
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
