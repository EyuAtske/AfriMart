# sql/Dockerfile
FROM golang:1.22-alpine

# Install goose ONCE during the build process
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.20.0

# Set goose as the default entrypoint
ENTRYPOINT ["goose"]