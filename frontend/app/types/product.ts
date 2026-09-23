export type ProductGender = 'Men' | 'Women' | 'Kids'
export type MainCategory = 'Clothing' | 'Accessories'
export type ProductCategory = MainCategory | 'Men' | 'Women' | 'Kids' | 'Shoes'

export const MAIN_CATEGORIES: MainCategory[] = ['Clothing', 'Accessories']
export const GENDERS: ProductGender[] = ['Men', 'Women', 'Kids']

export const HIERARCHICAL_ITEMS: Record<MainCategory, Record<ProductGender, ProductSubCategory[]>> = {
  Clothing: {
    Men: ['T-Shirts', 'Shirts', 'Trousers', 'Jeans', 'Jackets', 'Hoodies', 'Sweaters', 'Other'],
    Women: ['Dresses', 'Tops', 'T-Shirts', 'Shirts', 'Trousers', 'Jeans', 'Skirts', 'Jackets', 'Sweaters', 'Hoodies', 'Other'],
    Kids: ['T-Shirts', 'Dresses', 'Trousers', 'Jeans', 'Jackets', 'Hoodies', 'Sweaters', 'Other']
  },
  Accessories: {
    Men: ['Bags', 'Watches', 'Belts', 'Hats', 'Sneakers', 'Formal Shoes', 'Boots', 'Sandals', 'Other'],
    Women: ['Bags', 'Jewelry', 'Watches', 'Belts', 'Hats', 'Scarves', 'Sneakers', 'Heels', 'Boots', 'Sandals', 'Flats', 'Other'],
    Kids: ['Bags', 'Hats', 'Belts', 'Sneakers', 'Sandals', 'Boots', 'Other']
  }
}

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
  Clothing: ['T-Shirts', 'Shirts', 'Dresses', 'Tops', 'Trousers', 'Jeans', 'Skirts', 'Jackets', 'Sweaters', 'Hoodies', 'Other'],
  Accessories: ['Bags', 'Watches', 'Belts', 'Hats', 'Jewelry', 'Scarves', 'Sneakers', 'Other'],
  Men: ['T-Shirts', 'Shirts', 'Trousers', 'Jeans', 'Jackets', 'Hoodies', 'Sweaters', 'Sneakers', 'Formal Shoes', 'Boots', 'Sandals', 'Bags', 'Watches', 'Belts', 'Hats', 'Other'],
  Women: ['Dresses', 'Tops', 'T-Shirts', 'Shirts', 'Trousers', 'Jeans', 'Skirts', 'Jackets', 'Sweaters', 'Hoodies', 'Sneakers', 'Heels', 'Boots', 'Sandals', 'Flats', 'Bags', 'Jewelry', 'Belts', 'Hats', 'Scarves', 'Other'],
  Kids: ['T-Shirts', 'Dresses', 'Trousers', 'Jeans', 'Jackets', 'Hoodies', 'Sweaters', 'Sneakers', 'Sandals', 'Boots', 'Bags', 'Hats', 'Belts', 'Other'],
  Shoes: ['Sneakers', 'Formal Shoes', 'Heels', 'Boots', 'Sandals', 'Flats', 'Other']
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
  file?: File
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
  gender?: ProductGender
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
  gender?: string
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
  categoryId?: string
  gender?: ProductGender
  subCategory?: ProductSubCategory
  subcategoryId?: string
  brand?: string
  color?: string
  size?: string
  price: number
  stock: number
  image?: string
  status?: ProductStatus
  media?: ProductMedia[]
  files?: File[]
}

export interface UpdateProductDTO {
  name?: string
  description?: string
  category?: ProductCategory
  categoryId?: string
  gender?: ProductGender
  subCategory?: ProductSubCategory
  subcategoryId?: string
  brand?: string
  color?: string
  size?: string
  price?: number
  stock?: number
  image?: string
  status?: ProductStatus
  media?: ProductMedia[]
  files?: File[]
}

export interface BackendProductImage {
  ID: string
  ProductID: string
  ObjectKey: string
  DisplayOrder: number
  CreatedAt?: string
}

export interface BackendProductCreateResponse {
  product: {
    ID: string
    ShopID: string
    CategoryID: string
    SubcategoryID: string
    Name: string
    Description: { String: string; Valid: boolean } | string | null
    Brand?: { String: string; Valid: boolean } | string | null
    Color?: { String: string; Valid: boolean } | string | null
    Size?: { String: string; Valid: boolean } | string | null
    Price: string
    Stock: number
    Status: string
    CreatedAt?: string
    UpdatedAt?: string
  }
  images: BackendProductImage[]
}

