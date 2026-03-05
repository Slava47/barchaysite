# ── Stage 1: build ──────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Download dependencies first (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Build a statically linked binary (no libc dependency)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o barchay-api .

# ── Stage 2: minimal Alpine image ───────────────────────────────────────────
# We use Alpine (not scratch) so that the healthcheck has access to wget.
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /src/barchay-api .

# Data and uploads directories are mounted as volumes at runtime
EXPOSE 8080

ENV PORT=8080 \
    DATA_DIR=/app/data \
    UPLOADS_DIR=/app/uploads \
    PUBLIC_URL=http://localhost:8080 \
    ALLOWED_ORIGINS=http://localhost

CMD ["./barchay-api"]
