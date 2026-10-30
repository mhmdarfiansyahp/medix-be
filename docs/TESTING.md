# Testing Guide — medix-be

Testing guidelines for the medix-be project. Uses Go standard library `testing` package; no external frameworks.

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/user/model/dto/...

# Verbose output
go test -v ./internal/user/model/dto/...

# With coverage
go test -cover ./internal/user/model/dto/...
```

---

## Current Test Structure

```
internal/
└── user/
    └── model/
        └── dto/
            └── UserStatus_test.go    # only unit test in the project
```

### `UserStatus_test.go`

Verifies the custom `UserStatus` type unmarshals correctly from JSON:

| Input JSON | Expected Value |
|------------|----------------|
| `{"status":1}` | `"aktif"` |
| `{"status":"aktif"}` | `"aktif"` |
| `{"status":0}` | `"nonaktif"` |
| `{"status":"nonaktif"}` | `"nonaktif"` |

Ensures `CreateUserRequest` / `UpdateUserRequest` DTOs accept both numeric and string status values without error.

---

## Test Conventions

1. **File naming**: `<name>_test.go` in the same package as the code under test.
2. **Function naming**: `Test<Name>` (e.g. `TestUserStatusUnmarshal`).
3. **Table-driven tests** for multiple cases (see `UserStatus_test.go`).
4. **Response helpers** in `internal/common/response` are thin wrappers — no tests needed.
5. **API smoke tests** run manually via Postman collection in `Medix/`.

---

## Testing Strategy

### Unit Tests (High Priority)
- **DTO validation**: custom types, binding tags, enum `oneof`.
- **Service logic**: business rules (discount calculation, stock check, unique barcode validation).
- **Repository queries**: filters, pagination, search — use `sqlmock` or test DB.

### Integration Tests (Medium Priority)
- Full handler → service → repo flow with a test database (PostgreSQL container).
- Auth middleware: valid token, expired, invalid, role check.
- Migration + seed: verify schema and seed data load correctly.

### E2E / API Smoke Tests (Low Priority)
- Postman collection in `Medix/` folder for manual regression.
- Cover happy paths: login → CRUD medicines → transaction → reports.

---

## Test Database

Use Docker for an isolated test DB:

```yaml
# docker-compose.test.yml
services:
  postgres_test:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: medix_test
      POSTGRES_USER: medix
      POSTGRES_PASSWORD: medix
    ports: ["5433:5432"]
    tmpfs: /var/lib/postgresql/data
```

Run migrations before tests:

```bash
docker compose -f docker-compose.test.yml up -d
MIGRATION_PATH=file://migrations/data go run cmd/api/main.go
go test ./...
docker compose -f docker-compose.test.yml down
```

---

## CI/CD Integration (Example)

Add to `.github/workflows/test.yml`:

```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: medix_test
          POSTGRES_USER: medix
          POSTGRES_PASSWORD: medix
        ports: ["5432:5432"]
        options: >-
          --health-cmd="pg_isready -U medix"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.26' }
      - name: Run migrations
        run: |
          export DATABASE_URL="postgres://medix:medix@localhost:5432/medix_test?sslmode=disable"
          go run cmd/api/main.go & sleep 3 && kill $!
      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v4
```

---

## Checklist When Adding a New Feature

- [ ] Unit test for new DTO (validation tags, custom unmarshal).
- [ ] Unit test for new service logic (edge cases: empty, nil, boundary).
- [ ] Integration test for new handler (happy path + at least 1 error case).
- [ ] Update Postman collection `Medix/` for manual regression.
- [ ] `go test ./...` passes locally before pushing.

---

## Optional Tooling

| Tool | Purpose | Install |
|------|---------|---------|
| `testify` | Assertions (`assert`, `require`, `mock`) | `go get github.com/stretchr/testify` |
| `sqlmock` | Mock `database/sql` for repo tests | `go get github.com/DATA-DOG/go-sqlmock` |
| `httptest` | Stdlib HTTP test server (already available) | — |
| `go-cmp` | Deep equal for complex structs | `go get github.com/google/go-cmp/cmp` |

**Note**: The project currently uses **no external test dependencies**. Add only if needed (e.g. complex DB mocks).

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| `go test` hangs / timeout | Ensure DB is running and migrations are applied. Use `-timeout 30s`. |
| Flaky tests | Run `go test -race` and fix shared state or global variables. |
| Low coverage | Add tests at the service layer; handlers are thin by design. |
| `sql: database is closed` | Do not close `sql.DB` mid-test. Use `setup/teardown` per package. |

---

## References

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Testify Docs](https://github.com/stretchr/testify)
- [sqlmock Docs](https://github.com/DATA-DOG/go-sqlmock)
