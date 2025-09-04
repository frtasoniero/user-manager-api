# User Management API

A comprehensive REST API for user management built with Go, featuring Clean Architecture principles, MongoDB integration, and complete Swagger documentation.

## 🏗️ Architecture Overview

This project follows **Clean Architecture** (Hexagonal Architecture) principles with clear separation of concerns:

```
├── cmd/api/                    # Application entry point
├── internal/                   # Private application code
│   ├── core/                  # Business logic layer
│   │   ├── domain/           # Entities and business rules
│   │   ├── ports/            # Interfaces (contracts)
│   │   └── usecase/          # Application business logic
│   └── adapters/             # External layer
│       ├── handler/http/     # HTTP handlers (Gin)
│       └── repository/       # Data access (MongoDB)
├── pkg/                       # Reusable packages
│   ├── security/             # Password hashing utilities
│   └── logger/               # Logging utilities
├── database/                  # Database connection
├── routes/                    # Route definitions
├── scripts/                   # Database initialization
└── docs/                     # Generated Swagger documentation
```

### Architecture Layers:

1. **Domain Layer** (`internal/core/domain/`): Contains business entities and rules
2. **Application Layer** (`internal/core/usecase/`): Contains application-specific business logic
3. **Interface Layer** (`internal/core/ports/`): Defines contracts between layers
4. **Infrastructure Layer** (`internal/adapters/`): External concerns (HTTP, database, etc.)

## 🚀 Features

### Core Functionality
- ✅ **User Registration** with comprehensive validation
- ✅ **Advanced User Filtering** with pagination, search, and sorting
- ✅ **Secure Password Hashing** using bcrypt
- ✅ **UUID-based User IDs** for security and portability
- ✅ **MongoDB Integration** with schema validation
- ✅ **Comprehensive Input Validation** using Gin validator

### Technical Features
- ✅ **Clean Architecture** with dependency injection
- ✅ **Interactive Swagger Documentation** 
- ✅ **Hot Reload Development** with Air
- ✅ **Docker Compose** for database setup
- ✅ **Comprehensive Makefile** for development workflow
- ✅ **Environment-based Configuration**
- ✅ **Structured Logging** support
- ✅ **Code Quality Tools** (linting, formatting, vetting)

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/health` | Health check |
| `POST` | `/api/v1/users/register` | User registration |
| `GET` | `/api/v1/users` | Get users with filtering |
| `GET` | `/api/v1/users/{id}` | Get user by UUID |
| `GET` | `/swagger/index.html` | Interactive API documentation |

### Advanced Filtering Features
- **Pagination**: `?page=1&page_size=10`
- **Search**: `?search=john` (searches email, first_name, last_name)
- **Sorting**: `?sort=email&order=desc`
- **Field Selection**: `?fields=email,profile.first_name,created_at`

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.25.0
- **Web Framework**: Gin Gonic v1.10.1
- **Database**: MongoDB 7.0
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **UUID Generation**: Google UUID v1.6.0

### Development Tools
- **Hot Reload**: Air v1.x
- **Documentation**: Swagger/OpenAPI 3.0
- **Containerization**: Docker & Docker Compose
- **Code Quality**: golangci-lint, gofmt, go vet
- **Environment Management**: godotenv

### Database
- **Primary Database**: MongoDB 7.0
- **Schema Validation**: JSON Schema validation
- **Indexes**: Optimized for email, name searches, and sorting
- **Connection**: Official MongoDB Go Driver v1.17.4

## 📋 Prerequisites

- **Go**: 1.25.0 or later
- **Docker & Docker Compose**: Latest stable version
- **Make**: For using the provided Makefile commands

## 🚀 Quick Start

### 1. Clone and Setup
```bash
# Clone the repository
git clone https://github.com/your-username/user-management-api.git
cd user-management-api

# Complete project setup (installs tools, creates .env, starts database, generates docs)
make setup
```

### 2. Start Development Server
```bash
# Start with hot reload and Swagger documentation
make dev

# The API will be available at:
# - API: http://localhost:8080
# - Swagger UI: http://localhost:8080/swagger/index.html
```

### 3. Test the API
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Or use the provided HTTP test file
# Open api-tests.http in VS Code with REST Client extension
```

## 🔧 Development Commands

The project includes a comprehensive Makefile with the following commands:

### Setup & Dependencies
```bash
make setup          # Complete project setup
make install-tools  # Install development tools (air, swag, golangci-lint)
make deps          # Download and tidy Go dependencies
```

### Development
```bash
make dev           # Start with hot reload + Swagger generation
make run           # Run without hot reload
make build         # Build development binary
make build-prod    # Build production binary
```

### Database Operations
```bash
make docker-up     # Start MongoDB container
make docker-down   # Stop containers
make db-connect    # Connect to MongoDB shell
make db-backup     # Create database backup
make db-restore BACKUP=backup-name  # Restore from backup
```

### Code Quality
```bash
make check         # Run all checks (format, vet, lint)
make format        # Format code with gofmt and goimports
make lint          # Run golangci-lint
make vet           # Run go vet
```

