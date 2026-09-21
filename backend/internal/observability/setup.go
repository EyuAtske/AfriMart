package observability

import (
	"context"
	"log/slog"
)

func SetupObservability(ctx context.Context) func() {
	logger := NewLogger()
	slog.SetDefault(logger)

	shutdownTracer, err := InitTracer(ctx)
	if err != nil {
		slog.Error(
			"warning: failed to initialize tracing and metrics",
			"error", err,
		)
		return func() {}
	}

	return func() {
		if err := shutdownTracer(ctx); err != nil {
			slog.Error(
				"failed to shutdown tracer",
				"error", err,
			)
		}
	}
}
