import { defineConfig } from 'vitest/config'
import path from 'path'

export default defineConfig({
  test: {
    include: ['tests/unit/**/*.spec.ts'],
    globals: true,
  },
  define: {
    'import.meta.client': true,
  },
  resolve: {
    alias: {
      '~': path.resolve('./app'),
      '@': path.resolve('./app'),
    },
  },
})
