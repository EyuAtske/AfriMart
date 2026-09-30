package cart

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type mockCartQuerier struct {
	getCartByUserIDFn         func(context.Context, uuid.UUID) (database.Cart, error)
	createCartFn              func(context.Context, uuid.UUID) (database.Cart, error)
	getCartItemsFn            func(context.Context, uuid.UUID) ([]database.GetCartItemsRow, error)
	addCartItemFn             func(context.Context, database.AddCartItemParams) (database.CartItem, error)
	updateCartItemQuantityFn  func(context.Context, database.UpdateCartItemQuantityParams) (database.CartItem, error)
	deleteCartItemFn          func(context.Context, database.DeleteCartItemParams) error
	clearCartFn               func(context.Context, uuid.UUID) error
	updateCartTimestampFn     func(context.Context, uuid.UUID) error
	getProductFn              func(context.Context, uuid.UUID) (database.GetProductRow, error)
	getCartItemsOwnedByUserFn func(ctx context.Context, arg database.GetCartItemsOwnedByUserParams) ([]database.GetCartItemsOwnedByUserRow, error)
}

type mockShopQueries struct {
	deactivateShopFunc func(
		ctx context.Context,
		arg database.DeactivateShopParams,
	) (database.Shop, error)

	activateShopFunc func(
		ctx context.Context,
		arg database.ActivateShopParams,
	) (database.Shop, error)

	getShopByIDAndOwnerID func(
		ctx context.Context,
		arg database.GetShopByIDAndOwnerIDParams,
	) (database.Shop, error)
}

func (m *mockCartQuerier) GetCartByUserID(ctx context.Context, userID uuid.UUID) (database.Cart, error) {
	if m.getCartByUserIDFn != nil {
		return m.getCartByUserIDFn(ctx, userID)
	}
	return database.Cart{}, nil
}

func (m *mockCartQuerier) CreateCart(ctx context.Context, userID uuid.UUID) (database.Cart, error) {
	if m.createCartFn != nil {
		return m.createCartFn(ctx, userID)
	}
	return database.Cart{}, nil
}

func (m *mockCartQuerier) GetCartItems(ctx context.Context, cartID uuid.UUID) ([]database.GetCartItemsRow, error) {
	if m.getCartItemsFn != nil {
		return m.getCartItemsFn(ctx, cartID)
	}
	return nil, nil
}

func (m *mockCartQuerier) AddCartItem(ctx context.Context, arg database.AddCartItemParams) (database.CartItem, error) {
	if m.addCartItemFn != nil {
		return m.addCartItemFn(ctx, arg)
	}
	return database.CartItem{}, nil
}

func (m *mockCartQuerier) UpdateCartItemQuantity(ctx context.Context, arg database.UpdateCartItemQuantityParams) (database.CartItem, error) {
	if m.updateCartItemQuantityFn != nil {
		return m.updateCartItemQuantityFn(ctx, arg)
	}
	return database.CartItem{}, nil
}

func (m *mockCartQuerier) DeleteCartItem(ctx context.Context, arg database.DeleteCartItemParams) error {
	if m.deleteCartItemFn != nil {
		return m.deleteCartItemFn(ctx, arg)
	}
	return nil
}

func (m *mockCartQuerier) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	if m.clearCartFn != nil {
		return m.clearCartFn(ctx, cartID)
	}
	return nil
}

func (m *mockCartQuerier) UpdateCartTimestamp(ctx context.Context, cartID uuid.UUID) error {
	if m.updateCartTimestampFn != nil {
		return m.updateCartTimestampFn(ctx, cartID)
	}
	return nil
}

func (m *mockCartQuerier) GetProduct(ctx context.Context, id uuid.UUID) (database.GetProductRow, error) {
	if m.getProductFn != nil {
		return m.getProductFn(ctx, id)
	}
	return database.GetProductRow{}, nil
}

