# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /secretize ./cmd/secretize

# Final stage
FROM alpine:3.18

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

COPY --from=builder /secretize /usr/local/bin/secretize

# KRM functions run as nobody user
USER nobody

ENTRYPOINT ["/usr/local/bin/secretize"] 