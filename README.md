# go-tenant-service

> Part of the SaaS Framework - Extracted from monorepo

## Description

[Add service description here]

## Features

- Feature 1
- Feature 2
- Feature 3

## Prerequisites

- Go 1.25.5+
- MongoDB 4.4+ (if applicable)
- Redis 6.0+ (if applicable)
- RabbitMQ 3.9+ (if applicable)

**For Windows users:** See [docs/WINDOWS.md](docs/WINDOWS.md) for detailed Windows development setup.

## Installation

### Linux / macOS

```bash
# Clone the repository
git clone https://github.com/vhvplatform/go-tenant-service.git
cd go-tenant-service

# Install dependencies
go mod download
```

### Windows

```powershell
# Clone the repository
git clone https://github.com/vhvplatform/go-tenant-service.git
cd go-tenant-service

# Install dependencies
.\build.ps1 deps
```

See [docs/WINDOWS.md](docs/WINDOWS.md) for comprehensive Windows setup instructions.

## Configuration

Copy the example environment file and update with your values:

```bash
cp .env.example .env
```

See [DEPENDENCIES.md](docs/DEPENDENCIES.md) for a complete list of environment variables.

## Development

### Running Locally

**Linux / macOS:**
```bash
# Run the service
make run

# Or with go run
go run cmd/main.go
```

**Windows:**
```powershell
# Run the service
.\build.ps1 run

# Or with go run
go run .\cmd\main.go
```

### Running with Docker

**Linux / macOS:**
```bash
# Build and run
make docker-build
make docker-run
```

**Windows:**
```powershell
# Build and run
.\build.ps1 docker-build
.\build.ps1 docker-run
```

### Running Tests

**Linux / macOS:**
```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

**Windows:**
```powershell
# Run all tests
.\build.ps1 test

# Run with coverage
.\build.ps1 test-coverage
```

### Linting

**Linux / macOS:**
```bash
# Run linters
make lint

# Format code
make fmt
```

**Windows:**
```powershell
# Run linters
.\build.ps1 lint

# Format code
.\build.ps1 fmt
```

## API Documentation

Interactive API documentation is available via Swagger UI:

**Swagger UI:** http://localhost:8083/swagger/index.html

For detailed API documentation, see [docs/API.md](docs/API.md).

### Generating Swagger Documentation

After making changes to API endpoints or annotations, regenerate the Swagger documentation:

**Linux / macOS:**
```bash
make swagger
```

**Windows:**
```powershell
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g cmd/main.go -o docs --parseDependency --parseInternal
```

## Deployment

See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for deployment instructions.

## Architecture

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for architecture details.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## Related Repositories

- [go-shared](https://github.com/vhvplatform/go-shared) - Shared Go libraries

## License

MIT License - see [LICENSE](LICENSE) for details

## Support

- Documentation: [Wiki](https://github.com/vhvplatform/go-tenant-service/wiki)
- Issues: [GitHub Issues](https://github.com/vhvplatform/go-tenant-service/issues)
- Discussions: [GitHub Discussions](https://github.com/vhvplatform/go-tenant-service/discussions)
