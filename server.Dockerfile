# Build stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# Final stage
FROM alpine:latest

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /app/server .
COPY --chown=appuser:appuser web ./web
COPY --chown=appuser:appuser configs ./configs

USER appuser

EXPOSE 8080

CMD ["./server"]