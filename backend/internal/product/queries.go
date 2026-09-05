package product

import (
	"context"

	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type ProductQuerier interface {
	CreateProduct(
		ctx context.Context,
		arg database.CreateProductParams,
	) (database.Product, error)

	GetProduct(
		ctx context.Context,
		id uuid.UUID,
	) (database.Product, error)

	DeleteProduct(
		ctx context.Context,
		id uuid.UUID,
	) error

	UpdateProduct(
		ctx context.Context,
		arg database.UpdateProductParams,
	) (database.Product, error)

	ListProducts(
		ctx context.Context,
		arg database.ListProductsParams,
	) ([]database.Product, error)

	ListProductsByShop(
		ctx context.Context,
		arg database.ListProductsByShopParams,
	) ([]database.Product, error)

	ListProductsByCategory(
		ctx context.Context,
		arg database.ListProductsByCategoryParams,
	) ([]database.Product, error)

	ListProductsBySubcategory(
		ctx context.Context,
		arg database.ListProductsBySubcategoryParams,
	) ([]database.Product, error)
}