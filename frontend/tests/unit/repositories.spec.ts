import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MockProductRepository } from '../../app/repositories/mock/MockProductRepository'
import { MockShopRepository } from '../../app/repositories/mock/MockShopRepository'
import { MockOrderRepository } from '../../app/repositories/mock/MockOrderRepository'
import { MockAuthRepository } from '../../app/repositories/mock/MockAuthRepository'
import { ApiAuthRepository } from '../../app/repositories/api/ApiAuthRepository'
import { useMockDataStore } from '../../app/repositories/mock/MockDataStore'
import { useMarketplace } from '../../app/composables/useMarketplace'
import { STATIC_CATALOG } from '../../app/utils/categoryCatalog'

describe('Repository Layer Unit Tests', () => {
  const productRepo = new MockProductRepository()
  const shopRepo = new MockShopRepository()
  const orderRepo = new MockOrderRepository()
  const authRepo = new MockAuthRepository()

  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('ProductRepository: should query and filter products', async () => {
    await productRepo.createProduct('Atelier North', {
      name: 'Men Shirt',
      description: 'A test shirt for men',
      category: 'Men',
      price: 1500,
      stock: 10,
      status: 'Active'
    })
    const res = await productRepo.getProducts({ category: 'Men' })
    expect(res.data).toBeDefined()
    expect(res.data.every(p => p.category === 'Men')).toBe(true)
  })

  it('ProductRepository: should query products by shop', async () => {
    await productRepo.createProduct('Atelier North', {
      name: 'Shop Specific Shirt',
      description: 'A test shirt for Atelier North',
      category: 'Men',
      price: 1800,
      stock: 7,
      status: 'Active'
    })
    const shopProducts = await productRepo.getProductsByShop('Atelier North')
    expect(shopProducts.length).toBeGreaterThan(0)
    expect(shopProducts.some(p => p.name === 'Shop Specific Shirt')).toBe(true)
  })

  it('ProductRepository: should update and toggle product status', async () => {
    const created = await productRepo.createProduct('Atelier North', {
      name: 'Sample Item',
      description: 'Sample description',
      category: 'Women',
      price: 2000,
      stock: 5,
      status: 'Active'
    })
    const product = await productRepo.getProductById(created.id)
    expect(product).not.toBeNull()

    if (product) {
      const updated = await productRepo.updateProduct(created.id, { price: 1999 })
      expect(updated?.price).toBe(1999)

      const toggled = await productRepo.toggleProductStatus(created.id)
      expect(toggled?.status).toBe('Draft')
    }
  })

  it('ShopRepository: should retrieve and update seller shop profile', async () => {
    const newShop = await shopRepo.createShop('seller@afrimart.com', {
      name: 'Atelier North',
      description: 'Handcrafted goods'
    })
    expect(newShop.slug).toBe('atelier-north')

    const shop = await shopRepo.getShopBySlug('atelier-north')
    expect(shop).not.toBeNull()
    expect(shop?.name).toBe('Atelier North')

    const updatedShop = await shopRepo.updateShop('atelier-north', {
      description: 'Updated handcrafted goods description'
    })
    expect(updatedShop?.description).toBe('Updated handcrafted goods description')
  })

  it('OrderRepository: should create order and update delivery status', async () => {
    const created = await productRepo.createProduct('Atelier North', {
      name: 'Order Item',
      description: 'Desc',
      category: 'Men',
      price: 1400,
      stock: 10,
      status: 'Active'
    })
    const newOrder = await orderRepo.createOrder(
      [{ productId: created.id, quantity: 2 }],
      2800,
      {
        buyerName: 'Jane Doe',
        deliveryAddress: 'Kazanchis, Addis Ababa',
        phone: '+251911223344',
        paymentMethod: 'Telebirr'
      }
    )

    expect(newOrder.id).toBeDefined()
    expect(newOrder.buyerName).toBe('Jane Doe')
    expect(newOrder.status).toBe('Ordered')

    const updatedOrder = await orderRepo.updateOrderStatus(newOrder.id, 'Shipped')
    expect(updatedOrder?.status).toBe('Shipped')
  })

  it('AuthRepository: should authenticate login and registration', async () => {
    const session = await authRepo.login({
      email: 'john@example.com',
      password: 'password123'
    })

    expect(session.token).toContain('mock-jwt-token-')
    expect(session.user.email).toBe('john@example.com')

    await authRepo.logout()
    const currentSession = await authRepo.getCurrentSession()
    expect(currentSession).toBeNull()
  })

  it('ProductReviews: should save new review as pending and prevent duplicate submissions', async () => {
    const store = useMockDataStore()
    const created = await productRepo.createProduct('Atelier North', {
      name: 'Review Product',
      description: 'Desc',
      category: 'Men',
      price: 1000,
      stock: 5,
      status: 'Active'
    })

    const review1 = store.addReview(created.id, 5, 'Great quality!', 'Test Reviewer', 240824)
    expect(review1.status).toBe('pending')
    expect(review1.rating).toBe(5)

    // Attempting duplicate review for same order item & author returns existing review
    const review2 = store.addReview(created.id, 4, 'Duplicate review attempt', 'Test Reviewer', 240824)
    expect(review2.id).toBe(review1.id)
    expect(review2.comment).toBe('Great quality!')
  })

  it('ApiAuthRepository: should handle successful login and store token', async () => {
    const apiAuthRepo = new ApiAuthRepository()
    const mockResponse = {
      id: 101,
      email: 'backend@afrimart.com',
      token: 'jwt-access-token-123',
      refresh_token: 'jwt-refresh-token-456'
    }

    vi.stubGlobal('$fetch', vi.fn().mockResolvedValue(mockResponse))
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

    const session = await apiAuthRepo.login({ email: 'backend@afrimart.com', password: 'secretpassword' })
    expect(session.token).toBe('jwt-access-token-123')
    expect(session.user.email).toBe('backend@afrimart.com')

    const { isLoggedIn } = useMockDataStore()
    expect(isLoggedIn.value).toBe(true)

    await apiAuthRepo.logout()
    expect(isLoggedIn.value).toBe(false)
  })

  it('ApiAuthRepository: should handle registration without auto-login if token is omitted', async () => {
    const apiAuthRepo = new ApiAuthRepository()
    const mockRegisterRes = { id: 102, email: 'newuser@afrimart.com' }

    vi.stubGlobal('$fetch', vi.fn().mockResolvedValue(mockRegisterRes))
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

    const session = await apiAuthRepo.register({
      username: 'newuser',
      firstName: 'New',
      lastName: 'User',
      email: 'newuser@afrimart.com',
      password: 'password123'
    })

    expect(session.user.email).toBe('newuser@afrimart.com')
    expect(session.token).toBe('')

    const { isLoggedIn } = useMockDataStore()
    expect(isLoggedIn.value).toBe(false)
  })

  it('ApiAuthRepository: 401 Unauthorized must throw error and keep user logged out', async () => {
    const apiAuthRepo = new ApiAuthRepository()
    const error401 = { data: { message: 'Invalid credentials' }, statusCode: 401 }

    vi.stubGlobal('$fetch', vi.fn().mockRejectedValue(error401))
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

    await expect(apiAuthRepo.login({ email: 'wrong@example.com', password: 'wrong' }))
      .rejects.toThrow('Invalid credentials')

    const { isLoggedIn } = useMockDataStore()
    expect(isLoggedIn.value).toBe(false)
  })

  it('ApiAuthRepository: 400 Validation Error must throw error and keep user logged out', async () => {
    const apiAuthRepo = new ApiAuthRepository()
    const error400 = { data: { message: 'Email already exists' }, statusCode: 400 }

    vi.stubGlobal('$fetch', vi.fn().mockRejectedValue(error400))
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

    await expect(apiAuthRepo.register({
      username: 'existing',
      firstName: 'Test',
      lastName: 'User',
      email: 'existing@afrimart.com',
      password: 'short'
    })).rejects.toThrow('Email already exists')

    const { isLoggedIn } = useMockDataStore()
    expect(isLoggedIn.value).toBe(false)
  })

  it('ApiAuthRepository: Network failure must throw error and keep user logged out', async () => {
    const apiAuthRepo = new ApiAuthRepository()

    vi.stubGlobal('$fetch', vi.fn().mockRejectedValue(new Error('Network connection failed')))
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

    await expect(apiAuthRepo.login({ email: 'test@example.com', password: 'password' }))
      .rejects.toThrow('Network connection failed')

    const { isLoggedIn } = useMockDataStore()
    expect(isLoggedIn.value).toBe(false)
  })

  it('Self-Purchase Protection & Sold Out Validation: should block own products using the real shop ID', async () => {
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { authMode: 'mock' } }))
    const store = useMockDataStore()
    store.shop.value = {
      id: 'shop-101',
      backendId: 'shop-101',
      name: 'Atelier North',
      slug: 'atelier-north',
      description: 'Test shop',
      ownerEmail: 'seller@example.com',
      products: [],
      paymentMethods: [],
      status: 'active'
    }

    const ownProd = { id: 1, shop: 'Different display name', shopId: 'shop-101', stock: 5 } as any
    const otherProd = { id: 2, shop: 'Atelier North', shopId: 'shop-202', stock: 5 } as any
    const soldOutProd = { id: 3, shop: 'Urban Thread', stock: 0 } as any
    store.products.value = [ownProd, otherProd, soldOutProd]

    const { isOwnProduct, addToCart, filterProducts } = useMarketplace()

    expect(isOwnProduct(ownProd)).toBe(true)
    expect(isOwnProduct(otherProd)).toBe(false)
    expect(isOwnProduct(soldOutProd)).toBe(false)
    await expect(addToCart(ownProd.id)).rejects.toThrow('You cannot purchase items listed by your own shop.')
    expect(store.cart.value).toHaveLength(0)

    const clothing = [
      { ...ownProd, id: 'men-item', name: 'Men item', category: 'Clothing', gender: 'Men', status: 'Active' },
      { ...otherProd, id: 'women-item', name: 'Women item', category: 'Clothing', gender: 'Women', status: 'Active' },
      { ...soldOutProd, id: 'kids-item', name: 'Kids item', category: 'Clothing', gender: 'Kids', status: 'Active' }
    ] as any
    expect(filterProducts({ category: 'Men' }, clothing).map(product => product.name)).toEqual(['Men item'])
    expect(filterProducts({ category: 'Women' }, clothing).map(product => product.name)).toEqual(['Women item'])
    expect(filterProducts({ category: 'Kids' }, clothing).map(product => product.name)).toEqual(['Kids item'])
  })

  it('Marketplace filters compose using catalog IDs and keep public products available while logged out', () => {
    const store = useMockDataStore()
    store.shop.value = null
    store.isLoggedIn.value = false

    const menCatalog = STATIC_CATALOG.find(category => category.name === 'Men')!
    const womenCatalog = STATIC_CATALOG.find(category => category.name === 'Women')!
    const shoesCatalog = STATIC_CATALOG.find(category => category.name === 'Shoes')!
    const shirtId = menCatalog.subcategories.find(subcategory => subcategory.name === 'Shirts')?.id
    const hoodieId = menCatalog.subcategories.find(subcategory => subcategory.name === 'Hoodies')?.id
    const dressId = womenCatalog.subcategories.find(subcategory => subcategory.name === 'Dresses')?.id
    const sneakerId = shoesCatalog.subcategories.find(subcategory => subcategory.name === 'Sneakers')?.id
    const catalogProducts = [
      { id: 'men-shirt', shop: 'North', name: 'Linen shirt', description: '', category: 'Clothing', categoryId: menCatalog.id, gender: 'Men', subCategory: 'Shirts', subcategoryId: shirtId, price: 120, stock: 4, rating: 'New', image: '', status: 'Active' },
      { id: 'men-hoodie', shop: 'North', name: 'Cotton hoodie', description: '', category: 'Clothing', categoryId: menCatalog.id, gender: 'Men', subCategory: 'Hoodies', subcategoryId: hoodieId, price: 260, stock: 0, rating: 'New', image: '', status: 'Active' },
      { id: 'women-dress', shop: 'South', name: 'Summer dress', description: '', category: 'Clothing', categoryId: womenCatalog.id, gender: 'Women', subCategory: 'Dresses', subcategoryId: dressId, price: 240, stock: 2, rating: 'New', image: '', status: 'Active' },
      { id: 'sneakers', shop: 'East', name: 'Canvas sneakers', description: '', category: 'Shoes', categoryId: shoesCatalog.id, subCategory: 'Sneakers', subcategoryId: sneakerId, price: 90, stock: 3, rating: 'New', image: '', status: 'Active' },
      { id: 'draft-shirt', shop: 'North', name: 'Draft shirt', description: '', category: 'Clothing', categoryId: menCatalog.id, gender: 'Men', subCategory: 'Shirts', subcategoryId: shirtId, price: 80, stock: 1, rating: 'New', image: '', status: 'Draft' }
    ] as any

    const { filterProducts } = useMarketplace()

    expect(filterProducts({ category: 'Men', subCategory: 'Shirts' }, catalogProducts).map(product => product.id)).toEqual(['men-shirt'])
    expect(filterProducts({ category: 'Men', subCategory: 'Shirts', subcategoryId: shirtId }, catalogProducts).map(product => product.id)).toEqual(['men-shirt'])
    expect(filterProducts({ category: 'All', includeInactive: true }, catalogProducts)).toHaveLength(5)
    expect(filterProducts({ category: 'All', availability: 'available' }, catalogProducts).map(product => product.id)).toEqual(['men-shirt', 'women-dress', 'sneakers'])
    expect(filterProducts({ category: 'All', availability: 'sold-out' }, catalogProducts).map(product => product.id)).toEqual(['men-hoodie'])
    expect(filterProducts({ category: 'All', minPrice: 100, maxPrice: 240, sortBy: 'price-asc' }, catalogProducts).map(product => product.id)).toEqual(['men-shirt', 'women-dress'])
    expect(filterProducts({ category: 'All', sortBy: 'price-desc' }, catalogProducts).map(product => product.id)).toEqual(['men-hoodie', 'women-dress', 'men-shirt', 'sneakers'])
    expect(filterProducts({ category: 'Women', search: 'summer' }, catalogProducts).map(product => product.id)).toEqual(['women-dress'])
    expect(filterProducts({ category: 'All' }, catalogProducts).map(product => product.id)).toHaveLength(4)
  })
})
