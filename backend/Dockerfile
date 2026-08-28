FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install certificates and git
RUN apk add --no-cache ca-certificates git

# Copy dependency files first for Docker layer caching
COPY go.mod go.sum ./

RUN go mod download

# Copy the backend source
COPY . .

# Build the backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o afrimart-backend ./cmd/gateway

FROM alpine:latest

WORKDIR /app

# CA certificates are useful for HTTPS/API calls
RUN apk add --no-cache ca-certificates

# Copy compiled binary
COPY --from=builder /app/afrimart-backend .

# Backend API port
EXPOSE 8080

# Start backend
CMD ["./afrimart-backend"]