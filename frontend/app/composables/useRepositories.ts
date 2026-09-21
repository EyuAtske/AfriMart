import { MockProductRepository } from '~/repositories/mock/MockProductRepository'
import { MockShopRepository } from '~/repositories/mock/MockShopRepository'
import { MockOrderRepository } from '~/repositories/mock/MockOrderRepository'
import { MockAuthRepository } from '~/repositories/mock/MockAuthRepository'
import { MockMediaRepository } from '~/repositories/mock/MockMediaRepository'
import { MockCartRepository } from '~/repositories/mock/MockCartRepository'
import { ApiAuthRepository } from '~/repositories/api/ApiAuthRepository'
import { ApiShopRepository } from '~/repositories/api/ApiShopRepository'
import { ApiProductRepository } from '~/repositories/api/ApiProductRepository'
import { ApiCartRepository } from '~/repositories/api/ApiCartRepository'
import { ApiOrderRepository } from '~/repositories/api/ApiOrderRepository'
import type { IProductRepository } from '~/repositories/interfaces/IProductRepository'
import type { IShopRepository } from '~/repositories/interfaces/IShopRepository'
import type { IOrderRepository } from '~/repositories/interfaces/IOrderRepository'
import type { IAuthRepository } from '~/repositories/interfaces/IAuthRepository'
import type { IMediaRepository } from '~/repositories/interfaces/IMediaRepository'
import type { ICartRepository } from '~/repositories/interfaces/ICartRepository'

const mockProductRepository: IProductRepository = new MockProductRepository()
const apiProductRepository: IProductRepository = new ApiProductRepository()
const mockOrderRepository: IOrderRepository = new MockOrderRepository()
const apiOrderRepository: IOrderRepository = new ApiOrderRepository()
const mediaRepository: IMediaRepository = new MockMediaRepository()

const mockAuthRepository: IAuthRepository = new MockAuthRepository()
const apiAuthRepository: IAuthRepository = new ApiAuthRepository()

const mockShopRepository: IShopRepository = new MockShopRepository()
const apiShopRepository: IShopRepository = new ApiShopRepository()

const mockCartRepository: ICartRepository = new MockCartRepository()
const apiCartRepository: ICartRepository = new ApiCartRepository()

export const useRepositories = () => {
  const config = useRuntimeConfig()
  const authMode = (config?.public?.authMode as string) || 'mock'

  const authRepo = authMode === 'api' ? apiAuthRepository : mockAuthRepository
  const shopRepo = authMode === 'api' ? apiShopRepository : mockShopRepository
  const productRepo = authMode === 'api' ? apiProductRepository : mockProductRepository
  const cartRepo = authMode === 'api' ? apiCartRepository : mockCartRepository
  const orderRepo = authMode === 'api' ? apiOrderRepository : mockOrderRepository

  return {
    productRepo,
    shopRepo,
    cartRepo,
    orderRepo,
    mediaRepo: mediaRepository,
    authRepo
  }
}
