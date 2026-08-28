import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',

  css: ['~/assets/css/main.css'],
  devtools: {
    enabled: true
  },
  runtimeConfig: { 
    public: { 
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
      authMode: process.env.NUXT_PUBLIC_AUTH_MODE || 'mock'
    }
  },

  ssr: true,

  vite: {
    plugins: [
      
      tailwindcss()
    ]
  }
})