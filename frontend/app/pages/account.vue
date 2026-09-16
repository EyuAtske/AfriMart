<script setup lang="ts">
const { isLoggedIn } = useAuth()
const mode = ref<'login' | 'register'>('login')

const switchMode = (newMode: 'login' | 'register') => {
  mode.value = newMode
}

watch(isLoggedIn, (loggedIn) => {
  if (loggedIn) navigateTo('/profile')
}, { immediate: true })
</script>

<template>
  <main
    class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12"
  >
    <div
      v-if="!isLoggedIn"
      class="flex min-h-screen items-center justify-center"
    >
      <UiAppCard
        padding="large"
        class="w-full max-w-md"
      >
        <div class="mb-8 text-center">
          <h2 class="font-serif text-3xl text-[#211f1d]">
            {{ mode === 'login' ? 'Welcome back' : 'Create your account' }}
          </h2>

          <p class="mt-2 text-sm leading-6 text-[#756a60]">
            {{
              mode === 'login'
                ? 'Sign in to continue shopping.'
                : 'Join Afrimart and discover independent style.'
            }}
          </p>
        </div>

        <AuthLoginForm
          v-if="mode === 'login'"
          @switch-mode="switchMode('register')"
        />

        <AuthRegisterForm
          v-else
          @switch-mode="switchMode('login')"
        />
      </UiAppCard>
    </div>
  </main>
</template>

