import type { IOrderRepository } from '../interfaces/IOrderRepository'
import type { MarketplaceOrder, OrderStatus, CreateOrderDTO, CartItem, ChapaCheckoutRequest, ChapaCheckoutResponse } from '~/types/order'
import { useMockDataStore } from '../mock/MockDataStore'
import { authenticatedFetch, extractError, getAccessToken } from './apiHelpers'

export interface BackendOrderResponse {
  id?: string | number
  ID?: string | number
  user_id?: string
  UserID?: string
  subtotal?: string | number
  Subtotal?: string | number
  status?: string
  Status?: string
  recipient_name?: string
  RecipientName?: string
  phone?: string
  Phone?: string
  delivery_address?: string
  DeliveryAddress?: string
  delivery_city?: string
  DeliveryCity?: string
  delivery_notes?: string
  DeliveryNotes?: string | { String?: string; Valid?: boolean }
  created_at?: string
  CreatedAt?: string
  items?: any[]
  Items?: any[]
  order?: BackendOrderResponse
}

function extractNullableString(value: unknown): string {
  if (!value) return ''
  if (typeof value === 'string') return value
  if (typeof value === 'object') {
    const nullable = value as { String?: unknown; Valid?: unknown }
    return nullable.Valid && typeof nullable.String === 'string' ? nullable.String : ''
  }
  return String(value)
}

function mapBackendOrder(raw: any): MarketplaceOrder {
  const o = raw?.order ? raw.order : raw
  const rawItems = raw?.items || raw?.Items || o?.items || o?.Items || []

  const rawId = o.ID ?? o.id
  const strId = rawId !== undefined && rawId !== null ? String(rawId) : ''
  const id = strId || Date.now()

  const rawStatus = String(o.Status || o.status || 'pending').toLowerCase()
  let status: OrderStatus = 'Pending'
  if (rawStatus === 'confirmed') status = 'Confirmed'
  else if (rawStatus === 'processing') status = 'Processing'
  else if (rawStatus === 'shipped') status = 'Shipped'
  else if (rawStatus === 'delivered') status = 'Delivered'
  else if (rawStatus === 'cancelled') status = 'Cancelled'
  else if (rawStatus === 'ordered') status = 'Ordered'
  else status = 'Pending'

  const rawSubtotal = o.Subtotal ?? o.subtotal
  const subtotalNum = typeof rawSubtotal === 'number' ? rawSubtotal : parseFloat(String(rawSubtotal || 0)) || 0
  const rawCreatedAt = o.CreatedAt || o.created_at

  return {
    id,
    backendId: strId,
    buyerName: o.RecipientName || o.recipient_name || o.buyer_name || o.buyerName || 'Valued Customer',
    items: (rawItems || []).map((item: any) => ({
      productId: item.product_id || item.ProductID || item.productId || 0,
      quantity: item.quantity ?? item.Quantity ?? 1
    })),
    deliveryAddress: o.DeliveryAddress || o.delivery_address || o.deliveryAddress || '',
    deliveryCity: o.DeliveryCity || o.delivery_city || o.deliveryCity || 'Addis Ababa',
    deliveryNotes: extractNullableString(o.DeliveryNotes ?? o.delivery_notes ?? o.deliveryNotes),
    phone: o.Phone || o.phone || '',
    paymentMethod: 'Cash on delivery',
    paymentStatus: 'Pending',
    status,
    date: rawCreatedAt ? String(rawCreatedAt).split('T')[0] : new Date().toISOString().split('T')[0]!,
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

  async initiateChapaCheckout(request: ChapaCheckoutRequest): Promise<ChapaCheckoutResponse> {
    if (!getAccessToken()) {
      throw new Error('Authentication required to start online checkout.')
    }

    try {
      const response = await authenticatedFetch<ChapaCheckoutResponse>('api/orders/checkout', {
        method: 'POST',
        body: request
      })

      if (!response?.checkout_url || !response.order?.ID || !response.transaction_id) {
        throw new Error('Server returned an incomplete payment checkout response.')
      }

      return response
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to start online checkout'))
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
