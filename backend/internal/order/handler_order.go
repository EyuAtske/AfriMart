package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/comm"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/EyuAtske/AfriMart/backend/internal/payment"
)

type OrderHandler struct {
	Config       *config.ApiConfig
	Queries      OrderQuerier
	PaymentQuery PaymentUpdater
	Logger       *slog.Logger
	ChapaClient  *payment.ChapaClient
}

func NewOrderHandler(cfg *config.ApiConfig, queries OrderQuerier, paymentQuery PaymentUpdater, logger *slog.Logger, chapaClient *payment.ChapaClient) *OrderHandler {
	return &OrderHandler{
		Config:       cfg,
		Queries:      queries,
		PaymentQuery: paymentQuery,
		Logger:       logger,
		ChapaClient:  chapaClient,
	}
}

type checkoutResponse struct {
	Order         database.Order       `json:"order"`
	Items         []database.OrderItem `json:"items"`
	PaymentID     uuid.UUID            `json:"payment_id"`
	PaymentMethod string               `json:"payment_method"`
	TransactionID string               `json:"transaction_id"`
	CheckoutURL   string               `json:"checkout_url,omitempty"`
}

type checkoutRequest struct {
	PaymentMethod   string `json:"payment_method"`
	RecipientName   string `json:"recipient_name"`
	Phone           string `json:"phone"`
	DeliveryAddress string `json:"delivery_address"`
	DeliveryCity    string `json:"delivery_city"`
	DeliveryNotes   string `json:"delivery_notes"`
}

const maxChapaReferenceLen = 20

func newTransactionRef() string {
	const prefix = "AFR-"

	id := strings.ReplaceAll(uuid.New().String(), "-", "")

	return prefix + id[:maxChapaReferenceLen-len(prefix)]
}

