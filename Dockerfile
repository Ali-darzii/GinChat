# Use the official Golang Alpine image as the base image (lightweight)
FROM golang:1.22-alpine

# Install any dependencies needed (e.g., git for go modules)
RUN apk add --no-cache git

# Set the working directory to /app
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all Go dependencies
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application
CMD ["/cmd/main"]
