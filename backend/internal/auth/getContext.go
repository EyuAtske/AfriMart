package auth

import (
	"context"

	"github.com/google/uuid"
)

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}