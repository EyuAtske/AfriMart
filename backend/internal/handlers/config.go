package handlers

import (
	"database/sql"

	database "github.com/EyuAtske/AfriMart/backend/internal/database"
)

type ApiConfig struct {
	DB      *sql.DB
	Queries *database.Queries
	Secret 	string
}
