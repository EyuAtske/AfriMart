<script setup lang="ts">
defineProps<{ categories: any[] }>()
const emit = defineEmits(['close'])
</script>

<template>
  <div class="border-t border-[#ded8ce] bg-[#f5f1e9] lg:hidden">
    <div class="px-6 py-4">
      <NuxtLink to="/" class="block border-b border-[#ded8ce] py-4 text-sm text-[#302d29]" @click="emit('close')">
        NEW IN
      </NuxtLink>

      <div v-for="category in categories" :key="category.name" class="border-b border-[#ded8ce] py-4">
        <NuxtLink :to="category.to" class="block text-sm text-[#302d29]" @click="emit('close')">
          {{ category.name }}
        </NuxtLink>
        <div class="mt-3 grid grid-cols-3 gap-4">
          <div v-for="section in category.sections" :key="section.title">
            <p class="mb-2 text-[10px] text-[#806344]">{{ section.title }}</p>
            <div class="flex flex-col gap-1.5">
              <NuxtLink
                v-for="item in section.items"
                :key="item.name"
                :to="item.to"
                class="text-xs text-[#5d574f]"
                @click="emit('close')"
              >
                {{ item.name }}
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>

      <NuxtLink to="/account" class="block border-b border-[#ded8ce] py-4 text-sm text-[#302d29]" @click="emit('close')">
        ACCOUNT
      </NuxtLink>
      <NuxtLink to="/cart" class="block py-4 text-sm text-[#302d29]" @click="emit('close')">
        CART
      </NuxtLink>
    </div>
  </div>
</template>
