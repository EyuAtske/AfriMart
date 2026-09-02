package shop

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/config"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type ShopHandler struct {
	Config *config.ApiConfig
}

type createShopRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (apiCfg *ShopHandler) HandleCreateShop(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var params createShopRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if params.Name == "" {
		http.Error(w, "shop name is required", http.StatusBadRequest)
		return
	}

	shop, err := apiCfg.Config.Queries.CreateShop(r.Context(), database.CreateShopParams{
		ID:      uuid.New(),
		OwnerID: userID,
		Name:    params.Name,
		Description: sql.NullString{
			String: params.Description,
			Valid:  true,
		},
	})
	if err != nil {
		http.Error(w, "could not create shop", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(shop)
}
