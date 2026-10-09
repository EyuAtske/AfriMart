import type { IOrderRepository } from '../interfaces/IOrderRepository'
import type { MarketplaceOrder, OrderStatus, PaymentStatus, CreateOrderDTO, CartItem, ChapaCheckoutRequest, ChapaCheckoutResponse } from '~/types/order'
import { authenticatedFetch, extractError, getAccessToken, getApiBase } from './apiHelpers'

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
  method?: string
  Method?: string
  total?: string | number
  Total?: string | number
  created_at?: string
  CreatedAt?: string
  items?: any[]
  Items?: any[]
  order?: BackendOrderResponse
  payment?: any
  Payment?: any
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

function toNumber(value: unknown): number {
  if (typeof value === 'number') return value
  const parsed = parseFloat(String(value || 0))
  return Number.isFinite(parsed) ? parsed : 0
}

function pickString(...values: unknown[]): string {
  for (const value of values) {
    if (value && typeof value === 'object' && !('String' in value)) continue
    const str = extractNullableString(value).trim()
    if (str) return str
  }
  return ''
}

function pickNestedName(...values: unknown[]): string {
  for (const value of values) {
    const direct = pickString(value)
    if (direct) return direct

    if (value && typeof value === 'object') {
      const obj = value as Record<string, unknown>
      const nested = pickString(obj.Name, obj.name, obj.ShopName, obj.shop_name, obj.shopName)
      if (nested) return nested
    }
  }
  return ''
}

function resolveImageUrl(key: unknown): string {
  const raw = extractNullableString(key).trim()
  if (!raw) return ''

  let minioBase = 'http://localhost:9000/afrimart-images'
  try {
    const config = useRuntimeConfig()
    const configured = ((config.public as any)?.minioBaseUrl as string || (config.public as any)?.minioUrl as string || '').replace(/\/$/, '')
    if (configured) minioBase = configured
  } catch {
    // headless/test context
  }

  const cleaned = raw.replace(/^https?:\/\/(minio|afrimart-minio):9000/, 'http://localhost:9000')
  if (cleaned.startsWith('http://') || cleaned.startsWith('https://') || cleaned.startsWith('/images/')) {
    return cleaned
  }
  if (cleaned.startsWith('/')) return `${minioBase}${cleaned}`
  return `${minioBase}/${cleaned}`
}

function firstImageFromMedia(raw: any): string {
  const images = raw?.images || raw?.Images || raw?.media || raw?.Media || []
  if (!Array.isArray(images) || !images.length) return ''
  const primary = images.find((img: any) => img?.IsPrimary || img?.is_primary || img?.isPrimary) || images[0]
  return resolveImageUrl(primary?.ObjectKey ?? primary?.object_key ?? primary?.url ?? primary?.URL)
}

function normalizeStatus(raw: unknown): OrderStatus {
  const value = String(raw || 'pending').toLowerCase()
  if (value === 'confirmed') return 'Confirmed'
  if (value === 'processing') return 'Processing'
  if (value === 'shipped') return 'Shipped'
  if (value === 'delivered') return 'Delivered'
  if (value === 'cancelled') return 'Cancelled'
  if (value === 'ordered') return 'Ordered'
  return 'Pending'
}

function normalizePaymentStatus(raw: unknown): PaymentStatus {
  const value = String(raw || '').toLowerCase()
  if (value === 'successful' || value === 'success' || value === 'paid') return 'Paid'
  if (value === 'processing') return 'Processing'
  if (value === 'failed') return 'Failed'
  if (value === 'cancelled' || value === 'canceled') return 'Cancelled'
  return 'Pending'
}

function normalizePaymentMethod(raw: unknown): string {
  const label = extractNullableString(raw).trim()
  const value = label.toLowerCase().replace(/[\s-]+/g, '_')
  if (['online', 'online_payment', 'chapa', 'chapa_payment', 'card', 'bank_card'].includes(value)) return 'Online Payment'
  if (['cod', 'cash_on_delivery', 'cash_delivery', 'cash'].includes(value)) return 'Cash on delivery'
  if (value === 'telebirr') return 'Telebirr'
  if (['cbe', 'cbe_birr'].includes(value)) return 'CBE Birr'
  return label
}

function mapOrderItem(item: any): CartItem {
  const product = item.Product || item.product || item.productDetails || item.product_details || {}
  const shop = item.Shop || item.shop || product.Shop || product.shop
  const productId = item.ProductID ?? item.product_id ?? item.productId ?? product.ID ?? product.id ?? ''
  const quantity = toNumber(item.Quantity ?? item.quantity ?? 0)
  const rawPrice = item.Price ?? item.price ?? product.Price ?? product.price
  const price = rawPrice === undefined || rawPrice === null ? undefined : toNumber(rawPrice)
  const productImage = firstImageFromMedia(item) ||
    firstImageFromMedia(product) ||
    resolveImageUrl(item.ProductImage ?? item.product_image ?? item.productImage ?? item.Image ?? item.image ?? product.ProductImage ?? product.product_image ?? product.productImage ?? product.Image ?? product.image)

  return {
    productId,
    quantity,
    backendItemId: pickString(item.ID, item.id),
    productName: pickString(item.ProductName, item.product_name, item.productName, product.Name, product.name),
    productImage,
    productShop: pickNestedName(item.ShopName, item.shop_name, item.shopName, shop, product.ShopName, product.shop_name, product.shopName),
    price,
    lineTotal: price !== undefined ? price * quantity : (item.lineTotal ?? item.line_total ?? item.subtotal ?? item.Subtotal !== undefined ? toNumber(item.lineTotal ?? item.line_total ?? item.subtotal ?? item.Subtotal) : undefined)
  }
}

