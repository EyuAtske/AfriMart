package auth

import (
	"context"

	"github.com/EyuAtske/AfriMart/backend/internal/database"
)

type AuthQuerier interface {
	CreateUser(
		ctx context.Context,
		arg database.CreateUserParams,
	) (database.User, error)

	UpdateUserPassword(
		ctx context.Context,
		arg database.UpdateUserPasswordParams,
	) (database.User, error)

	UpdateUsername(
		ctx context.Context,
		arg database.UpdateUsernameParams,
	) (database.User, error)
}