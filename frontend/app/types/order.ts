import type { Product } from './product'

export type OrderStatus = 'Ordered' | 'Shipped' | 'Delivered'
export type PaymentStatus = 'Pending' | 'Paid' | 'Failed' | 'Cash on delivery'

export interface CartItem {
  productId: number
  quantity: number
}

export interface CartProductItem extends CartItem {
  product: Product
  lineTotal: number
}

export interface MarketplaceOrder {
  id: number
  buyerName: string
  items: CartItem[]
  deliveryAddress: string
  phone: string
  paymentMethod: string
  paymentStatus: PaymentStatus
  status: OrderStatus
  date: string
  total: number
}

export interface CreateOrderDTO {
  buyerName: string
  deliveryAddress: string
  phone: string
  paymentMethod: string
  paymentStatus?: PaymentStatus
  deliveryFee?: number
}
