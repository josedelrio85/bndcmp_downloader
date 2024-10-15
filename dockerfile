# Use the official Go image as the base image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o ./cmd/api/main ./cmd/api/main.go

# Set the entry point to run the application
ENTRYPOINT ["./cmd/api/main"]

# Command to run the application (can be overridden)
CMD []