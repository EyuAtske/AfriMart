package cart

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/comm"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type CartHandler struct {
	Queries CartQuerier
}

type addCartItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type updateCartItemRequest struct {
	Quantity int32 `json:"quantity"`
}

func (h *CartHandler) getOrCreateCart(
	r *http.Request,
	userID uuid.UUID,
) (database.Cart, error) {
	cart, err := h.Queries.GetCartByUserID(r.Context(), userID)
	if err == nil {
		return cart, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return database.Cart{}, err
	}

	return h.Queries.CreateCart(r.Context(), userID)
}

func (h *CartHandler) HandleGetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get cart",
			err,
		)
		return
	}

	items, err := h.Queries.GetCartItems(r.Context(), cart.ID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get cart items",
			err,
		)
		return
	}

	var subtotal float64

	for _, item := range items {
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

	response := map[string]interface{}{
		"id":       cart.ID,
		"user_id":  cart.UserID,
		"items":    items,
		"subtotal": subtotal,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *CartHandler) HandleAddCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	var params addCartItemRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Error decoding params",
			err,
		)
		return
	}

	productID, err := uuid.Parse(strings.TrimSpace(params.ProductID))
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid product ID",
			err,
		)
		return
	}

	if params.Quantity <= 0 {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Quantity must be greater than zero",
			nil,
		)
		return
	}

	product, err := h.Queries.GetProduct(r.Context(), productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comm.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"Product not found",
				err,
			)
			return
		}

		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get product",
			err,
		)
		return
	}

	if product.Status != "active" {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Product is not available",
			nil,
		)
		return
	}

	if params.Quantity > product.Stock {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Requested quantity exceeds available stock",
			nil,
		)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get cart",
			err,
		)
		return
	}

	item, err := h.Queries.AddCartItem(
		r.Context(),
		database.AddCartItemParams{
			CartID:    cart.ID,
			ProductID: productID,
			Quantity:  params.Quantity,
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not add item to cart",
			err,
		)
		return
	}

	_ = h.Queries.UpdateCartTimestamp(r.Context(), cart.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		return
	}
}

func (h *CartHandler) HandleUpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	itemID, err := uuid.Parse(strings.TrimSpace(r.PathValue("id")))
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid cart item ID",
			err,
		)
		return
	}

	var params updateCartItemRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Error decoding params",
			err,
		)
		return
	}

	if params.Quantity <= 0 {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Quantity must be greater than zero",
			nil,
		)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get cart",
			err,
		)
		return
	}

	items, err := h.Queries.GetCartItems(r.Context(), cart.ID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get cart items",
			err,
		)
		return
	}

	var productID uuid.UUID

	for _, item := range items {
		if item.ID == itemID {
			productID = item.ProductID
			break
		}
	}

	if productID == uuid.Nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusNotFound,
			"Cart item not found",
			nil,
		)
		return
	}

	product, err := h.Queries.GetProduct(r.Context(), productID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusNotFound,
			"Product not found",
			err,
		)
		return
	}

	if params.Quantity > product.Stock {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Requested quantity exceeds available stock",
			nil,
		)
		return
	}

	item, err := h.Queries.UpdateCartItemQuantity(
		r.Context(),
		database.UpdateCartItemQuantityParams{
			ID:       itemID,
			CartID:   cart.ID,
			Quantity: params.Quantity,
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not update cart item",
			err,
		)
		return
	}

	_ = h.Queries.UpdateCartTimestamp(r.Context(), cart.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		return
	}
}

func (h *CartHandler) HandleDeleteCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	itemID, err := uuid.Parse(strings.TrimSpace(r.PathValue("id")))
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusBadRequest,
			"Invalid cart item ID",
			err,
		)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get cart",
			err,
		)
		return
	}

	err = h.Queries.DeleteCartItem(
		r.Context(),
		database.DeleteCartItemParams{
			ID:     itemID,
			CartID: cart.ID,
		},
	)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not remove cart item",
			err,
		)
		return
	}

	_ = h.Queries.UpdateCartTimestamp(r.Context(), cart.ID)

	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) HandleClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"Error getting user id",
			nil,
		)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not get cart",
			err,
		)
		return
	}

	if err := h.Queries.ClearCart(r.Context(), cart.ID); err != nil {
		comm.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Could not clear cart",
			err,
		)
		return
	}

	_ = h.Queries.UpdateCartTimestamp(r.Context(), cart.ID)

	w.WriteHeader(http.StatusNoContent)
}
