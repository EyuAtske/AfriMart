package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"log/slog"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

// 1. Create the Mock Querier
type mockPaymentQuerier struct {
	getUserByIDFullFn         func(ctx context.Context, id uuid.UUID) (database.User, error)
	getCartByUserIDFn         func(ctx context.Context, id uuid.UUID) (database.Cart, error)
	getCartItemsFn            func(ctx context.Context, id uuid.UUID) ([]database.GetCartItemsRow, error)
	createOrderFn             func(ctx context.Context, arg database.CreateOrderParams) (database.Order, error)
	createOrderItemFn         func(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error)
	clearCartFn               func(ctx context.Context, id uuid.UUID) error
	createPaymentFn           func(ctx context.Context, arg database.CreatePaymentParams) (database.Payment, error)
	updatePaymentStatusCondFn func(ctx context.Context, arg database.UpdatePaymentStatusConditionalParams) (int64, error)
	updateOrderStatusFn       func(ctx context.Context, arg database.UpdateOrderStatusParams) (database.Order, error)
	getPaymentByProviderRefFn func(ctx context.Context, ref sql.NullString) (database.Payment, error)
}

// Implement all interface methods (only showing the ones needed for these tests)
func (m *mockPaymentQuerier) GetUserByIDFull(ctx context.Context, id uuid.UUID) (database.User, error) {
	if m.getUserByIDFullFn != nil {
		return m.getUserByIDFullFn(ctx, id)
	}
	return database.User{}, nil
}
func (m *mockPaymentQuerier) GetCartByUserID(ctx context.Context, id uuid.UUID) (database.Cart, error) {
	if m.getCartByUserIDFn != nil {
		return m.getCartByUserIDFn(ctx, id)
	}
	return database.Cart{}, nil
}
func (m *mockPaymentQuerier) GetCartItems(ctx context.Context, id uuid.UUID) ([]database.GetCartItemsRow, error) {
	if m.getCartItemsFn != nil {
		return m.getCartItemsFn(ctx, id)
	}
	return nil, nil
}
func (m *mockPaymentQuerier) CreateOrder(ctx context.Context, arg database.CreateOrderParams) (database.Order, error) {
	if m.createOrderFn != nil {
		return m.createOrderFn(ctx, arg)
	}
	return database.Order{}, nil
}
func (m *mockPaymentQuerier) CreateOrderItem(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error) {
	if m.createOrderItemFn != nil {
		return m.createOrderItemFn(ctx, arg)
	}
	return database.OrderItem{}, nil
}
func (m *mockPaymentQuerier) ClearCart(ctx context.Context, id uuid.UUID) error {
	if m.clearCartFn != nil {
		return m.clearCartFn(ctx, id)
	}
	return nil
}
func (m *mockPaymentQuerier) CreatePayment(ctx context.Context, arg database.CreatePaymentParams) (database.Payment, error) {
	if m.createPaymentFn != nil {
		return m.createPaymentFn(ctx, arg)
	}
	return database.Payment{}, nil
}
func (m *mockPaymentQuerier) UpdatePaymentStatusConditional(ctx context.Context, arg database.UpdatePaymentStatusConditionalParams) (int64, error) {
	if m.updatePaymentStatusCondFn != nil {
		return m.updatePaymentStatusCondFn(ctx, arg)
	}
	return 0, nil
}
func (m *mockPaymentQuerier) UpdateOrderStatus(ctx context.Context, arg database.UpdateOrderStatusParams) (database.Order, error) {
	if m.updateOrderStatusFn != nil {
		return m.updateOrderStatusFn(ctx, arg)
	}
	return database.Order{}, nil
}
func (m *mockPaymentQuerier) GetPaymentByProviderRef(ctx context.Context, ref sql.NullString) (database.Payment, error) {
	if m.getPaymentByProviderRefFn != nil {
		return m.getPaymentByProviderRefFn(ctx, ref)
	}
	return database.Payment{}, nil
}

// Helper to create request with user context
func paymentRequestWithUser(method, target string, body string, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	ctx := auth.ContextWithUserID(req.Context(), userID)
	return req.WithContext(ctx)
}

