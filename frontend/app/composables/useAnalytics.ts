import type { PostHog } from 'posthog-js'

export const useAnalytics = () => {
    const { $posthog } = useNuxtApp()
    const posthog = $posthog as PostHog | undefined

    const track = (
        event: string,
        properties?: Record<string, unknown>
    ) => {
        if (!import.meta.client || !posthog) return

        posthog.capture(event, properties)
    }

    const isFeatureEnabled = (flag: string) => {
        if (!import.meta.client || !posthog) return false

        return posthog.isFeatureEnabled(flag)
    }

    const getFeatureFlag = (flag: string) => {
        if (!import.meta.client || !posthog) return undefined

        return posthog.getFeatureFlag(flag)
    }

    return {
        track,
        isFeatureEnabled,
        getFeatureFlag,
    }
}