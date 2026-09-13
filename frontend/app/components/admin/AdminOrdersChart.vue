<script setup lang="ts">
import { computed, ref } from 'vue'
import type { MonthlyMetric } from '~/composables/useAdminDashboard'

const props = defineProps<{
  metrics: MonthlyMetric[]
}>()

const hoveredIndex = ref<number | null>(null)

// Dimensions for SVG viewport
const width = 560
const height = 240
const padLeft = 45
const padRight = 30
const padTop = 25
const padBottom = 35

const chartWidth = width - padLeft - padRight
const chartHeight = height - padTop - padBottom

const maxY = 240
const yTicks = [240, 180, 120, 60, 0]

const getY = (val: number) => {
  const clamped = Math.max(0, Math.min(val, maxY))
  return padTop + (1 - clamped / maxY) * chartHeight
}

const getX = (index: number) => {
  if (props.metrics.length <= 1) return padLeft + chartWidth / 2
  return padLeft + (index / (props.metrics.length - 1)) * chartWidth
}

const barWidth = 18

const bars = computed(() => {
  return props.metrics.map((m, i) => {
    const baseline = padTop + chartHeight
    const barTop = getY(m.orders)
    const barHeight = Math.max(2, baseline - barTop)

    return {
      x: getX(i) - barWidth / 2,
      y: barTop,
      width: barWidth,
      height: barHeight,
      data: m
    }
  })
})

const hoveredMetric = computed(() => {
  if (hoveredIndex.value === null) return null
  return props.metrics[hoveredIndex.value] ?? null
})
</script>

<template>
  <div class="rounded-lg border border-[#ede8e1] bg-white p-5 sm:p-6">
    <div class="mb-4 flex items-center justify-between">
      <h3 class="text-[10px] sm:text-[11px] font-bold tracking-[0.12em] text-[#85776a] uppercase">
        ORDERS OVERVIEW — LAST 6 MONTHS
      </h3>
      <span v-if="hoveredMetric" class="font-serif text-sm font-medium text-[#795736]">
        {{ hoveredMetric.fullName }}: {{ hoveredMetric.orders }} orders
      </span>
    </div>

    <div class="relative w-full">
      <svg
        :viewBox="`0 0 ${width} ${height}`"
        class="h-auto w-full overflow-visible font-sans text-xs select-none"
      >
        <!-- Horizontal Dotted Grid Lines & Y-Axis Labels -->
        <g v-for="tick in yTicks" :key="tick">
          <text
            :x="padLeft - 12"
            :y="getY(tick) + 3.5"
            text-anchor="end"
            class="fill-[#8a7a6c] text-[10px]"
          >
            {{ tick }}
          </text>

          <line
            :x1="padLeft"
            :y1="getY(tick)"
            :x2="width - padRight"
            :y2="getY(tick)"
            stroke="#e4dbce"
            stroke-width="1"
            stroke-dasharray="2 3"
          />
        </g>

        <!-- Vertical Dotted Grid Lines & X-Axis Labels -->
        <g v-for="(metric, i) in metrics" :key="metric.month">
          <line
            :x1="getX(i)"
            :y1="padTop"
            :x2="getX(i)"
            :y2="padTop + chartHeight"
            stroke="#e4dbce"
            stroke-width="1"
            stroke-dasharray="2 3"
          />

          <text
            :x="getX(i)"
            :y="padTop + chartHeight + 18"
            text-anchor="middle"
            class="fill-[#61574d] text-[11px]"
          >
            {{ metric.month }}
          </text>
        </g>

        <!-- Vertical Bars -->
        <g v-for="(bar, i) in bars" :key="i">
          <rect
            :x="bar.x"
            :y="bar.y"
            :width="bar.width"
            :height="bar.height"
            rx="2"
            ry="2"
            :fill="hoveredIndex === i ? '#624326' : '#7a5836'"
            class="transition-colors duration-150"
          />

          <!-- Invisible wider hit area for hover -->
          <rect
            :x="bar.x - 10"
            :y="padTop"
            :width="bar.width + 20"
            :height="chartHeight"
            fill="transparent"
            class="cursor-pointer"
            @mouseenter="hoveredIndex = i"
            @mouseleave="hoveredIndex = null"
          />
        </g>
      </svg>
    </div>
  </div>
</template>
