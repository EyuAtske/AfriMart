package observability

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "afrimart-http")
}