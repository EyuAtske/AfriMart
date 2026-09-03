package shop

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/commErr"
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
		commErr.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	var params createShopRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		commErr.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error decoding params", err)
		return
	}

	if params.Name == "" {
		commErr.RespondErrorWithJson(w, r, http.StatusBadRequest, "shop name is required", nil)
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
		commErr.RespondErrorWithJson(w, r, http.StatusInternalServerError, "could not create shop", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(shop)
}
