import type { Shop, CreateShopDTO, UpdateShopDTO } from '~/types/shop'

export interface IShopRepository {
  getShopBySlug(slug: string): Promise<Shop | null>
  getShopByOwner(email: string): Promise<Shop | null>
  createShop(ownerEmail: string, dto: CreateShopDTO): Promise<Shop>
  updateShop(slug: string, dto: UpdateShopDTO): Promise<Shop | null>
}
