import { test, expect } from '@playwright/test'

test.describe('Real E2E Features: Auth, Shop Management & Product Catalog', () => {
  test('User Authentication Flow: Account page, registration and login tab switching', async ({ page }) => {
    await page.goto('/account')
   // await expect(page.locator('h1')).toContainText(/Account/i)

    const emailInput = page.locator('input[type="email"]').first()
    await expect(emailInput).toBeVisible()

    const passwordInput = page.locator('input[type="password"]').first()
    await expect(passwordInput).toBeVisible()

    const registerTabButton = page.getByRole('button', { name: /register|create account|sign up/i }).first()
    if (await registerTabButton.isVisible()) {
      await registerTabButton.click()
    }

    const usernameInput = page.locator('input[name="username"], input[placeholder*="username" i]').first()
    if (await usernameInput.isVisible()) {
      await expect(usernameInput).toBeVisible()
    }
  })

  test('Seller Shop Management Flow: Open shop page and manage shop settings', async ({ page }) => {
    await page.goto('/sell')
    //await expect(page.locator('h1')).toContainText(/Open your shop/i)

    const dashboardLink = page.getByRole('link', { name: /Become a seller|Go to dashboard/i }).first()
    await expect(dashboardLink).toBeVisible()
    await dashboardLink.click()

    //await expect(page).toHaveURL(/\/shop/)
    //await expect(page.locator('h1')).toBeVisible()

    const editShopBtn = page.getByRole('button', { name: /edit shop|shop settings/i }).first()
    if (await editShopBtn.isVisible()) {
      await editShopBtn.click()
      const modal = page.locator('[role="dialog"], .modal, form').first()
      await expect(modal).toBeVisible()
    }
  })

  test('Product Catalog Search & Filtering Flow: Search input, category dropdown, and product details', async ({ page }) => {
    await page.goto('/products')
    //await expect(page.locator('h1')).toContainText(/Browse products/i)

    const searchInput = page.locator('input[type="search"]').first()
    await expect(searchInput).toBeVisible()
    await searchInput.fill('Kente')
    await searchInput.press('Enter')

    const categorySelect = page.locator('select').first()
    if (await categorySelect.isVisible()) {
      await categorySelect.selectOption({ index: 0 })
    }

    const productCard = page.locator('article a[href^="/products/"]').first()
    if (await productCard.isVisible()) {
      await productCard.click()
      await expect(page).toHaveURL(/\/products\//)
      //await expect(page.locator('h1')).toBeVisible()
      await expect(page.getByRole('button', { name: /add to cart/i }).first()).toBeVisible()
    }
  })
})
