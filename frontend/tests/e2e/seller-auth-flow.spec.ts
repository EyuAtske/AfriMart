import { test, expect } from '@playwright/test'

test.describe('Seller Authentication & Shop Setup E2E Flow', () => {
  test('Account Login / Registration and Shop Dashboard Access', async ({ page }) => {
    await page.goto('/account')
   // await expect(page.locator('h1')).toContainText(/Account/i)

    await page.goto('/sell')
    //await expect(page.locator('h1')).toContainText(/Open your shop/i)

    const sellerCta = page.getByRole('link', { name: /Become a seller|Go to dashboard/i }).first()
    await expect(sellerCta).toBeVisible()
    await sellerCta.click()

    
  })

  test('Public Home Page navigation to Catalog and Categories', async ({ page }) => {
    await page.goto('/')
    await expect(page.locator('body')).toBeVisible()

    const shopLink = page.locator('a[href="/products"]').first()
    await expect(shopLink).toBeVisible()
    await shopLink.click()
    await expect(page).toHaveURL('/products')
  })
})
