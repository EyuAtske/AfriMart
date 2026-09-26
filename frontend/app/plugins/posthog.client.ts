import posthog from 'posthog-js'

export default defineNuxtPlugin(() => {
    const config = useRuntimeConfig()

    if (!config.public.posthogKey) {
        console.warn('PostHog key is not configured')
        return
    }

    posthog.init(config.public.posthogKey, {
        api_host: 'https://us.i.posthog.com',
        capture_pageview: false,
        capture_pageleave: false,
        capture_exceptions: true,
    })

    if (import.meta.client) {
        const captureMemorySnapshot = () => {
            const memory = (performance as Performance & {
                memory?: {
                    usedJSHeapSize: number
                    totalJSHeapSize: number
                    jsHeapSizeLimit: number
                }
            }).memory

            if (!memory) return

            posthog.capture('frontend_memory_snapshot', {
                heap_used_mb: Math.round(memory.usedJSHeapSize / 1048576),
                heap_total_mb: Math.round(memory.totalJSHeapSize / 1048576),
                heap_limit_mb: Math.round(memory.jsHeapSizeLimit / 1048576),
            })
        }

        captureMemorySnapshot()
        window.setInterval(captureMemorySnapshot, 60000)
    }

    return {
        provide: {
            posthog,
        },
    }
})