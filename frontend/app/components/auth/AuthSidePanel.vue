<script setup lang="ts">
defineProps<{
  mode: 'login' | 'signup'
}>()

const image = '/images/1.jpg'
</script>

<template>
  <section
    class="relative hidden min-h-screen overflow-hidden lg:block lg:w-1/2"
  >
    <!-- Background image -->
    <img
      :src="image"
      alt="African fashion"
      class="absolute inset-0 h-full w-full object-cover"
    />

    <!-- Dark overlay -->
    <div
      class="absolute inset-0 bg-[#211f1d]/55"
    />

    <!-- Content -->
    <div class="relative z-10 flex min-h-screen items-center justify-center p-12">
      <div
        class="w-full max-w-lg border border-white/15 bg-[#211f1d]/80 px-10 py-12 text-center backdrop-blur-[2px] xl:px-16 xl:py-14"
      >
        <p
          class="text-xs font-medium uppercase tracking-[0.3em] text-[#d8c3a7]"
        >
          {{ mode === 'login'
            ? 'WELCOME BACK'
            : 'WELCOME TO AFRIMART' }}
        </p>

        <h2
          class="mt-7 font-serif text-5xl leading-[0.95] text-[#f5f1e9] xl:text-6xl"
        >
          {{ mode === 'login'
            ? 'Your style.'
            : 'Your style.' }}
          <br />

          <em>
            {{ mode === 'login'
              ? 'Your marketplace.'
              : 'Starts here.' }}
          </em>
        </h2>

        <p
          class="mx-auto mt-8 max-w-md text-base leading-7 text-[#e2d8c8]"
        >
          {{ mode === 'login'
            ? 'Sign in to access your orders, wishlist, and discover your next favorite pieces.'
            : 'Create an account and discover clothing from independent sellers across Africa.' }}
        </p>

        <!-- Switch action -->
        <div class="mt-10">
          <NuxtLink
            v-if="mode === 'login'"
            to="/auth/signup"
            class="inline-flex items-center justify-center border border-[#f5f1e9] px-10 py-4 text-sm font-medium tracking-[0.16em] text-[#f5f1e9] transition duration-300 hover:bg-[#f5f1e9] hover:text-[#211f1d]"
          >
            CREATE ACCOUNT
          </NuxtLink>

          <NuxtLink
            v-else
            to="/auth/login"
            class="inline-flex items-center justify-center border border-[#f5f1e9] px-10 py-4 text-sm font-medium tracking-[0.16em] text-[#f5f1e9] transition duration-300 hover:bg-[#f5f1e9] hover:text-[#211f1d]"
          >
            SIGN IN
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Bottom quote -->
    <p
      class="absolute bottom-8 left-0 right-0 z-10 text-center font-serif text-lg italic text-[#e2d8c8]/80"
    >
      "Your style. Your marketplace."
    </p>
  </section>
</template>