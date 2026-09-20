package order

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type mockOrderQuerier struct {
	createOrderFn                func(context.Context, database.CreateOrderParams) (database.Order, error)
	createOrderItemFn            func(context.Context, database.CreateOrderItemParams) (database.OrderItem, error)
	getOrderByIDFn               func(context.Context, database.GetOrderByIDParams) (database.Order, error)
	getOrderItemsFn              func(context.Context, uuid.UUID) ([]database.GetOrderItemsRow, error)
	listOrdersByUserFn           func(context.Context, database.ListOrdersByUserParams) ([]database.Order, error)
	reduceProductStockFn         func(context.Context, database.ReduceProductStockParams) (database.ReduceProductStockRow, error)
	updateOrderStatusFn          func(context.Context, database.UpdateOrderStatusParams) (database.Order, error)
	verifyOrderSellerOwnershipFn func(context.Context, database.VerifyOrderSellerOwnershipParams) (string, error)
	listOrdersBySellerFn         func(context.Context, database.ListOrdersBySellerParams) ([]database.Order, error)
	getCartByUserIDForUpdatefn   func(context.Context, uuid.UUID) (database.Cart, error)
}

func (m *mockOrderQuerier) CreateOrder(ctx context.Context, arg database.CreateOrderParams) (database.Order, error) {
	if m.createOrderFn != nil {
		return m.createOrderFn(ctx, arg)
	}
	return database.Order{}, nil
}

func (m *mockOrderQuerier) CreateOrderItem(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error) {
	if m.createOrderItemFn != nil {
		return m.createOrderItemFn(ctx, arg)
	}
	return database.OrderItem{}, nil
}

func (m *mockOrderQuerier) GetOrderByID(ctx context.Context, arg database.GetOrderByIDParams) (database.Order, error) {
	if m.getOrderByIDFn != nil {
		return m.getOrderByIDFn(ctx, arg)
	}
	return database.Order{}, nil
}

func (m *mockOrderQuerier) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]database.GetOrderItemsRow, error) {
	if m.getOrderItemsFn != nil {
		return m.getOrderItemsFn(ctx, orderID)
	}
	return nil, nil
}

func (m *mockOrderQuerier) ListOrdersByUser(ctx context.Context, arg database.ListOrdersByUserParams) ([]database.Order, error) {
	if m.listOrdersByUserFn != nil {
		return m.listOrdersByUserFn(ctx, arg)
	}
	return nil, nil
}

func (m *mockOrderQuerier) ReduceProductStock(ctx context.Context, arg database.ReduceProductStockParams) (database.ReduceProductStockRow, error) {
	if m.reduceProductStockFn != nil {
		return m.reduceProductStockFn(ctx, arg)
	}
	return database.ReduceProductStockRow{}, nil
}

func (m *mockOrderQuerier) UpdateOrderStatus(ctx context.Context, arg database.UpdateOrderStatusParams) (database.Order, error) {
	if m.updateOrderStatusFn != nil {
		return m.updateOrderStatusFn(ctx, arg)
	}
	return database.Order{}, nil
}

func (m *mockOrderQuerier) VerifyOrderSellerOwnership(ctx context.Context, arg database.VerifyOrderSellerOwnershipParams) (string, error) {
	if m.verifyOrderSellerOwnershipFn != nil {
		return m.verifyOrderSellerOwnershipFn(ctx, arg)
	}
	return "", nil
}

func (m *mockOrderQuerier) ListOrdersBySeller(ctx context.Context, arg database.ListOrdersBySellerParams) ([]database.Order, error) {
	if m.listOrdersBySellerFn != nil {
		return m.listOrdersBySellerFn(ctx, arg)
	}
	return nil, nil
}

func (m *mockOrderQuerier) GetCartByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (database.Cart, error) {
	if m.getCartByUserIDForUpdatefn != nil {
		return m.getCartByUserIDForUpdatefn(ctx, userID)
	}
	return database.Cart{}, nil
}

func orderRequestWithUser(method, target, body string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	ctx := auth.ContextWithUserID(req.Context(), userID)
	return req.WithContext(ctx)
}

