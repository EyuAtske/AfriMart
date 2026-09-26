import type { IShopRepository } from '../interfaces/IShopRepository'
import type { Shop, CreateShopDTO, UpdateShopDTO } from '~/types/shop'
import { useMockDataStore } from '../mock/MockDataStore'
import { authenticatedFetch, extractError, getAccessToken } from './apiHelpers'

/**
 * Raw shape returned by the Go backend for a Shop.
 * Go's sql.NullString serializes as { String, Valid }.
 * Fields use Go's default PascalCase since sqlc has no JSON tags configured.
 */
interface BackendShopResponse {
  ID: string
  OwnerID: string
  Name: string
  Description: { String: string; Valid: boolean }
  Status: string
  CreatedAt: string
  UpdatedAt: string
}

/**
 * Maps the raw backend shop response to the frontend Shop type.
 */
function mapBackendShop(raw: BackendShopResponse): Shop {
  return {
    id: raw.ID,
    name: raw.Name,
    slug: raw.Name.toLowerCase().replace(/\s+/g, '-'),
    description: raw.Description?.Valid ? raw.Description.String : '',
    ownerEmail: '',  // Backend doesn't return owner email, only OwnerID
    products: [],
    paymentMethods: [],
    status: raw.Status,
    backendId: raw.ID
  }
}

export class ApiShopRepository implements IShopRepository {

  async getMyShop(): Promise<Shop | null> {
    const token = getAccessToken()
    if (!token) {
      return null
    }

    try {
      const res = await authenticatedFetch<BackendShopResponse>('api/shops/me', {
        method: 'GET'
      })
      if (!res) return null
      return mapBackendShop(res)
    } catch (err: any) {
      const status = err?.response?.status || err?.statusCode || err?.status
      // 404 means user has no shop — not an error
      if (status === 404) return null
      if (status === 401) {
        throw new Error('Please log in to access your seller shop.')
      }
      throw new Error(extractError(err, 'Failed to fetch your shop'))
    }
  }

  async createShop(_ownerEmail: string, dto: CreateShopDTO): Promise<Shop> {
    const { shop, user } = useMockDataStore()

    try {
      const res = await authenticatedFetch<BackendShopResponse>('api/shops', {
        method: 'POST',
        body: {
          name: dto.name.trim(),
          description: dto.description.trim()
        }
      })

      const created = mapBackendShop(res)
      shop.value = created
      user.value.role = 'seller'
      return created
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to create shop'))
    }
  }

  async updateShop(slug: string, dto: UpdateShopDTO): Promise<Shop | null> {
    const { shop } = useMockDataStore()
    if (!shop.value) return null

    const shopId = shop.value.backendId || shop.value.id
    let updated: BackendShopResponse | null = null

    try {
      // Backend has separate endpoints for name and description
      if (dto.name) {
        updated = await authenticatedFetch<BackendShopResponse>(`api/shops/${shopId}/name`, {
          method: 'PATCH',
          body: { name: dto.name.trim() }
        })
      }

      if (dto.description !== undefined) {
        updated = await authenticatedFetch<BackendShopResponse>(`api/shops/${shopId}/description`, {
          method: 'PATCH',
          body: { description: dto.description.trim() }
        })
      }

      if (updated) {
        const mappedShop = mapBackendShop(updated)
        shop.value = mappedShop
        return mappedShop
      }

      return shop.value ? { ...shop.value } : null
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to update shop'))
    }
  }

  async deactivateShop(shopId: string): Promise<Shop | null> {
    const { shop } = useMockDataStore()

    try {
      const res = await authenticatedFetch<BackendShopResponse>(`api/shops/${shopId}/deactivate`, {
        method: 'PATCH'
      })

      const updated = mapBackendShop(res)
      if (shop.value) {
        shop.value = updated
      }
      return updated
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to deactivate shop'))
    }
  }

  async activateShop(shopId: string): Promise<Shop | null> {
    const { shop } = useMockDataStore()

    try {
      const res = await authenticatedFetch<BackendShopResponse>(`api/shops/${shopId}/activate`, {
        method: 'PATCH'
      })

      const updated = mapBackendShop(res)
      if (shop.value) {
        shop.value = updated
      }
      return updated
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to activate shop'))
    }
  }

  // Fallback methods — these don't have backend equivalents yet
  async getShopBySlug(_slug: string): Promise<Shop | null> {
    // No public shop-by-slug endpoint on the backend
    // Fall through to the mock for now
    return null
  }

  async getShopByOwner(_email: string): Promise<Shop | null> {
    // Equivalent to getMyShop when authenticated
    return this.getMyShop()
  }
}
