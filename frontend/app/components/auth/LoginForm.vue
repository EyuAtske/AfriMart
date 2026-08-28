<script setup lang="ts">
const emit = defineEmits<{
  switchMode: []
}>()

const router = useRouter()

const loginForm = reactive({
  email: '',
  password: ''
})

const loginError = ref('')
const isLoading = ref(false)

const { login } = useAuth()

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
    // Try the current login function.
    // For now, we are using mock authentication.
    await login({
      email: loginForm.email.trim(),
      password: loginForm.password
    })
  } catch (error) {
    // Backend authentication is not connected yet,
    // so don't stop the frontend login flow.
    console.log('Mock login:', loginForm.email)
  } finally {
    // For the frontend/mock stage,
    // always go to the profile page after pressing Sign in.
    isLoading.value = false

    await router.push('/profile')
  }
}
</script>

<template>
  <form
    class="space-y-5"
    @submit.prevent="submitLogin"
  >
    <!-- Email -->
    <AuthInput
      v-model="loginForm.email"
      label="Email address"
      placeholder="Enter your email"
      type="email"
      autocomplete="email"
      name="email"
    />

    <!-- Password -->
    <AuthInput
      v-model="loginForm.password"
      label="Password"
      placeholder="Enter your password"
      type="password"
      autocomplete="current-password"
      name="password"
    />

    <!-- Error -->
    <p
      v-if="loginError"
      class="text-sm font-medium text-red-600"
    >
      {{ loginError }}
    </p>

    <!-- Forgot password -->
    <div class="flex items-center justify-end">
      <button
        type="button"
        class="text-sm text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline"
      >
        Forgot password?
      </button>
    </div>

    <!-- Sign In -->
    <button
      type="submit"
      :disabled="isLoading"
      class="h-12 w-full rounded-full border border-[#806344] text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
    >
      {{ isLoading ? 'Signing in...' : 'Sign in' }}
    </button>

    <!-- Register -->
    <p class="pt-3 text-center text-sm text-[#756a60]">
      Don't have an account?

      <button
        type="button"
        class="ml-1 font-medium text-[#806344] underline-offset-4 hover:text-[#211f1d] hover:underline"
        @click="emit('switchMode')"
      >
        Create account
      </button>
    </p>
  </form>
</template>