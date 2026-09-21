export const useFlyToCart = () => {
  const isCartBouncing = useState('fly-to-cart-bouncing', () => false)

  const triggerCartBounce = () => {
    isCartBouncing.value = true
    setTimeout(() => {
      isCartBouncing.value = false
    }, 450)
  }

  const flyToCart = (
    source: HTMLElement | null,
    imageUrl?: string
  ) => {
    if (typeof window === 'undefined') return

    const dest = document.getElementById('header-cart-icon')
    if (!dest) {
      triggerCartBounce()
      return
    }

    // Respect reduced motion preferences
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      triggerCartBounce()
      return
    }

    // Determine starting bounds
    let startRect: DOMRect | null = null
    let imgSrc = imageUrl || ''

    if (source) {
      startRect = source.getBoundingClientRect()
      if (!imgSrc && source instanceof HTMLImageElement) {
        imgSrc = source.src
      } else if (!imgSrc) {
        const foundImg = source.querySelector('img')
        if (foundImg) {
          imgSrc = foundImg.src
          startRect = foundImg.getBoundingClientRect()
        }
      }
    }

    if (!startRect || startRect.width === 0 || startRect.height === 0) {
      triggerCartBounce()
      return
    }

    const destRect = dest.getBoundingClientRect()

    // Create the flying element
    const flyer = document.createElement(imgSrc ? 'img' : 'div') as HTMLImageElement
    if (imgSrc) {
      flyer.src = imgSrc
    }
    
    // Set initial flying clone styles
    flyer.style.position = 'fixed'
    flyer.style.zIndex = '9999'
    flyer.style.left = `${startRect.left}px`
    flyer.style.top = `${startRect.top}px`
    flyer.style.width = `${Math.min(startRect.width, 140)}px`
    flyer.style.height = `${Math.min(startRect.height, 140)}px`
    flyer.style.borderRadius = '12px'
    flyer.style.objectFit = 'cover'
    flyer.style.pointerEvents = 'none'
    flyer.style.boxShadow = '0 12px 30px rgba(33, 31, 29, 0.25)'
    flyer.style.border = '2px solid #806344'
    flyer.style.transition = 'all 0.65s cubic-bezier(0.2, 0.85, 0.3, 1)'
    flyer.style.transform = 'translate3d(0, 0, 0) scale(1)'

    document.body.appendChild(flyer)

    // Trigger flight animation on next tick
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        const targetX = destRect.left + (destRect.width / 2) - 14
        const targetY = destRect.top + (destRect.height / 2) - 14

        flyer.style.left = `${targetX}px`
        flyer.style.top = `${targetY}px`
        flyer.style.width = '28px'
        flyer.style.height = '28px'
        flyer.style.borderRadius = '9999px'
        flyer.style.opacity = '0.3'
        flyer.style.transform = 'translate3d(0, 0, 0) scale(0.6) rotate(20deg)'
      })
    })

    // Clean up and trigger cart icon bounce when it lands
    setTimeout(() => {
      if (flyer.parentNode) {
        flyer.parentNode.removeChild(flyer)
      }
      triggerCartBounce()
    }, 650)
  }

  return {
    isCartBouncing,
    triggerCartBounce,
    flyToCart
  }
}
