export type ProductCategory = 'Men' | 'Women' | 'Kids' | 'Shoes' | 'Accessories'
export type ProductSubCategory =
  | 'T-Shirts'
  | 'Shirts'
  | 'Dresses'
  | 'Tops'
  | 'Trousers'
  | 'Jeans'
  | 'Skirts'
  | 'Jackets'
  | 'Hoodies'
  | 'Sweaters'
  | 'Sneakers'
  | 'Formal Shoes'
  | 'Heels'
  | 'Boots'
  | 'Sandals'
  | 'Flats'
  | 'Bags'
  | 'Watches'
  | 'Belts'
  | 'Hats'
  | 'Jewelry'
  | 'Scarves'
  | 'Other'

export const CATEGORY_SUBCATEGORIES: Record<ProductCategory, ProductSubCategory[]> = {
  Men: ['T-Shirts', 'Shirts', 'Trousers', 'Jeans', 'Jackets', 'Hoodies', 'Sweaters', 'Sneakers', 'Formal Shoes', 'Boots', 'Sandals', 'Bags', 'Watches', 'Belts', 'Hats', 'Other'],
  Women: ['Dresses', 'Tops', 'T-Shirts', 'Shirts', 'Trousers', 'Jeans', 'Skirts', 'Jackets', 'Sweaters', 'Hoodies', 'Sneakers', 'Heels', 'Boots', 'Sandals', 'Flats', 'Bags', 'Jewelry', 'Belts', 'Hats', 'Scarves', 'Other'],
  Kids: ['T-Shirts', 'Dresses', 'Trousers', 'Jeans', 'Jackets', 'Hoodies', 'Sweaters', 'Sneakers', 'Sandals', 'Boots', 'Bags', 'Hats', 'Belts', 'Other'],
  Shoes: ['Sneakers', 'Formal Shoes', 'Heels', 'Boots', 'Sandals', 'Flats', 'Other'],
  Accessories: ['Bags', 'Watches', 'Belts', 'Hats', 'Jewelry', 'Scarves', 'Other']
}

export type ProductStatus = 'Active' | 'Draft'
export type MediaType = 'image' | 'video'

export interface ProductMedia {
  id: string
  type: MediaType
  url: string
  alt: string
  position: number
  isPrimary: boolean
  fileName: string
  fileSize: number
}

export interface ProductReview {
  id: number | string
  productId: number | string
  orderId?: number | string
  author: string
  rating: number
  comment: string
  date: string
  status: 'pending' | 'approved'
}

export interface Product {
  id: number | string
  backendId?: string
  shop: string
  name: string
  description: string
  category: ProductCategory
  subCategory?: ProductSubCategory
  price: number
  stock: number
  rating: string
  image: string
  status: ProductStatus
  media?: ProductMedia[]
  reviewsCount?: number
  reviews?: ProductReview[]
}

export interface ProductFilterParams {
  search?: string
  category?: string
  subCategory?: string
  shop?: string
  status?: ProductStatus
  page?: number
  pageSize?: number
}

export interface CreateProductDTO {
  name: string
  description: string
  category: ProductCategory
  subCategory?: ProductSubCategory
  price: number
  stock: number
  image: string
  status?: ProductStatus
  media?: ProductMedia[]
}

export interface UpdateProductDTO {
  name?: string
  description?: string
  category?: ProductCategory
  subCategory?: ProductSubCategory
  price?: number
  stock?: number
  image?: string
  status?: ProductStatus
  media?: ProductMedia[]
}
