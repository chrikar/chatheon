# Chatheon

[![Go](https://github.com/chrikar/chatheon/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/chatheon/actions)
[![codecov](https://codecov.io/gh/chrikar/chatheon/graph/badge.svg?token=3D0K04DH8Q)](https://codecov.io/gh/chrikar/chatheon)

A modern, Go-powered chat application built with Hexagonal architecture and full test coverage.

## Table of Contents
- [Chatheon](#chatheon)
  - [Table of Contents](#table-of-contents)
  - [Project Structure](#project-structure)
  - [Developer setup](#developer-setup)
  - [Makefile commands](#makefile-commands)
    - [Build the application](#build-the-application)
    - [Run the application](#run-the-application)
    - [Run unit tests](#run-unit-tests)
    - [Check test coverage](#check-test-coverage)
    - [Build and run locally](#build-and-run-locally)
    - [Run with Docker](#run-with-docker)
    - [Stop Docker services](#stop-docker-services)
    - [Run linter](#run-linter)
  - [Features](#features)
  - [Roadmap](#roadmap)
  - [Testing](#testing)
    - [curl commands](#curl-commands)
      - [Register a new user](#register-a-new-user)
      - [Login to get JWT token](#login-to-get-jwt-token)
      - [Send a message (use the preiously obtained JWT token)](#send-a-message-use-the-preiously-obtained-jwt-token)
      - [Get messages](#get-messages)
      - [Get messages with pagination](#get-messages-with-pagination)
  - [Contributing](#contributing)
  - [License](#license)

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

## Developer setup
To set up pre-commit hooks for linting and formatting, install pre-commit hooks:

```bash
pre-commit install
```

Then, you can run the pre-commit hooks manually with:
```bash
pre-commit run --all-files
```

## Makefile commands
The Makefile provides a convenient way to run common tasks. Here are some of the commands you can use:

### Build the application
```bash
make build
```

### Run the application
```bash
make run
```

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

## Testing
### curl commands

#### Register a new user
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password"}'
```

#### Login to get JWT token
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password"}'
```

Response:
```json
{
  "token": "your-jwt-token-here"
}
```

#### Send a message (use the preiously obtained JWT token)
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token-here" \
  -d '{"receiver_id":"user2","content":"Hello, user2!"}'
```

Expected result:
```bash
HTTP/1.1 201 Created
```

#### Get messages
```bash
curl -X GET http://localhost:8080/messages \
  -H "Authorization: Bearer your-jwt-token-here"
```

#### Get messages with pagination
```bash
curl -X GET "http://localhost:8080/messages?limit=5&offset=0" -H "Authorization: Bearer your-token"
```

Response:
```json
[
  {
    "ID": "some-uuid",
    "SenderID": "user1",
    "ReceiverID": "user2",
    "Content": "Hello, user2!"
  }
]
```

## Contributing

Pull requests are welcome! Please open an issue first to discuss changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
