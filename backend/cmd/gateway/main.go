package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/EyuAtske/AfriMart/backend/internal/observability"

	database "github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/EyuAtske/AfriMart/backend/internal/handlers"
	"github.com/XSAM/otelsql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	logger := observability.NewLogger()
	slog.SetDefault(logger)

	shutdownTracer, err := observability.InitTracer(ctx)
	if err != nil {
		slog.Error("warning: failed to initialize tracing and metrics",
			"error", err,
		)
	} else {
		defer func() {
			if err := shutdownTracer(ctx); err != nil {
				slog.Error(
					"failed to shutdown tracer",
					"error", err,
				)
			}
		}()
	}

	slog.Info("starting AfriMart backend")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		slog.Error("DB_URL must be set")
		os.Exit(1)
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		slog.Error("SECRET_KEY environment variable is required")
		os.Exit(1)
	}

	dbConn, err := otelsql.Open("postgres", dbURL)
	if err != nil {
		slog.Error(
			"Error opening database",
			"error", err,
		)
		os.Exit(1)
	}

	defer dbConn.Close()

	if _, err := otelsql.RegisterDBStatsMetrics(dbConn); err != nil {
		slog.Warn(
			"failed to register DB metrics",
			"error", err,
		)
	}

	if err := dbConn.Ping(); err != nil {
		slog.Error(
			"database connection failed",
			"error", err,
		)
		os.Exit(1)
	}
	slog.Info("database connected")
	dbQueries := database.New(dbConn)
	apicfg := handlers.ApiConfig{
		DB:      dbConn,
		Queries: dbQueries,
		Secret:  secretKey,
	}
	servermux := http.NewServeMux()
	tracedHandler := observability.TraceMiddleware(servermux)
	server := &http.Server{
		Handler: tracedHandler,
		Addr:    ":8080",
	}
	slog.Info(
		"server started",
		"address", server.Addr,
	)
	servermux.HandleFunc("GET /api/health", handlers.HandelHealth)
	servermux.HandleFunc("POST /api/auth/register", apicfg.HandleRegister)
	servermux.HandleFunc("POST /api/auth/login", apicfg.HandleLogIn)
	servermux.HandleFunc("POST /api/auth/logout", apicfg.HandleRevoke)
	servermux.HandleFunc("PUT /api/auth/password", apicfg.HandleUpdatePassword)
	servermux.HandleFunc("PUT /api/auth/username", apicfg.HandleUpdateUsername)
	servermux.HandleFunc("POST /api/refresh", apicfg.HandleRefresh)
	servermux.HandleFunc("POST /api/shops", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/shops", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/shops/{id}", handlers.HandelProducts)
	servermux.HandleFunc("PUT /api/shops/{id}", handlers.HandelProducts)
	servermux.HandleFunc("DELETE /api/shops/{id}", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/products", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/products/{id}", handlers.HandelProducts)
	servermux.HandleFunc("POST /api/products", handlers.HandelProducts)
	servermux.HandleFunc("PUT /api/products/{id}", handlers.HandelProducts)
	servermux.HandleFunc("DELETE /api/products/{id}", handlers.HandelProducts)
	servermux.HandleFunc("POST /api/products/{id}/images", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/cart", handlers.HandelProducts)
	servermux.HandleFunc("POST /api/cart/items", handlers.HandelProducts)
	servermux.HandleFunc("PUT /api/cart/items/{id}", handlers.HandelProducts)
	servermux.HandleFunc("DELETE /api/cart/items/{id}", handlers.HandelProducts)
	servermux.HandleFunc("POST /api/checkout", handlers.HandelProducts)
	servermux.HandleFunc("POST /api/payments", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/payments/{id}", handlers.HandelProducts)
	servermux.HandleFunc("POST /api/payments/{id}/verify", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/orders", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/orders/{id}", handlers.HandelProducts)
	servermux.HandleFunc("GET /api/seller/orders", handlers.HandelProducts)
	servermux.HandleFunc("PATCH /api/orders/{id}/status", handlers.HandelProducts)
	servermux.HandleFunc("POST /api/orders/{id}/cancel", handlers.HandelProducts)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
