# Stage 1: Build the Go app
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tononkira cmd/server/main.go

# Stage 2: Create a small final image
FROM alpine:3.18

WORKDIR /root/

# Copy the Go binary from the builder
COPY --from=builder /app/tononkira .

# Expose the application's port
EXPOSE 8080

# Start the application
CMD ["./tononkira"]
