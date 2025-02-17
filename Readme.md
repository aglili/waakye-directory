# Waakye Directory

An golang api built with Gin

## Prerequisites

Before you begin, ensure you have the following installed:
- Go (latest version)
- Docker and Docker Compose
- PostgreSQL
- `migrate` CLI tool for database migrations

## Environment Setup

Create a `.env` file in the root directory with the following variables:

```env
DB_USER=your_username
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=waakye_directory
PORT=":8080"
```

## Getting Started

### Building and Running

1. Build the project:
```bash
make build
```

2. Run the application:
```bash
make run
```

### Development Commands

Format the code:
```bash
make fmt
```

Run linting:
```bash
make lint
```

Run tests:
```bash
make test
```

Clean build artifacts:
```bash
make clean
```

## Docker Operations

Start the containers:
```bash
make docker-up
```

Stop the containers:
```bash
make docker-down
```

Restart the containers:
```bash
make docker-restart
```

View container logs:
```bash
make docker-logs
```

## Database Migrations

Create a new migration:
```bash
make migrate-create name=your_migration_name
```

Apply all migrations:
```bash
make migrate-up
```

Revert migrations:
```bash
make migrate-down
```

Force a specific migration version:
```bash
make migrate-force version=version_number
```

## Available Make Commands

Run `make help` to see all available commands:

```bash
make help
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
