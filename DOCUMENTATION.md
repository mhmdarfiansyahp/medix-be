# Medix BE - Documentation

## 1. Background
Problem: manage pharmacy inventory, sales, users, reports. Prior: manual spreadsheets, no centralized API. Urgency: need real‑time stock, sales tracking, role‑based access. Trend: microservice APIs, Go Gin, PostgreSQL.

## 2. Purpose of Creation
Primary goal: provide REST API for medicines, drug types, users, transactions, reports. Secondary goals: JWT auth, role gating, data validation, export reports. Targets: <ul><li>CRUD ops for each domain</li><li>JWT‑protected routes</li><li>PDF/Excel export</li></ul>Scope: API only, no UI, no third‑party integrations.

## 3. Benefits
### 3.1 Functional Benefits
- CRUD endpoints simplify pharmacy staff tasks.
- Alerts for low stock and expiring medicines.
- Transaction history & receipt generation.
- Role‑based access (admin, kasir, owner).
- Export sales summary & drug ranking.
### 3.2 Non‑Functional Benefits
- **Performance**: Go + Gin, low latency.
- **Scalability**: Horizontal scaling via Docker, stateless handlers.
- **Reliability**: PostgreSQL with `uuid‑ossp` extension, migrations run on startup.
- **Security**: JWT auth, CORS limited to frontend URL, role middleware.
- **Maintainability**: Layered structure (handler → service → repo).
- **Usability**: Consistent JSON envelope via `response` helper.
- **Compatibility**: Works on any platform supporting Go 1.26+, PostgreSQL.
- **Portability**: Docker compose for dev environment.

## 4. Conclusion
Background shows need for centralized pharmacy API. Purpose delivers CRUD, auth, reporting. Value: streamlined operations, data integrity, extensible foundation. Future: add pagination, analytics, multi‑tenant support. Recommend deploying via Docker, use JWT from login endpoint.

## 5. Registered APIs
| No | Method | Endpoint | Description | Auth | Parameters | Response |
|---|--------|----------|-------------|------|------------|----------|
| 1 | POST | /api/v1/users/login | Authenticate user, return JWT | No | `username`, `password` (JSON) | 200 `{token, user}`; 401 error |
| 2 | POST | /api/v1/users | Create new user | Yes (admin) | `CreateUserRequest` JSON | 201 user created |
| 3 | GET | /api/v1/users | List users (paginated) | Yes (admin/kasir/owner) | `page`, `limit`, `search`, `role`, `status` (query) | 200 `UserListResponse` |
| 4 | GET | /api/v1/users/:id | Get user by ID | Yes | path ID | 200 `UserResponse` |
| 5 | PUT | /api/v1/users/:id | Update user | Yes (admin) | `UpdateUserRequest` JSON | 200 `UserResponse` |
| 6 | DELETE | /api/v1/users/:id | Delete user | Yes (admin) | path ID | 200 success |
| 7 | GET | /api/v1/users/profile | Get own profile | Yes | – | 200 `UserResponse` |
| 8 | PUT | /api/v1/users/profile | Update own profile | Yes | `UpdateProfileRequest` multipart/form-data | 200 `UserResponse` |
| 9 | POST | /api/v1/users/profile/photo | Upload profile photo | Yes | file multipart | 200 success |
| 10 | GET | /api/v1/medicines | List medicines (filters) | Yes | `page`, `limit`, `search`, `jenis_obat_id`, `status`, `status_stok`, `barcode` | 200 `[]MedicineResponse` |
| 11 | GET | /api/v1/medicines/:id | Get medicine by ID | Yes | path ID | 200 `MedicineResponse` |
| 12 | GET | /api/v1/medicines/barcode/:barcode | Get medicine by barcode | Yes | path barcode | 200 `MedicineResponse` |
| 13 | POST | /api/v1/medicines | Create medicine | Yes (admin/kasir) | `CreateMedicineRequest` JSON | 201 medicine created |
| 14 | PUT | /api/v1/medicines/:id | Update medicine | Yes (admin/kasir) | `UpdateMedicineRequest` JSON | 200 updated |
| 15 | PATCH | /api/v1/medicines/:id/status | Toggle active status | Yes (admin) | `{is_active:bool}` JSON | 200 status updated |
| 16 | DELETE | /api/v1/medicines/:id | Delete medicine | Yes (admin) | path ID | 200 success |
| 17 | GET | /api/v1/medicines/alerts/low-stock | Low‑stock alerts | Yes | – | 200 list |
| 18 | GET | /api/v1/medicines/alerts/expiring | Expiring alerts | Yes | – | 200 list |
| 19 | GET | /api/v1/medicines/alerts/summary | Alerts summary | Yes | – | 200 summary |
| 20 | GET | /api/v1/type-drugs | List drug types | Yes | – | 200 `[]TypeDrugResponse` |
| 21 | GET | /api/v1/type-drugs/:id | Get drug type by ID | Yes | path ID | 200 `TypeDrugResponse` |
| 22 | POST | /api/v1/type-drugs | Create drug type | Yes (admin) | `CreateTypeDrugRequest` JSON | 201 created |
| 23 | PUT | /api/v1/type-drugs/:id | Update drug type | Yes (admin) | `UpdateTypeDrugRequest` JSON | 200 updated |
| 24 | DELETE | /api/v1/type-drugs/:id | Delete drug type | Yes (admin) | path ID | 200 success |
| 25 | POST | /api/v1/transactions | Create transaction | Yes (kasir) | `CreateTransactionRequest` JSON | 201 transaction saved |
| 26 | GET | /api/v1/transactions | List all transactions | Yes (admin/kasir) | – | 200 list |
| 27 | GET | /api/v1/transactions/today | Today's transactions for logged user | Yes (kasir) | – | 200 list |
| 28 | GET | /api/v1/transactions/:id | Get transaction by ID | Yes | path ID | 200 transaction |
| 29 | PATCH | /api/v1/transactions/:id/cancel | Cancel transaction | Yes (kasir) | – | 200 cancelled |
| 30 | GET | /api/v1/transactions/:id/receipt | Download receipt PDF | Yes | path ID | 200 PDF stream |
| 31 | GET | /api/v1/reports/sales-summary | Sales summary (date range, group) | Yes (admin) | `start_date`, `end_date`, `group_by` query | 200 summary |
| 32 | GET | /api/v1/reports/drug-ranking | Drug ranking (date range) | Yes (admin) | `start_date`, `end_date` query | 200 ranking |
| 33 | GET | /api/v1/reports/export/excel | Export sales to Excel | Yes (admin) | same as sales‑summary | 200 Excel file |
| 34 | GET | /api/v1/reports/export/pdf | Export sales to PDF | Yes (admin) | same as sales‑summary | 200 PDF file |