func (h *OrderHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "checkout failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	h.Logger.InfoContext(ctx, "Checkout initiated", "user_id", userID)

	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.WarnContext(ctx, "checkout failed: invalid request body", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = "cod"
	}

	if req.PaymentMethod != "online" && req.PaymentMethod != "cod" {
		h.Logger.WarnContext(ctx, "checkout failed: invalid payment method", "method", req.PaymentMethod)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid payment method. Use 'online' or 'cod'", nil)
		return
	}

	req.RecipientName = strings.TrimSpace(req.RecipientName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.DeliveryAddress = strings.TrimSpace(req.DeliveryAddress)
	req.DeliveryCity = strings.TrimSpace(req.DeliveryCity)
	req.DeliveryNotes = strings.TrimSpace(req.DeliveryNotes)

	if req.RecipientName == "" || req.Phone == "" || req.DeliveryAddress == "" || req.DeliveryCity == "" {
		h.Logger.WarnContext(ctx, "checkout failed: missing required delivery information", "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Missing required delivery information", nil)
		return
	}

	usr, err := h.Queries.GetUserByIDFull(ctx, userID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: failed to retrieve user information", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to retrieve user information", err)
		return
	}
	h.Logger.InfoContext(ctx, "User retrieved for checkout", "user_id", userID)

	var customerPhone, firstName, lastName string
	if req.PaymentMethod == "online" {
		if usr.PhoneNumber.Valid {
			customerPhone = strings.TrimSpace(usr.PhoneNumber.String)
		}
		if customerPhone == "" {
			customerPhone = req.Phone
		}
		if customerPhone == "" {
			h.Logger.WarnContext(ctx, "checkout failed: phone number required for online payment", "user_id", userID)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Phone number is required for online payment", nil)
			return
		}
		if !strings.HasPrefix(customerPhone, "+") {
			customerPhone = "251" + strings.TrimLeft(customerPhone, "0")
		}
		firstName = usr.FirstName.String
		if !usr.FirstName.Valid || strings.TrimSpace(firstName) == "" {
			firstName = "Customer"
		}
		lastName = usr.LastName.String
	}

	tx, err := h.Config.DB.BeginTx(ctx, nil)
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: failed to start transaction", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to start checkout", err)
		return
	}
	defer tx.Rollback()

	txQueries := database.New(tx)

	cart, err := txQueries.GetCartByUserIDForUpdate(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "checkout failed: cart not found", "user_id", userID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Cart not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "checkout failed: failed to get cart", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get cart", err)
		return
	}
	h.Logger.InfoContext(ctx, "Cart retrieved for checkout", "user_id", userID, "cart_id", cart.ID)

	ownedItems, err := txQueries.GetCartItemsOwnedByUser(ctx, database.GetCartItemsOwnedByUserParams{CartID: cart.ID, OwnerID: userID})
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout warning: failed to check owned items", "cart_id", cart.ID, "error", err)
	} else if len(ownedItems) > 0 {
		h.Logger.WarnContext(ctx, "checkout failed: user attempting to purchase own product", "user_id", userID, "cart_id", cart.ID)
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You cannot purchase your own product", nil)
		return
	}

	cartItems, err := txQueries.GetCartItems(ctx, cart.ID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: GetCartItems error", "error", err, "cart_id", cart.ID, "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to load cart items", err)
		return
	}
	if len(cartItems) == 0 {
		h.Logger.WarnContext(ctx, "checkout failed: cart is empty", "cart_id", cart.ID, "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Cart is empty", nil)
		return
	}
	h.Logger.InfoContext(ctx, "Validating cart items", "cart_id", cart.ID, "item_count", len(cartItems))

	var subtotal float64
	for _, item := range cartItems {
		if item.Status != "active" || item.Quantity > item.Stock {
			h.Logger.WarnContext(ctx, "checkout failed: insufficient stock or product unavailable",
				"product_id", item.ProductID, "product_name", item.ProductName, "requested", item.Quantity, "available", item.Stock)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Insufficient stock or product unavailable: "+item.ProductName, nil)
			return
		}
		price, _ := strconv.ParseFloat(item.Price, 64)
		subtotal += price * float64(item.Quantity)
	}
	subtotalString := strconv.FormatFloat(subtotal, 'f', 2, 64)
	h.Logger.InfoContext(ctx, "Cart validation successful", "cart_id", cart.ID, "subtotal", subtotalString)

	// UPDATED: Added Method field to match new schema
	order, err := txQueries.CreateOrder(ctx, database.CreateOrderParams{
		UserID:          userID,
		Subtotal:        subtotalString,
		Method:          req.PaymentMethod, // <-- NEW FIELD
		RecipientName:   req.RecipientName,
		Phone:           req.Phone,
		DeliveryAddress: req.DeliveryAddress,
		DeliveryCity:    req.DeliveryCity,
		DeliveryNotes:   sql.NullString{String: req.DeliveryNotes, Valid: req.DeliveryNotes != ""},
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: failed to create order", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to create order", err)
		return
	}
	h.Logger.InfoContext(ctx, "Order created successfully", "order_id", order.ID, "user_id", userID, "method", req.PaymentMethod)

	orderItems := make([]database.OrderItem, 0, len(cartItems))
	for _, item := range cartItems {
		orderItem, err := txQueries.CreateOrderItem(ctx, database.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
		if err != nil {
			h.Logger.ErrorContext(ctx, "checkout failed: failed to create order item", "order_id", order.ID, "product_id", item.ProductID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to create order item", err)
			return
		}

		_, err = txQueries.ReduceProductStock(ctx, database.ReduceProductStockParams{ID: item.ProductID, Stock: item.Quantity})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.Logger.WarnContext(ctx, "checkout failed: concurrent stock depletion", "order_id", order.ID, "product_id", item.ProductID, "product_name", item.ProductName)
				comm.RespondErrorWithJson(w, r, http.StatusConflict, "Insufficient stock for product: "+item.ProductName, nil)
				return
			}
			h.Logger.ErrorContext(ctx, "checkout failed: failed to update stock", "order_id", order.ID, "product_id", item.ProductID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to update stock", err)
			return
		}
		h.Logger.DebugContext(ctx, "Product stock reduced", "order_id", order.ID, "product_id", item.ProductID, "quantity", item.Quantity)
		orderItems = append(orderItems, orderItem)
	}

	transactionID := newTransactionRef()
	paymentRecord, err := txQueries.CreatePayment(ctx, database.CreatePaymentParams{
		OrderID:       order.ID,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: "pending",
		Amount:        subtotalString,
		Provider:      sql.NullString{String: "chapa", Valid: req.PaymentMethod == "online"},
		TransactionID: sql.NullString{String: transactionID, Valid: true},
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: failed to create payment record", "order_id", order.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to create payment record", err)
		return
	}
	h.Logger.InfoContext(ctx, "Payment record created", "payment_id", paymentRecord.ID, "transaction_id", transactionID, "order_id", order.ID)

	if err := txQueries.ClearCart(ctx, cart.ID); err != nil {
		h.Logger.ErrorContext(ctx, "checkout warning: failed to clear cart", "cart_id", cart.ID, "order_id", order.ID, "error", err)
		// Note: We don't return here because the order is already created, but we log it for investigation
	}

	if err := tx.Commit(); err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: failed to commit transaction", "order_id", order.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to complete checkout", err)
		return
	}
	h.Logger.InfoContext(ctx, "Checkout transaction committed successfully", "order_id", order.ID)

	if req.PaymentMethod == "cod" {
		h.Logger.InfoContext(ctx, "COD checkout completed", "order_id", order.ID, "transaction_id", transactionID)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(checkoutResponse{
			Order:         order,
			Items:         orderItems,
			PaymentID:     paymentRecord.ID,
			PaymentMethod: "cod",
			TransactionID: transactionID,
		})
		return
	}

	h.Logger.InfoContext(ctx, "Initializing Chapa payment", "order_id", order.ID, "amount", subtotal, "customer_phone", customerPhone)
	chapaReq := payment.InitializePaymentRequest{
		Amount:            subtotal,
		Currency:          "ETB",
		MerchantReference: transactionID,
		Customer: payment.Customer{
			FirstName:   firstName,
			LastName:    lastName,
			Email:       usr.Email,
			PhoneNumber: customerPhone,
		},
		Meta: payment.Meta{
			OrderID: order.ID.String(),
			Notes:   "AfriMart Order",
		},
		ReturnURL: h.Config.ReturnURL,
	}

	chapaResp, err := h.ChapaClient.InitializePayment(ctx, chapaReq)
	h.Logger.InfoContext(ctx, "Chapa Raw Response",
		"order_id", order.ID,
		"status", chapaResp != nil && chapaResp.Status != "", // safe check
		"message", chapaResp != nil && chapaResp.Message != "",
		"error", err,
	)

	if err != nil {
		h.Logger.ErrorContext(ctx, "Chapa initialization failed", "order_id", order.ID, "transaction_id", transactionID, "error", err)
		applied, cancelErr := payment.CancelUnpaidOrder(ctx, h.Config.DB, order.ID, paymentRecord.ID, "failed", "Failed to initialize Chapa payment")
		if cancelErr != nil || !applied {
			h.Logger.ErrorContext(ctx, "CRITICAL: Compensation failed! Order needs manual cleanup",
				"order_id", order.ID,
				"payment_id", paymentRecord.ID,
				"error", cancelErr)
		} else {
			h.Logger.InfoContext(ctx, "Compensation successful: stock restored and order cancelled",
				"order_id", order.ID, "payment_id", paymentRecord.ID)
		}
		comm.RespondErrorWithJson(w, r, http.StatusBadGateway, "Failed to initialize payment", err)
		return
	}

	if chapaResp == nil || chapaResp.Data.CheckoutURL == "" {
		h.Logger.ErrorContext(ctx, "Chapa returned empty checkout URL", "order_id", order.ID, "response", chapaResp)
		applied, cancelErr := payment.CancelUnpaidOrder(ctx, h.Config.DB, order.ID, paymentRecord.ID, "failed", "Chapa returned empty checkout URL")
		if cancelErr != nil || !applied {
			h.Logger.ErrorContext(ctx, "CRITICAL: Compensation failed! Order needs manual cleanup",
				"order_id", order.ID,
				"payment_id", paymentRecord.ID,
				"error", cancelErr)
		} else {
			h.Logger.InfoContext(ctx, "Compensation successful: stock restored and order cancelled",
				"order_id", order.ID, "payment_id", paymentRecord.ID)
		}
		comm.RespondErrorWithJson(w, r, http.StatusBadGateway, "Payment gateway did not return a checkout URL. Please check your Chapa configuration.", nil)
		return
	}

	providerRef := chapaResp.Data.Reference
	if providerRef == "" {
		providerRef = transactionID
	}

	h.Logger.InfoContext(ctx, "Saving Chapa provider reference", "order_id", order.ID, "payment_id", paymentRecord.ID, "provider_ref", providerRef)
	if _, err := h.PaymentQuery.UpdatePaymentStatusConditional(ctx, database.UpdatePaymentStatusConditionalParams{
		ID:                paymentRecord.ID,
		CurrentStatus:     "pending",
		NewStatus:         "pending",
		ProviderReference: sql.NullString{String: providerRef, Valid: true},
	}); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to save payment reference", "order_id", order.ID, "payment_id", paymentRecord.ID, "error", err)
		// Attempt compensation if we can't even save the reference
		payment.CancelUnpaidOrder(ctx, h.Config.DB, order.ID, paymentRecord.ID, "failed", "Failed to save payment reference")
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to record payment", err)
		return
	}

	h.Logger.InfoContext(ctx, "Online checkout completed successfully", "order_id", order.ID, "checkout_url", chapaResp.Data.CheckoutURL)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(checkoutResponse{
		Order:         order,
		Items:         orderItems,
		PaymentID:     paymentRecord.ID,
		PaymentMethod: "online",
		TransactionID: transactionID,
		CheckoutURL:   chapaResp.Data.CheckoutURL,
	})
}

