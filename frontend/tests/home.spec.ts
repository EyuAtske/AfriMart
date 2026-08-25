import {test,expect} from '@playwright/test'

test('navbar displays correctly', async ({ page }) => {
  await page.goto('/')
  await expect( page.getByRole('link', { name: 'AFRIMART' }).first()).toBeVisible()
  await expect (page.getByRole('link',{name: 'NEW IN'}).first()).toBeVisible()
  await expect (page.getByRole('link',{name: 'MEN'}).first()).toBeVisible()
  await expect (page.getByRole('link',{name: 'WOMEN'}).first()).toBeVisible()
  await expect (page.getByRole('link',{name: 'KIDS'}).first()).toBeVisible()
  await expect(page.getByPlaceholder('Search products or ask AI ✦...')).toBeVisible()
  await expect( page.getByRole('link', { name: 'Account' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Shopping cart' }).first()).toBeVisible()
  await expect(page.getByText('YOUR STYLE.').first()).toBeVisible()
  await expect(page.getByRole('link',{name:'EXPLORE ITEMS'})).toBeVisible()
})
