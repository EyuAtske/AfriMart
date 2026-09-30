import { computed } from 'vue'
import type { Product, ProductFilterParams } from '~/types/product'
import type { MarketplaceOrder, OrderStatus, PaymentStatus, CartItem, CartProductItem, CreateOrderDTO } from '~/types/order'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import { useRepositories } from '~/composables/useRepositories'
import { STATIC_CATALOG, ensureCategoryCatalog } from '~/utils/categoryCatalog'

interface MarketplaceProductFilters extends ProductFilterParams {
  subcategoryId?: string
  minPrice?: number
  maxPrice?: number
  availability?: 'all' | 'available' | 'sold-out'
  sortBy?: 'default' | 'price-asc' | 'price-desc'
  includeInactive?: boolean
}

export const formatPrice = (amount: number) =>
  `${amount.toLocaleString()} ETB`

export const useMarketplace = () => {
  const { products, cart, orders, reviews, addReview: addReviewToStore } = useMockDataStore()
  const { productRepo, cartRepo, orderRepo, authRepo } = useRepositories()
  let gtag: any = () => { }
  try {
    const g = useGtag()
    if (g?.gtag) gtag = g.gtag
  } catch {
    // headless/test context
  }

  const categories = computed(() => {
    const list: string[] = ['All']
    for (const cat of STATIC_CATALOG) {
      if (cat.name && !list.includes(cat.name)) {
        list.push(cat.name)
      }
    }
    for (const product of products.value) {
      if (product.category && !list.includes(product.category)) {
        list.push(product.category)
      }
    }
    return list
  })

  const getProduct = (id: number | string) =>
    products.value.find(product => String(product.id) === String(id)) || null

  const getProductReviews = (productId: number | string) =>
    reviews.value.filter(review => String(review.productId) === String(productId))

  const getUserReviewForProduct = (productId: number | string, orderId?: number | string) => {
    return reviews.value.find(
      r => String(r.productId) === String(productId) && (orderId ? String(r.orderId) === String(orderId) : true)
    )
  }

  const submitProductReview = (productId: number | string, rating: number, comment: string, authorName?: string, orderId?: number | string) => {
    return addReviewToStore(typeof productId === 'number' ? productId : (parseInt(String(productId), 10) || Date.now()), rating, comment, authorName, orderId)
  }

  const fetchProducts = async (filterParams?: ProductFilterParams) => {
    try {
      await ensureCategoryCatalog()
      const res = await productRepo.getProducts(filterParams || { page: 1, pageSize: 50 })
      if (res?.data) {
        products.value = res.data
      }
      return res?.data || []
    } catch (err: any) {
      console.warn('Failed to fetch marketplace products:', err?.message || err)
      return []
    }
  }

  try {
    if (import.meta.client) {
      // Public pages restore the session too, so own-shop cart protection works after refresh.
      const sessionRestored = useState<boolean>('marketplace-auth-session-restored', () => false)
      if (!sessionRestored.value) {
        sessionRestored.value = true
        authRepo.getCurrentSession().catch(() => {})
      }

      try {
        useSellerShop()
      } catch {
        // Headless/test context without Nuxt composable state
      }

      const productsHydrated = useState<boolean>('marketplace-products-hydrated', () => false)
      if (!productsHydrated.value) {
        productsHydrated.value = true
        fetchProducts()
      }
    }
  } catch {
    // Vitest/headless context
  }

  const filterProducts = (filters: MarketplaceProductFilters, customList?: Product[]) => {
    const search = filters.search?.trim().toLowerCase() || ''
    const category = filters.category || 'All'
    const targetList = customList || products.value
    const categoryEntry = STATIC_CATALOG.find(cat => cat.name.toLowerCase() === category.toLowerCase())
    const selectedSubcategory = filters.subCategory && filters.subCategory !== 'All'
      ? (categoryEntry?.subcategories.find(sub => sub.name.toLowerCase() === filters.subCategory?.toLowerCase()) ||
        STATIC_CATALOG.flatMap(cat => cat.subcategories).find(sub => sub.name.toLowerCase() === filters.subCategory?.toLowerCase()))
      : undefined

    const filtered = targetList.filter((product) => {
      const matchesSearch = !search ||
        product.name.toLowerCase().includes(search) ||
        product.shop.toLowerCase().includes(search) ||
        product.description.toLowerCase().includes(search) ||
        (product.subCategory && product.subCategory.toLowerCase().includes(search))

      const isGenderCategory = category === 'Men' || category === 'Women' || category === 'Kids'
      const matchesCategory = category === 'All' || (isGenderCategory
        ? product.category.toLowerCase() === category.toLowerCase() || product.gender?.toLowerCase() === category.toLowerCase()
        : product.category.toLowerCase() === category.toLowerCase() || Boolean(categoryEntry?.id && product.categoryId === categoryEntry.id))
      const matchesSubcategory = (!filters.subCategory || filters.subCategory === 'All') ||
        (filters.subcategoryId && product.subcategoryId === filters.subcategoryId) ||
        (selectedSubcategory?.id && product.subcategoryId === selectedSubcategory.id) ||
        product.subCategory?.toLowerCase() === filters.subCategory?.toLowerCase()
      const matchesMinPrice = filters.minPrice === undefined || product.price >= filters.minPrice
      const matchesMaxPrice = filters.maxPrice === undefined || product.price <= filters.maxPrice
      const matchesAvailability = filters.availability === undefined || filters.availability === 'all' ||
        (filters.availability === 'available' && product.stock > 0) ||
        (filters.availability === 'sold-out' && product.stock <= 0)

      return (filters.includeInactive || product.status === 'Active') && matchesSearch && matchesCategory &&
        matchesSubcategory && matchesMinPrice && matchesMaxPrice && matchesAvailability
    })

    if (filters.sortBy === 'price-asc') return filtered.sort((a, b) => a.price - b.price)
    if (filters.sortBy === 'price-desc') return filtered.sort((a, b) => b.price - a.price)
    return filtered
  }

  const syncCartFromBackend = async () => {
    try {
      const backendCart = await cartRepo.getCart()
      if (backendCart?.items && Array.isArray(backendCart.items)) {
        cart.value = backendCart.items.map((item: any) => ({
          productId: item.ProductID || item.product_id || item.productId,
          quantity: item.Quantity ?? item.quantity ?? 1,
          backendItemId: item.ID || item.id
        }))
      } else {
        cart.value = []
      }
    } catch (err: any) {
      console.warn('Cart sync warning:', err?.message || err)
    }
  }

  try {
    if (import.meta.client) {
      const cartHydrated = useState<boolean>('marketplace-cart-hydrated', () => false)

      if (!cartHydrated.value) {
        cartHydrated.value = true
        syncCartFromBackend()
      }
    }
  } catch (err: any) {
    console.warn('Cart sync warning:', err?.message || err)
  }

const isOwnProduct = (target: number | string | Product): boolean => {
  const { shop } = useMockDataStore()
  if (!shop.value) return false

  let targetShopName = ''
  let targetShopId = ''

  if (typeof target === 'object' && target !== null) {
    targetShopName = target.shop || ''
    targetShopId = String((target as any).shopId || (target as any).ShopID || '')
  } else {
    const found = getProduct(target)
    if (found) {
      targetShopName = found.shop || ''
      targetShopId = String((found as any).shopId || (found as any).ShopID || '')
    }
  }

  const sellerShopId = String(shop.value.backendId || shop.value.id || '')
  if (targetShopId) {
    return sellerShopId === targetShopId
  }

  const sellerShopName = (shop.value.name || '').trim().toLowerCase()
  const pShopName = targetShopName.trim().toLowerCase()

  if (sellerShopName && pShopName && sellerShopName === pShopName) {
    return true
  }

  return false
}

const addToCart = async (productId: number | string) => {
  const product = getProduct(productId)
  if (!product) return

  if (isOwnProduct(product)) {
    throw new Error('You cannot purchase items listed by your own shop.')
  }

  if (product.stock < 1) {
    throw new Error('This item is currently sold out.')
  }

  const existingIndex = cart.value.findIndex(item => String(item.productId) === String(productId))
  const prevCart = [...cart.value]

  if (existingIndex !== -1 && cart.value[existingIndex]) {
    const updatedItem = {
      productId,
      quantity: Math.min(cart.value[existingIndex].quantity + 1, product.stock),
      backendItemId: cart.value[existingIndex].backendItemId
    }

    const nextCart = [...cart.value]
    nextCart[existingIndex] = updatedItem
    cart.value = nextCart
  } else {
    cart.value = [...cart.value, { productId, quantity: 1 }]
  }

    try {
      const res = await cartRepo.addItem(String(productId), 1)
      const backendId = (res as any)?.ID || (res as any)?.id
      if (backendId) {
        const item = cart.value.find(i => String(i.productId) === String(productId))
        if (item) item.backendItemId = backendId
      }
    } catch (err: any) {
    cart.value = prevCart
    throw err
  }

  if (import.meta.client) {
    gtag('event', 'add_to_cart', {
      currency: 'ETB',
      value: Number(product.price),
      items: [
        {
          item_id: String(product.id),
          item_name: product.name,
          item_category: product.category,
          price: Number(product.price),
          quantity: 1
        }
      ]
    })
  }
}

const updateCartQuantity = async (productId: number | string, quantity: number) => {
  const product = getProduct(productId)

  if (quantity < 1) {
    await removeFromCart(productId)
    return
  }

  const existingIndex = cart.value.findIndex(cartItem => String(cartItem.productId) === String(productId))
  if (existingIndex !== -1 && product && cart.value[existingIndex]) {
    const backendItemId = cart.value[existingIndex].backendItemId
    const newQty = Math.min(quantity, product.stock)

    if (!backendItemId) {
      console.error('Cart integrity error: missing backend item ID for product', productId)
      await syncCartFromBackend()
      throw new Error('Cart data is out of sync. Cart has been refreshed — please try again.')
    }

    const nextCart = [...cart.value]
    nextCart[existingIndex] = {
      productId,
      quantity: newQty,
      backendItemId
    }
    cart.value = nextCart

    await cartRepo.updateItemQuantity(backendItemId, newQty)
  }
}

const removeFromCart = async (productId: number | string) => {
  const product = getProduct(productId)
  if (!product) return

  const existingItem = cart.value.find(item => String(item.productId) === String(productId))
  if (!existingItem) return

  const backendItemId = existingItem.backendItemId

  if (!backendItemId) {
    console.error('Cart integrity error: missing backend item ID for product', productId)
    await syncCartFromBackend()
    throw new Error('Cart data is out of sync. Cart has been refreshed — please try again.')
  }

  cart.value = cart.value.filter(item => String(item.productId) !== String(productId))

  await cartRepo.removeItem(backendItemId)

  if (import.meta.client) {
    gtag('event', 'remove_from_cart', {
      currency: 'ETB',
      value: Number(product.price) * existingItem.quantity,
      items: [
        {
          item_id: String(product.id),
          item_name: product.name,
          item_category: product.category,
          price: Number(product.price),
          quantity: existingItem.quantity
        }
      ]
    })
  }
}

const cartProducts = computed<CartProductItem[]>(() =>
  cart.value
    .map((item) => {
      const product = getProduct(item.productId)
      if (!product) return null

      return {
        ...item,
        product,
        lineTotal: product.price * item.quantity
      }
    })
    .filter((item): item is CartProductItem => item !== null)
)

const cartSubtotal = computed(() =>
  cartProducts.value.reduce((total, item) => total + item.lineTotal, 0)
)

const fetchUserOrders = async () => {
  try {
    const userOrders = await orderRepo.getOrders()
    if (Array.isArray(userOrders)) {
      orders.value = userOrders
    }
  } catch (err: any) {
    console.warn('User orders fetch warning:', err?.message || err)
  }
}

const fetchSellerOrders = async () => {
  try {
    const sellerOrders = await orderRepo.getSellerOrders()
    if (Array.isArray(sellerOrders)) {
      orders.value = sellerOrders
    }
  } catch (err: any) {
    console.warn('Seller orders fetch warning:', err?.message || err)
  }
}

const createOrder = async (details: CreateOrderDTO): Promise<{ order: MarketplaceOrder; refreshError: string | null } | null> => {
  if (!cart.value.length) return null
  const newOrder = await orderRepo.createOrder(cart.value, cartSubtotal.value, details)

  // Refetch real cart and buyer orders from the API after confirmed checkout
  const failures: string[] = []

    try {
      const backendCart = await cartRepo.getCart()
      if (backendCart?.items && Array.isArray(backendCart.items)) {
        cart.value = backendCart.items.map((item: any) => ({
          productId: item.ProductID || item.product_id || item.productId,
          quantity: item.Quantity ?? item.quantity ?? 1,
          backendItemId: item.ID || item.id
        }))
      } else {
        cart.value = []
      }
    } catch {
      failures.push('cart')
    }

    try {
      const userOrders = await orderRepo.getOrders()
      if (Array.isArray(userOrders)) {
        orders.value = userOrders
      }
    } catch {
      failures.push('orders')
    }

    try {
      await productRepo.getProducts({ page: 1, pageSize: 50 })
    } catch {
      failures.push('products')
    }

    const refreshError = failures.length
      ? `Order placed, but latest ${failures.join(' and ')} data could not be refreshed.`
      : null

    return { order: newOrder, refreshError }
  }

  /**
   * Retry cart + order + product refetch after a successful checkout.
   * Throws if any refetch still fails.
   */
  const retryPostCheckoutRefresh = async () => {
    const errors: string[] = []

    try {
      const backendCart = await cartRepo.getCart()
      if (backendCart?.items && Array.isArray(backendCart.items)) {
        cart.value = backendCart.items.map((item: any) => ({
          productId: item.ProductID || item.product_id || item.productId,
          quantity: item.Quantity ?? item.quantity ?? 1,
          backendItemId: item.ID || item.id
        }))
      } else {
        cart.value = []
      }
    } catch {
      errors.push('cart')
    }

  try {
    const userOrders = await orderRepo.getOrders()
    if (Array.isArray(userOrders)) {
      orders.value = userOrders
    }
  } catch {
    errors.push('orders')
  }

  try {
    await productRepo.getProducts({ page: 1, pageSize: 50 })
  } catch {
    errors.push('products')
  }

  if (errors.length) {
    throw new Error(`Could not refresh ${errors.join(' and ')} data.`)
  }
}

const updateOrderStatus = async (orderId: number | string, status: OrderStatus) => {
  const res = await orderRepo.updateOrderStatus(orderId, status)
  if (res) {
    const orderIndex = orders.value.findIndex(o => String(o.id) === String(orderId) || o.backendId === String(orderId))
    if (orderIndex !== -1 && orders.value[orderIndex]) {
      orders.value[orderIndex] = res
    }
  }
  return res
}

const getOrderProducts = (order: MarketplaceOrder): CartProductItem[] =>
  order.items
    .map((item) => {
      const product = getProduct(item.productId)
      if (!product) return null

      return {
        ...item,
        product,
        lineTotal: product.price * item.quantity
      }
    })
    .filter((item): item is CartProductItem => item !== null)

const updateProduct = (id: number | string, updates: Partial<Product>) => {
  return productRepo.updateProduct(id, updates)
}

const deleteProduct = (id: number | string) => {
  return productRepo.deleteProduct(id)
}

const toggleProductStatus = (id: number | string) => {
  return productRepo.toggleProductStatus(id)
}

const updateProductStock = (id: number | string, newStock: number) => {
  return productRepo.updateProduct(id, { stock: Math.max(0, newStock) })
}

return {
  products,
  categories,
  cart,
  orders,
  reviews,
  cartProducts,
  cartSubtotal,
  formatPrice,
  getProduct,
  getProductReviews,
  getUserReviewForProduct,
  submitProductReview,
  filterProducts,
  addToCart,
  updateCartQuantity,
  removeFromCart,
  syncCartFromBackend,
  createOrder,
  retryPostCheckoutRefresh,
  fetchUserOrders,
  fetchSellerOrders,
  fetchProducts,
  updateOrderStatus,
  getOrderProducts,
  isOwnProduct,
  updateProduct,
  deleteProduct,
  toggleProductStatus,
  updateProductStock
}
}
