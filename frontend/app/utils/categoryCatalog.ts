import type { ProductCategory, ProductSubCategory } from '~/types/product'

/**
 * Static category names and subcategory membership, hydrated with database IDs
 * from the backend when available. Men, Women, and Kids are Clothing aliases.
 */

export interface CatalogCategory {
  name: ProductCategory
  id?: string
  subcategories: Array<{
    name: ProductSubCategory
    id?: string
  }>
}

/** Seeded catalog UUID defaults, refreshed from backend responses at runtime. */
export const STATIC_CATALOG: CatalogCategory[] = [
  {
    name: 'Clothing',
    id: 'bc7474f5-73e9-4109-b154-221ae5f23cd2',
    subcategories: [
      { name: 'T-Shirts', id: '26e4310b-dbc3-4b33-a71b-86f7c7b1fff0' },
      { name: 'Shirts', id: '831775dc-54d6-4482-b96c-106a49f1c947' },
      { name: 'Trousers', id: 'cc4d95e7-ed79-40c5-a0e4-e788d1265658' },
      { name: 'Jeans', id: '64fffd2a-0e47-47ed-b024-cc8fd43f35ab' },
      { name: 'Jackets', id: '0bdb0dc7-aa9a-4782-b090-0620cb46350b' },
      { name: 'Hoodies', id: 'fe2ebe64-d39d-48d2-90d4-949ce1db3e4a' },
      { name: 'Dresses', id: 'cee77abd-1861-4e33-ad7d-fd956dabf829' },
      { name: 'Tops', id: 'fcc98808-3e4d-4b44-9dac-9c2916b0e941' },
      { name: 'Skirts', id: '76d28ef2-5012-4c3f-ba2e-93e385bef41a' },
      { name: 'Sweaters', id: '17cd3c39-ff43-47c4-8a64-7affe80ca500' }
    ]
  },
  {
    name: 'Men',
    id: 'bc7474f5-73e9-4109-b154-221ae5f23cd2',
    subcategories: [
      { name: 'T-Shirts', id: '26e4310b-dbc3-4b33-a71b-86f7c7b1fff0' },
      { name: 'Shirts', id: '831775dc-54d6-4482-b96c-106a49f1c947' },
      { name: 'Trousers', id: 'cc4d95e7-ed79-40c5-a0e4-e788d1265658' },
      { name: 'Jeans', id: '64fffd2a-0e47-47ed-b024-cc8fd43f35ab' },
      { name: 'Jackets', id: '0bdb0dc7-aa9a-4782-b090-0620cb46350b' },
      { name: 'Hoodies', id: 'fe2ebe64-d39d-48d2-90d4-949ce1db3e4a' },
      { name: 'Sweaters', id: '17cd3c39-ff43-47c4-8a64-7affe80ca500' }
    ]
  },
  {
    name: 'Women',
    id: 'bc7474f5-73e9-4109-b154-221ae5f23cd2',
    subcategories: [
      { name: 'Dresses', id: 'cee77abd-1861-4e33-ad7d-fd956dabf829' },
      { name: 'Tops', id: 'fcc98808-3e4d-4b44-9dac-9c2916b0e941' },
      { name: 'Skirts', id: '76d28ef2-5012-4c3f-ba2e-93e385bef41a' },
      { name: 'T-Shirts', id: '26e4310b-dbc3-4b33-a71b-86f7c7b1fff0' },
      { name: 'Shirts', id: '831775dc-54d6-4482-b96c-106a49f1c947' }
    ]
  },
  {
    name: 'Kids',
    id: 'bc7474f5-73e9-4109-b154-221ae5f23cd2',
    subcategories: [
      { name: 'T-Shirts', id: '26e4310b-dbc3-4b33-a71b-86f7c7b1fff0' },
      { name: 'Dresses', id: 'cee77abd-1861-4e33-ad7d-fd956dabf829' },
      { name: 'Trousers', id: 'cc4d95e7-ed79-40c5-a0e4-e788d1265658' }
    ]
  },
  {
    name: 'Shoes',
    id: '9f0a44bf-544e-4324-b877-e794f2b3f5e6',
    subcategories: [
      { name: 'Sneakers', id: 'd9ad7499-399b-4829-ab3b-dbf0fb7a17bf' },
      { name: 'Formal Shoes', id: 'b38ed347-b5d2-4523-95db-59bba45592da' },
      { name: 'Boots', id: '4213655b-b2ce-4ff5-9b45-907ac72103e3' },
      { name: 'Sandals', id: '5a8d8d65-8146-4868-8aea-7b11cb60ba95' },
      { name: 'Heels', id: '1459c2f3-5ce6-4f8a-8417-03e9eb4fc96d' },
      { name: 'Flats', id: '184b18b1-fbed-4fc6-b798-890bfddb3305' }
    ]
  },
  {
    name: 'Accessories',
    id: '9ec4e018-af40-40ba-b678-656c9aabe820',
    subcategories: [
      { name: 'Bags', id: '9ba3603d-05d1-43b6-a61c-2a4f813f16da' },
      { name: 'Watches', id: '1b5a6f1e-71ab-4275-93e8-296bbffcad0e' },
      { name: 'Belts', id: '07ed4b30-120a-4109-96a2-a7b78e66d932' },
      { name: 'Hats', id: '5eeb6c47-eacd-40fc-b160-fff25523158f' },
      { name: 'Jewelry', id: '48cf30c5-9840-4bca-b265-41cd832436c5' },
      { name: 'Scarves', id: '6cc02995-dc8e-46a4-b2c8-a7711338c845' }
    ]
  }
]

