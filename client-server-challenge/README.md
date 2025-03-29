# Currency Exchange Rate Application

This is a Go-based client-server application that fetches, stores, and displays currency exchange rates for USD to BRL (US Dollar to Brazilian Real). The application consists of a server component that fetches data from an external API and a client component that retrieves this data from the server.

## Architecture

### Server Component

The server component (`server/main.go`) is responsible for:

1. Fetching real-time USD to BRL exchange rate data from an external API (`https://economia.awesomeapi.com.br/json/last/USD-BRL`)
2. Storing this data in a SQLite database
3. Providing an HTTP endpoint for clients to retrieve the current exchange rate

Key features:
- Built with pure Go using the standard library
- Implements strict timeout handling (200ms for API requests, 10ms for database operations)
- Provides robust error handling with fallbacks for API and database failures
- Uses SQLite for persistent storage

### Client Component

The client component (`client/client.go`) is responsible for:

1. Periodically fetching the current USD to BRL exchange rate from the server
2. Saving the exchange rate to a text file (`cotacao.txt`)
3. Displaying the information to the user

Key features:
- Simple, focused functionality
- Configurable server URL via environment variables
- Implements error handling for network failures

## Database Schema

The server uses a SQLite database with the following schema:

```sql
CREATE TABLE IF NOT EXISTS currency_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL,
    codein TEXT NOT NULL,
    name TEXT NOT NULL,
    high REAL NOT NULL,
    low REAL NOT NULL,
    varBid REAL NOT NULL,
    pctChange REAL NOT NULL,
    bid REAL NOT NULL,
    ask REAL NOT NULL,
    timestamp INTEGER NOT NULL,
    create_date TEXT NOT NULL
);
```

## API Endpoints

### Server API

- `GET /cotacao`: Returns the current USD to BRL exchange rate as JSON
  - Response format: `{"bid": "5.7314"}`

## Running the Application

### Using Docker Compose

The easiest way to run the application is using Docker Compose:

```bash
docker-compose up
```

This will:
1. Build the Docker image with both server and client components
2. Start the container with both services running
3. Expose the server on port 8080
4. Mount volumes for persistent database storage and client output

### Configuration

The following environment variables can be set to configure the application:

- `SERVER_URL`: The URL where the server is running (default: `http://localhost:8080`)
- `GO_ENV`: Environment setting (`development` or `production`)

## Docker Setup

The application uses a Docker-based setup with:

1. A single Dockerfile that builds both client and server components
2. A custom startup script that launches both applications in the correct order
3. Volume mounts for database persistence and client output
4. Health checks to ensure the application is running correctly

## Timeout Constraints

The application implements strict timeout constraints:

- The server allows a maximum of 200ms for API requests to the external currency service
- Database operations have a strict timeout of 10ms

These constraints ensure the application remains responsive even when external services or the database experience slowdowns.

## Output

The client produces a file named `cotacao.txt` in the client directory, containing:

```
Dólar: 5.7314
```

The value is updated every 10 seconds from the latest exchange rate data.

## Error Handling

Both client and server components implement robust error handling:

- The server continues to operate and return data even if database operations fail
- The client retries operations if the server is temporarily unavailable
- All errors are properly logged for troubleshooting 