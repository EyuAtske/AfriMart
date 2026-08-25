<script setup lang="ts">
definePageMeta({
  layout: 'auth'
})
import { ref } from 'vue'

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const showPassword = ref(false)
const showConfirmPassword = ref(false)

const handleSignup = () => {
  console.log('Signup:', {
    name: name.value,
    email: email.value,
    password: password.value,
    confirmPassword: confirmPassword.value
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

        <!-- Form -->
        <div class="my-auto w-full max-w-xl py-12">

          <!-- Heading -->
          <div class="mb-10">
            <p
              class="mb-4 text-xs font-medium uppercase tracking-[0.22em] text-[#806344]"
            >
              JOIN AFRIMART
            </p>

            <h1
              class="font-serif text-5xl leading-none tracking-[-0.025em] text-[#211f1d] sm:text-6xl"
            >
              Create your account.
            </h1>

            <p class="mt-5 text-base leading-7 text-[#806344]">
              Join our marketplace and discover clothing from independent sellers.
            </p>
          </div>

          <form
            class="space-y-6"
            @submit.prevent="handleSignup"
          >

            <!-- Name -->
            <div>
              <label
                for="name"
                class="mb-3 block text-xs font-medium uppercase tracking-[0.18em] text-[#4c4945]"
              >
                Full Name
              </label>

              <input
                id="name"
                v-model="name"
                type="text"
                autocomplete="name"
                placeholder="Your full name"
                required
                class="h-16 w-full border border-[#d8cbbb] bg-[#eee5d5]/60 px-5 text-base text-[#211f1d] outline-none transition placeholder:text-[#a8957d] focus:border-[#806344] focus:bg-[#f5f1e9]"
              />
            </div>

            <!-- Email -->
            <div>
              <label
                for="signup-email"
                class="mb-3 block text-xs font-medium uppercase tracking-[0.18em] text-[#4c4945]"
              >
                Email Address
              </label>

              <input
                id="signup-email"
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
              <label
                for="signup-password"
                class="mb-3 block text-xs font-medium uppercase tracking-[0.18em] text-[#4c4945]"
              >
                Password
              </label>

              <div class="relative">
                <input
                  id="signup-password"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="new-password"
                  placeholder="Create a password"
                  required
                  class="h-16 w-full border border-[#d8cbbb] bg-[#eee5d5]/60 px-5 pr-14 text-base text-[#211f1d] outline-none transition placeholder:text-[#a8957d] focus:border-[#806344] focus:bg-[#f5f1e9]"
                />

                <button
                  type="button"
                  aria-label="Toggle password visibility"
                  class="absolute right-5 top-1/2 -translate-y-1/2 text-[#806344]"
                  @click="showPassword = !showPassword"
                >
                  <svg
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
                </button>
              </div>
            </div>

            <!-- Confirm password -->
            <div>
              <label
                for="confirm-password"
                class="mb-3 block text-xs font-medium uppercase tracking-[0.18em] text-[#4c4945]"
              >
                Confirm Password
              </label>

              <div class="relative">
                <input
                  id="confirm-password"
                  v-model="confirmPassword"
                  :type="showConfirmPassword ? 'text' : 'password'"
                  autocomplete="new-password"
                  placeholder="Confirm your password"
                  required
                  class="h-16 w-full border border-[#d8cbbb] bg-[#eee5d5]/60 px-5 pr-14 text-base text-[#211f1d] outline-none transition placeholder:text-[#a8957d] focus:border-[#806344] focus:bg-[#f5f1e9]"
                />

                <button
                  type="button"
                  aria-label="Toggle password visibility"
                  class="absolute right-5 top-1/2 -translate-y-1/2 text-[#806344]"
                  @click="showConfirmPassword = !showConfirmPassword"
                >
                  <svg
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
                </button>
              </div>
            </div>

            <!-- Terms -->
            <label class="flex items-start gap-3 pt-1 text-sm leading-6 text-[#806344]">
              <input
                type="checkbox"
                required
                class="mt-1 h-4 w-4 accent-[#211f1d]"
              />

              <span>
                I agree to AfriMart's
                <NuxtLink
                  to="/terms"
                  class="text-[#211f1d] underline underline-offset-4"
                >
                  terms
                </NuxtLink>
                and
                <NuxtLink
                  to="/privacy"
                  class="text-[#211f1d] underline underline-offset-4"
                >
                  privacy policy
                </NuxtLink>.
              </span>
            </label>

            <!-- Submit -->
            <button
              type="submit"
              class="h-16 w-full bg-[#211f1d] text-sm font-medium tracking-[0.18em] text-[#f5f1e9] transition duration-300 hover:bg-[#806344]"
            >
              CREATE ACCOUNT
            </button>
          </form>

          <!-- Divider -->
          <div class="my-8 flex items-center gap-5">
            <div class="h-px flex-1 bg-[#d8cbbb]" />
            <span
              class="text-xs uppercase tracking-[0.16em] text-[#806344]"
            >
              OR
            </span>
            <div class="h-px flex-1 bg-[#d8cbbb]" />
          </div>

          <!-- Login -->
          <p class="text-center text-sm text-[#806344]">
            Already have an account?

            <NuxtLink
              to="/auth/login"
              class="ml-1 font-medium text-[#211f1d] underline decoration-[#806344] underline-offset-4 transition hover:text-[#806344]"
            >
              Sign in
            </NuxtLink>
          </p>
        </div>

        <!-- Footer -->
        <p class="text-center text-xs text-[#806344]/70 lg:text-left">
          © {{ new Date().getFullYear() }} AfriMart. All rights reserved.
        </p>
      </section>

      <!-- Image side -->
      <AuthSidePanel mode="signup" />

    </div>
  </main>
</template>