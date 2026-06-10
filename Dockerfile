# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git and certificates if needed
RUN apk add --no-cache git ca-certificates

# Copy dependency files and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o goserver ./cmd/main.go

# Production stage
FROM alpine:3.19

WORKDIR /app

# Copy ca-certificates for outgoing HTTPS requests (if any)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary and migrations
COPY --from=builder /app/goserver .
COPY --from=builder /app/migrations ./migrations

# Expose port 8080
EXPOSE 8080

# Run the server
ENTRYPOINT ["./goserver"]