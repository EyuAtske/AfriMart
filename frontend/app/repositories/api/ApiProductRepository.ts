import type { IProductRepository } from '../interfaces/IProductRepository'
import type { Product, ProductFilterParams, CreateProductDTO, UpdateProductDTO, ProductCategory } from '~/types/product'
import type { PaginatedResponse } from '~/types/api'
import { useMockDataStore } from '../mock/MockDataStore'
import { authenticatedFetch, extractError, getApiBase } from './apiHelpers'

/**
 * Shape returned by Go backend for a Product.
 * sqlc serializes NullString as { String, Valid }.
 */
interface BackendProductResponse {
  ID: string
  ShopID: string
  CategoryID: string
  SubcategoryID: string
  Name: string
  Description: { String: string; Valid: boolean } | string | null
  Brand?: { String: string; Valid: boolean } | string | null
  Color?: { String: string; Valid: boolean } | string | null
  Size?: { String: string; Valid: boolean } | string | null
  Price: string
  Stock: number
  Image: { String: string; Valid: boolean } | string | null
  Status: string
  CreatedAt?: string
  UpdatedAt?: string
}

function extractString(val: { String: string; Valid: boolean } | string | null | undefined): string {
  if (!val) return ''
  if (typeof val === 'string') return val
  return val.Valid ? val.String : ''
}

function mapBackendProduct(raw: BackendProductResponse, shopName: string = 'Shop'): Product {
  const desc = extractString(raw.Description)
  const img = extractString(raw.Image) || '/images/herotemp.png'
  const priceNum = parseFloat(raw.Price) || 0
  const statusFormatted = (raw.Status === 'active' || raw.Status === 'Active') ? 'Active' : 'Draft'

  return {
    id: raw.ID,
    backendId: raw.ID,
    shop: shopName,
    name: raw.Name,
    description: desc,
    category: 'Men' as ProductCategory,
    price: priceNum,
    stock: raw.Stock ?? 0,
    rating: '4.8',
    image: img,
    status: statusFormatted
  }
}

// Fallback UUIDs for required backend category & subcategory parameters
const DEFAULT_CATEGORY_ID = '00000000-0000-0000-0000-000000000001'
const DEFAULT_SUBCATEGORY_ID = '00000000-0000-0000-0000-000000000002'

export class ApiProductRepository implements IProductRepository {

  async getProducts(params: ProductFilterParams = {}): Promise<PaginatedResponse<Product>> {
    const page = params.page || 1
    const pageSize = params.pageSize || 20
    const offset = (page - 1) * pageSize

    const queryParams: Record<string, string> = {
      limit: String(pageSize),
      offset: String(offset)
    }

    if (params.search?.trim()) {
      queryParams.search = params.search.trim()
    }

    try {
      const apiBase = getApiBase()
      const urlParams = new URLSearchParams(queryParams).toString()
      const url = `${apiBase}/api/products${urlParams ? `?${urlParams}` : ''}`
      
      const rawProducts = await $fetch<BackendProductResponse[]>(url, {
        method: 'GET'
      })

      const data = (rawProducts || []).map(p => mapBackendProduct(p, params.shop || 'Shop'))

      return {
        data,
        total: data.length,
        page,
        pageSize,
        totalPages: Math.max(1, Math.ceil(data.length / pageSize))
      }
    } catch (err: any) {
      // In API mode, return real API response shape (empty data if backend unreachable/error)
      return {
        data: [],
        total: 0,
        page,
        pageSize,
        totalPages: 1
      }
    }
  }

  async getProductById(id: number | string): Promise<Product | null> {
    const strId = String(id).trim()
    try {
      const apiBase = getApiBase()
      const raw = await $fetch<BackendProductResponse>(`${apiBase}/api/products/${encodeURIComponent(strId)}`, {
        method: 'GET'
      })
      if (!raw) return null
      return mapBackendProduct(raw)
    } catch (err: any) {
      const status = err?.response?.status || err?.statusCode || err?.status
      if (status === 404) return null
      throw new Error(extractError(err, 'Failed to fetch product details'))
    }
  }

  async createProduct(shopName: string, dto: CreateProductDTO): Promise<Product> {
    const { shop, products } = useMockDataStore()

    let shopId = shop.value?.backendId || shop.value?.id
    if (!shopId) {
      // Get or create shop via shopRepo if not set
      const apiShopRepo = new (await import('./ApiShopRepository')).ApiShopRepository()
      const myShop = await apiShopRepo.getMyShop()
      if (myShop) {
        shopId = myShop.backendId || myShop.id
      }
    }

    if (!shopId || typeof shopId !== 'string') {
      throw new Error('You must create a shop before creating products.')
    }

    try {
      const res = await authenticatedFetch<BackendProductResponse>('api/products', {
        method: 'POST',
        body: {
          shop_id: shopId,
          category_id: DEFAULT_CATEGORY_ID,
          subcategory_id: DEFAULT_SUBCATEGORY_ID,
          name: dto.name.trim(),
          description: dto.description.trim(),
          price: String(dto.price),
          stock: dto.stock,
          image: dto.image || '',
          status: dto.status === 'Active' ? 'active' : 'inactive'
        }
      })

      const created = mapBackendProduct(res, shopName)
      
      // Update local reactive store as well
      products.value.unshift(created)
      if (shop.value) {
        shop.value.products.unshift(created)
      }

      return created
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to create product'))
    }
  }

  async updateProduct(id: number | string, dto: UpdateProductDTO): Promise<Product | null> {
    const strId = String(id)
    const { products, shop } = useMockDataStore()

    try {
      const existing = await this.getProductById(strId)
      if (!existing) return null

      const res = await authenticatedFetch<BackendProductResponse>(`api/products/${strId}`, {
        method: 'PUT',
        body: {
          category_id: DEFAULT_CATEGORY_ID,
          subcategory_id: DEFAULT_SUBCATEGORY_ID,
          name: dto.name !== undefined ? dto.name.trim() : existing.name,
          description: dto.description !== undefined ? dto.description.trim() : existing.description,
          price: dto.price !== undefined ? String(dto.price) : String(existing.price),
          stock: dto.stock !== undefined ? dto.stock : existing.stock,
          image: dto.image !== undefined ? dto.image : existing.image,
          status: (dto.status || existing.status) === 'Active' ? 'active' : 'inactive'
        }
      })

      const updated = mapBackendProduct(res, existing.shop)

      // Sync local store
      const index = products.value.findIndex(p => String(p.id) === strId)
      if (index !== -1) {
        products.value[index] = updated
      }
      if (shop.value) {
        const shopIdx = shop.value.products.findIndex(p => String(p.id) === strId)
        if (shopIdx !== -1) {
          shop.value.products[shopIdx] = updated
        }
      }

      return updated
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to update product'))
    }
  }

  async deleteProduct(id: number | string): Promise<boolean> {
    const strId = String(id)
    const { products, shop } = useMockDataStore()

    try {
      await authenticatedFetch(`api/products/${strId}`, {
        method: 'DELETE'
      })

      products.value = products.value.filter(p => String(p.id) !== strId)
      if (shop.value) {
        shop.value.products = shop.value.products.filter(p => String(p.id) !== strId)
      }

      return true
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to delete product'))
    }
  }

  async toggleProductStatus(id: number | string): Promise<Product | null> {
    const product = await this.getProductById(id)
    if (!product) return null

    const nextStatus = product.status === 'Active' ? 'Draft' : 'Active'
    return this.updateProduct(id, { status: nextStatus })
  }
}
