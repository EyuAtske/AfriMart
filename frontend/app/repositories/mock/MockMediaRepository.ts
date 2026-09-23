import type { IMediaRepository, MediaUploadResult } from '../interfaces/IMediaRepository'
import type { ProductMedia, MediaType } from '~/types/product'

const ACCEPTED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp']
const ACCEPTED_VIDEO_TYPES = ['video/mp4', 'video/webm']
const ALL_ACCEPTED_TYPES = [...ACCEPTED_IMAGE_TYPES, ...ACCEPTED_VIDEO_TYPES]

const MAX_FILES = 10
const MAX_IMAGE_SIZE = 5 * 1024 * 1024    // 5 MB (matches backend limit)
const MAX_VIDEO_SIZE = 50 * 1024 * 1024   // 50 MB

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function generateId(): string {
  return `media-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

function getMediaType(mimeType: string): MediaType {
  return ACCEPTED_IMAGE_TYPES.includes(mimeType) ? 'image' : 'video'
}

export class MockMediaRepository implements IMediaRepository {
  async uploadMedia(files: File[]): Promise<MediaUploadResult> {
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 300))

    const uploaded: ProductMedia[] = []
    const errors: string[] = []

    for (const file of files) {
      // Validate MIME type
      if (!ALL_ACCEPTED_TYPES.includes(file.type)) {
        errors.push(
          `"${file.name}" — unsupported format. Accepted: JPEG, PNG, WebP, MP4, WebM.`
        )
        continue
      }

      const mediaType = getMediaType(file.type)

      // Validate file size
      const maxSize = mediaType === 'image' ? MAX_IMAGE_SIZE : MAX_VIDEO_SIZE
      if (file.size > maxSize) {
        errors.push(
          `"${file.name}" — too large (${formatFileSize(file.size)}). Max ${mediaType === 'image' ? '5 MB' : '50 MB'}.`
        )
        continue
      }

      // Validate total count
      if (uploaded.length >= MAX_FILES) {
        errors.push(
          `"${file.name}" — limit reached. Maximum ${MAX_FILES} files per product.`
        )
        continue
      }

      // Create object URL (client-side only)
      let url = ''
      if (import.meta.client) {
        url = URL.createObjectURL(file)
      }

      uploaded.push({
        id: generateId(),
        type: mediaType,
        url,
        alt: '',
        position: uploaded.length,
        isPrimary: uploaded.length === 0 && mediaType === 'image',
        fileName: file.name,
        fileSize: file.size,
        file
      })
    }

    return { uploaded, errors }
  }

  async deleteMedia(mediaId: string): Promise<boolean> {
    // In a real implementation this would call the API
    // Object URLs are revoked by the component when items are removed
    return true
  }
}

export { MAX_FILES, MAX_IMAGE_SIZE, MAX_VIDEO_SIZE, ACCEPTED_IMAGE_TYPES, ACCEPTED_VIDEO_TYPES, ALL_ACCEPTED_TYPES, formatFileSize }
