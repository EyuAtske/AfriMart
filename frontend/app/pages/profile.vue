<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'

definePageMeta({
  middleware: 'auth'
})

const { user } = useAuth()

const isEditingUsername = ref(false)
const isChangingPassword = ref(false)
const isSaving = ref(false)
const passwordError = ref('')

const profileUser = reactive({
  username: '',
  email: ''
})

watch(
  user,
  (currentUser) => {
    profileUser.username = currentUser.username || currentUser.name || 'User'
    profileUser.email = currentUser.email || 'user@example.com'
  },
  { immediate: true }
)

const usernameForm = reactive({
  username: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const actionButtonClass = 'h-12 rounded-full border border-[#806344] px-6 text-sm font-medium uppercase tracking-[0.14em] text-[#5d4b37] transition-all duration-300 hover:bg-[#806344] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#806344] focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50'

const startEditingUsername = () => {
  usernameForm.username = profileUser.username
  isEditingUsername.value = true
}

const cancelEditingUsername = () => {
  usernameForm.username = ''
  isEditingUsername.value = false
}

const saveUsername = async () => {
  const username = usernameForm.username.trim()

  if (!username) return

  isSaving.value = true

  try {
    profileUser.username = username
    user.value.username = username
    user.value.name = username
    usernameForm.username = ''
    isEditingUsername.value = false
  } finally {
    isSaving.value = false
  }
}

const startChangingPassword = () => {
  passwordError.value = ''
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  isChangingPassword.value = true
}

const cancelChangingPassword = () => {
  passwordError.value = ''
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  isChangingPassword.value = false
}

const changePassword = async () => {
  passwordError.value = ''

  if (
    !passwordForm.currentPassword.trim() ||
    !passwordForm.newPassword.trim() ||
    !passwordForm.confirmPassword.trim()
  ) {
    passwordError.value = 'Please fill in all password fields.'
    return
  }

  if (passwordForm.newPassword.length < 6) {
    passwordError.value = 'Password must be at least 6 characters long.'
    return
  }

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    passwordError.value = 'Passwords do not match.'
    return
  }

  isSaving.value = true

  try {
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    isChangingPassword.value = false
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-8">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 lg:flex-row">
      <AccountSidebar active="profile" />

      <div class="min-w-0 flex-1">
        <div class="mb-8">
          <h1 class="font-serif text-4xl tracking-[-0.025em] text-[#211f1d]">
            My Profile
          </h1>

          <p class="mt-2 text-base text-[#756a60]">
            Update your username and keep your password secure.
          </p>
        </div>

        <section class="rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
          <div class="flex items-center gap-5">
            <div class="flex h-20 w-20 shrink-0 items-center justify-center rounded-full bg-[#211f1d] text-2xl font-medium text-white">
              {{ profileUser.username ? profileUser.username.charAt(0).toUpperCase() : 'U' }}
            </div>

            <div class="min-w-0">
              <h2 class="truncate text-2xl font-medium text-[#211f1d]">
                {{ profileUser.username }}
              </h2>

              <p class="mt-1 truncate text-base text-[#756a60]">
                {{ profileUser.email }}
              </p>
            </div>
          </div>
        </section>

        <section class="mt-6 rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
          <div class="flex items-start justify-between gap-4">
            <div>
              <h2 class="text-xl font-medium text-[#211f1d]">
                Username
              </h2>

              <p class="mt-1 text-sm text-[#756a60]">
                This is the name shown on your account.
              </p>
            </div>

            <button
              v-if="!isEditingUsername"
              type="button"
              class="text-sm font-medium text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline"
              @click="startEditingUsername"
            >
              Edit
            </button>
          </div>

          <p
            v-if="!isEditingUsername"
            class="mt-6 text-base text-[#211f1d]"
          >
            {{ profileUser.username }}
          </p>

          <form
            v-else
            class="mt-6 space-y-5"
            @submit.prevent="saveUsername"
          >
            <AuthInput
              v-model="usernameForm.username"
              label="Username"
              placeholder="Enter your username"
              autocomplete="username"
              name="profile-username"
            />

            <div class="flex flex-wrap gap-3">
              <button
                type="submit"
                :disabled="isSaving"
                :class="actionButtonClass"
              >
                {{ isSaving ? 'Saving...' : 'Save changes' }}
              </button>

              <button
                type="button"
                :disabled="isSaving"
                :class="actionButtonClass"
                @click="cancelEditingUsername"
              >
                Cancel
              </button>
            </div>
          </form>
        </section>

        <section class="mt-6 rounded-[12px] border border-[#d9d0c4] bg-[#faf8f4] p-6 shadow-[0_20px_70px_rgba(33,31,29,0.06)] sm:p-8">
          <div class="flex items-start justify-between gap-4">
            <div>
              <h2 class="text-xl font-medium text-[#211f1d]">
                Password
              </h2>

              <p class="mt-1 text-sm text-[#756a60]">
                Change your account password.
              </p>
            </div>

            <button
              v-if="!isChangingPassword"
              type="button"
              class="text-sm font-medium text-[#806344] underline-offset-4 transition hover:text-[#211f1d] hover:underline"
              @click="startChangingPassword"
            >
              Change
            </button>
          </div>

          <p
            v-if="!isChangingPassword"
            class="mt-6 text-base tracking-[0.2em] text-[#211f1d]"
          >
            ********
          </p>

          <form
            v-else
            class="mt-6 space-y-5"
            @submit.prevent="changePassword"
          >
            <AuthInput
              v-model="passwordForm.currentPassword"
              label="Current password"
              placeholder="Enter your current password"
              type="password"
              autocomplete="current-password"
              name="current-password"
            />

            <AuthInput
              v-model="passwordForm.newPassword"
              label="New password"
              placeholder="Enter your new password"
              type="password"
              autocomplete="new-password"
              name="new-password"
            />

            <AuthInput
              v-model="passwordForm.confirmPassword"
              label="Confirm new password"
              placeholder="Confirm your new password"
              type="password"
              autocomplete="new-password"
              name="confirm-new-password"
            />

            <p
              v-if="passwordError"
              class="text-sm font-medium text-red-600"
            >
              {{ passwordError }}
            </p>

            <div class="flex flex-wrap gap-3">
              <button
                type="submit"
                :disabled="isSaving"
                :class="actionButtonClass"
              >
                {{ isSaving ? 'Updating...' : 'Save changes' }}
              </button>

              <button
                type="button"
                :disabled="isSaving"
                :class="actionButtonClass"
                @click="cancelChangingPassword"
              >
                Cancel
              </button>
            </div>
          </form>
        </section>
      </div>
    </div>
  </main>
</template>


