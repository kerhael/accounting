# Go accounting API

## Structure
- `cmd/api` → entry point
- `internal/router` → HTTP routes
- `internal/handler` → logic
- `pkg/` → common middlewares

## Use the service
### Installation
```bash
docker compose build
```

### Launch
```bash
docker compose up
```

## Test the service
```bash
go test ./...
```

## Database
### Access PostgreSQL
```bash
docker compose exec db psql -U accounting accounting
``` 

### Cancel migration
```bash
docker compose run --rm migrate down 1
```

## Routes

Swagger documentation is available at `http://localhost:8080/swagger/index.html`

* healthcheck: 
```bash
 curl http://localhost:8080/api/v1/health   -H "Content-Type: application/json"
 ```
* create category: 
```bash
 curl -X POST http://localhost:8080/api/v1/categories   -H "Content-Type: application/json"   -d '{"label":"Food"}'
 ```
* retrieve a category by ID: 
```bash
 curl http://localhost:8080/api/v1/categories/1   -H "Content-Type: application/json"
 ```