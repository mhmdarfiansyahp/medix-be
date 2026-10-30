# Medix Backend System Architecture

## Overview

This document describes the architecture of the Medix Backend (medix-be) application based on a system knowledge graph mapping (using `graphify`), covering modular structure, data flow, service layers, and database management.

---

## 1. Architecture Overview

The Medix Backend is built using **Go (v1.26.5)** with the **Gin** web framework for HTTP routing and **GORM** as an ORM for interacting with **PostgreSQL** database.

### Key Design Principles

Each business domain is isolated into independent modules under the `internal/` directory. Each module implements the standard pattern:

1. **Model / DTO**: Request/response data structures and database entities.
2. **Repository**: Data access layer (GORM) communicating directly with PostgreSQL.
3. **Service**: Business logic layer.
4. **Handler**: HTTP controller layer (Gin router/handlers) using standard response envelope.
5. **App Wiring (`app.go`)**: `StartApp(cfg)` function that wires dependencies (dependency injection) from repo → service → handler.

---

## 2. Main Domain Modules

Based on system connectivity analysis, the architecture is divided into several main modules:

- **Medicine Module (`internal/medicine`)**: Manages medicine data, barcode search, low stock alerts, expiration dates, and active/inactive status.
- **Drug Type Module (`internal/drug`)**: Manages drug categories/types (`type-drugs`).
- **Transaction Module (`internal/transaction`)**: The largest module handling shopping cart, payment processing, transaction history, PDF/WhatsApp receipt printing, and return flow with approval threshold.
- **User Module (`internal/user`)**: Manages authentication (JWT), profile management, roles, and password updates.
- **Report Module (`internal/report`)**: Generates daily sales reports, drug ranking, and exports reports to Excel and PDF formats.

---

## 3. Communication Flow & Entry Point

- **Main Entry Point**: `cmd/api/main.go`
- **Startup Initialization**:
  1. Load configuration from `config/config.yaml` using **Viper**.
  2. Run database migrations automatically via **golang-migrate** from `migrations/data`.
  3. Connect to PostgreSQL and ensure `uuid-ossp` extension is active.
  4. Set up CORS and register API route groups for each module (`/api/v1/...`).

---

## 4. Core Abstraction & God Nodes (Central Components)

Based on system connectivity metrics, components with the highest dependencies (God Nodes) that form the application pillars include:

1. **`Success()` (`internal/common/response`)**: Global response wrapper helper used extensively across all handlers.
2. **`TransactionRepository` / `TransactionService` / `TransactionHandler`**: The backbone of cashier operations and shopping cart.
3. **`UserService` / `UserHandler`**: Central authentication, session management, and security.
4. **`MedicineService` / `MedicineHandler`**: Medicine management, search, and inventory operations.
5. **`AuthMiddleware`**: Token validation and authentication middleware.

---

## 5. Security & Infrastructure

### Authentication & Authorization

- **Authentication**: JWT-based authentication with bcrypt password hashing
- **Authorization**: Role-based access control (RBAC) for admin, kasir, owner
- **CORS**: Configured for frontend application (`http://localhost:5173`)
- **HTTPS**: TLS encryption for production deployments

### Infrastructure

- **Containerization**: Docker support for consistent deployment
- **Configuration Management**: Viper for environment-based configuration
- **Logging**: Structured logging with logrus
- **Monitoring**: Logging and health check endpoints

---

## 6. Testing & Quality Assurance

### Testing Strategy

- **Unit Tests**: Individual function testing with mocking
- **Integration Tests**: End-to-end scenario testing
- **Test Coverage**: Comprehensive coverage of core functionality
- **CI/CD**: Automated testing in CI/CD pipelines

### Quality Assurance

- **Code Style**: Consistent formatting and naming conventions
- **Code Reviews**: Peer review process for all changes
- **Static Analysis**: Linting and code quality checks
- **Performance Testing**: Load and performance testing

---

## 7. Deployment & Scalability

### Deployment Strategies

- **Local Development**: `go run cmd/api/main.go`
- **Docker Containerization**: Multi-stage Docker builds
- **Cloud Deployment**: Kubernetes, AWS, GCP support
- **Scaling**: Horizontal scaling with load balancing

### Backup & Recovery

- **Database Backups**: Automated PostgreSQL backups
- **Disaster Recovery**: Backup and restore procedures
- **High Availability**: Multi-zone deployments
- **Monitoring**: Health checks and alerts

---

## 8. Implementation Details

Based on graphify-out analysis of the codebase, the main architectural components and their connectivity:

- **Central Components**: `Success()`, `TransactionRepository`, `TransactionService`, `TransactionHandler`, `UserService`, `UserHandler`
- **Key Modules**: `medicine`, `drug`, `user`, `transaction`, `report`
- **Integration Points**: Auth middleware, service dependencies, repository connections
- **Data Flow**: Frontend → HTTP → Handler → Service → Repository → Database

This architecture emphasizes:
1. **Separation of Concerns**: Clear boundaries between layers
2. **Dependency Injection**: Clean dependency management
3. **Testability**: Mockable components for unit testing
4. **Maintainability**: Modular structure for easy maintenance
5. **Performance**: Optimized data access patterns

---

## 9. Future Considerations

### Scalability

- **Microservices Architecture**: Potential future microservices decomposition
- **API Gateway**: Centralized API management
- **Service Mesh**: Advanced service communication
- **Event-Driven Architecture**: Asynchronous communication patterns

### Technology Updates

- **New Frameworks**: Updates to Go runtime and web frameworks
- **Database Optimization**: PostgreSQL features and performance tuning
- **Security Enhancements**: Advanced authentication mechanisms
- **Monitoring Tools**: Enhanced observability and monitoring

---

## Summary

The Medix Backend architecture is designed for scalability, maintainability, and performance. It follows modern software engineering practices with:

1. **Clean Architecture**: Clear separation of concerns
2. **Test-Driven Development**: Comprehensive testing strategy
3. **Continuous Integration**: Automated CI/CD pipelines
4. **Security-First**: Built-in security measures
5. **Performance-Optimized**: Efficient data access patterns

This architecture provides a solid foundation for building scalable, secure, and maintainable applications while maintaining code quality and developer productivity.

---

**Last Updated:** October 2024
**Version:** 1.0
**Author:** Medix Backend Team