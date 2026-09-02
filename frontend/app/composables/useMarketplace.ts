import type { Product, ProductFilterParams } from '~/types/product'
import type { MarketplaceOrder, OrderStatus, PaymentStatus, CartItem, CartProductItem, CreateOrderDTO } from '~/types/order'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import { useRepositories } from '~/composables/useRepositories'

export const formatPrice = (amount: number) =>
  `${amount.toLocaleString()} ETB`

export const useMarketplace = () => {
  const { products, cart, orders, reviews, addReview: addReviewToStore } = useMockDataStore()
  const { productRepo, orderRepo } = useRepositories()

  const categories = computed(() => [
    'All',
    ...Array.from(new Set(products.value.map(product => product.category)))
  ])

  const getProduct = (id: number) =>
    products.value.find(product => product.id === id) || null

  const getProductReviews = (productId: number) =>
    reviews.value.filter(review => review.productId === productId)

  const getUserReviewForProduct = (productId: number, orderId?: number) => {
    return reviews.value.find(
      r => r.productId === productId && (orderId ? r.orderId === orderId : true)
    )
  }

  const submitProductReview = (productId: number, rating: number, comment: string, authorName?: string, orderId?: number) => {
    return addReviewToStore(productId, rating, comment, authorName, orderId)
  }

  const filterProducts = (filters: ProductFilterParams) => {
    const search = filters.search?.trim().toLowerCase() || ''
    const category = filters.category || 'All'

    return products.value.filter((product) => {
      const matchesSearch = !search ||
        product.name.toLowerCase().includes(search) ||
        product.shop.toLowerCase().includes(search) ||
        product.description.toLowerCase().includes(search)

      const matchesCategory = category === 'All' || product.category === category

      return product.status === 'Active' && matchesSearch && matchesCategory
    })
  }

  const addToCart = (productId: number) => {
    const product = getProduct(productId)
    if (!product || product.stock < 1) return

    const existingIndex = cart.value.findIndex(item => item.productId === productId)

    if (existingIndex !== -1 && cart.value[existingIndex]) {
      const updatedItem = {
        productId,
        quantity: Math.min(cart.value[existingIndex].quantity + 1, product.stock)
      }
      const nextCart = [...cart.value]
      nextCart[existingIndex] = updatedItem
      cart.value = nextCart
      return
    }

    cart.value = [...cart.value, { productId, quantity: 1 }]
  }

  const updateCartQuantity = (productId: number, quantity: number) => {
    const product = getProduct(productId)

    if (quantity < 1) {
      cart.value = cart.value.filter(item => item.productId !== productId)
      return
    }

    const existingIndex = cart.value.findIndex(cartItem => cartItem.productId === productId)
    if (existingIndex !== -1 && product && cart.value[existingIndex]) {
      const nextCart = [...cart.value]
      nextCart[existingIndex] = {
        productId,
        quantity: Math.min(quantity, product.stock)
      }
      cart.value = nextCart
    }
  }

  const removeFromCart = (productId: number) => {
    cart.value = cart.value.filter(item => item.productId !== productId)
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

  const createOrder = (details: CreateOrderDTO) => {
    if (!cart.value.length) return null
    return orderRepo.createOrder(cart.value, cartSubtotal.value, details)
  }

  const updateOrderStatus = (orderId: number, status: OrderStatus) => {
    return orderRepo.updateOrderStatus(orderId, status)
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

  const updateProduct = (id: number, updates: Partial<Product>) => {
    return productRepo.updateProduct(id, updates)
  }

  const deleteProduct = (id: number) => {
    return productRepo.deleteProduct(id)
  }

  const toggleProductStatus = (id: number) => {
    return productRepo.toggleProductStatus(id)
  }

  const updateProductStock = (id: number, newStock: number) => {
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
    createOrder,
    updateOrderStatus,
    getOrderProducts,
    updateProduct,
    deleteProduct,
    toggleProductStatus,
    updateProductStock
  }
}
