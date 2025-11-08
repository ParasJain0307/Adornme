# ================================
# Stage 1 — Build the Go binary
# ================================
FROM golang:1.24.3 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project (your Go source code + config files)
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o adornme ./cmd/adronme-code-server/main.go

# ================================
# Stage 2 — Create a minimal image
# ================================
FROM alpine:3.20

# Set working directory inside the container
WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/adornme .

# Copy config files
COPY --from=builder /app/config ./config

# Expose your application port
EXPOSE 9090

# Run the binary
ENTRYPOINT ["./adornme"]

