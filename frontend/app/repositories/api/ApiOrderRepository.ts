import type { IOrderRepository } from '../interfaces/IOrderRepository'
import type { MarketplaceOrder, OrderStatus, CreateOrderDTO, CartItem } from '~/types/order'
import { useMockDataStore } from '../mock/MockDataStore'
import { authenticatedFetch, extractError, getAccessToken } from './apiHelpers'

export interface BackendOrderResponse {
  id: string | number
  user_id?: string
  subtotal?: string | number
  status?: string
  recipient_name?: string
  phone?: string
  delivery_address?: string
  delivery_city?: string
  delivery_notes?: string
  created_at?: string
  items?: any[]
  order?: BackendOrderResponse
}

function mapBackendOrder(raw: any): MarketplaceOrder {
  const o = raw?.order ? raw.order : raw
  const rawItems = raw?.items || o?.items || []

  const strId = String(o.id || '')
  const numId = typeof o.id === 'number' ? o.id : (parseInt(strId.replace(/\D/g, ''), 10) || Date.now())

  const rawStatus = (o.status || 'pending').toLowerCase()
  let status: OrderStatus = 'Pending'
  if (rawStatus === 'confirmed') status = 'Confirmed'
  else if (rawStatus === 'processing') status = 'Processing'
  else if (rawStatus === 'shipped') status = 'Shipped'
  else if (rawStatus === 'delivered') status = 'Delivered'
  else if (rawStatus === 'cancelled') status = 'Cancelled'
  else if (rawStatus === 'ordered') status = 'Ordered'
  else status = 'Pending'

  const subtotalNum = typeof o.subtotal === 'number' ? o.subtotal : parseFloat(String(o.subtotal || 0)) || 0

  return {
    id: numId,
    backendId: strId,
    buyerName: o.recipient_name || o.buyer_name || o.buyerName || 'Valued Customer',
    items: (rawItems || []).map((item: any) => ({
      productId: item.product_id || item.ProductID || item.productId || 0,
      quantity: item.quantity || item.Quantity || 1
    })),
    deliveryAddress: o.delivery_address || o.deliveryAddress || '',
    deliveryCity: o.delivery_city || o.deliveryCity || 'Addis Ababa',
    deliveryNotes: o.delivery_notes || o.deliveryNotes || '',
    phone: o.phone || '',
    paymentMethod: 'Cash on delivery',
    paymentStatus: 'Pending',
    status,
    date: o.created_at ? o.created_at.split('T')[0] : new Date().toISOString().split('T')[0]!,
    total: subtotalNum
  }
}

export class ApiOrderRepository implements IOrderRepository {
  /**
   * Fetch buyer orders via GET /api/orders
   */
  async getOrders(): Promise<MarketplaceOrder[]> {
    const token = getAccessToken()
    if (!token) return []

    try {
      const res = await authenticatedFetch<any>('api/orders', {
        method: 'GET'
      })
      const list = Array.isArray(res) ? res : (res?.orders || [])
      return list.map(mapBackendOrder)
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to fetch user orders'))
    }
  }

  /**
   * Fetch seller orders via GET /api/orders/seller
   */
  async getSellerOrders(): Promise<MarketplaceOrder[]> {
    const token = getAccessToken()
    if (!token) return []

    try {
      const res = await authenticatedFetch<any>('api/orders/seller', {
        method: 'GET'
      })
      const list = Array.isArray(res) ? res : (res?.orders || [])
      return list.map(mapBackendOrder)
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to fetch seller orders'))
    }
  }

  /**
   * Fetch single order detail via GET /api/orders/{id}
   */
  async getOrderById(id: number | string): Promise<MarketplaceOrder | null> {
    const token = getAccessToken()
    if (!token) return null

    try {
      const res = await authenticatedFetch<any>(`api/orders/${id}`, {
        method: 'GET'
      })
      if (res) return mapBackendOrder(res)
      return null
    } catch (err: any) {
      const status = err?.response?.status || err?.statusCode || err?.status
      if (status === 404) return null
      throw new Error(extractError(err, 'Failed to fetch order details'))
    }
  }

  /**
   * Submit COD checkout & create order via POST /api/orders/checkout
   */
  async createOrder(cartItems: CartItem[], cartSubtotal: number, dto: CreateOrderDTO): Promise<MarketplaceOrder> {
    const token = getAccessToken()

    if (!token) {
      throw new Error('Authentication required to place an order.')
    }

    const payload = {
      recipient_name: dto.buyerName.trim(),
      phone: dto.phone.trim(),
      delivery_address: dto.deliveryAddress.trim(),
      delivery_city: (dto.deliveryCity || 'Addis Ababa').trim(),
      delivery_notes: (dto.deliveryNotes || '').trim()
    }

    try {
      const res = await authenticatedFetch<BackendOrderResponse>('api/orders/checkout', {
        method: 'POST',
        body: payload
      })

      if (!res) {
        throw new Error('Server returned empty order response.')
      }

      return mapBackendOrder(res)
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to create order'))
    }
  }

  /**
   * Update order status via PATCH /api/orders/{id}/status
   */
  async updateOrderStatus(orderId: number | string, status: OrderStatus): Promise<MarketplaceOrder | null> {
    const { orders } = useMockDataStore()
    const token = getAccessToken()

    if (!token) {
      throw new Error('Authentication required to update order status.')
    }

    const lowerStatus = status.toLowerCase()

    try {
      const res = await authenticatedFetch<BackendOrderResponse>(`api/orders/${orderId}/status`, {
        method: 'PATCH',
        body: { status: lowerStatus }
      })
      const updated = mapBackendOrder(res)
      const index = orders.value.findIndex(o => String(o.id) === String(orderId) || o.backendId === String(orderId))
      if (index !== -1) {
        orders.value[index] = updated
      }
      return updated
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to update order status'))
    }
  }
}
