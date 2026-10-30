# Environment Variables Documentation

## Overview

This document documents the environment variables used in the Medix Backend (medix-be) project. Proper environment variable configuration is essential for secure deployment and operation of the application.

---

## Required Environment Variables

The following environment variables must be set for the application to function correctly:

### JWT Configuration

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `JWT_SECRET` | Not set | Secret key for signing JWT tokens. **Required** for authentication. | `your-secret-key-1234567890` |

> **Security Note:** This should be a strong, random string (32+ characters) stored securely in production environments.

### Database Configuration

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `DATABASE_HOST` | Not set | Host server database PostgreSQL. | `localhost` |
| `DATABASE_PORT` | Not set | Port server database PostgreSQL. | `5432` |
| `DATABASE_USER` | Not set | Username database PostgreSQL. | `medix` |
| `DATABASE_PASSWORD` | Not set | Password database PostgreSQL. | `strong-password-123` |
| `DATABASE_NAME` | Not set | Nama database PostgreSQL. | `medix` |
| `DATABASE_SSL` | Not set | Mode SSL database PostgreSQL. Values: `disable`, `require`, `verify-full`. | `disable` |

> **Note:** The `DATABASE_URL` environment variable can be used as an alternative to individual database variables.

### Mail/SMTP Configuration

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `MAIL_HOST` | Not set | Host server SMTP. | `smtp.gmail.com` |
| `MAIL_PORT` | Not set | Port server SMTP. | `587` |
| `MAIL_USERNAME` | Not set | Username autentikasi SMTP. | `your-email@gmail.com` |
| `MAIL_PASSWORD` | Not set | Password autentikasi SMTP. | `your-smtp-password` |
| `MAIL_FROM` | Not set | Alamat pengirim email. | `no-reply@medix.local` |

> **Note:** For Gmail, you may need to use an App Password or enable 2-Step Authentication.

### Application Configuration

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `APP_NAME` | Not set | Nama aplikasi untuk logging dan branding. | `Medix BE` |
| `APP_PORT` | Not set | Port yang didengarkan aplikasi. | `8080` |
| `FRONTEND_URL` | Not set | URL aplikasi frontend untuk CORS. | `http://localhost:5173` |

---

## Optional Environment Variables

The following environment variables are optional and provide additional configuration:

### JWT Configuration (Optional)

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `JWT_EXPIRY_HOURS` | Not set | JWT token expiry in hours. | `24` |
| `JWT_REFRESH_EXPIRY_HOURS` | Not set | JWT refresh token expiry in hours. | `168` |

### Logging Configuration

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error`. | `info` |
| `LOG_FORMAT` | `json` | Log format: `json` or `text`. | `json` |

### Rate Limiting (if implemented)

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `RATE_LIMIT_REQUESTS` | Not set | Maximum requests per minute per IP. | `100` |
| `RATE_LIMIT_WINDOW` | Not set | Rate limiting time window in seconds. | `60` |

### Security Headers

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `SECURITY_CORS_ORIGINS` | `*` | CORS allowed origins. | `http://localhost:5173` |

---

## Environment File Examples

### Development Environment (`.env.development`)

```env
# JWT Configuration
JWT_SECRET=dev-secret-key-1234567890

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=medix
DATABASE_PASSWORD=medix123
DATABASE_NAME=medix
DATABASE_SSL=disable

# Mail
MAIL_HOST=sandbox.smtp.mailtrap.io
MAIL_PORT=2525
MAIL_USERNAME=41f48730974e0c
MAIL_PASSWORD=96c0802df0f6f1
MAIL_FROM=no-reply@medix.local

# Aplikasi
APP_NAME=Medix BE
APP_PORT=8080
FRONTEND_URL=http://localhost:5173

# Logging
LOG_LEVEL=debug
```

### Production Environment (`.env.production`)

```env
# JWT Configuration
JWT_SECRET=your-super-secret-key-32-chars-min
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168

# Database
DATABASE_HOST=medix-db.internal
DATABASE_PORT=5432
DATABASE_USER=medix
DATABASE_PASSWORD=${MEDIX_DB_PASSWORD}
DATABASE_NAME=medix
DATABASE_SSL=require

# Mail
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=${GMAIL_APP_PASSWORD}
MAIL_FROM=notifications@medix.com

# Aplikasi
APP_NAME=Medix BE
APP_PORT=8080
FRONTEND_URL=https://medix.com

# Keamanan
SECURITY_CORS_ORIGINS=https://medix.com

# Logging
LOG_LEVEL=info
```

