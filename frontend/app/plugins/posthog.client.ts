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
    })

    return {
        provide: {
            posthog,
        },
    }
})