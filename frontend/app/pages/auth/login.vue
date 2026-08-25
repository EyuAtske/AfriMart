<script setup lang="ts">
definePageMeta({
  layout: 'auth'
})
import { ref } from 'vue'

const email = ref('')
const password = ref('')
const showPassword = ref(false)

const handleLogin = () => {
  console.log('Login:', {
    email: email.value,
    password: password.value
  })

  // Backend will be connected here later.
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9]">
    <div class="flex min-h-screen">

      <!-- Form side -->
      <section
        class="flex min-h-screen w-full flex-col px-6 py-10 sm:px-10 lg:w-1/2 lg:px-16 xl:px-24"
      >

        <!-- Logo -->
        <NuxtLink
          to="/"
          class="font-serif text-2xl tracking-[0.35em] text-[#211f1d]"
        >
          AFRIMART
        </NuxtLink>

        <!-- Form container -->
        <div class="my-auto w-full max-w-xl py-12">

          <!-- Heading -->
          <div class="mb-12">
            <p
              class="mb-4 text-xs font-medium uppercase tracking-[0.22em] text-[#806344]"
            >
              ACCOUNT
            </p>

            <h1
              class="font-serif text-5xl leading-none tracking-[-0.025em] text-[#211f1d] sm:text-6xl"
            >
              Welcome back.
            </h1>

            <p class="mt-5 text-base leading-7 text-[#806344]">
              Sign in to continue shopping your favorite styles.
            </p>
          </div>

          <!-- Form -->
          <form
            class="space-y-7"
            @submit.prevent="handleLogin"
          >

            <!-- Email -->
            <div>
              <label
                for="email"
                class="mb-3 block text-xs font-medium uppercase tracking-[0.18em] text-[#4c4945]"
              >
                Email Address
              </label>

              <input
                id="email"
                v-model="email"
                type="email"
                autocomplete="email"
                placeholder="you@example.com"
                required
                class="h-16 w-full border border-[#d8cbbb] bg-[#eee5d5]/60 px-5 text-base text-[#211f1d] outline-none transition placeholder:text-[#a8957d] focus:border-[#806344] focus:bg-[#f5f1e9]"
              />
            </div>

            <!-- Password -->
            <div>
              <div class="mb-3 flex items-center justify-between">
                <label
                  for="password"
                  class="text-xs font-medium uppercase tracking-[0.18em] text-[#4c4945]"
                >
                  Password
                </label>

                <button
                  type="button"
                  class="text-xs tracking-[0.08em] text-[#806344] transition hover:text-[#211f1d]"
                >
                  Forgot password?
                </button>
              </div>

              <div class="relative">
                <input
                  id="password"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  placeholder="••••••••"
                  required
                  class="h-16 w-full border border-[#d8cbbb] bg-[#eee5d5]/60 px-5 pr-14 text-base text-[#211f1d] outline-none transition placeholder:text-[#a8957d] focus:border-[#806344] focus:bg-[#f5f1e9]"
                />

                <button
                  type="button"
                  aria-label="Toggle password visibility"
                  class="absolute right-5 top-1/2 -translate-y-1/2 text-[#806344] transition hover:text-[#211f1d]"
                  @click="showPassword = !showPassword"
                >
                  <svg
                    v-if="!showPassword"
                    xmlns="http://www.w3.org/2000/svg"
                    width="21"
                    height="21"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>

                  <svg
                    v-else
                    xmlns="http://www.w3.org/2000/svg"
                    width="21"
                    height="21"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path d="M3 3l18 18" />
                    <path d="M10.6 10.6a2 2 0 0 0 2.8 2.8" />
                    <path d="M9.9 5.1A10.8 10.8 0 0 1 12 5c6.5 0 10 7 10 7a18.3 18.3 0 0 1-3.2 4.1" />
                    <path d="M6.1 6.1C3.3 8.1 2 12 2 12s3.5 7 10 7a9.8 9.8 0 0 0 4.1-.9" />
                  </svg>
                </button>
              </div>
            </div>

            <!-- Submit -->
            <button
              type="submit"
              class="h-16 w-full bg-[#211f1d] text-sm font-medium tracking-[0.18em] text-[#f5f1e9] transition duration-300 hover:bg-[#806344]"
            >
              SIGN IN
            </button>
          </form>

          <!-- Divider -->
          <div class="my-9 flex items-center gap-5">
            <div class="h-px flex-1 bg-[#d8cbbb]" />
            <span
              class="text-xs uppercase tracking-[0.16em] text-[#806344]"
            >
              OR
            </span>
            <div class="h-px flex-1 bg-[#d8cbbb]" />
          </div>

          <!-- Signup -->
          <p class="text-center text-sm text-[#806344]">
            Don't have an account?

            <NuxtLink
              to="/auth/signup"
              class="ml-1 font-medium text-[#211f1d] underline decoration-[#806344] underline-offset-4 transition hover:text-[#806344]"
            >
              Create one
            </NuxtLink>
          </p>
        </div>

        <!-- Footer -->
        <p class="text-center text-xs text-[#806344]/70 lg:text-left">
          © {{ new Date().getFullYear() }} AfriMart. All rights reserved.
        </p>
      </section>

      <!-- Image side -->
      <AuthSidePanel mode="login" />

    </div>
  </main>
</template>