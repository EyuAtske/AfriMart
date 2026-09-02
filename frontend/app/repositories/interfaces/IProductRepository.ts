import type { Product, ProductFilterParams, CreateProductDTO, UpdateProductDTO } from '~/types/product'
import type { PaginatedResponse } from '~/types/api'

export interface IProductRepository {
  getProducts(params?: ProductFilterParams): Promise<PaginatedResponse<Product>>
  getProductById(id: number): Promise<Product | null>
  createProduct(shopName: string, dto: CreateProductDTO): Promise<Product>
  updateProduct(id: number, dto: UpdateProductDTO): Promise<Product | null>
  deleteProduct(id: number): Promise<boolean>
  toggleProductStatus(id: number): Promise<Product | null>
}
