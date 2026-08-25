# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: home.spec.ts >> AfriMart Home Page >> home page loads successfully
- Location: frontend\tests\home.spec.ts:8:3

# Error details

```
Error: expect(page).toHaveTitle(expected) failed

Expected pattern: /AfriMart/i
Received string:  ""
Timeout: 5000ms

Call log:
  - Expect "toHaveTitle" with timeout 5000ms
    13 × locator resolved to <html>…</html>
       - unexpected value ""

```

```yaml
- banner:
  - navigation:
    - link "NEW IN":
      - /url: /
    - link "MEN":
      - /url: /men
    - link "WOMEN":
      - /url: /women
    - link "KIDS":
      - /url: /kids
    - link "AFRIMART":
      - /url: /
    - img
    - textbox "Search products or ask AI ✦..."
    - link "Account":
      - /url: /account
      - img
    - link "Shopping cart":
      - /url: /cart
      - img
      - text: "2"
- main:
  - main:
    - img "Afrimart fashion collection"
    - heading "Your style. Your marketplace." [level=1]
    - paragraph: Discover clothing from independent sellers. Quality pieces. Great prices. One marketplace.
    - link "SHOP NOW":
      - /url: /products
      - text: SHOP NOW
      - img
    - link "EXPLORE ITEMS":
      - /url: /products
      - text: EXPLORE ITEMS
      - img
    - link "Men Men EXPLORE NOW":
      - /url: /men
      - img "Men"
      - heading "Men" [level=2]
      - text: EXPLORE NOW
    - link "Women Women EXPLORE NOW":
      - /url: /women
      - img "Women"
      - heading "Women" [level=2]
      - text: EXPLORE NOW
    - link "Kids Kids EXPLORE NOW":
      - /url: /kids
      - img "Kids"
      - heading "Kids" [level=2]
      - text: EXPLORE NOW
    - heading "Featured Products" [level=2]
    - link "VIEW ALL":
      - /url: /products
      - text: VIEW ALL
      - img
    - article:
      - img "Shirt for men - image 1"
      - button "Add to cart":
        - img
      - link "Atelier North":
        - /url: /shops/atelier-north
      - heading "Shirt for men" [level=3]
      - paragraph: 1,400 ETB
      - text: ★ 4.8
    - article:
      - img "Tank Tops - image 1"
      - button "Add to cart":
        - img
      - link "Beyond Score":
        - /url: /shops/beyond-score
      - heading "Tank Tops" [level=3]
      - paragraph: 2,500 ETB
      - text: ★ 4.9
    - article:
      - img "Casual Outfit Set - image 1"
      - button "Add to cart":
        - img
      - link "True Form":
        - /url: /shops/true-form
      - heading "Casual Outfit Set" [level=3]
      - paragraph: 3,500 ETB
      - text: ★ 4.6
    - article:
      - img "Everyday fit for kids - image 1"
      - button "Add to cart":
        - img
      - link "Minimal Studio":
        - /url: /shops/minimal-studio
      - heading "Everyday fit for kids" [level=3]
      - paragraph: 2,900 ETB
      - text: ★ 4.8
    - article:
      - img "Cute dress for kids - image 1"
      - button "Add to cart":
        - img
      - link "Atelier North":
        - /url: /shops/atelier-north
      - heading "Cute dress for kids" [level=3]
      - paragraph: 3,800 ETB
      - text: ★ 4.7
    - article:
      - img "Hoodie - image 1"
      - button "Add to cart":
        - img
      - link "Urban Thread":
        - /url: /shops/urban-thread
      - heading "Hoodie" [level=3]
      - paragraph: 4,700 ETB
      - text: ★ 4.8
    - article:
      - img "Shirt for men - image 1"
      - button "Add to cart":
        - img
      - link "Mara Studio":
        - /url: /shops/mara-studio
      - heading "Shirt for men" [level=3]
      - paragraph: 1,550 ETB
      - text: ★ 4.9
    - article:
      - img "Relaxed Denim - image 1"
      - button "Add to cart":
        - img
      - link "ute dress":
        - /url: /shops/ute-dress
      - heading "Relaxed Denim" [level=3]
      - paragraph: 5,200 ETB
      - text: ★ 4.7
    - article:
      - img "Watch for women - image 1"
      - button "Add to cart":
        - img
      - link "Minimal Studio":
        - /url: /shops/minimal-studio
      - heading "Watch for women" [level=3]
      - paragraph: 2,500 ETB
      - text: ★ 4.9
    - article:
      - img "Hat - image 1"
      - button "Add to cart":
        - img
      - link "Mara Studio":
        - /url: /shops/mara-studio
      - heading "Hat" [level=3]
      - paragraph: 1,100 ETB
      - text: ★ 4.8
    - button "LOAD MORE":
      - text: LOAD MORE
      - img
    - text: ✦
    - heading "Looking for something specific?" [level=2]
    - paragraph: Tell us what you're looking for and our AI assistant will help you find the perfect items.
    - textbox "What are you looking for?"
    - button "Search with AI":
      - img
- contentinfo:
  - link "AFRIMART":
    - /url: /
  - paragraph: Discover clothing from independent sellers. Buy, sell, and find your style in one marketplace.
  - heading "SHOP" [level=3]
  - navigation:
    - link "New Arrivals":
      - /url: /products
    - link "Men":
      - /url: /men
    - link "Women":
      - /url: /women
    - link "Kids":
      - /url: /kids
    - link "Shopping Cart":
      - /url: /cart
  - heading "SELL" [level=3]
  - navigation:
    - link "Start Selling":
      - /url: /sell
    - link "My Shop":
      - /url: /shop
    - link "My Products":
      - /url: /seller/products
    - link "Seller Orders":
      - /url: /seller/orders
  - heading "ACCOUNT" [level=3]
  - navigation:
    - link "My Profile":
      - /url: /profile
    - link "My Orders":
      - /url: /orders
    - link "AI Shopping Assistant":
      - /url: /ai
- img
- button "Toggle Nuxt DevTools":
  - img
- text: 165 ms
- button "Toggle Component Inspector":
  - img
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
  8  |   test('home page loads successfully', async ({ page }) => {
> 9  |     await expect(page).toHaveTitle(/AfriMart/i)
     |                        ^ Error: expect(page).toHaveTitle(expected) failed
  10 |   })
  11 | 
  12 |   test('displays the hero section', async ({ page }) => {
  13 |     await expect(
  14 |       page.getByRole('heading', {
  15 |         name: /Your style\. Your marketplace\./i,
  16 |       }),
  17 |     ).toBeVisible()
  18 | 
  19 |     await expect(
  20 |       page.getByText(
  21 |         /Discover clothing from independent sellers\./i,
  22 |       ),
  23 |     ).toBeVisible()
  24 | 
  25 |     await expect(
  26 |       page.getByAltText('Afrimart fashion collection'),
  27 |     ).toBeVisible()
  28 |   })
  29 | 
  30 |   test('displays the shop now and explore items links', async ({ page }) => {
  31 |     const shopNow = page.getByRole('link', { name: /shop now/i })
  32 |     const exploreItems = page.getByRole('link', { name: /explore items/i })
  33 | 
  34 |     await expect(shopNow).toBeVisible()
  35 |     await expect(shopNow).toHaveAttribute('href', '/products')
  36 | 
  37 |     await expect(exploreItems).toBeVisible()
  38 |     await expect(exploreItems).toHaveAttribute('href', '/products')
  39 |   })
  40 | 
  41 |   test('displays the main clothing categories', async ({ page }) => {
  42 |     await expect(page.getByText('Men', { exact: true })).toBeVisible()
  43 |     await expect(page.getByText('Women', { exact: true })).toBeVisible()
  44 |     await expect(page.getByText('Kids', { exact: true })).toBeVisible()
  45 |   })
  46 | 
  47 |   test('displays the featured products section', async ({ page }) => {
  48 |     await expect(
  49 |       page.getByText(/featured products/i),
  50 |     ).toBeVisible()
  51 |   })
  52 | 
  53 |   test('displays the load more button', async ({ page }) => {
  54 |     await expect(
  55 |       page.getByRole('button', { name: /load more/i }),
  56 |     ).toBeVisible()
  57 |   })
  58 | 
  59 |   test('displays the AI search section', async ({ page }) => {
  60 |     await expect(
  61 |       page.getByText(/AI/i).first(),
  62 |     ).toBeVisible()
  63 |   })
  64 | 
  65 |   test('displays the footer', async ({ page }) => {
  66 |     await expect(
  67 |       page.getByRole('contentinfo'),
  68 |     ).toBeVisible()
  69 |   })
  70 | })
  71 | 
  72 | 
```