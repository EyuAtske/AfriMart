package product

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/EyuAtske/AfriMart/backend/internal/auth"
// 	"github.com/EyuAtske/AfriMart/backend/internal/database"
// 	"github.com/google/uuid"
// )

// type mockProductQueries struct {
// 	getShopByIDAndOwnerIDFunc func(
// 		context.Context,
// 		database.GetShopByIDAndOwnerIDParams,
// 	) (database.Shop, error)

// 	createProductFunc func(
// 		context.Context,
// 		database.CreateProductParams,
// 	) (database.Product, error)

// 	getProductFunc func(
// 		context.Context,
// 		uuid.UUID,
// 	) (database.Product, error)

// 	deleteProductFunc func(
// 		context.Context,
// 		uuid.UUID,
// 	) error

// 	updateProductFunc func(
// 		context.Context,
// 		database.UpdateProductParams,
// 	) (database.Product, error)

// 	listProductsFunc func(
// 		context.Context,
// 		database.ListProductsParams,
// 	) ([]database.Product, error)

// 	listProductsByShopFunc func(
// 		context.Context,
// 		database.ListProductsByShopParams,
// 	) ([]database.Product, error)

// 	listProductsByCategoryFunc func(
// 		context.Context,
// 		database.ListProductsByCategoryParams,
// 	) ([]database.Product, error)

// 	listProductsBySubcategoryFunc func(
// 		context.Context,
// 		database.ListProductsBySubcategoryParams,
// 	) ([]database.Product, error)

// 	createProductImageFunc func(
// 		ctx context.Context,
// 		arg database.CreateProductImageParams,
// 	) (database.ProductImage, error)

// 	getProductImagesFunc func(
// 		ctx context.Context,
// 		productID uuid.UUID,
// 	) ([]database.ProductImage, error)

// 	getProductImageFunc func(
// 		ctx context.Context,
// 		arg database.GetProductImageParams,
// 	) (database.ProductImage, error)

// 	updateProductImageFunc func(
// 		ctx context.Context,
// 		arg database.UpdateProductImageParams,
// 	) (database.ProductImage, error)

// 	deleteProductImageFunc func(
// 		ctx context.Context,
// 		arg database.DeleteProductImageParams,
// 	) (database.ProductImage, error)
// }

// func (m *mockProductQueries) GetShopByIDAndOwnerID(
// 	ctx context.Context,
// 	arg database.GetShopByIDAndOwnerIDParams,
// ) (database.Shop, error) {
// 	return m.getShopByIDAndOwnerIDFunc(ctx, arg)
// }

// func (m *mockProductQueries) CreateProduct(
// 	ctx context.Context,
// 	arg database.CreateProductParams,
// ) (database.Product, error) {
// 	return m.createProductFunc(ctx, arg)
// }

// func (m *mockProductQueries) GetProduct(
// 	ctx context.Context,
// 	id uuid.UUID,
// ) (database.Product, error) {
// 	return m.getProductFunc(ctx, id)
// }

// func (m *mockProductQueries) DeleteProduct(
// 	ctx context.Context,
// 	id uuid.UUID,
// ) error {
// 	return m.deleteProductFunc(ctx, id)
// }

// func (m *mockProductQueries) UpdateProduct(
// 	ctx context.Context,
// 	arg database.UpdateProductParams,
// ) (database.Product, error) {
// 	return m.updateProductFunc(ctx, arg)
// }

// func (m *mockProductQueries) ListProducts(
// 	ctx context.Context,
// 	arg database.ListProductsParams,
// ) ([]database.Product, error) {
// 	return m.listProductsFunc(ctx, arg)
// }

// func (m *mockProductQueries) ListProductsByShop(
// 	ctx context.Context,
// 	arg database.ListProductsByShopParams,
// ) ([]database.Product, error) {
// 	return m.listProductsByShopFunc(ctx, arg)
// }

// func (m *mockProductQueries) ListProductsByCategory(
// 	ctx context.Context,
// 	arg database.ListProductsByCategoryParams,
// ) ([]database.Product, error) {
// 	return m.listProductsByCategoryFunc(ctx, arg)
// }

// func (m *mockProductQueries) ListProductsBySubcategory(
// 	ctx context.Context,
// 	arg database.ListProductsBySubcategoryParams,
// ) ([]database.Product, error) {
// 	return m.listProductsBySubcategoryFunc(ctx, arg)
// }

// func (m *mockProductQueries) CreateProductImage(
// 		ctx context.Context,
// 		arg database.CreateProductImageParams,
// 	) (database.ProductImage, error){
// 	return m.createProductImageFunc(ctx, arg)
// }

// func (m *mockProductQueries) GetProductImages(
// 		ctx context.Context,
// 		productID uuid.UUID,
// 	) ([]database.ProductImage, error){
// 	return m.getProductImagesFunc(ctx, productID)
// }

// func (m *mockProductQueries) GetProductImage(
// 		ctx context.Context,
// 		arg database.GetProductImageParams,
// 	) (database.ProductImage, error){
// 	return m.getProductImageFunc(ctx, arg)
// }

// func (m *mockProductQueries) UpdateProductImage(
// 		ctx context.Context,
// 		arg database.UpdateProductImageParams,
// 	) (database.ProductImage, error){
// 	return m.updateProductImageFunc(ctx, arg)
// }

