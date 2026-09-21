package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/cart"
	"github.com/EyuAtske/AfriMart/backend/internal/health"
	"github.com/EyuAtske/AfriMart/backend/internal/observability"
	"github.com/EyuAtske/AfriMart/backend/internal/order"
	"github.com/EyuAtske/AfriMart/backend/internal/product"
	"github.com/EyuAtske/AfriMart/backend/internal/shop"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	shutdownObservability := observability.SetupObservability(ctx)
	defer shutdownObservability()
	slog.Info("starting AfriMart backend")
	apicfg := config.SetupAPIConfig(ctx)
	authHandler := auth.NewAuthHandler(apicfg, apicfg.Queries, slog.Default())
	shopHandler := shop.NewShopHandler(apicfg, apicfg.Queries, slog.Default())
	productHandler := product.NewProductHandler(apicfg, apicfg.Queries, apicfg.ImageStorage, slog.Default())
	cartHandler := cart.NewCartHandler(apicfg.Queries, slog.Default())
	orderHandler := order.NewOrderHandler(apicfg, apicfg.Queries, slog.Default())
	servermux := http.NewServeMux()
	tracedHandler := observability.TraceMiddleware(servermux)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "traceparent", "tracestate", "baggage"},
		ExposedHeaders:   []string{"traceresponse"}, // Allows frontend to read the trace response
		AllowCredentials: true,
	})
	server := &http.Server{
		Handler: c.Handler(tracedHandler),
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
	servermux.Handle("POST /api/products", protected(http.HandlerFunc(productHandler.HandleCreateProduct)))
	servermux.Handle("PUT /api/products/{id}", protected(http.HandlerFunc(productHandler.HandleUpdateProduct)))
	servermux.Handle("DELETE /api/products/{id}", protected(http.HandlerFunc(productHandler.HandleDeleteProduct)))
	servermux.HandleFunc("GET /api/products/{id}", productHandler.HandleGetProduct)
	servermux.HandleFunc("GET /api/products", productHandler.HandleListProducts)
	servermux.Handle("GET /api/shops/{shop_id}/products", protected(http.HandlerFunc(productHandler.HandleListProductsByShop)))
	servermux.HandleFunc("GET /api/categories/{category_id}/products", productHandler.HandleListProductsByCategory)
	servermux.HandleFunc("GET /api/subcategories/{subcategory_id}/products", productHandler.HandleListProductsBySubcategory)
	servermux.Handle("GET /api/cart", protected(http.HandlerFunc(cartHandler.HandleGetCart)))
	servermux.Handle("POST /api/cart/items", protected(http.HandlerFunc(cartHandler.HandleAddCartItem)))
	servermux.Handle("PATCH /api/cart/items/{id}", protected(http.HandlerFunc(cartHandler.HandleUpdateCartItem)))
	servermux.Handle("DELETE /api/cart/items/{id}", protected(http.HandlerFunc(cartHandler.HandleDeleteCartItem)))
	servermux.Handle("DELETE /api/cart", protected(http.HandlerFunc(cartHandler.HandleClearCart)))
	servermux.Handle("POST /api/orders/checkout", protected(http.HandlerFunc(orderHandler.HandleCheckout)))
	servermux.Handle("GET /api/orders", protected(http.HandlerFunc(orderHandler.HandleListOrders)))
	servermux.Handle("GET /api/orders/{id}", protected(http.HandlerFunc(orderHandler.HandleGetOrder)))
	servermux.Handle("PATCH /api/orders/{id}/status", protected(http.HandlerFunc(orderHandler.HandleUpdateOrderStatus)))
	servermux.Handle("GET /api/orders/seller", protected(http.HandlerFunc(orderHandler.HandleListSellerOrders)))
	servermux.Handle("PATCH /api/products/{id}/images/{imageID}", protected(http.HandlerFunc(productHandler.HandleUpdateProductImage)))
	servermux.Handle("DELETE /api/products/{id}/images/{imageID}", protected(http.HandlerFunc(productHandler.HandleDeleteProductImage)))
	servermux.Handle("GET /api/categories", http.HandlerFunc(productHandler.HandleListCategories))
	servermux.Handle("GET /api/categories/{category_id}/subcategories", http.HandlerFunc(productHandler.HandleListSubcategories))
	// servermux.HandleFunc("POST /api/orders/{id}/cancel", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/payments", handlers.HandelProducts)
	// servermux.HandleFunc("GET /api/payments/{id}", handlers.HandelProducts)
	// servermux.HandleFunc("POST /api/payments/{id}/verify", handlers.HandelProducts)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
