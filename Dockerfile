# Use the official Go image as the base image
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY backend/go.mod ./
COPY backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY backend .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main main.go

# Use a minimal image for the final stage
FROM alpine:3.14

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage into /app
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["/main"]
