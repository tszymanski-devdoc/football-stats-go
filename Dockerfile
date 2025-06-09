FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# Build the Go app
RUN go build -o football-stats-go main.go

# Use a minimal image for running
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/football-stats-go .

EXPOSE 8080

# Run the binary
CMD ["./football-stats-go"]
