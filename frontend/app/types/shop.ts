import type { Product } from './product'

export interface PaymentMethod {
  id: number
  type: 'Telebirr' | 'CBE'
  accountName: string
  accountNumber: string
}

export interface Shop {
  id: string
  name: string
  slug: string
  description: string
  ownerEmail: string
  products: Product[]
  paymentMethods: PaymentMethod[]
}

export interface CreateShopDTO {
  name: string
  description: string
}

export interface UpdateShopDTO {
  name?: string
  description?: string
}
