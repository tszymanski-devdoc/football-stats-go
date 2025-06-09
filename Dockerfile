
# Start from the official Golang image for building
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum if present
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o football-stats-go main.go

# Use a minimal image for running
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/football-stats-go .

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./football-stats-go"]
