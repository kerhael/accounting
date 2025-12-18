# Go accounting API

## Structure
- `cmd/api` → entry point
- `internal/router` → HTTP routes
- `internal/handler` → logic
- `pkg/` → common middlewares

## Launch the service
```bash
go run ./cmd/api
```

## Launch tests
```bash
go test ./...
```

## Access PostgreSQL : 
```bash
docker compose exec db psql -U accounting accounting
``` 