// func (m *mockProductQueries)	DeleteProductImage(
// 		ctx context.Context,
// 		arg database.DeleteProductImageParams,
// 	) (database.ProductImage, error){
// 	return m.deleteProductImageFunc(ctx, arg)
// }

// func TestHandleCreateProduct_Success(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()
// 	categoryID := uuid.New()
// 	subcategoryID := uuid.New()
// 	productID := uuid.New()

// 	expectedProduct := database.Product{
// 		ID:            productID,
// 		ShopID:        shopID,
// 		CategoryID:    categoryID,
// 		SubcategoryID: subcategoryID,
// 		Name:          "Nike Air Max",
// 		Price:         "120.00",
// 		Stock:         10,
// 		Status:        "active",
// 	}

// 	mockQueries := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			if arg.ID != shopID {
// 				t.Errorf("expected shop ID %v, got %v", shopID, arg.ID)
// 			}

// 			if arg.OwnerID != userID {
// 				t.Errorf("expected owner ID %v, got %v", userID, arg.OwnerID)
// 			}

// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 				Name:    "My Shop",
// 			}, nil
// 		},

// 		createProductFunc: func(
// 			ctx context.Context,
// 			arg database.CreateProductParams,
// 		) (database.Product, error) {
// 			if arg.ShopID != shopID {
// 				t.Errorf("expected shop ID %v, got %v", shopID, arg.ShopID)
// 			}

// 			if arg.CategoryID != categoryID {
// 				t.Errorf("expected category ID %v, got %v", categoryID, arg.CategoryID)
// 			}

// 			if arg.SubcategoryID != subcategoryID {
// 				t.Errorf(
// 					"expected subcategory ID %v, got %v",
// 					subcategoryID,
// 					arg.SubcategoryID,
// 				)
// 			}

// 			if arg.Name != "Nike Air Max" {
// 				t.Errorf("expected product name Nike Air Max, got %v", arg.Name)
// 			}

// 			if arg.Price != "120.00" {
// 				t.Errorf("expected price 120.00, got %v", arg.Price)
// 			}

// 			if arg.Stock != 10 {
// 				t.Errorf("expected stock 10, got %v", arg.Stock)
// 			}

// 			if arg.Status != "active" {
// 				t.Errorf("expected status active, got %v", arg.Status)
// 			}

// 			return expectedProduct, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	body := `{
// 		"shop_id": "` + shopID.String() + `",
// 		"category_id": "` + categoryID.String() + `",
// 		"subcategory_id": "` + subcategoryID.String() + `",
// 		"name": "Nike Air Max",
// 		"description": "Running shoes",
// 		"brand": "Nike",
// 		"color": "Black",
// 		"size": "42",
// 		"price": "120.00",
// 		"stock": 10,
// 		"image": "https://example.com/nike.jpg"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPost,
// 		"/products",
// 		strings.NewReader(body),
// 	)

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	recorder := httptest.NewRecorder()

// 	handler.HandleCreateProduct(recorder, req)

// 	if recorder.Code != http.StatusCreated {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusCreated,
// 			recorder.Code,
// 		)
// 	}
// }

// func TestHandleCreateProduct_ShopNotOwned(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()
// 	categoryID := uuid.New()
// 	subcategoryID := uuid.New()

// 	createProductCalled := false

// 	mockQueries := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{}, sql.ErrNoRows
// 		},

