# web-state

## TODO
- list down how to run extension
- make public

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

## Extension Dev

Use the AHK window switcher.

`cd browser-extension`
`npm run watch`

Set the reload keybind from manifest.json in extension settings.

Mouse 1 + W switches to Chrome and reloads the extension.
