# Multi-stage build for inventory service

# Build stage
FROM golang:1.16-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventory-service cmd/main.go

# Final stage
FROM alpine:3.14

# Set working directory
WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Shanghai

# Copy the binary from builder
COPY --from=builder /app/inventory-service .
COPY --from=builder /app/configs/config.yaml ./configs/

# Expose ports
EXPOSE 8000 50053

# Command to run
CMD ["./inventory-service"]
