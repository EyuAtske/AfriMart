import type { ICartRepository, BackendCartResponse, BackendCartItem } from '../interfaces/ICartRepository'
import { authenticatedFetch, extractError, getAccessToken } from './apiHelpers'

export class ApiCartRepository implements ICartRepository {
  /**
   * Fetch user's active cart from GET /api/cart
   */
  async getCart(): Promise<BackendCartResponse | null> {
    const token = getAccessToken()
    if (!token) return null

    try {
      const res = await authenticatedFetch<BackendCartResponse>('api/cart', {
        method: 'GET'
      })
      return res || null
    } catch (err: any) {
      const status = err?.response?.status || err?.statusCode || err?.status
      if (status === 404) {
        return null
      }
      throw new Error(extractError(err, 'Failed to fetch cart from server'))
    }
  }

  /**
   * Add a product item to user's cart via POST /api/cart/items
   */
  async addItem(productId: string, quantity: number = 1): Promise<BackendCartItem> {
    const token = getAccessToken()
    if (!token) {
      throw new Error('User must be logged in to add items to cart.')
    }

    try {
      const res = await authenticatedFetch<BackendCartItem>('api/cart/items', {
        method: 'POST',
        body: {
          product_id: String(productId).trim(),
          quantity: Math.max(1, Math.floor(quantity))
        }
      })
      return res
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to add item to cart'))
    }
  }

  /**
   * Update quantity of a cart item via PUT /api/cart/items/{id}
   */
  async updateItemQuantity(itemId: string, quantity: number): Promise<BackendCartItem> {
    const token = getAccessToken()
    if (!token) {
      throw new Error('User must be logged in to update cart items.')
    }

    try {
      const res = await authenticatedFetch<BackendCartItem>(`api/cart/items/${itemId}`, {
        method: 'PATCH',
        body: {
          quantity: Math.max(1, Math.floor(quantity))
        }
      })
      return res
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to update cart item quantity'))
    }
  }

  /**
   * Remove single item from cart via DELETE /api/cart/items/{id}
   */
  async removeItem(itemId: string): Promise<boolean> {
    const token = getAccessToken()
    if (!token) {
      throw new Error('User must be logged in to remove cart items.')
    }

    try {
      await authenticatedFetch(`api/cart/items/${itemId}`, {
        method: 'DELETE'
      })
      return true
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to remove item from cart'))
    }
  }

  /**
   * Clear all cart items via DELETE /api/cart
   */
  async clearCart(): Promise<boolean> {
    const token = getAccessToken()
    if (!token) return false

    try {
      await authenticatedFetch('api/cart', {
        method: 'DELETE'
      })
      return true
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to clear cart'))
    }
  }
}
