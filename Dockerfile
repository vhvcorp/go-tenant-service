# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.work and modules
COPY go.work go.work
COPY pkg/go.mod pkg/go.sum pkg/
COPY services/tenant-service/go.mod services/tenant-service/go.sum services/tenant-service/

# Download dependencies
WORKDIR /app/services/tenant-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY pkg/ pkg/
COPY services/tenant-service/ services/tenant-service/

# Build the application
WORKDIR /app/services/tenant-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/tenant-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/bin/tenant-service .

# Expose ports
EXPOSE 50053 8083

CMD ["./tenant-service"]
