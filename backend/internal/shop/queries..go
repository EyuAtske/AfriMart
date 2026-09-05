package shop

import (
	"context"

	"github.com/EyuAtske/AfriMart/backend/internal/database"
)

type ShopQuerier interface {
	DeactivateShop(
		ctx context.Context,
		arg database.DeactivateShopParams,
	) (database.Shop, error)

	ActivateShop(
		ctx context.Context,
		arg database.ActivateShopParams,
	) (database.Shop, error)

	GetShopByIDAndOwnerID(
    	ctx context.Context,
    	arg database.GetShopByIDAndOwnerIDParams,
	) (database.Shop, error)
}