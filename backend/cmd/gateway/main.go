package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/health"
	"github.com/EyuAtske/AfriMart/backend/internal/observability"
	"github.com/EyuAtske/AfriMart/backend/internal/shop"

	"github.com/EyuAtske/AfriMart/backend/config"
	database "github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/EyuAtske/AfriMart/backend/internal/handlers"
	"github.com/XSAM/otelsql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	shutdownObservability := observability.SetupObservability(ctx)
	defer shutdownObservability()

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
	apicfg := config.ApiConfig{
		DB:      dbConn,
		Queries: dbQueries,
		Secret:  secretKey,
	}
	authHandler := &auth.AuthHandler{
		Config: &apicfg,
	}
	shopHandler := &shop.ShopHandler{
		Config: &apicfg,
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
	protected := auth.AuthMiddleware(apicfg.Secret)
	servermux.HandleFunc("GET /api/health", health.HandelHealth)
	servermux.HandleFunc("POST /api/auth/register", authHandler.HandleRegister)
	servermux.HandleFunc("POST /api/auth/login", authHandler.HandleLogIn)
	servermux.HandleFunc("POST /api/auth/logout", authHandler.HandleRevoke)
	servermux.Handle("PUT /api/auth/password", protected(http.HandlerFunc(authHandler.HandleUpdatePassword)),)
	servermux.Handle("PUT /api/auth/username", protected(http.HandlerFunc(authHandler.HandleUpdateUsername)))
	servermux.Handle("GET /api/user/profile", protected(http.HandlerFunc(authHandler.HandleGetProfile)))
	servermux.HandleFunc("POST /api/refresh", authHandler.HandleRefresh)
	servermux.Handle("POST /api/shops",protected(http.HandlerFunc(shopHandler.HandleCreateShop)),)
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
