package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/comm"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/EyuAtske/AfriMart/backend/internal/storage"
	"github.com/google/uuid"
)

type ProductHandler struct {
	Config       *config.ApiConfig
	Queries      ProductQuerier
	ShopQueries  ShopOwnershipQuerier
	ImageStorage storage.ImageStorage
	Logger       *slog.Logger 
}

type ProductWithImages struct {
	Product database.Product        `json:"product"`
	Images  []database.ProductImage `json:"images"`
}

// enrichProductsWithImages fetches all images for a slice of products in a
// single query and returns them paired together. This avoids the N+1 problem.
func (h *ProductHandler) enrichProductsWithImages(ctx context.Context, products []database.Product) ([]ProductWithImages, error) {
	if len(products) == 0 {
		return []ProductWithImages{}, nil
	}

	productIDs := make([]uuid.UUID, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
	}

	allImages, err := h.Queries.GetProductImagesByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	// Group images by product_id
	imageMap := make(map[uuid.UUID][]database.ProductImage, len(products))
	for _, img := range allImages {
		imageMap[img.ProductID] = append(imageMap[img.ProductID], img)
	}

	result := make([]ProductWithImages, len(products))
	for i, p := range products {
		images := imageMap[p.ID]
		if images == nil {
			images = []database.ProductImage{} // never null in JSON
		}
		result[i] = ProductWithImages{
			Product: p,
			Images:  images,
		}
	}

	return result, nil
}

func NewProductHandler(cfg *config.ApiConfig, queries ProductQuerier, imageStorage storage.ImageStorage, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		Config:       cfg,
		Queries:      queries,
		ImageStorage: imageStorage,
		Logger:       logger,
	}
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
	Status        string `json:"status"`
}

