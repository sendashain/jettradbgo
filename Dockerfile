# Build stage
FROM golang:1.19-alpine AS builder

# Install git (needed for go mod download)
RUN apk add --no-cache git

WORKDIR /app
COPY . .

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o multimodel-db .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/multimodel-db .

# Expose the default port
EXPOSE 8080

# Run the binary
CMD ["./multimodel-db"]