<script setup lang="ts">
defineProps<{ categories: any[] }>()
</script>

<template>
  <div class="hidden items-center gap-5 lg:flex xl:gap-7">
    <!-- NEW IN -->
    <NuxtLink
      to="/"
      class="py-5 text-xs tracking-[0.08em] text-[#302d29] hover:opacity-60 xl:text-sm"
    >
      NEW IN
    </NuxtLink>

    <!-- Categories -->
    <div v-for="category in categories" :key="category.name" class="group static">
      <NuxtLink
        :to="category.to"
        class="block py-5 text-xs tracking-[0.08em] text-[#302d29] hover:opacity-60 xl:text-sm"
      >
        {{ category.name }}
      </NuxtLink>

      <!-- Mega Menu -->
      <div
        class="invisible absolute left-0 right-0 top-full border-t border-[#ded8ce] bg-[#f5f1e9] opacity-0 shadow-[0_12px_30px_rgba(48,45,41,0.08)] transition-all duration-200 group-hover:visible group-hover:opacity-100"
      >
        <div class="mx-auto max-w-6xl px-12 py-8">
          <div class="mb-7 flex items-center justify-between border-b border-[#ded8ce] pb-5">
            <div>
              <p class="text-xs tracking-[0.18em] text-[#806344]">SHOP</p>
              <h3 class="mt-1 font-serif text-2xl text-[#211f1d]">{{ category.name }}</h3>
            </div>
            <NuxtLink :to="category.to" class="flex items-center gap-2 text-xs text-[#806344]">
              VIEW ALL →
            </NuxtLink>
          </div>

          <!-- Sections -->
          <div class="grid grid-cols-3 gap-12">
            <div v-for="section in category.sections" :key="section.title">
              <h4 class="mb-4 text-xs font-medium tracking-[0.14em] text-[#806344]">
                {{ section.title }}
              </h4>
              <div class="flex flex-col gap-2.5">
                <NuxtLink
                  v-for="item in section.items"
                  :key="item.name"
                  :to="item.to"
                  class="text-sm text-[#302d29] hover:text-[#806344]"
                >
                  {{ item.name }}
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
