package cart

import (
	"context"
	"github.com/google/uuid"

	"github.com/EyuAtske/AfriMart/backend/internal/database"
)

type CartQuerier interface {
	GetCartByUserID(
		ctx context.Context, 
		userID uuid.UUID,
	)(database.Cart, error)

	CreateCart(
		ctx context.Context, 
		userID uuid.UUID,
	) (database.Cart, error)

	GetCartItems(
		ctx context.Context, 
		cartID uuid.UUID,
	) ([]database.GetCartItemsRow, error)

	AddCartItem(
		ctx context.Context, 
		arg database.AddCartItemParams,
	) (database.CartItem, error)

	UpdateCartItemQuantity(
		ctx context.Context, 
		arg database.UpdateCartItemQuantityParams,
	) (database.CartItem, error)

	DeleteCartItem(
		ctx context.Context, 
		arg database.DeleteCartItemParams,
	) error

	ClearCart(
		ctx context.Context, 
		cartID uuid.UUID,
	) error

	UpdateCartTimestamp(
		ctx context.Context, 
		id uuid.UUID,
	) error

	GetProduct(
		ctx context.Context, 
		id uuid.UUID,
	) (database.Product, error)
}
