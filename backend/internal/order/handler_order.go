package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

	cart, err := h.Config.Queries.GetCartByUserID(r.Context(), userID)
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

	order, err := txQueries.CreateOrder(
		r.Context(),
		database.CreateOrderParams{
			UserID:   userID,
			Subtotal: subtotalString,
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
