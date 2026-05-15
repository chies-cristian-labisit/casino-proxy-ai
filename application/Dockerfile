# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Download dependencies first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a statically-linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o api \
    ./cmd/api

# Stage 2: Runtime (minimal image)
FROM alpine:3.20

RUN apk --no-cache add \
    ca-certificates \
    tzdata

# Run as non-root
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/api .

USER app

EXPOSE 8081

ENTRYPOINT ["/app/api"]
