package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"

	database "github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type updateParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type profile struct{
	Email      string    `json:"email"`
	Username string `json:"username"`
}

func DecodeUpdateParams(r *http.Request) (updateParams, error) {
	var params updateParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return updateParams{}, errors.New("Something went wrong parsing the request body")
	}

	return params, nil
}

func ValidateUpdateParams(params updateParams) error {
	if params.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(params.Email); err != nil {
		return errors.New("invalid email")
	}

	if params.Password == "" {
		return errors.New("password is required")
	}

	if len(params.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

func RespondWithUpdatedUser(w http.ResponseWriter, usr database.User) {
	resp := user{
		Userid:     usr.ID,
		Created_at: usr.CreatedAt,
		Updated_at: usr.UpdatedAt,
		Email:      usr.Email,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}

func RespondWithUserProfile(w http.ResponseWriter, usr database.GetUserByIDRow) {
	resp := profile{
		Email:      usr.Email,
		Username: usr.Username.String,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}
