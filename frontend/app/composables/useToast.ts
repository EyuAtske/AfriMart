export type ToastType = 'success' | 'info' | 'warning' | 'error'

export type ToastMessage = {
  id: number
  message: string
  type: ToastType
}

export const useToast = () => {
  const toasts = useState<ToastMessage[]>('app-toasts', () => [])

  const showToast = (message: string, type: ToastType = 'success', duration = 3500) => {
    const id = Date.now() + Math.random()
    toasts.value.push({ id, message, type })

    setTimeout(() => {
      removeToast(id)
    }, duration)
  }

  const removeToast = (id: number) => {
    toasts.value = toasts.value.filter(toast => toast.id !== id)
  }

  return {
    toasts,
    showToast,
    removeToast
  }
}
