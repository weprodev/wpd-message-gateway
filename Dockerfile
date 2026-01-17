# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies (if needed)
# RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# CGO_ENABLED=0: Disable CGO for static binary
# -ldflags="-w -s": Strip debug information for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Final Stage
FROM gcr.io/distroless/static-debian12:nonroot AS runner

WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Copy configuration files
COPY configs/ ./configs/

# Expose port
EXPOSE 10101

# Command to run the executable
ENTRYPOINT ["/app/server"]
