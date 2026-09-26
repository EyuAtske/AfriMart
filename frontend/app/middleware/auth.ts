export default defineNuxtRouteMiddleware(async () => {
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
