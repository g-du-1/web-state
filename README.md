# page-state-saver

A service for saving and restoring page states.

## Database Setup

This application uses PostgreSQL as its database. You can run the database using Docker Compose:

```bash
# Start the database
docker-compose up -d

# Stop the database
docker-compose down

# View logs
docker-compose logs postgres
```

The database will be available at:

- Host: `localhost`
- Port: `5432`
- Database: `page-state-saver`
- Username: `postgres`
- Password: `postgres`

The connection string is: `postgres://postgres:postgres@localhost:5432/page-state-saver?sslmode=disable`

## Running the Application

Set the `DATABASE_URL` environment variable (or it will use the default localhost connection):

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/page-state-saver?sslmode=disable"
go run main.go
```

## Testing

```go
go test -v ./...

go test -coverprofile=coverage.out ./...

go tool cover --html=coverage
```
