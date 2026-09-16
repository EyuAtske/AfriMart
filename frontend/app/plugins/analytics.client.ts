export default defineNuxtPlugin(() => {
  const router = useRouter()
  const { gtag } = useGtag()

  let startedAt = Date.now()
  let accumulatedSeconds = 0
  let isActive = !document.hidden

  const getPageData = () => ({
    page_path: window.location.pathname + window.location.search,
    page_title: document.title
  })

  const startTimer = () => {
    startedAt = Date.now()
    isActive = true
  }

  const pauseTimer = () => {
    if (!isActive) return

    accumulatedSeconds += Math.floor((Date.now() - startedAt) / 1000)
    isActive = false
  }

  const getActiveSeconds = () => {
    if (!isActive) {
      return accumulatedSeconds
    }

    return accumulatedSeconds + Math.floor((Date.now() - startedAt) / 1000)
  }

  const sendPageTime = () => {
    const durationSeconds = getActiveSeconds()

    // Don't send extremely short visits.
    if (durationSeconds < 1) return

    const page = getPageData()

    gtag('event', 'page_active_time', {
      page_path: page.page_path,
      page_title: page.page_title,
      duration_seconds: durationSeconds
    })
  }

  const resetTimer = () => {
    accumulatedSeconds = 0
    startTimer()
  }

  const handleVisibilityChange = () => {
    if (document.hidden) {
      pauseTimer()
    } else {
      startTimer()
    }
  }

  document.addEventListener('visibilitychange', handleVisibilityChange)

  router.afterEach(() => {
    sendPageTime()
    resetTimer()
  })

  window.addEventListener('beforeunload', sendPageTime)

  return {
    provide: {
      analytics: {
        sendPageTime
      }
    }
  }
})