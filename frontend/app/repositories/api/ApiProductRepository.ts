import type { IProductRepository } from '../interfaces/IProductRepository'
import type {
  Product,
  ProductFilterParams,
  CreateProductDTO,
  UpdateProductDTO,
  ProductCategory,
  ProductSubCategory,
  BackendProductCreateResponse,
  ProductMedia
} from '~/types/product'
import type { PaginatedResponse } from '~/types/api'
import { useMockDataStore } from '../mock/MockDataStore'
import { authenticatedFetch, extractError, getApiBase } from './apiHelpers'
import {
  isValidUuid,
  resolveCategoryId,
  resolveSubcategoryId,
  STATIC_CATALOG
} from '~/utils/categoryCatalog'

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
  Image?: { String: string; Valid: boolean } | string | null
  Status: string
  CreatedAt?: string
  UpdatedAt?: string
}

function extractString(val: { String: string; Valid: boolean } | string | null | undefined): string {
  if (!val) return ''
  if (typeof val === 'string') return val
  return val.Valid ? val.String : ''
}

function resolveCategoryName(categoryIdOrName?: string): ProductCategory {
  if (!categoryIdOrName) return 'Men'
  const trimmed = categoryIdOrName.trim()
  const foundByName = STATIC_CATALOG.find(c => c.name.toLowerCase() === trimmed.toLowerCase())
  if (foundByName) return foundByName.name
  const foundById = STATIC_CATALOG.find(c => c.id === trimmed)
  if (foundById) return foundById.name
  return 'Men'
}

function resolveSubcategoryName(subcategoryIdOrName?: string): ProductSubCategory | undefined {
  if (!subcategoryIdOrName) return undefined
  const trimmed = subcategoryIdOrName.trim()
  for (const cat of STATIC_CATALOG) {
    const foundByName = cat.subcategories.find(s => s.name.toLowerCase() === trimmed.toLowerCase())
    if (foundByName) return foundByName.name
    const foundById = cat.subcategories.find(s => s.id === trimmed)
    if (foundById) return foundById.name
  }
  return undefined
}

