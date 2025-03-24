# Use the official Go image as a base
FROM golang:1.21 AS builder

# Set the working directory
WORKDIR /app

# Install necessary tools
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Create db directory
RUN mkdir -p /app/db

# Copy server files and build
COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/main.go .
RUN go build -o /app/server main.go

# Copy client files and build
RUN mkdir -p /app/client
COPY client/go.mod /app/client/
RUN cd /app/client && go mod tidy
COPY client/client.go /app/client/
WORKDIR /app/client
RUN go build -o /app/client_bin client.go
WORKDIR /app

# Set permissions
RUN chmod +x /app/server /app/client_bin

# Create a startup script
RUN echo '#!/bin/bash\necho "Starting server..."\n/app/server & SERVER_PID=$!\necho "Server started with PID $SERVER_PID"\nsleep 5\necho "Starting client..."\ncd /app/client && SERVER_URL=http://localhost:8080 /app/client_bin & CLIENT_PID=$!\necho "Client started with PID $CLIENT_PID"\nwait' > /app/start.sh && \
    chmod +x /app/start.sh

# Expose port
EXPOSE 8080

# Run both applications
CMD ["/app/start.sh"]