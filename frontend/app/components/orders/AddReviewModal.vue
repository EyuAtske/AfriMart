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
  <UiAppModal
    :is-open="isOpen && !!product"
    :title="product ? `Review ${product.name}` : 'Add Review'"
    kicker="Product Review"
    @close="emit('close')"
  >
    <template v-if="product" #header>
      <p class="mt-1 text-xs text-[#756a60]">
        From shop: <strong class="text-[#211f1d]">{{ product.shop }}</strong>
      </p>
    </template>

    <form v-if="product" class="space-y-5" @submit.prevent="submitReview">
      <!-- Product Preview Header -->
      <div class="flex items-center gap-4 rounded-lg border border-[#ded6cc] bg-[#f5f1e9] p-3">
        <img
          :src="product.image"
          :alt="product.name"
          class="h-14 w-12 rounded-md object-cover object-top"
        />
        <div>
          <p class="text-sm font-semibold text-[#211f1d]">{{ product.name }}</p>
          <p class="text-xs text-[#756a60]">Delivered item</p>
        </div>
      </div>

      <!-- Star Rating Selector -->
      <div class="space-y-2">
        <label class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
          Rating
        </label>
        <div class="flex items-center gap-2">
          <button
            v-for="star in 5"
            :key="star"
            type="button"
            :aria-label="`Rate ${star} out of 5 stars`"
            class="text-2xl transition hover:scale-110 focus:outline-none focus-visible:ring-2 focus-visible:ring-[#806344] focus-visible:ring-offset-1 rounded"
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
        <label for="review-comment" class="block text-xs font-medium uppercase tracking-[0.16em] text-[#4d4035]">
          Your Review & Comments
        </label>
        <textarea
          id="review-comment"
          v-model="comment"
          rows="4"
          placeholder="Share details about fit, material, comfort, or delivery experience..."
          class="w-full rounded-md border border-[#cfc4b5] bg-[#faf8f4] px-4 py-3 text-sm text-[#211f1d] outline-none transition-all placeholder:text-[#92877b] hover:border-[#9e8b77] focus:border-[#806344] focus:ring-2 focus:ring-[#806344]/15"
        />
      </div>

      <UiAppAlert v-if="errorMsg">
        {{ errorMsg }}
      </UiAppAlert>

      <div class="flex justify-end gap-3 pt-2">
        <UiAppButton
          variant="ghost"
          @click="emit('close')"
        >
          Cancel
        </UiAppButton>

        <UiAppButton
          type="submit"
          variant="primary"
        >
          Submit Review
        </UiAppButton>
      </div>
    </form>
  </UiAppModal>
</template>
