# Go Accounting API

A REST API for managing personal accounting data including categories and outcomes, built with Go and PostgreSQL.

## Features

- **Categories Management**: Create, read, and delete expense categories
- **Outcomes Tracking**: Record financial outcomes with amounts, categories, and timestamps
- **Health Check**: API health monitoring endpoint
- **Swagger Documentation**: Interactive API documentation
- **PostgreSQL Database**: Robust data persistence with migrations
- **Docker Support**: Easy containerized deployment

## Tech Stack

- **Go 1.25.1** - Programming language
- **PostgreSQL 18** - Database
- **Docker & Docker Compose** - Containerization
- **Swagger/OpenAPI** - API documentation

## Project Structure

```
.
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── db/               # Database connection
│   ├── domain/           # Business entities and errors
│   ├── handler/          # HTTP handlers and DTOs
│   ├── infrastructure/   # Repository layer
│   ├── router/           # Route definitions
│   └── service/          # Business logic
├── pkg/                  # Shared packages (logger, etc.)
├── migrations/           # Database migrations
├── docs/                 # Generated API documentation
└── docker-compose.yml    # Docker services
```

## Prerequisites

- Docker and Docker Compose
- Go 1.25.1+ (for local development)

## Quick Start

### 1. Environment Setup

Copy the example environment file and configure your settings:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Environment
APP_ENV=development

# PostgreSQL Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=accounting
DB_PASSWORD=your_password
DB_NAME=accounting
DB_SSLMODE=disable

# Application
LOG_LEVEL=info
```

### 2. Build and Run with Docker

```bash
# Build the services
docker compose build

# Start all services (database, migration, and API)
docker compose up
```

The API will be available at `http://localhost:8080`

### 3. Run Tests

```bash
go test ./...
```

## API Documentation

### Interactive Documentation

Swagger UI is available at: `http://localhost:8080/swagger/index.html`

### API Endpoints

#### Health Check

**GET** `/api/v1/health`

Check API health status.

```bash
curl http://localhost:8080/api/v1/health
```

#### Categories

**POST** `/api/v1/categories/`

Create a new category.

```bash
curl -X POST http://localhost:8080/api/v1/categories/ \
  -H "Content-Type: application/json" \
  -d '{"label":"Food"}'
```

**GET** `/api/v1/categories/`

Retrieve all categories.

```bash
curl http://localhost:8080/api/v1/categories/
```

**GET** `/api/v1/categories/{id}`

Retrieve a specific category by ID.

```bash
curl http://localhost:8080/api/v1/categories/1
```

**DELETE** `/api/v1/categories/{id}`

Delete a category by ID.

```bash
curl -X DELETE http://localhost:8080/api/v1/categories/1
```

#### Outcomes

**POST** `/api/v1/outcomes/`

Create a new outcome.

```bash
curl -X POST http://localhost:8080/api/v1/outcomes/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Restaurant",
    "amount": 1999,
    "categoryId": 1,
    "createdAt": "2026-01-01T00:00:00Z"
  }'
```

**GET** `/api/v1/outcomes/`

Retrieve all outcomes.

```bash
curl http://localhost:8080/api/v1/outcomes/
curl http://localhost:8080/api/v1/outcomes/?from=2025-01-01T00:00:00Z&to=2026-01-01T00:00:00Z
```

## Database Management

### Access PostgreSQL

```bash
docker compose exec db psql -U accounting accounting
```

### Run Migrations Manually

```bash
# Apply migrations
docker compose run --rm migrate up

# Rollback last migration
docker compose run --rm migrate down 1
```

### View Migration Status

```bash
docker compose run --rm migrate version
```

## Development

### Local Development Setup

1. Install Go 1.25.1+
2. Set up PostgreSQL (or use Docker)
3. Copy and configure `.env`
4. Run migrations
5. Start the API:

```bash
go run cmd/api/main.go
```

### Generate Swagger Documentation

```bash
swag init -g cmd/api/main.go
```

## Error Responses

The API returns standard HTTP status codes:

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `500` - Internal Server Error

Error response format:
```json
{
  "error": "error message description"
}
```
