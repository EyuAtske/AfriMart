import type { IShopRepository } from '../interfaces/IShopRepository'
import type { Shop, CreateShopDTO, UpdateShopDTO } from '~/types/shop'
import { useMockDataStore } from './MockDataStore'

export class MockShopRepository implements IShopRepository {
  private slugify(value: string): string {
    return value.toLowerCase().replace(/\s+/g, '-')
  }

  async getShopBySlug(slug: string): Promise<Shop | null> {
    const { shop, products } = useMockDataStore()
    const targetSlug = slug.toLowerCase()

    if (shop.value && this.slugify(shop.value.name) === targetSlug) {
      return { ...shop.value }
    }

    const matchingProducts = products.value.filter(
      p => this.slugify(p.shop) === targetSlug
    )

    const firstMatch = matchingProducts[0]
    if (firstMatch) {
      const realShopName = firstMatch.shop
      return {
        id: targetSlug,
        name: realShopName,
        slug: targetSlug,
        description: `Independent verified boutique featuring curated pieces by ${realShopName}.`,
        ownerEmail: 'seller@afrimart.com',
        products: matchingProducts,
        paymentMethods: [
          { id: 1, type: 'Telebirr', accountName: realShopName, accountNumber: '0911223344' }
        ]
      }
    }

    return null
  }

  async getShopByOwner(email: string): Promise<Shop | null> {
    const { shop } = useMockDataStore()
    if (shop.value && shop.value.ownerEmail === email) {
      return { ...shop.value }
    }
    return shop.value ? { ...shop.value } : null
  }

  async createShop(ownerEmail: string, dto: CreateShopDTO): Promise<Shop> {
    const { shop, user } = useMockDataStore()
    const name = dto.name.trim()
    const slug = this.slugify(name)

    const newShop: Shop = {
      id: slug,
      name,
      slug,
      description: dto.description.trim(),
      ownerEmail,
      products: [],
      paymentMethods: []
    }

    shop.value = newShop
    user.value.role = 'seller'

    return { ...newShop }
  }

  async updateShop(slug: string, dto: UpdateShopDTO): Promise<Shop | null> {
    const { shop, products } = useMockDataStore()
    if (!shop.value) return null

    const oldName = shop.value.name
    if (dto.name) {
      const newName = dto.name.trim()
      shop.value.name = newName
      shop.value.slug = this.slugify(newName)

      products.value.forEach((p) => {
        if (p.shop === oldName) {
          p.shop = newName
        }
      })
    }

    if (dto.description) {
      shop.value.description = dto.description.trim()
    }

    return { ...shop.value }
  }
}
