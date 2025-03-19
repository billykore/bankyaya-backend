# Simple Backend Banking System

The repository provides a **Simple Backend Banking System**
designed with hexagonal architecture principles and modern Go development practices.

# Hexagonal Architecture

https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749

# Project Structure

```text
├── cmd/                # Entry points for the application (e.g., HTTP server, CLI, etc.).
├── internal/           # Application core logic.
│   ├── adapter/        # Frameworks, database, and external APIs implementations.
│   │   ├── api/        # External API handlers.
│   │   ├── email/      # Email service handlers.
│   │   ├── http/       # HTTP server (Handlers, Routers).
│   │   ├── messaging/  # Message broker  implementation (RabbitMQ, Kafka, etc.).
│   │   ├── storage/    # Database implementation (Postgres, Redis, etc.).
│   │   └── ...         # Other adapters.
│   ├── application/    # Frameworks, database, and external APIs implementations.
│   │   ├── qris/       # QRIS payment domain.
│   │   ├── scheduler/  # Transaction scheduler domain.
│   │   ├── transfer/   # Intra-bank transfer domain.
│   │   ├── user/       # User domain.
│   │   └── ...         # Other domains.
│   └── core/           # Core domain logic (Entities, Value Objects, Aggregates, Interfaces).
│       ├── qris/       # QRIS payment domain.
│       ├── scheduler/  # Transaction scheduler domain.
│       ├── transfer/   # Intra-bank transfer domain.
│       ├── user/       # User domain.
│       └── ...         # Other domains.
├── pkg/                # Shared libraries or utilities.
├── script/             # Utility scripts.
├── .gitignore          # .gitignore file.
├── Dockerfile          # Application Dockerfile.
├── go.mod              # Go module definition.
├── go.sum              # Go module sum definition.
├── Makefile            # Makefile.
└── README.md           # Project documentation.
```

# Modules

Some of the open-source modules we used are:

- [Echo](https://echo.labstack.com) for http routing.
- [GORM](https://gorm.io) for database ORM.
- [Validator](https://github.com/go-playground/validator) for request validation.
- [zap](https://github.com/uber-go/zap) for logging.
- [envconfig](https://github.com/kelseyhightower/envconfig)
  and [godotenv](https://github.com/joho/godotenv) for loading env variables.
- [JWT](https://github.com/golang-jwt/jwt) for generate and validate authorization token.
- [swag](https://github.com/swaggo/swag) and [echo-swagger](https://github.com/swaggo/echo-swagger)
  for generate API documentation.
- [ecszap](https://github.com/elastic/ecs-logging-go-zap) to support ECS for zap logger.
- [Google Wire](https://github.com/google/wire) for dependency injection.
- [testify](https://github.com/stretchr/testify) for unit testing.
