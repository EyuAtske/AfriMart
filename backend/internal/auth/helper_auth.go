package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/EyuAtske/AfriMart/backend/internal/commErr"
	database "github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type updateParams struct {
	Password string `json:"password"`
}

type profile struct {
	Email    string `json:"email"`
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
		Email:    usr.Email,
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

func getUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"User not authenticated",
			errors.New("user ID missing from context"),
		)
		return uuid.Nil, false
	}

	return userID, true
}

func decodeAndValidateUsername(r *http.Request) (updateUsername, error) {
	var params updateUsername

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return params, errors.New("invalid request body")
	}

	params.Username = strings.TrimSpace(params.Username)

	if params.Username == "" {
		return params, errors.New("username is required")
	}

	if len(params.Username) < 3 {
		return params, errors.New("username must be at least 3 characters long")
	}

	if len(params.Username) > 50 {
		return params, errors.New("username must not exceed 50 characters")
	}

	return params, nil
}

func (apicfg *AuthHandler) updateUsername(
	r *http.Request,
	userID uuid.UUID,
	username string,
) (database.User, error) {
	return apicfg.Config.Queries.UpdateUsername(
		r.Context(),
		database.UpdateUsernameParams{
			Username: sql.NullString{
				String: username,
				Valid:  true,
			},
			ID: userID,
		},
	)
}

func validateRegistration(req register) error {
	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email address")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if strings.TrimSpace(req.Username) == "" {
		return errors.New("username is required")
	}

	return nil
}