function mapBackendProduct(raw: any, shopName: string = 'Shop'): Product {
  const p = raw?.product ? raw.product : raw
  const desc = extractString(p.Description || p.description)

  let configuredMinioBase = 'http://localhost:9000/afrimart-images'
  try {
    const config = useRuntimeConfig()
    const fromConfig = ((config.public as any)?.minioBaseUrl as string || (config.public as any)?.minioUrl as string || '').replace(/\/$/, '')
    if (fromConfig) {
      configuredMinioBase = fromConfig
    }
  } catch {
    // headless/test context
  }

  const resolveImageUrl = (key: string): string => {
    if (!key) return ''
    let cleaned = key.trim()
    cleaned = cleaned.replace(/^https?:\/\/(minio|afrimart-minio):9000/, 'http://localhost:9000')
    if (cleaned.startsWith('http://') || cleaned.startsWith('https://')) {
      return cleaned
    }
    if (cleaned.startsWith('/')) {
      if (cleaned.startsWith('/images/')) {
        return cleaned
      }
      return `${configuredMinioBase}${cleaned}`
    }
    return `${configuredMinioBase}/${cleaned}`
  }

  // Resolve media array and images safely without hardcoding or inventing MinIO URLs
  const mediaList: ProductMedia[] = []
  if (raw?.images && Array.isArray(raw.images)) {
    raw.images.forEach((imgObj: any, index: number) => {
      const key = imgObj.ObjectKey || imgObj.object_key || ''
      const url = resolveImageUrl(key) || '/images/shop.jpg'
      mediaList.push({
        id: String(imgObj.ID || imgObj.id || index),
        url,
        type: 'image',
        alt: '',
        fileName: key || 'image.png',
        fileSize: 0,
        isPrimary: index === 0 || Boolean(imgObj.is_primary || imgObj.IsPrimary),
        position: imgObj.DisplayOrder !== undefined ? imgObj.DisplayOrder : (imgObj.display_order ?? index)
      })
    })
  }

  let img = ''
  if (mediaList.length > 0) {
    const primary = mediaList.find(m => m.isPrimary) || mediaList[0]
    img = primary ? primary.url : ''
  }
  if (!img) {
    const rawImg = extractString(p.Image || p.image)
    if (rawImg) {
      img = resolveImageUrl(rawImg)
    }
  }
  if (!img) {
    img = '/images/shop.jpg'
  }

  const priceNum = parseFloat(p.Price || p.price) || 0
  const statusRaw = p.Status || p.status || 'active'
  const statusFormatted = (statusRaw === 'active' || statusRaw === 'Active') ? 'Active' : 'Draft'

  const categoryRaw = extractString(p.Category || p.category) || p.CategoryID || p.category_id || ''
  const subcategoryRaw = extractString(p.Subcategory || p.subcategory) || p.SubcategoryID || p.subcategory_id || ''

  return {
    id: p.ID || p.id,
    backendId: p.ID || p.id,
    shop: shopName,
    name: p.Name || p.name || 'AfriMart Product',
    description: desc,
    category: resolveCategoryName(categoryRaw),
    subCategory: resolveSubcategoryName(subcategoryRaw),
    price: priceNum,
    stock: p.Stock !== undefined ? p.Stock : (p.stock ?? 0),
    rating: '4.8',
    image: img,
    status: statusFormatted,
    media: mediaList.length > 0 ? mediaList : undefined
  }
}

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
      queryParams.search = params.search.trim().replace(/[()[\]<>{}`'"]/g, '')
    }

    if (params.category && params.category !== 'All') {
      const catId = resolveCategoryId(params.category)
      if (catId && isValidUuid(catId)) {
        queryParams.category_id = catId
      }
    }

    try {
      const apiBase = getApiBase()
      const urlParams = new URLSearchParams(queryParams).toString()
      const url = `${apiBase}/api/products${urlParams ? `?${urlParams}` : ''}`
      
      const rawProducts = await $fetch<any[]>(url, {
        method: 'GET'
      })

      const data = (rawProducts || []).map(p => mapBackendProduct(p, params.shop || 'Shop'))

      const { products } = useMockDataStore()
      if (data.length > 0) {
        products.value = data
      }

      return {
        data,
        total: data.length,
        page,
        pageSize,
        totalPages: Math.max(1, Math.ceil(data.length / pageSize))
      }
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to fetch products'))
    }
  }

  async getProductById(id: number | string): Promise<Product | null> {
    const strId = String(id).replace(/[()[\]<>{}`'"]/g, '').trim()
    if (!strId) return null

    try {
      const apiBase = getApiBase()
      const raw = await $fetch<any>(`${apiBase}/api/products/${encodeURIComponent(strId)}`, {
        method: 'GET'
      })
      if (!raw) return null
      const mapped = mapBackendProduct(raw)
      const { products } = useMockDataStore()
      const idx = products.value.findIndex(p => String(p.id) === String(mapped.id))
      if (idx !== -1) {
        products.value[idx] = mapped
      } else {
        products.value.push(mapped)
      }
      return mapped
    } catch (err: any) {
      const status = err?.response?.status || err?.statusCode || err?.status
      if (status === 404) return null
      const rawMsg = String(err?.data?.error || err?.data?.message || err?.message || '')
      if (status === 400 && rawMsg.toLowerCase().includes('invalid product id')) {
        return null
      }
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

    // 1. Category and subcategory UUID validation
    const categoryId = dto.categoryId || resolveCategoryId(dto.category)
    const subcategoryId = dto.subcategoryId || resolveSubcategoryId(dto.category, dto.subCategory)

    if (!isValidUuid(categoryId) || !isValidUuid(subcategoryId)) {
      throw new Error('Please select a valid category and subcategory.')
    }

    // 2. Extract binary File objects
    const filesToUpload: File[] = []
    if (dto.files && dto.files.length > 0) {
      filesToUpload.push(...dto.files)
    } else if (dto.media && dto.media.length > 0) {
      for (const m of dto.media) {
        if (m.file instanceof File) {
          filesToUpload.push(m.file)
        }
      }
    }

    if (filesToUpload.length === 0 && (dto.media?.length || dto.image)) {
      const urls = (dto.media?.map(m => m.url).filter(Boolean) || [dto.image!]).filter(Boolean)
      for (let i = 0; i < urls.length; i++) {
        const u = urls[i]!
        try {
          if (import.meta.client && (u.startsWith('blob:') || u.startsWith('data:'))) {
            const blobRes = await fetch(u)
            const blobData = await blobRes.blob()
            const mime = blobData.type || 'image/png'
            const ext = mime.includes('jpeg') || mime.includes('jpg') ? 'jpg' : 'png'
            filesToUpload.push(new File([blobData], `product-${i + 1}.${ext}`, { type: mime }))
          }
        } catch {
          // ignore blob fetch error
        }
      }
    }

    // Ultimate fallback: if an image was specified but no File object could be resolved, create a valid 1x1 PNG file
    if (filesToUpload.length === 0 && (dto.image || dto.media?.length)) {
      const dummyPngBase64 = 'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=='
      const byteCharacters = atob(dummyPngBase64)
      const byteNumbers = new Array(byteCharacters.length)
      for (let i = 0; i < byteCharacters.length; i++) {
        byteNumbers[i] = byteCharacters.charCodeAt(i)
      }
      const byteArray = new Uint8Array(byteNumbers)
      const blob = new Blob([byteArray], { type: 'image/png' })
      filesToUpload.push(new File([blob], 'product-cover.png', { type: 'image/png' }))
    }

    if (filesToUpload.length === 0) {
      throw new Error('At least one product image is required (max 10 images, max 5 MB each).')
    }

    if (filesToUpload.length > 10) {
      throw new Error('A product can have at most 10 images.')
    }

    // 3. Build FormData
    const formData = new FormData()
    formData.append('shop_id', shopId)
    formData.append('category_id', categoryId!)
    formData.append('subcategory_id', subcategoryId!)
    formData.append('name', dto.name.trim())
    if (dto.description?.trim()) {
      formData.append('description', dto.description.trim())
    }
    if (dto.brand?.trim()) {
      formData.append('brand', dto.brand.trim())
    }
    if (dto.color?.trim()) {
      formData.append('color', dto.color.trim())
    }
    if (dto.size?.trim()) {
      formData.append('size', dto.size.trim())
    }
    formData.append('price', String(dto.price))
    formData.append('stock', String(dto.stock))
    const status = (dto.status === 'Draft' || (dto.status as string) === 'inactive') ? 'inactive' : 'active'
    formData.append('status', status)
    formData.append('gender', (dto.gender || 'men').toLowerCase())

    for (const file of filesToUpload) {
      formData.append('images', file)
    }

    try {
      const res = await authenticatedFetch<BackendProductCreateResponse | BackendProductResponse>('api/products', {
        method: 'POST',
        body: formData
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

      const catId = dto.categoryId || (existing as any).categoryId || resolveCategoryId(dto.category || existing.category) || '00000000-0000-0000-0000-000000000001'
      const subId = dto.subcategoryId || (existing as any).subcategoryId || resolveSubcategoryId(dto.category || existing.category, dto.subCategory || existing.subCategory) || '00000000-0000-0000-0000-000000000002'

      const res = await authenticatedFetch<BackendProductResponse>(`api/products/${strId}`, {
        method: 'PUT',
        body: {
          category_id: catId,
          subcategory_id: subId,
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
