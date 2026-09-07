# AGENTS.md — medix-be

Compact guide for OpenCode sessions. If a fact is obvious from filenames, it is not here.

## Stack & entrypoint

- Go 1.26.5, Gin, GORM, PostgreSQL, golang-migrate.
- Entrypoint: `cmd/api/main.go`.
- Run local server: `go run cmd/api/main.go` (defaults to `:8080`).

## Configuration

- Primary config is `config/config.yaml` (loaded by Viper), not `.env`.
- Required keys: `database.*`, `jwt.secret`, `migration_path`.
- `JWT_SECRET` is also read from `os.Getenv` in the auth middleware; `config.LoadEnv()` copies `jwt.secret` to the env var if the env var is empty.
- CORS is hardcoded to `http://localhost:5173` (Vite dev server) in `main.go`.

## Database & migrations

- Migrations run automatically on startup via `migrations.RunMigrations()`.
- Migration source: `migrations/data` (golang-migrate file source).
- There are no `.down.sql` files; schema changes are forward-only.
- `config.ConnectDatabase()` creates the `uuid-ossp` extension on first connect.

## Response envelope convention

- Use the helper in `internal/common/response` for **all** handler responses:

  ```go
  response.Success(c, http.StatusOK, "Pesan sukses", data)
  response.Error(c, http.StatusBadRequest, "Pesan error")
  ```

- Shape is always `{status, message, data}`. For errors, `data` is `nil`.
- **Exception:** `internal/middleware/AuthMiddleware.go` still returns raw `gin.H{"error": "..."}` — leave it as-is or update globally; do not mix raw and helper responses in the same handler.
- `MedicineHandler.go` already returns `http.StatusConflict` for duplicate-barcode errors; that special-case is fine, but keep using the helper for the body.

## Module wiring

- Each domain (`medicine`, `drug`, `user`, `transaction`, `report`) exposes `StartApp(cfg)` from `app.go`.
- Inside `app.go`: instantiate repo → service → handler with `HandlerContract{Logger, Router}`.
- Handler structs are created by `StartXxxHandler(contract, props)` and register their own route groups (e.g., `/medicines`, `/type-drugs`).
- Gotcha: module config struct names are inconsistent. `medicine` uses `ModuleConfig`; `drug`, `transaction`, `report` reuse their handler struct name (`TypeDrugHandler`, `TransactionHandler`, `ReportHandler`); `user` uses `UserHandler` and also accepts `AuthRouter` for split public/protected routes.

## Auth & routes

- `middleware.Auth()` reads `Authorization: Bearer <jwt>` and sets:
  - `c.Get("user_id")` as `uint`
  - `c.Get("username")`
  - `c.Get("role")` as `string`
- Use `middleware.RequireRoles("admin", "kasir", ...)` for role-gated routes if needed.
- User routes are split:
  - Public (`/api/v1/users`): `POST /login`, `POST /` (create user).
  - Protected (`/api/v1/users`): everything else (`GET`, `GET /:id`, `PUT /:id`, `DELETE /:id`, `/profile`).

## Status-field gotcha

- User `status` is stored in the DB as `SMALLINT` (`1` = aktif, `0` = nonaktif) but exposed as the strings `"aktif"` / `"nonaktif"` in `UserResponse`.
- `CreateUserRequest` and `UpdateUserRequest` use a custom `UserStatus` type that accepts **either** `"aktif"`/`"nonaktif"` strings **or** numeric `1`/`0`, so both shapes now unmarshal successfully.
- In contrast, `UpdateMedicineRequest.status` is typed `*int` with `oneof=0 1`. Status semantics differ between modules — verify the DTO type before changing request contracts.

## Testing

- One focused test exists: `internal/user/model/dto/UserStatus_test.go` verifies that `UserStatus` unmarshals from both numeric and string JSON values.
- Run it: `go test ./internal/user/model/dto/...`
- Other modules rely on manual/API smoke tests (see `Medix/` collection files).

## Commands

```powershell
# Verify compile / typecheck
go build ./...

# Build binary
go build -o medix-be.exe ./cmd/api

# Run PostgreSQL via Docker (matches default config/config.yaml)
docker compose up -d

# Run server (requires config.yaml and running PostgreSQL)
go run cmd/api/main.go

# No lint or test commands configured — add only if you also add tooling config.
```
