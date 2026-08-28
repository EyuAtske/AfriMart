import { MockProductRepository } from '~/repositories/mock/MockProductRepository'
import { MockShopRepository } from '~/repositories/mock/MockShopRepository'
import { MockOrderRepository } from '~/repositories/mock/MockOrderRepository'
import { MockAuthRepository } from '~/repositories/mock/MockAuthRepository'
import { ApiAuthRepository } from '~/repositories/api/ApiAuthRepository'
import type { IProductRepository } from '~/repositories/interfaces/IProductRepository'
import type { IShopRepository } from '~/repositories/interfaces/IShopRepository'
import type { IOrderRepository } from '~/repositories/interfaces/IOrderRepository'
import type { IAuthRepository } from '~/repositories/interfaces/IAuthRepository'

const productRepository: IProductRepository = new MockProductRepository()
const shopRepository: IShopRepository = new MockShopRepository()
const orderRepository: IOrderRepository = new MockOrderRepository()

const mockAuthRepository: IAuthRepository = new MockAuthRepository()
const apiAuthRepository: IAuthRepository = new ApiAuthRepository()

export const useRepositories = () => {
  const config = useRuntimeConfig()
  const authMode = (config?.public?.authMode as string) || 'mock'

  const authRepo = authMode === 'api' ? apiAuthRepository : mockAuthRepository

  return {
    productRepo: productRepository,
    shopRepo: shopRepository,
    orderRepo: orderRepository,
    authRepo
  }
}
