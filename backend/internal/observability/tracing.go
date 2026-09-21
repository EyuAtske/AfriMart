package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracer(ctx context.Context) (func(context.Context) error, error) {
	endpoint := getOTLPEndpoint()

	res, err := createResource(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider, err := createTracerProvider(ctx, endpoint, res)
	if err != nil {
		return nil, err
	}

	meterProvider, err := createMeterProvider(ctx, endpoint, res)
	if err != nil {
		_ = tracerProvider.Shutdown(ctx)
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)

	startRuntimeMetrics(meterProvider)

	return createShutdown(tracerProvider, meterProvider), nil
}

func getOTLPEndpoint() string {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return "otel-collector:4317"
	}

	return endpoint
}

func createResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName("afrimart-backend"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return res, nil
}

func createTracerProvider(
	ctx context.Context,
	endpoint string,
	res *resource.Resource,
) (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	), nil
}

func createMeterProvider(
	ctx context.Context,
	endpoint string,
	res *resource.Resource,
) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	return metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(15*time.Second),
			),
		),
		metric.WithResource(res),
	), nil
}

func startRuntimeMetrics(meterProvider *metric.MeterProvider) {
	if err := runtime.Start(
		runtime.WithMeterProvider(meterProvider),
		runtime.WithMinimumReadMemStatsInterval(5*time.Second),
	); err != nil {
		slog.Warn(
			"failed to start runtime metrics",
			"error", err,
		)
	}
}

func createShutdown(
	tracerProvider *sdktrace.TracerProvider,
	meterProvider *metric.MeterProvider,
) func(context.Context) error {
	return func(ctx context.Context) error {
		var errs []error

		if err := tracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("tracer shutdown: %w", err))
		}

		if err := meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("meter shutdown: %w", err))
		}

		if len(errs) > 0 {
			return fmt.Errorf("observability shutdown errors: %v", errs)
		}

		return nil
	}
}