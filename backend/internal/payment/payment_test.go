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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EyuAtske/AfriMart/backend/config"
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

// 3. Test Callback

func TestHandleCallback_Success(t *testing.T) {
	orderID := uuid.New()
	paymentID := uuid.New()
	chapaRef := "CHAPA-MOCK-123"
	transactionID := "AFR-123"

	// 1. Setup sqlmock
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	// 2. Tell sqlmock what queries to expect from ConfirmPaidOrder's transaction
	mockDB.ExpectBegin()

	// UpdatePaymentStatusConditional is :execrows, so it uses ExecContext
	mockDB.ExpectExec(`UPDATE payments`).
		WithArgs("successful", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), paymentID, "pending").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// UpdateOrderStatus is :one, so it uses QueryRowContext.
	// We must return all columns defined in the Order struct for the RETURNING * clause.
	mockDB.ExpectQuery(`UPDATE orders`).
		WithArgs(orderID, "confirmed").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "subtotal", "Status", "created_at", "updated_at",
			"recipient_name", "phone", "delivery_address", "delivery_city", "delivery_notes",
		}).AddRow(
			orderID, uuid.New(), "1000.00", "confirmed", time.Now(), time.Now(),
			"Test User", "+251911111111", "Test Address", "Test City", sql.NullString{},
		))

	mockDB.ExpectCommit()

	// 3. Setup the mock Queries
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
	}

	// 4. Setup Chapa Mock Server
	chapaMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
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

	// 5. Initialize Handler WITH the mocked Config
	handler := &PaymentHandler{
		Config: &config.ApiConfig{
			DB: db, // 🚨 THIS FIXES THE NIL POINTER PANIC
		},
		Queries:     mock,
		Logger:      slog.Default(),
		ChapaClient: NewClient("TEST_KEY", "http://localhost:8080/api/payments/callback"),
	}
	// Override BaseURL to point to the local mock server
	handler.ChapaClient.BaseURL = chapaMockServer.URL

	// 6. Execute Callback
	req := httptest.NewRequest(http.MethodGet, "/api/payments/callback?reference="+chapaRef+"&status=success", nil)
	rec := httptest.NewRecorder()
	handler.HandleCallback(rec, req)

	// 7. Assert
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

	// Ensure all sqlmock expectations were met
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
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
