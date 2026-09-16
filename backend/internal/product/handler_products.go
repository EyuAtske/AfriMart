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
	"github.com/EyuAtske/AfriMart/backend/internal/comm"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type ProductHandler struct {
	Config      *config.ApiConfig
	Queries     ProductQuerier
	ShopQueries ShopOwnershipQuerier
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

type updateProductRequest struct {
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
		comm.RespondErrorWithJson(
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
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid product status",
			nil,
		)
		return
	}

	if params.Name == "" {
		comm.RespondErrorWithJson(
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
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Price must be a valid non-negative number",
			nil,
		)
		return
	}

	if params.Stock < 0 {
		comm.RespondErrorWithJson(
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
		comm.RespondErrorWithJson(
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
		comm.RespondErrorWithJson(
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
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid subcategory ID",
			err,
		)
		return
	}
	_, err = apiCfg.ShopQueries.GetShopByIDAndOwnerID(
		r.Context(),
		database.GetShopByIDAndOwnerIDParams{
			ID:      shopID,
			OwnerID: userID,
		},
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusForbidden,
				"You do not own this shop",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
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
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not create product",
			err,
		)
		return
	}
	comm.RespondwithJson(w, r, product)
}

func (apiCfg *ProductHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productIDString := strings.TrimSpace(r.PathValue("id"))

	productID, err := uuid.Parse(productIDString)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid product ID",
			err,
		)
		return
	}

	product, err := apiCfg.Queries.GetProduct(
		r.Context(),
		productID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Product not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get product",
			err,
		)
		return
	}

	comm.RespondwithJson(w, r, product)
}

func (apiCfg *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	productIDString := strings.TrimSpace(r.PathValue("id"))

	productID, err := uuid.Parse(productIDString)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid product ID",
			err,
		)
		return
	}

	existingProduct, err := apiCfg.Queries.GetProduct(
		r.Context(),
		productID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Product not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get product",
			err,
		)
		return
	}

	_, err = apiCfg.ShopQueries.GetShopByIDAndOwnerID(
		r.Context(),
		database.GetShopByIDAndOwnerIDParams{
			ID:      existingProduct.ShopID,
			OwnerID: userID,
		},
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusForbidden,
				"You do not own this product's shop",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not verify shop ownership",
			err,
		)
		return
	}

	var params updateProductRequest

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
	params.Description = strings.TrimSpace(params.Description)
	params.Brand = strings.TrimSpace(params.Brand)
	params.Color = strings.TrimSpace(params.Color)
	params.Size = strings.TrimSpace(params.Size)
	params.Price = strings.TrimSpace(params.Price)
	params.Image = strings.TrimSpace(params.Image)
	params.Status = strings.TrimSpace(params.Status)
	params.CategoryID = strings.TrimSpace(params.CategoryID)
	params.SubcategoryID = strings.TrimSpace(params.SubcategoryID)

	if params.Status == "" {
		params.Status = "active"
	}

	if params.Status != "active" && params.Status != "inactive" {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid product status",
			nil,
		)
		return
	}

	if params.Name == "" {
		comm.RespondErrorWithJson(
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
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Price must be a valid non-negative number",
			err,
		)
		return
	}

	if params.Stock < 0 {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Stock cannot be negative",
			nil,
		)
		return
	}

	categoryID, err := uuid.Parse(params.CategoryID)
	if err != nil {
		comm.RespondErrorWithJson(
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
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid subcategory ID",
			err,
		)
		return
	}

	product, err := apiCfg.Queries.UpdateProduct(
		r.Context(),
		database.UpdateProductParams{
			ID:            productID,
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
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not update product",
			err,
		)
		return
	}

	comm.RespondwithJson(w, r, product)
}

