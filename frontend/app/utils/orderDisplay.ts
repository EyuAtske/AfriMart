import type { CartItem, MarketplaceOrder } from '~/types/order'
import type { Product } from '~/types/product'

export const getOrderNumber = (order: MarketplaceOrder) =>
  order.backendId || String(order.id)

export const getOrderItems = (order: MarketplaceOrder) =>
  order.items || []

export const getPrimaryOrderItem = (order: MarketplaceOrder) =>
  getOrderItems(order)[0] || null

export const getOrderTitle = (order: MarketplaceOrder) => {
  const item = getPrimaryOrderItem(order)
  if (!item?.productName) return 'Product details unavailable'

  const extraCount = getOrderItems(order).length - 1
  return extraCount > 0 ? `${item.productName} + ${extraCount} more` : item.productName
}

export const formatPaymentSummary = (order: MarketplaceOrder) => {
  const parts = [order.paymentMethod, order.paymentStatus]
    .filter(Boolean)
    .filter((value, index, array) => array.indexOf(value) === index)
  return parts.length ? parts.join(' - ') : 'Payment details unavailable'
}

export const hasReviewableProduct = (item: CartItem) =>
  Boolean(item.productId && item.productName)

export const toReviewProduct = (item: CartItem): Product | null => {
  if (!hasReviewableProduct(item)) return null

  return {
    id: item.productId,
    name: item.productName || 'Product details unavailable',
    image: item.productImage || '',
    shop: item.productShop || 'Shop unavailable',
    price: item.price ?? 0,
    category: 'Men',
    stock: 0,
    rating: '0',
    description: '',
    status: 'Active'
  } as Product
}
