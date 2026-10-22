# Multi-Stage Production Build for NESCO Go Core API Gateway
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile optimized static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/nesco-core-api ./cmd/server

# Final minimal runtime
FROM alpine:3.19

WORKDIR /app
RUN apk add --no-cache ca-certificates curl tzdata

COPY --from=builder /bin/nesco-core-api /app/nesco-core-api

EXPOSE 8080

ENTRYPOINT ["/app/nesco-core-api"]

