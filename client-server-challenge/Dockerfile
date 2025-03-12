# Use the official Go image as a base
FROM golang:1.21

# Set the working directory
WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@v1.49.0

# Copy go.mod and go.sum first to cache dependencies
COPY server/go.mod server/go.sum ./

# Download dependencies
RUN go mod tidy

# Create the db directory
RUN mkdir -p db

# Expose port 8080 for the Go application
EXPOSE 8080

# Run air for hot reloading
CMD ["air", "-c", ".air.toml"]