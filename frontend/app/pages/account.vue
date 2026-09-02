<script setup lang="ts">
const { isLoggedIn, user, logout } = useAuth()
const mode = ref<'login' | 'register'>('login')

const switchMode = (newMode: 'login' | 'register') => {
  mode.value = newMode
}
</script>

<template>
  <main
    class="min-h-screen bg-[#f5f1e9] px-4 py-12 sm:px-6 lg:px-8"
  >
    <!-- Logged in -->
    <div
      v-if="isLoggedIn"
      class="mx-auto w-full max-w-5xl"
    >
      <div
        class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-8 shadow-[0_20px_70px_rgba(33,31,29,0.08)] sm:p-10"
      >
        <h1 class="font-serif text-3xl text-[#211f1d]">
          My Account
        </h1>

        <p class="mt-2 text-sm text-[#756a60]">
          Welcome back, {{ user.name }}.
        </p>

        <div class="mt-8 border-t border-[#ded6cc] pt-6">
          <h2 class="text-sm font-medium uppercase tracking-[0.12em] text-[#211f1d]">
            Profile
          </h2>

          <div class="mt-4">
            <p class="text-sm text-[#756a60]">
              Name
            </p>

            <p class="mt-1 text-sm text-[#211f1d]">
              {{ user.name }}
            </p>
          </div>

          <div class="mt-4">
            <p class="text-sm text-[#756a60]">
              Email
            </p>

            <p class="mt-1 text-sm text-[#211f1d]">
              {{ user.email }}
            </p>
          </div>
        </div>

        <button
          type="button"
          class="mt-8 rounded-full bg-[#211f1d] px-6 py-3 text-xs font-medium uppercase tracking-[0.15em] text-white transition hover:bg-[#3b3733]"
          @click="logout"
        >
          Logout
        </button>
      </div>
    </div>

    <!-- Logged out -->
    <div
      v-else
      class="flex min-h-screen items-center justify-center"
    >
      <div
        class="w-full max-w-md rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-8 shadow-[0_20px_70px_rgba(33,31,29,0.08)] sm:p-10"
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
      </div>
    </div>
  </main>
</template>
