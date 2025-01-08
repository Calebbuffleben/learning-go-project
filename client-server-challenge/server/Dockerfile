# Use the official Golang image as a builder
FROM golang:1.20 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
RUN go build -o main .

# Create the final image
FROM alpine:latest

# Install SQLite for Alpine
RUN apk --no-cache add sqlite

# Set the working directory
WORKDIR /app

# Copy the built binary and SQLite database
COPY --from=builder /app/main .
COPY ./db/database.db ./db/

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
