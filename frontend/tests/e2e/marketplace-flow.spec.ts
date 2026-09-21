import { test, expect } from '@playwright/test'

test.describe('End-to-End Buyer Journey Flow', () => {
  test.describe.configure({ mode: 'serial' })

  test('Complete flow: Browse -> Add to Cart -> View Cart -> Checkout -> Mock Payment -> Order Tracking', async ({ page }) => {
    await page.goto('/products')
   // await expect(page.locator('h1')).toContainText('Browse products')

    const addToCartButton = page.getByRole('button', { name: /add to cart/i }).first()
    await expect(addToCartButton).toBeVisible()
    await addToCartButton.click()

    const cartHeaderLink = page.locator('header a[href="/cart"]').first()
    await expect(cartHeaderLink).toBeVisible()
    await cartHeaderLink.click()

    await expect(page).toHaveURL('/cart')
    const checkoutLink = page.locator('a[href="/checkout"]').first()
    await expect(checkoutLink).toBeVisible({ timeout: 10000 })
    await checkoutLink.click()
    await expect(page).toHaveURL('/checkout')

    await page.fill('input[name="checkout-name"]', 'Abebe Bikila')
    await page.fill('input[name="checkout-phone"]', '+251 911 123 456')
    await page.fill('input[name="checkout-address"]', 'Bole Atlas')
    await page.fill('input[name="checkout-city"]', 'Addis Ababa')

    const confirmButton = page.getByRole('button', { name: /place order/i })
    await confirmButton.click()

    await expect(page).toHaveURL(/\/payment\/success/, { timeout: 10000 })
  //  await expect(page.locator('h1')).toContainText('Order confirmed')

    const trackOrderLink = page.locator('a[href="/orders"]').first()
    await expect(trackOrderLink).toBeVisible()
    await trackOrderLink.click()
    await expect(page).toHaveURL('/orders')
   // await expect(page.locator('h1')).toContainText('My Orders')

    const modalHeading = page.getByRole('heading', { name: /shipping & order tracking/i })
    await expect(modalHeading).not.toBeVisible()

    const trackShippingBtn = page.getByRole('button', { name: /track shipping/i }).first()
    await expect(trackShippingBtn).toBeVisible()
    await trackShippingBtn.click()

    await expect(modalHeading).toBeVisible()
    await expect(page.locator('text=Current Delivery Progress')).toBeVisible()

    const closeModalBtn = page.getByRole('button', { name: /close modal/i })
    await closeModalBtn.click()
    await expect(modalHeading).not.toBeVisible()
  })

  test('Product Card click navigates to /products/[id] and displays only that product details', async ({ page }) => {
    await page.goto('/products')
   // await expect(page.locator('h1')).toContainText('Browse products')

    const firstProductLink = page.locator('article a[href^="/products/"]').first()
    await expect(firstProductLink).toBeVisible()
    await firstProductLink.click()

    await expect(page).toHaveURL(/\/products\/\d+/)

    await expect(page.locator('h1')).not.toContainText('Browse products')
    await expect(page.getByRole('button', { name: /add to cart/i }).first()).toBeVisible()
    await expect(page.locator('text=Customer Reviews')).toBeVisible()
  })
})