func (h *OrderHandler) HandleListOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "list orders failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	h.Logger.InfoContext(ctx, "handling list orders request", "user_id", userID)

	page := 1
	limit := 20
	query := r.URL.Query()

	if value := query.Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			h.Logger.WarnContext(ctx, "list orders failed: invalid page param", "user_id", userID, "page", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid page", err)
			return
		}
		page = parsed
	}

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 50 {
			h.Logger.WarnContext(ctx, "list orders failed: invalid limit param", "user_id", userID, "limit", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid limit", err)
			return
		}
		limit = parsed
	}

	offset := (page - 1) * limit

	orders, err := h.Queries.ListOrdersByUser(ctx, database.ListOrdersByUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "list orders failed: database error", "user_id", userID, "page", page, "limit", limit, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get orders", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := map[string]interface{}{
		"page":   page,
		"limit":  limit,
		"orders": orders,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "list orders succeeded but failed to encode response", "user_id", userID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "orders listed successfully", "user_id", userID, "page", page, "limit", limit, "count", len(orders))
}

func (h *OrderHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "get order failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	orderIDStr := r.PathValue("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "get order failed: invalid order ID format", "user_id", userID, "order_id", orderIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	h.Logger.InfoContext(ctx, "handling get order request", "user_id", userID, "order_id", orderID)

	// 1. Try fetching as the buyer first
	orderRecord, err := h.Queries.GetOrderByID(ctx, database.GetOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 2. If not found as buyer, check if the user is the seller
			sellerOrder, sellerErr := h.Queries.GetOrderByIDForSeller(ctx, database.GetOrderByIDForSellerParams{
				ID:      orderID,
				OwnerID: userID,
			})
			if sellerErr != nil {
				if errors.Is(sellerErr, sql.ErrNoRows) {
					h.Logger.WarnContext(ctx, "get order failed: order not found for buyer or seller", "user_id", userID, "order_id", orderID)
					comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Order not found", sellerErr)
					return
				}
				h.Logger.ErrorContext(ctx, "get order failed: database error checking seller ownership", "user_id", userID, "order_id", orderID, "error", sellerErr)
				comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get order", sellerErr)
				return
			}
			// Successfully fetched as seller
			orderRecord = sellerOrder
		} else {
			h.Logger.ErrorContext(ctx, "get order failed: database error", "user_id", userID, "order_id", orderID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get order", err)
			return
		}
	}

	items, err := h.Queries.GetOrderItems(ctx, orderID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "get order failed: could not get order items", "user_id", userID, "order_id", orderID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get order items", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := map[string]interface{}{
		"order": orderRecord,
		"items": items,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "get order succeeded but failed to encode response", "user_id", userID, "order_id", orderID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "order retrieved successfully", "user_id", userID, "order_id", orderID, "item_count", len(items))
}

