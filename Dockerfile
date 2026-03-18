# Build stage
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o .build/rest-api ./cmd/rest

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/.build/rest-api ./rest-api

# Expose port 3000
EXPOSE 3000

# Run the application
CMD ["./rest-api"]
