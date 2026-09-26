import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ApiAuthRepository } from '../../app/repositories/api/ApiAuthRepository'
import { ApiShopRepository } from '../../app/repositories/api/ApiShopRepository'
import { ApiProductRepository } from '../../app/repositories/api/ApiProductRepository'
import { ApiCartRepository } from '../../app/repositories/api/ApiCartRepository'
import { ApiOrderRepository } from '../../app/repositories/api/ApiOrderRepository'
import { useMockDataStore } from '../../app/repositories/mock/MockDataStore'
import * as apiHelpers from '../../app/repositories/api/apiHelpers'

describe('Real Data Frontend API Repositories (Auth, Shop, Product)', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
    apiHelpers.clearTokens()
    const { isLoggedIn, user, shop } = useMockDataStore()
    isLoggedIn.value = false
    user.value = {
      name: 'Guest User',
      email: 'guest@afrimart.com',
      avatar: '/images/default-avatar.png',
      role: 'buyer'
    }
    shop.value = null
  })

  describe('ApiAuthRepository', () => {
    const authRepo = new ApiAuthRepository()

    it('login: should authenticate with real backend API format and create user session', async () => {
      const mockLoginResponse = {
        id: 'user-uuid-101',
        email: 'seller@afrimart.com',
        username: 'afri_seller',
        token: 'jwt-access-token-xyz',
        refresh_token: 'jwt-refresh-token-123'
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/api/auth/login')) {
          return Promise.resolve(mockLoginResponse)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const session = await authRepo.login({
        email: 'seller@afrimart.com',
        password: 'securepassword123'
      })

      expect(session.token).toBe('jwt-access-token-xyz')
      expect(session.user.email).toBe('seller@afrimart.com')

      const { isLoggedIn } = useMockDataStore()
      expect(isLoggedIn.value).toBe(true)
    })

    it('register: should call POST /api/auth/register with user payload', async () => {
      const mockRegisterResponse = {
        id: 'user-uuid-102',
        email: 'newbuyer@afrimart.com',
        username: 'newbuyer'
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/api/auth/register')) {
          return Promise.resolve(mockRegisterResponse)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const session = await authRepo.register({
        username: 'newbuyer',
        firstName: 'New',
        lastName: 'Buyer',
        email: 'newbuyer@afrimart.com',
        password: 'password123'
      })

      expect(session.user.email).toBe('newbuyer@afrimart.com')
    })

    it('login failure (401 Unauthorized): should throw error and remain logged out', async () => {
      const error401 = { data: { message: 'Invalid email or password' }, statusCode: 401 }

      vi.stubGlobal('$fetch', vi.fn().mockRejectedValue(error401))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      await expect(
        authRepo.login({ email: 'wrong@example.com', password: 'badpassword' })
      ).rejects.toThrow('Invalid email or password')

      const { isLoggedIn } = useMockDataStore()
      expect(isLoggedIn.value).toBe(false)
    })

    it('logout: should invoke POST /api/auth/logout and clear session', async () => {
      apiHelpers.setAccessToken('active-access-token')
      const { isLoggedIn } = useMockDataStore()
      isLoggedIn.value = true

      vi.stubGlobal('$fetch', vi.fn().mockResolvedValue({ status: 'ok' }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      await authRepo.logout()

      expect(isLoggedIn.value).toBe(false)
    })

    it('authenticatedFetch: should sanitize endpoint with Markdown brackets or parentheses', async () => {
      apiHelpers.setAccessToken('active-access-token')
      let calledUrl = ''

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        calledUrl = url
        return Promise.resolve({ ok: true })
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      await apiHelpers.authenticatedFetch('[api]/(user)/profile')
      expect(calledUrl).toBe('http://localhost:8080/api/user/profile')
    })
  })

  describe('ApiShopRepository', () => {
    const shopRepo = new ApiShopRepository()

    const mockBackendShop = {
      ID: 'shop-uuid-888',
      OwnerID: 'user-uuid-101',
      Name: 'Addis Craft Market',
      Description: { String: 'Traditional Ethiopian crafts and textiles', Valid: true },
      Status: 'active',
      CreatedAt: '2026-09-01T10:00:00Z',
      UpdatedAt: '2026-09-01T10:00:00Z'
    }

    it('getMyShop: should return null when not authenticated', async () => {
      apiHelpers.clearTokens()
      const shop = await shopRepo.getMyShop()
      expect(shop).toBeNull()
    })

    it('getMyShop: should fetch and map Go backend PascalCase shop response when authenticated', async () => {
      apiHelpers.setAccessToken('valid-seller-token')
      let capturedHeader = ''

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/shops/me')) {
          capturedHeader = opts?.headers?.Authorization || ''
          return Promise.resolve(mockBackendShop)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const shop = await shopRepo.getMyShop()

      expect(shop).not.toBeNull()
      expect(shop?.id).toBe('shop-uuid-888')
      expect(shop?.name).toBe('Addis Craft Market')
      expect(shop?.description).toBe('Traditional Ethiopian crafts and textiles')
      expect(shop?.status).toBe('active')
      expect(capturedHeader).toBe('Bearer valid-seller-token')
    })

    it('getMyShop: should display clear login message on 401 response', async () => {
      apiHelpers.setAccessToken('expired-or-invalid-token')

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/api/shops/me')) {
          const err: any = new Error('Unauthorized')
          err.statusCode = 401
          return Promise.reject(err)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      await expect(shopRepo.getMyShop()).rejects.toThrow('Please log in to access your seller shop.')
    })

    it('createShop: should POST /api/shops and update user role to seller', async () => {
      const mockCreatedShop = {
        ID: 'shop-uuid-999',
        OwnerID: 'user-uuid-101',
        Name: 'AfriMart Artisan Studio',
        Description: { String: 'Handcrafted goods from East Africa', Valid: true },
        Status: 'active',
        CreatedAt: '2026-09-18T10:00:00Z',
        UpdatedAt: '2026-09-18T10:00:00Z'
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/shops') && opts?.method === 'POST') {
          return Promise.resolve(mockCreatedShop)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const newShop = await shopRepo.createShop('seller@afrimart.com', {
        name: 'AfriMart Artisan Studio',
        description: 'Handcrafted goods from East Africa'
      })

      expect(newShop.id).toBe('shop-uuid-999')
      expect(newShop.name).toBe('AfriMart Artisan Studio')
      
      const { user } = useMockDataStore()
      expect(user.value.role).toBe('seller')
    })

    it('updateShop: should PATCH shop name and description endpoints', async () => {
      const mockUpdatedShop = {
        ID: 'shop-uuid-999',
        OwnerID: 'user-uuid-101',
        Name: 'Updated Studio Name',
        Description: { String: 'Updated studio description', Valid: true },
        Status: 'active',
        CreatedAt: '2026-09-18T10:00:00Z',
        UpdatedAt: '2026-09-18T11:00:00Z'
      }

      const { shop } = useMockDataStore()
      shop.value = {
        id: 'shop-uuid-999',
        name: 'Old Studio Name',
        slug: 'old-studio-name',
        description: 'Old desc',
        ownerEmail: 'seller@afrimart.com',
        products: [],
        paymentMethods: [],
        status: 'active',
        backendId: 'shop-uuid-999'
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/api/shops/shop-uuid-999/name') || url.includes('/api/shops/shop-uuid-999/description')) {
          return Promise.resolve(mockUpdatedShop)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const updated = await shopRepo.updateShop('old-studio-name', {
        name: 'Updated Studio Name',
        description: 'Updated studio description'
      })

      expect(updated).not.toBeNull()
      expect(updated?.name).toBe('Updated Studio Name')
      expect(updated?.description).toBe('Updated studio description')
    })

    it('activateShop and deactivateShop: should PATCH shop status endpoints', async () => {
      const mockDeactivatedShop = {
        ID: 'shop-uuid-999',
        OwnerID: 'user-uuid-101',
        Name: 'Addis Market',
        Description: { String: 'Desc', Valid: true },
        Status: 'deactivated',
        CreatedAt: '2026-09-01T10:00:00Z',
        UpdatedAt: '2026-09-18T10:00:00Z'
      }

      const mockActivatedShop = {
        ...mockDeactivatedShop,
        Status: 'active'
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/deactivate')) return Promise.resolve(mockDeactivatedShop)
        if (url.includes('/activate')) return Promise.resolve(mockActivatedShop)
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const deactivated = await shopRepo.deactivateShop('shop-uuid-999')
      expect(deactivated?.status).toBe('deactivated')

      const activated = await shopRepo.activateShop('shop-uuid-999')
      expect(activated?.status).toBe('active')
    })
  })

  describe('ApiProductRepository', () => {
    const productRepo = new ApiProductRepository()

    const mockBackendShop = {
      ID: 'shop-uuid-999',
      OwnerID: 'user-uuid-101',
      Name: 'AfriMart Artisan Studio',
      Description: { String: 'Studio desc', Valid: true },
      Status: 'active'
    }

    const mockBackendProduct = {
      ID: 'prod-uuid-555',
      ShopID: 'shop-uuid-999',
      CategoryID: '00000000-0000-0000-0000-000000000001',
      SubcategoryID: '00000000-0000-0000-0000-000000000002',
      Name: 'Handwoven Kente Scarf',
      Description: { String: 'Authentic handwoven silk and cotton Kente scarf', Valid: true },
      Price: '250.50',
      Stock: 15,
      Image: { String: '/images/kente.png', Valid: true },
      Status: 'active'
    }

    it('getProducts: should fetch products list from GET /api/products with search & pagination', async () => {
      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/api/products')) {
          return Promise.resolve([mockBackendProduct])
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const result = await productRepo.getProducts({ search: 'Kente', page: 1, pageSize: 10 })

      expect(result.data.length).toBe(1)
      expect(result.data[0].id).toBe('prod-uuid-555')
      expect(result.data[0].name).toBe('Handwoven Kente Scarf')
      expect(result.data[0].price).toBe(250.50)
      expect(result.data[0].stock).toBe(15)
    })

    it('getProducts: should map backend ProductWithImages response wrapper with media list', async () => {
      const wrappedProduct = {
        product: mockBackendProduct,
        images: [
          {
            ID: 'img-1',
            ProductID: 'prod-uuid-555',
            ObjectKey: 'https://images.example.com/scarf-1.jpg',
            DisplayOrder: 0
          },
          {
            ID: 'img-2',
            ProductID: 'prod-uuid-555',
            ObjectKey: 'products/prod-uuid-555/scarf-2.jpg',
            DisplayOrder: 1
          }
        ]
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/api/products')) {
          return Promise.resolve([wrappedProduct])
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const result = await productRepo.getProducts()

      expect(result.data.length).toBe(1)
      expect(result.data[0].id).toBe('prod-uuid-555')
      expect(result.data[0].image).toBe('https://images.example.com/scarf-1.jpg')
      expect(result.data[0].media?.length).toBe(2)
      // ObjectKey is mapped to accessible MinIO image URL:
      expect(result.data[0].media?.[1].url).toBe('http://localhost:9000/afrimart-images/products/prod-uuid-555/scarf-2.jpg')
    })

    it('getProducts: should throw clear error on API failure', async () => {
      vi.stubGlobal('$fetch', vi.fn().mockRejectedValue({
        statusCode: 500,
        message: 'Internal Server Error'
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      await expect(productRepo.getProducts()).rejects.toThrow('Something went wrong. Please try again later.')
    })

    it('getProductById: should fetch single product from GET /api/products/{id}', async () => {
      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/api/products/prod-uuid-555')) {
          return Promise.resolve(mockBackendProduct)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const product = await productRepo.getProductById('prod-uuid-555')

      expect(product).not.toBeNull()
      expect(product?.name).toBe('Handwoven Kente Scarf')
      expect(product?.price).toBe(250.50)
    })

    it('getProductById: should return null when backend returns 400 Invalid product ID', async () => {
      vi.stubGlobal('$fetch', vi.fn().mockRejectedValue({
        statusCode: 400,
        data: { error: 'Invalid product ID' }
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const product = await productRepo.getProductById('not-a-uuid')
      expect(product).toBeNull()
    })

    it('createProduct: should POST /api/products as multipart FormData and map response', async () => {
      const { shop } = useMockDataStore()
      shop.value = {
        id: 'shop-uuid-999',
        name: 'AfriMart Artisan Studio',
        slug: 'afrimart-artisan-studio',
        description: 'Studio desc',
        ownerEmail: 'seller@afrimart.com',
        products: [],
        paymentMethods: [],
        status: 'active',
        backendId: 'shop-uuid-999'
      }

      let capturedBody: any = null
      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/shops/me')) {
          return Promise.resolve(mockBackendShop)
        }
        if (url.includes('/api/products') && opts?.method === 'POST') {
          capturedBody = opts?.body
          return Promise.resolve({
            product: mockBackendProduct,
            images: [
              {
                ID: 'img-uuid-1',
                ProductID: 'prod-uuid-555',
                ObjectKey: 'products/prod-uuid-555/image.png',
                DisplayOrder: 0
              }
            ]
          })
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const testFile = new File(['fake image bytes'], 'kente.png', { type: 'image/png' })

      const newProduct = await productRepo.createProduct('seller@afrimart.com', {
        name: 'Handwoven Kente Scarf',
        description: 'Authentic handwoven silk and cotton Kente scarf',
        category: 'Men',
        categoryId: '11111111-1111-4111-8111-111111111111',
        subcategoryId: '22222222-2222-4222-8222-222222222222',
        price: 250.50,
        stock: 15,
        files: [testFile]
      })

      expect(newProduct.id).toBe('prod-uuid-555')
      expect(newProduct.name).toBe('Handwoven Kente Scarf')
      expect(newProduct.price).toBe(250.50)
      expect(capturedBody).toBeInstanceOf(FormData)
      expect(capturedBody.get('name')).toBe('Handwoven Kente Scarf')
      expect(capturedBody.get('category_id')).toBe('11111111-1111-4111-8111-111111111111')
      expect(capturedBody.get('price')).toBe('250.5')
      expect(capturedBody.get('status')).toBe('active')
    })

    it('createProduct: should throw when category setup is not ready', async () => {
      const { shop } = useMockDataStore()
      shop.value = {
        id: 'shop-uuid-999',
        name: 'AfriMart Artisan Studio',
        slug: 'afrimart-artisan-studio',
        description: 'Studio desc',
        ownerEmail: 'seller@afrimart.com',
        products: [],
        paymentMethods: [],
        status: 'active',
        backendId: 'shop-uuid-999'
      }

      const testFile = new File(['fake image'], 'kente.png', { type: 'image/png' })

      await expect(
        productRepo.createProduct('seller@afrimart.com', {
          name: 'Handwoven Kente Scarf',
          description: 'Description',
          category: 'Men',
          price: 100,
          stock: 5,
          files: [testFile]
        })
      ).rejects.toThrow('Please select a valid category and subcategory.')
    })

    it('createProduct: should throw when no image files are provided', async () => {
      const { shop } = useMockDataStore()
      shop.value = {
        id: 'shop-uuid-999',
        name: 'AfriMart Artisan Studio',
        slug: 'afrimart-artisan-studio',
        description: 'Studio desc',
        ownerEmail: 'seller@afrimart.com',
        products: [],
        paymentMethods: [],
        status: 'active',
        backendId: 'shop-uuid-999'
      }

      await expect(
        productRepo.createProduct('seller@afrimart.com', {
          name: 'Handwoven Kente Scarf',
          description: 'Description',
          category: 'Men',
          categoryId: '11111111-1111-4111-8111-111111111111',
          subcategoryId: '22222222-2222-4222-8222-222222222222',
          price: 100,
          stock: 5,
          files: []
        })
      ).rejects.toThrow('At least one product image is required (max 10 images, max 5 MB each).')
    })

    it('updateProduct: should PUT /api/products/{id} with updated product attributes', async () => {
      const mockUpdatedBackendProduct = {
        ...mockBackendProduct,
        Price: '299.99',
        Stock: 20
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/products/prod-uuid-555')) {
          if (opts?.method === 'PUT') {
            return Promise.resolve(mockUpdatedBackendProduct)
          }
          return Promise.resolve(mockBackendProduct)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const updated = await productRepo.updateProduct('prod-uuid-555', {
        price: 299.99,
        stock: 20
      })

      expect(updated).not.toBeNull()
      expect(updated?.price).toBe(299.99)
      expect(updated?.stock).toBe(20)
    })

    it('deleteProduct: should call DELETE /api/products/{id}', async () => {
      let deleteCalled = false

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/products/prod-uuid-555') && opts?.method === 'DELETE') {
          deleteCalled = true
          return Promise.resolve({ status: 'deleted' })
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const success = await productRepo.deleteProduct('prod-uuid-555')

      expect(success).toBe(true)
      expect(deleteCalled).toBe(true)
    })
  })

  describe('ApiCartRepository', () => {
    const cartRepo = new ApiCartRepository()

    const mockCartResponse = {
      id: 'cart-uuid-111',
      user_id: 'user-uuid-101',
      items: [
        {
          id: 'cart-item-uuid-1',
          cart_id: 'cart-uuid-111',
          product_id: 'prod-uuid-555',
          quantity: 2,
          name: 'Handwoven Kente Scarf',
          price: '250.50',
          stock: 15
        }
      ],
      subtotal: 501.00
    }

    it('getCart: should fetch cart with Authorization Bearer header when logged in', async () => {
      apiHelpers.setAccessToken('test-access-token')
      let capturedHeader = ''

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/cart')) {
          capturedHeader = opts?.headers?.Authorization || ''
          return Promise.resolve(mockCartResponse)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const cart = await cartRepo.getCart()

      expect(cart).not.toBeNull()
      expect(cart?.id).toBe('cart-uuid-111')
      expect(cart?.items.length).toBe(1)
      expect(capturedHeader).toBe('Bearer test-access-token')
    })

    it('getCart: should return null when user is not logged in', async () => {
      apiHelpers.clearTokens()

      const cart = await cartRepo.getCart()
      expect(cart).toBeNull()
    })

    it('addItem: should POST /api/cart/items with product_id and quantity', async () => {
      apiHelpers.setAccessToken('test-access-token')
      let postedBody: any = null

      const mockAddedItem = {
        id: 'cart-item-uuid-2',
        cart_id: 'cart-uuid-111',
        product_id: 'prod-uuid-777',
        quantity: 3
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/cart/items') && opts?.method === 'POST') {
          postedBody = opts.body
          return Promise.resolve(mockAddedItem)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const added = await cartRepo.addItem('prod-uuid-777', 3)

      expect(added.id).toBe('cart-item-uuid-2')
      expect(postedBody).toEqual({ product_id: 'prod-uuid-777', quantity: 3 })
    })

    it('updateItemQuantity: should PATCH /api/cart/items/{id} with new quantity', async () => {
      apiHelpers.setAccessToken('test-access-token')
      let patchBody: any = null

      const mockUpdatedItem = {
        id: 'cart-item-uuid-1',
        cart_id: 'cart-uuid-111',
        product_id: 'prod-uuid-555',
        quantity: 5
      }

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/cart/items/cart-item-uuid-1') && opts?.method === 'PATCH') {
          patchBody = opts.body
          return Promise.resolve(mockUpdatedItem)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const updated = await cartRepo.updateItemQuantity('cart-item-uuid-1', 5)

      expect(updated.quantity).toBe(5)
      expect(patchBody).toEqual({ quantity: 5 })
    })

    it('removeItem: should DELETE /api/cart/items/{id}', async () => {
      apiHelpers.setAccessToken('test-access-token')
      let deleteCalled = false

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/cart/items/cart-item-uuid-1') && opts?.method === 'DELETE') {
          deleteCalled = true
          return Promise.resolve()
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const success = await cartRepo.removeItem('cart-item-uuid-1')

      expect(success).toBe(true)
      expect(deleteCalled).toBe(true)
    })

    it('clearCart: should DELETE /api/cart', async () => {
      apiHelpers.setAccessToken('test-access-token')
      let clearCalled = false

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/cart') && opts?.method === 'DELETE') {
          clearCalled = true
          return Promise.resolve()
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const success = await cartRepo.clearCart()

      expect(success).toBe(true)
      expect(clearCalled).toBe(true)
    })
  })

  describe('ApiOrderRepository', () => {
    const orderRepo = new ApiOrderRepository()

    const mockOrderResponse = {
      order: {
        id: 'order-uuid-9901',
        recipient_name: 'Aynalem Sitotaw',
        delivery_address: 'Bole Road, Addis Ababa',
        delivery_city: 'Addis Ababa',
        phone: '+251911223344',
        status: 'pending',
        subtotal: 750,
        created_at: '2026-09-20T10:00:00Z'
      },
      items: [{ product_id: 'prod-uuid-555', quantity: 2 }]
    }

    it('createOrder: should POST /api/orders/checkout with recipient payload and Bearer token', async () => {
      apiHelpers.setAccessToken('test-access-token')
      let capturedBody: any = null
      let capturedHeader = ''

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string, opts: any) => {
        if (url.includes('/api/orders/checkout') && opts?.method === 'POST') {
          capturedBody = opts.body
          capturedHeader = opts?.headers?.Authorization || ''
          return Promise.resolve(mockOrderResponse)
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const order = await orderRepo.createOrder(
        [{ productId: 555, quantity: 2 }],
        500,
        {
          buyerName: 'Aynalem Sitotaw',
          deliveryAddress: 'Bole Road, Addis Ababa',
          deliveryCity: 'Addis Ababa',
          phone: '+251911223344'
        }
      )

      expect(order.buyerName).toBe('Aynalem Sitotaw')
      expect(capturedHeader).toBe('Bearer test-access-token')
      expect(capturedBody.recipient_name).toBe('Aynalem Sitotaw')
      expect(capturedBody.delivery_city).toBe('Addis Ababa')
    })

    it('getOrders: should fetch list from GET /api/orders when authenticated', async () => {
      apiHelpers.setAccessToken('test-access-token')

      vi.stubGlobal('$fetch', vi.fn().mockImplementation((url: string) => {
        if (url.includes('/api/orders')) {
          return Promise.resolve({ page: 1, limit: 20, orders: [mockOrderResponse.order] })
        }
        return Promise.reject(new Error(`Unexpected fetch URL: ${url}`))
      }))
      vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'http://localhost:8080' } }))

      const ordersList = await orderRepo.getOrders()

      expect(ordersList.length).toBeGreaterThan(0)
      expect(ordersList[0].buyerName).toBe('Aynalem Sitotaw')
    })
  })
})
