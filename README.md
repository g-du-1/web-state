# page-state-saver

## Running

```bash
docker-compose up -d
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
