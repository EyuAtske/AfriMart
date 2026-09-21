package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/EyuAtske/AfriMart/backend/internal/comm"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Config  *config.ApiConfig
	Queries AuthQuerier
	Logger  *slog.Logger
}

func NewAuthHandler(cfg *config.ApiConfig, queries AuthQuerier, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		Config:  cfg,
		Queries: queries,
		Logger:  logger,
	}
}

type register struct {
	First    string `json:"firstname"`
	Last     string `json:"Lastname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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
	ctx := r.Context()
	apicfg.Logger.InfoContext(ctx, "handling register request", "method", r.Method, "path", r.URL.Path)

	var reg register
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		apicfg.Logger.WarnContext(ctx, "failed to decode registration request", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error while decoding", err)
		return
	}

	if err := validateRegistration(&reg); err != nil {
		apicfg.Logger.WarnContext(ctx, "registration validation failed", "email", reg.Email, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	hashedPassword, err := HashPassword(reg.Password)
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to hash password", "email", reg.Email, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while processing request", err)
		return
	}

	users, err := apicfg.Queries.CreateUser(ctx, database.CreateUserParams{
		FirstName: sql.NullString{
			String: reg.First,
			Valid:  strings.TrimSpace(reg.First) != "",
		},
		LastName: sql.NullString{
			String: reg.Last,
			Valid:  strings.TrimSpace(reg.Last) != "",
		},
		Username: sql.NullString{
			String: reg.Username,
			Valid:  true,
		},
		Email:        reg.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to create user in database", "email", reg.Email, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while creating user", err)
		return
	}

	respUser := user{
		Userid:     users.ID,
		Created_at: users.CreatedAt,
		Updated_at: users.UpdatedAt,
		Email:      users.Email,
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	
	if err := json.NewEncoder(w).Encode(respUser); err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to encode registration response", "user_id", users.ID, "error", err)
		return
	}

	apicfg.Logger.InfoContext(ctx, "user registered successfully", "user_id", users.ID, "email", users.Email)
}

func (apicfg *AuthHandler) HandleLogIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apicfg.Logger.InfoContext(ctx, "handling login request", "method", r.Method, "path", r.URL.Path)

	var reg login
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		apicfg.Logger.WarnContext(ctx, "failed to decode login request", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Error while decoding request", err)
		return
	}

	usr, err := apicfg.Config.Queries.GetUserByEmail(ctx, reg.Email)
	if err != nil {
		// Log as Warn to track failed attempts, but don't leak if the email exists or not in the response
		apicfg.Logger.WarnContext(ctx, "login failed: user not found or db error", "email", reg.Email, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	check, err := CheckPasswordHash(reg.Password, usr.PasswordHash)
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "error checking password hash", "user_id", usr.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error checking password", err)
		return
	}

	if !check {
		apicfg.Logger.WarnContext(ctx, "login failed: incorrect password", "user_id", usr.ID, "email", usr.Email)
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	expiresInSeconds := 3600
	token, err := MakeJWT(usr.ID, apicfg.Config.Secret, time.Duration(expiresInSeconds)*time.Second)
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to create JWT", "user_id", usr.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while creating token", err)
		return
	}

	refToken := MakeRefreshToken()
	refTokenHash := HashRefreshToken(refToken)

	_, err = apicfg.Config.Queries.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		TokenHash: refTokenHash,
		UserID:    usr.ID,
	})
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to save refresh token to database", "user_id", usr.ID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while creating refresh token", err)
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

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to encode login response", "user_id", usr.ID, "error", err)
		return
	}

	apicfg.Logger.InfoContext(ctx, "user logged in successfully", "user_id", usr.ID, "email", usr.Email)
}

func (apicfg *AuthHandler) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apicfg.Logger.InfoContext(ctx, "handling revoke request", "method", r.Method, "path", r.URL.Path)

	bearerToken, err := GetBearerToken(r.Header)
	if err != nil {
		apicfg.Logger.WarnContext(ctx, "revoke failed: missing or invalid authorization header", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Missing or invalid Authorization header", err)
		return
	}

	tokenHash := HashRefreshToken(bearerToken)

	_, err = apicfg.Config.Queries.RevokeRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			apicfg.Logger.WarnContext(ctx, "revoke failed: refresh token not found", "token_hash", tokenHash)
			comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}
		apicfg.Logger.ErrorContext(ctx, "failed to revoke refresh token in database", "token_hash", tokenHash, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while revoking refresh token", err)
		return
	}

	apicfg.Logger.InfoContext(ctx, "refresh token revoked successfully")
	w.WriteHeader(http.StatusNoContent)
}

func (apicfg *AuthHandler) HandleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apicfg.Logger.InfoContext(ctx, "handling update password request", "method", r.Method, "path", r.URL.Path)

	userID, ok := UserIDFromContext(ctx)
	if !ok {
		apicfg.Logger.WarnContext(ctx, "update password failed: user not authenticated")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "User not authenticated", errors.New("user ID missing from context"))
		return
	}

	params, err := DecodeUpdateParams(r)
	if err != nil {
		apicfg.Logger.WarnContext(ctx, "update password failed: invalid request params", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := ValidateUpdateParams(params); err != nil {
		apicfg.Logger.WarnContext(ctx, "update password failed: validation error", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	hashedPassword, err := HashPassword(params.Password)
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to hash new password", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, err.Error(), err)
		return
	}

	usr, err := apicfg.Queries.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		PasswordHash: hashedPassword,
		ID:           userID,
	})
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to update password in database", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while updating user", err)
		return
	}

	apicfg.Logger.InfoContext(ctx, "password updated successfully", "user_id", userID)
	RespondWithUpdatedUser(w, usr)
}

func (apicfg *AuthHandler) HandleUpdateUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apicfg.Logger.InfoContext(ctx, "handling update username request", "method", r.Method, "path", r.URL.Path)

	userID, ok := getUserID(w, r)
	if !ok {
		// getUserID likely already logged and responded, but we return early
		return
	}

	params, err := decodeAndValidateUsername(r)
	if err != nil {
		apicfg.Logger.WarnContext(ctx, "update username failed: invalid request params", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	usr, err := apicfg.updateUsername(r, userID, params.Username)
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to update username in database", "user_id", userID, "new_username", params.Username, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while updating username", err)
		return
	}

	apicfg.Logger.InfoContext(ctx, "username updated successfully", "user_id", userID, "new_username", params.Username)
	RespondWithUpdatedUser(w, usr)
}

func (apicfg *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apicfg.Logger.InfoContext(ctx, "handling refresh token request", "method", r.Method, "path", r.URL.Path)

	bearerToken, err := GetBearerToken(r.Header)
	if err != nil {
		apicfg.Logger.WarnContext(ctx, "refresh failed: missing or invalid authorization header", "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Missing or invalid Authorization header", err)
		return
	}

	tokenHash := HashRefreshToken(bearerToken)

	refToken, err := apicfg.Config.Queries.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		apicfg.Logger.WarnContext(ctx, "refresh failed: invalid refresh token in database", "token_hash", tokenHash, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if refToken.ExpiresAt.Before(time.Now()) || refToken.RevokedAt.Valid {
		apicfg.Logger.WarnContext(ctx, "refresh failed: token expired or revoked", "user_id", refToken.UserID, "token_hash", tokenHash)
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "Refresh token has expired or been revoked", errors.New("token expired or revoked"))
		return
	}

	newToken, err := MakeJWT(refToken.UserID, apicfg.Config.Secret, time.Hour)
	if err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to create new JWT during refresh", "user_id", refToken.UserID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while creating token", err)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: newToken,
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		apicfg.Logger.ErrorContext(ctx, "failed to encode refresh response", "user_id", refToken.UserID, "error", err)
		return
	}

	apicfg.Logger.InfoContext(ctx, "token refreshed successfully", "user_id", refToken.UserID)
}

func (apicfg *AuthHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apicfg.Logger.InfoContext(ctx, "handling get profile request", "method", r.Method, "path", r.URL.Path)

	userID, ok := UserIDFromContext(ctx)
	if !ok {
		apicfg.Logger.WarnContext(ctx, "get profile failed: user not authenticated")
		comm.RespondErrorWithJson(w, r, http.StatusUnauthorized, "User not authenticated", errors.New("user ID missing from context"))
		return
	}

	usr, err := apicfg.Config.Queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			apicfg.Logger.WarnContext(ctx, "get profile failed: user not found in database", "user_id", userID)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "User not found", err)
			return
		}
		apicfg.Logger.ErrorContext(ctx, "failed to get user from database", "user_id", userID, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error while retrieving user", err)
		return
	}

	apicfg.Logger.InfoContext(ctx, "profile retrieved successfully", "user_id", userID, "email", usr.Email)
	RespondWithUserProfile(w, usr)
}