import type { ProductMedia } from '~/types/product'

export interface MediaUploadResult {
  uploaded: ProductMedia[]
  errors: string[]
}

export interface IMediaRepository {
  uploadMedia(files: File[]): Promise<MediaUploadResult>
  deleteMedia(mediaId: string): Promise<boolean>
}
