# Chatheon

[![Go](https://github.com/chrikar/chatheon/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/chatheon/actions)
[![codecov](https://codecov.io/gh/chrikar/chatheon/graph/badge.svg?token=3D0K04DH8Q)](https://codecov.io/gh/chrikar/chatheon)

A modern, Go-powered chat application built with Hexagonal architecture and full test coverage.

## Project Structure

```bash
chatheon/
├── cmd/
│   └── chatheon/
│       └── main.go           # Application entry point
├── application/
│   ├── ports/                # Application ports (interfaces)
│   └── user_service.go
│   └── user_service_test.go
├── domain/                   # Domain models
├── adapters/                 # Adapters (database, etc.)
├── internal/
│   ├── auth/                 # JWT helper and middleware
│   └── config/               # Configuration loader
├── migrations/               # SQL migration scripts
├── tests/                    # External test helpers
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## How to Run

### Run unit tests
```bash
make test
```

### Check test coverage
```bash
make coverage
```

### Build and run locally
```bash
make build
make run
```

### Run with Docker
```bash
make docker-up
```

### Stop Docker services
```bash
make docker-down
```

### Run linter
```bash
make lint
```

## Features

- ✅ Hexagonal architecture (ports & adapters)
- ✅ JWT authentication
- ✅ User registration & login
- ✅ Full unit test coverage with GitHub Actions CI
- ✅ Test coverage reports uploaded to Codecov
- ✅ Clean, idiomatic Go project structure
- ✅ Dockerized development environment

## Roadmap

- 🔒 Secure endpoints with JWT middleware (✅ done)
- 💬 Add messaging service
- 📦 Persist messages to database
- 🚀 Add WebSocket support
- 🧩 Extend to conversations and channels
- 🧩 Improve CI with coverage thresholds

## Contributing

Pull requests are welcome! Please open an issue first to discuss changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
