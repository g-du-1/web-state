# web-state

## Development

### Backend

```bash
docker-compose up -d
```

### Client (Browser Extension)

`cd browser-extension`
`npm install` (when first cloned)
`npm run watch`

Load the unpacked extension (dist folder) in Chrome (chrome://extensions). 
Fill in the API URL in the extension (popup) window. It should be http://localhost:8080/api/v1 when running locally.
Whitelist the sites.

Use the AutoHotkey window switcher script (window-switcher.ahk).

Set the reload keybind from manifest.json in the Keyboard Shortcuts settings on the Chrome Extensions page.

Mouse 1 + W switches to Chrome and reloads the extension.














For solving the problem

## TODO
- list down how to run extension
- Rename extension
- make public

## Usage
























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

