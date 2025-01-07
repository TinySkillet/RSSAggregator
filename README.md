# RSS Feed Aggregator

A high-performance RSS feed aggregator built with Go that fetches and aggregates RSS feeds in real-time using concurrent processing.

## Features

- Real-time RSS feed aggregation using goroutines for concurrent processing
- RESTful API with authentication
- PostgreSQL database integration
- Database migrations using Goose
- Type-safe SQL queries using SQLC
- CORS support
- API versioning
- User management and feed subscription system

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Router:** Chi
- **SQL Tools:**
  - SQLC for type-safe query generation
  - Goose for database migrations
- **Authentication:** Custom middleware with API key support

## Environment Configuration

Create a `.env` file in the root directory with the following configuration:

```env
PORT=6000
DBCONN_STR=postgres://postgres:goapi@localhost:3000/postgres?sslmode=disable
```

### Environment Variables

- `PORT`: The port number on which the server will listen (default: 6000)
- `DBCONN_STR`: PostgreSQL connection string in the format:
  `postgres://[username]:[password]@[host]:[port]/[database]?sslmode=disable`

## API Endpoints

### Authentication

All endpoints marked with 🔒 require authentication via API key.

### User Management

- `POST /v1/user` - Create a new user
- `GET /v1/user` 🔒 - Get user details by API key
- `GET /v1/users` - Get all users

### Feed Management

- `POST /v1/feed` 🔒 - Create a new feed
- `GET /v1/feeds` - Get all available feeds
- `GET /v1/posts` 🔒 - Get posts for authenticated user

### Feed Follows

- `POST /v1/feedfollow` 🔒 - Follow a feed
- `GET /v1/feedfollows` 🔒 - Get all followed feeds
- `DELETE /v1/feedfollow/{feedID}` 🔒 - Unfollow a feed

### System

- `GET /v1/healthz` - Health check endpoint
- `GET /v1/error` - Error handling test endpoint

## Configuration

### CORS Configuration

```go
AllowedOrigins:   []string{"http://*", https://*"}
AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "UPDATE", "OPTIONS"}
AllowCredentials: false
AllowedHeaders:   []string{"*"}
ExposedHeaders:   []string{"Link"}
MaxAge:           300
```

## Installation

1. Clone the repository

```bash
git clone [repository-url]
cd rss-aggregator
```

2. Install dependencies

```bash
go mod download
```

3. Set up PostgreSQL database

4. Create and configure the `.env` file as described in the Environment Configuration section

5. Run database migrations

```bash
goose up
```

6. Start the server

```bash
go run main.go
```

## Database

### Migrations

This project uses Goose for database migrations. To run or rollback migrations, make sure you are in the sql/schema directory first.

To run migrations:

```bash
goose postgres <DB_URL> up
```

To rollback:

```bash
goose postgres <DB_URL> down
```

### Query Generation

SQLC is used to generate type-safe Go code from SQL queries. To regenerate queries:

```bash
sqlc generate
```

## Concurrency

The aggregator uses Go's goroutines to fetch multiple RSS feeds simultaneously, providing efficient performance and scalability. The concurrent processing ensures that feed updates are fetched and processed in real-time without blocking operations.

## Performance

The application is designed to scale horizontally, with concurrent feed processing and efficient database operations. The use of goroutines allows for parallel processing of multiple feeds, while connection pooling in PostgreSQL ensures efficient database operations.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