func (m *mockCartQuerier) GetCartItemsOwnedByUser(ctx context.Context, arg database.GetCartItemsOwnedByUserParams) ([]database.GetCartItemsOwnedByUserRow, error) {
	if m.getProductFn != nil {
		return m.getCartItemsOwnedByUserFn(ctx, arg)
	}
	return []database.GetCartItemsOwnedByUserRow{}, nil
}
func (m *mockShopQueries) DeactivateShop(
	ctx context.Context,
	arg database.DeactivateShopParams,
) (database.Shop, error) {
	if m.deactivateShopFunc != nil {
		return m.deactivateShopFunc(ctx, arg)
	}
	return database.Shop{}, sql.ErrNoRows // Default: user does not own the shop
}

func (m *mockShopQueries) ActivateShop(
	ctx context.Context,
	arg database.ActivateShopParams,
) (database.Shop, error) {
	if m.activateShopFunc != nil {
		return m.activateShopFunc(ctx, arg)
	}
	return database.Shop{}, sql.ErrNoRows // Default: user does not own the shop
}

func (m *mockShopQueries) GetShopByIDAndOwnerID(
	ctx context.Context,
	arg database.GetShopByIDAndOwnerIDParams,
) (database.Shop, error) {
	if m.getShopByIDAndOwnerID != nil {
		return m.getShopByIDAndOwnerID(ctx, arg)
	}
	return database.Shop{}, sql.ErrNoRows // Default: user does not own the shop
}

func cartRequestWithUser(method, target string, body string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	ctx := auth.ContextWithUserID(req.Context(), userID)
	return req.WithContext(ctx)
}

