# Build the Go application
FROM golang:1.22 AS builder
WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . ./

# Build the Go application (use ./main.go instead of ./cmd/main.go)
RUN CGO_ENABLED=0 go build -o admin-routes-service ./main.go

# Minimal image for execution
FROM debian:bullseye-slim

# Install necessary dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/admin-routes-service /app/admin-routes-service
COPY --from=builder /app/.env /app/.env

# Expose gRPC port 50053
EXPOSE 50053

# Run the service
CMD ["/app/admin-routes-service"]
