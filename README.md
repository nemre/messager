# Messager - Fast & Strong Messaging Tool 🚀

[![Go Reference](https://pkg.go.dev/badge/github.com/nemre/messager.svg)](https://pkg.go.dev/github.com/nemre/messager)
[![Go Report Card](https://goreportcard.com/badge/github.com/nemre/messager)](https://goreportcard.com/report/github.com/nemre/messager)
[![License: BSD-3-Clause](https://img.shields.io/badge/license-BSD--3--Clause-blue)](https://opensource.org/license/bsd-3-clause)
[![GitHub Release](https://img.shields.io/github/release/nemre/messager.svg)](https://github.com/nemre/messager/releases)

Messager is a high-performance, scalable messaging service built with Go. It provides a robust platform for message queuing and delivery, featuring real-time status updates, persistent storage, and reliable message processing capabilities.

![Banner](https://github.com/nemre/messager/blob/main/.github/assets/banner.png)

## 📑 Table of Contents

- [Features](#-features)
- [System Architecture](#%EF%B8%8F-system-architecture)
- [Getting Started](#-getting-started)
- [API Reference](#-api-reference)
- [Configuration](#-configuration)
- [Development](#-development)
- [Message Flow](#-message-flow)
- [Monitoring & Logging](#-monitoring--logging)
- [Security](#-security)
- [Contributing](#-contributing)
- [License](#-license)

## Big Picture

![Big Picture](https://github.com/nemre/messager/blob/main/.github/assets/big-picture.png)

The application exposes a REST API that allows users to create a message by providing content and phone data. Users can then retrieve a list of their messages and control the execution of jobs—starting or stopping them—through the same API.

Once a job is initiated, it periodically updates the status of messages that are in the pending state. These database changes are captured by Debezium and published to a Kafka topic. A Kafka consumer within the application listens for these changes and triggers an HTTP request to the corresponding client.

The metadata returned from the client is then stored in Redis. This eventual consistency architecture ensures resilience against common trade-offs such as:
	•	The database being updated, but the HTTP request not being sent.
	•	The HTTP request being sent, but the database not being updated.
	•	Duplicate requests being triggered.

Thanks to the consumer-based design, the application can scale horizontally by running multiple replicas, enabling faster message processing. The overall architecture is designed with high availability in mind.

The application follows Domain-Driven Design (DDD) principles and applies SOLID principles effectively, using appropriate abstractions and design patterns to ensure extensibility and maintainability. Additionally, by avoiding third-party libraries—including HTTP frameworks—the system reduces external dependencies and increases robustness. Any component can be replaced or modified without disrupting the integrity of other application layers.

## 🌟 Features

### Core Features
- **Message Management**
  - Create and queue messages with validation
  - Track message status (PENDING → SENT)
  - Phone number validation with international format
  - Message content validation (10-255 characters)
  
### Technical Features
- **High Performance**
  - Asynchronous message processing
  - Redis caching for sent message info
  - Kafka-based message queue
  - PostgreSQL for persistent storage
  
### Integration Features
- **Real-time CDC with Debezium**
  - Capture database changes in real-time
  - Automatic status updates via Kafka
  - Event-driven architecture
  
### Operational Features
- **Monitoring & Management**
  - Health check endpoints
  - Structured JSON logging
  - Correlation ID tracking


## 🏗️ System Architecture

### Clean Architecture Implementation
```
┌─────────────────────────────────────────────────────┐
│                   Presentation Layer                │
│   ┌─────────────┐    ┌──────────┐    ┌─────────┐    │
│   │  REST API   │    │   Jobs   │    │ Kafka   │    │
│   │  Handlers   │    │ Processor│    │Consumer │    │
│   └─────────────┘    └──────────┘    └─────────┘    │
├─────────────────────────────────────────────────────┤
│                  Application Layer                  │
│   ┌─────────────┐    ┌──────────┐    ┌─────────┐    │
│   │  Message    │    │ Business │    │ Service │    │
│   │  Services   │    │  Logic   │    │ Layer   │    │
│   └─────────────┘    └──────────┘    └─────────┘    │
├─────────────────────────────────────────────────────┤
│                    Domain Layer                     │
│   ┌─────────────┐    ┌──────────┐    ┌─────────┐    │
│   │  Entities   │    │Repository│    │ Domain  │    │
│   │  & Models   │    │Interface │    │Services │    │
│   └─────────────┘    └──────────┘    └─────────┘    │
├─────────────────────────────────────────────────────┤
│                Infrastructure Layer                 │
│┌────────┐ ┌─────┐ ┌────────┐ ┌────-─┐ ┌──────────┐  │
││Postgres│ │Redis│ │ Kafka  │ │HTTP  │ │ Logger   │  │
││  DB    │ │Cache│ │ Queue  │ │Client│ │& Monitor │  │
│└────────┘ └─────┘ └────────┘ └─────-┘ └──────────┘  │
└─────────────────────────────────────────────────────┘
```

## 🚀 Getting Started

### Prerequisites
```bash
# Check Go version (requires 1.24+)
go version

# Check Docker version
docker --version
docker-compose --version
```

### Detailed Installation Steps

1. **Clone and Setup**
   ```bash
   # Clone repository
   git clone https://github.com/nemre/messager.git
   cd messager
   
   # Create environment file
   cp .env.example .env
   
   # Initialize Go modules
   go mod tidy
   ```

2. **Configure Environment**
   ```bash
   # Edit .env file with your settings
   nano .env
   
   # Required settings:
   # - Server configuration (SERVER_*)
   # - Database credentials (POSTGRESQL_*)
   # - Redis settings (REDIS_*)
   # - Kafka configuration (KAFKA_*)
   # - Client settings (CLIENT_*)
   ```

3. **Start Services**
   ```bash
   # Start all services
   docker-compose up -d
   
   # Verify services are running
   docker-compose ps
   
   # Check logs
   docker-compose logs -f
   ```

4. **Verify Installation**
   ```bash
   # Check API health
   curl http://localhost:2025/health
   
   # Should return:
   # {"status":"green"}
   ```

## 📚 API Reference

### Create Message
```bash
curl -X POST http://localhost:2025/messages \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Your message content",
    "phone": "+905321234567"
  }'
```

### List Messages
```bash
# Get PENDING messages
curl http://localhost:2025/messages?status=PENDING

# Get SENT messages
curl http://localhost:2025/messages?status=SENT
```

### Manage Message Processing
```bash
# Start processing
curl -X POST http://localhost:2025/messages/jobs

# Stop processing
curl -X DELETE http://localhost:2025/messages/jobs
```

## 🔧 Configuration

### Environment Variables
```dotenv
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=2025
SERVER_ID_HEADER=X-Correlation-ID

# PostgreSQL Configuration
POSTGRESQL_HOST=postgres
POSTGRESQL_PORT=5432
POSTGRESQL_USER=messager
POSTGRESQL_PASSWORD=messager
POSTGRESQL_NAME=messager

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_DB=0

# Job Configuration
JOB_INTERVAL=2m

# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=messager.public.messages
KAFKA_GROUP_ID=messager

# Client Configuration
CLIENT_URL=https://api.example.com
CLIENT_TOKEN=your-token
CLIENT_TIMEOUT=5s
```

## 💻 Development

### Project Structure
```
messager/
├── application/                 # Application Services
│   └── service/
│       └── message/            # Message Service Implementation
├── domain/                     # Domain Layer
│   └── message/               
│       ├── entity.go          # Message Entity & Validation
│       ├── repository.go      # Repository Interface
│       └── service.go         # Service Interface
├── infrastructure/            # Infrastructure Layer
│   ├── client/               # HTTP Client
│   ├── config/               # Configuration
│   ├── database/             # Database Implementations
│   ├── logger/               # Structured Logger
│   ├── persistence/          # Repository Implementations
│   └── server/               # HTTP Server
└── presentation/             # Presentation Layer
    ├── consumer/             # Kafka Consumers
    ├── handler/              # HTTP Handlers
    └── job/                  # Background Jobs
```

## 📊 Monitoring & Logging

### Logging
- Structured JSON logs
- Log levels: DEBUG, INFO, WARNING, ERROR, FATAL
- Correlation ID tracking
- Separate stdout/stderr streams

## 🔐 Security

### Security Features
- TLS support
- Token-based authentication
- Input validation
- Rate limiting
- Secure defaults

### Security Policy
See [SECURITY.md](SECURITY.md) for:
- Supported versions
- Reporting vulnerabilities
- Security update policy

## 👥 Contributing

We welcome contributions! Please see:
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- [SECURITY.md](SECURITY.md)

### Development Process
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📝 License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

### Technologies
- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Apache Kafka](https://kafka.apache.org/)
- [Debezium](https://debezium.io/)
- [Docker](https://www.docker.com/)

### Libraries
- [pgx](https://github.com/jackc/pgx)
- [go-redis](https://github.com/redis/go-redis)
- [kafka-go](https://github.com/segmentio/kafka-go)
- [phonenumbers](https://github.com/nyaruka/phonenumbers)
- [uuid](https://github.com/google/uuid)
- [env](https://github.com/caarlos0/env)