func (h *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "create product failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}
	h.Logger.InfoContext(ctx, "handling create product request", "user_id", userID)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.Logger.WarnContext(ctx, "create product failed: invalid multipart form", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid multipart form", err)
		return
	}

	params := createProductRequest{
		ShopID:        strings.TrimSpace(r.FormValue("shop_id")),
		CategoryID:    strings.TrimSpace(r.FormValue("category_id")),
		SubcategoryID: strings.TrimSpace(r.FormValue("subcategory_id")),
		Name:          strings.TrimSpace(r.FormValue("name")),
		Description:   strings.TrimSpace(r.FormValue("description")),
		Brand:         strings.TrimSpace(r.FormValue("brand")),
		Color:         strings.TrimSpace(r.FormValue("color")),
		Size:          strings.TrimSpace(r.FormValue("size")),
		Price:         strings.TrimSpace(r.FormValue("price")),
		Status:        strings.TrimSpace(r.FormValue("status")),
	}

	stockString := strings.TrimSpace(r.FormValue("stock"))
	stock, err := strconv.ParseInt(stockString, 10, 32)
	if err != nil || stock < 0 {
		h.Logger.WarnContext(ctx, "create product failed: invalid stock", "user_id", userID, "stock", stockString, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Stock must be a valid non-negative integer", nil)
		return
	}
	params.Stock = int32(stock)

	if params.Status == "" {
		params.Status = "active"
	}
	if params.Status != "active" && params.Status != "inactive" {
		h.Logger.WarnContext(ctx, "create product failed: invalid status", "user_id", userID, "status", params.Status)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid product status", nil)
		return
	}
	if params.Name == "" {
		h.Logger.WarnContext(ctx, "create product failed: name required", "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Product name is required", nil)
		return
	}

	price, err := strconv.ParseFloat(params.Price, 64)
	if err != nil || price < 0 {
		h.Logger.WarnContext(ctx, "create product failed: invalid price", "user_id", userID, "price", params.Price, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Price must be a valid non-negative number", nil)
		return
	}

	shopID, err := uuid.Parse(params.ShopID)
	if err != nil {
		h.Logger.WarnContext(ctx, "create product failed: invalid shop ID", "user_id", userID, "shop_id", params.ShopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid shop ID", err)
		return
	}

	categoryID, err := uuid.Parse(params.CategoryID)
	if err != nil {
		h.Logger.WarnContext(ctx, "create product failed: invalid category ID", "user_id", userID, "category_id", params.CategoryID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	subcategoryID, err := uuid.Parse(params.SubcategoryID)
	if err != nil {
		h.Logger.WarnContext(ctx, "create product failed: invalid subcategory ID", "user_id", userID, "subcategory_id", params.SubcategoryID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid subcategory ID", err)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		h.Logger.WarnContext(ctx, "create product failed: no images provided", "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "At least one product image is required", nil)
		return
	}
	if len(files) > maxProductImageCount {
		h.Logger.WarnContext(ctx, "create product failed: too many images", "user_id", userID, "count", len(files), "max", maxProductImageCount)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "A product can have at most 10 images", nil)
		return
	}

	for _, file := range files {
		if err := validateProductImage(file); err != nil {
			h.Logger.WarnContext(ctx, "create product failed: image validation error", "user_id", userID, "filename", file.Filename, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}

	_, err = h.ShopQueries.GetShopByIDAndOwnerID(ctx, database.GetShopByIDAndOwnerIDParams{
		ID:      shopID,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "create product failed: user does not own shop", "user_id", userID, "shop_id", shopID)
			comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not own this shop", err)
			return
		}
		h.Logger.ErrorContext(ctx, "create product failed: could not verify shop ownership", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not verify shop ownership", err)
		return
	}

	product, err := h.Queries.CreateProduct(ctx, database.CreateProductParams{
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
		Price:  params.Price,
		Stock:  params.Stock,
		Status: params.Status,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "create product failed: database error creating product record", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not create product", err)
		return
	}

	uploadedObjects := make([]string, 0, len(files))
	
	// Enhanced cleanup function that logs the trigger and the result of every step
	cleanup := func(originalErr error) {
		h.Logger.ErrorContext(ctx, "initiating rollback/cleanup due to previous error", 
			"user_id", userID, 
			"product_id", product.ID, 
			"original_error", originalErr)

		for _, objectKey := range uploadedObjects {
			delErr := h.ImageStorage.Delete(context.Background(), objectKey)
			if delErr != nil {
				h.Logger.ErrorContext(ctx, "cleanup failed: could not delete orphaned image from storage", 
					"object_key", objectKey, 
					"error", delErr)
			} else {
				h.Logger.InfoContext(ctx, "cleanup successful: deleted orphaned image from storage", 
					"object_key", objectKey)
			}
		}

		delErr := h.Queries.DeleteProduct(context.Background(), product.ID)
		if delErr != nil {
			h.Logger.ErrorContext(ctx, "cleanup failed: could not delete orphaned product from database", 
				"product_id", product.ID, 
				"error", delErr)
		} else {
			h.Logger.InfoContext(ctx, "cleanup successful: deleted orphaned product from database", 
				"product_id", product.ID)
		}
	}

	for index, file := range files {
		src, err := file.Open()
		if err != nil {
			h.Logger.ErrorContext(ctx, "create product failed: could not open image file", "user_id", userID, "product_id", product.ID, "filename", file.Filename, "error", err)
			cleanup(err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not open product image", err)
			return
		}

		contentTypeBuffer := make([]byte, 512)
		n, err := src.Read(contentTypeBuffer)
		if err != nil {
			src.Close()
			h.Logger.ErrorContext(ctx, "create product failed: could not read image file", "user_id", userID, "product_id", product.ID, "filename", file.Filename, "error", err)
			cleanup(err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not read product image", err)
			return
		}

		contentType := http.DetectContentType(contentTypeBuffer[:n])
		if _, err := src.Seek(0, io.SeekStart); err != nil {
			src.Close()
			h.Logger.ErrorContext(ctx, "create product failed: could not reset image file reader", "user_id", userID, "product_id", product.ID, "error", err)
			cleanup(err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not reset product image", err)
			return
		}

		extension := imageExtension(contentType)
		objectKey := fmt.Sprintf("products/%s/%s%s", product.ID.String(), uuid.New().String(), extension)

		err = h.ImageStorage.Upload(ctx, objectKey, src, file.Size, contentType)
		src.Close()

		if err != nil {
			h.Logger.ErrorContext(ctx, "create product failed: storage upload error", "user_id", userID, "product_id", product.ID, "object_key", objectKey, "error", err)
			cleanup(err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not upload product image", err)
			return
		}

		uploadedObjects = append(uploadedObjects, objectKey)

		_, err = h.Queries.CreateProductImage(ctx, database.CreateProductImageParams{
			ProductID:    product.ID,
			ObjectKey:    objectKey,
			DisplayOrder: int32(index),
		})
		if err != nil {
			h.Logger.ErrorContext(ctx, "create product failed: database error saving image record", "user_id", userID, "product_id", product.ID, "object_key", objectKey, "error", err)
			cleanup(err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not save product image", err)
			return
		}
	}

	response := struct {
		Product database.Product        `json:"product"`
		Images  []database.ProductImage `json:"images"`
	}{
		Product: product,
	}

	response.Images, err = h.Queries.GetProductImages(ctx, product.ID)
	if err != nil {
		// Note: We don't call cleanup() here because the product and images ARE successfully saved.
		// We just log the fetch error and still return a 201 Created with the product data.
		h.Logger.ErrorContext(ctx, "create product succeeded but failed to fetch images for response", "user_id", userID, "product_id", product.ID, "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "create product succeeded but failed to encode response", "user_id", userID, "product_id", product.ID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "product created successfully", "user_id", userID, "product_id", product.ID, "shop_id", shopID, "image_count", len(response.Images))
}

func (h *ProductHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productIDString := strings.TrimSpace(r.PathValue("id"))

	productID, err := uuid.Parse(productIDString)
	if err != nil {
		h.Logger.WarnContext(ctx, "get product failed: invalid product ID", "product_id", productIDString, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling get product request", "product_id", productID)

	product, err := h.Queries.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "get product failed: product not found", "product_id", productID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "get product failed: database error", "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product", err)
		return
	}

	images, err := h.Queries.GetProductImages(ctx, productID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "get product failed: could not fetch images", "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product images", err)
		return
	}

	response := struct {
		Product database.Product        `json:"product"`
		Images  []database.ProductImage `json:"images"`
	}{
		Product: product,
		Images:  images,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "get product succeeded but failed to encode response", "product_id", productID, "error", err)
	}
}

func (h *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "update product failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	productIDString := strings.TrimSpace(r.PathValue("id"))
	productID, err := uuid.Parse(productIDString)
	if err != nil {
		h.Logger.WarnContext(ctx, "update product failed: invalid product ID", "user_id", userID, "product_id", productIDString, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling update product request", "user_id", userID, "product_id", productID)

	existingProduct, err := h.Queries.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "update product failed: product not found", "user_id", userID, "product_id", productID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "update product failed: database error fetching product", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product", err)
		return
	}

	_, err = h.ShopQueries.GetShopByIDAndOwnerID(ctx, database.GetShopByIDAndOwnerIDParams{
		ID:      existingProduct.ShopID,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "update product failed: user does not own shop", "user_id", userID, "product_id", productID, "shop_id", existingProduct.ShopID)
			comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not own this product's shop", err)
			return
		}
		h.Logger.ErrorContext(ctx, "update product failed: could not verify shop ownership", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not verify shop ownership", err)
		return
	}

	var params updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.Logger.WarnContext(ctx, "update product failed: invalid request body", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error decoding params", err)
		return
	}

	params.Name = strings.TrimSpace(params.Name)
	params.Description = strings.TrimSpace(params.Description)
	params.Brand = strings.TrimSpace(params.Brand)
	params.Color = strings.TrimSpace(params.Color)
	params.Size = strings.TrimSpace(params.Size)
	params.Price = strings.TrimSpace(params.Price)
	params.Status = strings.TrimSpace(params.Status)
	params.CategoryID = strings.TrimSpace(params.CategoryID)
	params.SubcategoryID = strings.TrimSpace(params.SubcategoryID)

	if params.Status == "" {
		params.Status = "active"
	}
	if params.Status != "active" && params.Status != "inactive" {
		h.Logger.WarnContext(ctx, "update product failed: invalid status", "user_id", userID, "product_id", productID, "status", params.Status)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid product status", nil)
		return
	}
	if params.Name == "" {
		h.Logger.WarnContext(ctx, "update product failed: name required", "user_id", userID, "product_id", productID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Product name is required", nil)
		return
	}

	price, err := strconv.ParseFloat(params.Price, 64)
	if err != nil || price < 0 {
		h.Logger.WarnContext(ctx, "update product failed: invalid price", "user_id", userID, "product_id", productID, "price", params.Price, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Price must be a valid non-negative number", nil)
		return
	}
	if params.Stock < 0 {
		h.Logger.WarnContext(ctx, "update product failed: invalid stock", "user_id", userID, "product_id", productID, "stock", params.Stock)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Stock cannot be negative", nil)
		return
	}

	categoryID, err := uuid.Parse(params.CategoryID)
	if err != nil {
		h.Logger.WarnContext(ctx, "update product failed: invalid category ID", "user_id", userID, "product_id", productID, "category_id", params.CategoryID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	subcategoryID, err := uuid.Parse(params.SubcategoryID)
	if err != nil {
		h.Logger.WarnContext(ctx, "update product failed: invalid subcategory ID", "user_id", userID, "product_id", productID, "subcategory_id", params.SubcategoryID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid subcategory ID", err)
		return
	}

	product, err := h.Queries.UpdateProduct(ctx, database.UpdateProductParams{
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
		Price:  params.Price,
		Stock:  params.Stock,
		Status: params.Status,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "update product failed: database error", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not update product", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		h.Logger.ErrorContext(ctx, "update product succeeded but failed to encode response", "user_id", userID, "product_id", productID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "product updated successfully", "user_id", userID, "product_id", productID)
}

func (h *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "delete product failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	productIDString := strings.TrimSpace(r.PathValue("id"))
	productID, err := uuid.Parse(productIDString)
	if err != nil {
		h.Logger.WarnContext(ctx, "delete product failed: invalid product ID", "user_id", userID, "product_id", productIDString, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling delete product request", "user_id", userID, "product_id", productID)

	existingProduct, err := h.Queries.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "delete product failed: product not found", "user_id", userID, "product_id", productID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "delete product failed: database error fetching product", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product", err)
		return
	}

	_, err = h.ShopQueries.GetShopByIDAndOwnerID(ctx, database.GetShopByIDAndOwnerIDParams{
		ID:      existingProduct.ShopID,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "delete product failed: user does not own shop", "user_id", userID, "product_id", productID, "shop_id", existingProduct.ShopID)
			comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not own this product's shop", err)
			return
		}
		h.Logger.ErrorContext(ctx, "delete product failed: could not verify shop ownership", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not verify shop ownership", err)
		return
	}

	images, err := h.Queries.GetProductImages(ctx, productID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "delete product failed: could not fetch product images", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product images", err)
		return
	}

	for _, image := range images {
		if err := h.ImageStorage.Delete(ctx, image.ObjectKey); err != nil {
			h.Logger.ErrorContext(ctx, "delete product: failed to delete image from storage", "user_id", userID, "product_id", productID, "object_key", image.ObjectKey, "error", err)
			// Continue to try deleting the DB record anyway
		}
	}

	err = h.Queries.DeleteProduct(ctx, productID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "delete product failed: database error", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not delete product", err)
		return
	}

	h.Logger.InfoContext(ctx, "product deleted successfully", "user_id", userID, "product_id", productID, "images_deleted", len(images))
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "handling list products request")

	query := r.URL.Query()
	search := strings.TrimSpace(query.Get("search"))
	brand := strings.TrimSpace(query.Get("brand"))
	color := strings.TrimSpace(query.Get("color"))
	size := strings.TrimSpace(query.Get("size"))

	var categoryID uuid.NullUUID
	if category := query.Get("category_id"); category != "" {
		id, err := uuid.Parse(category)
		if err != nil {
			h.Logger.WarnContext(ctx, "list products failed: invalid category_id", "category_id", category, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid category_id", err)
			return
		}
		categoryID = uuid.NullUUID{UUID: id, Valid: true}
	}

	var subcategoryID uuid.NullUUID
	if subcategory := query.Get("subcategory_id"); subcategory != "" {
		id, err := uuid.Parse(subcategory)
		if err != nil {
			h.Logger.WarnContext(ctx, "list products failed: invalid subcategory_id", "subcategory_id", subcategory, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid subcategory_id", err)
			return
		}
		subcategoryID = uuid.NullUUID{UUID: id, Valid: true}
	}

	var minPrice sql.NullString
	if value := query.Get("min_price"); value != "" {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			h.Logger.WarnContext(ctx, "list products failed: invalid min_price", "min_price", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid min_price", err)
			return
		}
		minPrice = sql.NullString{String: value, Valid: true}
	}

	var maxPrice sql.NullString
	if value := query.Get("max_price"); value != "" {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			h.Logger.WarnContext(ctx, "list products failed: invalid max_price", "max_price", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid max_price", err)
			return
		}
		maxPrice = sql.NullString{String: value, Valid: true}
	}

	limit := 20
	offset := 0

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			h.Logger.WarnContext(ctx, "list products failed: invalid limit", "limit", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid limit", err)
			return
		}
		limit = parsed
	}

	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			h.Logger.WarnContext(ctx, "list products failed: invalid offset", "offset", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid offset", err)
			return
		}
		offset = parsed
	}

	products, err := h.Queries.ListProducts(ctx, database.ListProductsParams{
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
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "list products failed: database error", "limit", limit, "offset", offset, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to retrieve products", err)
		return
	}

	response, err := h.enrichProductsWithImages(ctx, products)
	if err != nil {
		h.Logger.ErrorContext(ctx, "list products failed: could not fetch images", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to retrieve product images", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "list products succeeded but failed to encode response", "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "products listed successfully", "count", len(response), "limit", limit, "offset", offset)
}

func (h *ProductHandler) HandleListProductsByShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "list products by shop failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	shopIDString := strings.TrimSpace(r.PathValue("shop_id"))
	shopID, err := uuid.Parse(shopIDString)
	if err != nil {
		h.Logger.WarnContext(ctx, "list products by shop failed: invalid shop ID", "user_id", userID, "shop_id", shopIDString, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid shop ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling list products by shop request", "user_id", userID, "shop_id", shopID)

	_, err = h.ShopQueries.GetShopByIDAndOwnerID(ctx, database.GetShopByIDAndOwnerIDParams{
		ID:      shopID,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "list products by shop failed: user does not own shop", "user_id", userID, "shop_id", shopID)
			comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not own this shop", err)
			return
		}
		h.Logger.ErrorContext(ctx, "list products by shop failed: could not verify shop ownership", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not verify shop ownership", err)
		return
	}

	limit := 20
	offset := 0

	if value := strings.TrimSpace(r.URL.Query().Get("limit")); value != "" {
		parsedLimit, err := strconv.Atoi(value)
		if err != nil || parsedLimit <= 0 {
			h.Logger.WarnContext(ctx, "list products by shop failed: invalid limit", "user_id", userID, "shop_id", shopID, "limit", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid limit", err)
			return
		}
		limit = parsedLimit
	}

	if value := strings.TrimSpace(r.URL.Query().Get("offset")); value != "" {
		parsedOffset, err := strconv.Atoi(value)
		if err != nil || parsedOffset < 0 {
			h.Logger.WarnContext(ctx, "list products by shop failed: invalid offset", "user_id", userID, "shop_id", shopID, "offset", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid offset", err)
			return
		}
		offset = parsedOffset
	}

	products, err := h.Queries.ListProductsByShop(ctx, database.ListProductsByShopParams{
		ShopID: shopID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "list products by shop failed: database error", "user_id", userID, "shop_id", shopID, "limit", limit, "offset", offset, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not list products by shop", err)
		return
	}

	response, err := h.enrichProductsWithImages(ctx, products)
	if err != nil {
		h.Logger.ErrorContext(ctx, "list products by shop failed: could not fetch images", "user_id", userID, "shop_id", shopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not retrieve product images", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "list products by shop succeeded but failed to encode response", "user_id", userID, "shop_id", shopID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "products listed by shop successfully", "user_id", userID, "shop_id", shopID, "count", len(response))
}

func (h *ProductHandler) HandleListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryIDString := strings.TrimSpace(r.PathValue("category_id"))

	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		h.Logger.WarnContext(ctx, "list products by category failed: invalid category ID", "category_id", categoryIDString, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling list products by category request", "category_id", categoryID)

	limit := 20
	offset := 0

	if value := strings.TrimSpace(r.URL.Query().Get("limit")); value != "" {
		parsedLimit, err := strconv.Atoi(value)
		if err != nil || parsedLimit <= 0 {
			h.Logger.WarnContext(ctx, "list products by category failed: invalid limit", "category_id", categoryID, "limit", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid limit", err)
			return
		}
		limit = parsedLimit
	}

	if value := strings.TrimSpace(r.URL.Query().Get("offset")); value != "" {
		parsedOffset, err := strconv.Atoi(value)
		if err != nil || parsedOffset < 0 {
			h.Logger.WarnContext(ctx, "list products by category failed: invalid offset", "category_id", categoryID, "offset", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid offset", err)
			return
		}
		offset = parsedOffset
	}

	products, err := h.Queries.ListProductsByCategory(ctx, database.ListProductsByCategoryParams{
		CategoryID: categoryID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "list products by category failed: database error", "category_id", categoryID, "limit", limit, "offset", offset, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not list products by category", err)
		return
	}

	response, err := h.enrichProductsWithImages(ctx, products)
	if err != nil {
		h.Logger.ErrorContext(ctx, "list products by category failed: could not fetch images", "category_id", categoryID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not retrieve product images", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "list products by category succeeded but failed to encode response", "category_id", categoryID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "products listed by category successfully", "category_id", categoryID, "count", len(response))
}

func (h *ProductHandler) HandleListProductsBySubcategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subcategoryIDString := strings.TrimSpace(r.PathValue("subcategory_id"))

	subcategoryID, err := uuid.Parse(subcategoryIDString)
	if err != nil {
		h.Logger.WarnContext(ctx, "list products by subcategory failed: invalid subcategory ID", "subcategory_id", subcategoryIDString, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid subcategory ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling list products by subcategory request", "subcategory_id", subcategoryID)

	limit := 20
	offset := 0

	if value := strings.TrimSpace(r.URL.Query().Get("limit")); value != "" {
		parsedLimit, err := strconv.Atoi(value)
		if err != nil || parsedLimit <= 0 {
			h.Logger.WarnContext(ctx, "list products by subcategory failed: invalid limit", "subcategory_id", subcategoryID, "limit", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid limit", err)
			return
		}
		limit = parsedLimit
	}

	if value := strings.TrimSpace(r.URL.Query().Get("offset")); value != "" {
		parsedOffset, err := strconv.Atoi(value)
		if err != nil || parsedOffset < 0 {
			h.Logger.WarnContext(ctx, "list products by subcategory failed: invalid offset", "subcategory_id", subcategoryID, "offset", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid offset", err)
			return
		}
		offset = parsedOffset
	}

	products, err := h.Queries.ListProductsBySubcategory(ctx, database.ListProductsBySubcategoryParams{
		SubcategoryID: subcategoryID,
		Limit:         int32(limit),
		Offset:        int32(offset),
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "list products by subcategory failed: database error", "subcategory_id", subcategoryID, "limit", limit, "offset", offset, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not list products by subcategory", err)
		return
	}

	response, err := h.enrichProductsWithImages(ctx, products)
	if err != nil {
		h.Logger.ErrorContext(ctx, "list products by subcategory failed: could not fetch images", "subcategory_id", subcategoryID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not retrieve product images", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "list products by subcategory succeeded but failed to encode response", "subcategory_id", subcategoryID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "products listed by subcategory successfully", "subcategory_id", subcategoryID, "count", len(response))
}

func (h *ProductHandler) HandleUpdateProductImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "update product image failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	productIDStr := r.PathValue("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "update product image failed: invalid product ID", "user_id", userID, "product_id", productIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	imageIDStr := r.PathValue("imageID")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "update product image failed: invalid image ID", "user_id", userID, "product_id", productID, "image_id", imageIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid image ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling update product image request", "user_id", userID, "product_id", productID, "image_id", imageID)

	product, err := h.Queries.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "update product image failed: product not found", "user_id", userID, "product_id", productID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "update product image failed: database error fetching product", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product", err)
		return
	}

	_, err = h.ShopQueries.GetShopByIDAndOwnerID(ctx, database.GetShopByIDAndOwnerIDParams{
		ID:      product.ShopID,
		OwnerID: userID,
	})
	if err != nil {
		h.Logger.WarnContext(ctx, "update product image failed: user does not own shop", "user_id", userID, "product_id", productID, "shop_id", product.ShopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not own this product", err)
		return
	}

	oldImage, err := h.Queries.GetProductImage(ctx, database.GetProductImageParams{
		ID:        imageID,
		ProductID: productID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "update product image failed: image not found", "user_id", userID, "product_id", productID, "image_id", imageID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product image not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "update product image failed: database error fetching image", "user_id", userID, "product_id", productID, "image_id", imageID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product image", err)
		return
	}

	if err := r.ParseMultipartForm(8 << 20); err != nil {
		h.Logger.WarnContext(ctx, "update product image failed: invalid multipart form", "user_id", userID, "product_id", productID, "image_id", imageID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid multipart form", err)
		return
	}

	files := r.MultipartForm.File["image"]
	if len(files) != 1 {
		h.Logger.WarnContext(ctx, "update product image failed: expected exactly one image", "user_id", userID, "product_id", productID, "image_id", imageID, "count", len(files))
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Exactly one image is required", errors.New("expected one image file"))
		return
	}

	file := files[0]
	if err := validateProductImage(file); err != nil {
		h.Logger.WarnContext(ctx, "update product image failed: image validation error", "user_id", userID, "product_id", productID, "image_id", imageID, "filename", file.Filename, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	src, err := file.Open()
	if err != nil {
		h.Logger.ErrorContext(ctx, "update product image failed: could not open image file", "user_id", userID, "product_id", productID, "image_id", imageID, "filename", file.Filename, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not open image", err)
		return
	}
	defer src.Close()

	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil {
		h.Logger.ErrorContext(ctx, "update product image failed: could not read image file", "user_id", userID, "product_id", productID, "image_id", imageID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not read image", err)
		return
	}

	contentType := http.DetectContentType(buffer[:n])
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		h.Logger.ErrorContext(ctx, "update product image failed: could not reset image file", "user_id", userID, "product_id", productID, "image_id", imageID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not reset image", err)
		return
	}

	extension := imageExtension(contentType)
	if extension == "" {
		h.Logger.WarnContext(ctx, "update product image failed: unsupported image type", "user_id", userID, "product_id", productID, "image_id", imageID, "content_type", contentType)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Unsupported image type", errors.New("unsupported image type"))
		return
	}

	newObjectKey := fmt.Sprintf("products/%s/%s%s", productID.String(), uuid.New().String(), extension)

	if err := h.ImageStorage.Upload(ctx, newObjectKey, src, file.Size, contentType); err != nil {
		h.Logger.ErrorContext(ctx, "update product image failed: storage upload error", "user_id", userID, "product_id", productID, "image_id", imageID, "object_key", newObjectKey, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not upload image", err)
		return
	}

	updatedImage, err := h.Queries.UpdateProductImage(ctx, database.UpdateProductImageParams{
		ID:        imageID,
		ProductID: productID,
		ObjectKey: newObjectKey,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "update product image failed: database error, rolling back storage", "user_id", userID, "product_id", productID, "image_id", imageID, "object_key", newObjectKey, "error", err)
		if delErr := h.ImageStorage.Delete(ctx, newObjectKey); delErr != nil {
			h.Logger.ErrorContext(ctx, "rollback failed: could not delete newly uploaded image", "object_key", newObjectKey, "error", delErr)
		}
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not update product image", err)
		return
	}

	if err := h.ImageStorage.Delete(ctx, oldImage.ObjectKey); err != nil {
		h.Logger.ErrorContext(ctx, "update product image succeeded but failed to delete old image from storage", "user_id", userID, "product_id", productID, "image_id", imageID, "old_object_key", oldImage.ObjectKey, "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"image": updatedImage}); err != nil {
		h.Logger.ErrorContext(ctx, "update product image succeeded but failed to encode response", "user_id", userID, "product_id", productID, "image_id", imageID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "product image updated successfully", "user_id", userID, "product_id", productID, "image_id", imageID, "new_object_key", newObjectKey)
}

func (h *ProductHandler) HandleDeleteProductImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "delete product image failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	productIDStr := r.PathValue("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "delete product image failed: invalid product ID", "user_id", userID, "product_id", productIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	imageIDStr := r.PathValue("imageID")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "delete product image failed: invalid image ID", "user_id", userID, "product_id", productID, "image_id", imageIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid image ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling delete product image request", "user_id", userID, "product_id", productID, "image_id", imageID)

	product, err := h.Queries.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "delete product image failed: product not found", "user_id", userID, "product_id", productID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "delete product image failed: database error fetching product", "user_id", userID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product", err)
		return
	}

	_, err = h.ShopQueries.GetShopByIDAndOwnerID(ctx, database.GetShopByIDAndOwnerIDParams{
		ID:      product.ShopID,
		OwnerID: userID,
	})
	if err != nil {
		h.Logger.WarnContext(ctx, "delete product image failed: user does not own shop", "user_id", userID, "product_id", productID, "shop_id", product.ShopID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not own this product", err)
		return
	}

	image, err := h.Queries.GetProductImage(ctx, database.GetProductImageParams{
		ID:        imageID,
		ProductID: productID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "delete product image failed: image not found", "user_id", userID, "product_id", productID, "image_id", imageID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product image not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "delete product image failed: database error fetching image", "user_id", userID, "product_id", productID, "image_id", imageID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product image", err)
		return
	}

	_, err = h.Queries.DeleteProductImage(ctx, database.DeleteProductImageParams{
		ID:        imageID,
		ProductID: productID,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "delete product image failed: database error", "user_id", userID, "product_id", productID, "image_id", imageID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not delete product image", err)
		return
	}

	if err := h.ImageStorage.Delete(ctx, image.ObjectKey); err != nil {
		h.Logger.ErrorContext(ctx, "delete product image succeeded in DB but failed to delete from storage", "user_id", userID, "product_id", productID, "image_id", imageID, "object_key", image.ObjectKey, "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"message": "Product image deleted successfully"}); err != nil {
		h.Logger.ErrorContext(ctx, "delete product image succeeded but failed to encode response", "user_id", userID, "product_id", productID, "image_id", imageID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "product image deleted successfully", "user_id", userID, "product_id", productID, "image_id", imageID)
}