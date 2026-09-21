package shop

import (
	"database/sql"
	"encoding/json"
	"errors"
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

func (apiCfg *ShopHandler) HandleCreateShop(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	var params createShopRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error decoding params", err)
		return
	}

	if err := validateCreateShopRequest(&params); err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	shop, err := apiCfg.Config.Queries.CreateShop(
		r.Context(),
		database.CreateShopParams{
			ID:      uuid.New(),
			OwnerID: userID,
			Name:    params.Name,
			Description: sql.NullString{
				String: params.Description,
				Valid:  params.Description != "",
			},
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "could not create shop", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(shop)
}

func (apiCfg *ShopHandler) HandleGetMyShop(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	shop, err := apiCfg.Config.Queries.GetShopByOwnerID(
		r.Context(),
		userID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Shop not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get shop",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		return
	}
}

func (apiCfg *ShopHandler) HandleUpdateShopName(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	shopID, err := uuid.Parse(r.PathValue("shopID"))
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid shop ID",
			err,
		)
		return
	}

	var params updateShopNameRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Error decoding params",
			err,
		)
		return
	}

	params.Name = strings.TrimSpace(params.Name)

	if params.Name == "" {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Shop name is required",
			nil,
		)
		return
	}

	shop, err := apiCfg.Config.Queries.UpdateShopName(
		r.Context(),
		database.UpdateShopNameParams{
			ID:      shopID,
			Name:    params.Name,
			OwnerID: userID,
		},
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Shop not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not update shop name",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		return
	}
}

func (apiCfg *ShopHandler) HandleUpdateShopDescription(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	shopID, err := uuid.Parse(r.PathValue("shopID"))
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid shop ID",
			err,
		)
		return
	}

	var params updateShopDescriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Error decoding params",
			err,
		)
		return
	}

	params.Description = strings.TrimSpace(params.Description)

	shop, err := apiCfg.Config.Queries.UpdateShopDescription(
		r.Context(),
		database.UpdateShopDescriptionParams{
			ID: shopID,
			Description: sql.NullString{
				String: params.Description,
				Valid:  params.Description != "",
			},
			OwnerID: userID,
		},
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Shop not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not update shop description",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shop); err != nil {
		return
	}
}

func (apiCfg *ShopHandler) HandleDeactivateShop(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	shopId, err := uuid.Parse(r.PathValue("shopID"))
	if err != nil{
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid shop Id",
			err,
		)
		return
	} 

	shop, err := apiCfg.Queries.DeactivateShop(
		r.Context(),
		database.DeactivateShopParams{
			ID: shopId,
			OwnerID: userId,
		},
	)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Shop not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not deactivate shop",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(shop); err != nil{
		return
	}
}

func (apiCfg *ShopHandler) HandleActivateShop(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	shopId, err := uuid.Parse(r.PathValue("shopID"))
	if err != nil{
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid shop Id",
			err,
		)
		return
	} 

	shop, err := apiCfg.Queries.ActivateShop(
		r.Context(),
		database.ActivateShopParams{
			ID: shopId,
			OwnerID: userId,
		},
	)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Shop not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not activate shop",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(shop); err != nil{
		return
	}
}