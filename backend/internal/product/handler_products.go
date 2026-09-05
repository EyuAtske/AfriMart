package product

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/commErr"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type ProductHandler struct {
	Config  *config.ApiConfig
	Queries ProductQuerier
}

type createProductRequest struct {
	ShopID        string `json:"shop_id"`
	CategoryID    string `json:"category_id"`
	SubcategoryID string `json:"subcategory_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Brand         string `json:"brand"`
	Color         string `json:"color"`
	Size          string `json:"size"`
	Price         string `json:"price"`
	Stock         int32  `json:"stock"`
	Image         string `json:"image"`
	Status        string `json:"status"`
}

func (apiCfg *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	var params createProductRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Error decoding params",
			err,
		)
		return
	}

	params.Name = strings.TrimSpace(params.Name)
	params.Description = strings.TrimSpace(params.Description)
	params.Brand = strings.TrimSpace(params.Brand)
	params.Color = strings.TrimSpace(params.Color)
	params.Size = strings.TrimSpace(params.Size)
	params.Price = strings.TrimSpace(params.Price)
	params.Image = strings.TrimSpace(params.Image)
	params.Status = strings.TrimSpace(params.Status)
	params.ShopID = strings.TrimSpace(params.ShopID)
	params.CategoryID = strings.TrimSpace(params.CategoryID)
	params.SubcategoryID = strings.TrimSpace(params.SubcategoryID)

	if params.Status == "" {
		params.Status = "active"
	}
	if params.Status != "active" && params.Status != "inactive" {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid product status",
			nil,
		)
		return
	}

	if params.Name == "" {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Product name is required",
			nil,
		)
		return
	}

	price, err := strconv.ParseFloat(params.Price, 64)
	if err != nil || price < 0 {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Price must be a valid non-negative number",
			nil,
		)
		return
	}

	if params.Stock < 0 {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Stock cannot be negative",
			nil,
		)
		return
	}

	shopID, err := uuid.Parse(params.ShopID)
	if err != nil {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid shop ID",
			err,
		)
		return
	}

	categoryID, err := uuid.Parse(params.CategoryID)
	if err != nil {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid category ID",
			err,
		)
		return
	}

	subcategoryID, err := uuid.Parse(params.SubcategoryID)
	if err != nil {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid subcategory ID",
			err,
		)
		return
	}
	_, err = apiCfg.Config.Queries.GetShopByIDAndOwnerID(
		r.Context(),
		database.GetShopByIDAndOwnerIDParams{
			ID:      shopID,
			OwnerID: userID,
		},
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			commErr.RespondErrorWithJson(
				w,
				r,
				http.StatusForbidden,
				"You do not own this shop",
				err,
			)
			return
		}

		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not verify shop ownership",
			err,
		)
		return
	}
	product, err := apiCfg.Queries.CreateProduct(
		r.Context(),
		database.CreateProductParams{
			ShopID:        shopID,
			CategoryID:    categoryID,
			SubcategoryID: subcategoryID,
			Name:          params.Name,
			Description: sql.NullString{
				String: params.Description,
				Valid:  params.Description != "",
			},
			Brand: sql.NullString{
				String: params.Brand,
				Valid:  params.Brand != "",
			},
			Color: sql.NullString{
				String: params.Color,
				Valid:  params.Color != "",
			},
			Size: sql.NullString{
				String: params.Size,
				Valid:  params.Size != "",
			},
			Price: params.Price,
			Stock: params.Stock,
			Image: sql.NullString{
				String: params.Image,
				Valid:  params.Image != "",
			},
			Status: params.Status,
		},
	)
	if err != nil {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not create product",
			err,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		return
	}
}