*`[TODO: needs confirmation]` placeholders removed where info available.*

## 6. Project Metadata
- **Project Name**: Medix BE
- **Version**: 1.0.0
- **Last Updated**: 2026-09-18
- **Base URL**: `http://localhost:8080/api/v1`
- **Tech Stack**: Go 1.26, Gin, GORM, PostgreSQL, Docker, golang-migrate
- **Repository**: https://github.com/yourorg/medix-be

## 7. Getting Started
### 7.1 Prerequisites
- Go >= 1.26
- PostgreSQL >= 13
- Docker & Docker Compose
- Git

### 7.2 Installation Steps
1. Clone repo: `git clone https://github.com/yourorg/medix-be.git`
2. `cd medix-be`
3. Copy config template: `cp config/config.yaml.example config/config.yaml`
4. Edit `config/config.yaml` (DB credentials, JWT secret, migration_path)
5. Run migrations: `docker compose up -d` (starts PostgreSQL) then `go run cmd/api/main.go` (auto‑runs migrations) or `go run ./cmd/api` after DB ready

### 7.3 Environment Variables
| Variable | Description | Example | Required |
|----------|-------------|---------|----------|
| `JWT_SECRET` | Secret for signing JWTs | `supersecretkey` | Yes |
| `DB_HOST` | PostgreSQL host | `localhost` | Yes |
| `DB_PORT` | PostgreSQL port | `5432` | Yes |
| `DB_USER` | DB username | `medix` | Yes |
| `DB_PASSWORD` | DB password | `password` | Yes |
| `DB_NAME` | Database name | `medixdb` | Yes |
| `MIGRATION_PATH` | Path to migration files | `migrations/data` | Yes |

### 7.4 Running the Application
- **Development**: `go run cmd/api/main.go` (uses config.yaml, hot‑reload not built‑in)
- **Production (Docker)**: `docker build -t medix-be . && docker run -p 8080:8080 medix-be`
- **Testing**: `go test ./...`

