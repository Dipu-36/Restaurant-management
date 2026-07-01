# -----------------------------
# Build Stage
# -----------------------------
FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# -----------------------------
# Runtime Stage
# -----------------------------
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 4000

CMD ["./api"]
