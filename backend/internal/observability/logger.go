package observability

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type TraceHandler struct {
	slog.Handler
}

func (h *TraceHandler) Handle(
	ctx context.Context,
	r slog.Record,
) error {
	span := trace.SpanFromContext(ctx)

	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, r)
}

func NewLogger() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	handler := &TraceHandler{
		Handler: jsonHandler,
	}

	return slog.New(handler)
}