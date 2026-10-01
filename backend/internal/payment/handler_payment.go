package payment

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/comm"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

const maxChapaReferenceLen = 20

func newTransactionRef() string {
	const prefix = "AFR-"
	hexID := strings.ReplaceAll(uuid.New().String(), "-", "")
	return prefix + hexID[:maxChapaReferenceLen-len(prefix)]
}

type PaymentHandler struct {
	Queries     PaymentQuerier
	Logger      *slog.Logger
	ChapaClient *ChapaClient
}

func NewPaymentHandler(queries PaymentQuerier, logger *slog.Logger) *PaymentHandler {
	secretKey := os.Getenv("CHAPA_SECRET_KEY")
	callbackURL := os.Getenv("CHAPA_CALLBACK_URL")

	if secretKey == "" {
		logger.Warn("CHAPA_SECRET_KEY is not set in environment")
	}

	return &PaymentHandler{
		Queries:     queries,
		Logger:      logger,
		ChapaClient: NewClient(secretKey, callbackURL),
	}
}

type CheckoutRequest struct {
	PaymentMethod   string `json:"payment_method"`
	RecipientName   string `json:"recipient_name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	DeliveryAddress string `json:"delivery_address"`
	DeliveryCity    string `json:"delivery_city"`
	DeliveryNotes   string `json:"delivery_notes"`
}

type CheckoutResponse struct {
	PaymentMethod string    `json:"payment_method"`
	OrderID       uuid.UUID `json:"order_id"`
	TransactionID string    `json:"transaction_id"`
	CheckoutURL   string    `json:"checkout_url,omitempty"`
	Message       string    `json:"message"`
}

func (h *PaymentHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "handling checkout request", "method", r.Method)

	// 1. Authenticate
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "checkout failed: user not authenticated")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "User not authenticated", errors.New("missing user ID"))
		return
	}

	// 2. Decode and Validate Request
	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.WarnContext(ctx, "checkout failed: invalid request payload", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.PaymentMethod != "online" && req.PaymentMethod != "cod" {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid payment method. Use 'online' or 'cod'", nil)
		return
	}

	// 3. Get User to retrieve registered phone number and email
	usr, err := h.Queries.GetUserByIDFull(ctx, userID) 
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get user", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error retrieving user", err)
		return
	}

	// 4. Get Cart and Calculate Authoritative Total
	cart, err := h.Queries.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Cart is empty", nil)
			return
		}
		h.Logger.ErrorContext(ctx, "failed to get cart", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error retrieving cart", err)
		return
	}

	items, err := h.Queries.GetCartItems(ctx, cart.ID)
	if err != nil || len(items) == 0 {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Cart is empty", nil)
		return
	}

	var totalAmount float64
	for _, item := range items {
		// 1. Validate Stock
		if item.Quantity > item.Stock {
			h.Logger.WarnContext(ctx, "insufficient stock", "product_id", item.ProductID, "requested", item.Quantity, "available", item.Stock)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, fmt.Sprintf("Insufficient stock for %s", item.ProductName), nil)
			return
		}

		// 2. Check if product is active
		if item.Status != "active" {
			h.Logger.WarnContext(ctx, "inactive product in cart", "product_id", item.ProductID, "status", item.Status)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, fmt.Sprintf("Product %s is no longer available", item.ProductName), nil)
			return
		}

		// 3. Parse Price (NUMERIC maps to string in Go)
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			h.Logger.ErrorContext(ctx, "failed to parse product price", "product_id", item.ProductID, "price", item.Price, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Invalid product price in database", err)
			return
		}

		// 4. Calculate Total
		totalAmount += price * float64(item.Quantity)
	}

	amountStr := fmt.Sprintf("%.2f", totalAmount)
	transactionID := newTransactionRef()

	// 5. Create Order (Delivery info comes from the request)
	order, err := h.Queries.CreateOrder(ctx, database.CreateOrderParams{
		UserID:          userID,
		Subtotal:        amountStr,
		RecipientName:   req.RecipientName,
		Phone:           req.Phone, // Delivery phone number saved to order
		DeliveryAddress: req.DeliveryAddress,
		DeliveryCity:    req.DeliveryCity,
		DeliveryNotes:   sql.NullString{String: req.DeliveryNotes, Valid: req.DeliveryNotes != ""},
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to create order", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error creating order", err)
		return
	}

	// 6. Create Order Items
	for _, item := range items {
		_, err = h.Queries.CreateOrderItem(ctx, database.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
		if err != nil {
			h.Logger.ErrorContext(ctx, "failed to create order item", "order_id", order.ID, "product_id", item.ProductID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error saving order items", err)
			return
		}
	}

	// 7. Clear the user's cart
	err = h.Queries.ClearCart(ctx, cart.ID)
	if err != nil {
		h.Logger.WarnContext(ctx, "failed to clear cart after order creation", "cart_id", cart.ID, "error", err)
	}

	// 8. Create Payment Record
	payment, err := h.Queries.CreatePayment(ctx, database.CreatePaymentParams{
		OrderID:       order.ID,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: "pending",
		Amount:        amountStr,
		Provider:      sql.NullString{String: req.PaymentMethod, Valid: true},
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to create payment record", "order_id", order.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error creating payment record", err)
		return
	}

	// 9. Prepare Chapa Request using the REGISTERED user data
	// Determine the phone number to send to Chapa (prefer registered, fallback to checkout)
	var customerPhone string
	if usr.PhoneNumber.Valid {
		customerPhone = usr.PhoneNumber.String
	}

	// Fallback to checkout phone if profile phone is empty
	if customerPhone == "" {
		if req.Phone == "" {
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Phone number is required. Please update your profile or provide it at checkout.", nil)
			return
		}
		customerPhone = req.Phone
	}

	// Ensure it's in international format for Chapa
	if !strings.HasPrefix(customerPhone, "+") {
		customerPhone = "+251" + strings.TrimLeft(customerPhone, "0")
	}

	// Safely extract first/last name (FirstName and LastName are sql.NullString)
	firstName := usr.FirstName.String
	if !usr.FirstName.Valid {
		firstName = "Customer"
	}
	
	lastName := usr.LastName.String

	chapaReq := InitializePaymentRequest{
		Amount:            totalAmount,
		Currency:          "ETB",
		MerchantReference: transactionID, 
		Customer: Customer{
			FirstName:   firstName,
			LastName:    lastName,
			Email:       usr.Email, // Email is a plain string (NOT NULL in DB)
			PhoneNumber: customerPhone,
		},
		Meta: Meta{
			OrderID: order.ID.String(),
			Notes:   "AfriMart Order",
		},
	}

	// 10. Initialize Chapa
	chapaResp, err := h.ChapaClient.InitializePayment(ctx, chapaReq)
	if err != nil {
		h.Logger.ErrorContext(ctx, "chapa initialization failed", "order_id", order.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadGateway, "Failed to connect to payment gateway", err)
		return
	}

	// 11. Update Payment Record with Chapa's Reference
	_, err = h.Queries.UpdatePaymentStatusConditional(ctx, database.UpdatePaymentStatusConditionalParams{
		ID:                payment.ID,
		NewStatus:         "pending",
		PaidAt:            sql.NullTime{},
		ProviderReference: sql.NullString{String: chapaResp.Data.Reference, Valid: true},
		TransactionID:     sql.NullString{String: transactionID, Valid: true},
		FailureReason:     sql.NullString{},
		CurrentStatus:     "pending",
	})
	if err != nil {
		h.Logger.WarnContext(ctx, "failed to link chapa reference to payment", "payment_id", payment.ID, "error", err)
	}

	h.Logger.InfoContext(ctx, "checkout initialized successfully", "order_id", order.ID, "transaction_id", transactionID)

	resp := CheckoutResponse{
		PaymentMethod: "online",
		OrderID:       order.ID,
		TransactionID: transactionID,
		CheckoutURL:   chapaResp.Data.CheckoutURL,
		Message:       "Redirect to checkout_url to complete payment",
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "handling payment callback", "method", r.Method)

	reference := r.URL.Query().Get("reference")
	if reference == "" {
		h.Logger.WarnContext(ctx, "callback failed: missing reference")
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Missing payment reference", nil)
		return
	}

	// 1. Fetch our internal payment record using Chapa's provider reference
	payment, err := h.Queries.GetPaymentByProviderRef(ctx, sql.NullString{String: reference, Valid: reference != ""})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "callback failed: payment record not found", "reference", reference)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Payment record not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "failed to get payment record", "reference", reference, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Database error", err)
		return
	}

	// 2. Verify with Chapa API
	verifyResp, err := h.ChapaClient.VerifyPayment(ctx, reference)
	if err != nil {
		h.Logger.ErrorContext(ctx, "chapa verification failed", "reference", reference, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadGateway, "Failed to verify payment with gateway", err)
		return
	}

	// 3. Validate payload matches our internal record
	// verifyResp.Data.MerchantReference is what we sent them (our transactionID)
	if verifyResp.Data.MerchantReference != payment.TransactionID.String {
		h.Logger.ErrorContext(ctx, "merchant reference mismatch", "expected", payment.TransactionID, "got", verifyResp.Data.MerchantReference)
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "Payment verification failed: reference mismatch", nil)
		return
	}

	// Compare amounts (Chapa returns float64, our DB is string)
	chapaAmountStr := fmt.Sprintf("%.2f", verifyResp.Data.Amount)
	if chapaAmountStr != payment.Amount {
		h.Logger.ErrorContext(ctx, "amount mismatch", "expected", payment.Amount, "got", chapaAmountStr)
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "Payment verification failed: amount mismatch", nil)
		return
	}

	// 4. Idempotent State Transition
	var newPaymentStatus string
	var newOrderStatus string
	var httpStatus int
	var message string
	var paidAt sql.NullTime

	switch verifyResp.Data.Status {
	case "success":
		newPaymentStatus = "successful" // Matches your DB CHECK constraint
		newOrderStatus = "confirmed"
		httpStatus = http.StatusOK
		message = "Payment successful. Order confirmed."
		paidAt = sql.NullTime{Time: time.Now(), Valid: true}

	case "failed", "cancelled":
		newPaymentStatus = verifyResp.Data.Status
		newOrderStatus = "cancelled"
		httpStatus = http.StatusPaymentRequired
		message = "Payment was not successful."

	case "pending":
		h.Logger.InfoContext(ctx, "payment still pending", "reference", reference)
		resp := map[string]string{"status": "pending", "message": "Payment is still being processed"}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return

	default:
		h.Logger.WarnContext(ctx, "unknown chapa status", "status", verifyResp.Data.Status, "reference", reference)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Unknown payment status", nil)
		return
	}

	// 5. Conditional Update (Idempotency)
	rowsAffected, err := h.Queries.UpdatePaymentStatusConditional(ctx, database.UpdatePaymentStatusConditionalParams{
		ID:                payment.ID,
		NewStatus:         newPaymentStatus,
		PaidAt:            paidAt,
		ProviderReference: sql.NullString{String: reference, Valid: true},
		TransactionID:     payment.TransactionID,
		FailureReason:     sql.NullString{String: verifyResp.Message, Valid: newPaymentStatus == "failed" || newPaymentStatus == "cancelled"},
		CurrentStatus:     "pending",
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to update payment status", "payment_id", payment.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error updating payment status", err)
		return
	}

	if rowsAffected == 0 {
		h.Logger.InfoContext(ctx, "callback ignored: payment already processed", "payment_id", payment.ID)
		resp := map[string]string{"status": payment.PaymentStatus, "message": "Payment already processed"}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// 6. Update Order Status
	_, err = h.Queries.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     payment.OrderID,
		Status: newOrderStatus,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to update order status", "order_id", payment.OrderID, "error", err)
	}

	h.Logger.InfoContext(ctx, "payment processed successfully", "order_id", payment.OrderID, "new_status", newPaymentStatus)

	resp := map[string]string{
		"status":   newPaymentStatus,
		"message":  message,
		"order_id": payment.OrderID.String(),
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(resp)
}