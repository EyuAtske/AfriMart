package observability

import (
	"log/slog"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming Request", 
			"method", r.Method,
			"path", r.URL.Path, 
			"traceparent", r.Header.Get("traceparent"),
		)
		otelhttp.NewMiddleware("afrimart-backend")(next).ServeHTTP(w, r)
	})
}
