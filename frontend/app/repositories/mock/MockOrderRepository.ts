import type { IOrderRepository } from '../interfaces/IOrderRepository'
import type { MarketplaceOrder, OrderStatus, CreateOrderDTO, CartItem } from '~/types/order'
import { useMockDataStore } from './MockDataStore'

export class MockOrderRepository implements IOrderRepository {
  async getOrders(): Promise<MarketplaceOrder[]> {
    const { orders } = useMockDataStore()
    return [...orders.value]
  }

  async getOrderById(id: number): Promise<MarketplaceOrder | null> {
    const { orders } = useMockDataStore()
    const found = orders.value.find(o => o.id === id)
    return found ? { ...found } : null
  }

  async createOrder(cartItems: CartItem[], cartSubtotal: number, dto: CreateOrderDTO): Promise<MarketplaceOrder> {
    const { orders, cart } = useMockDataStore()
    const id = Date.now()

    const order: MarketplaceOrder = {
      id,
      buyerName: dto.buyerName,
      items: [...cartItems],
      deliveryAddress: dto.deliveryAddress,
      phone: dto.phone,
      paymentMethod: dto.paymentMethod,
      paymentStatus: dto.paymentStatus || (dto.paymentMethod === 'Cash on delivery' ? 'Cash on delivery' : 'Paid'),
      status: 'Ordered',
      date: new Date().toLocaleDateString('en-US', {
        month: 'long',
        day: 'numeric',
        year: 'numeric'
      }),
      total: cartSubtotal + (dto.deliveryFee || 0)
    }

    orders.value.unshift(order)
    cart.value = []

    return { ...order }
  }

  async updateOrderStatus(orderId: number, status: OrderStatus): Promise<MarketplaceOrder | null> {
    const { orders } = useMockDataStore()
    const order = orders.value.find(o => o.id === orderId)
    if (order) {
      order.status = status
      return { ...order }
    }
    return null
  }
}