func TestHandleGetCartSuccess(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()

	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{
				ID:     cartID,
				UserID: userID,
			}, nil
		},
		getCartItemsFn: func(ctx context.Context, id uuid.UUID) ([]database.GetCartItemsRow, error) {
			return []database.GetCartItemsRow{
				{
					ID:          uuid.New(),
					CartID:      cartID,
					ProductID:   uuid.New(),
					Quantity:    2,
					ProductName: "T-Shirt",
					Price:       "100.00",
					Stock:       10,
					Status:      "active",
				},
			}, nil
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := cartRequestWithUser(http.MethodGet, "/api/cart", "", userID)
	rec := httptest.NewRecorder()

	handler.HandleGetCart(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["subtotal"] != 200.0 {
		t.Fatalf("expected subtotal 200, got %v", response["subtotal"])
	}
}

func TestHandleGetCartUnauthorized(t *testing.T) {
	handler := &CartHandler{
		Queries: &mockCartQuerier{},
		Logger:  slog.Default(),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)
	rec := httptest.NewRecorder()

	handler.HandleGetCart(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandleGetCartCreatesCart(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()
	shopID := uuid.New()

	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{}, sql.ErrNoRows
		},
		createCartFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{
				ID:     cartID,
				UserID: userID,
			}, nil
		},
		getCartItemsFn: func(ctx context.Context, id uuid.UUID) ([]database.GetCartItemsRow, error) {
			return []database.GetCartItemsRow{}, nil
		},
	}

	mockShop := &mockShopQueries{
		getShopByIDAndOwnerID: func(ctx context.Context, arg database.GetShopByIDAndOwnerIDParams) (database.Shop, error) {
			return database.Shop{
				ID:      shopID,
				OwnerID: userID, // Make sure the owner matches the user in the test
				Status:  "active",
			}, nil
		},
	}

	handler := &CartHandler{
		Queries:     mock,
		ShopQueries: mockShop,
		Logger:      slog.Default(),
	}

	req := cartRequestWithUser(http.MethodGet, "/api/cart", "", userID)
	rec := httptest.NewRecorder()

	handler.HandleGetCart(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleAddCartItemSuccess(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	cartID := uuid.New()

	// 1. Setup Cart Mock
	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{ID: cartID, UserID: userID}, nil
		},
		getProductFn: func(ctx context.Context, id uuid.UUID) (database.GetProductRow, error) {
			return database.GetProductRow{
				ID:     productID,
				ShopID: shopID,
				Price:  "100.00",
				Stock:  10,
				Status: "active",
			}, nil
		},
		addCartItemFn: func(ctx context.Context, arg database.AddCartItemParams) (database.CartItem, error) {
			return database.CartItem{ID: uuid.New(), CartID: arg.CartID, ProductID: arg.ProductID, Quantity: arg.Quantity}, nil
		},
	}

	// 2. Buyer does not own the shop -> lookup returns no rows
	mockShop := &mockShopQueries{
		getShopByIDAndOwnerID: func(ctx context.Context, arg database.GetShopByIDAndOwnerIDParams) (database.Shop, error) {
			return database.Shop{}, sql.ErrNoRows
		},
	}

	// 3. Initialize Handler with BOTH mocks explicitly assigned
	handler := &CartHandler{
		Queries:     mock,
		ShopQueries: mockShop, // <-- This MUST be here, otherwise it is nil and panics
		Logger:      slog.Default(),
	}

	// 4. Create request with body
	body := `{"product_id": "` + productID.String() + `", "quantity": 2}`
	req := cartRequestWithUser(http.MethodPost, "/api/cart/items", body, userID)
	rec := httptest.NewRecorder()

	// 5. Execute
	handler.HandleAddCartItem(rec, req)

	// 6. Assert
	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated {
		t.Fatalf("expected status 200 or 201, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleAddCartItemInvalidProductID(t *testing.T) {
	userID := uuid.New()

	handler := &CartHandler{
		Queries: &mockCartQuerier{},
		Logger:  slog.Default(),
	}

	req := cartRequestWithUser(
		http.MethodPost,
		"/api/cart/items",
		`{"product_id":"invalid","quantity":1}`,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleAddCartItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleAddCartItemInvalidQuantity(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()

	handler := &CartHandler{
		Queries: &mockCartQuerier{},
		Logger:  slog.Default(),
	}

	req := cartRequestWithUser(
		http.MethodPost,
		"/api/cart/items",
		`{"product_id":"`+productID.String()+`","quantity":0}`,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleAddCartItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleAddCartItemProductNotFound(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()

	mock := &mockCartQuerier{
		getProductFn: func(ctx context.Context, id uuid.UUID) (database.GetProductRow, error) {
			return database.GetProductRow{}, sql.ErrNoRows
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := cartRequestWithUser(
		http.MethodPost,
		"/api/cart/items",
		`{"product_id":"`+productID.String()+`","quantity":1}`,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleAddCartItem(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandleAddCartItemInactiveProduct(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()

	mock := &mockCartQuerier{
		getProductFn: func(ctx context.Context, id uuid.UUID) (database.GetProductRow, error) {
			return database.GetProductRow{
				ID:     productID,
				Status: "inactive",
				Stock:  10,
			}, nil
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := cartRequestWithUser(
		http.MethodPost,
		"/api/cart/items",
		`{"product_id":"`+productID.String()+`","quantity":1}`,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleAddCartItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleAddCartItemInsufficientStock(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()

	mock := &mockCartQuerier{
		getProductFn: func(ctx context.Context, id uuid.UUID) (database.GetProductRow, error) {
			return database.GetProductRow{
				ID:     productID,
				Status: "active",
				Stock:  2,
			}, nil
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := cartRequestWithUser(
		http.MethodPost,
		"/api/cart/items",
		`{"product_id":"`+productID.String()+`","quantity":5}`,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleAddCartItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleUpdateCartItemSuccess(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()
	itemID := uuid.New()
	productID := uuid.New()

	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{
				ID:     cartID,
				UserID: userID,
			}, nil
		},
		getCartItemsFn: func(ctx context.Context, id uuid.UUID) ([]database.GetCartItemsRow, error) {
			return []database.GetCartItemsRow{
				{
					ID:        itemID,
					CartID:    cartID,
					ProductID: productID,
					Quantity:  1,
				},
			}, nil
		},
		getProductFn: func(ctx context.Context, id uuid.UUID) (database.GetProductRow, error) {
			return database.GetProductRow{
				ID:    productID,
				Stock: 10,
			}, nil
		},
		updateCartItemQuantityFn: func(ctx context.Context, arg database.UpdateCartItemQuantityParams) (database.CartItem, error) {
			if arg.ID != itemID {
				t.Fatalf("unexpected item ID")
			}

			if arg.Quantity != 5 {
				t.Fatalf("expected quantity 5, got %d", arg.Quantity)
			}

			return database.CartItem{
				ID:       itemID,
				CartID:   cartID,
				Quantity: 5,
			}, nil
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/cart/items/"+itemID.String(),
		strings.NewReader(`{"quantity":5}`),
	)

	req.SetPathValue("id", itemID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateCartItem(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleUpdateCartItemNotFound(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()
	itemID := uuid.New()

	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{
				ID:     cartID,
				UserID: userID,
			}, nil
		},
		getCartItemsFn: func(ctx context.Context, id uuid.UUID) ([]database.GetCartItemsRow, error) {
			return []database.GetCartItemsRow{}, nil
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/cart/items/"+itemID.String(),
		strings.NewReader(`{"quantity":2}`),
	)

	req.SetPathValue("id", itemID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateCartItem(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandleUpdateCartItemInsufficientStock(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()
	itemID := uuid.New()
	productID := uuid.New()

	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{ID: cartID, UserID: userID}, nil
		},
		getCartItemsFn: func(ctx context.Context, id uuid.UUID) ([]database.GetCartItemsRow, error) {
			return []database.GetCartItemsRow{
				{
					ID:        itemID,
					CartID:    cartID,
					ProductID: productID,
				},
			}, nil
		},
		getProductFn: func(ctx context.Context, id uuid.UUID) (database.GetProductRow, error) {
			return database.GetProductRow{
				ID:    productID,
				Stock: 2,
			}, nil
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/cart/items/"+itemID.String(),
		strings.NewReader(`{"quantity":10}`),
	)

	req.SetPathValue("id", itemID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateCartItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleDeleteCartItemSuccess(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()
	itemID := uuid.New()

	deleted := false

	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{ID: cartID, UserID: userID}, nil
		},
		deleteCartItemFn: func(ctx context.Context, arg database.DeleteCartItemParams) error {
			deleted = true

			if arg.ID != itemID {
				t.Fatalf("unexpected item ID")
			}

			if arg.CartID != cartID {
				t.Fatalf("unexpected cart ID")
			}

			return nil
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/cart/items/"+itemID.String(),
		nil,
	)

	req.SetPathValue("id", itemID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleDeleteCartItem(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if !deleted {
		t.Fatal("expected DeleteCartItem to be called")
	}
}

func TestHandleClearCartSuccess(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()

	cleared := false

	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{ID: cartID, UserID: userID}, nil
		},
		clearCartFn: func(ctx context.Context, id uuid.UUID) error {
			cleared = true

			if id != cartID {
				t.Fatalf("unexpected cart ID")
			}

			return nil
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := cartRequestWithUser(
		http.MethodDelete,
		"/api/cart",
		"",
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleClearCart(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if !cleared {
		t.Fatal("expected ClearCart to be called")
	}
}

func TestHandleClearCartUnauthorized(t *testing.T) {
	handler := &CartHandler{
		Queries: &mockCartQuerier{},
		Logger:  slog.Default(),
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/cart", nil)
	rec := httptest.NewRecorder()

	handler.HandleClearCart(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestGetOrCreateCartDatabaseError(t *testing.T) {
	userID := uuid.New()

	mock := &mockCartQuerier{
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{}, errors.New("database failure")
		},
	}

	handler := &CartHandler{Queries: mock, Logger: slog.Default()}

	req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)

	_, err := handler.getOrCreateCart(req, userID)

	if err == nil {
		t.Fatal("expected database error")
	}
}

func TestHandleAddCartItemOwnProductForbidden(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	mock := &mockCartQuerier{
		getProductFn: func(ctx context.Context, id uuid.UUID) (database.GetProductRow, error) {
			return database.GetProductRow{ID: productID, ShopID: shopID, Price: "100.00", Stock: 10, Status: "active"}, nil
		},
	}
	mockShop := &mockShopQueries{
		getShopByIDAndOwnerID: func(ctx context.Context, arg database.GetShopByIDAndOwnerIDParams) (database.Shop, error) {
			return database.Shop{ID: shopID, OwnerID: userID, Status: "active"}, nil
		},
	}
	handler := &CartHandler{Queries: mock, ShopQueries: mockShop, Logger: slog.Default()}

	body := `{"product_id": "` + productID.String() + `", "quantity": 1}`
	req := cartRequestWithUser(http.MethodPost, "/api/cart/items", body, userID)
	rec := httptest.NewRecorder()
	handler.HandleAddCartItem(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}