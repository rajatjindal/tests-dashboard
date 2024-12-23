// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  css: ['~/assets/css/main.css'],
  components: [
    {
      path: '~/components',
      pathPrefix: false,
    },
  ],
  extends: [
    'github:rajatjindal/nuxt-components'
  ],
  devtools: { enabled: true },

  imports: {
    dirs: [
      // Scan top-level modules
      // 'composables',
      // ... or scan modules nested one level deep with a specific name and file extension
      // 'composables/*/index.{ts,js,mjs,mts}',
      // ... or scan all modules within given directory
      'composables/**',
      'node_modules/.c12'
    ]
  },

  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },

  runtimeConfig: {
    public: {
      baseURL: ''
      // baseURL: 'http://localhost:3000'
      // baseURL: 'https://tests-dashboard.rajatjindal.com'
    }
  },

  ssr: false,

  typescript: {
    strict: true,
    typeCheck: true,
  },

  vite: {
    clearScreen: false
  },

  compatibilityDate: '2024-11-09'
})