### Documentation
```bash
make swagger       # Generate Swagger documentation
make swagger-fmt   # Format Swagger comments
make swagger-clean # Clean generated docs
```

### Testing & Utilities
```bash
make test-api      # Test API endpoints (requires running server)
make status        # Show project status
make clean         # Clean build artifacts
```

## 📝 API Documentation

### Interactive Documentation
Access the complete interactive API documentation at:
**http://localhost:8080/swagger/index.html**

### User Registration Example
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securePassword123",
    "profile": {
      "first_name": "John",
      "last_name": "Doe",
      "address": {
        "street": "123 Main St",
        "city": "New York",
        "state": "NY",
        "country": "USA",
        "zip_code": "10001"
      },
      "phone": "+1-555-123-4567",
      "birthdate": "1990-05-15",
      "nin": "123-45-6789"
    }
  }'
```

### Advanced User Filtering
```bash
# Search users with pagination and sorting
curl "http://localhost:8080/api/v1/users?page=1&page_size=5&search=john&sort=email&order=desc"

# Get specific fields only
curl "http://localhost:8080/api/v1/users?fields=email,profile.first_name,profile.last_name"
```

## ⚙️ Configuration

### Environment Variables
The application uses environment variables defined in `.env`:

```env
# MongoDB Configuration
MONGODB_URI=mongodb://api_user:api_password@localhost:27017/user_management?authSource=user_management
MONGODB_DB_NAME=user_management

# Server Configuration
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info

# JWT Configuration (for future authentication)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Environment
ENV=development
```

### Database Schema
The MongoDB collection uses strict schema validation:

- **Required fields**: `_id`, `email`, `password_hash`, `profile`, `created_at`, `updated_at`
- **Unique constraints**: `email`, `nin` (National Identification Number)
- **Indexed fields**: `email`, `profile.first_name`, `profile.last_name`, `created_at`
- **UUID format**: String-based UUIDs for better portability

## 🧪 Testing

The project includes comprehensive HTTP tests in `api-tests.http`:

### Test Categories
1. **User Registration**: Valid, invalid, edge cases
2. **User Filtering**: Pagination, search, sorting, field selection
3. **Error Handling**: Invalid inputs, duplicates, validation errors
4. **International Support**: Multi-language and locale testing

### Running Tests
```bash
# Using VS Code REST Client extension
# Open api-tests.http and click "Send Request"

# Or using curl commands
make test-api
```

## 🐳 Docker & Database

### MongoDB Setup
The project uses Docker Compose for MongoDB:

```yaml
# docker-compose.yml
services:
  mongodb:
    image: mongo:7.0
    container_name: user-management-mongodb
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password123
      MONGO_INITDB_DATABASE: user_management
```

### Database Features
- **Automatic Initialization**: Schema validation and indexes
- **User Management**: Application-specific database user
- **Backup & Restore**: Built-in backup/restore commands
- **Data Persistence**: Volume mounting for data persistence

## 📊 Monitoring & Logging

### Health Check
- **Endpoint**: `GET /api/v1/health`
- **Response**: API status and version information

### Logging
- Structured logging support (ready for implementation)
- Request/response logging via Gin middleware
- Error tracking and debugging

## 🔐 Security Features

### Implemented
- **Password Security**: bcrypt hashing with salt
- **Input Validation**: Comprehensive request validation
- **UUID IDs**: Non-predictable user identifiers
- **Schema Validation**: MongoDB-level data validation

### Future Enhancements
- JWT authentication (configuration ready)
- Rate limiting
- CORS configuration
- Request sanitization

## 🚧 Future Improvements

### Short-term
1. **Authentication & Authorization**
   - JWT-based authentication
   - Role-based access control (RBAC)
   - Password reset functionality

2. **Enhanced User Management**
   - User profile updates
   - Account deactivation/deletion
   - Email verification

3. **API Enhancements**
   - Rate limiting
   - Response caching
   - Request/response compression

### Long-term
1. **Scalability**
   - Database connection pooling optimization
   - Horizontal scaling support
   - Caching layer (Redis)

2. **Monitoring & Observability**
   - Structured logging implementation
   - Metrics collection (Prometheus)
   - Health checks enhancement
   - Distributed tracing

3. **Testing & Quality**
   - Unit test coverage (currently at 0%)
   - Integration tests
   - Load testing
   - Security testing (OWASP)

4. **DevOps & Deployment**
   - CI/CD pipeline
   - Production Docker image
   - Kubernetes manifests
   - Environment-specific configurations

5. **Business Features**
   - User activity logging
   - Advanced search capabilities
   - Bulk operations
   - Data export/import

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

For questions, issues, or contributions, please:
- Open an issue on GitHub
- Check the interactive API documentation at `/swagger/index.html`
- Review the comprehensive test cases in `api-tests.http`

---

**Project Status**: ⚡ Active Development
**API Version**: v1.0
**Go Version**: 1.25.0
**Last Updated**: September 2024