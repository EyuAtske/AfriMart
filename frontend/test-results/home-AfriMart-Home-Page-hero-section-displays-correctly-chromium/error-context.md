# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: home.spec.ts >> AfriMart Home Page >> hero section displays correctly
- Location: frontend\tests\home.spec.ts:16:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/Discover clothing from independent sellers\./i)
Expected: visible
Error: strict mode violation: getByText(/Discover clothing from independent sellers\./i) resolved to 2 elements:
    1) <p class="mt-6 max-w-md text-sm leading-6 text-[#4c4945] sm:mt-8 sm:text-base sm:leading-7 lg:mt-10 lg:text-lg lg:leading-8">…</p> aka getByText('Discover clothing from independent sellers. Quality pieces. Great prices. One')
    2) <p class="mt-7 max-w-xs text-base leading-8 text-[#998d82]"> Discover clothing from independent sellers. Buy,…</p> aka getByText('Discover clothing from independent sellers. Buy, sell, and find your style in')

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/Discover clothing from independent sellers\./i)

```

# Page snapshot

```yaml
- generic [ref=e3]:
  - banner [ref=e4]:
    - navigation [ref=e5]:
      - generic [ref=e6]:
        - link "NEW IN" [ref=e7] [cursor=pointer]:
          - /url: /
        - link "MEN" [ref=e9] [cursor=pointer]:
          - /url: /men
        - link "WOMEN" [ref=e11] [cursor=pointer]:
          - /url: /women
        - link "KIDS" [ref=e13] [cursor=pointer]:
          - /url: /kids
      - link "AFRIMART" [ref=e14] [cursor=pointer]:
        - /url: /
      - generic [ref=e15]:
        - textbox "Search products or ask AI ✦..." [ref=e17]
        - link "Account" [ref=e18] [cursor=pointer]:
          - /url: /account
        - link "Shopping cart" [ref=e22] [cursor=pointer]:
          - /url: /cart
          - generic [ref=e26]: "2"
  - main [ref=e27]:
    - main [ref=e28]:
      - generic [ref=e29]:
        - img "Afrimart fashion collection" [ref=e30]
        - generic [ref=e33]:
          - heading "Your style. Your marketplace." [level=1] [ref=e34]
          - paragraph [ref=e35]: Discover clothing from independent sellers. Quality pieces. Great prices. One marketplace.
          - generic [ref=e36]:
            - link "SHOP NOW" [ref=e37] [cursor=pointer]:
              - /url: /products
            - link "EXPLORE ITEMS" [ref=e41] [cursor=pointer]:
              - /url: /products
      - generic [ref=e45]:
        - link "Men Men EXPLORE NOW" [ref=e46] [cursor=pointer]:
          - /url: /men
          - img "Men" [ref=e48]
          - generic [ref=e49]:
            - heading "Men" [level=2] [ref=e50]
            - generic [ref=e51]: EXPLORE NOW
        - link "Women Women EXPLORE NOW" [ref=e54] [cursor=pointer]:
          - /url: /women
          - img "Women" [ref=e56]
          - generic [ref=e57]:
            - heading "Women" [level=2] [ref=e58]
            - generic [ref=e59]: EXPLORE NOW
        - link "Kids Kids EXPLORE NOW" [ref=e62] [cursor=pointer]:
          - /url: /kids
          - img "Kids" [ref=e64]
          - generic [ref=e65]:
            - heading "Kids" [level=2] [ref=e66]
            - generic [ref=e67]: EXPLORE NOW
      - generic [ref=e70]:
        - generic [ref=e71]:
          - heading "Featured Products" [level=2] [ref=e73]
          - link "VIEW ALL" [ref=e75] [cursor=pointer]:
            - /url: /products
        - generic [ref=e78]:
          - article [ref=e80]:
            - generic [ref=e81]:
              - img "Shirt for men - image 1" [ref=e82]
              - button "Add to cart" [ref=e83]
            - generic [ref=e87]:
              - link "Atelier North" [ref=e88] [cursor=pointer]:
                - /url: /shops/atelier-north
              - heading "Shirt for men" [level=3] [ref=e89]
              - generic [ref=e90]:
                - paragraph [ref=e91]: 1,400 ETB
                - generic [ref=e92]:
                  - generic [ref=e93]: ★
                  - generic [ref=e94]: "4.8"
          - article [ref=e96]:
            - generic [ref=e97]:
              - img "Tank Tops - image 1" [ref=e98]
              - button "Add to cart" [ref=e99]
            - generic [ref=e103]:
              - link "Beyond Score" [ref=e104] [cursor=pointer]:
                - /url: /shops/beyond-score
              - heading "Tank Tops" [level=3] [ref=e105]
              - generic [ref=e106]:
                - paragraph [ref=e107]: 2,500 ETB
                - generic [ref=e108]:
                  - generic [ref=e109]: ★
                  - generic [ref=e110]: "4.9"
          - article [ref=e112]:
            - generic [ref=e113]:
              - img "Casual Outfit Set - image 1" [ref=e114]
              - button "Add to cart" [ref=e115]
            - generic [ref=e119]:
              - link "True Form" [ref=e120] [cursor=pointer]:
                - /url: /shops/true-form
              - heading "Casual Outfit Set" [level=3] [ref=e121]
              - generic [ref=e122]:
                - paragraph [ref=e123]: 3,500 ETB
                - generic [ref=e124]:
                  - generic [ref=e125]: ★
                  - generic [ref=e126]: "4.6"
          - article [ref=e128]:
            - generic [ref=e129]:
              - img "Everyday fit for kids - image 1" [ref=e130]
              - button "Add to cart" [ref=e131]
            - generic [ref=e135]:
              - link "Minimal Studio" [ref=e136] [cursor=pointer]:
                - /url: /shops/minimal-studio
              - heading "Everyday fit for kids" [level=3] [ref=e137]
              - generic [ref=e138]:
                - paragraph [ref=e139]: 2,900 ETB
                - generic [ref=e140]:
                  - generic [ref=e141]: ★
                  - generic [ref=e142]: "4.8"
          - article [ref=e144]:
            - generic [ref=e145]:
              - img "Cute dress for kids - image 1" [ref=e146]
              - button "Add to cart" [ref=e147]
            - generic [ref=e151]:
              - link "Atelier North" [ref=e152] [cursor=pointer]:
                - /url: /shops/atelier-north
              - heading "Cute dress for kids" [level=3] [ref=e153]
              - generic [ref=e154]:
                - paragraph [ref=e155]: 3,800 ETB
                - generic [ref=e156]:
                  - generic [ref=e157]: ★
                  - generic [ref=e158]: "4.7"
          - article [ref=e160]:
            - generic [ref=e161]:
              - img "Hoodie - image 1" [ref=e162]
              - button "Add to cart" [ref=e163]
            - generic [ref=e167]:
              - link "Urban Thread" [ref=e168] [cursor=pointer]:
                - /url: /shops/urban-thread
              - heading "Hoodie" [level=3] [ref=e169]
              - generic [ref=e170]:
                - paragraph [ref=e171]: 4,700 ETB
                - generic [ref=e172]:
                  - generic [ref=e173]: ★
                  - generic [ref=e174]: "4.8"
          - article [ref=e176]:
            - generic [ref=e177]:
              - img "Shirt for men - image 1" [ref=e178]
              - button "Add to cart" [ref=e179]
            - generic [ref=e183]:
              - link "Mara Studio" [ref=e184] [cursor=pointer]:
                - /url: /shops/mara-studio
              - heading "Shirt for men" [level=3] [ref=e185]
              - generic [ref=e186]:
                - paragraph [ref=e187]: 1,550 ETB
                - generic [ref=e188]:
                  - generic [ref=e189]: ★
                  - generic [ref=e190]: "4.9"
          - article [ref=e192]:
            - generic [ref=e193]:
              - img "Relaxed Denim - image 1" [ref=e194]
              - button "Add to cart" [ref=e195]
            - generic [ref=e199]:
              - link "ute dress" [ref=e200] [cursor=pointer]:
                - /url: /shops/ute-dress
              - heading "Relaxed Denim" [level=3] [ref=e201]
              - generic [ref=e202]:
                - paragraph [ref=e203]: 5,200 ETB
                - generic [ref=e204]:
                  - generic [ref=e205]: ★
                  - generic [ref=e206]: "4.7"
          - article [ref=e208]:
            - generic [ref=e209]:
              - img "Watch for women - image 1" [ref=e210]
              - button "Add to cart" [ref=e211]
            - generic [ref=e215]:
              - link "Minimal Studio" [ref=e216] [cursor=pointer]:
                - /url: /shops/minimal-studio
              - heading "Watch for women" [level=3] [ref=e217]
              - generic [ref=e218]:
                - paragraph [ref=e219]: 2,500 ETB
                - generic [ref=e220]:
                  - generic [ref=e221]: ★
                  - generic [ref=e222]: "4.9"
          - article [ref=e224]:
            - generic [ref=e225]:
              - img "Hat - image 1" [ref=e226]
              - button "Add to cart" [ref=e227]
            - generic [ref=e231]:
              - link "Mara Studio" [ref=e232] [cursor=pointer]:
                - /url: /shops/mara-studio
              - heading "Hat" [level=3] [ref=e233]
              - generic [ref=e234]:
                - paragraph [ref=e235]: 1,100 ETB
                - generic [ref=e236]:
                  - generic [ref=e237]: ★
                  - generic [ref=e238]: "4.8"
      - button "LOAD MORE" [ref=e240]
      - generic [ref=e245]:
        - generic [ref=e246]:
          - generic [ref=e247]: ✦
          - generic [ref=e248]:
            - heading "Looking for something specific?" [level=2] [ref=e249]
            - paragraph [ref=e250]: Tell us what you're looking for and our AI assistant will help you find the perfect items.
        - generic [ref=e251]:
          - textbox "What are you looking for?" [ref=e252]
          - button "Search with AI" [ref=e253]
  - contentinfo [ref=e256]:
    - generic [ref=e258]:
      - generic [ref=e259]:
        - link "AFRIMART" [ref=e260] [cursor=pointer]:
          - /url: /
        - paragraph [ref=e261]: Discover clothing from independent sellers. Buy, sell, and find your style in one marketplace.
      - generic [ref=e262]:
        - heading "SHOP" [level=3] [ref=e263]
        - navigation [ref=e264]:
          - link "New Arrivals" [ref=e265] [cursor=pointer]:
            - /url: /products
          - link "Men" [ref=e266] [cursor=pointer]:
            - /url: /men
          - link "Women" [ref=e267] [cursor=pointer]:
            - /url: /women
          - link "Kids" [ref=e268] [cursor=pointer]:
            - /url: /kids
          - link "Shopping Cart" [ref=e269] [cursor=pointer]:
            - /url: /cart
      - generic [ref=e270]:
        - heading "SELL" [level=3] [ref=e271]
        - navigation [ref=e272]:
          - link "Start Selling" [ref=e273] [cursor=pointer]:
            - /url: /sell
          - link "My Shop" [ref=e274] [cursor=pointer]:
            - /url: /shop
          - link "My Products" [ref=e275] [cursor=pointer]:
            - /url: /seller/products
          - link "Seller Orders" [ref=e276] [cursor=pointer]:
            - /url: /seller/orders
      - generic [ref=e277]:
        - heading "ACCOUNT" [level=3] [ref=e278]
        - navigation [ref=e279]:
          - link "My Profile" [ref=e280] [cursor=pointer]:
            - /url: /profile
          - link "My Orders" [ref=e281] [cursor=pointer]:
            - /url: /orders
          - link "AI Shopping Assistant" [ref=e282] [cursor=pointer]:
            - /url: /ai
