package cart

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
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
	Logger  *slog.Logger // Added logger
}

// NewCartHandler is a helper to initialize the handler with dependencies
func NewCartHandler(queries CartQuerier, logger *slog.Logger) *CartHandler {
	return &CartHandler{
		Queries: queries,
		Logger:  logger,
	}
}

type addCartItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type updateCartItemRequest struct {
	Quantity int32 `json:"quantity"`
}

func (h *CartHandler) getOrCreateCart(r *http.Request, userID uuid.UUID) (database.Cart, error) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "attempting to get or create cart", "user_id", userID)

	cart, err := h.Queries.GetCartByUserID(ctx, userID)
	if err == nil {
		h.Logger.InfoContext(ctx, "retrieved existing cart", "cart_id", cart.ID)
		return cart, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		h.Logger.ErrorContext(ctx, "failed to get cart from database", "user_id", userID, "error", err)
		return database.Cart{}, err
	}

	h.Logger.InfoContext(ctx, "cart not found, creating new cart", "user_id", userID)
	cart, err = h.Queries.CreateCart(ctx, userID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to create new cart", "user_id", userID, "error", err)
		return database.Cart{}, err
	}

	h.Logger.InfoContext(ctx, "successfully created new cart", "cart_id", cart.ID, "user_id", userID)
	return cart, nil
}

func (h *CartHandler) HandleGetCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "handling get cart request")

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "unauthorized: missing user id in context")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get or create cart", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get cart", err)
		return
	}

	items, err := h.Queries.GetCartItems(ctx, cart.ID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get cart items", "cart_id", cart.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get cart items", err)
		return
	}

	var subtotal float64
	for _, item := range items {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			h.Logger.ErrorContext(ctx, "invalid product price in database", "product_id", item.ProductID, "price", item.Price, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Invalid product price", err)
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

	h.Logger.InfoContext(ctx, "successfully retrieved cart", "cart_id", cart.ID, "item_count", len(items), "subtotal", subtotal)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Changed from StatusCreated to StatusOK for a GET request

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.ErrorContext(ctx, "failed to encode cart response", "error", err)
	}
}

func (h *CartHandler) HandleAddCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "handling add cart item request")

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "unauthorized: missing user id in context")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	var params addCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.Logger.WarnContext(ctx, "failed to decode request body", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error decoding params", err)
		return
	}

	productID, err := uuid.Parse(strings.TrimSpace(params.ProductID))
	if err != nil {
		h.Logger.WarnContext(ctx, "invalid product ID format", "product_id", params.ProductID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	if params.Quantity <= 0 {
		h.Logger.WarnContext(ctx, "invalid quantity requested", "product_id", productID, "quantity", params.Quantity)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Quantity must be greater than zero", nil)
		return
	}

	product, err := h.Queries.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "product not found", "product_id", productID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "failed to get product from database", "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get product", err)
		return
	}

	if product.Status != "active" {
		h.Logger.WarnContext(ctx, "attempted to add inactive product to cart", "product_id", productID, "status", product.Status)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Product is not available", nil)
		return
	}

	if params.Quantity > product.Stock {
		h.Logger.WarnContext(ctx, "requested quantity exceeds stock", "product_id", productID, "requested", params.Quantity, "available", product.Stock)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Requested quantity exceeds available stock", nil)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get or create cart", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get cart", err)
		return
	}

	item, err := h.Queries.AddCartItem(ctx, database.AddCartItemParams{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  params.Quantity,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to add item to cart in database", "cart_id", cart.ID, "product_id", productID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not add item to cart", err)
		return
	}

	if err := h.Queries.UpdateCartTimestamp(ctx, cart.ID); err != nil {
		h.Logger.WarnContext(ctx, "failed to update cart timestamp", "cart_id", cart.ID, "error", err)
		// Non-fatal, continue
	}

	h.Logger.InfoContext(ctx, "successfully added item to cart", "cart_id", cart.ID, "item_id", item.ID, "quantity", params.Quantity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		h.Logger.ErrorContext(ctx, "failed to encode add cart item response", "error", err)
	}
}

