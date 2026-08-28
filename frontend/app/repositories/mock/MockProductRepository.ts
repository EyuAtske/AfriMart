import type { IProductRepository } from '../interfaces/IProductRepository'
import type { Product, ProductFilterParams, CreateProductDTO, UpdateProductDTO } from '~/types/product'
import type { PaginatedResponse } from '~/types/api'
import { useMockDataStore } from './MockDataStore'

export class MockProductRepository implements IProductRepository {
  async getProducts(params: ProductFilterParams = {}): Promise<PaginatedResponse<Product>> {
    const { products } = useMockDataStore()
    const search = params.search?.trim().toLowerCase() || ''
    const category = params.category || 'All'
    const shop = params.shop?.trim().toLowerCase()
    const status = params.status
    const page = params.page || 1
    const pageSize = params.pageSize || 50

    const filtered = products.value.filter((p) => {
      const matchesSearch = !search ||
        p.name.toLowerCase().includes(search) ||
        p.shop.toLowerCase().includes(search) ||
        p.description.toLowerCase().includes(search)

      const matchesCategory = category === 'All' || p.category === category
      const matchesShop = !shop || p.shop.toLowerCase() === shop
      const matchesStatus = !status ? p.status === 'Active' : p.status === status

      return matchesSearch && matchesCategory && matchesShop && matchesStatus
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

  async getProductById(id: number): Promise<Product | null> {
    const { products } = useMockDataStore()
    const found = products.value.find(p => p.id === id)
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
      price: dto.price,
      stock: dto.stock,
      rating: 'New',
      image: dto.image,
      status: dto.status || 'Active'
    }

    products.value.unshift(newProduct)

    if (shop.value && shop.value.name === shopName) {
      shop.value.products.unshift(newProduct)
    }

    return { ...newProduct }
  }

  async updateProduct(id: number, dto: UpdateProductDTO): Promise<Product | null> {
    const { products, shop } = useMockDataStore()
    const index = products.value.findIndex(p => p.id === id)
    const current = products.value[index]

    if (index === -1 || !current) return null

    const updated: Product = {
      ...current,
      ...dto
    }

    products.value[index] = updated

    if (shop.value) {
      const shopItemIndex = shop.value.products.findIndex(p => p.id === id)
      if (shopItemIndex !== -1) {
        shop.value.products[shopItemIndex] = { ...updated }
      }
    }

    return { ...updated }
  }

  async deleteProduct(id: number): Promise<boolean> {
    const { products, cart, shop } = useMockDataStore()
    const initialLen = products.value.length

    products.value = products.value.filter(p => p.id !== id)
    cart.value = cart.value.filter(item => item.productId !== id)

    if (shop.value) {
      shop.value.products = shop.value.products.filter(p => p.id !== id)
    }

    return products.value.length < initialLen
  }

  async toggleProductStatus(id: number): Promise<Product | null> {
    const { products } = useMockDataStore()
    const product = products.value.find(p => p.id === id)
    if (!product) return null

    const nextStatus = product.status === 'Active' ? 'Draft' : 'Active'
    return this.updateProduct(id, { status: nextStatus })
  }
}
