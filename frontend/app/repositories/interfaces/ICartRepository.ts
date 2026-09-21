import type { CartItem, CartProductItem } from '~/types/order'

export interface BackendCartItem {
  id: string
  cart_id: string
  product_id: string
  quantity: number
  name?: string
  price?: string | number
  image?: { String: string; Valid: boolean } | string | null
  stock?: number
}

export interface BackendCartResponse {
  id: string
  user_id: string
  items: BackendCartItem[]
  subtotal: number
}

export interface ICartRepository {
  getCart(): Promise<BackendCartResponse | null>
  addItem(productId: string, quantity: number): Promise<BackendCartItem>
  updateItemQuantity(itemId: string, quantity: number): Promise<BackendCartItem>
  removeItem(itemId: string): Promise<boolean>
  clearCart(): Promise<boolean>
}
