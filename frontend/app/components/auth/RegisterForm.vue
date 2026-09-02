<script setup lang="ts">
const emit = defineEmits<{
  switchMode: []
}>()

const registerForm = reactive({
  username: '',
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const registerError = ref('')
const isLoading = ref(false)

const { register } = useAuth()

const submitRegister = async () => {
  registerError.value = ''

  if (
    !registerForm.username.trim() ||
    !registerForm.firstName.trim() ||
    !registerForm.lastName.trim() ||
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
    const session = await register({
      username: registerForm.username.trim(),
      firstName: registerForm.firstName.trim(),
      lastName: registerForm.lastName.trim(),
      email: registerForm.email.trim(),
      password: registerForm.password
    })

    if (!session?.token) {
      const { showToast } = useToast()
      showToast('Account created successfully! Please sign in with your credentials.')
      emit('switchMode')
    }
  } catch (error: any) {
    registerError.value =
      error?.message || 'Failed to create account. Please try again.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <form
    class="space-y-5"
    @submit.prevent="submitRegister"
  >

    <!-- Username -->
    <AuthInput
      v-model="registerForm.username"
      label="Username"
      placeholder="Choose a username"
      type="text"
      autocomplete="username"
      name="register-username"
    />

    <!-- First + Last Name -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
      <AuthInput
        v-model="registerForm.firstName"
        label="First name"
        placeholder="First name"
        type="text"
        autocomplete="given-name"
        name="register-first-name"
      />

      <AuthInput
        v-model="registerForm.lastName"
        label="Last name"
        placeholder="Last name"
        type="text"
        autocomplete="family-name"
        name="register-last-name"
      />
    </div>

    <!-- Email -->
    <AuthInput
      v-model="registerForm.email"
      label="Email address"
      placeholder="Enter your email"
      type="email"
      autocomplete="email"
      name="register-email"
    />

    <!-- Password -->
    <AuthInput
      v-model="registerForm.password"
      label="Password"
      placeholder="Create a password"
      type="password"
      autocomplete="new-password"
      name="register-password"
    />

    <!-- Confirm Password -->
    <AuthInput
      v-model="registerForm.confirmPassword"
      label="Confirm password"
      placeholder="Confirm your password"
      type="password"
      autocomplete="new-password"
      name="register-confirm-password"
    />

    <!-- Error -->
    <p
      v-if="registerError"
      class="text-sm font-medium text-red-600"
    >
      {{ registerError }}
    </p>

    
    <!-- Submit -->
    <button
      type="submit"
      :disabled="isLoading"
      class="h-14 w-full rounded-full border border-[#806344] text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
    >
      {{ isLoading ? 'Creating account...' : 'Create account' }}
    </button>

    <!-- Login -->
    <p class="pt-2 text-center text-sm text-[#756a60]">
      Already have an account?

      <button
        type="button"
        class="ml-1 font-medium text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline"
        @click="emit('switchMode')"
      >
        Sign in
      </button>
    </p>

  </form>
</template>