# go-tenant-service

> Part of the SaaS Framework - Tenant Management Service

## Description

This repository contains the tenant management service with multiple components:

- **server/** - Golang backend microservice
- **client/** - ReactJS frontend microservice
- **flutter/** - Flutter mobile application
- **docs/** - Project documentation

## Getting Started

### Server (Golang Backend)

See [server/README.md](server/README.md) for detailed instructions on setting up and running the backend service.

### Client (ReactJS Frontend)

See [client/README.md](client/README.md) for detailed instructions on setting up and running the frontend application.

### Flutter (Mobile App)

See [flutter/README.md](flutter/README.md) for detailed instructions on setting up and running the mobile application.

## Documentation

Complete project documentation is available in the [docs/](docs/) directory:

- [CHANGELOG.md](docs/CHANGELOG.md) - Version history and changes
- [CONTRIBUTING.md](docs/CONTRIBUTING.md) - Contribution guidelines
- [DEPENDENCIES.md](docs/DEPENDENCIES.md) - Dependencies and environment variables
- [WINDOWS.md](docs/WINDOWS.md) - Windows development setup
- [WINDOWS_COMPATIBILITY_SUMMARY.md](docs/WINDOWS_COMPATIBILITY_SUMMARY.md) - Windows compatibility notes

## Repository Structure

```
.
├── client/          # ReactJS frontend microservice
├── server/          # Golang backend microservice
├── flutter/         # Flutter mobile application
├── docs/            # Project documentation
└── README.md        # This file
```

## Prerequisites

- For backend: Go 1.25.5+, MongoDB 4.4+, Redis 6.0+, RabbitMQ 3.9+
- For frontend: Node.js 16+, npm/yarn
- For mobile: Flutter 3.0+, Android SDK, Xcode (for iOS)

## Quick Start

```bash
# Clone the repository
git clone https://github.com/vhvplatform/go-tenant-service.git
cd go-tenant-service

# Setup backend
cd server
go mod download
make run

# Setup frontend (in another terminal)
cd client
npm install
npm start

# Setup mobile app
cd flutter
flutter pub get
flutter run
```

## Contributing

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for contribution guidelines.

## Related Repositories

- [go-shared](https://github.com/vhvplatform/go-shared) - Shared Go libraries

## License

MIT License - see [LICENSE](LICENSE) for details

## Support

- Documentation: [Wiki](https://github.com/vhvplatform/go-tenant-service/wiki)
- Issues: [GitHub Issues](https://github.com/vhvplatform/go-tenant-service/issues)
- Discussions: [GitHub Discussions](https://github.com/vhvplatform/go-tenant-service/discussions)
