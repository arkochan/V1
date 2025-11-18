# V1 Development Environment

This repository contains a full-stack development environment with a Go backend and Next.js frontend, managed by a comprehensive CLI tool.

## Table of Contents
- [Quick Start](#quick-start)
- [Development Environment](#development-environment)
- [Database & Migrations](#database--migrations)
- [SQL Runner](#sql-runner)
- [Docker Setup](#docker-setup)

## Quick Start
### TLDR
```sh
# clone
git clone git@github.com:arkochan/V1.git 
cd V1

# generate .env
cp .env.example .env

# First make sure DB is up or run 
./run up

# Then make sure the DB is migrated or run the migrations using
./run migrate up

# Start both Go backend and Next.js frontend with TUI
./run dev
```

### Initial Setup
1. Clone the repository
```sh
git clone git@github.com:arkochan/V1.git 
```
2. Copy the example environment file: 

```sh
cp .env.example .env
```

3. Update the `.env` file with your specific configuration
4. Install dependencies for both Go and Next.js apps

### Running the Development Environment

There are multiple ways to start the development environment:

```bash
# First make sure DB is up or run
./run up

# Then make sure the DB is migrated or run the migrations using
./run migrate up

# Start both Go backend and Next.js frontend with TUI
./run dev

# Start only the Go backend
./run go

# Start only the Next.js frontend
./run next
```

The TUI (Terminal User Interface) provides a visual interface with:
- Service status indicators
- Keyboard shortcuts for management
- Real-time log output
- Service restart capabilities (g→r for Go, n→r for Next.js, r for both)

## Development Environment

### The `./run` Script

The `./run` script is a comprehensive development environment manager with the following commands:

#### Development Commands
- `./run dev` - Start both Go and Next.js services with TUI
- `./run tui` - Same as `./run dev`, starts TUI interface
- `./run go` - Start only the Go backend service
- `./run next` - Start only the Next.js frontend service

#### Docker Commands
- `./run up` - Start Docker development environment
- `./run down` - Stop Docker development environment
- `./run prod-up` - Start Docker production environment
- `./run prod-down` - Stop Docker production environment

#### Database Migration Commands
- `./run migrate up` - Apply pending migrations
- `./run migrate down` - Rollback migrations
- `./run migrate create <migration_name>` - Create new migration file
- `./run migrate drop` - Drop all migrations
- `./run migrate force <version>` - Force migration version
- `./run migrate version` - Show current migration version

### Environment Variables

The script automatically loads environment variables from a `.env` file in the project root. This file supports:
- Variable expansion with `${VAR_NAME}` syntax
- Comments with `#`
- Standard environment variable formats

## Database & Migrations

### Migration System

The project uses a database migration system with the following structure:

```
db/
├── pg/
│   └── migrations/     # SQL migration files
```

Migration commands:
```bash
# Create a new migration file
./run migrate create add_users_table

# Apply all pending migrations
./run migrate up

# Rollback migrations
./run migrate down

# Check current migration version
./run migrate version
```

When creating migrations with `./run migrate create <name>`, the system will:
- Generate sequential migration files in `/home/arkochan/Repositories/V1/db/pg/migrations`
- Use the format `<sequence_number>_<name>.(up|down).sql`
- Create both up and down migration files

### Environment Configuration

Migration commands use the following environment variables from your `.env` file:
- `MIGRATIONS` - Path to migration files directory
- `DATABASE_URL` - Database connection URL

## SQL Runner

The repository includes a powerful SQL runner script located at `./sql` for executing SQL queries against your database.

### Usage

```bash
# Interactive psql session in Docker container
./sql

# Run a SQL query in Docker container
./sql "SELECT * FROM users;"

# Run a SQL file locally
./sql -c query.sql

# Read from stdin
cat query.sql | ./sql

# Use different .env file
./sql -e /path/to/.env "SELECT 1;"

# Verbose mode
./sql -v "SELECT * FROM users;"
```

### Modes

The SQL runner supports two execution modes:

#### Attach Mode (Default)
- Executes SQL inside the Docker container
- Uses `docker exec` to connect to the PostgreSQL container
- Accesses psql directly in the container

#### Connect Mode
- Executes SQL using local psql client
- Connects to the database through exposed ports
- Uses `-c` flag to enable this mode

### Environment Variables

The SQL runner requires the following environment variables in your `.env` file:
- `POSTGRES_USER` - Database user
- `POSTGRES_PASSWORD` - Database password
- `POSTGRES_DB` - Database name
- `POSTGRES_PORT` - Database port
- `POSTGRES_CONTAINER_NAME` - Docker container name

## Docker Setup

### Development Environment

The development environment is managed through Docker Compose with configurations in:
- `docker-compose.dev.yml` - Development setup
- `docker-compose.prod.yml` - Production setup

### Starting Services

```bash
# Start development environment
./run up

# Start production environment
./run prod-up

# Stop development environment
./run down

# Stop production environment
./run prod-down
```

## Project Structure

```
.
├── run                 # Main development environment manager
├── sql                 # SQL runner script
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── apps/
│   ├── go-app/         # Go backend application
│   └── next-app/       # Next.js frontend application
├── db/
│   └── pg/
│       └── migrations/ # Database migration files
└── .env               # Environment configuration
```

## Troubleshooting

### Common Issues

1. **Services not starting**: Ensure Docker is running and all required environment variables are set
2. **Database connection errors**: Verify the PostgreSQL container is running and environment variables are correct
3. **Migration errors**: Check that the migration directory exists and has proper permissions

### Debugging

Use the verbose flag (`-v`) with the SQL runner or enable debug logging in your applications to get more detailed output.