func (apiCfg *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	productIDString := strings.TrimSpace(r.PathValue("id"))

	productID, err := uuid.Parse(productIDString)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid product ID",
			err,
		)
		return
	}

	existingProduct, err := apiCfg.Queries.GetProduct(
		r.Context(),
		productID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Product not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get product",
			err,
		)
		return
	}

	_, err = apiCfg.ShopQueries.GetShopByIDAndOwnerID(
		r.Context(),
		database.GetShopByIDAndOwnerIDParams{
			ID:      existingProduct.ShopID,
			OwnerID: userID,
		},
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusForbidden,
				"You do not own this product's shop",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not verify shop ownership",
			err,
		)
		return
	}

	err = apiCfg.Queries.DeleteProduct(
		r.Context(),
		productID,
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not delete product",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := strings.TrimSpace(query.Get("search"))
	brand := strings.TrimSpace(query.Get("brand"))
	color := strings.TrimSpace(query.Get("color"))
	size := strings.TrimSpace(query.Get("size"))
	var categoryID uuid.NullUUID
	if category := query.Get("category_id"); category != "" {
		id, err := uuid.Parse(category)
		if err != nil {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid category_id",
				err,
			)
			return
		}

		categoryID = uuid.NullUUID{
			UUID:  id,
			Valid: true,
		}
	}

	// Subcategory
	var subcategoryID uuid.NullUUID
	if subcategory := query.Get("subcategory_id"); subcategory != "" {
		id, err := uuid.Parse(subcategory)
		if err != nil {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid subcategory_id",
				err,
			)
			return
		}

		subcategoryID = uuid.NullUUID{
			UUID:  id,
			Valid: true,
		}
	}

	// Minimum price
	var minPrice sql.NullString
	if value := query.Get("min_price"); value != "" {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid min_price",
				err,
			)
			return
		}

		minPrice = sql.NullString{
			String: value,
			Valid:  true,
		}
	}

	// Maximum price
	var maxPrice sql.NullString
	if value := query.Get("max_price"); value != "" {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid max_price",
				err,
			)
			return
		}

		maxPrice = sql.NullString{
			String: value,
			Valid:  true,
		}
	}

	// Pagination
	limit := 20
	offset := 0

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid limit",
				err,
			)
			return
		}

		limit = parsed
	}

	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid offset",
				err,
			)
			return
		}

		offset = parsed
	}

	products, err := h.Queries.ListProducts(
		r.Context(),
		database.ListProductsParams{
			Search:        search,
			CategoryID:    categoryID,
			SubcategoryID: subcategoryID,
			Brand:         brand,
			Color:         color,
			Size:          size,
			MinPrice:      minPrice,
			MaxPrice:      maxPrice,
			PageOffset:    int32(offset),
			PageLimit:     int32(limit),
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to retrieve products",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(products); err != nil {
		return
	}
}

func (apiCfg *ProductHandler) HandleListProductsByShop(w http.ResponseWriter, r *http.Request) {
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

	shopIDString := strings.TrimSpace(r.PathValue("shop_id"))

	shopID, err := uuid.Parse(shopIDString)
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

	_, err = apiCfg.ShopQueries.GetShopByIDAndOwnerID(
		r.Context(),
		database.GetShopByIDAndOwnerIDParams{
			ID:      shopID,
			OwnerID: userID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusForbidden,
				"You do not own this shop",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not verify shop ownership",
			err,
		)
		return
	}

	limit := 20
	offset := 0

	if value := strings.TrimSpace(r.URL.Query().Get("limit")); value != "" {
		parsedLimit, err := strconv.Atoi(value)
		if err != nil || parsedLimit <= 0 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid limit",
				err,
			)
			return
		}

		limit = parsedLimit
	}

	if value := strings.TrimSpace(r.URL.Query().Get("offset")); value != "" {
		parsedOffset, err := strconv.Atoi(value)
		if err != nil || parsedOffset < 0 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid offset",
				err,
			)
			return
		}

		offset = parsedOffset
	}

	products, err := apiCfg.Queries.ListProductsByShop(
		r.Context(),
		database.ListProductsByShopParams{
			ShopID: shopID,
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not list products by shop",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(products); err != nil {
		return
	}
}

func (apiCfg *ProductHandler) HandleListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDString := strings.TrimSpace(r.PathValue("category_id"))

	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid category ID",
			err,
		)
		return
	}

	limit := 20
	offset := 0

	if value := strings.TrimSpace(r.URL.Query().Get("limit")); value != "" {
		parsedLimit, err := strconv.Atoi(value)
		if err != nil || parsedLimit <= 0 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid limit",
				err,
			)
			return
		}

		limit = parsedLimit
	}

	if value := strings.TrimSpace(r.URL.Query().Get("offset")); value != "" {
		parsedOffset, err := strconv.Atoi(value)
		if err != nil || parsedOffset < 0 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid offset",
				err,
			)
			return
		}

		offset = parsedOffset
	}

	products, err := apiCfg.Queries.ListProductsByCategory(
		r.Context(),
		database.ListProductsByCategoryParams{
			CategoryID: categoryID,
			Limit:      int32(limit),
			Offset:     int32(offset),
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not list products by category",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(products); err != nil {
		return
	}
}

func (apiCfg *ProductHandler) HandleListProductsBySubcategory(w http.ResponseWriter, r *http.Request) {
	subcategoryIDString := strings.TrimSpace(r.PathValue("subcategory_id"))

	subcategoryID, err := uuid.Parse(subcategoryIDString)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid subcategory ID",
			err,
		)
		return
	}

	limit := 20
	offset := 0

	if value := strings.TrimSpace(r.URL.Query().Get("limit")); value != "" {
		parsedLimit, err := strconv.Atoi(value)
		if err != nil || parsedLimit <= 0 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid limit",
				err,
			)
			return
		}

		limit = parsedLimit
	}

	if value := strings.TrimSpace(r.URL.Query().Get("offset")); value != "" {
		parsedOffset, err := strconv.Atoi(value)
		if err != nil || parsedOffset < 0 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid offset",
				err,
			)
			return
		}

		offset = parsedOffset
	}

	products, err := apiCfg.Queries.ListProductsBySubcategory(
		r.Context(),
		database.ListProductsBySubcategoryParams{
			SubcategoryID: subcategoryID,
			Limit:         int32(limit),
			Offset:        int32(offset),
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not list products by subcategory",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(products); err != nil {
		return
	}
}
