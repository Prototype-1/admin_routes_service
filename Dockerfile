 # Use official Golang image as base
FROM golang:1.20

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy the Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o admin_routes_service .

# Command to run the executable
CMD ["/app/admin_routes_service"]

