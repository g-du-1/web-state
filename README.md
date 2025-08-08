# page-state-saver

## Running

```bash
docker-compose up -d

docker-compose down
```

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

## Userscript

Use Tampermonkey.
