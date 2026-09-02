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
    price: 2600,
    stock: 7,
    rating: '4.6',
    image: '/images/product3.jpg',
    status: 'Active'
  }
]

export const initialMockOrders: MarketplaceOrder[] = [
  {
    id: 240824,
    buyerName: 'Test User',
    items: [
      {
        productId: 1,
        quantity: 1
      }
    ],
    deliveryAddress: 'Bole, Addis Ababa',
    phone: '+251 911 000 000',
    paymentMethod: 'Cash on delivery',
    paymentStatus: 'Cash on delivery',
    status: 'Shipped',
    date: 'August 24, 2026',
    total: 1400
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
    quantity: 1
  }
]

export const useMockDataStore = () => {
  const products = safeState<Product[]>('mock-ds-products', () => [...initialMockProducts])
  const orders = safeState<MarketplaceOrder[]>('mock-ds-orders', () => [...initialMockOrders])
  const cart = safeState<CartItem[]>('mock-ds-cart', () => [...initialMockCart])
  const shop = safeState<Shop | null>('mock-ds-shop', () => null)
  const user = safeState<User>('mock-ds-user', () => ({
    username: '',
    name: '',
    email: '',
    role: 'buyer'
  }))
  const isLoggedIn = safeState<boolean>('mock-ds-is-logged-in', () => false)
  const reviews = safeState<ProductReview[]>('mock-ds-reviews', () => [...initialMockReviews])

  const addReview = (productId: number, rating: number, comment: string, authorName?: string, orderId?: number) => {
    const author = authorName || user.value.name || user.value.username || 'Verified Buyer'

    // Prevent duplicate review for the same product by the same author
    const existing = reviews.value.find(
      r => r.productId === productId && r.author === author && (orderId ? r.orderId === orderId : true)
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
