<script setup lang="ts">
import App from '~/app.vue';

</script>

<template>
  <div class="min-h-screen bg-[#f5f1e9]">
   <LayoutAppHeader />

    <main>
      <slot />
    </main>
   <LayoutAppFooter />
   <UiAppToast />
  </div>
</template>