# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: home.spec.ts >> homepage loads
- Location: frontend\tests\home.spec.ts:2:1

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByRole('heading', { name: 'AFRIMART' })
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByRole('heading', { name: 'AFRIMART' })

```

# Test source

```ts
  1 | import { test, expect } from '@playwright/test'
  2 | test('homepage loads', async ({ page }) => {
  3 |    await expect(
> 4 |     page.getByRole('heading', { name: 'AFRIMART' })).toBeVisible()
    |                                                      ^ Error: expect(locator).toBeVisible() failed
  5 | })
  6 | 
```