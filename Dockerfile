# Use the official golang image as a build stage
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

# Use a minimal image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose the port on which the app will run
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
