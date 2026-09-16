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
    await login({
      email: loginForm.email.trim(),
      password: loginForm.password
    })
  } catch (error: any) {
    loginError.value =
      error?.message || 'Login failed. Please check your credentials.'
  } finally {
    isLoading.value = false
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
    <UiAppAlert v-if="loginError">
      {{ loginError }}
    </UiAppAlert>

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
    <UiAppButton
      type="submit"
      :disabled="isLoading"
      class="w-full"
    >
      {{ isLoading ? 'Signing in...' : 'Sign in' }}
    </UiAppButton>

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