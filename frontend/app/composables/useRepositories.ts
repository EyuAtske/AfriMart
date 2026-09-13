import { MockProductRepository } from '~/repositories/mock/MockProductRepository'
import { MockShopRepository } from '~/repositories/mock/MockShopRepository'
import { MockOrderRepository } from '~/repositories/mock/MockOrderRepository'
import { MockAuthRepository } from '~/repositories/mock/MockAuthRepository'
import { MockMediaRepository } from '~/repositories/mock/MockMediaRepository'
import { ApiAuthRepository } from '~/repositories/api/ApiAuthRepository'
import { ApiShopRepository } from '~/repositories/api/ApiShopRepository'
import type { IProductRepository } from '~/repositories/interfaces/IProductRepository'
import type { IShopRepository } from '~/repositories/interfaces/IShopRepository'
import type { IOrderRepository } from '~/repositories/interfaces/IOrderRepository'
import type { IAuthRepository } from '~/repositories/interfaces/IAuthRepository'
import type { IMediaRepository } from '~/repositories/interfaces/IMediaRepository'

const productRepository: IProductRepository = new MockProductRepository()
const orderRepository: IOrderRepository = new MockOrderRepository()
const mediaRepository: IMediaRepository = new MockMediaRepository()

const mockAuthRepository: IAuthRepository = new MockAuthRepository()
const apiAuthRepository: IAuthRepository = new ApiAuthRepository()

const mockShopRepository: IShopRepository = new MockShopRepository()
const apiShopRepository: IShopRepository = new ApiShopRepository()

export const useRepositories = () => {
  const config = useRuntimeConfig()
  const authMode = (config?.public?.authMode as string) || 'mock'

  const authRepo = authMode === 'api' ? apiAuthRepository : mockAuthRepository
  const shopRepo = authMode === 'api' ? apiShopRepository : mockShopRepository

  return {
    productRepo: productRepository,
    shopRepo,
    orderRepo: orderRepository,
    mediaRepo: mediaRepository,
    authRepo
  }
}
