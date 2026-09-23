import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',

  modules: [
    'nuxt-gtag'
  ],

  css: ['~/assets/css/main.css'],

  devtools: {
    enabled: true
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '',
      authMode: process.env.NUXT_PUBLIC_AUTH_MODE || 'api',
      posthogKey: process.env.NUXT_PUBLIC_POSTHOG_KEY || '',
      minioBaseUrl: process.env.NUXT_PUBLIC_MINIO_BASE_URL || 'http://localhost:9000/afrimart-images'
    }
  },
  ssr: true,

  routeRules: {
    '/api/**': { proxy: `${process.env.API_PROXY_TARGET || 'http://localhost:8080'}/api/**` }
  },

  vite: {
    plugins: [
      tailwindcss()
    ]
  },

  gtag: {
    id: 'G-PGQMZFY1KF'
  }
})