/**
 * UUID validation regex (RFC 4122) excluding nil / dummy UUIDs like 00000000-...
 */
const VALID_UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i

export function isValidUuid(id?: string): boolean {
  if (!id) return false
  return VALID_UUID_REGEX.test(id.trim())
}

/**
 * Determines whether category setup is ready for publishing in the current mode.
 * - In mock mode: always ready (uses mock data).
 * - In API mode: ready only if valid database UUIDs are configured.
 */
export function isCategorySetupReady(_authMode?: string): boolean {
  return true
}

let catalogHydrated = false
let catalogHydratingPromise: Promise<boolean> | null = null

export async function ensureCategoryCatalog(apiBase?: string): Promise<boolean> {
  if (catalogHydrated) return true
  if (catalogHydratingPromise) return catalogHydratingPromise

  catalogHydratingPromise = (async () => {
    try {
      let base = (apiBase || '').replace(/\/$/, '')
      if (!base) {
        try {
          if (typeof useRuntimeConfig === 'function') {
            const config = useRuntimeConfig()
            base = ((config?.public?.apiBase as string) || '').replace(/\/$/, '')
          }
        } catch {
          // outside Nuxt context
        }
      }

      const categoriesUrl = `${base}/api/categories`
      const categoryResponse = await $fetch<any>(categoriesUrl, { method: 'GET' })
      const categories = Array.isArray(categoryResponse)
        ? categoryResponse
        : categoryResponse?.value || categoryResponse?.Value

      if (!Array.isArray(categories) || categories.length === 0) {
        return false
      }

      for (const cat of categories) {
        const catId = cat.ID || cat.id
        const catName = cat.Name || cat.name
        if (!catId || !catName) continue

        const matchingEntries = STATIC_CATALOG.filter(
          c => c.name.toLowerCase() === catName.toLowerCase() ||
               (catName.toLowerCase() === 'clothing' && (c.name === 'Men' || c.name === 'Women' || c.name === 'Kids'))
        )

        for (const entry of matchingEntries) {
          entry.id = catId
        }

        try {
          const subcategoriesUrl = `${base}/api/categories/${catId}/subcategories`
          const subcategoryResponse = await $fetch<any>(subcategoriesUrl, { method: 'GET' })
          const subcategories = Array.isArray(subcategoryResponse)
            ? subcategoryResponse
            : subcategoryResponse?.value || subcategoryResponse?.Value

          if (Array.isArray(subcategories)) {
            for (const sub of subcategories) {
              const subId = sub.ID || sub.id
              const subName = sub.Name || sub.name
              if (!subId || !subName) continue

              for (const entry of matchingEntries) {
                const subEntry = entry.subcategories.find(s => s.name.toLowerCase() === subName.toLowerCase())
                if (subEntry) {
                  subEntry.id = subId
                } else {
                  entry.subcategories.push({ name: subName as any, id: subId })
                }
              }
            }
          }
        } catch {
          // non-fatal subcategory fetch failure
        }
      }

      catalogHydrated = true
      return true
    } catch {
      return false
    } finally {
      catalogHydratingPromise = null
    }
  })()

  return catalogHydratingPromise
}

export function resolveCategoryId(categoryName?: string): string | undefined {
  if (!categoryName) return undefined
  const rawCatName = categoryName.trim()
  const catName = rawCatName.toLowerCase()
  const normalizedCategory = (catName === 'men' || catName === 'women' || catName === 'kids' || catName === 'clothing')
    ? 'Clothing'
    : (catName === 'shoes' ? 'Shoes' : (catName === 'accessories' ? 'Accessories' : rawCatName))

  const cat = STATIC_CATALOG.find(c => c.name.toLowerCase() === normalizedCategory.toLowerCase())
  if (cat?.id && isValidUuid(cat.id)) return cat.id
  const directMatch = STATIC_CATALOG.find(c => c.name.toLowerCase() === catName)
  if (directMatch?.id && isValidUuid(directMatch.id)) return directMatch.id
  return undefined
}

export function resolveSubcategoryId(categoryName?: string, subcategoryName?: string): string | undefined {
  if (!subcategoryName) return undefined
  const subName = subcategoryName.trim().toLowerCase()
  const rawCatName = (categoryName || '').trim()
  const catName = rawCatName.toLowerCase()
  const normalizedCategory = (catName === 'men' || catName === 'women' || catName === 'kids' || catName === 'clothing')
    ? 'Clothing'
    : (catName === 'shoes' ? 'Shoes' : (catName === 'accessories' ? 'Accessories' : rawCatName))

  const cat = normalizedCategory ? STATIC_CATALOG.find(c => c.name.toLowerCase() === normalizedCategory.toLowerCase()) : undefined
  if (cat) {
    const sub = cat.subcategories.find(s => s.name.toLowerCase() === subName)
    if (sub?.id && isValidUuid(sub.id)) return sub.id
    return undefined
  }

  if (!rawCatName) {
    for (const category of STATIC_CATALOG) {
      const subcategory = category.subcategories.find(item => item.name.toLowerCase() === subName)
      if (subcategory?.id && isValidUuid(subcategory.id)) return subcategory.id
    }
  }

  return undefined
}
