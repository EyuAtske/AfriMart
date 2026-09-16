package auth

import (
	"context"
	"net/http"

	"github.com/EyuAtske/AfriMart/backend/internal/commErr"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := GetBearerToken(r.Header)
			if err != nil {
				commErr.RespondErrorWithJson(
					w,
					r,
					http.StatusUnauthorized,
					"Missing or invalid Authorization header",
					err,
				)
				return
			}

			userID, err := ValidateJWT(token, secret)
			if err != nil {
				commErr.RespondErrorWithJson(
					w,
					r,
					http.StatusUnauthorized,
					"Invalid or expired token",
					err,
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				userIDKey,
				userID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
