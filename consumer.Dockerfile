# Build stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o consumer ./cmd/consumer

# Final stage
FROM alpine:latest

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /app/consumer .
COPY --chown=appuser:appuser configs ./configs

USER appuser

CMD ["./consumer"]