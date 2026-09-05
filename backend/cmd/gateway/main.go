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
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	shutdownObservability := observability.SetupObservability(ctx)
	defer shutdownObservability()
	slog.Info("starting AfriMart backend")
	apicfg := config.SetupAPIConfig(ctx)
	authHandler := &auth.AuthHandler{
		Config: apicfg,
		Queries: apicfg.Queries,
	}
	shopHandler := &shop.ShopHandler{
		Config:  apicfg,
		Queries: apicfg.Queries,
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
	servermux.Handle("PUT /api/auth/password", protected(http.HandlerFunc(authHandler.HandleUpdatePassword)))
	servermux.Handle("PUT /api/auth/username", protected(http.HandlerFunc(authHandler.HandleUpdateUsername)))
	servermux.Handle("GET /api/user/profile", protected(http.HandlerFunc(authHandler.HandleGetProfile)))
	servermux.HandleFunc("POST /api/refresh", authHandler.HandleRefresh)
	servermux.Handle("POST /api/shops", protected(http.HandlerFunc(shopHandler.HandleCreateShop)))
	servermux.Handle("GET /api/shops/me", protected(http.HandlerFunc(shopHandler.HandleGetMyShop)))
	servermux.Handle("PATCH /api/shops/{shopID}/name", protected(http.HandlerFunc(shopHandler.HandleUpdateShopName)))
	servermux.Handle("PATCH /api/shops/{shopID}/description", protected(http.HandlerFunc(shopHandler.HandleUpdateShopDescription)))
	servermux.Handle("PATCH /api/shops/{shopID}/deactivate", protected(http.HandlerFunc(shopHandler.HandleDeactivateShop)))
	servermux.Handle("PATCH /api/shops/{shopID}/activate", protected(http.HandlerFunc(shopHandler.HandleActivateShop)))
	// servermux.HandleFunc("GET /api/products", handlers.HandelProducts)
	// servermux.HandleFunc("GET /api/products/{id}", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/products", handlers.HandelProducts)
	// servermux.HandleFunc("PUT /api/products/{id}", handlers.HandelProducts)
	// servermux.HandleFunc("DELETE /api/products/{id}", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/products/{id}/images", handlers.HandelProducts)
	// servermux.HandleFunc("GET /api/cart", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/cart/items", handlers.HandelProducts)
	// servermux.HandleFunc("PUT /api/cart/items/{id}", handlers.HandelProducts)
	// servermux.HandleFunc("DELETE /api/cart/items/{id}", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/checkout", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/payments", handlers.HandelProducts)
	// servermux.HandleFunc("GET /api/payments/{id}", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/payments/{id}/verify", handlers.HandelProducts)
	// servermux.HandleFunc("GET /api/orders", handlers.HandelProducts)
	// servermux.HandleFunc("GET /api/orders/{id}", handlers.HandelProducts)
	// servermux.HandleFunc("GET /api/seller/orders", handlers.HandelProducts)
	// servermux.HandleFunc("PATCH /api/orders/{id}/status", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/orders/{id}/cancel", handlers.HandelProducts)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