---

## Environment Variable Usage in Code

The application reads environment variables using:

### Go Code Example (`config/config.go`)

```go
import (
    "os"
    "github.com/spf13/viper"
)

func init() {
    // Load environment variables
    viper.AutomaticEnv()
    
    // Set defaults
    viper.SetDefault("jwt.secret", "change-this-secret-in-production")
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", "5432")
    
    // Application configuration
    appConfig := struct {
        Name string `mapstructure:"name"`
        Port string `mapstructure:"port"`
        FrontendURL string `mapstructure:"frontend_url"`
    }{}
    
    if err := viper.Unmarshal(&appConfig); err != nil {
        panic(err)
    }
}
```

### Shell Check Example (`#!/bin/sh`)

```sh
#!/bin.sh

// Check required environment variables
required_vars=("JWT_SECRET" "DATABASE_PASSWORD" "MAIL_PASSWORD")

for var in "${required_vars[@]}"; do
    if [[ -z "${!var}" ]]; then
        echo "Error: $var is not set"
        exit 1
    fi
    echo "$var is set"
done

// Check database connection
if ! psql "host=localhost dbname=medix user=medix password=${DATABASE_PASSWORD}" -c "SELECT 1" >/dev/null 2>&1; then
    echo "Error: Cannot connect to database"
    exit 1
fi

echo "All required environment variables are set and validated"
```

### Go Validation Example (`internal/config/validate.go`)

```go
package config

import (
    "os"
    "fmt"
)

func ValidateEnvironment() error {
    requiredVars := []string{"JWT_SECRET", "DATABASE_PASSWORD", "MAIL_PASSWORD"}
    
    for _, varName var range requiredVars {
        if value := os.Getenv(varName); value == "" {
            return fmt.Errorf("required environment variable %s is not set", varName)
        }
    }
    
    // Validate database connection
    if err := validateDatabaseConnection(); err != nil {
        return fmt.Errorf("database connection validation failed: %w", err)
    }
    
    return nil
}
```

---

## Environment Variable Documentation Template

### Creating a New Environment Variable

When adding a new environment variable:

1. **Update `config/config.yaml`** with default values
2. **Add documentation to `docs/ENV.md`**
3. **Add validation logic** if needed
4. **Update `.env.example`** file
5. **Update documentation** in project documentation

Example:

```yaml
// config/config.yaml
new_feature:
  enabled: false
  api_key: ""
  timeout: "30s"
```

```env
// .env.example
NEW_FEATURE_ENABLED=true
NEW_FEATURE_API_KEY=your-api-key-here
NEW_FEATURE_TIMEOUT=30s
```

---

## Environment Variable Usage Summary

| Environment Variable | Category | Required | Default | Description |
|----------------------|----------|----------|---------|-------------|
| `JWT_SECRET` | JWT | Yes | Not set | Secret key for JWT signing |
| `DATABASE_HOST` | Database | Yes | Not set | Database host |
| `DATABASE_PORT` | NAFTA | Yes | Not set | Database port |
| `DATABASE_NAME` | Database | Yes | Not set | Database name |
| `DATABASE_PASSWORD` | Database | Yes | Not set | Database password |
| `DATABASE_SSL` | Database | No | `disable` | SSL mode |
| `MAIL_HOST` | Mail | Yes | Not set | SMTP host |
| `MAIL_PORT` | Mail | Yes | Not set | SMTP port |
| `MAIL_USERNAME` | Mail | Yes | Not set | SMTP username |
| `MAIL_PASSWORD` | Mail | Yes | Not set | SMTP password |
| `MAIL_FROM` | Mail | Yes | Not set | Sender email address |
| `APP_NAME` | App | No | Not set | Application name |
| `APP_PORT` | App | No | Not set | Application port |

---

## Additional Resources

- [Viper Documentation](https://github.com/spf13/viper)
- [JWT Go Documentation](https://github.com/golang-jwt/jwt)
- [PostgreSQL Environment Variables](https://www.postgresql.org/docs/current/libpq-envars.html)
- [Security Best Practices](https://owasp.org/www-project-top-ten/)
- [Docker Environment Management](https://docs.docker.com/compose/environment-variables/)

---

**Last Updated:** October 2024
**Version:** 1.0
**Author:** Medix Backend Team