function mapBackendOrder(raw: any): MarketplaceOrder {
  const o = raw?.order ? raw.order : raw
  const rawItems = raw?.items || raw?.Items || o?.items || o?.Items || []
  const payment = raw?.payment || raw?.Payment || o?.payment || o?.Payment || null

  const rawId = o.ID ?? o.id
  const strId = rawId !== undefined && rawId !== null ? String(rawId) : ''
  const id = strId

  const rawSubtotal = o.Subtotal ?? o.subtotal
  const subtotalNum = toNumber(rawSubtotal)
  const rawCreatedAt = o.CreatedAt || o.created_at
  const paymentAmount = payment ? toNumber(payment.Amount ?? payment.amount) : undefined

  return {
    id,
    backendId: strId,
    buyerName: o.RecipientName || o.recipient_name || o.buyer_name || o.buyerName || '',
    items: (rawItems || []).map(mapOrderItem),
    deliveryAddress: o.DeliveryAddress || o.delivery_address || o.deliveryAddress || '',
    deliveryCity: o.DeliveryCity || o.delivery_city || o.deliveryCity || '',
    deliveryNotes: extractNullableString(o.DeliveryNotes ?? o.delivery_notes ?? o.deliveryNotes),
    phone: o.Phone || o.phone || '',
    paymentMethod: normalizePaymentMethod(payment?.PaymentMethod ?? payment?.payment_method ?? o.PaymentMethod ?? o.payment_method ?? o.paymentMethod ?? o.Method ?? o.method ?? raw?.PaymentMethod ?? raw?.payment_method ?? raw?.paymentMethod),
    paymentStatus: normalizePaymentStatus(payment?.PaymentStatus ?? payment?.payment_status ?? o.PaymentStatus ?? o.payment_status ?? o.paymentStatus ?? raw?.PaymentStatus ?? raw?.payment_status ?? raw?.paymentStatus),
    paymentAmount,
    transactionId: extractNullableString(payment?.TransactionID ?? payment?.transaction_id ?? raw?.TransactionID ?? raw?.transaction_id),
    status: normalizeStatus(o.Status || o.status),
    date: rawCreatedAt ? (String(rawCreatedAt).split('T')[0] || '') : '',
    total: subtotalNum
  }
}

export class ApiOrderRepository implements IOrderRepository {
  private async fetchMappedOrderById(id: number | string): Promise<MarketplaceOrder | null> {
    const res = await authenticatedFetch<any>(`api/orders/${id}`, {
      method: 'GET'
    })
    if (!res) return null
    return this.enrichOrderProductFields(mapBackendOrder(res))
  }

  private async enrichOrderProductFields(order: MarketplaceOrder): Promise<MarketplaceOrder> {
    if (!order.items.length) return order

    const apiBase = getApiBase()
    const items = await Promise.all(order.items.map(async (item) => {
      if (item.productName && item.productImage && item.productShop && item.price !== undefined) return item
      if (!item.productId) return item

      try {
        const raw = await $fetch<any>(`${apiBase}/api/products/${encodeURIComponent(String(item.productId))}`, {
          method: 'GET'
        })
        const product = raw?.product || raw || {}
        const enriched = mapOrderItem({
          ProductID: item.productId,
          Quantity: item.quantity,
          Price: item.price,
          ProductName: item.productName,
          ProductImage: item.productImage,
          ShopName: item.productShop,
          Product: product,
          images: raw?.images || raw?.Images
        })
        const price = item.price ?? enriched.price
        return {
          ...item,
          productName: item.productName || enriched.productName,
          productImage: item.productImage || enriched.productImage,
          productShop: item.productShop || enriched.productShop,
          price,
          lineTotal: price === undefined ? item.lineTotal : price * item.quantity
        }
      } catch {
        return item
      }
    }))

    return {
      ...order,
      items
    }
  }

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
      const mapped: MarketplaceOrder[] = list.map(mapBackendOrder)
      return Promise.all(mapped.map(async (order) => {
        if (!order.backendId) return this.enrichOrderProductFields(order)
        try {
          return await this.fetchMappedOrderById(order.backendId) || order
        } catch {
          return this.enrichOrderProductFields(order)
        }
      }))
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
      const mapped: MarketplaceOrder[] = list.map(mapBackendOrder)
      return Promise.all(mapped.map(async (order) => {
        if (!order.backendId) return this.enrichOrderProductFields(order)
        try {
          return await this.fetchMappedOrderById(order.backendId) || order
        } catch {
          return this.enrichOrderProductFields(order)
        }
      }))
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
      return await this.fetchMappedOrderById(id)
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

      return this.enrichOrderProductFields(mapBackendOrder(res))
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
      return this.enrichOrderProductFields(mapBackendOrder(res))
    } catch (err: any) {
      throw new Error(extractError(err, 'Failed to update order status'))
    }
  }
}
