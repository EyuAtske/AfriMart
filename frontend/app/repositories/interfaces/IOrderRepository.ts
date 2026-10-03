import type { MarketplaceOrder, OrderStatus, CreateOrderDTO, CartItem, ChapaCheckoutRequest, ChapaCheckoutResponse } from '~/types/order'

export interface IOrderRepository {
  getOrders(): Promise<MarketplaceOrder[]>
  getSellerOrders(): Promise<MarketplaceOrder[]>
  getOrderById(id: number | string): Promise<MarketplaceOrder | null>
  createOrder(cartItems: CartItem[], cartSubtotal: number, dto: CreateOrderDTO): Promise<MarketplaceOrder>
  initiateChapaCheckout(request: ChapaCheckoutRequest): Promise<ChapaCheckoutResponse>
  updateOrderStatus(orderId: number | string, status: OrderStatus): Promise<MarketplaceOrder | null>
}