func TestIsValidStatusTransition(t *testing.T) {
	tests := []struct {
		name    string
		current string
		next    string
		valid   bool
	}{
		{"pending to confirmed", "pending", "confirmed", true},
		{"pending to cancelled", "pending", "cancelled", true},
		{"pending to processing", "pending", "processing", false},

		{"confirmed to processing", "confirmed", "processing", true},
		{"confirmed to cancelled", "confirmed", "cancelled", true},
		{"confirmed to shipped", "confirmed", "shipped", false},

		{"processing to shipped", "processing", "shipped", true},
		{"processing to cancelled", "processing", "cancelled", true},
		{"processing to delivered", "processing", "delivered", false},

		{"shipped to delivered", "shipped", "delivered", true},
		{"shipped to cancelled", "shipped", "cancelled", false},

		{"delivered cannot change", "delivered", "cancelled", false},
		{"cancelled cannot change", "cancelled", "confirmed", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidStatusTransition(tt.current, tt.next)

			if got != tt.valid {
				t.Fatalf(
					"expected %v, got %v for %s -> %s",
					tt.valid,
					got,
					tt.current,
					tt.next,
				)
			}
		})
	}
}

func TestHandleCheckoutUnauthorized(t *testing.T) {
	handler := &OrderHandler{
		Config:  &config.ApiConfig{},
		Queries: &mockOrderQuerier{},
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/orders/checkout",
		strings.NewReader(`{}`),
	)

	rec := httptest.NewRecorder()

	handler.HandleCheckout(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestHandleCheckoutInvalidBody(t *testing.T) {
	handler := &OrderHandler{
		Config:  &config.ApiConfig{},
		Queries: &mockOrderQuerier{},
	}

	userID := uuid.New()

	req := orderRequestWithUser(
		http.MethodPost,
		"/api/orders/checkout",
		`invalid json`,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleCheckout(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestHandleCheckoutMissingRecipientName(t *testing.T) {
	handler := &OrderHandler{
		Config:  &config.ApiConfig{},
		Queries: &mockOrderQuerier{},
	}

	userID := uuid.New()

	body := `{
		"phone":"0912345678",
		"delivery_address":"Bole",
		"delivery_city":"Addis Ababa"
	}`

	req := orderRequestWithUser(
		http.MethodPost,
		"/api/orders/checkout",
		body,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleCheckout(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleCheckoutMissingPhone(t *testing.T) {
	handler := &OrderHandler{
		Config:  &config.ApiConfig{},
		Queries: &mockOrderQuerier{},
	}

	userID := uuid.New()

	body := `{
		"recipient_name":"Test User",
		"delivery_address":"Bole",
		"delivery_city":"Addis Ababa"
	}`

	req := orderRequestWithUser(
		http.MethodPost,
		"/api/orders/checkout",
		body,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleCheckout(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleCheckoutMissingAddress(t *testing.T) {
	handler := &OrderHandler{
		Config:  &config.ApiConfig{},
		Queries: &mockOrderQuerier{},
	}

	userID := uuid.New()

	body := `{
		"recipient_name":"Test User",
		"phone":"0912345678",
		"delivery_city":"Addis Ababa"
	}`

	req := orderRequestWithUser(
		http.MethodPost,
		"/api/orders/checkout",
		body,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleCheckout(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleCheckoutMissingCity(t *testing.T) {
	handler := &OrderHandler{
		Config:  &config.ApiConfig{},
		Queries: &mockOrderQuerier{},
	}

	userID := uuid.New()

	body := `{
		"recipient_name":"Test User",
		"phone":"0912345678",
		"delivery_address":"Bole"
	}`

	req := orderRequestWithUser(
		http.MethodPost,
		"/api/orders/checkout",
		body,
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleCheckout(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleListOrdersUnauthorized(t *testing.T) {
	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/orders",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.HandleListOrders(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandleListOrdersSuccess(t *testing.T) {
	userID := uuid.New()

	called := false

	mock := &mockOrderQuerier{
		listOrdersByUserFn: func(
			ctx context.Context,
			arg database.ListOrdersByUserParams,
		) ([]database.Order, error) {
			called = true

			if arg.UserID != userID {
				t.Fatalf("unexpected user ID")
			}

			if arg.Limit != 10 {
				t.Fatalf("expected limit 10, got %d", arg.Limit)
			}

			if arg.Offset != 10 {
				t.Fatalf("expected offset 10, got %d", arg.Offset)
			}

			return []database.Order{
				{},
			}, nil
		},
	}

	handler := &OrderHandler{
		Queries: mock,
	}

	req := orderRequestWithUser(
		http.MethodGet,
		"/api/orders?page=2&limit=10",
		"",
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleListOrders(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !called {
		t.Fatal("expected ListOrdersByUser to be called")
	}
}

func TestHandleListOrdersInvalidPage(t *testing.T) {
	userID := uuid.New()

	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := orderRequestWithUser(
		http.MethodGet,
		"/api/orders?page=0",
		"",
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleListOrders(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleListOrdersInvalidLimit(t *testing.T) {
	userID := uuid.New()

	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := orderRequestWithUser(
		http.MethodGet,
		"/api/orders?limit=51",
		"",
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleListOrders(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleGetOrderSuccess(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	called := false

	mock := &mockOrderQuerier{
		getOrderByIDFn: func(
			ctx context.Context,
			arg database.GetOrderByIDParams,
		) (database.Order, error) {
			called = true

			if arg.ID != orderID {
				t.Fatalf("unexpected order ID")
			}

			if arg.UserID != userID {
				t.Fatalf("unexpected user ID")
			}

			return database.Order{}, nil
		},
		getOrderItemsFn: func(
			ctx context.Context,
			id uuid.UUID,
		) ([]database.GetOrderItemsRow, error) {
			return []database.GetOrderItemsRow{}, nil
		},
	}

	handler := &OrderHandler{
		Queries: mock,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/orders/"+orderID.String(),
		nil,
	)

	req.SetPathValue("id", orderID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleGetOrder(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !called {
		t.Fatal("expected GetOrderByID to be called")
	}
}

func TestHandleGetOrderInvalidID(t *testing.T) {
	userID := uuid.New()

	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := orderRequestWithUser(
		http.MethodGet,
		"/api/orders/invalid",
		"",
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleGetOrder(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleGetOrderNotFound(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	mock := &mockOrderQuerier{
		getOrderByIDFn: func(
			ctx context.Context,
			arg database.GetOrderByIDParams,
		) (database.Order, error) {
			return database.Order{}, sql.ErrNoRows
		},
	}

	handler := &OrderHandler{
		Queries: mock,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/orders/"+orderID.String(),
		nil,
	)

	req.SetPathValue("id", orderID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleGetOrder(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandleUpdateOrderStatusUnauthorized(t *testing.T) {
	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/orders/"+uuid.New().String()+"/status",
		strings.NewReader(`{"status":"confirmed"}`),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateOrderStatus(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandleUpdateOrderStatusSuccess(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	updated := false

	mock := &mockOrderQuerier{
		verifyOrderSellerOwnershipFn: func(
			ctx context.Context,
			arg database.VerifyOrderSellerOwnershipParams,
		) (string, error) {
			if arg.ID != orderID {
				t.Fatalf("unexpected order ID")
			}

			if arg.OwnerID != userID {
				t.Fatalf("unexpected seller ID")
			}

			return "pending", nil
		},
		updateOrderStatusFn: func(
			ctx context.Context,
			arg database.UpdateOrderStatusParams,
		) (database.Order, error) {
			updated = true

			if arg.ID != orderID {
				t.Fatalf("unexpected order ID")
			}

			if arg.Status != "confirmed" {
				t.Fatalf("expected confirmed status, got %s", arg.Status)
			}

			return database.Order{}, nil
		},
	}

	handler := &OrderHandler{
		Queries: mock,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/orders/"+orderID.String()+"/status",
		strings.NewReader(`{"status":"confirmed"}`),
	)

	req.SetPathValue("id", orderID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateOrderStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !updated {
		t.Fatal("expected UpdateOrderStatus to be called")
	}
}

func TestHandleUpdateOrderStatusInvalidStatus(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := orderRequestWithUser(
		http.MethodPatch,
		"/api/orders/"+orderID.String()+"/status",
		`{"status":"random"}`,
		userID,
	)

	req.SetPathValue("id", orderID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateOrderStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleUpdateOrderStatusSellerNotOwner(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	mock := &mockOrderQuerier{
		verifyOrderSellerOwnershipFn: func(
			ctx context.Context,
			arg database.VerifyOrderSellerOwnershipParams,
		) (string, error) {
			return "", errors.New("not owner")
		},
	}

	handler := &OrderHandler{
		Queries: mock,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/orders/"+orderID.String()+"/status",
		strings.NewReader(`{"status":"confirmed"}`),
	)

	req.SetPathValue("id", orderID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateOrderStatus(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestHandleUpdateOrderStatusInvalidTransition(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	updated := false

	mock := &mockOrderQuerier{
		verifyOrderSellerOwnershipFn: func(
			ctx context.Context,
			arg database.VerifyOrderSellerOwnershipParams,
		) (string, error) {
			return "pending", nil
		},
		updateOrderStatusFn: func(
			ctx context.Context,
			arg database.UpdateOrderStatusParams,
		) (database.Order, error) {
			updated = true
			return database.Order{}, nil
		},
	}

	handler := &OrderHandler{
		Queries: mock,
	}

	req := orderRequestWithUser(
		http.MethodPatch,
		"/api/orders/"+orderID.String()+"/status",
		`{"status":"shipped"}`,
		userID,
	)

	req.SetPathValue("id", orderID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateOrderStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if updated {
		t.Fatal("UpdateOrderStatus should not be called for invalid transition")
	}
}

func TestHandleUpdateOrderStatusDeliveredCannotChange(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	mock := &mockOrderQuerier{
		verifyOrderSellerOwnershipFn: func(
			ctx context.Context,
			arg database.VerifyOrderSellerOwnershipParams,
		) (string, error) {
			return "delivered", nil
		},
	}

	handler := &OrderHandler{
		Queries: mock,
	}

	req := orderRequestWithUser(
		http.MethodPatch,
		"/api/orders/"+orderID.String()+"/status",
		`{"status":"cancelled"}`,
		userID,
	)

	req.SetPathValue("id", orderID.String())
	req = req.WithContext(
		auth.ContextWithUserID(req.Context(), userID),
	)

	rec := httptest.NewRecorder()

	handler.HandleUpdateOrderStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleListSellerOrdersUnauthorized(t *testing.T) {
	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/orders/seller",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.HandleListSellerOrders(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandleListSellerOrdersSuccess(t *testing.T) {
	userID := uuid.New()

	called := false

	mock := &mockOrderQuerier{
		listOrdersBySellerFn: func(
			ctx context.Context,
			arg database.ListOrdersBySellerParams,
		) ([]database.Order, error) {
			called = true

			if arg.OwnerID != userID {
				t.Fatalf("unexpected seller ID")
			}

			if arg.Limit != 10 {
				t.Fatalf("expected limit 10, got %d", arg.Limit)
			}

			if arg.Offset != 10 {
				t.Fatalf("expected offset 10, got %d", arg.Offset)
			}

			return []database.Order{
				{},
			}, nil
		},
	}

	handler := &OrderHandler{
		Queries: mock,
	}

	req := orderRequestWithUser(
		http.MethodGet,
		"/api/orders/seller?page=2&limit=10",
		"",
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleListSellerOrders(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !called {
		t.Fatal("expected ListOrdersBySeller to be called")
	}
}

func TestHandleListSellerOrdersInvalidPage(t *testing.T) {
	userID := uuid.New()

	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := orderRequestWithUser(
		http.MethodGet,
		"/api/orders/seller?page=0",
		"",
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleListSellerOrders(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleListSellerOrdersInvalidLimit(t *testing.T) {
	userID := uuid.New()

	handler := &OrderHandler{
		Queries: &mockOrderQuerier{},
	}

	req := orderRequestWithUser(
		http.MethodGet,
		"/api/orders/seller?limit=51",
		"",
		userID,
	)

	rec := httptest.NewRecorder()

	handler.HandleListSellerOrders(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
