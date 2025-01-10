# Use the official Go image as a base
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY server/go.mod server/go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire server directory into the container
COPY server/ .

# Ensure SQLite database directory exists
RUN mkdir -p db && touch db/database.db

# Expose port 8080 for the Go application
EXPOSE 8080

# Run the Go application
CMD ["go", "run", "main.go"]