# web-state

## Stack

- Go
- PostgreSQL
- testcontainers
- Chrome Extension
- Termux Android, Edge

## Development

### Backend

```bash
docker-compose up -d
```

### Client (Browser Extension)

- `cd browser-extension`
- `npm install` (when first cloned)
- `npm run watch`

Load the unpacked extension (dist folder) in Chrome (chrome://extensions).
Fill in the API URL in the extension (popup) window. It should be http://localhost:8080/api/v1 when running locally.
Whitelist the sites.

Use the AutoHotkey window switcher script (window-switcher.ahk).

Set the reload keybind from manifest.json in the Keyboard Shortcuts settings on the Chrome Extensions page.

Mouse 1 + W switches to Chrome and reloads the extension.

## Testing

```go
go test -v./...

go test -coverprofile= coverage.out./...

go tool cover --html = coverage
```
