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
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Unauthorized",
			nil,
		)
		return
	}

	var req checkoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid request body",
			err,
		)
		return
	}

	req.RecipientName = strings.TrimSpace(req.RecipientName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.DeliveryAddress = strings.TrimSpace(req.DeliveryAddress)
	req.DeliveryCity = strings.TrimSpace(req.DeliveryCity)
	req.DeliveryNotes = strings.TrimSpace(req.DeliveryNotes)

	if req.RecipientName == "" {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Recipient name is required",
			nil,
		)
		return
	}

	if req.Phone == "" {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Phone is required",
			nil,
		)
		return
	}

	if req.DeliveryAddress == "" {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Delivery address is required",
			nil,
		)
		return
	}

	if req.DeliveryCity == "" {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Delivery city is required",
			nil,
		)
		return
	}

	tx, err := h.Config.DB.BeginTx(r.Context(), nil)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to start checkout",
			err,
		)
		return
	}

	defer tx.Rollback()

	txQueries := database.New(tx)

	cart, err := txQueries.GetCartByUserIDForUpdate(
		r.Context(),
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Cart not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to get cart",
			err,
		)
		return
	}

	cartItems, err := h.Config.Queries.GetCartItems(r.Context(), cart.ID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to get cart items",
			err,
		)
		return
	}

	if len(cartItems) == 0 {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Cart is empty",
			nil,
		)
		return
	}

	// Calculate the subtotal using the current product prices.
	var subtotal float64

	for _, item := range cartItems {
		if item.Status != "active" {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"One or more products are no longer available",
				nil,
			)
			return
		}

		if item.Quantity > item.Stock {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Insufficient stock for product: "+item.ProductName,
				nil,
			)
			return
		}

		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusInternalServerError,
				"Invalid product price",
				err,
			)
			return
		}

		subtotal += price * float64(item.Quantity)
	}

	subtotalString := strconv.FormatFloat(subtotal, 'f', 2, 64)

	order, err := txQueries.CreateOrder(
		r.Context(),
		database.CreateOrderParams{
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
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to create order",
			err,
		)
		return
	}

	orderItems := make([]database.OrderItem, 0, len(cartItems))

	for _, item := range cartItems {
		orderItem, err := txQueries.CreateOrderItem(
			r.Context(),
			database.CreateOrderItemParams{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			},
		)
		if err != nil {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusInternalServerError,
				"Failed to create order item",
				err,
			)
			return
		}

		_, err = txQueries.ReduceProductStock(
			r.Context(),
			database.ReduceProductStockParams{
				ID:    item.ProductID,
				Stock: item.Quantity,
			},
		)
		if err != nil {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Insufficient stock for product: "+item.ProductName,
				nil,
			)
			return
		}

		orderItems = append(orderItems, orderItem)
	}

	if err := txQueries.ClearCart(r.Context(), cart.ID); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to clear cart",
			err,
		)
		return
	}

	if err := tx.Commit(); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to complete checkout",
			err,
		)
		return
	}

	response := checkoutResponse{
		Order: order,
		Items: orderItems,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) HandleListOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Unauthorized",
			nil,
		)
		return
	}

	page := 1
	limit := 20

	query := r.URL.Query()

	if value := query.Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid page",
				err,
			)
			return
		}
		page = parsed
	}

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 50 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid limit",
				err,
			)
			return
		}
		limit = parsed
	}

	offset := (page - 1) * limit

	orders, err := h.Queries.ListOrdersByUser(
		r.Context(),
		database.ListOrdersByUserParams{
			UserID: userID,
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to get orders",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"page":   page,
		"limit":  limit,
		"orders": orders,
	})
}

func (h *OrderHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Unauthorized",
			nil,
		)
		return
	}

	orderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid order ID",
			err,
		)
		return
	}

	orderRecord, err := h.Queries.GetOrderByID(
		r.Context(),
		database.GetOrderByIDParams{
			ID:     orderID,
			UserID: userID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Order not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to get order",
			err,
		)
		return
	}

	items, err := h.Queries.GetOrderItems(
		r.Context(),
		orderID,
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to get order items",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"order": orderRecord,
		"items": items,
	})
}

func (h *OrderHandler) HandleUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Unauthorized",
			nil,
		)
		return
	}

	orderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid order ID",
			err,
		)
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid request body",
			err,
		)
		return
	}

	switch request.Status {
	case "confirmed", "processing", "shipped", "delivered", "cancelled":
		// Valid status.
	default:
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid order status",
			nil,
		)
		return
	}

	currentStatus, err := h.Queries.VerifyOrderSellerOwnership(
		r.Context(),
		database.VerifyOrderSellerOwnershipParams{
			ID:      orderID,
			OwnerID: userID,
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusForbidden,
			"You do not own this order",
			err,
		)
		return
	}

	if !isValidStatusTransition(currentStatus, request.Status) {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid order status transition",
			nil,
		)
		return
	}

	updatedOrder, err := h.Queries.UpdateOrderStatus(
		r.Context(),
		database.UpdateOrderStatusParams{
			ID:     orderID,
			Status: request.Status,
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to update order status",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"order": updatedOrder,
	})
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
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Authentication required",
			nil,
		)
		return
	}

	page := 1
	limit := 20

	if value := r.URL.Query().Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid page",
				err,
			)
			return
		}
		page = parsed
	}

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 50 {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusBadRequest,
				"Invalid limit",
				err,
			)
			return
		}
		limit = parsed
	}

	offset := (page - 1) * limit

	orders, err := h.Queries.ListOrdersBySeller(
		r.Context(),
		database.ListOrdersBySellerParams{
			OwnerID: userID,
			Limit:   int32(limit),
			Offset:  int32(offset),
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Failed to get seller orders",
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := map[string]interface{}{
		"page":   page,
		"limit":  limit,
		"orders": orders,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode seller orders", "error", err)
	}
}
