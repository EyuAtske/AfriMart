package payment

import (
	"context"
	"database/sql"

	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type PaymentQuerier interface {
	GetUserByIDFull(ctx context.Context, id uuid.UUID) (database.User, error)
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (database.Cart, error)
	GetCartItems(ctx context.Context, cartID uuid.UUID) ([]database.GetCartItemsRow, error)
	CreateOrder(ctx context.Context, arg database.CreateOrderParams) (database.Order, error)
	CreateOrderItem(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error)
	ClearCart(ctx context.Context, cartID uuid.UUID) error
	CreatePayment(ctx context.Context, arg database.CreatePaymentParams) (database.Payment, error)
	GetPaymentByProviderRef(ctx context.Context, providerRef sql.NullString) (database.Payment, error)
	UpdatePaymentStatusConditional(ctx context.Context, arg database.UpdatePaymentStatusConditionalParams) (int64, error)
	UpdateOrderStatus(ctx context.Context, arg database.UpdateOrderStatusParams) (database.Order, error)
}
