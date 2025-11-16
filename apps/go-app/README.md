# user-review-ingest

User review ingestion API service.

## Getting Started

### Prerequisites

- Go 1.21+
- Docker
- Docker Compose
- golangci-lint
- goose
- sqlc

### Installation

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Run database migrations: `goose -dir "migrations" postgres "postgres://user:password@localhost:5432/db?sslmode=disable" up`
4. Generate SQLC code: `sqlc generate`

### Running the application

Using Docker Compose:
```bash
docker-compose up --build
```

Locally:
```bash
go run ./cmd/api
```

### Hot Reload

Using `air`:
```bash
air
```

## API Endpoints

- `POST /v1/reviews`: Create a new review.
- `GET /v1/reviews/:id`: Get a review by ID.
- `GET /v1/reviews`: List reviews with pagination.
- `GET /health`: Health check.
