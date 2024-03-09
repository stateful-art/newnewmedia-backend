# Use the official Golang image as a base image
FROM golang:1.21.0-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Set the environment variables
ENV ENVIRONMENT="dev"

# Build the Go app
RUN go build -o main .

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
