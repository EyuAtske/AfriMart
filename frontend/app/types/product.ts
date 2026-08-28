export type ProductCategory = 'Men' | 'Women' | 'Kids' | 'Shoes' | 'Accessories'
export type ProductStatus = 'Active' | 'Draft'

export interface ProductReview {
  id: number
  productId: number
  orderId?: number
  author: string
  rating: number
  comment: string
  date: string
  status: 'pending' | 'approved'
}

export interface Product {
  id: number
  shop: string
  name: string
  description: string
  category: ProductCategory
  price: number
  stock: number
  rating: string
  image: string
  status: ProductStatus
  reviewsCount?: number
  reviews?: ProductReview[]
}

export interface ProductFilterParams {
  search?: string
  category?: string
  shop?: string
  status?: ProductStatus
  page?: number
  pageSize?: number
}

export interface CreateProductDTO {
  name: string
  description: string
  category: ProductCategory
  price: number
  stock: number
  image: string
  status?: ProductStatus
}

export interface UpdateProductDTO {
  name?: string
  description?: string
  category?: ProductCategory
  price?: number
  stock?: number
  image?: string
  status?: ProductStatus
}