// 2. Test Checkout
func TestHandleCheckout_Success(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()
	productID := uuid.New()

	// Setup Mocks
	mock := &mockPaymentQuerier{
		getUserByIDFullFn: func(ctx context.Context, id uuid.UUID) (database.User, error) {
			return database.User{
				ID:          userID,
				Email:       "test@test.com",
				FirstName:   sql.NullString{String: "Test", Valid: true},
				LastName:    sql.NullString{String: "User", Valid: true},
				PhoneNumber: sql.NullString{String: "+251911223344", Valid: true},
			}, nil
		},
		getCartByUserIDFn: func(ctx context.Context, id uuid.UUID) (database.Cart, error) {
			return database.Cart{ID: cartID, UserID: userID}, nil
		},
		getCartItemsFn: func(ctx context.Context, id uuid.UUID) ([]database.GetCartItemsRow, error) {
			return []database.GetCartItemsRow{
				{
					ID:          uuid.New(),
					CartID:      cartID,
					ProductID:   productID,
					Quantity:    2,
					ProductName: "Test Product",
					Price:       "500.00",
					Stock:       10,
					Status:      "active",
				},
			}, nil
		},
		createOrderFn: func(ctx context.Context, arg database.CreateOrderParams) (database.Order, error) {
			return database.Order{ID: uuid.New(), UserID: userID, Subtotal: arg.Subtotal, Status: "pending"}, nil
		},
		createOrderItemFn: func(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error) {
			return database.OrderItem{ID: uuid.New(), OrderID: arg.OrderID, ProductID: arg.ProductID}, nil
		},
		createPaymentFn: func(ctx context.Context, arg database.CreatePaymentParams) (database.Payment, error) {
			return database.Payment{ID: uuid.New(), OrderID: arg.OrderID, Amount: arg.Amount}, nil
		},
		updatePaymentStatusCondFn: func(ctx context.Context, arg database.UpdatePaymentStatusConditionalParams) (int64, error) {
			return 1, nil
		},
	}

	// Mock Chapa API using httptest
	chapaMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Payment initialized",
			"data": map[string]string{
				"reference":    "CHAPA-MOCK-123",
				"checkout_url": "https://checkout.chapa.co/mock",
				"created_at":   time.Now().String(),
				"expires_at":   time.Now().Add(time.Hour).String(),
			},
		})
	}))
	defer chapaMockServer.Close()

	// Initialize Handler
	handler := &PaymentHandler{
		Queries:     mock,
		Logger:      slog.Default(),
		ChapaClient: NewClient("TEST_KEY", chapaMockServer.URL), // Point client to mock server
	}

	handler.ChapaClient.BaseURL = chapaMockServer.URL

	// Execute
	body := `{"payment_method": "online", "recipient_name": "Abebe", "phone": "+251911223344", "email": "a@b.com", "delivery_address": "Bole", "delivery_city": "Addis"}`
	req := paymentRequestWithUser(http.MethodPost, "/api/payments/checkout", body, userID)
	rec := httptest.NewRecorder()

	handler.HandleCheckout(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var resp CheckoutResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.CheckoutURL == "" {
		t.Fatal("expected checkout_url in response")
	}
}

// 3. Test Callback
func TestHandleCallback_Success(t *testing.T) {
	orderID := uuid.New()
	paymentID := uuid.New()
	chapaRef := "CHAPA-MOCK-123"
	transactionID := "AFR-123"

	mock := &mockPaymentQuerier{
		getPaymentByProviderRefFn: func(ctx context.Context, ref sql.NullString) (database.Payment, error) {
			return database.Payment{
				ID:            paymentID,
				OrderID:       orderID,
				TransactionID: sql.NullString{String: transactionID, Valid: true},
				Amount:        "1000.00",
				PaymentStatus: "pending",
			}, nil
		},
		updatePaymentStatusCondFn: func(ctx context.Context, arg database.UpdatePaymentStatusConditionalParams) (int64, error) {
			return 1, nil // Simulate successful idempotent update
		},
		updateOrderStatusFn: func(ctx context.Context, arg database.UpdateOrderStatusParams) (database.Order, error) {
			return database.Order{ID: orderID, Status: "confirmed"}, nil
		},
	}

	// 1. Mock Chapa Verify API
	chapaMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Verified",
			"data": map[string]interface{}{
				"reference":          chapaRef,
				"amount":             1000.00,
				"currency":           "ETB",
				"status":             "success",
				"merchant_reference": transactionID,
			},
		})
	}))
	defer chapaMockServer.Close()

	// 2. Initialize Handler
	handler := &PaymentHandler{
		Queries:     mock,
		Logger:      slog.Default(),
		ChapaClient: NewClient("TEST_KEY", "http://localhost:8080/api/payments/callback"), // Realistic callback URL
	}

	// 🚨 CRITICAL FIX: Explicitly override BaseURL to point to the local mock server
	handler.ChapaClient.BaseURL = chapaMockServer.URL

	// 3. Execute Callback
	req := httptest.NewRequest(http.MethodGet, "/api/payments/callback?reference="+chapaRef+"&status=success", nil)
	rec := httptest.NewRecorder()

	handler.HandleCallback(rec, req)

	// 4. Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "successful" {
		t.Fatalf("expected status 'successful', got '%s'", resp["status"])
	}
}

func TestNewTransactionRef_FitsChapaLimitAndIsUnique(t *testing.T) {
	format := regexp.MustCompile(`^AFR-[0-9a-f]{16}$`)
	seen := make(map[string]struct{}, 10000)

	for i := 0; i < 10000; i++ {
		ref := newTransactionRef()

		if len(ref) > maxChapaReferenceLen {
			t.Fatalf("reference %q is %d chars, Chapa allows at most %d", ref, len(ref), maxChapaReferenceLen)
		}
		if !format.MatchString(ref) {
			t.Fatalf("reference %q does not match expected format", ref)
		}
		if _, dup := seen[ref]; dup {
			t.Fatalf("duplicate reference generated: %q", ref)
		}
		seen[ref] = struct{}{}
	}
}
