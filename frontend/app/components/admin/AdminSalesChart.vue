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
const padLeft = 55
const padRight = 30
const padTop = 25
const padBottom = 35

const chartWidth = width - padLeft - padRight
const chartHeight = height - padTop - padBottom

const maxY = 12000
const yTicks = [12000, 9000, 6000, 3000, 0]

const getY = (val: number) => {
  const clamped = Math.max(0, Math.min(val, maxY))
  return padTop + (1 - clamped / maxY) * chartHeight
}

const getX = (index: number) => {
  if (props.metrics.length <= 1) return padLeft + chartWidth / 2
  return padLeft + (index / (props.metrics.length - 1)) * chartWidth
}

const points = computed(() => {
  return props.metrics.map((m, i) => ({
    x: getX(i),
    y: getY(m.sales),
    data: m
  }))
})

const hoveredMetric = computed(() => {
  if (hoveredIndex.value === null) return null
  return props.metrics[hoveredIndex.value] ?? null
})

// Generate smooth cubic Bézier curve
const pathD = computed(() => {
  const pts = points.value
  if (!pts || pts.length === 0) return ''
  const start = pts[0]
  if (!start) return ''
  if (pts.length === 1) return `M ${start.x} ${start.y}`

  let d = `M ${start.x} ${start.y}`

  for (let i = 0; i < pts.length - 1; i++) {
    const p0 = pts[Math.max(0, i - 1)] ?? start
    const p1 = pts[i] ?? start
    const p2 = pts[i + 1] ?? p1
    const p3 = pts[Math.min(pts.length - 1, i + 2)] ?? p2

    const cp1x = p1.x + (p2.x - p0.x) / 6
    const cp1y = p1.y + (p2.y - p0.y) / 6

    const cp2x = p2.x - (p3.x - p1.x) / 6
    const cp2y = p2.y - (p3.y - p1.y) / 6

    d += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${p2.x} ${p2.y}`
  }

  return d
})

// Generate area path closed to bottom baseline
const areaD = computed(() => {
  const line = pathD.value
  const pts = points.value
  if (!line || !pts || pts.length === 0) return ''
  const first = pts[0]
  const last = pts[pts.length - 1]
  if (!first || !last) return ''
  const baseline = getY(0)

  return `${line} L ${last.x} ${baseline} L ${first.x} ${baseline} Z`
})
</script>

<template>
  <div class="rounded-lg border border-[#ede8e1] bg-white p-5 sm:p-6">
    <div class="mb-4 flex items-center justify-between">
      <h3 class="text-[10px] sm:text-[11px] font-bold tracking-[0.12em] text-[#85776a] uppercase">
        SALES OVERVIEW — LAST 6 MONTHS
      </h3>
      <span v-if="hoveredMetric" class="font-serif text-sm font-medium text-[#795736]">
        {{ hoveredMetric.fullName }}: ${{ hoveredMetric.sales.toLocaleString() }}
      </span>
    </div>

    <div class="relative w-full">
      <svg
        :viewBox="`0 0 ${width} ${height}`"
        class="h-auto w-full overflow-visible font-sans text-xs select-none"
      >
        <defs>
          <linearGradient id="salesGradient" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#7a5836" stop-opacity="0.18" />
            <stop offset="75%" stop-color="#7a5836" stop-opacity="0.03" />
            <stop offset="100%" stop-color="#7a5836" stop-opacity="0.0" />
          </linearGradient>
        </defs>

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

        <!-- Filled Gradient Area -->
        <path
          :d="areaD"
          fill="url(#salesGradient)"
        />

        <!-- Smooth Spline Stroke -->
        <path
          :d="pathD"
          fill="none"
          stroke="#7a5836"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />

        <!-- Hover Interactive Touch/Mouse Columns and Dots -->
        <g v-for="(pt, i) in points" :key="i">
          <!-- Invisible wide hover column -->
          <rect
            :x="pt.x - chartWidth / (metrics.length * 2)"
            :y="padTop"
            :width="chartWidth / metrics.length"
            :height="chartHeight"
            fill="transparent"
            class="cursor-pointer"
            @mouseenter="hoveredIndex = i"
            @mouseleave="hoveredIndex = null"
          />

          <!-- Point Dot (glow on hover) -->
          <circle
            v-if="hoveredIndex === i"
            :cx="pt.x"
            :cy="pt.y"
            r="5"
            fill="#7a5836"
            stroke="#ffffff"
            stroke-width="2.5"
            class="transition-all"
          />
        </g>
      </svg>
    </div>
  </div>
</template>
