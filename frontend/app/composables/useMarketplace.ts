import type { Product, ProductFilterParams } from '~/types/product'
import type { MarketplaceOrder, OrderStatus, PaymentStatus, CartItem, CartProductItem, CreateOrderDTO } from '~/types/order'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import { useRepositories } from '~/composables/useRepositories'

export const formatPrice = (amount: number) =>
  `${amount.toLocaleString()} ETB`

export const useMarketplace = () => {
  const { products, cart, orders, reviews, addReview: addReviewToStore } = useMockDataStore()
  const { productRepo, orderRepo, cartRepo } = useRepositories()
  const { gtag } = useGtag()

  const categories = computed(() => [
    'All',
    ...Array.from(new Set(products.value.map(product => product.category)))
  ])

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

  const filterProducts = (filters: ProductFilterParams, customList?: Product[]) => {
    const search = filters.search?.trim().toLowerCase() || ''
    const category = filters.category || 'All'
    const targetList = customList || products.value

    return targetList.filter((product) => {
      const matchesSearch = !search ||
        product.name.toLowerCase().includes(search) ||
        product.shop.toLowerCase().includes(search) ||
        product.description.toLowerCase().includes(search) ||
        (product.subCategory && product.subCategory.toLowerCase().includes(search))

      const matchesCategory = category === 'All' ||
        product.category === category ||
        product.subCategory === category

      return product.status === 'Active' && matchesSearch && matchesCategory
    })
  }

  const syncCartFromBackend = async () => {
    try {
      const backendCart = await cartRepo.getCart()
      if (backendCart?.items) {
        cart.value = backendCart.items.map(item => ({
          productId: item.product_id,
          quantity: item.quantity,
          backendItemId: item.id
        }))
      } else {
        cart.value = []
      }
    } catch (err: any) {
      console.warn('Cart sync warning:', err?.message || err)
    }
  }

  const addToCart = async (productId: number | string) => {
    const product = getProduct(productId)
    if (!product || product.stock < 1) return

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
      if (res?.id) {
        const item = cart.value.find(i => String(i.productId) === String(productId))
        if (item) item.backendItemId = res.id
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
      if (backendCart?.items) {
        cart.value = backendCart.items.map(item => ({
          productId: item.product_id,
          quantity: item.quantity,
          backendItemId: item.id
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

    const refreshError = failures.length
      ? `Order placed, but latest ${failures.join(' and ')} data could not be refreshed.`
      : null

    return { order: newOrder, refreshError }
  }

  /**
   * Retry cart + order refetch after a successful checkout.
   * Throws if any refetch still fails.
   */
  const retryPostCheckoutRefresh = async () => {
    const errors: string[] = []

    try {
      const backendCart = await cartRepo.getCart()
      if (backendCart?.items) {
        cart.value = backendCart.items.map(item => ({
          productId: item.product_id,
          quantity: item.quantity,
          backendItemId: item.id
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
    updateOrderStatus,
    getOrderProducts,
    updateProduct,
    deleteProduct,
    toggleProductStatus,
    updateProductStock
  }
}
