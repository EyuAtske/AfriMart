package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/commErr"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type email struct {
	Password string `json:"password"`
	Email    string `json:"email"`
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

func (apicfg *ApiConfig) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var reqEmail email
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqEmail)
	if err != nil {
		commErr.RespondErrorWithJson(w, 400, "Error while decoding")
		return
	}
	hashedPassword, err := auth.HashPassword(reqEmail.Password)
	if err != nil {
		commErr.RespondErrorWithJson(w, 500, "Error while decoding request")
		return
	}
	users, err := apicfg.Queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:        reqEmail.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		commErr.RespondErrorWithJson(w, 500, "Error while creating user")
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
		commErr.RespondErrorWithJson(w, 500, "Error while encoding response")
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (apicfg *ApiConfig) HandleLogIn(w http.ResponseWriter, r *http.Request) {
	var reqEmail email
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqEmail)
	if err != nil {
		commErr.RespondErrorWithJson(w, 400, "Error while decoding request")
		return
	}
	expires_in_seconds := 3600
	usr, err := apicfg.Queries.GetUserByEmail(r.Context(), reqEmail.Email)
	if err != nil {
		commErr.RespondErrorWithJson(w, 401, "Incorrect email or password")
		return
	}
	check, _ := auth.CheckPasswordHash(reqEmail.Password, usr.PasswordHash)
	if !check {
		commErr.RespondErrorWithJson(w, 401, "Incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(usr.ID, apicfg.Secret, time.Duration(expires_in_seconds)*time.Second)
	if err != nil {
		commErr.RespondErrorWithJson(w, 500, "Error while creating token")
		return
	}

	refToken := auth.MakeRefreshToken()
	refTokenHash := auth.HashRefreshToken(refToken)

	_, err = apicfg.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		TokenHash: refTokenHash,
		UserID:    usr.ID,
	})
	if err != nil {
		commErr.RespondErrorWithJson(w, 500, "Error while creating refresh token")
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

func (apicfg *ApiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		commErr.RespondErrorWithJson(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	tokenHash := auth.HashRefreshToken(bearerToken)

	_, err = apicfg.Queries.RevokeRefreshToken(r.Context(), tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			commErr.RespondErrorWithJson(w, http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		commErr.RespondErrorWithJson(w, http.StatusInternalServerError, "Error while revoking refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (apicfg *ApiConfig) HandleUpdates(w http.ResponseWriter, r *http.Request) {
	type updateParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		commErr.RespondErrorWithJson(w, 401, "Missing or invalid Authorization header")
		return
	}
	id, err := auth.ValidateJWT(bearerToken, apicfg.Secret)
	if err != nil {
		commErr.RespondErrorWithJson(w, 401, fmt.Sprintf("Invalid token %v", err))
		return
	}
	var params updateParams
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		commErr.RespondErrorWithJson(w, 400, "Something went wrong parsing the request body")
		return
	}
	if params.Email == "" {
		commErr.RespondErrorWithJson(w, http.StatusBadRequest, "email is required")
		return
	}

	if _, err := mail.ParseAddress(params.Email); err != nil {
		commErr.RespondErrorWithJson(w, http.StatusBadRequest, "invalid email")
		return
	}

	if params.Password == "" {
		commErr.RespondErrorWithJson(w, http.StatusBadRequest, "password is required")
		return
	}

	if len(params.Password) < 8 {
    	commErr.RespondErrorWithJson(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		commErr.RespondErrorWithJson(w, 500, "Error while hashing password")
		return
	}
	usr, err := apicfg.Queries.UpdateUserPasswordAndEmail(r.Context(), database.UpdateUserPasswordAndEmailParams{
		PasswordHash: hashedPassword,
		Email:        params.Email,
		ID:           id,
	})
	if err != nil {
		commErr.RespondErrorWithJson(w, 500, "Error while updating user")
		return
	}
	resp := user{
		Userid:     usr.ID,
		Created_at: usr.CreatedAt,
		Updated_at: usr.UpdatedAt,
		Email:      usr.Email,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		commErr.RespondErrorWithJson(w, 500, "Error while encoding response")
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (apicfg *ApiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		commErr.RespondErrorWithJson(w, 401, "Missing or invalid Authorization header")
		return
	}
	tokenHash := auth.HashRefreshToken(bearerToken)

	refToken, err := apicfg.Queries.GetRefreshToken(r.Context(), tokenHash)
	if err != nil {
		commErr.RespondErrorWithJson(w, 401, "Invalid refresh token")
		return
	}
	if refToken.ExpiresAt.Before(time.Now()) || refToken.RevokedAt.Valid {
		commErr.RespondErrorWithJson(w, 401, "Refresh token has expired or been revoked")
		return
	}
	newToken, err := auth.MakeJWT(refToken.UserID, apicfg.Secret, time.Hour)
	if err != nil {
		commErr.RespondErrorWithJson(w, 500, "Error while creating token")
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
