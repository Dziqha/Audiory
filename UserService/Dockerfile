# Build stage
FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy binary and config files
COPY --from=builder /build/main /app/main
COPY .env.example /app/.env

# Make sure the binary is executable
RUN chmod +x /app/main

# Expose port
EXPOSE 3000

CMD ["/app/main"]