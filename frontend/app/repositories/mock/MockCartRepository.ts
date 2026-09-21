import type { ICartRepository, BackendCartResponse, BackendCartItem } from '../interfaces/ICartRepository'
import { useMockDataStore } from './MockDataStore'

export class MockCartRepository implements ICartRepository {
  async getCart(): Promise<BackendCartResponse | null> {
    const { cart, products } = useMockDataStore()
    const items: BackendCartItem[] = cart.value.map(item => {
      const p = products.value.find(prod => Number(prod.id) === Number(item.productId))
      return {
        id: `mock-item-${item.productId}`,
        cart_id: 'mock-cart-id',
        product_id: String(item.productId),
        quantity: item.quantity,
        name: p?.name || 'Mock Product',
        price: p?.price || 0,
        image: p?.image || '',
        stock: p?.stock || 10
      }
    })

    const subtotal = items.reduce((acc, item) => acc + Number(item.price) * item.quantity, 0)

    return {
      id: 'mock-cart-id',
      user_id: 'mock-user-id',
      items,
      subtotal
    }
  }

  async addItem(productId: string, quantity: number = 1): Promise<BackendCartItem> {
    const { cart, products } = useMockDataStore()
    const numId = Number(productId)
    const existing = cart.value.find(i => Number(i.productId) === numId)
    if (existing) {
      existing.quantity += quantity
    } else {
      cart.value.push({ productId: numId, quantity })
    }

    const p = products.value.find(prod => Number(prod.id) === numId)
    return {
      id: `mock-item-${productId}`,
      cart_id: 'mock-cart-id',
      product_id: productId,
      quantity: existing ? existing.quantity : quantity,
      name: p?.name || 'Mock Product',
      price: p?.price || 0,
      image: p?.image || '',
      stock: p?.stock || 10
    }
  }

  async updateItemQuantity(itemId: string, quantity: number): Promise<BackendCartItem> {
    const { cart, products } = useMockDataStore()
    const prodId = itemId.replace('mock-item-', '')
    const numId = Number(prodId)

    const existing = cart.value.find(i => Number(i.productId) === numId)
    if (existing) {
      existing.quantity = quantity
    }

    const p = products.value.find(prod => Number(prod.id) === numId)
    return {
      id: itemId,
      cart_id: 'mock-cart-id',
      product_id: prodId,
      quantity,
      name: p?.name || 'Mock Product',
      price: p?.price || 0,
      image: p?.image || '',
      stock: p?.stock || 10
    }
  }

  async removeItem(itemId: string): Promise<boolean> {
    const { cart } = useMockDataStore()
    const prodId = itemId.replace('mock-item-', '')
    const numId = Number(prodId)
    cart.value = cart.value.filter(i => Number(i.productId) !== numId)
    return true
  }

  async clearCart(): Promise<boolean> {
    const { cart } = useMockDataStore()
    cart.value = []
    return true
  }
}
