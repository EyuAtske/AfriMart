import type { Product, ProductFilterParams, CreateProductDTO, UpdateProductDTO } from '~/types/product'
import type { PaginatedResponse } from '~/types/api'

export interface IProductRepository {
  getProducts(params?: ProductFilterParams): Promise<PaginatedResponse<Product>>
  getProductById(id: number | string): Promise<Product | null>
  createProduct(shopName: string, dto: CreateProductDTO): Promise<Product>
  updateProduct(id: number | string, dto: UpdateProductDTO): Promise<Product | null>
  deleteProduct(id: number | string): Promise<boolean>
  toggleProductStatus(id: number | string): Promise<Product | null>
}
