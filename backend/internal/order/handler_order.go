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
)

type OrderHandler struct {
	Config  *config.ApiConfig
	Queries OrderQuerier
	Logger  *slog.Logger
}

func NewOrderHandler(cfg *config.ApiConfig, queries OrderQuerier, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{
		Config:  cfg,
		Queries: queries,
		Logger:  logger,
	}
}

type checkoutResponse struct {
	Order database.Order       `json:"order"`
	Items []database.OrderItem `json:"items"`
}

type checkoutRequest struct {
	RecipientName   string `json:"recipient_name"`
	Phone           string `json:"phone"`
	DeliveryAddress string `json:"delivery_address"`
	DeliveryCity    string `json:"delivery_city"`
	DeliveryNotes   string `json:"delivery_notes"`
}

func (h *OrderHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "checkout failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	h.Logger.InfoContext(ctx, "handling checkout request", "user_id", userID)

	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.WarnContext(ctx, "checkout failed: invalid request body", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	req.RecipientName = strings.TrimSpace(req.RecipientName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.DeliveryAddress = strings.TrimSpace(req.DeliveryAddress)
	req.DeliveryCity = strings.TrimSpace(req.DeliveryCity)
	req.DeliveryNotes = strings.TrimSpace(req.DeliveryNotes)

	if req.RecipientName == "" {
		h.Logger.WarnContext(ctx, "checkout failed: recipient name required", "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Recipient name is required", nil)
		return
	}
	if req.Phone == "" {
		h.Logger.WarnContext(ctx, "checkout failed: phone required", "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Phone is required", nil)
		return
	}
	if req.DeliveryAddress == "" {
		h.Logger.WarnContext(ctx, "checkout failed: delivery address required", "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Delivery address is required", nil)
		return
	}
	if req.DeliveryCity == "" {
		h.Logger.WarnContext(ctx, "checkout failed: delivery city required", "user_id", userID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Delivery city is required", nil)
		return
	}

	tx, err := h.Config.DB.BeginTx(ctx, nil)
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: could not start transaction", "user_id", userID, "error", err)
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
		h.Logger.ErrorContext(ctx, "checkout failed: could not get cart for update", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get cart", err)
		return
	}

	// Note: Using h.Config.Queries here instead of txQueries as per original code
	cartItems, err := h.Config.Queries.GetCartItems(ctx, cart.ID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: could not get cart items", "user_id", userID, "cart_id", cart.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to get cart items", err)
		return
	}

	if len(cartItems) == 0 {
		h.Logger.WarnContext(ctx, "checkout failed: cart is empty", "user_id", userID, "cart_id", cart.ID)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Cart is empty", nil)
		return
	}

	var subtotal float64
	for _, item := range cartItems {
		if item.Status != "active" {
			h.Logger.WarnContext(ctx, "checkout failed: product no longer available", "user_id", userID, "product_id", item.ProductID, "product_name", item.ProductName)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "One or more products are no longer available", nil)
			return
		}

		if item.Quantity > item.Stock {
			h.Logger.WarnContext(ctx, "checkout failed: insufficient stock", "user_id", userID, "product_id", item.ProductID, "product_name", item.ProductName, "requested", item.Quantity, "available", item.Stock)
			comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Insufficient stock for product: "+item.ProductName, nil)
			return
		}

		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			h.Logger.ErrorContext(ctx, "checkout failed: invalid product price in db", "user_id", userID, "product_id", item.ProductID, "price", item.Price, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Invalid product price", err)
			return
		}

		subtotal += price * float64(item.Quantity)
	}

	subtotalString := strconv.FormatFloat(subtotal, 'f', 2, 64)

	order, err := txQueries.CreateOrder(ctx, database.CreateOrderParams{
		UserID:          userID,
		Subtotal:        subtotalString,
		RecipientName:   req.RecipientName,
		Phone:           req.Phone,
		DeliveryAddress: req.DeliveryAddress,
		DeliveryCity:    req.DeliveryCity,
		DeliveryNotes: sql.NullString{
			String: req.DeliveryNotes,
			Valid:  true,
		},
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: could not create order", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to create order", err)
		return
	}

	orderItems := make([]database.OrderItem, 0, len(cartItems))

	for _, item := range cartItems {
		orderItem, err := txQueries.CreateOrderItem(ctx, database.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
		if err != nil {
			h.Logger.ErrorContext(ctx, "checkout failed: could not create order item", "user_id", userID, "order_id", order.ID, "product_id", item.ProductID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to create order item", err)
			return
		}

		_, err = txQueries.ReduceProductStock(ctx, database.ReduceProductStockParams{
			ID:    item.ProductID,
			Stock: item.Quantity,
		})
		if err != nil {
			h.Logger.ErrorContext(ctx, "checkout failed: could not reduce product stock", "user_id", userID, "order_id", order.ID, "product_id", item.ProductID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to update stock for product: "+item.ProductName, err)
			return
		}

		orderItems = append(orderItems, orderItem)
	}

	if err := txQueries.ClearCart(ctx, cart.ID); err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: could not clear cart", "user_id", userID, "cart_id", cart.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to clear cart", err)
		return
	}

	if err := tx.Commit(); err != nil {
		h.Logger.ErrorContext(ctx, "checkout failed: could not commit transaction", "user_id", userID, "order_id", order.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to complete checkout", err)
		return
	}

	response := checkoutResponse{
		Order: order,
		Items: orderItems,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "checkout succeeded but failed to encode response", "user_id", userID, "order_id", order.ID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "checkout completed successfully", "user_id", userID, "order_id", order.ID, "subtotal", subtotal)
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
		h.Logger.WarnContext(ctx, "update order status failed: unauthorized")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	orderIDStr := r.PathValue("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.Logger.WarnContext(ctx, "update order status failed: invalid order ID format", "user_id", userID, "order_id", orderIDStr, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var request struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.WarnContext(ctx, "update order status failed: invalid request body", "user_id", userID, "order_id", orderID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	switch request.Status {
	case "confirmed", "processing", "shipped", "delivered", "cancelled":
		// Valid status.
	default:
		h.Logger.WarnContext(ctx, "update order status failed: invalid status value", "user_id", userID, "order_id", orderID, "status", request.Status)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid order status", nil)
		return
	}

	h.Logger.InfoContext(ctx, "handling update order status request", "user_id", userID, "order_id", orderID, "requested_status", request.Status)

	currentStatus, err := h.Queries.VerifyOrderSellerOwnership(ctx, database.VerifyOrderSellerOwnershipParams{
		ID:      orderID,
		OwnerID: userID,
	})
	if err != nil {
		h.Logger.WarnContext(ctx, "update order status failed: user does not own this order", "user_id", userID, "order_id", orderID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "You do not own this order", err)
		return
	}

	if !isValidStatusTransition(currentStatus, request.Status) {
		h.Logger.WarnContext(ctx, "update order status failed: invalid status transition", "user_id", userID, "order_id", orderID, "current_status", currentStatus, "requested_status", request.Status)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid order status transition", nil)
		return
	}

	updatedOrder, err := h.Queries.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     orderID,
		Status: request.Status,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "update order status failed: database error", "user_id", userID, "order_id", orderID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Failed to update order status", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"order": updatedOrder,
	}); err != nil {
		h.Logger.ErrorContext(ctx, "update order status succeeded but failed to encode response", "user_id", userID, "order_id", orderID, "error", err)
		return
	}

	h.Logger.InfoContext(ctx, "order status updated successfully", "user_id", userID, "order_id", orderID, "new_status", request.Status)
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