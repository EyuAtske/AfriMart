package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	database "github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

func decodeUpdateParams(r *http.Request) (updateParams, error) {
	var params updateParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return updateParams{}, errors.New("Something went wrong parsing the request body")
	}

	return params, nil
}

func validateUpdateParams(params updateParams) error {
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

func (apicfg *ApiConfig) updateUser(
	r *http.Request,
	userID uuid.UUID,
	params updateParams,
) (database.User, error) {
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		return database.User{}, err
	}

	return apicfg.Queries.UpdateUserPasswordAndEmail(
		r.Context(),
		database.UpdateUserPasswordAndEmailParams{
			PasswordHash: hashedPassword,
			Email:        params.Email,
			ID:           userID,
		},
	)
}

func respondWithUpdatedUser(w http.ResponseWriter, usr database.User) {
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