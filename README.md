# page-state-saver

## Running

```bash
docker-compose up -d

docker-compose down
```

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/page_state_saver?sslmode=disable"

go run main.go
```

## Testing

```go
go test -v ./...

go test -coverprofile=coverage.out ./...

go tool cover --html=coverage
```

## Debugging

Use the VSCode launch config.

## Userscript

Use Tampermonkey.

## DB

Connect e.g. via Docker Desktop.

```bash
psql -h localhost -p 5432 -U postgres -d page_state_saver

\c page_state_saver

TRUNCATE pagestates;
```
