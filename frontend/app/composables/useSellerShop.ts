import type { Product, ProductCategory, CreateProductDTO, UpdateProductDTO } from '~/types/product'
import type { Shop, CreateShopDTO, UpdateShopDTO, PaymentMethod } from '~/types/shop'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import { useRepositories } from '~/composables/useRepositories'

export type SellerProduct = Product
export type SellerShop = Shop

export const useSellerShop = () => {
  const { shop, user } = useMockDataStore()
  const { shopRepo, productRepo } = useRepositories()

  const hasShop = computed(() => Boolean(shop.value))

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
    price: number
    stock: number
    image: string
  }) => {
    if (!shop.value) return null
    const dto: CreateProductDTO = {
      name: product.name,
      description: product.description,
      category: product.category,
      price: product.price,
      stock: product.stock,
      image: product.image,
      status: 'Active'
    }
    return productRepo.createProduct(shop.value.name, dto)
  }

  const updateSellerProduct = (id: number, updates: UpdateProductDTO) => {
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
    addPaymentMethod
  }
}
