# Use the official Go image as the base image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o golang-kafka-kubernetes-microservice cmd/server/main.go

# Expose the port on which the application will run
EXPOSE 8080

# Run the application when the container starts
CMD ["./golang-kafka-kubernetes-microservice"]