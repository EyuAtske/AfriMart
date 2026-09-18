package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/EyuAtske/AfriMart/backend/internal/database"
)

type OrderQuerier interface {
	CreateOrder(
		ctx context.Context,
		arg database.CreateOrderParams,
	) (database.Order, error)

	CreateOrderItem(
		ctx context.Context,
		arg database.CreateOrderItemParams,
	) (database.OrderItem, error)

	GetOrderByID(
		ctx context.Context,
		arg database.GetOrderByIDParams,
	) (database.Order, error)

	GetOrderItems(
		ctx context.Context,
		orderID uuid.UUID,
	) ([]database.GetOrderItemsRow, error)

	ListOrdersByUser(
		ctx context.Context,
		arg database.ListOrdersByUserParams,
	) ([]database.Order, error)

	ReduceProductStock(
		ctx context.Context,
		arg database.ReduceProductStockParams,
	) (database.ReduceProductStockRow, error)

	UpdateOrderStatus(
		ctx context.Context,
		arg database.UpdateOrderStatusParams,
	) (database.Order, error)

	VerifyOrderSellerOwnership(
		ctx context.Context,
		arg database.VerifyOrderSellerOwnershipParams,
	) (string, error)

	ListOrdersBySeller(
		ctx context.Context,
		arg database.ListOrdersBySellerParams,
	) ([]database.Order, error)
}
