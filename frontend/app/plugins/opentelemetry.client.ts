import { metrics, propagation } from '@opentelemetry/api'
import { W3CTraceContextPropagator } from '@opentelemetry/core'
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'
import { OTLPMetricExporter } from '@opentelemetry/exporter-metrics-otlp-http'
import { OTLPLogExporter } from '@opentelemetry/exporter-logs-otlp-http'
import { ZoneContextManager } from '@opentelemetry/context-zone'
import { resourceFromAttributes } from '@opentelemetry/resources'
import { registerInstrumentations } from '@opentelemetry/instrumentation'
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch'
import {
  BatchSpanProcessor,
  WebTracerProvider,
} from '@opentelemetry/sdk-trace-web'
import {
  MeterProvider,
  PeriodicExportingMetricReader,
} from '@opentelemetry/sdk-metrics'
import {
  LoggerProvider,
  BatchLogRecordProcessor,
} from '@opentelemetry/sdk-logs'
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions'
import { onCLS, onINP, onLCP } from 'web-vitals'

export default defineNuxtPlugin((nuxtApp) => {
  const resource = resourceFromAttributes({
    [ATTR_SERVICE_NAME]: 'afrimart-frontend',
  })

  // ---------------------------------------------------------------------------
  // Tracing
  // ---------------------------------------------------------------------------

  const traceExporter = new OTLPTraceExporter({
    url: 'http://localhost:4318/v1/traces',
  })

  const tracerProvider = new WebTracerProvider({
    resource,
    spanProcessors: [
      new BatchSpanProcessor(traceExporter),
    ],
  })

  tracerProvider.register({
    contextManager: new ZoneContextManager(),
    propagator: new W3CTraceContextPropagator(),
  })

  propagation.setGlobalPropagator(
    new W3CTraceContextPropagator(),
  )

  const tracer = tracerProvider.getTracer('afrimart-frontend')

  // ---------------------------------------------------------------------------
  // Metrics
  // ---------------------------------------------------------------------------

  const metricExporter = new OTLPMetricExporter({
    url: 'http://localhost:4318/v1/metrics',
  })

  const metricReader = new PeriodicExportingMetricReader({
    exporter: metricExporter,
    exportIntervalMillis: 10000,
  })

  const meterProvider = new MeterProvider({
    resource,
    readers: [metricReader],
  })

  metrics.setGlobalMeterProvider(meterProvider)

  const meter = meterProvider.getMeter('afrimart-frontend')



  // ---------------------------------------------------------------------------
  // Frontend errors
  // ---------------------------------------------------------------------------

  const frontendErrorCounter = meter.createCounter(
    'frontend_errors',
    {
      description: 'Number of unhandled frontend JavaScript errors',
    },
  )

  const lcpHistogram = meter.createHistogram(
    'web_vital_lcp',
    {
      description: 'Largest Contentful Paint',
      unit: 'ms',
    },
  )

  const inpHistogram = meter.createHistogram(
    'web_vital_inp',
    {
      description: 'Interaction to Next Paint',
      unit: 'ms',
    },
  )

  const clsHistogram = meter.createHistogram(
    'web_vital_cls',
    {
      description: 'Cumulative Layout Shift',
    },
  )

  // ---------------------------------------------------------------------------
  // Logging
  // ---------------------------------------------------------------------------

  const logExporter = new OTLPLogExporter({
    url: 'http://localhost:4318/v1/logs',
  })

  const loggerProvider = new LoggerProvider({
    resource,
    processors: [
      new BatchLogRecordProcessor({
        exporter: logExporter,
      }),
    ],
  })

  const logger = loggerProvider.getLogger('afrimart-frontend')

  const emitCorrelatedLog = (
    span: ReturnType<typeof tracer.startSpan>,
    severityText: 'INFO' | 'ERROR',
    body: string,
    attributes: Record<string, string>,
  ) => {
    const spanContext = span.spanContext()

    const logData = {
      time: new Date().toISOString(),
      level: severityText,
      msg: body,
      ...(attributes.error
        ? { error: attributes.error }
        : {}),
      trace_id: spanContext.traceId,
      span_id: spanContext.spanId,
    }

    logger.emit({
      severityText,
      body: JSON.stringify(logData),
    })
  }

  // ---------------------------------------------------------------------------
  // Global frontend errors
  // ---------------------------------------------------------------------------

  window.addEventListener('error', (event) => {
    const error =
      event.error instanceof Error
        ? event.error
        : new Error(event.message)

    frontendErrorCounter.add(1, {
      route: window.location.pathname,
      error_type: error.name || 'Error',
    })

    tracer.startActiveSpan(
      'frontend.error',
      (span) => {
        span.recordException(error)

        emitCorrelatedLog(
          span,
          'ERROR',
          event.message || 'Unhandled frontend error',
          {
            route: window.location.pathname,
            error_type: error.name || 'Error',
            error: error.message || event.message || 'Unknown frontend error',
          },
        )

        span.end()
      },
    )
  })

  // ---------------------------------------------------------------------------
  // Unhandled promise rejections
  // ---------------------------------------------------------------------------

  window.addEventListener('unhandledrejection', (event) => {
    const error =
      event.reason instanceof Error
        ? event.reason
        : new Error(String(event.reason))

    frontendErrorCounter.add(1, {
      route: window.location.pathname,
      error_type: 'UnhandledPromiseRejection',
    })

    tracer.startActiveSpan(
      'frontend.unhandled_rejection',
      (span) => {
        span.recordException(error)

        emitCorrelatedLog(
          span,
          'ERROR',
          'Unhandled promise rejection',
          {
            route: window.location.pathname,
            error_type: 'UnhandledPromiseRejection',
            error: error.message,
          },
        )

        span.end()
      },
    )
  })

  // ---------------------------------------------------------------------------
  // Web Vitals
  // ---------------------------------------------------------------------------

  onLCP((metric) => {
    lcpHistogram.record(metric.value, {
      route: window.location.pathname,
    })
  })

  onINP((metric) => {
    inpHistogram.record(metric.value, {
      route: window.location.pathname,
    })
  })

  onCLS((metric) => {
    clsHistogram.record(metric.value, {
      route: window.location.pathname,
    })
  })

  // ---------------------------------------------------------------------------
  // Vue errors
  // ---------------------------------------------------------------------------

  nuxtApp.hook('vue:error', (err, _instance, info) => {
    const route =
      typeof window !== 'undefined'
        ? window.location.pathname
        : ''

    const errorObj =
      err instanceof Error
        ? err
        : new Error(String(err))

    frontendErrorCounter.add(1, {
      route,
      error_type: 'VueError',
    })

    tracer.startActiveSpan(
      'frontend.vue_error',
      (span) => {
        span.recordException(errorObj)

        emitCorrelatedLog(
          span,
          'ERROR',
          `Vue Error: ${info}`,
          {
            route,
            info: String(info),
            error: errorObj.message,
          },
        )

        span.end()
      },
    )
  })

  // ---------------------------------------------------------------------------
  // Fetch tracing / trace propagation
  // ---------------------------------------------------------------------------

  registerInstrumentations({
    instrumentations: [
      new FetchInstrumentation({
        propagateTraceHeaderCorsUrls: [
          /http:\/\/localhost:8080/,
          /http:\/\/localhost:3000/,
          /^\/api\//,
        ],
      }),
    ],
  })

  // ---------------------------------------------------------------------------
  // OTel initialization test
  // ---------------------------------------------------------------------------

  tracer.startActiveSpan(
    'frontend.log-trace-test',
    (span) => {
      emitCorrelatedLog(
        span,
        'INFO',
        'OpenTelemetry frontend logging initialized',
        {
          route: window.location.pathname,
        },
      )

      console.info(
        '[OpenTelemetry] Frontend tracing, metrics, and logging initialized',
      )

      span.end()
    },
  )
})