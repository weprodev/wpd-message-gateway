# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Final Stage
FROM gcr.io/distroless/static-debian12:nonroot AS runner

WORKDIR /app

COPY --from=builder /app/server .
COPY configs/ ./configs/

# Default to example config (memory providers for all message types)
ENV CONFIG_PATH=/app/configs/local.example.yml

EXPOSE 10101

ENTRYPOINT ["/app/server"]
