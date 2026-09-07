# Medix BE

[![Go](https://img.shields.io/badge/Go-1.26.5-00ADD8?logo=go)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-1.12-00ADD8)](https://github.com/gin-gonic/gin)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

REST API for the **Medix** medical service system.

---

## Quick Start

```bash
# Start PostgreSQL
docker compose up -d

# Run server (migrations + seed run automatically)
go run cmd/api/main.go
```

Server runs at **http://localhost:8080**.

---

## Prerequisites

- Go 1.26+
- PostgreSQL 16+ (or Docker)

---

## Configuration

Main config: `config/config.yaml`. Override via `.env`:

```bash
cp .env.example .env
```

| Key | Default |
|-----|---------|
| `app.port` | `8080` |
| `database.name` | `medix` |
| `jwt.secret` | `change-this-secret-in-production` |
| `migration_path` | `file://migrations/data` |

---

## API Endpoints

All protected endpoints require `Authorization: Bearer <token>` header.

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/users/login` | Login, returns JWT (8h expiry) |
| `POST` | `/api/v1/users` | Create user (public) |

### User Profile

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/users/profile` | Get current user profile |
| `PUT` | `/api/v1/users/profile` | Update profile (name, phone, username) |
| `POST` | `/api/v1/users/profile/photo` | Upload profile photo (max 2MB, jpg/png/webp) |

### User Management (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/users` | List users (paginated, filterable) |
| `GET` | `/api/v1/users/:id` | Get user by ID |
| `PUT` | `/api/v1/users/:id` | Update user |
| `DELETE` | `/api/v1/users/:id` | Delete user |

### Drug Types

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/type-drugs` | List all drug types |
| `POST` | `/api/v1/type-drugs` | Create drug type |
| `GET` | `/api/v1/type-drugs/:id` | Get drug type |
| `PUT` | `/api/v1/type-drugs/:id` | Update drug type |
| `DELETE` | `/api/v1/type-drugs/:id` | Delete drug type |

### Medicines

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/medicines` | List medicines (search, filter, pagination) |
| `POST` | `/api/v1/medicines` | Create medicine |
| `GET` | `/api/v1/medicines/barcode/:barcode` | Find by barcode |
| `GET` | `/api/v1/medicines/:id` | Get medicine by ID |
| `PUT` | `/api/v1/medicines/:id` | Update medicine |
| `PATCH` | `/api/v1/medicines/:id/status` | Toggle active status |
| `GET` | `/api/v1/medicines/alerts/low-stock` | Low stock alerts |
| `GET` | `/api/v1/medicines/alerts/expiring?days=30` | Expiring soon alerts |
| `GET` | `/api/v1/medicines/alerts/summary` | Notification summary |

### Transactions

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/transactions` | Create transaction (multi-item, auto stock deduction) |
| `GET` | `/api/v1/transactions` | List all transactions |
| `GET` | `/api/v1/transactions/today` | Current user's transactions today + summary |
| `GET` | `/api/v1/transactions/:id` | Transaction detail |
| `PATCH` | `/api/v1/transactions/:id/cancel` | Cancel (same day only, restores stock) |
| `GET` | `/api/v1/transactions/:id/receipt` | Digital receipt |

### Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/reports/sales-summary?group_by=daily\|weekly\|monthly` | Sales chart data |
| `GET` | `/api/v1/reports/drug-ranking` | Top 10 / Bottom 10 selling medicines |
| `GET` | `/api/v1/reports/export/excel` | Export to Excel |
| `GET` | `/api/v1/reports/export/pdf` | Export to PDF |

---

## Test Accounts (Seeded)

| Username | Password | Role | Status |
|----------|----------|------|--------|
| `admin` | `admin123` | admin | Active |
| `kasir1` | `kasir123` | cashier | Active |
| `kasir2` | `kasir123` | cashier | Inactive |
| `kasir3` | `kasir123` | cashier | Active |
| `owner` | `owner123` | owner | Active |

---

## Project Structure

```
medix-be/
├── cmd/api/           # Entry point
├── config/            # Configuration (YAML + env)
├── internal/
│   ├── common/        # Response helpers
│   ├── drug/          # Drug type module
│   ├── medicine/      # Medicine module
│   ├── middleware/    # Auth & role middleware
│   ├── report/        # Reports & dashboard
│   ├── transaction/   # Transaction module
│   └── user/          # Auth & profile
├── migrations/
│   └── data/          # SQL migrations + seed
├── docs/              # Architecture, ER, OpenAPI
├── Medix/             # Postman collections
├── .env.example
├── config.yaml
├── docker-compose.yml
└── README.md
```

---

## Build

```bash
go build ./...
```

---

## Notes

- Migrations are **forward-only** (no `.down.sql`)
- User status stored as `SMALLINT` (1 = active, 0 = inactive), exposed as strings `"active"` / `"inactive"`
- CORS hardcoded to `http://localhost:5173` (Vite dev server)
- Price snapshots stored in transaction details for historical accuracy