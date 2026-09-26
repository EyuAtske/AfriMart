import type { Product, ProductReview } from '~/types/product'
import type { Shop } from '~/types/shop'
import type { MarketplaceOrder, CartItem } from '~/types/order'
import type { User } from '~/types/auth'

export const initialMockReviews: ProductReview[] = [
  {
    id: 1,
    productId: 1,
    author: 'Taye Bekele',
    rating: 5,
    comment: 'Exceptional cotton quality and perfect relaxed fit. Will buy again!',
    date: 'August 14, 2026',
    status: 'approved'
  },
  {
    id: 2,
    productId: 1,
    author: 'Selam A.',
    rating: 4,
    comment: 'Very comfortable material. Sizing is accurate.',
    date: 'August 20, 2026',
    status: 'approved'
  },
  {
    id: 3,
    productId: 2,
    author: 'Marta G.',
    rating: 5,
    comment: 'Love the fabric, super soft and breathable for warm days.',
    date: 'August 18, 2026',
    status: 'approved'
  }
]

export const initialMockProducts: Product[] = []

export const initialMockUsers: User[] = [
  { id: 'usr-1', username: 'tayeb', name: 'Taye Bekele', email: 'taye@example.com', role: 'buyer', created_at: '2026-03-02' },
  { id: 'usr-2', username: 'selama', name: 'Selam Assefa', email: 'selam@example.com', role: 'buyer', created_at: '2026-03-15' },
  { id: 'usr-3', username: 'martag', name: 'Marta Girma', email: 'marta@example.com', role: 'buyer', created_at: '2026-04-10' },
  { id: 'usr-4', username: 'abebek', name: 'Abebe Kebede', email: 'abebe@atelierno.com', role: 'seller', created_at: '2026-04-20' },
  { id: 'usr-5', username: 'hanad', name: 'Hana Daniel', email: 'hana@beyondscore.com', role: 'seller', created_at: '2026-05-05' },
  { id: 'usr-6', username: 'dawitm', name: 'Dawit Mengistu', email: 'dawit@trueform.com', role: 'seller', created_at: '2026-05-18' },
  { id: 'usr-7', username: 'bethelhemh', name: 'Bethelhem Haile', email: 'beth@minimalstudio.com', role: 'seller', created_at: '2026-06-02' },
  { id: 'usr-8', username: 'robelt', name: 'Robel Tesfaye', email: 'robel@urbanthread.com', role: 'seller', created_at: '2026-06-21' },
  { id: 'usr-9', username: 'saray', name: 'Sara Yohannes', email: 'sara@marastudio.com', role: 'seller', created_at: '2026-07-11' },
  { id: 'usr-10', username: 'admin', name: 'System Admin', email: 'admin@platform.com', role: 'buyer', created_at: '2026-08-01' }
]

export const initialMockOrders: MarketplaceOrder[] = [
  {
    id: 240801,
    buyerName: 'Taye Bekele',
    items: [{ productId: 1, quantity: 1 }],
    deliveryAddress: 'Bole, Addis Ababa',
    phone: '+251 911 112 233',
    paymentMethod: 'Telebirr',
    paymentStatus: 'Paid',
    status: 'Delivered',
    date: 'March 12, 2026',
    total: 45
  },
  {
    id: 240802,
    buyerName: 'Selam Assefa',
    items: [{ productId: 2, quantity: 1 }],
    deliveryAddress: 'Kazanchis, Addis Ababa',
    phone: '+251 912 334 455',
    paymentMethod: 'CBE',
    paymentStatus: 'Paid',
    status: 'Delivered',
    date: 'April 4, 2026',
    total: 70
  },
  {
    id: 240803,
    buyerName: 'Marta Girma',
    items: [{ productId: 3, quantity: 1 }],
    deliveryAddress: 'CMC, Addis Ababa',
    phone: '+251 913 556 677',
    paymentMethod: 'Telebirr',
    paymentStatus: 'Paid',
    status: 'Delivered',
    date: 'April 22, 2026',
    total: 65
  },
  {
    id: 240804,
    buyerName: 'Taye Bekele',
    items: [{ productId: 4, quantity: 1 }],
    deliveryAddress: 'Bole, Addis Ababa',
    phone: '+251 911 112 233',
    paymentMethod: 'Telebirr',
    paymentStatus: 'Paid',
    status: 'Delivered',
    date: 'May 15, 2026',
    total: 55
  },
  {
    id: 240805,
    buyerName: 'Selam Assefa',
    items: [{ productId: 6, quantity: 1 }],
    deliveryAddress: 'Kazanchis, Addis Ababa',
    phone: '+251 912 334 455',
    paymentMethod: 'CBE',
    paymentStatus: 'Paid',
    status: 'Delivered',
    date: 'June 8, 2026',
    total: 95
  },
  {
    id: 240806,
    buyerName: 'Marta Girma',
    items: [{ productId: 7, quantity: 1 }],
    deliveryAddress: 'CMC, Addis Ababa',
    phone: '+251 913 556 677',
    paymentMethod: 'Telebirr',
    paymentStatus: 'Paid',
    status: 'Delivered',
    date: 'June 27, 2026',
    total: 80
  },
  {
    id: 240807,
    buyerName: 'Taye Bekele',
    items: [{ productId: 8, quantity: 1 }],
    deliveryAddress: 'Bole, Addis Ababa',
    phone: '+251 911 112 233',
    paymentMethod: 'Telebirr',
    paymentStatus: 'Paid',
    status: 'Delivered',
    date: 'July 19, 2026',
    total: 110
  },
  {
    id: 240808,
    buyerName: 'Selam Assefa',
    items: [{ productId: 1, quantity: 1 }],
    deliveryAddress: 'Bole, Addis Ababa',
    phone: '+251 911 000 000',
    paymentMethod: 'Cash on delivery',
    paymentStatus: 'Cash on delivery',
    status: 'Shipped',
    date: 'August 24, 2026',
    total: 130
  }
]

