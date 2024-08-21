# Stage 1: Build the Go binary with CGO enabled
FROM golang:1.21-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Enable CGO
ENV CGO_ENABLED=1

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app with CGO enabled
RUN go build -o svc .

# Stage 2: Run the Go binary in a minimal container
FROM debian:bullseye-slim

# Install SQLite to ensure runtime compatibility
RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/svc .

# Expose the application port (change if necessary)
EXPOSE 8000

# Command to run the executable
CMD ["./svc"]
