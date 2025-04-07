
FROM golang:1.16-alpine AS builder

# Set the working directory
WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache gcc musl-dev git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the applications
RUN go build -o /app/bin/user_web web_api/user_web/main.go
RUN go build -o /app/bin/goods_web web_api/goods_web/main.go
RUN go build -o /app/bin/order_web web_api/order_web/main.go
RUN go build -o /app/bin/userop_web web_api/userop_web/main.go
RUN go build -o /app/bin/oss_web web_api/oss_web/main.go

RUN go build -o /app/bin/user_srv nd/user_srv/main.go
RUN go build -o /app/bin/goods_srv nd/goods_srv/main.go
RUN go build -o /app/bin/inventory_srv nd/inventory_srv/main.go
RUN go build -o /app/bin/order_srv nd/order_srv/main.go
RUN go build -o /app/bin/userop_srv nd/userop_srv/main.go

# Final stage
FROM alpine:3.14

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Shanghai

# Create app directories
WORKDIR /app
RUN mkdir -p /app/config /app/logs /app/tmp/nacos/log /app/tmp/nacos/cache

# Copy the built binaries from the builder stage
COPY --from=builder /app/bin /app/bin

# Copy configuration files
COPY --from=builder /app/*/config*.yaml /app/config/
COPY --from=builder /app/scripts/start.sh /app/

# Make the start script executable
RUN chmod +x /app/start.sh

# Expose necessary ports
EXPOSE 8080 8021 8022 8023 8024 8025 9091 9092 9093 9094 9095

# Set environment variable
ENV DEV_CONFIG=1

# Start the application
CMD ["/app/start.sh"]