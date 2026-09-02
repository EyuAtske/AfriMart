package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/EyuAtske/AfriMart/backend/internal/commErr"
	"github.com/EyuAtske/AfriMart/backend/internal/config"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Config *config.ApiConfig
}

type register struct {
	First    string `json:"firstname"`
	Last     string `json:"Lastname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type user struct {
	Userid     uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

type loginResponse struct {
	Userid       uuid.UUID `json:"id"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type updateUsername struct {
	Username string `json:"username"`
}

func (apicfg *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(
		r.Context(),
		"request received",
		"method", r.Method,
		"path", r.URL.Path,
	)
	var reqEmail register
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqEmail)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 400, "Error while decoding", err)
		return
	}
	hashedPassword, err := HashPassword(reqEmail.Password)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 500, "Error while decoding request", err)
		return
	}
	users, err := apicfg.Config.Queries.CreateUser(r.Context(), database.CreateUserParams{
		FirstName: sql.NullString{
			String: reqEmail.First,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: reqEmail.Last,
			Valid:  true,
		},
		Username: sql.NullString{
			String: reqEmail.Username,
			Valid:  true,
		},
		Email:        reqEmail.Email,
		PasswordHash: hashedPassword,
		Role:         reqEmail.Role,
	})
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 500, "Error while creating user", err)
		return
	}
	respUser := user{
		Userid:     users.ID,
		Created_at: users.CreatedAt,
		Updated_at: users.UpdatedAt,
		Email:      users.Email,
	}
	data, err := json.Marshal(respUser)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 500, "Error while encoding response", err)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (apicfg *AuthHandler) HandleLogIn(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(
		r.Context(),
		"request received",
		"method", r.Method,
		"path", r.URL.Path,
	)
	var reqEmail login
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqEmail)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 400, "Error while decoding request", err)
		return
	}
	expires_in_seconds := 3600
	usr, err := apicfg.Config.Queries.GetUserByEmail(r.Context(), reqEmail.Email)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 401, "Incorrect email or password", err)
		return
	}
	check, err := CheckPasswordHash(reqEmail.Password, usr.PasswordHash)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 500, "Error checking password", err)
		return
	}
	if !check {
		commErr.RespondErrorWithJson(w, r, 401, "Incorrect email or password", nil)
		return
	}
	token, err := MakeJWT(usr.ID, apicfg.Config.Secret, time.Duration(expires_in_seconds)*time.Second)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 500, "Error while creating token", err)
		return
	}

	refToken := MakeRefreshToken()
	refTokenHash := HashRefreshToken(refToken)

	_, err = apicfg.Config.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		TokenHash: refTokenHash,
		UserID:    usr.ID,
	})
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 500, "Error while creating refresh token", err)
		return
	}

	resp := loginResponse{
		Token:        token,
		Userid:       usr.ID,
		Created_at:   usr.CreatedAt,
		Updated_at:   usr.UpdatedAt,
		Email:        usr.Email,
		RefreshToken: refToken,
	}
	respData, _ := json.Marshal(resp)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}

func (apicfg *AuthHandler) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(
		r.Context(),
		"request received",
		"method", r.Method,
		"path", r.URL.Path,
	)
	bearerToken, err := GetBearerToken(r.Header)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Missing or invalid Authorization header", err)
		return
	}

	tokenHash := HashRefreshToken(bearerToken)

	_, err = apicfg.Config.Queries.RevokeRefreshToken(r.Context(), tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			commErr.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}

		commErr.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (apicfg *AuthHandler) HandleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(
		r.Context(),
		"request received",
		"method", r.Method,
		"path", r.URL.Path,
	)
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"User not authenticated",
			errors.New("user ID missing from context"),
		)
		return
	}

	params, err := DecodeUpdateParams(r)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := ValidateUpdateParams(params); err != nil {
		commErr.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	hashedPassword, err := HashPassword(params.Password)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, http.StatusInternalServerError, err.Error(), err)
		return
	}

	usr, err := apicfg.Config.Queries.UpdateUserPassword(
		r.Context(),
		database.UpdateUserPasswordParams{
			PasswordHash: hashedPassword,
			ID:           userID,
		},
	)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while updating user", err)
		return
	}

	RespondWithUpdatedUser(w, usr)
}

func (apicfg *AuthHandler) HandleUpdateUsername(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(
		r.Context(),
		"request received",
		"method", r.Method,
		"path", r.URL.Path,
	)
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"User not authenticated",
			errors.New("user ID missing from context"),
		)
		return
	}

	var params updateUsername

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		commErr.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	usr, err := apicfg.Config.Queries.UpdateUsername(
		r.Context(),
		database.UpdateUsernameParams{
			Username: sql.NullString{
				String: params.Username,
				Valid:  true,
			},
			ID: userID,
		},
	)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while updating user", err)
		return
	}

	RespondWithUpdatedUser(w, usr)
}

func (apicfg *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(
		r.Context(),
		"request received",
		"method", r.Method,
		"path", r.URL.Path,
	)
	bearerToken, err := GetBearerToken(r.Header)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 401, "Missing or invalid Authorization header", err)
		return
	}
	tokenHash := HashRefreshToken(bearerToken)

	refToken, err := apicfg.Config.Queries.GetRefreshToken(r.Context(), tokenHash)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 401, "Invalid refresh token", err)
		return
	}
	if refToken.ExpiresAt.Before(time.Now()) || refToken.RevokedAt.Valid {
		commErr.RespondErrorWithJson(w, r, 401, "Refresh token has expired or been revoked", errors.New("token expired or revoked"))
		return
	}
	newToken, err := MakeJWT(refToken.UserID, apicfg.Config.Secret, time.Hour)
	if err != nil {
		commErr.RespondErrorWithJson(w, r, 500, "Error while creating token", err)
		return
	}
	resp := struct {
		Token string `json:"token"`
	}{
		Token: newToken,
	}
	respData, _ := json.Marshal(resp)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}

func (apicfg *AuthHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusUnauthorized,
			"User not authenticated",
			errors.New("user ID missing from context"),
		)
		return
	}

	usr, err := apicfg.Config.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			commErr.RespondErrorWithJson(
				w,
				r,
				http.StatusNotFound,
				"User not found",
				err,
			)
			return
		}

		commErr.RespondErrorWithJson(
			w,
			r,
			http.StatusInternalServerError,
			"Error while retrieving user",
			err,
		)
		return
	}

	RespondWithUserProfile(w, usr)
}