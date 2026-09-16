<script setup lang="ts">
import AccountSidebar from '~/components/account/AccountSidebar.vue'

definePageMeta({
  middleware: 'auth'
})

const { user, updateUsername, updatePassword, fetchProfile } = useAuth()
const { showToast } = useToast()

const isEditingUsername = ref(false)
const isChangingPassword = ref(false)
const isSaving = ref(false)
const usernameError = ref('')
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

onMounted(async () => {
  try {
    const profile = await fetchProfile()
    if (profile) {
      if (profile.username) profileUser.username = profile.username
      if (profile.email) profileUser.email = profile.email
    }
  } catch {
    // Keep local session info if remote fetch is not reachable
  }
})

const usernameForm = reactive({
  username: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const startEditingUsername = () => {
  usernameError.value = ''
  usernameForm.username = profileUser.username
  isEditingUsername.value = true
}

const cancelEditingUsername = () => {
  usernameError.value = ''
  usernameForm.username = ''
  isEditingUsername.value = false
}

const saveUsername = async () => {
  usernameError.value = ''
  const username = usernameForm.username.trim()

  if (!username) {
    usernameError.value = 'Username is required.'
    return
  }
  if (username.length < 3) {
    usernameError.value = 'Username must be at least 3 characters long.'
    return
  }
  if (username.length > 50) {
    usernameError.value = 'Username must not exceed 50 characters.'
    return
  }

  isSaving.value = true

  try {
    const updated = await updateUsername(username)
    profileUser.username = updated.username
    usernameForm.username = ''
    isEditingUsername.value = false
    showToast('Username updated successfully', 'success')
  } catch (err: any) {
    usernameError.value = err?.message || 'Failed to update username'
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
    !passwordForm.newPassword.trim() ||
    !passwordForm.confirmPassword.trim()
  ) {
    passwordError.value = 'Please fill in all password fields.'
    return
  }

  if (passwordForm.newPassword.length < 8) {
    passwordError.value = 'Password must be at least 8 characters long.'
    return
  }

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    passwordError.value = 'Passwords do not match.'
    return
  }

  isSaving.value = true

  try {
    await updatePassword(passwordForm.newPassword)
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    isChangingPassword.value = false
    showToast('Password updated successfully', 'success')
  } catch (err: any) {
    passwordError.value = err?.message || 'Failed to update password'
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#f5f1e9] px-4 py-20 sm:px-6 lg:px-12">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 lg:flex-row">
      <AccountSidebar active="profile" />

      <div class="min-w-0 flex-1 space-y-6">
        <div>
          <h1 class="font-serif text-4xl tracking-[-0.025em] text-[#211f1d]">
            My Profile
          </h1>

          <p class="mt-2 text-base text-[#756a60]">
            Update your username and keep your password secure.
          </p>
        </div>

        <UiAppCard padding="large">
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
        </UiAppCard>

        <UiAppCard padding="large">
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

            <UiAppAlert v-if="usernameError">
              {{ usernameError }}
            </UiAppAlert>

            <div class="flex flex-wrap gap-3">
              <UiAppButton
                type="submit"
                variant="secondary"
                :disabled="isSaving"
              >
                {{ isSaving ? 'Saving...' : 'Save changes' }}
              </UiAppButton>

              <UiAppButton
                variant="ghost"
                :disabled="isSaving"
                @click="cancelEditingUsername"
              >
                Cancel
              </UiAppButton>
            </div>
          </form>
        </UiAppCard>

        <UiAppCard padding="large">
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

            <UiAppAlert v-if="passwordError">
              {{ passwordError }}
            </UiAppAlert>

            <div class="flex flex-wrap gap-3">
              <UiAppButton
                type="submit"
                variant="secondary"
                :disabled="isSaving"
              >
                {{ isSaving ? 'Updating...' : 'Save changes' }}
              </UiAppButton>

              <UiAppButton
                variant="ghost"
                :disabled="isSaving"
                @click="cancelChangingPassword"
              >
                Cancel
              </UiAppButton>
            </div>
          </form>
        </UiAppCard>
      </div>
    </div>
  </main>
</template>


