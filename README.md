# page-state-saver

## Usage

## Dev

```bash
docker-compose up -d
```

## Testing

```go
go test -v./...

go test -coverprofile= coverage.out./...

go tool cover --html = coverage
```

## DB

Connect e.g. via Docker Desktop.

```bash
psql -h localhost -p 5432 -U postgres -d page_state_saver

\c page_state_saver

TRUNCATE pagestates;
```

## Stack

- Go
- PostgreSQL
- testcontainers
- Chrome Extension (JSDoc types)
- Termux Android, Edge