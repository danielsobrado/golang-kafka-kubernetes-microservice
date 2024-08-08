# Use the official Golang 1.22 image as the base image (Debian-based)
FROM golang:1.22 AS builder

# Install build essentials and Kafka dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    librdkafka-dev \
    git \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies with retries and fallback
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download || \
    (go env -w GOPROXY=direct && go mod download) || \
    (git config --global url."https://".insteadOf git:// && \
     go mod download -x)

# Copy the source code
COPY . .

# Build the application
# Note: Go 1.22 has improved build caching, so we don't need to explicitly disable it
RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/server/main.go

# Use a smaller base image for the final stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/application.properties .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]