import { ref, watch, type Ref } from 'vue'

const vitestStore = new Map<string, Ref<any>>()

const safeState = <T>(key: string, init: () => T): Ref<T> => {
  try {
    if (typeof useState === 'function') {
      return useState<T>(key, init)
    }
  } catch {
    // Fallback for non-Nuxt test runner environments (Vitest)
  }
  if (!vitestStore.has(key)) {
    vitestStore.set(key, ref<T>(init()))
  }
  return vitestStore.get(key)! as Ref<T>
}

export const initialMockCart: CartItem[] = []

const STORAGE_KEYS = {
  USERS: 'afrimart_mock_users',
  PRODUCTS: 'afrimart_mock_products',
  ORDERS: 'afrimart_mock_orders',
  CART: 'afrimart_mock_cart',
  SHOP: 'afrimart_mock_shop',
  USER: 'afrimart_mock_user',
  IS_LOGGED_IN: 'afrimart_mock_is_logged_in',
  REVIEWS: 'afrimart_mock_reviews'
}

let isClientInitialized = false

function getItemFromStorage<T>(key: string, fallback: T): T {
  if (typeof window === 'undefined' || !window.localStorage) return fallback
  try {
    const item = localStorage.getItem(key)
    if (item !== null) {
      return JSON.parse(item)
    }
  } catch (e) {
    console.warn(`Error reading ${key} from localStorage:`, e)
  }
  return fallback
}

function setItemInStorage<T>(key: string, value: T): void {
  if (typeof window === 'undefined' || !window.localStorage) return
  try {
    localStorage.setItem(key, JSON.stringify(value))
  } catch (e) {
    console.warn(`Error writing ${key} to localStorage:`, e)
  }
}