```

# Test source

```ts
  1  | import { test, expect } from '@playwright/test'
  2  | 
  3  | test.describe('AfriMart Home Page', () => {
  4  |   test.beforeEach(async ({ page }) => {
  5  |     await page.goto('/')
  6  |   })
  7  | 
  8  |   test('homepage loads', async ({ page }) => {
  9  |     await expect(
  10 |       page.getByRole('heading', {
  11 |         name: /Your style.*Your.*marketplace/i,
  12 |       }),
  13 |     ).toBeVisible()
  14 |   })
  15 | 
  16 |   test('hero section displays correctly', async ({ page }) => {
  17 |     await expect(
  18 |       page.getByAltText('Afrimart fashion collection'),
  19 |     ).toBeVisible()
  20 | 
  21 |     await expect(
  22 |       page.getByText(
  23 |         /Discover clothing from independent sellers\./i,
  24 |       ),
> 25 |     ).toBeVisible()
     |       ^ Error: expect(locator).toBeVisible() failed
  26 |   })
  27 | 
  28 |   test('shop now link is visible and points to products', async ({ page }) => {
  29 |     const shopNow = page.getByRole('link', {
  30 |       name: /shop now/i,
  31 |     })
  32 | 
  33 |     await expect(shopNow).toBeVisible()
  34 |     await expect(shopNow).toHaveAttribute('href', '/products')
  35 |   })
  36 | 
  37 |   test('explore items link is visible and points to products', async ({
  38 |     page,
  39 |   }) => {
  40 |     const exploreItems = page.getByRole('link', {
  41 |       name: /explore items/i,
  42 |     })
  43 | 
  44 |     await expect(exploreItems).toBeVisible()
  45 |     await expect(exploreItems).toHaveAttribute('href', '/products')
  46 |   })
  47 | 
  48 |   test('main clothing categories are visible', async ({ page }) => {
  49 |     await expect(
  50 |       page.getByText('Men', { exact: true }),
  51 |     ).toBeVisible()
  52 | 
  53 |     await expect(
  54 |       page.getByText('Women', { exact: true }),
  55 |     ).toBeVisible()
  56 | 
  57 |     await expect(
  58 |       page.getByText('Kids', { exact: true }),
  59 |     ).toBeVisible()
  60 |   })
  61 | 
  62 |   test('load more button is visible', async ({ page }) => {
  63 |     await expect(
  64 |       page.getByRole('button', {
  65 |         name: /load more/i,
  66 |       }),
  67 |     ).toBeVisible()
  68 |   })
  69 | })
```