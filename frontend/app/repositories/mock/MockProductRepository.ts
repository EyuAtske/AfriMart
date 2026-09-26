import type { IProductRepository } from '../interfaces/IProductRepository'
import type { Product, ProductFilterParams, CreateProductDTO, UpdateProductDTO } from '~/types/product'
import type { PaginatedResponse } from '~/types/api'
import { useMockDataStore } from './MockDataStore'

export class MockProductRepository implements IProductRepository {
  async getProducts(params: ProductFilterParams = {}): Promise<PaginatedResponse<Product>> {
    const { products } = useMockDataStore()
    const search = params.search?.trim().toLowerCase() || ''
    const category = params.category || 'All'
    const subCategory = params.subCategory || 'All'
    const shop = params.shop?.trim().toLowerCase()
    const status = params.status
    const page = params.page || 1
    const pageSize = params.pageSize || 50

    const filtered = products.value.filter((p) => {
      const matchesSearch = !search ||
        p.name.toLowerCase().includes(search) ||
        p.shop.toLowerCase().includes(search) ||
        p.description.toLowerCase().includes(search) ||
        (p.subCategory && p.subCategory.toLowerCase().includes(search))

      const matchesCategory = category === 'All' || p.category === category
      const matchesSubCategory = subCategory === 'All' || p.subCategory === subCategory
      const matchesShop = !shop || p.shop.toLowerCase() === shop
      const matchesStatus = !status ? p.status === 'Active' : p.status === status

      return matchesSearch && matchesCategory && matchesSubCategory && matchesShop && matchesStatus
    })

    const total = filtered.length
    const totalPages = Math.max(1, Math.ceil(total / pageSize))
    const start = (page - 1) * pageSize
    const data = filtered.slice(start, start + pageSize)

    return {
      data,
      total,
      page,
      pageSize,
      totalPages
    }
  }

  async getProductById(id: number | string): Promise<Product | null> {
    const { products } = useMockDataStore()
    const strId = String(id)
    const found = products.value.find(p => String(p.id) === strId)
    return found ? { ...found } : null
  }

  async createProduct(shopName: string, dto: CreateProductDTO): Promise<Product> {
    const { products, shop } = useMockDataStore()
    const id = Date.now() + Math.floor(Math.random() * 1000)

    const newProduct: Product = {
      id,
      shop: shopName,
      name: dto.name.trim(),
      description: dto.description.trim(),
      category: dto.category,
      subCategory: dto.subCategory,
      price: dto.price,
      stock: dto.stock,
      rating: 'New',
      image: dto.image || '/images/shop.jpg',
      status: dto.status || 'Active',
      media: dto.media
    }

    products.value.unshift(newProduct)

    if (shop.value && shop.value.name === shopName) {
      shop.value.products.unshift(newProduct)
    }

    return { ...newProduct }
  }

  async updateProduct(id: number | string, dto: UpdateProductDTO): Promise<Product | null> {
    const { products, shop } = useMockDataStore()
    const strId = String(id)
    const index = products.value.findIndex(p => String(p.id) === strId)
    const current = products.value[index]

    if (index === -1 || !current) return null

    const updated: Product = {
      ...current,
      ...dto
    }

    products.value[index] = updated

    if (shop.value) {
      const shopItemIndex = shop.value.products.findIndex(p => String(p.id) === strId)
      if (shopItemIndex !== -1) {
        shop.value.products[shopItemIndex] = { ...updated }
      }
    }

    return { ...updated }
  }

  async deleteProduct(id: number | string): Promise<boolean> {
    const { products, cart, shop } = useMockDataStore()
    const strId = String(id)
    const initialLen = products.value.length

    products.value = products.value.filter(p => String(p.id) !== strId)
    cart.value = cart.value.filter(item => String(item.productId) !== strId)

    if (shop.value) {
      shop.value.products = shop.value.products.filter(p => String(p.id) !== strId)
    }

    return products.value.length < initialLen
  }

  async toggleProductStatus(id: number | string): Promise<Product | null> {
    const product = await this.getProductById(id)
    if (!product) return null

    const nextStatus = product.status === 'Active' ? 'Draft' : 'Active'
    return this.updateProduct(id, { status: nextStatus })
  }
}
