package shop

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/comm"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type ShopHandler struct {
	Config  *config.ApiConfig
	Queries ShopQuerier
	Logger  *slog.Logger
}

func NewShopHandler(cfg *config.ApiConfig, queries ShopQuerier, logger *slog.Logger) *ShopHandler {
	return &ShopHandler{
		Config:  cfg,
		Queries: queries,
		Logger:  logger,
	}
}

type createShopRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateShopNameRequest struct {
	Name string `json:"name"`
}

type updateShopDescriptionRequest struct {
	Description string `json:"description"`
}

func (h *ShopHandler) HandleCreateShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "create shop failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}
	h.Logger.InfoContext(ctx, "handling create shop request", "user_id", userID)

	var params createShopRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.Logger.WarnContext(ctx, "create shop failed: invalid request body", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error decoding params", err)
		return
	}

	if err := validateCreateShopRequest(&params); err != nil {
		h.Logger.WarnContext(ctx, "create shop failed: validation error", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	shop, err := h.Config.Queries.CreateShop(ctx, database.CreateShopParams{
		ID:      uuid.New(),
		OwnerID: userID,
		Name:    params.Name,
		Description: sql.NullString{
			String: params.Description,
			Valid:  params.Description != "",
		},
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "create shop failed: database error", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "could not create shop", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		h.Logger.ErrorContext(ctx, "create shop succeeded but failed to encode response", "user_id", userID, "shop_id", shop.ID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "shop created successfully", "user_id", userID, "shop_id", shop.ID)
}

func (h *ShopHandler) HandleGetMyShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "get my shop failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}
	h.Logger.InfoContext(ctx, "handling get my shop request", "user_id", userID)

	shop, err := h.Config.Queries.GetShopByOwnerID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "get my shop failed: shop not found", "user_id", userID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Shop not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "get my shop failed: database error", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get shop", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		h.Logger.ErrorContext(ctx, "get my shop succeeded but failed to encode response", "user_id", userID, "shop_id", shop.ID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "shop retrieved successfully", "user_id", userID, "shop_id", shop.ID)
}

func (h *ShopHandler) HandleUpdateShopName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "update shop name failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	shopIDStr := r.PathValue("shopID")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "update shop name failed: invalid shop ID", "user_id", userID, "shop_id", shopIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid shop ID", err)
		return
	}
	h.Logger.InfoContext(ctx, "handling update shop name request", "user_id", userID, "shop_id", shopID)

	var params updateShopNameRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.Logger.WarnContext(ctx, "update shop name failed: invalid request body", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error decoding params", err)
		return
	}

	params.Name = strings.TrimSpace(params.Name)
	if params.Name == "" {
		h.Logger.WarnContext(ctx, "update shop name failed: name is required", "user_id", userID, "shop_id", shopID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Shop name is required", nil)
		return
	}

	shop, err := h.Config.Queries.UpdateShopName(ctx, database.UpdateShopNameParams{
		ID:      shopID,
		Name:    params.Name,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "update shop name failed: shop not found", "user_id", userID, "shop_id", shopID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Shop not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "update shop name failed: database error", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not update shop name", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		h.Logger.ErrorContext(ctx, "update shop name succeeded but failed to encode response", "user_id", userID, "shop_id", shopID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "shop name updated successfully", "user_id", userID, "shop_id", shopID, "new_name", params.Name)
}

func (h *ShopHandler) HandleUpdateShopDescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "update shop description failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	shopIDStr := r.PathValue("shopID")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "update shop description failed: invalid shop ID", "user_id", userID, "shop_id", shopIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid shop ID", err)
		return
	}
	h.Logger.InfoContext(ctx, "handling update shop description request", "user_id", userID, "shop_id", shopID)

	var params updateShopDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.Logger.WarnContext(ctx, "update shop description failed: invalid request body", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error decoding params", err)
		return
	}

	params.Description = strings.TrimSpace(params.Description)

	shop, err := h.Config.Queries.UpdateShopDescription(ctx, database.UpdateShopDescriptionParams{
		ID: shopID,
		Description: sql.NullString{
			String: params.Description,
			Valid:  params.Description != "",
		},
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "update shop description failed: shop not found", "user_id", userID, "shop_id", shopID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Shop not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "update shop description failed: database error", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not update shop description", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		h.Logger.ErrorContext(ctx, "update shop description succeeded but failed to encode response", "user_id", userID, "shop_id", shopID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "shop description updated successfully", "user_id", userID, "shop_id", shopID)
}

func (h *ShopHandler) HandleDeactivateShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "deactivate shop failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	shopIDStr := r.PathValue("shopID")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "deactivate shop failed: invalid shop ID", "user_id", userID, "shop_id", shopIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid shop ID", err)
		return
	}
	h.Logger.InfoContext(ctx, "handling deactivate shop request", "user_id", userID, "shop_id", shopID)

	shop, err := h.Queries.DeactivateShop(ctx, database.DeactivateShopParams{
		ID:      shopID,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "deactivate shop failed: shop not found", "user_id", userID, "shop_id", shopID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Shop not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "deactivate shop failed: database error", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not deactivate shop", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		h.Logger.ErrorContext(ctx, "deactivate shop succeeded but failed to encode response", "user_id", userID, "shop_id", shopID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "shop deactivated successfully", "user_id", userID, "shop_id", shopID)
}

func (h *ShopHandler) HandleActivateShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "activate shop failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	shopIDStr := r.PathValue("shopID")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "activate shop failed: invalid shop ID", "user_id", userID, "shop_id", shopIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid shop ID", err)
		return
	}
	h.Logger.InfoContext(ctx, "handling activate shop request", "user_id", userID, "shop_id", shopID)

	shop, err := h.Queries.ActivateShop(ctx, database.ActivateShopParams{
		ID:      shopID,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "activate shop failed: shop not found", "user_id", userID, "shop_id", shopID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Shop not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "activate shop failed: database error", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not activate shop", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		h.Logger.ErrorContext(ctx, "activate shop succeeded but failed to encode response", "user_id", userID, "shop_id", shopID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "shop activated successfully", "user_id", userID, "shop_id", shopID)
}