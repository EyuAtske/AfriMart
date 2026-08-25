<script setup lang="ts">
const mode = ref<'login' | 'register'>('login')

const loginForm = reactive({
  email: '',
  password: '',
  remember: false
})

const registerForm = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const loginError = ref('')
const registerError = ref('')
const isLoading = ref(false)

const { login, register } = useAuth()

const switchMode = (newMode: 'login' | 'register') => {
  mode.value = newMode
  loginError.value = ''
  registerError.value = ''
}

const submitLogin = async () => {
  loginError.value = ''

  if (!loginForm.email.trim() || !loginForm.password.trim()) {
    loginError.value = 'Please fill in all fields.'
    return
  }

  if (!loginForm.email.includes('@')) {
    loginError.value = 'Please enter a valid email address.'
    return
  }

  isLoading.value = true

  try {
    await login({
      email: loginForm.email,
      password: loginForm.password
    })
  } catch (error: any) {
    loginError.value = error?.message || 'Failed to sign in. Please check your credentials.'
  } finally {
    isLoading.value = false
  }
}

const submitRegister = async () => {
  registerError.value = ''

  if (
    !registerForm.name.trim() ||
    !registerForm.email.trim() ||
    !registerForm.password.trim() ||
    !registerForm.confirmPassword.trim()
  ) {
    registerError.value = 'All fields are required and cannot be empty.'
    return
  }

  if (!registerForm.email.includes('@')) {
    registerError.value = 'Please enter a valid email address.'
    return
  }

  if (registerForm.password.length < 6) {
    registerError.value = 'Password must be at least 6 characters long.'
    return
  }

  if (registerForm.password !== registerForm.confirmPassword) {
    registerError.value = 'Passwords do not match.'
    return
  }

  isLoading.value = true

  try {
    await register({
      name: registerForm.name,
      email: registerForm.email,
      password: registerForm.password
    })
  } catch (error: any) {
    registerError.value = error?.message || 'Failed to create account. Please try again.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <main class="flex min-h-screen items-center justify-center bg-[#f5f1e9] px-4 py-12 sm:px-6 lg:px-8">
    <div class="w-full max-w-md rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-8 shadow-[0_20px_70px_rgba(33,31,29,0.08)] sm:p-10 transition-all duration-300">

      <!-- Header -->
      <div class="mb-8 text-center">
      

        <div>
          <h2 class="font-serif text-3xl tracking-[-0.025em] text-[#211f1d] transition-all duration-300">
            {{ mode === 'login' ? 'Welcome back' : 'Create your account' }}
          </h2>

          <p class="mt-2 text-sm leading-6 text-[#756a60] transition-all duration-300">
            {{
              mode === 'login'
                ? 'Sign in to continue shopping.'
                : 'Join Afrimart and discover independent style.'
            }}
          </p>
        </div>
      </div>

      <!-- Forms Section -->
      <div class="transition-all duration-300">
        <!-- Login -->
        <form
          v-if="mode === 'login'"
          class="space-y-5"
          @submit.prevent="submitLogin"
        >
          <AuthInput
            v-model="loginForm.email"
            label="Email address"
            placeholder="Enter your email"
            type="email"
            autocomplete="email"
            name="email"
          />

          <AuthInput
            v-model="loginForm.password"
            label="Password"
            placeholder="Enter your password"
            type="password"
            autocomplete="current-password"
            name="password"
          />

          <p v-if="loginError" class="text-xs text-red-600 font-medium">
            {{ loginError }}
          </p>

          <div class="flex items-center justify-between gap-4">
            <label class="flex cursor-pointer items-center gap-2 text-xs text-[#665c53]">
              <input
                v-model="loginForm.remember"
                type="checkbox"
                class="h-4 w-4 rounded border-[#c7bdb0] text-[#211f1d] focus:ring-[#806344]"
              />
              Remember me
            </label>

            <button
              type="button"
              class="text-xs text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline focus:outline-none focus:ring-2 focus:ring-[#806344]/20"
            >
              Forgot password?
            </button>
          </div>
          <button
            type="submit"
            :disabled="isLoading"
            class="h-12 w-full rounded-full bg-[#211f1d] text-xs font-medium uppercase tracking-[0.18em] text-white transition hover:bg-[#3b3733] focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {{ isLoading ? 'Signing in...' : 'Sign in' }}
          </button>

          <div class="flex items-center gap-4 py-2">
            <div class="h-px flex-1 bg-[#ded6cc]" />
            <span class="text-[10px] uppercase tracking-[0.18em] text-[#94877a]">
              or
            </span>
            <div class="h-px flex-1 bg-[#ded6cc]" />
          </div>

          <SocialButtons />

          <p class="pt-3 text-center text-xs text-[#756a60]">
            Don't have an account?

            <button
              type="button"
              class="ml-1 font-medium text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline focus:outline-none focus:ring-2 focus:ring-[#806344]/20"
              @click="switchMode('register')"
            >
              Create account
            </button>
          </p>
        </form>

        <!-- Register -->
        <form
          v-else
          class="space-y-5"
          @submit.prevent="submitRegister"
        >
          <AuthInput
            v-model="registerForm.name"
            label="Full name"
            placeholder="Enter your full name"
            autocomplete="name"
            name="name"
          />

          <AuthInput
            v-model="registerForm.email"
            label="Email address"
            placeholder="Enter your email"
            type="email"
            autocomplete="email"
            name="register-email"
          />

          <AuthInput
            v-model="registerForm.password"
            label="Password"
            placeholder="Create a password"
            type="password"
            autocomplete="new-password"
            name="register-password"
          />

          <AuthInput
            v-model="registerForm.confirmPassword"
            label="Confirm password"
            placeholder="Confirm your password"
            type="password"
            autocomplete="new-password"
            name="confirm-password"
          />

          <p v-if="registerError" class="text-xs text-red-600 font-medium">
            {{ registerError }}
          </p>
          <button
            type="submit"
            :disabled="isLoading"
            class="h-12 w-full rounded-full bg-[#211f1d] text-xs font-medium uppercase tracking-[0.18em] text-white transition hover:bg-[#3b3733] focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {{ isLoading ? 'Creating account...' : 'Create account' }}
          </button>

          <p class="text-center text-[11px] leading-5 text-[#81766c]">
            By creating an account, you agree to our
            <button
              type="button"
              class="text-[#806344] underline-offset-4 hover:underline"
            >
              Terms of Service
            </button>
            and
            <button type="button" class="text-[#806344] underline-offset-4 hover:underline">
              Privacy Policy
            </button>
          </p>

          <p class="pt-2 text-center text-xs text-[#756a60]">
            Already have an account?
            <button
              type="button"
              class="ml-1 font-medium text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline focus:outline-none focus:ring-2 focus:ring-[#806344]/20"
              @click="switchMode('login')"
            >
              Sign in
            </button>
          </p>
        </form>
      </div>

    </div>
  </main>
</template>