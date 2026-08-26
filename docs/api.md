# AfriMart API Documentation

## Overview

The AfriMart backend exposes HTTP APIs through the API Gateway.

Current development server:

    http://localhost:8080

The API uses JSON for request and response bodies.

## Common Error Response

Errors use the following format:

{
  "error": "error message"
}

Content-Type:

    application/json; charset=utf-8

## Health

### GET /api/health

Checks whether the API Gateway is running.

### Response

Status: `200 OK`

{
  "status": "ok",
  "service": "api-gateway"
}

---

# Authentication

## POST /api/auth/register

Creates a new user.

### Request

{
  "email": "user@example.com",
  "password": "password"
}

### Success

Status: `201 Created`

{
  "id": "user-id",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "user@example.com"
}

### Errors

`400 Bad Request`

Returned when the request body cannot be decoded.

`500 Internal Server Error`

Returned when password hashing or user creation fails.

---

## POST /api/auth/login

Authenticates a user and creates an access token and refresh token.

### Request

{
  "email": "user@example.com",
  "password": "password"
}

### Success

Status: `200 OK`

{
  "id": "user-id",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "user@example.com",
  "token": "jwt-access-token",
  "refresh_token": "refresh-token"
}

### Errors

`400 Bad Request`

Returned when the request body cannot be decoded.

`401 Unauthorized`

Returned when the email or password is incorrect.

`500 Internal Server Error`

Returned when token or refresh-token creation fails.

---

## POST /api/auth/logout

Revokes the supplied refresh token.

### Authentication

Requires:

    Authorization: Bearer <refresh-token>

### Success

Status: `204 No Content`

### Errors

`401 Unauthorized`

Returned when the Authorization header is missing or invalid, or when the refresh token is invalid.

`500 Internal Server Error`

Returned when the refresh token cannot be revoked.

---

## PUT /api/auth/user

Updates the authenticated user's email or password.

### Authentication

Requires:

    Authorization: Bearer <access-token>

### Request

{
  "email": "new@example.com",
  "password": "new-password"
}

### Success

Returns the updated user.

### Errors

`400 Bad Request`

Returned when the request is invalid.

`401 Unauthorized`

Returned when the access token is missing or invalid.

`500 Internal Server Error`

Returned when the user cannot be updated.

---

## POST /api/refresh

Creates a new access token using a valid refresh token.

### Authentication

Requires:

    Authorization: Bearer <refresh-token>

### Success

Status: `200 OK`

{
  "token": "new-jwt-access-token"
}

### Errors

`401 Unauthorized`

Returned when the refresh token is missing, invalid, expired, or revoked.

`500 Internal Server Error`

Returned when a new access token cannot be created.

---

# Phase 1 API Scaffolding

The following routes are registered in the API Gateway but are currently scaffolding for Phase 1. They are not yet implemented as marketplace functionality.

## Shops

- `POST /api/shops`
- `GET /api/shops`
- `GET /api/shops/{id}`
- `PUT /api/shops/{id}`
- `DELETE /api/shops/{id}`

## Products

- `GET /api/products`
- `GET /api/products/{id}`
- `POST /api/products`
- `PUT /api/products/{id}`
- `DELETE /api/products/{id}`
- `POST /api/products/{id}/images`

## Cart

- `GET /api/cart`
- `POST /api/cart/items`
- `PUT /api/cart/items/{id}`
- `DELETE /api/cart/items/{id}`

## Checkout

- `POST /api/checkout`

## Payments

- `POST /api/payments`
- `GET /api/payments/{id}`
- `POST /api/payments/{id}/verify`

## Orders

- `GET /api/orders`
- `GET /api/orders/{id}`
- `GET /api/seller/orders`
- `PATCH /api/orders/{id}/status`
- `POST /api/orders/{id}/cancel`

These routes currently use a placeholder handler and will be implemented during Phase 1.

---

# Authentication Strategy

The backend currently uses:

- Password hashing for stored passwords.
- JWT access tokens.
- Refresh tokens stored as hashed values.
- Bearer-token authentication through the `Authorization` header.

Access tokens currently have a one-hour expiration.

Refresh tokens are checked for expiration and revocation before a new access token is issued.

# API Gateway

The API Gateway currently listens on:

    :8080

It is responsible for exposing the backend HTTP API and routing requests to the appropriate handlers.

# Development Status

Implemented:

- Health endpoint
- User registration
- User login
- User logout/token revocation
- User update
- Access-token refresh
- Common error responses

Scaffolded for Phase 1:

- Shop management
- Product management
- Cart
- Checkout
- Payments
- Orders