// 		createProductFunc: func(
// 			ctx context.Context,
// 			arg database.CreateProductParams,
// 		) (database.Product, error) {
// 			createProductCalled = true
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	body := `{
// 		"shop_id": "` + shopID.String() + `",
// 		"category_id": "` + categoryID.String() + `",
// 		"subcategory_id": "` + subcategoryID.String() + `",
// 		"name": "Nike Air Max",
// 		"price": "120.00",
// 		"stock": 10
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPost,
// 		"/api/products",
// 		strings.NewReader(body),
// 	)

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	recorder := httptest.NewRecorder()

// 	handler.HandleCreateProduct(recorder, req)

// 	if recorder.Code != http.StatusForbidden {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusForbidden,
// 			recorder.Code,
// 		)
// 	}

// 	if createProductCalled {
// 		t.Fatal("CreateProduct should not be called when user does not own the shop")
// 	}
// }

// func TestHandleCreateProduct_Unauthenticated(t *testing.T) {
// 	mockQueries := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	body := `{
//         "shop_id": "` + uuid.New().String() + `",
//         "category_id": "` + uuid.New().String() + `",
//         "subcategory_id": "` + uuid.New().String() + `",
//         "name": "Nike Air Max",
//         "price": "120.00",
//         "stock": 10
//     }`

// 	req := httptest.NewRequest(
// 		http.MethodPost,
// 		"/products",
// 		strings.NewReader(body),
// 	)

// 	recorder := httptest.NewRecorder()

// 	handler.HandleCreateProduct(recorder, req)

// 	if recorder.Code != http.StatusUnauthorized {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusUnauthorized,
// 			recorder.Code,
// 		)
// 	}
// }

// func TestHandleCreateProduct_InvalidPrice(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()
// 	categoryID := uuid.New()
// 	subcategoryID := uuid.New()

// 	createProductCalled := false

// 	mockQueries := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 				Name:    "My Shop",
// 			}, nil
// 		},

// 		createProductFunc: func(
// 			ctx context.Context,
// 			arg database.CreateProductParams,
// 		) (database.Product, error) {
// 			createProductCalled = true
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	body := `{
// 		"shop_id": "` + shopID.String() + `",
// 		"category_id": "` + categoryID.String() + `",
// 		"subcategory_id": "` + subcategoryID.String() + `",
// 		"name": "Nike Air Max",
// 		"price": "abc",
// 		"stock": 10
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPost,
// 		"/products",
// 		strings.NewReader(body),
// 	)

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	recorder := httptest.NewRecorder()

// 	handler.HandleCreateProduct(recorder, req)

// 	if recorder.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			recorder.Code,
// 		)
// 	}

// 	if createProductCalled {
// 		t.Fatal("CreateProduct should not be called when price is invalid")
// 	}
// }

// func TestHandleCreateProduct_MissingName(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()
// 	categoryID := uuid.New()
// 	subcategoryID := uuid.New()

// 	createProductCalled := false

// 	mockQueries := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 				Name:    "My Shop",
// 			}, nil
// 		},

// 		createProductFunc: func(
// 			ctx context.Context,
// 			arg database.CreateProductParams,
// 		) (database.Product, error) {
// 			createProductCalled = true
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	body := `{
// 		"shop_id": "` + shopID.String() + `",
// 		"category_id": "` + categoryID.String() + `",
// 		"subcategory_id": "` + subcategoryID.String() + `",
// 		"price": "120.00",
// 		"stock": 10
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPost,
// 		"/products",
// 		strings.NewReader(body),
// 	)

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	recorder := httptest.NewRecorder()

// 	handler.HandleCreateProduct(recorder, req)

// 	if recorder.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			recorder.Code,
// 		)
// 	}

// 	if createProductCalled {
// 		t.Fatal("CreateProduct should not be called when name is missing")
// 	}
// }

// func TestHandleCreateProduct_NegativeStock(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()
// 	categoryID := uuid.New()
// 	subcategoryID := uuid.New()

// 	createProductCalled := false

// 	mockQueries := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 				Name:    "My Shop",
// 			}, nil
// 		},

// 		createProductFunc: func(
// 			ctx context.Context,
// 			arg database.CreateProductParams,
// 		) (database.Product, error) {
// 			createProductCalled = true
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	body := `{
// 		"shop_id": "` + shopID.String() + `",
// 		"category_id": "` + categoryID.String() + `",
// 		"subcategory_id": "` + subcategoryID.String() + `",
// 		"name": "Nike Air Max",
// 		"price": "120.00",
// 		"stock": -5
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPost,
// 		"/products",
// 		strings.NewReader(body),
// 	)

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	recorder := httptest.NewRecorder()

// 	handler.HandleCreateProduct(recorder, req)

// 	if recorder.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			recorder.Code,
// 		)
// 	}

// 	if createProductCalled {
// 		t.Fatal("CreateProduct should not be called when stock is negative")
// 	}
// }

// func TestHandleGetProduct_Success(t *testing.T) {
// 	productID := uuid.New()
// 	shopID := uuid.New()
// 	categoryID := uuid.New()
// 	subcategoryID := uuid.New()

// 	expectedProduct := database.Product{
// 		ID:            productID,
// 		ShopID:        shopID,
// 		CategoryID:    categoryID,
// 		SubcategoryID: subcategoryID,
// 		Name:          "Nike Air Max",
// 		Price:         "120.00",
// 		Stock:         10,
// 		Status:        "active",
// 	}

// 	mockQueries := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			if id != productID {
// 				t.Fatalf(
// 					"expected product ID %v, got %v",
// 					productID,
// 					id,
// 				)
// 			}

// 			return expectedProduct, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	recorder := httptest.NewRecorder()

// 	handler.HandleGetProduct(recorder, req)

// 	if recorder.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			recorder.Code,
// 		)
// 	}

// 	var gotProduct database.Product

// 	if err := json.NewDecoder(recorder.Body).Decode(&gotProduct); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if gotProduct.ID != expectedProduct.ID {
// 		t.Fatalf(
// 			"expected product ID %v, got %v",
// 			expectedProduct.ID,
// 			gotProduct.ID,
// 		)
// 	}

// 	if gotProduct.Name != expectedProduct.Name {
// 		t.Fatalf(
// 			"expected product name %q, got %q",
// 			expectedProduct.Name,
// 			gotProduct.Name,
// 		)
// 	}

// 	if gotProduct.Price != expectedProduct.Price {
// 		t.Fatalf(
// 			"expected price %q, got %q",
// 			expectedProduct.Price,
// 			gotProduct.Price,
// 		)
// 	}

// 	if gotProduct.Stock != expectedProduct.Stock {
// 		t.Fatalf(
// 			"expected stock %d, got %d",
// 			expectedProduct.Stock,
// 			gotProduct.Stock,
// 		)
// 	}
// }

// func TestHandleGetProduct_NotFound(t *testing.T) {
// 	productID := uuid.New()

// 	mockQueries := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			if id != productID {
// 				t.Fatalf(
// 					"expected product ID %v, got %v",
// 					productID,
// 					id,
// 				)
// 			}

// 			return database.Product{}, sql.ErrNoRows
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	recorder := httptest.NewRecorder()

// 	handler.HandleGetProduct(recorder, req)

// 	if recorder.Code != http.StatusNotFound {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusNotFound,
// 			recorder.Code,
// 		)
// 	}
// }

// func TestHandleGetProduct_InvalidID(t *testing.T) {
// 	getProductCalled := false

// 	mockQueries := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			getProductCalled = true
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mockQueries,
// 		ShopQueries: mockQueries,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products/not-a-uuid",
// 		nil,
// 	)

// 	req.SetPathValue("id", "not-a-uuid")

// 	recorder := httptest.NewRecorder()

// 	handler.HandleGetProduct(recorder, req)

// 	if recorder.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			recorder.Code,
// 		)
// 	}

// 	if getProductCalled {
// 		t.Fatal("GetProduct should not be called with an invalid product ID")
// 	}
// }

// func TestHandleUpdateProduct_Success(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()
// 	categoryID := uuid.New()
// 	subcategoryID := uuid.New()

// 	existingProduct := database.Product{
// 		ID:            productID,
// 		ShopID:        shopID,
// 		CategoryID:    categoryID,
// 		SubcategoryID: subcategoryID,
// 		Name:          "Old Product",
// 		Price:         "100.00",
// 		Stock:         10,
// 		Status:        "active",
// 	}

// 	updatedProduct := existingProduct
// 	updatedProduct.Name = "Updated Product"
// 	updatedProduct.Price = "150.00"
// 	updatedProduct.Stock = 20

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			if id != productID {
// 				t.Errorf("expected product ID %v, got %v", productID, id)
// 			}

// 			return existingProduct, nil
// 		},

// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			if arg.ID != shopID {
// 				t.Errorf("expected shop ID %v, got %v", shopID, arg.ID)
// 			}

// 			if arg.OwnerID != userID {
// 				t.Errorf("expected owner ID %v, got %v", userID, arg.OwnerID)
// 			}

// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},

// 		updateProductFunc: func(
// 			ctx context.Context,
// 			arg database.UpdateProductParams,
// 		) (database.Product, error) {
// 			if arg.ID != productID {
// 				t.Errorf("expected product ID %v, got %v", productID, arg.ID)
// 			}

// 			if arg.Name != "Updated Product" {
// 				t.Errorf("expected updated name %q, got %q", "Updated Product", arg.Name)
// 			}

// 			if arg.Price != "150.00" {
// 				t.Errorf("expected price %q, got %q", "150.00", arg.Price)
// 			}

// 			if arg.Stock != 20 {
// 				t.Errorf("expected stock %d, got %d", 20, arg.Stock)
// 			}

// 			return updatedProduct, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	body := `{
// 		"category_id": "` + categoryID.String() + `",
// 		"subcategory_id": "` + subcategoryID.String() + `",
// 		"name": "Updated Product",
// 		"description": "Updated description",
// 		"brand": "Updated Brand",
// 		"color": "Black",
// 		"size": "L",
// 		"price": "150.00",
// 		"stock": 20,
// 		"image": "updated-image.jpg",
// 		"status": "active"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPut,
// 		"/products/"+productID.String(),
// 		strings.NewReader(body),
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleUpdateProduct(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}

// 	var response database.Product

// 	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if response.Name != "Updated Product" {
// 		t.Errorf(
// 			"expected response name %q, got %q",
// 			"Updated Product",
// 			response.Name,
// 		)
// 	}

// 	if response.Price != "150.00" {
// 		t.Errorf(
// 			"expected response price %q, got %q",
// 			"150.00",
// 			response.Price,
// 		)
// 	}
// }

// func TestHandleUpdateProduct_ShopNotOwned(t *testing.T) {
// 	userID := uuid.New()
// 	ownerID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()

// 	existingProduct := database.Product{
// 		ID:     productID,
// 		ShopID: shopID,
// 		Name:   "Existing Product",
// 		Price:  "100.00",
// 		Stock:  10,
// 		Status: "active",
// 	}

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return existingProduct, nil
// 		},

// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			if arg.ID != shopID {
// 				t.Errorf("expected shop ID %v, got %v", shopID, arg.ID)
// 			}

// 			if arg.OwnerID != userID {
// 				t.Errorf("expected owner ID %v, got %v", userID, arg.OwnerID)
// 			}

// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: ownerID,
// 			}, sql.ErrNoRows
// 		},

// 		updateProductFunc: func(
// 			ctx context.Context,
// 			arg database.UpdateProductParams,
// 		) (database.Product, error) {
// 			t.Fatal("UpdateProduct should not be called when shop is not owned")
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	body := `{
// 		"category_id": "` + uuid.New().String() + `",
// 		"subcategory_id": "` + uuid.New().String() + `",
// 		"name": "Attempted Update",
// 		"price": "150.00",
// 		"stock": 20,
// 		"status": "active"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPut,
// 		"/products/"+productID.String(),
// 		strings.NewReader(body),
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleUpdateProduct(rec, req)

// 	if rec.Code != http.StatusForbidden {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusForbidden,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleUpdateProduct_Unauthenticated(t *testing.T) {
// 	productID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	body := `{
// 		"category_id": "` + uuid.New().String() + `",
// 		"subcategory_id": "` + uuid.New().String() + `",
// 		"name": "Updated Product",
// 		"price": "150.00",
// 		"stock": 20,
// 		"status": "active"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPut,
// 		"/products/"+productID.String(),
// 		strings.NewReader(body),
// 	)

// 	req.SetPathValue("id", productID.String())

// 	// Deliberately do NOT add a user ID to the context.
// 	rec := httptest.NewRecorder()

// 	handler.HandleUpdateProduct(rec, req)

// 	if rec.Code != http.StatusUnauthorized {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusUnauthorized,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleUpdateProduct_InvalidID(t *testing.T) {
// 	userID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	body := `{
// 		"name": "Updated Product",
// 		"price": "150.00",
// 		"stock": 20,
// 		"status": "active"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPut,
// 		"/products/not-a-uuid",
// 		strings.NewReader(body),
// 	)

// 	req.SetPathValue("id", "not-a-uuid")

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleUpdateProduct(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleUpdateProduct_NotFound(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return database.Product{}, sql.ErrNoRows
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	body := `{
// 		"name": "Updated Product",
// 		"price": "150.00",
// 		"stock": 20,
// 		"status": "active"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPut,
// 		"/products/"+productID.String(),
// 		strings.NewReader(body),
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleUpdateProduct(rec, req)

// 	if rec.Code != http.StatusNotFound {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusNotFound,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleUpdateProduct_InvalidPrice(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()

// 	existingProduct := database.Product{
// 		ID:     productID,
// 		ShopID: shopID,
// 		Name:   "Existing Product",
// 		Price:  "100.00",
// 		Stock:  10,
// 		Status: "active",
// 	}

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return existingProduct, nil
// 		},

// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},

// 		updateProductFunc: func(
// 			ctx context.Context,
// 			arg database.UpdateProductParams,
// 		) (database.Product, error) {
// 			t.Fatal("UpdateProduct should not be called with invalid price")
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	body := `{
// 		"category_id": "` + uuid.New().String() + `",
// 		"subcategory_id": "` + uuid.New().String() + `",
// 		"name": "Updated Product",
// 		"price": "invalid-price",
// 		"stock": 20,
// 		"status": "active"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPut,
// 		"/products/"+productID.String(),
// 		strings.NewReader(body),
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleUpdateProduct(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleUpdateProduct_MissingName(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()

// 	existingProduct := database.Product{
// 		ID:     productID,
// 		ShopID: shopID,
// 		Name:   "Existing Product",
// 		Price:  "100.00",
// 		Stock:  10,
// 		Status: "active",
// 	}

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return existingProduct, nil
// 		},

// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},

// 		updateProductFunc: func(
// 			ctx context.Context,
// 			arg database.UpdateProductParams,
// 		) (database.Product, error) {
// 			t.Fatal("UpdateProduct should not be called when name is missing")
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	body := `{
// 		"category_id": "` + uuid.New().String() + `",
// 		"subcategory_id": "` + uuid.New().String() + `",
// 		"name": "",
// 		"price": "150.00",
// 		"stock": 20,
// 		"status": "active"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPut,
// 		"/products/"+productID.String(),
// 		strings.NewReader(body),
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleUpdateProduct(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleUpdateProduct_NegativeStock(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()

// 	existingProduct := database.Product{
// 		ID:     productID,
// 		ShopID: shopID,
// 		Name:   "Existing Product",
// 		Price:  "100.00",
// 		Stock:  10,
// 		Status: "active",
// 	}

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return existingProduct, nil
// 		},

// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},

// 		updateProductFunc: func(
// 			ctx context.Context,
// 			arg database.UpdateProductParams,
// 		) (database.Product, error) {
// 			t.Fatal("UpdateProduct should not be called with negative stock")
// 			return database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	body := `{
// 		"category_id": "` + uuid.New().String() + `",
// 		"subcategory_id": "` + uuid.New().String() + `",
// 		"name": "Updated Product",
// 		"price": "150.00",
// 		"stock": -5,
// 		"status": "active"
// 	}`

// 	req := httptest.NewRequest(
// 		http.MethodPut,
// 		"/products/"+productID.String(),
// 		strings.NewReader(body),
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleUpdateProduct(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleDeleteProduct_Success(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()

// 	existingProduct := database.Product{
// 		ID:     productID,
// 		ShopID: shopID,
// 		Name:   "Product To Delete",
// 		Price:  "100.00",
// 		Stock:  10,
// 		Status: "active",
// 	}

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			if id != productID {
// 				t.Errorf("expected product ID %v, got %v", productID, id)
// 			}

// 			return existingProduct, nil
// 		},

// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			if arg.ID != shopID {
// 				t.Errorf("expected shop ID %v, got %v", shopID, arg.ID)
// 			}

// 			if arg.OwnerID != userID {
// 				t.Errorf("expected owner ID %v, got %v", userID, arg.OwnerID)
// 			}

// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},

// 		deleteProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) error {
// 			if id != productID {
// 				t.Errorf("expected product ID %v, got %v", productID, id)
// 			}

// 			return nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodDelete,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleDeleteProduct(rec, req)

// 	if rec.Code != http.StatusNoContent {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusNoContent,
// 			rec.Code,
// 		)
// 	}

// 	if rec.Body.Len() != 0 {
// 		t.Errorf("expected empty response body, got %q", rec.Body.String())
// 	}
// }

// func TestHandleDeleteProduct_ShopNotOwned(t *testing.T) {
// 	userID := uuid.New()
// 	ownerID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()

// 	existingProduct := database.Product{
// 		ID:     productID,
// 		ShopID: shopID,
// 		Name:   "Product To Delete",
// 		Price:  "100.00",
// 		Stock:  10,
// 		Status: "active",
// 	}

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return existingProduct, nil
// 		},

// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			if arg.ID != shopID {
// 				t.Errorf("expected shop ID %v, got %v", shopID, arg.ID)
// 			}

// 			if arg.OwnerID != userID {
// 				t.Errorf("expected owner ID %v, got %v", userID, arg.OwnerID)
// 			}

// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: ownerID,
// 			}, sql.ErrNoRows
// 		},

// 		deleteProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) error {
// 			t.Fatal("DeleteProduct should not be called when shop is not owned")
// 			return nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodDelete,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleDeleteProduct(rec, req)

// 	if rec.Code != http.StatusForbidden {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusForbidden,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleDeleteProduct_Unauthenticated(t *testing.T) {
// 	productID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodDelete,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleDeleteProduct(rec, req)

// 	if rec.Code != http.StatusUnauthorized {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusUnauthorized,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleDeleteProduct_InvalidID(t *testing.T) {
// 	userID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodDelete,
// 		"/products/not-a-uuid",
// 		nil,
// 	)

// 	req.SetPathValue("id", "not-a-uuid")

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleDeleteProduct(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleDeleteProduct_NotFound(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return database.Product{}, sql.ErrNoRows
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodDelete,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleDeleteProduct(rec, req)

// 	if rec.Code != http.StatusNotFound {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusNotFound,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleDeleteProduct_GetProductError(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return database.Product{}, errors.New("database error")
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodDelete,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleDeleteProduct(rec, req)

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleDeleteProduct_ShopOwnershipError(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return database.Product{
// 				ID:     productID,
// 				ShopID: shopID,
// 			}, nil
// 		},
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{}, errors.New("database error")
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodDelete,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleDeleteProduct(rec, req)

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleDeleteProduct_DeleteError(t *testing.T) {
// 	userID := uuid.New()
// 	productID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) (database.Product, error) {
// 			return database.Product{
// 				ID:     productID,
// 				ShopID: shopID,
// 			}, nil
// 		},
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 		deleteProductFunc: func(
// 			ctx context.Context,
// 			id uuid.UUID,
// 		) error {
// 			return errors.New("database error")
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodDelete,
// 		"/products/"+productID.String(),
// 		nil,
// 	)

// 	req.SetPathValue("id", productID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleDeleteProduct(rec, req)

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProducts_Success(t *testing.T) {
// 	product1 := database.Product{
// 		ID:   uuid.New(),
// 		Name: "T-Shirt",
// 	}

// 	product2 := database.Product{
// 		ID:   uuid.New(),
// 		Name: "Jeans",
// 	}

// 	mock := &mockProductQueries{
// 		listProductsFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsParams,
// 		) ([]database.Product, error) {
// 			if arg.PageLimit != 20 {
// 				t.Fatalf("expected limit 20, got %d", arg.PageLimit)
// 			}

// 			if arg.PageOffset != 0 {
// 				t.Fatalf("expected offset 0, got %d", arg.PageOffset)
// 			}

// 			return []database.Product{
// 				product1,
// 				product2,
// 			}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}

// 	var products []database.Product

// 	if err := json.NewDecoder(rec.Body).Decode(&products); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if len(products) != 2 {
// 		t.Fatalf(
// 			"expected 2 products, got %d",
// 			len(products),
// 		)
// 	}

// 	if products[0].ID != product1.ID {
// 		t.Fatalf("expected first product %s, got %s", product1.ID, products[0].ID)
// 	}

// 	if products[1].ID != product2.ID {
// 		t.Fatalf("expected second product %s, got %s", product2.ID, products[1].ID)
// 	}
// }

// func TestHandleListProducts_CustomPagination(t *testing.T) {
// 	mock := &mockProductQueries{
// 		listProductsFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsParams,
// 		) ([]database.Product, error) {
// 			if arg.PageLimit != 10 {
// 				t.Fatalf("expected limit 10, got %d", arg.PageLimit)
// 			}

// 			if arg.PageOffset != 20 {
// 				t.Fatalf("expected offset 20, got %d", arg.PageOffset)
// 			}

// 			return []database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products?limit=10&offset=20",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProducts_InvalidLimit(t *testing.T) {
// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products?limit=abc",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProducts_ZeroLimit(t *testing.T) {
// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products?limit=0",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProducts_NegativeLimit(t *testing.T) {
// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products?limit=-1",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProducts_InvalidOffset(t *testing.T) {
// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products?offset=abc",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProducts_NegativeOffset(t *testing.T) {
// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products?offset=-1",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProducts_DatabaseError(t *testing.T) {
// 	mock := &mockProductQueries{
// 		listProductsFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsParams,
// 		) ([]database.Product, error) {
// 			return nil, errors.New("database error")
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProducts_Empty(t *testing.T) {
// 	mock := &mockProductQueries{
// 		listProductsFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsParams,
// 		) ([]database.Product, error) {
// 			return []database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/products",
// 		nil,
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProducts(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}

// 	var products []database.Product

// 	if err := json.NewDecoder(rec.Body).Decode(&products); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if len(products) != 0 {
// 		t.Fatalf(
// 			"expected 0 products, got %d",
// 			len(products),
// 		)
// 	}
// }

// func TestHandleListProductsByShop_Success(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	product1 := database.Product{
// 		ID:     uuid.New(),
// 		ShopID: shopID,
// 		Name:   "T-Shirt",
// 	}

// 	product2 := database.Product{
// 		ID:     uuid.New(),
// 		ShopID: shopID,
// 		Name:   "Jeans",
// 	}

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			if arg.ID != shopID {
// 				t.Fatalf("expected shop ID %s, got %s", shopID, arg.ID)
// 			}

// 			if arg.OwnerID != userID {
// 				t.Fatalf("expected owner ID %s, got %s", userID, arg.OwnerID)
// 			}

// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 		listProductsByShopFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsByShopParams,
// 		) ([]database.Product, error) {
// 			if arg.ShopID != shopID {
// 				t.Fatalf("expected shop ID %s, got %s", shopID, arg.ShopID)
// 			}

// 			if arg.Limit != 20 {
// 				t.Fatalf("expected limit 20, got %d", arg.Limit)
// 			}

// 			if arg.Offset != 0 {
// 				t.Fatalf("expected offset 0, got %d", arg.Offset)
// 			}

// 			return []database.Product{
// 				product1,
// 				product2,
// 			}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}

// 	var products []database.Product

// 	if err := json.NewDecoder(rec.Body).Decode(&products); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if len(products) != 2 {
// 		t.Fatalf(
// 			"expected 2 products, got %d",
// 			len(products),
// 		)
// 	}
// }

// func TestHandleListProductsByShop_Unauthenticated(t *testing.T) {
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusUnauthorized {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusUnauthorized,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_InvalidShopID(t *testing.T) {
// 	userID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/not-a-uuid/products",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", "not-a-uuid")

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_ShopNotOwned(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{}, sql.ErrNoRows
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusForbidden {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusForbidden,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_ShopOwnershipError(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{}, errors.New("database error")
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_DatabaseError(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 		listProductsByShopFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsByShopParams,
// 		) ([]database.Product, error) {
// 			return nil, errors.New("database error")
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_CustomPagination(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	var receivedArgs database.ListProductsByShopParams

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 		listProductsByShopFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsByShopParams,
// 		) ([]database.Product, error) {
// 			receivedArgs = arg
// 			return []database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products?limit=10&offset=20",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}

// 	if receivedArgs.ShopID != shopID {
// 		t.Fatalf("expected shop ID %v, got %v", shopID, receivedArgs.ShopID)
// 	}

// 	if receivedArgs.Limit != 10 {
// 		t.Fatalf("expected limit 10, got %d", receivedArgs.Limit)
// 	}

// 	if receivedArgs.Offset != 20 {
// 		t.Fatalf("expected offset 20, got %d", receivedArgs.Offset)
// 	}
// }

// func TestHandleListProductsByShop_InvalidLimit(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products?limit=abc",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_ZeroLimit(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products?limit=0",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_NegativeLimit(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products?limit=-1",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_InvalidOffset(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products?offset=abc",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_NegativeOffset(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products?offset=-1",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByShop_Empty(t *testing.T) {
// 	userID := uuid.New()
// 	shopID := uuid.New()

// 	mock := &mockProductQueries{
// 		getShopByIDAndOwnerIDFunc: func(
// 			ctx context.Context,
// 			arg database.GetShopByIDAndOwnerIDParams,
// 		) (database.Shop, error) {
// 			return database.Shop{
// 				ID:      shopID,
// 				OwnerID: userID,
// 			}, nil
// 		},
// 		listProductsByShopFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsByShopParams,
// 		) ([]database.Product, error) {
// 			return []database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries:     mock,
// 		ShopQueries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/shops/"+shopID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("shop_id", shopID.String())

// 	req = req.WithContext(
// 		auth.ContextWithUserID(req.Context(), userID),
// 	)

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByShop(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_Success(t *testing.T) {
// 	categoryID := uuid.New()

// 	expectedProducts := []database.Product{
// 		{
// 			ID:         uuid.New(),
// 			CategoryID: categoryID,
// 			Name:       "Test Product",
// 			Status:     "active",
// 		},
// 	}

// 	mock := &mockProductQueries{
// 		listProductsByCategoryFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsByCategoryParams,
// 		) ([]database.Product, error) {
// 			if arg.CategoryID != categoryID {
// 				t.Fatalf("expected category ID %v, got %v", categoryID, arg.CategoryID)
// 			}

// 			if arg.Limit != 20 {
// 				t.Fatalf("expected default limit 20, got %d", arg.Limit)
// 			}

// 			if arg.Offset != 0 {
// 				t.Fatalf("expected default offset 0, got %d", arg.Offset)
// 			}

// 			return expectedProducts, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_InvalidID(t *testing.T) {
// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/invalid/products",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", "invalid")

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_DatabaseError(t *testing.T) {
// 	categoryID := uuid.New()

// 	mock := &mockProductQueries{
// 		listProductsByCategoryFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsByCategoryParams,
// 		) ([]database.Product, error) {
// 			return nil, errors.New("database error")
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_CustomPagination(t *testing.T) {
// 	categoryID := uuid.New()

// 	var receivedArgs database.ListProductsByCategoryParams

// 	mock := &mockProductQueries{
// 		listProductsByCategoryFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsByCategoryParams,
// 		) ([]database.Product, error) {
// 			receivedArgs = arg
// 			return []database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products?limit=10&offset=20",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}

// 	if receivedArgs.CategoryID != categoryID {
// 		t.Fatalf(
// 			"expected category ID %v, got %v",
// 			categoryID,
// 			receivedArgs.CategoryID,
// 		)
// 	}

// 	if receivedArgs.Limit != 10 {
// 		t.Fatalf("expected limit 10, got %d", receivedArgs.Limit)
// 	}

// 	if receivedArgs.Offset != 20 {
// 		t.Fatalf("expected offset 20, got %d", receivedArgs.Offset)
// 	}
// }

// func TestHandleListProductsByCategory_InvalidLimit(t *testing.T) {
// 	categoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products?limit=abc",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_ZeroLimit(t *testing.T) {
// 	categoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products?limit=0",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_NegativeLimit(t *testing.T) {
// 	categoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products?limit=-1",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_InvalidOffset(t *testing.T) {
// 	categoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products?offset=abc",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_NegativeOffset(t *testing.T) {
// 	categoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products?offset=-1",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsByCategory_Empty(t *testing.T) {
// 	categoryID := uuid.New()

// 	mock := &mockProductQueries{
// 		listProductsByCategoryFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsByCategoryParams,
// 		) ([]database.Product, error) {
// 			return []database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/categories/"+categoryID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("category_id", categoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsByCategory(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsBySubcategory_Success(t *testing.T) {
// 	subcategoryID := uuid.New()

// 	expectedProducts := []database.Product{
// 		{
// 			ID:            uuid.New(),
// 			SubcategoryID: subcategoryID,
// 			Name:          "Test Product",
// 			Status:        "active",
// 		},
// 	}

// 	mock := &mockProductQueries{
// 		listProductsBySubcategoryFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsBySubcategoryParams,
// 		) ([]database.Product, error) {
// 			if arg.SubcategoryID != subcategoryID {
// 				t.Fatalf(
// 					"expected subcategory ID %v, got %v",
// 					subcategoryID,
// 					arg.SubcategoryID,
// 				)
// 			}

// 			if arg.Limit != 20 {
// 				t.Fatalf("expected default limit 20, got %d", arg.Limit)
// 			}

// 			if arg.Offset != 0 {
// 				t.Fatalf("expected default offset 0, got %d", arg.Offset)
// 			}

// 			return expectedProducts, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/subcategories/"+subcategoryID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("subcategory_id", subcategoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsBySubcategory(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsBySubcategory_InvalidID(t *testing.T) {
// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/subcategories/invalid/products",
// 		nil,
// 	)

// 	req.SetPathValue("subcategory_id", "invalid")

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsBySubcategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsBySubcategory_DatabaseError(t *testing.T) {
// 	subcategoryID := uuid.New()

// 	mock := &mockProductQueries{
// 		listProductsBySubcategoryFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsBySubcategoryParams,
// 		) ([]database.Product, error) {
// 			return nil, errors.New("database error")
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/subcategories/"+subcategoryID.String()+"/products",
// 		nil,
// 	)

// 	req.SetPathValue("subcategory_id", subcategoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsBySubcategory(rec, req)

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsBySubcategory_CustomPagination(t *testing.T) {
// 	subcategoryID := uuid.New()

// 	var receivedArgs database.ListProductsBySubcategoryParams

// 	mock := &mockProductQueries{
// 		listProductsBySubcategoryFunc: func(
// 			ctx context.Context,
// 			arg database.ListProductsBySubcategoryParams,
// 		) ([]database.Product, error) {
// 			receivedArgs = arg
// 			return []database.Product{}, nil
// 		},
// 	}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/subcategories/"+subcategoryID.String()+"/products?limit=10&offset=20",
// 		nil,
// 	)

// 	req.SetPathValue("subcategory_id", subcategoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsBySubcategory(rec, req)

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusOK,
// 			rec.Code,
// 		)
// 	}

// 	if receivedArgs.SubcategoryID != subcategoryID {
// 		t.Fatalf(
// 			"expected subcategory ID %v, got %v",
// 			subcategoryID,
// 			receivedArgs.SubcategoryID,
// 		)
// 	}

// 	if receivedArgs.Limit != 10 {
// 		t.Fatalf("expected limit 10, got %d", receivedArgs.Limit)
// 	}

// 	if receivedArgs.Offset != 20 {
// 		t.Fatalf("expected offset 20, got %d", receivedArgs.Offset)
// 	}
// }

// func TestHandleListProductsBySubcategory_InvalidLimit(t *testing.T) {
// 	subcategoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/subcategories/"+subcategoryID.String()+"/products?limit=abc",
// 		nil,
// 	)

// 	req.SetPathValue("subcategory_id", subcategoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsBySubcategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsBySubcategory_ZeroLimit(t *testing.T) {
// 	subcategoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/subcategories/"+subcategoryID.String()+"/products?limit=0",
// 		nil,
// 	)

// 	req.SetPathValue("subcategory_id", subcategoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsBySubcategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsBySubcategory_NegativeLimit(t *testing.T) {
// 	subcategoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/subcategories/"+subcategoryID.String()+"/products?limit=-1",
// 		nil,
// 	)

// 	req.SetPathValue("subcategory_id", subcategoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsBySubcategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }

// func TestHandleListProductsBySubcategory_InvalidOffset(t *testing.T) {
// 	subcategoryID := uuid.New()

// 	mock := &mockProductQueries{}

// 	handler := &ProductHandler{
// 		Queries: mock,
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		"/subcategories/"+subcategoryID.String()+"/products?offset=abc",
// 		nil,
// 	)

// 	req.SetPathValue("subcategory_id", subcategoryID.String())

// 	rec := httptest.NewRecorder()

// 	handler.HandleListProductsBySubcategory(rec, req)

// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			rec.Code,
// 		)
// 	}
// }
