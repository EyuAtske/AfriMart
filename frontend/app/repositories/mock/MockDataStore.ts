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

export const initialMockProducts: Product[] = [
  {
    id: 1,
    shop: 'Atelier North',
    name: 'Shirt for men',
    description: 'A clean everyday shirt with a relaxed fit and soft cotton feel.',
    category: 'Men',
    subCategory: 'Shirts',
    price: 1400,
    stock: 12,
    rating: '4.8',
    image: '/images/product1.jpg',
    status: 'Active'
  },
  {
    id: 2,
    shop: 'Beyond Score',
    name: 'Tank Tops',
    description: 'Lightweight tank tops made for warm days and easy layering.',
    category: 'Women',
    subCategory: 'Tops',
    price: 2500,
    stock: 8,
    rating: '4.9',
    image: '/images/product2.jpg',
    status: 'Active'
  },
  {
    id: 3,
    shop: 'True Form',
    name: 'Casual Outfit Set',
    description: 'Matched casual set with a tidy silhouette for daily wear.',
    category: 'Women',
    subCategory: 'Dresses',
    price: 3500,
    stock: 6,
    rating: '4.6',
    image: '/images/product3.jpg',
    status: 'Active'
  },
  {
    id: 4,
    shop: 'Minimal Studio',
    name: 'Everyday fit for kids',
    description: 'Comfortable kids outfit built for school days and weekends.',
    category: 'Kids',
    subCategory: 'T-Shirts',
    price: 2900,
    stock: 14,
    rating: '4.8',
    image: '/images/product4.jpg',
    status: 'Active'
  },
  {
    id: 5,
    shop: 'Atelier North',
    name: 'Cute dress for kids',
    description: 'Soft dress with a cheerful cut and easy movement.',
    category: 'Kids',
    subCategory: 'Dresses',
    price: 3800,
    stock: 5,
    rating: '4.7',
    image: '/images/product5.jpg',
    status: 'Active'
  },
  {
    id: 6,
    shop: 'Urban Thread',
    name: 'Hoodie',
    description: 'Warm hoodie with a soft inner layer and simple streetwear shape.',
    category: 'Men',
    subCategory: 'Hoodies',
    price: 4700,
    stock: 9,
    rating: '4.8',
    image: '/images/product6.jpg',
    status: 'Active'
  },
  {
    id: 7,
    shop: 'Mara Studio',
    name: 'Classic cotton shirt',
    description: 'Crisp shirt with a polished collar and breathable fabric.',
    category: 'Men',
    subCategory: 'Shirts',
    price: 1550,
    stock: 11,
    rating: '4.9',
    image: '/images/product7.jpg',
    status: 'Active'
  },
  {
    id: 8,
    shop: 'Saba Edit',
    name: 'Relaxed Denim',
    description: 'Easy denim piece with a flattering relaxed shape.',
    category: 'Women',
    subCategory: 'Jeans',
    price: 5200,
    stock: 4,
    rating: '4.7',
    image: '/images/product8.jpg',
    status: 'Active'
  },
  {
    id: 9,
    shop: 'Minimal Studio',
    name: 'Watch for women',
    description: 'A slim everyday watch with a clean face and subtle finish.',
    category: 'Accessories',
    subCategory: 'Watches',
    price: 2500,
    stock: 7,
    rating: '4.9',
    image: '/images/product9.jpg',
    status: 'Active'
  },
  {
    id: 10,
    shop: 'Mara Studio',
    name: 'Hat',
    description: 'Simple everyday hat for sun coverage and finishing an outfit.',
    category: 'Accessories',
    subCategory: 'Hats',
    price: 1100,
    stock: 16,
    rating: '4.8',
    image: '/images/product10.jpg',
    status: 'Active'
  },
  {
    id: 11,
    shop: 'Urban Thread',
    name: 'Clean everyday sneakers',
    description: 'Low-profile sneakers that pair easily with relaxed denim and casual outfits.',
    category: 'Shoes',
    subCategory: 'Sneakers',
    price: 4300,
    stock: 10,
    rating: '4.7',
    image: '/images/product6.jpg',
    status: 'Active'
  },
  {
    id: 12,
    shop: 'Saba Edit',
    name: 'Soft city sandals',
    description: 'Comfortable sandals for warm days, errands, and weekend styling.',
    category: 'Shoes',
    subCategory: 'Sandals',
    price: 2600,
    stock: 7,
    rating: '4.6',
    image: '/images/product3.jpg',
    status: 'Active'
  }
]

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

import { ref, type Ref } from 'vue'

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

export const initialMockCart: CartItem[] = [
  {
    productId: 1,
    quantity: 1,
    backendItemId: 'mock-item-1'
  }
]

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

  const addReview = (productId: number, rating: number, comment: string, authorName?: string, orderId?: number | string) => {
    const author = authorName || user.value.name || user.value.username || 'Verified Buyer'

    // Prevent duplicate review for the same product by the same author
    const existing = reviews.value.find(
      r => r.productId === productId && r.author === author && (orderId ? String(r.orderId) === String(orderId) : true)
    )
    if (existing) {
      return existing
    }

    const newReview: ProductReview = {
      id: Date.now(),
      productId,
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