export const useMockDataStore = () => {
  let isApiMode = false
  try {
    const config = useRuntimeConfig()
    isApiMode = config?.public?.authMode === 'api'
  } catch {
    // Non-Nuxt test runner environments (Vitest) fall back to mock mode
  }

  const users = safeState<User[]>('mock-ds-all-users', () => (isApiMode ? [] : [...initialMockUsers]))
  const products = safeState<Product[]>('mock-ds-products', () => (isApiMode ? [] : [...initialMockProducts]))
  const orders = safeState<MarketplaceOrder[]>('mock-ds-orders', () => (isApiMode ? [] : [...initialMockOrders]))
  const cart = safeState<CartItem[]>('mock-ds-cart', () => (isApiMode ? [] : [...initialMockCart]))
  const shop = safeState<Shop | null>('mock-ds-shop', () => null)
  const user = safeState<User>('mock-ds-user', () => ({
    username: '',
    name: '',
    email: '',
    role: 'buyer'
  }))
  const isLoggedIn = safeState<boolean>('mock-ds-is-logged-in', () => false)
  const reviews = safeState<ProductReview[]>('mock-ds-reviews', () => (isApiMode ? [] : [...initialMockReviews]))

  if (typeof window !== 'undefined' && !isClientInitialized) {
    isClientInitialized = true

    // Restore stored values if available in browser localStorage
    const savedProducts = getItemFromStorage<Product[] | null>(STORAGE_KEYS.PRODUCTS, null)
    if (savedProducts && Array.isArray(savedProducts)) {
      const userProducts = savedProducts.filter(p => typeof p.id === 'number' ? p.id > 12 : true)
      products.value = userProducts
      setItemInStorage(STORAGE_KEYS.PRODUCTS, userProducts)
    } else {
      products.value = []
      setItemInStorage(STORAGE_KEYS.PRODUCTS, [])
    }

    const savedShop = getItemFromStorage<Shop | null>(STORAGE_KEYS.SHOP, null)
    if (savedShop) {
      shop.value = savedShop
    }

    const savedOrders = getItemFromStorage<MarketplaceOrder[] | null>(STORAGE_KEYS.ORDERS, null)
    if (savedOrders && Array.isArray(savedOrders)) {
      orders.value = savedOrders
    }

    const savedCart = getItemFromStorage<CartItem[] | null>(STORAGE_KEYS.CART, null)
    if (savedCart && Array.isArray(savedCart)) {
      const validCart = savedCart.filter(item => products.value.some(p => String(p.id) === String(item.productId)))
      cart.value = validCart
      setItemInStorage(STORAGE_KEYS.CART, validCart)
    } else {
      cart.value = []
      setItemInStorage(STORAGE_KEYS.CART, [])
    }

    const savedUser = getItemFromStorage<User | null>(STORAGE_KEYS.USER, null)
    if (savedUser && savedUser.username) {
      user.value = savedUser
    }

    const savedLoggedIn = getItemFromStorage<boolean | null>(STORAGE_KEYS.IS_LOGGED_IN, null)
    if (savedLoggedIn !== null) {
      isLoggedIn.value = savedLoggedIn
    }

    const savedReviews = getItemFromStorage<ProductReview[] | null>(STORAGE_KEYS.REVIEWS, null)
    if (savedReviews && Array.isArray(savedReviews)) {
      reviews.value = savedReviews
    }

    // Persist changes to localStorage automatically
    watch(products, (val) => setItemInStorage(STORAGE_KEYS.PRODUCTS, val), { deep: true })
    watch(shop, (val) => setItemInStorage(STORAGE_KEYS.SHOP, val), { deep: true })
    watch(orders, (val) => setItemInStorage(STORAGE_KEYS.ORDERS, val), { deep: true })
    watch(cart, (val) => setItemInStorage(STORAGE_KEYS.CART, val), { deep: true })
    watch(user, (val) => setItemInStorage(STORAGE_KEYS.USER, val), { deep: true })
    watch(isLoggedIn, (val) => setItemInStorage(STORAGE_KEYS.IS_LOGGED_IN, val))
    watch(reviews, (val) => setItemInStorage(STORAGE_KEYS.REVIEWS, val), { deep: true })

    // Synchronize changes across multiple open tabs
    if (typeof window.addEventListener === 'function') {
      window.addEventListener('storage', (e) => {
        if (!e.key || e.newValue === null) return
        try {
          if (e.key === STORAGE_KEYS.PRODUCTS) {
            products.value = JSON.parse(e.newValue)
          } else if (e.key === STORAGE_KEYS.SHOP) {
            shop.value = JSON.parse(e.newValue)
          } else if (e.key === STORAGE_KEYS.ORDERS) {
            orders.value = JSON.parse(e.newValue)
          } else if (e.key === STORAGE_KEYS.CART) {
            cart.value = JSON.parse(e.newValue)
          } else if (e.key === STORAGE_KEYS.REVIEWS) {
            reviews.value = JSON.parse(e.newValue)
          } else if (e.key === STORAGE_KEYS.USER) {
            user.value = JSON.parse(e.newValue)
          } else if (e.key === STORAGE_KEYS.IS_LOGGED_IN) {
            isLoggedIn.value = JSON.parse(e.newValue)
          }
        } catch {
          // ignore invalid JSON parsing
        }
      })
    }
  }

  const addReview = (productId: number | string, rating: number, comment: string, authorName?: string, orderId?: number | string) => {
    const author = authorName || user.value.name || user.value.username || 'Verified Buyer'

    // Prevent duplicate review for the same product by the same author
    const existing = reviews.value.find(
      r => String(r.productId) === String(productId) && r.author === author && (orderId ? String(r.orderId) === String(orderId) : true)
    )
    if (existing) {
      return existing
    }

    const newReview: ProductReview = {
      id: Date.now(),
      productId: typeof productId === 'number' ? productId : (parseInt(String(productId), 10) || Date.now()),
      orderId,
      author,
      rating,
      comment,
      date: new Date().toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' }),
      status: 'pending'
    }
    reviews.value.unshift(newReview)
    return newReview
  }

  return {
    users,
    products,
    orders,
    cart,
    shop,
    user,
    isLoggedIn,
    reviews,
    addReview
  }
}


