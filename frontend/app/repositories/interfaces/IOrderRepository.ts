import type { MarketplaceOrder, OrderStatus, CreateOrderDTO, CartItem } from '~/types/order'

export interface IOrderRepository {
  getOrders(): Promise<MarketplaceOrder[]>
  getOrderById(id: number): Promise<MarketplaceOrder | null>
  createOrder(cartItems: CartItem[], cartSubtotal: number, dto: CreateOrderDTO): Promise<MarketplaceOrder>
  updateOrderStatus(orderId: number, status: OrderStatus): Promise<MarketplaceOrder | null>
}
