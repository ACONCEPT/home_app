# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server

# Run stage
FROM alpine:latest

WORKDIR /

# Copy the binary from builder
COPY --from=builder /api-server /api-server

EXPOSE 8080

ENTRYPOINT ["/api-server"]
