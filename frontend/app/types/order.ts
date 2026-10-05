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

export interface ChapaCheckoutRequest {
  payment_method: 'cod' | 'online'
  recipient_name: string
  phone: string
  email: string
  delivery_address: string
  delivery_city: string
  delivery_notes: string
}

export interface ChapaCheckoutResponse {
  order: {
    ID: string
  }
  items: Array<{
    ID: string
    OrderID: string
    ProductID: string
    Quantity: number
    Price: string
    CreatedAt: string
  }>
  payment_id: string
  payment_method: string
  transaction_id: string
  checkout_url: string
}
