package config

import (
	"context"
	"log/slog"
	"os"

	database "github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/EyuAtske/AfriMart/backend/internal/storage"
	"github.com/XSAM/otelsql"
)

func SetupAPIConfig(ctx context.Context) *ApiConfig {
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
			"error opening database",
			"error", err,
		)
		os.Exit(1)
	}

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

	imageStorage, err := storage.NewMinioStorage(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_BUCKET"),
		false,
	)

	if err != nil {
		slog.Error("failed to initialize MinIO storage", "error", err)
		return nil
	}

	chapaSecretKey := os.Getenv("CHAPA_SECRET_KEY")
	chapaCallbackURL := os.Getenv("CHAPA_CALLBACK_URL")

	return &ApiConfig{
		DB:               dbConn,
		Queries:          dbQueries,
		Secret:           secretKey,
		ImageStorage:     imageStorage,
		ChapaSecretKey:   chapaSecretKey,
		ChapaCallbackURL: chapaCallbackURL,
	}
}
