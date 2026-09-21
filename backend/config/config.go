package config

import (
	"database/sql"

	database "github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/EyuAtske/AfriMart/backend/internal/storage"
)

type ApiConfig struct {
	DB           *sql.DB
	Queries      *database.Queries
	Secret       string
	ImageStorage storage.ImageStorage
}