func (h *CartHandler) HandleUpdateCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "handling update cart item request")

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "unauthorized: missing user id in context")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	itemID, err := uuid.Parse(strings.TrimSpace(r.PathValue("id")))
	if err != nil {
		h.Logger.WarnContext(ctx, "invalid cart item ID format", "item_id", r.PathValue("id"), "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid cart item ID", err)
		return
	}

	var params updateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.Logger.WarnContext(ctx, "failed to decode request body", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error decoding params", err)
		return
	}

	if params.Quantity <= 0 {
		h.Logger.WarnContext(ctx, "invalid quantity requested", "item_id", itemID, "quantity", params.Quantity)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Quantity must be greater than zero", nil)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get or create cart", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get cart", err)
		return
	}

	items, err := h.Queries.GetCartItems(ctx, cart.ID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get cart items", "cart_id", cart.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get cart items", err)
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
		h.Logger.WarnContext(ctx, "cart item not found in user's cart", "item_id", itemID, "cart_id", cart.ID)
		comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Cart item not found", nil)
		return
	}

	product, err := h.Queries.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "product associated with cart item not found", "product_id", productID)
		} else {
			h.Logger.ErrorContext(ctx, "failed to get product from database", "product_id", productID, "error", err)
		}
		comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Product not found", err)
		return
	}

	if params.Quantity > product.Stock {
		h.Logger.WarnContext(ctx, "requested quantity exceeds stock", "product_id", productID, "requested", params.Quantity, "available", product.Stock)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Requested quantity exceeds available stock", nil)
		return
	}

	item, err := h.Queries.UpdateCartItemQuantity(ctx, database.UpdateCartItemQuantityParams{
		ID:       itemID,
		CartID:   cart.ID,
		Quantity: params.Quantity,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to update cart item quantity in database", "item_id", itemID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not update cart item", err)
		return
	}

	if err := h.Queries.UpdateCartTimestamp(ctx, cart.ID); err != nil {
		h.Logger.WarnContext(ctx, "failed to update cart timestamp", "cart_id", cart.ID, "error", err)
	}

	h.Logger.InfoContext(ctx, "successfully updated cart item", "item_id", itemID, "new_quantity", params.Quantity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		h.Logger.ErrorContext(ctx, "failed to encode update cart item response", "error", err)
	}
}

func (h *CartHandler) HandleDeleteCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "handling delete cart item request")

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "unauthorized: missing user id in context")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	itemID, err := uuid.Parse(strings.TrimSpace(r.PathValue("id")))
	if err != nil {
		h.Logger.WarnContext(ctx, "invalid cart item ID format", "item_id", r.PathValue("id"), "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Invalid cart item ID", err)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get or create cart", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get cart", err)
		return
	}

	err = h.Queries.DeleteCartItem(ctx, database.DeleteCartItemParams{
		ID:     itemID,
		CartID: cart.ID,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to delete cart item from database", "item_id", itemID, "cart_id", cart.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not remove cart item", err)
		return
	}

	if err := h.Queries.UpdateCartTimestamp(ctx, cart.ID); err != nil {
		h.Logger.WarnContext(ctx, "failed to update cart timestamp", "cart_id", cart.ID, "error", err)
	}

	h.Logger.InfoContext(ctx, "successfully deleted cart item", "item_id", itemID, "cart_id", cart.ID)
	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) HandleClearCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.Logger.InfoContext(ctx, "handling clear cart request")

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.Logger.WarnContext(ctx, "unauthorized: missing user id in context")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Error getting user id", nil)
		return
	}

	cart, err := h.getOrCreateCart(r, userID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get or create cart", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not get cart", err)
		return
	}

	if err := h.Queries.ClearCart(ctx, cart.ID); err != nil {
		h.Logger.ErrorContext(ctx, "failed to clear cart in database", "cart_id", cart.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Could not clear cart", err)
		return
	}

	if err := h.Queries.UpdateCartTimestamp(ctx, cart.ID); err != nil {
		h.Logger.WarnContext(ctx, "failed to update cart timestamp", "cart_id", cart.ID, "error", err)
	}

	h.Logger.InfoContext(ctx, "successfully cleared cart", "cart_id", cart.ID)
	w.WriteHeader(http.StatusNoContent)
}