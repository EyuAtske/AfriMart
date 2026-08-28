import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MockProductRepository } from '../../app/repositories/mock/MockProductRepository'
import { MockShopRepository } from '../../app/repositories/mock/MockShopRepository'
import { MockOrderRepository } from '../../app/repositories/mock/MockOrderRepository'
import { MockAuthRepository } from '../../app/repositories/mock/MockAuthRepository'
import { ApiAuthRepository } from '../../app/repositories/api/ApiAuthRepository'
import { useMockDataStore } from '../../app/repositories/mock/MockDataStore'

describe('Repository Layer Unit Tests', () => {
  const productRepo = new MockProductRepository()
  const shopRepo = new MockShopRepository()
  const orderRepo = new MockOrderRepository()
  const authRepo = new MockAuthRepository()

  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('ProductRepository: should query and filter products', async () => {
    const res = await productRepo.getProducts({ category: 'Men' })
    expect(res.data).toBeDefined()
    expect(res.data.every(p => p.category === 'Men')).toBe(true)
  })

  it('ProductRepository: should update and toggle product status', async () => {
    const product = await productRepo.getProductById(1)
    expect(product).not.toBeNull()

    if (product) {
      const updated = await productRepo.updateProduct(1, { price: 1999 })
      expect(updated?.price).toBe(1999)

      const toggled = await productRepo.toggleProductStatus(1)
      expect(toggled?.status).toBe('Draft')
    }
  })

  it('ShopRepository: should retrieve and update seller shop profile', async () => {
    const shop = await shopRepo.getShopBySlug('atelier-north')
    expect(shop).not.toBeNull()
    expect(shop?.name).toBe('Atelier North')

    const newShop = await shopRepo.createShop('seller@afrimart.com', {
      name: 'Test Artisan Studio',
      description: 'Handcrafted goods'
    })
    expect(newShop.slug).toBe('test-artisan-studio')

    const updatedShop = await shopRepo.updateShop('test-artisan-studio', {
      description: 'Updated handcrafted goods description'
    })
    expect(updatedShop?.description).toBe('Updated handcrafted goods description')
  })

  it('OrderRepository: should create order and update delivery status', async () => {
    const newOrder = await orderRepo.createOrder(
      [{ productId: 1, quantity: 2 }],
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

  it('ProductReviews: should save new review as pending and prevent duplicate submissions', () => {
    const store = useMockDataStore()

    const review1 = store.addReview(1, 5, 'Great quality!', 'Test Reviewer', 240824)
    expect(review1.status).toBe('pending')
    expect(review1.rating).toBe(5)

    // Attempting duplicate review for same order item & author returns existing review
    const review2 = store.addReview(1, 4, 'Duplicate review attempt', 'Test Reviewer', 240824)
    expect(review2.id).toBe(review1.id)
    expect(review2.comment).toBe('Great quality!')
  })

  it('ApiAuthRepository: should handle successful login and store token in memory', async () => {
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
    expect(apiAuthRepo.getMemoryToken()).toBe('jwt-access-token-123')

    const { isLoggedIn } = useMockDataStore()
    expect(isLoggedIn.value).toBe(true)

    await apiAuthRepo.logout()
    expect(apiAuthRepo.getMemoryToken()).toBeNull()
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
    expect(apiAuthRepo.getMemoryToken()).toBeNull()

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
    expect(apiAuthRepo.getMemoryToken()).toBeNull()
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
    expect(apiAuthRepo.getMemoryToken()).toBeNull()
  })

  it('ApiAuthRepository: Network failure must throw error and keep user logged out', async () => {
    const apiAuthRepo = new ApiAuthRepository()

    vi.stubGlobal('$fetch', vi.fn().mockRejectedValue(new Error('Network connection failed')))
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

    await expect(apiAuthRepo.login({ email: 'test@example.com', password: 'password' }))
      .rejects.toThrow('Network connection failed')

    const { isLoggedIn } = useMockDataStore()
    expect(isLoggedIn.value).toBe(false)
    expect(apiAuthRepo.getMemoryToken()).toBeNull()
  })
})
