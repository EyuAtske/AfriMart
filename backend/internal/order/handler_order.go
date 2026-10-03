package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sort"
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

	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// EDIT 1: Temporary default so the current frontend keeps working
	if req.PaymentMethod == "" {
		req.PaymentMethod = "cod"
	}

	if req.PaymentMethod != "online" && req.PaymentMethod != "cod" {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid payment method. Use 'online' or 'cod'", nil)
		return
	}

	req.RecipientName = strings.TrimSpace(req.RecipientName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.DeliveryAddress = strings.TrimSpace(req.DeliveryAddress)
	req.DeliveryCity = strings.TrimSpace(req.DeliveryCity)
	req.DeliveryNotes = strings.TrimSpace(req.DeliveryNotes)

	if req.RecipientName == "" || req.Phone == "" || req.DeliveryAddress == "" || req.DeliveryCity == "" {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Missing required delivery information", nil)
		return
	}

	usr, err := h.Queries.GetUserByIDFull(ctx, userID)
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to retrieve user information", err)
		return
	}

	// EDIT 2: Validate online phone BEFORE starting the transaction
	var customerPhone, firstName, lastName string
	if req.PaymentMethod == "online" {
		if usr.PhoneNumber.Valid {
			customerPhone = strings.TrimSpace(usr.PhoneNumber.String)
		}
		if customerPhone == "" {
			customerPhone = req.Phone
		}
		if customerPhone == "" {
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Phone number is required for online payment", nil)
			return
		}
		if !strings.HasPrefix(customerPhone, "+") {
			customerPhone = "+251" + strings.TrimLeft(customerPhone, "0")
		}
		firstName = usr.FirstName.String
		if !usr.FirstName.Valid || strings.TrimSpace(firstName) == "" {
			firstName = "Customer"
		}
		lastName = usr.LastName.String
	}

	tx, err := h.Config.DB.BeginTx(ctx, nil)
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to start checkout", err)
		return
	}
	defer tx.Rollback()

	txQueries := database.New(tx)
	cart, err := txQueries.GetCartByUserIDForUpdate(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Cart not found", err)
			return
		}
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get cart", err)
		return
	}

	// Prevent buying own products
	ownedItems, _ := txQueries.GetCartItemsOwnedByUser(ctx, database.GetCartItemsOwnedByUserParams{CartID: cart.ID, OwnerID: userID})
	if len(ownedItems) > 0 {
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You cannot purchase your own product", nil)
		return
	}

	cartItems, err := txQueries.GetCartItems(ctx, cart.ID)
	if err != nil || len(cartItems) == 0 {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Cart is empty or could not be loaded", err)
		return
	}

	var subtotal float64
	for _, item := range cartItems {
		if item.Status != "active" || item.Quantity > item.Stock {
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Insufficient stock or product unavailable", nil)
			return
		}
		price, _ := strconv.ParseFloat(item.Price, 64)
		subtotal += price * float64(item.Quantity)
	}
	subtotalString := strconv.FormatFloat(subtotal, 'f', 2, 64)

	order, err := txQueries.CreateOrder(ctx, database.CreateOrderParams{
		UserID: userID, Subtotal: subtotalString, RecipientName: req.RecipientName,
		Phone: req.Phone, DeliveryAddress: req.DeliveryAddress, DeliveryCity: req.DeliveryCity,
		DeliveryNotes: sql.NullString{String: req.DeliveryNotes, Valid: req.DeliveryNotes != ""},
	})
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to create order", err)
		return
	}

	orderItems := make([]database.OrderItem, 0, len(cartItems))
	for _, item := range cartItems {
		orderItem, _ := txQueries.CreateOrderItem(ctx, database.CreateOrderItemParams{
			OrderID: order.ID, ProductID: item.ProductID, Quantity: item.Quantity, Price: item.Price,
		})

		_, err = txQueries.ReduceProductStock(ctx, database.ReduceProductStockParams{ID: item.ProductID, Stock: item.Quantity})
		if err != nil {
			// EDIT 4: Stock race should return 409 Conflict, not 500
			if errors.Is(err, sql.ErrNoRows) {
				comm.RespondErrorWithJson(w, r, http.StatusConflict, "Insufficient stock for product: "+item.ProductName, nil)
				return
			}
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to update stock", err)
			return
		}
		orderItems = append(orderItems, orderItem)
	}

	transactionID := newTransactionRef()
	paymentRecord, _ := txQueries.CreatePayment(ctx, database.CreatePaymentParams{
		OrderID: order.ID, PaymentMethod: req.PaymentMethod, PaymentStatus: "pending",
		Amount: subtotalString, Provider: sql.NullString{String: "chapa", Valid: req.PaymentMethod == "online"},
		TransactionID: sql.NullString{String: transactionID, Valid: true},
	})

	_ = txQueries.ClearCart(ctx, cart.ID)
	if err := tx.Commit(); err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to complete checkout", err)
		return
	}

	if req.PaymentMethod == "cod" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(checkoutResponse{
			Order: order, Items: orderItems, PaymentID: paymentRecord.ID,
			PaymentMethod: "cod", TransactionID: transactionID,
		})
		return
	}

	// Online Payment Flow
	chapaReq := payment.InitializePaymentRequest{
		Amount: subtotal, Currency: "ETB", MerchantReference: transactionID,
		Customer: payment.Customer{FirstName: firstName, LastName: lastName, Email: usr.Email, PhoneNumber: customerPhone},
		Meta:     payment.Meta{OrderID: order.ID.String(), Notes: "AfriMart Order"},
	}

	chapaResp, err := h.ChapaClient.InitializePayment(ctx, chapaReq)
	if err != nil {
		// Compensation for Chapa init failure
		payment.CancelUnpaidOrder(ctx, h.Config.DB, order.ID, paymentRecord.ID, "failed", "Failed to initialize Chapa payment")
		comm.RespondErrorWithJson(w, r, http.StatusBadGateway, "Failed to initialize payment", err)
		return
	}

	// EDIT 3: Save Chapa's provider reference so the callback can find it
	providerRef := chapaResp.Data.Reference
	if providerRef == "" {
		providerRef = transactionID
	}
	if _, err := h.PaymentQuery.UpdatePaymentStatusConditional(ctx, database.UpdatePaymentStatusConditionalParams{
		ID: paymentRecord.ID, CurrentStatus: "pending", NewStatus: "pending",
		ProviderReference: sql.NullString{String: providerRef, Valid: true},
	}); err != nil {
		payment.CancelUnpaidOrder(ctx, h.Config.DB, order.ID, paymentRecord.ID, "failed", "Failed to save payment reference")
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to record payment", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(checkoutResponse{
		Order: order, Items: orderItems, PaymentID: paymentRecord.ID,
		PaymentMethod: "online", TransactionID: transactionID, CheckoutURL: chapaResp.Data.CheckoutURL,
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

	orderRecord, err := h.Queries.GetOrderByID(ctx, database.GetOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "get order failed: order not found", "user_id", userID, "order_id", orderID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Order not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "get order failed: database error", "user_id", userID, "order_id", orderID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get order", err)
		return
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
		tx, err := h.Config.DB.BeginTx(ctx, nil)
		if err != nil {
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to start transaction", err)
			return
		}
		defer tx.Rollback()
		txQueries := database.New(tx)

		updatedOrder, err := txQueries.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
			ID: orderID, Status: request.Status,
		})
		if err != nil {
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to update order status", err)
			return
		}

		// Restore stock
		items, err := txQueries.GetOrderItems(ctx, orderID)
		if err != nil {
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get order items", err)
			return
		}

		// Sort to prevent deadlocks
		sort.Slice(items, func(i, j int) bool {
			return items[i].ProductID.String() < items[j].ProductID.String()
		})

		for _, item := range items {
			if _, err := txQueries.RestoreProductStock(ctx, database.RestoreProductStockParams{
				ID: item.ProductID, Stock: item.Quantity,
			}); err != nil {
				h.Logger.ErrorContext(ctx, "failed to restore stock during cancel", "product_id", item.ProductID, "error", err)
			}
		}

		if err := tx.Commit(); err != nil {
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to commit cancellation", err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"order": updatedOrder})
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
		ID: orderID,
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

	if order.UserID != userID {
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not have permission to cancel this order", nil)
		return
	}

	if order.Status != "pending" && order.Status != "confirmed" {
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Order cannot be cancelled in its current status: "+order.Status, nil)
		return
	}
	tx, err := h.Config.DB.BeginTx(ctx, nil)
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to start transaction", err)
		return
	}
	defer tx.Rollback()

	txQueries := database.New(tx)

	updatedOrder, err := txQueries.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     orderID,
		Status: "cancelled",
	})
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to update order status", err)
		return
	}

	items, err := txQueries.GetOrderItems(ctx, orderID)
	if err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get order items", err)
		return
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ProductID.String() < items[j].ProductID.String()
	})

	for _, item := range items {
		_, err := txQueries.RestoreProductStock(ctx, database.RestoreProductStockParams{
			ID:    item.ProductID,
			Stock: item.Quantity,
		})
		if err != nil {
			h.Logger.ErrorContext(ctx, "failed to restore stock during buyer cancel",
				"product_id", item.ProductID, "quantity", item.Quantity, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to restore product stock", err)
			return
		}
	}

	_, _ = txQueries.UpdatePaymentStatusConditional(ctx, database.UpdatePaymentStatusConditionalParams{
		CurrentStatus: "pending", // Only update if it's still pending
		NewStatus:     "cancelled",
		ID:            orderID, // Note: Adjust this if your payment query expects PaymentID instead of OrderID
	})

	// 8. Commit the transaction
	if err := tx.Commit(); err != nil {
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to commit cancellation", err)
		return
	}

	h.Logger.InfoContext(ctx, "buyer cancelled order and stock restored",
		"order_id", orderID, "user_id", userID)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order cancelled successfully",
		"order":   updatedOrder,
	})
}
