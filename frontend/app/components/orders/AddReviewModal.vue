<script setup lang="ts">
import type { Product } from '~/types/product'
import { useMarketplace } from '~/composables/useMarketplace'
import { useToast } from '~/composables/useToast'

const props = defineProps<{
  isOpen: boolean
  product: Product | null
  orderId?: number
}>()

const emit = defineEmits<{
  close: []
  submitted: []
}>()

const { submitProductReview } = useMarketplace()
const { showToast } = useToast()

const rating = ref(5)
const comment = ref('')
const errorMsg = ref('')

const submitReview = () => {
  errorMsg.value = ''

  if (!comment.value.trim()) {
    errorMsg.value = 'Please write a brief comment for your review.'
    return
  }

  if (!props.product) return

  submitProductReview(props.product.id, rating.value, comment.value.trim(), undefined, props.orderId)
  showToast('Review submitted! It will appear after admin approval.')

  comment.value = ''
  rating.value = 5
  emit('submitted')
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isOpen && product"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm transition-opacity"
      @click.self="emit('close')"
    >
      <div class="relative w-full max-w-lg rounded-[16px] border border-[#d9d0c4] bg-[#faf8f4] p-6 sm:p-8 shadow-[0_25px_80px_rgba(0,0,0,0.25)]">
        <!-- Close Button -->
        <button
          type="button"
          class="absolute top-5 right-5 flex h-9 w-9 items-center justify-center rounded-full bg-[#eee8df] text-[#665c53] transition hover:bg-[#211f1d] hover:text-white"
          @click="emit('close')"
        >
          <span class="sr-only">Close modal</span>
          ✕
        </button>

        <div class="border-b border-[#ded6cc] pb-4">
          <p class="text-xs font-medium uppercase tracking-[0.18em] text-[#806344]">
            Product Review
          </p>
          <h2 class="mt-1 font-serif text-2xl text-[#211f1d]">
            Review {{ product.name }}
          </h2>
          <p class="mt-1 text-xs text-[#756a60]">
            From shop: <strong class="text-[#211f1d]">{{ product.shop }}</strong>
          </p>
        </div>

        <form class="mt-6 space-y-5" @submit.prevent="submitReview">
          <!-- Product Preview Header -->
          <div class="flex items-center gap-4 rounded-[10px] border border-[#ded6cc] bg-[#f5f1e9] p-3">
            <img
              :src="product.image"
              :alt="product.name"
              class="h-14 w-12 rounded object-cover object-top"
            />
            <div>
              <p class="text-sm font-semibold text-[#211f1d]">{{ product.name }}</p>
              <p class="text-xs text-[#756a60]">Delivered item</p>
            </div>
          </div>

          <!-- Star Rating Selector -->
          <div class="space-y-2">
            <label class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
              Rating
            </label>
            <div class="flex items-center gap-2">
              <button
                v-for="star in 5"
                :key="star"
                type="button"
                class="text-2xl transition hover:scale-110 focus:outline-none"
                :class="star <= rating ? 'text-amber-500' : 'text-gray-300'"
                @click="rating = star"
              >
                ★
              </button>
              <span class="ml-2 text-sm font-medium text-[#211f1d]">
                {{ rating }} / 5 Stars
              </span>
            </div>
          </div>

          <!-- Review Comment Input -->
          <div class="space-y-2">
            <label for="review-comment" class="block text-[11px] font-medium uppercase tracking-[0.16em] text-[#4d4035]">
              Your Review & Comments
            </label>
            <textarea
              id="review-comment"
              v-model="comment"
              rows="4"
              placeholder="Share details about fit, material, comfort, or delivery experience..."
              class="w-full rounded-[5px] border border-[#cfc4b5] bg-white px-4 py-3 text-sm text-[#211f1d] outline-none transition placeholder:text-[#92877b] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
            />
          </div>

          <p v-if="errorMsg" class="text-xs font-medium text-red-600">
            {{ errorMsg }}
          </p>

          <div class="flex justify-end gap-3 pt-2">
            <button
              type="button"
              class="h-11 rounded-full border border-[#cfc4b5] px-5 text-xs font-medium uppercase tracking-[0.14em] text-[#756a60] transition hover:bg-[#eee8df]"
              @click="emit('close')"
            >
              Cancel
            </button>

            <button
              type="submit"
              class="h-11 rounded-full bg-[#211f1d] px-6 text-xs font-medium uppercase tracking-[0.14em] text-white transition hover:bg-[#3b3733]"
            >
              Submit Review
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>
