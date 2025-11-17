# RSS Feed Aggregator Server

A production-ready REST API server for an RSS feed aggregator built in Go. This application allows users to manage RSS feeds, follow feeds of interest, and retrieve aggregated posts from their subscriptions.

## Features

- **User Management**: Create accounts and manage API keys for authentication
- **Feed Management**: Add, list, and manage RSS feeds
- **Feed Following**: Subscribe to or unsubscribe from RSS feeds
- **Automatic Scraping**: Background worker continuously fetches and parses RSS feeds
- **Post Aggregation**: Retrieve posts from feeds you follow
- **API Key Authentication**: Secure endpoints with API key-based authentication
- **CORS Support**: Cross-origin requests enabled for frontend integration

## Technology Stack

- **Language**: Go 1.25.4
- **HTTP Router**: [Chi](https://github.com/go-chi/chi) - Fast and flexible HTTP router
- **Database**: PostgreSQL
- **SQL Tool**: [sqlc](https://sqlc.dev/) - Type-safe SQL query generation
- **Other Libraries**:
  - `github.com/lib/pq` - PostgreSQL driver
  - `github.com/google/uuid` - UUID generation
  - `github.com/joho/godotenv` - Environment variable loading
  - `github.com/go-chi/cors` - CORS middleware

## Project Structure

```
├── main.go                       # Application entry point and server setup
├── models.go                     # HTTP response models
├── json.go                       # JSON helper functions
├── rss.go                        # RSS XML parsing logic
├── scraper.go                    # Background feed scraping orchestration
├── middlewareAuth.go             # API key authentication middleware
├── handlers/
│   ├── handlerUser.go           # User creation and profile endpoints
│   ├── handlerFeed.go           # Feed creation and retrieval
│   ├── handlerFeedFollows.go    # Feed follow/unfollow management
│   ├── handlerPost.go           # Post retrieval endpoints
│   ├── handlerError.go          # Error handler
│   └── handlerReadiness.go      # Health check endpoint
├── internal/
│   ├── auth/
│   │   └── auth.go              # API key extraction and validation
│   └── db/
│       ├── db.go                # Database connection management
│       ├── models.go            # Auto-generated database models
│       └── *.sql.go             # Auto-generated SQL query methods
├── sql/
│   ├── schema/                  # Database migration files
│   └── queries/                 # SQL query definitions for code generation
├── go.mod                        # Go module dependencies
├── go.sum                        # Dependency checksums
└── .env                          # Environment configuration
```

## Prerequisites

- Go 1.25.4 or higher
- PostgreSQL 12 or higher
- Git

## Local Setup & Running

### 1. Clone the Repository

```bash
git clone <repository-url>
cd rss
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Set Up PostgreSQL Database

#### Option A: Using PostgreSQL Command Line

```bash
# Create the database
createdb rss

# Run migrations (in order)
psql rss < sql/schema/001_users.sql
psql rss < sql/schema/002_users_apikey.sql
psql rss < sql/schema/003_feed.sql
psql rss < sql/schema/004_feed_follows.sql
psql rss < sql/schema/005_feeds_lastfetchedat.sql
psql rss < sql/schema/006_posts.sql
```

#### Option B: Using Goose Migration Tool

[Goose](https://github.com/pressly/goose) is a database migration tool that manages schema changes. To use goose:

```bash
# Install goose (if not already installed)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations up
goose postgres "postgres://username:password@localhost:5432/rss?sslmode=disable" up

# Check migration status
goose postgres "postgres://username:password@localhost:5432/rss?sslmode=disable" status

# Rollback last migration
goose postgres "postgres://username:password@localhost:5432/rss?sslmode=disable" down
```

**Common Goose Command Line Options:**
- `up` - Migrate forward
- `down` - Rollback last migration
- `status` - Check which migrations have been applied
- `version` - Display current migration version
- `reset` - Rollback all migrations
- `-v` - Verbose output

For more information on goose options, run:
```bash
goose -h
```

#### Option C: Using PostgreSQL GUI Tools

Use tools like pgAdmin or DBeaver to:
1. Create a new database named `rss`
2. Execute SQL migration files in `sql/schema/` in numerical order

### 4. Configure Environment Variables

Create a `.env` file in the project root:

```env
PORT=8080
DB_URL=postgres://username:password@localhost:5432/rss?sslmode=disable
```

Replace `username` and `password` with your PostgreSQL credentials.

### 5. Run the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

**Expected output:**
```
Starting server on port 8080...
Starting scraper...
```

### 6. Verify Server is Running

```bash
curl http://localhost:8080/v1/healthz
```

Expected response: `ok`

## Building for Production

### Build Binary

```bash
go build -o rss main.go
```

### Run Binary

```bash
./rss
```

## API Endpoints

### Health Check

```
GET /v1/healthz
```

Returns: `ok`

### Users

**Create a new user:**
```
POST /v1/users
Content-Type: application/json

{
  "name": "John Doe"
}
```

Response:
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "name": "John Doe",
  "api_key": "64-character-api-key"
}
```

**Get current user (requires authentication):**
```
GET /v1/users
Authorization: ApiKey <your-api-key>
```

### Feeds

**Create a new feed (requires authentication):**
```
POST /v1/feeds
Authorization: ApiKey <your-api-key>
Content-Type: application/json

{
  "name": "Example Feed",
  "url": "https://example.com/feed.xml"
}
```

**Get all feeds:**
```
GET /v1/feeds
```

### Feed Follows

**Follow a feed (requires authentication):**
```
POST /v1/feed_follows
Authorization: ApiKey <your-api-key>
Content-Type: application/json

{
  "feed_id": "feed-uuid"
}
```

**Get your feed follows (requires authentication):**
```
GET /v1/feed_follows
Authorization: ApiKey <your-api-key>
```

**Unfollow a feed (requires authentication):**
```
DELETE /v1/feed_follows/{feed_follow_id}
Authorization: ApiKey <your-api-key>
```

### Posts

**Get posts from feeds you follow (requires authentication, limited to 10):**
```
GET /v1/posts
Authorization: ApiKey <your-api-key>
```

## How It Works

### Authentication

1. When you create a user, an API key is automatically generated (64-character SHA256 hash)
2. Include this key in the `Authorization` header for protected endpoints:
   ```
   Authorization: ApiKey your-64-character-api-key
   ```

### Feed Scraping

- A background scraper runs every 60 seconds
- It fetches feeds that haven't been updated recently
- Uses concurrent requests (10 goroutines by default) for efficiency
- Parses RSS 2.0 format and extracts post information
- Automatically prevents duplicate posts in the database
- Updates the `last_fetched_at` timestamp for each feed

### Data Flow

1. **User** creates an account → receives API key
2. **User** creates or discovers a feed → feed is stored in database
3. **User** follows a feed → creates a subscription record
4. **Scraper** (runs continuously) → fetches all feeds → stores new posts
5. **User** queries `/v1/posts` → retrieves posts from followed feeds

## Database Schema

### Users Table
- `id` (UUID) - Primary key
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `name` (string)
- `api_key` (string, 64 chars, unique)

### Feeds Table
- `id` (UUID) - Primary key
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `name` (string)
- `url` (string, unique)
- `user_id` (UUID) - Foreign key to users
- `last_fetched_at` (timestamp, nullable)

### Feed Follows Table
- `id` (UUID) - Primary key
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `user_id` (UUID) - Foreign key to users
- `feed_id` (UUID) - Foreign key to feeds
- Unique constraint on (user_id, feed_id)

### Posts Table
- `id` (UUID) - Primary key
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `title` (string)
- `description` (string)
- `url` (string, unique)
- `published_at` (timestamp, nullable)
- `feed_id` (UUID) - Foreign key to feeds

## Development

### Regenerating Database Code

If you modify SQL queries in `sql/queries/`, regenerate the Go code:

```bash
sqlc generate
```

This will update the auto-generated files in `internal/db/`.

## Troubleshooting

### Database Connection Error

**Error**: `pq: role "username" does not exist`

**Solution**: Ensure your PostgreSQL user exists and credentials in `.env` are correct.

```bash
psql -U postgres -c "CREATE USER username WITH PASSWORD 'password';"
psql -U postgres -c "ALTER USER username CREATEDB;"
```

### Port Already in Use

**Error**: `listen tcp :8080: bind: address already in use`

**Solution**: Change the port in `.env` file or kill the process using port 8080:

```bash
lsof -i :8080  # Find process ID
kill -9 <PID>  # Kill the process
```

### Migrations Not Applied

Ensure you run the migration files in order (001-006). You can verify by checking the database schema:

```bash
psql rss -c "\dt"  # List all tables
```

## Performance Considerations

- Feed scraping runs with 10 concurrent goroutines by default
- Each feed request has a 10-second timeout
- Scraper runs every 60 seconds
- Posts are limited to 10 per query to reduce memory usage
- API key authentication is validated on every protected request

## Future Enhancements

- Add pagination to posts endpoint
- Implement feed categories
- Add user preferences (feed ordering, filtering)
- Add search functionality for posts
- Implement rate limiting
- Add webhook support for feed updates
- Create admin dashboard

## Contributing

1. Create a new branch for your feature
2. Make your changes
3. Test thoroughly
4. Submit a pull request

## Support

For issues or questions, please open an issue on the repository.