## 8. Authentication & Authorization
### 8.1 Authentication Flow
1. Client POST `/users/login` with credentials.
2. Server validates, returns JWT in response body.
3. Client stores token, includes `Authorization: Bearer <jwt>` on subsequent requests.
4. Middleware `Auth()` verifies token, extracts `user_id`, `username`, `role` into Gin context.

### 8.2 Role‑Based Access Control
| Endpoint | admin | kasir | owner |
|----------|-------|-------|-------|
| `/users/*` | ✅ | ✅ (except create/delete) | ✅ |
| `/medicines/*` | ✅ | ✅ (create/update) | ❌ |
| `/type-drugs/*` | ✅ | ❌ | ❌ |
| `/transactions/*` | ✅ | ✅ | ❌ |
| `/reports/*` | ✅ | ❌ | ❌ |

### 8.3 Token Structure
```json
{
  "sub": "123",            // user ID
  "username": "jdoe",
  "role": "admin",
  "exp": 1735689600,      // expiration epoch
  "iat": 1735603200
}
```

## 9. API Conventions
### 9.1 Standard Response Format
**Success**:
```json
{
  "status": "success",
  "message": "<human readable>",
  "data": <payload>
}
```
**Error**:
```json
{
  "status": "error",
  "message": "<error description>",
  "data": null
}
```
### 9.2 HTTP Status Codes
| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET/PUT/PATCH/DELETE |
| 201 | Created | Resource created |
| 400 | Bad Request | Validation error |
| 401 | Unauthorized | Missing/invalid token |
| 403 | Forbidden | Role not permitted |
| 404 | Not Found | Resource missing |
| 409 | Conflict | Duplicate unique field |
| 500 | Internal Server Error | Unexpected failure |

### 9.3 Pagination
- Query params: `page` (default 1), `limit` (default 10)
- Response includes `pagination` object with `current_page`, `total_pages`, `total_items`, `items_per_page`.

### 9.4 Filtering & Sorting
- Filtering via query parameters (`search`, `status`, `start_date`, `end_date`, etc.)
- Sorting not currently implemented; future extensions may add `sort_by` & `order`.

## 10. API Examples
### Login (POST /api/v1/users/login)
**Request:**
```http
POST /api/v1/users/login HTTP/1.1
Content-Type: application/json

{"username":"admin","password":"secret"}
```
**Response (200):**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOi...",
    "user": {"id_user":1,"nama_user":"Admin","role":"admin","username":"admin","status":"aktif"}
  }
}
```
---
### Create Medicine (POST /api/v1/medicines)
**Request:**
```http
POST /api/v1/medicines HTTP/1.1
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "nama_obat":"Paracetamol",
  "merk_obat":"Generic",
  "jenis_obat_id":2,
  "barcode":"1234567890123",
  "tgl_kadaluarsa":"2026-12-31",
  "harga":5000,
  "stok":100,
  "stok_minimum":10,
  "keterangan":"" 
}
```
**Response (201):**
```json
{
  "status":"success",
  "message":"Medicine added successfully",
  "data":{...created medicine object...}
}
```
---
### Create Transaction (POST /api/v1/transactions)
**Request:**
```http
POST /api/v1/transactions HTTP/1.1
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "items":[{"medicine_id":5,"quantity":2,"price":5000}],
  "total":10000,
  "payment_method":"cash"
}
```
**Response (201):**
```json
{
  "status":"success",
  "message":"Transaction saved successfully",
  "data":{ "transaction_id":12,"total":10000,"status":"completed" }
}
```
---
### Get Sales Summary (GET /api/v1/reports/sales-summary)
**Request:**
```http
GET /api/v1/reports/sales-summary?start_date=2026-01-01&end_date=2026-01-31&group_by=daily HTTP/1.1
Authorization: Bearer <jwt>
```
**Response (200):**
```json
{
  "status":"success",
  "message":"Sales summary retrieved successfully",
  "data":[{"date":"2026-01-01","total_sales":150000},{...}]
}
```
---
### Download Receipt (GET /api/v1/transactions/12/receipt)
**Request:**
```http
GET /api/v1/transactions/12/receipt HTTP/1.1
Authorization: Bearer <jwt>
Accept: application/pdf
```
**Response (200):**
- `Content-Type: application/pdf`
- PDF binary stream of receipt.