func (h *OrderHandler) HandleUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	orderIDStr := r.PathValue("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var request struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate status transition
	currentStatus, err := h.Queries.VerifyOrderSellerOwnership(ctx, database.VerifyOrderSellerOwnershipParams{
		ID: orderID, OwnerID: userID,
	})
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not own this order", err)
		return
	}

	if !isValidStatusTransition(currentStatus, request.Status) {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid order status transition", nil)
		return
	}

	// If cancelling, we MUST restore stock in a transaction
	if request.Status == "cancelled" && currentStatus != "cancelled" {
		// Fetch payment record
		paymentRecord, err := h.Queries.GetPaymentByOrderID(ctx, orderID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			h.Logger.ErrorContext(ctx, "failed to get payment record for seller cancellation", "order_id", orderID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to retrieve payment record", err)
			return
		}

		var paymentID uuid.UUID
		if err == nil {
			paymentID = paymentRecord.ID
		}

		// Use the atomic helper
		applied, err := payment.CancelUnpaidOrder(ctx, h.Config.DB, orderID, paymentID, "cancelled", "Cancelled by seller")
		if err != nil {
			h.Logger.ErrorContext(ctx, "failed to cancel unpaid order (seller)", "order_id", orderID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to cancel order", err)
			return
		}

		if !applied {
			comm.RespondErrorWithJson(w, r, http.StatusConflict, "Order is no longer in a cancellable state", nil)
			return
		}

		h.Logger.InfoContext(ctx, "seller cancelled order and stock restored", "order_id", orderID, "seller_id", userID)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Order cancelled successfully",
		})
		return
	}

	// Normal status update (no stock restoration needed)
	updatedOrder, err := h.Queries.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID: orderID, Status: request.Status,
	})
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to update order status", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"order": updatedOrder})
}

