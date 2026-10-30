# Code Style Guide

## Overview

This document defines the coding standards and best practices used in the Medix Backend (medix-be) project. Adhering to these guidelines ensures code consistency, readability, and maintainability across the codebase.

---

## 🏗️ Language & Framework

- **Language**: Go 1.26.5
- **Web Framework**: Gin (github.com/gin-gonic/gin v1.12.0)
- **ORM**: GORM (gorm.io/gorm v1.31.2)
- **Database**: PostgreSQL 14+
- **Migration**: golang-migrate/v4 (migrations/data/*.sql)
- **Authentication**: JWT (golang-jwt/v5 v5.3.1)
- **Logging**: logrus (github.com/sirupsen/logrus v1.9.4)
- **Configuration**: Viper (github.com/spf13/viper v1.21.0)
- **Security**: bcrypt for password hashing (golang.org/x/crypto/bcrypt)

---

## 📁 Project Structure

The project follows a clean architecture pattern:

```
medix-be/
├── cmd/api/
│   └── main.go              # Application entry point
├── config/
│   └── config.yaml         # Application configuration
├── internal/
│   ├── common/              # Shared utilities
│   │   ├── response/        # Standard API responses
│   │   ├── mailer.go        # Email utilities
│   │   └── email_template.go
│   ├── medicine/            # Medicine domain
│   │   ├── app.go           # Module entry point
│   │   ├── handler/         # HTTP handlers
│   │   ├── repository/      # Database repositories
│   │   ├── service/         # Business logic
│   │   └── model/           # Data models
│   ├── drug/                # Drug type domain
│   ├── user/                # User domain (auth, profiles)
│   ├── transaction/         # Transaction domain
│   └── report/              # Reporting domain
├── migrations/
│   └── data/                # Database migration scripts
└── docs/                    # Documentation
    ├── API.md
    ├── DATABASE.md
    ├── SECURITY.md
    ├── README.md
    └── ...
```

---

## 📝 File Structure & Conventions

### Package Organization
- Each domain (`medicine`, `drug`, `user`, `transaction`, `report`) is in its own package
- Handler files end with `Handler.go` (e.g., `MedicineHandler.go`)
- Repository files end with `Repo.go` (e.g., `medicineRepo.go`)
- Service files end with `Service.go` (e.g., `medicineService.go`)
- Models/DTOs are in `model/dto/` and `model/entities/` directories
- Common utilities are in `internal/common/`

### File Naming
- Use `kebab-case` for filenames (e.g., `medicine_handler.go`, `user_handler.go`)
- Use `PascalCase` for struct names and exported types
- Use `camelCase` for variable, function, and method names

### Package Naming
- Use single-word lowercase packages for handlers (e.g., `package handler`)
- Domain packages use the domain name in lowercase (e.g., `package medicine`, `package user`)

---

## 💻 Code Style

### Import Organization
Imports follow this order:
1. Standard library imports
2. Third-party imports (alphabetical order)
3. Local imports (alphabetical order)

Example:
```go
import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"medix-be/internal/common/response"
	"medix-be/internal/medicine/model/dto"
	"medix-be/internal/medicine/service"
	transDto "medix-be/internal/transaction/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)
```

### Method & Function Naming
- Use `camelCase` for methods and functions
- Use `PascalCase` for struct methods when they are part of the public interface
- Function names should be descriptive and actionable

### Variable & Constant Naming
- Use `camelCase` for variables
- Use `PascalCase` for exported constants
- Use `ALL_CAPS` for package-level constants

### Commenting Style
- Use Go's standard comment format (`//` for single-line, `/* ... */` for multi-line)
- Document exported functions, types, and methods with comments explaining their purpose
- Use `TODO` and `FIXME` comments for incomplete code or known issues
- Avoid inline comments unless necessary; prefer clear, self-documenting code

### Error Handling
- Return errors as values (Go's idiomatic way)
- Use `error` types from the standard library or third-party packages
- For HTTP handlers, use the `response` helper from `medix-be/internal/common/response`
- Never panic in production code unless it's for unrecoverable errors

### API Response Pattern
All HTTP handlers follow this pattern:

```go
func (h *MedicineHandler) CreateMedicine() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dto.CreateMedicineRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.medicineService.CreateMedicine(payload)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusCreated, "Medicine added successfully", res)
	}
}
```

---

## 📊 Code Quality & Best Practices

### Dependency Injection
- Use dependency injection for service wiring
- Handler constructors take contract interfaces as parameters
- Example: `StartMedicineHandler(contract *HandlerContract, props *MedicineHandlerProps)`

### Validation
- Validate all input data at the trust boundary (HTTP handlers)
- Use Go's `binding` tags for request validation
- Return descriptive error messages for validation failures

### Security
- Never hardcode secrets; use environment variables or secure configuration
- Hash passwords with bcrypt
- Use JWT for authentication with appropriate expiration times
- Implement proper CORS policies
- Sanitize and validate all user inputs

### Logging
- Use logrus for structured logging
- Include relevant context in log messages (request ID, user ID, etc.)
- Log errors with appropriate levels (error, warn, info)

### Testing
- Write focused unit tests for individual functions
- Use table-driven tests for testing various scenarios
- Mock dependencies when testing
- Aim for high test coverage (ideally >80%)

### Performance
- Use efficient database queries with proper indexing
- Implement pagination for list endpoints
- Use connection pooling for database connections
- Profile and optimize bottlenecks

---

## 🎨 Formatting

### Code Formatting
- Use `gofmt` for consistent formatting
- Set `tabwidth` to 4 spaces in editor configuration
- Use line length of 120 characters

### Line Length
- Keep lines shorter than 120 characters when possible
- Break long lines logically (e.g., function calls, type assertions)

### Whitespace
- Use spaces, not tabs
- Add blank lines between logical sections in functions
- Use blank lines after imports and before package declarations

---

## 🏷️ Naming Conventions

### Structs
- Use `PascalCase` (e.g., `MedicineHandler`, `UserHandler`)
- Use meaningful names that describe the purpose
- For structs representing database entities, use the database table name in singular form (e.g., `Obat` for `obat`)

### Interfaces
- Use `PascalCase` for interface names (e.g., `TransactionService`)
- Prefix interfaces with `I` when used as types (e.g., `IHandler`)

### Variables
- Use `camelCase` for variables (e.g., `userID`, `medicineService`)
- Use `camelCase` for loop variables (e.g., `i`, `j` for indices)
- Use `camelCase` for boolean flags (e.g., `isActive`, `shouldValidate`)

### Constants
- Use `ALL_CAPS` for package-level constants (e.g., `MAX_LIMIT`, `DEFAULT_PAGE`)
- Use `camelCase` for exported constants (e.g., `ErrorNotFound`)

### Functions
- Use `camelCase` for function names (e.g., `getAllMedicines`, `createMedicine`)
- Use descriptive names that indicate the action performed

### Types
- Use `PascalCase` for custom types (e.g., `UserStatus`, `MedicineFilterParams`)

---

## 🔧 Tools & Linting

### Code Formatting & Linting
Use `gofmt` for formatting and `golint` or `revive` for linting. Consider using CI/CD to enforce code quality.

### IDE Configuration
- Tab width: 4 spaces
- Line length: 120 characters
- Auto-insert final newline
- Trim trailing whitespace

### Pre-commit Hooks
- Run `gofmt -d` before committing
- Run `golint ./...` to check for linting issues
- Run tests before committing

---

## 📚 Documentation

- Document public APIs with Go doc comments
- Document complex business logic with comments
- Keep documentation in sync with code changes
- Use Markdown for project documentation (in the `docs/` directory)

---

## 🚀 Deployment & Configuration

### Environment Variables
- Use `.env` files for local development (see `config/config.yaml.example`)
- Never commit secrets or API keys to version control
- Use environment-specific configurations

### Configuration
- Use `config/config.yaml` for application configuration
- Use `Viper` for configuration management
- Support environment variables and configuration files

---

## 🚨 Common Pitfalls

1. **Not validating input data**: Always validate user input
2. **Hardcoding secrets**: Use environment variables or secure configuration management
3. **Not handling errors**: Always handle and return errors appropriately
4. **Not closing resources**: Always close database connections, files, etc.
5. **Not using context**: Use context for request-scoped values
6. **Not testing**: Always write tests for your code
7. **Not following Go conventions**: Follow Go's idiomatic patterns

---

## 📝 Version Control

- Use feature branches for new development
- Use conventional commit messages (e.g., `feat: add new endpoint`)
- Write descriptive commit messages
- Rebase feature branches before merging to main

---

## 🎯 Summary

This code style guide ensures that the Medix Backend maintains high code quality, consistency, and maintainability. By following these guidelines, developers can work more efficiently, produce cleaner code, and reduce technical debt.

---

**Last Updated:** October 2024
**Version:** 1.0
