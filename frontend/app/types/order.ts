import type { Product } from './product'

export type OrderStatus = 'Pending' | 'Confirmed' | 'Processing' | 'Shipped' | 'Delivered' | 'Cancelled' | 'Ordered'
export type PaymentStatus = 'Pending' | 'Paid' | 'Failed' | 'Cash on delivery'

export interface CartItem {
  productId: number | string
  quantity: number
  backendItemId?: string
}

export interface CartProductItem extends CartItem {
  product: Product
  lineTotal: number
}

export interface MarketplaceOrder {
  id: number | string
  backendId?: string
  buyerName: string
  items: CartItem[]
  deliveryAddress: string
  deliveryCity?: string
  deliveryNotes?: string
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
  deliveryCity: string
  deliveryNotes?: string
  phone: string
  paymentMethod?: string
  paymentStatus?: PaymentStatus
  deliveryFee?: number
}
