import type { MarketplaceOrder, OrderStatus, CreateOrderDTO, CartItem } from '~/types/order'

export interface IOrderRepository {
  getOrders(): Promise<MarketplaceOrder[]>
  getSellerOrders(): Promise<MarketplaceOrder[]>
  getOrderById(id: number | string): Promise<MarketplaceOrder | null>
  createOrder(cartItems: CartItem[], cartSubtotal: number, dto: CreateOrderDTO): Promise<MarketplaceOrder>
  updateOrderStatus(orderId: number | string, status: OrderStatus): Promise<MarketplaceOrder | null>
}
