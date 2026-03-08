# Auth Service

A comprehensive authentication and authorization microservice built with Go, Gin, and PostgreSQL, featuring JWT tokens, RBAC (Role-Based Access Control), and event-driven architecture with Kafka.

## Features

- **Authentication**: JWT-based authentication with access and refresh tokens
- **Authorization**: Comprehensive RBAC system with roles and permissions
- **User Management**: Complete user profile management
- **Security**: Password hashing, input validation, security headers
- **Event-Driven**: Kafka integration for user events (registration, updates, role changes)
- **Caching**: Redis integration for performance optimization
- **Testing**: Comprehensive unit tests with mocks
- **Docker**: Containerized deployment with health checks

## Architecture

### Tech Stack
- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Message Queue**: Kafka
- **Testing**: Testify with mocks
- **Containerization**: Docker

### Project Structure
```
auth-service/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── controllers/     # HTTP request handlers
│   ├── repositories/    # Data access layer
│   └── services/        # Business logic layer
├── migrations/          # Database migrations
├── shared/              # Shared packages (linked)
└── tests/               # Unit and integration tests
```

## API Endpoints

### Public Endpoints
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Token refresh

### Protected Endpoints (Require Authentication)
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/profile` - Get user profile
- `PUT /api/v1/auth/profile` - Update user profile
- `POST /api/v1/auth/change-password` - Change password

### Admin Endpoints (Require Admin Role)
- `POST /api/v1/auth/users/:user_id/roles` - Assign role to user
- `DELETE /api/v1/auth/users/:user_id/roles` - Revoke role from user

## RBAC System

### Default Roles
- **admin**: Full system access
- **user**: Basic user access
- **seller**: Product management access
- **customer**: Shopping access
- **moderator**: Content moderation access
- **support**: Customer support access

### Default Permissions
- **User Management**: user:view, user:create, user:update, user:delete
- **Product Management**: product:view, product:create, product:update, product:delete
- **Order Management**: order:view, order:create, order:update, order:delete
- **Payment Management**: payment:view, payment:create, payment:refund
- **Admin Access**: admin:view, admin:manage

## Configuration

### Environment Variables
```bash
# Service Configuration
SERVICE_NAME=auth-service
SERVICE_PORT=8081
ENVIRONMENT=development

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_service_db
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=1

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC_USER_REGISTERED=user.registered
KAFKA_TOPIC_USER_UPDATED=user.updated
KAFKA_TOPIC_USER_DELETED=user.deleted
KAFKA_TOPIC_ROLE_ASSIGNED=role.assigned
KAFKA_TOPIC_ROLE_REVOKED=role.revoked
```

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Redis 6+
- Kafka 2.8+
- Docker and Docker Compose (optional)

### Local Development
1. Clone the repository
2. Copy `.env.example` to `.env` and configure
3. Run database migrations
4. Start the service

```bash
# Copy environment configuration
cp .env.example .env

# Run migrations (requires PostgreSQL)
# Use your preferred migration tool or run SQL scripts manually

# Install dependencies
go mod download

# Run the service
go run cmd/server/main.go
```

### Docker Development
```bash
# Build and run with Docker Compose
docker-compose up -d

# Check service health
curl http://localhost:8081/health
```

## Testing

### Unit Tests
```bash
go test ./internal/services/... -v
```

### Integration Tests
```bash
go test ./tests/... -v
```

### Test Coverage
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Security Features

- **Password Security**: Bcrypt hashing with default cost
- **JWT Security**: HS256 signing with configurable secret
- **Input Validation**: Comprehensive request validation
- **Rate Limiting**: Ready for rate limiting middleware
- **CORS**: Configurable CORS policies
- **Security Headers**: X-Content-Type-Options, X-Frame-Options, etc.
- **SQL Injection Prevention**: Parameterized queries with GORM
- **Data Sanitization**: Input sanitization middleware

## Event System

The service publishes events to Kafka for:
- User registration
- User profile updates
- Role assignments
- Role revocations

Event format:
```json
{
  "type": "user.registered",
  "data": {
    "user_id": "uuid",
    "email": "user@example.com",
    "username": "username",
    "created_at": "timestamp"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Monitoring and Health Checks

- Health check endpoint: `GET /health`
- Service status and dependencies
- Database connectivity check
- Redis connectivity check
- Kafka connectivity check

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new features
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License.