package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EyuAtske/AfriMart/backend/internal/observability"

	database "github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/EyuAtske/AfriMart/backend/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	godotenv.Load()
	shutdownTracer, err := observability.InitTracer(context.Background())

	if err != nil {
		log.Printf("warning: failed to initialize tracing: %v", err)
	} else {
		defer func() {
			if err := shutdownTracer(context.Background()); err != nil {
				log.Printf("failed to shutdown tracer: %v", err)
			}
		}()
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable is required")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)
	apicfg := handlers.ApiConfig{
		DB:      dbConn,
		Queries: dbQueries,
		Secret:  secretKey,
	}
	servermux := http.NewServeMux()
	server := &http.Server{
		Handler: otelhttp.NewHandler(servermux, "afrimart-backend"),
		Addr:    ":8080",
	}
	servermux.HandleFunc("GET /api/health", handlers.HandelHealth)
	servermux.HandleFunc("POST /api/auth/register", apicfg.HandleRegister)
	servermux.HandleFunc("POST /api/auth/login", apicfg.HandleLogIn)
	servermux.HandleFunc("POST /api/auth/logout", apicfg.HandleRevoke)
	servermux.HandleFunc("PUT /api/auth/user", apicfg.HandleUpdates)
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
	err = server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
	}
}