func isValidStatusTransition(current, next string) bool {
	switch current {
	case "pending":
		return next == "confirmed" || next == "cancelled"
	case "confirmed":
		return next == "processing" || next == "cancelled"
	case "processing":
		return next == "shipped" || next == "cancelled"
	case "shipped":
		return next == "delivered"
	case "delivered", "cancelled":
		return false
	default:
		return false
	}
}

func (h *OrderHandler) HandleListSellerOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "list seller orders failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Authentication required", nil)
		return
	}
	h.Logger.InfoContext(ctx, "handling list seller orders request", "user_id", userID)

	page := 1
	limit := 20

	if value := r.URL.Query().Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			h.Logger.WarnContext(ctx, "list seller orders failed: invalid page param", "user_id", userID, "page", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid page", err)
			return
		}
		page = parsed
	}

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 50 {
			h.Logger.WarnContext(ctx, "list seller orders failed: invalid limit param", "user_id", userID, "limit", value, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid limit", err)
			return
		}
		limit = parsed
	}

	offset := (page - 1) * limit

	orders, err := h.Queries.ListOrdersBySeller(ctx, database.ListOrdersBySellerParams{
		OwnerID: userID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "list seller orders failed: database error", "user_id", userID, "page", page, "limit", limit, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get seller orders", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := map[string]interface{}{
		"page":   page,
		"limit":  limit,
		"orders": orders,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "list seller orders succeeded but failed to encode response", "user_id", userID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "seller orders listed successfully", "user_id", userID, "page", page, "limit", limit, "count", len(orders))
}

func (h *OrderHandler) HandleCancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	orderIDStr := r.PathValue("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	order, err := h.Queries.GetOrderByID(ctx, database.GetOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Order not found", err)
			return
		}
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to retrieve order", err)
		return
	}

	if order.Status != "pending" && order.Status != "confirmed" {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Order cannot be cancelled in its current status: "+order.Status, nil)
		return
	}

	// 1. Fetch the associated payment record to get the paymentID
	paymentRecord, err := h.Queries.GetPaymentByOrderID(ctx, orderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.Logger.ErrorContext(ctx, "failed to get payment record for cancellation", "order_id", orderID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to retrieve payment record", err)
		return
	}

	var paymentID uuid.UUID
	if err == nil {
		paymentID = paymentRecord.ID
	}

	// 2. Use the existing, battle-tested helper for atomic cancellation
	applied, err := payment.CancelUnpaidOrder(ctx, h.Config.DB, orderID, paymentID, "cancelled", "Cancelled by buyer")
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to cancel unpaid order", "order_id", orderID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to cancel order", err)
		return
	}

	if !applied {
		// This safely prevents cancelling an order that has already been successfully paid for
		comm.RespondErrorWithJson(w, r, http.StatusConflict, "Order is no longer in a cancellable state (may be already paid)", nil)
		return
	}

	// 3. Fetch the updated order to return in the response
	updatedOrder, _ := h.Queries.GetOrderByID(ctx, database.GetOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	})

	h.Logger.InfoContext(ctx, "buyer cancelled order and stock restored", "order_id", orderID, "user_id", userID)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order cancelled successfully",
		"order":   updatedOrder,
	})
}
