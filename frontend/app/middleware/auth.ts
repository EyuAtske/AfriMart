export default defineNuxtRouteMiddleware(async () => {
  if (import.meta.server && useCookie<string | null>('afrimart_access_token').value) {
    return
  }

  const { isLoggedIn } = useAuth()
  const { authRepo } = useRepositories()

  if (!isLoggedIn.value) {
    try {
      const restoredUser = await authRepo.getCurrentSession()
      if (restoredUser) {
        return
      }
    } catch {
      // Fallback
    }

    return navigateTo('/account')
  }
})
