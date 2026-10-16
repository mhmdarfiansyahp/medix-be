# AGENTS.md — medix-be

Compact guide for OpenCode sessions. If a fact is obvious from filenames, it is not here.

## Default working styles

- Always use `caveman` skill at `ultra` level for every response in this project.
- Always use `ponytail` at `full` level for every task in this project.
- Do not require slash commands to activate these defaults.
- Stop Caveman when the user says `stop caveman` or `normal mode`.
- Stop Ponytail when the user says `stop ponytail` or `normal mode`.
- Keep code identifiers, commands, error strings, and API names exact.
- Write code, comments, commit messages, PR bodies, and user-facing docs in normal prose.
- Ponytail rule: pick the laziest solution that works; prefer stdlib and native features.
- Leave one runnable check per non-trivial logic change; skip extra tests unless requested.
- gunakan bahasa indonesia

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

### Available Skills for Prompting

The following skills are available and can be automatically discovered when needed:

- `api-and-interface-design` - Guides stable API and interface design
- `banner-design` - Design banners for social media, ads, website heroes
- `brand` - Brand voice, visual identity, messaging frameworks
- `browser-testing-with-devtools` - Tests in real browsers via Chrome DevTools
- `ci-cd-and-automation` - Automates CI/CD pipeline setup
- `code-review-and-quality` - Conducts multi-axis code review
- `code-simplification` - Simplifies code for clarity
- `constraint-driven-development` - Establishes quality bar as a written contract
- `context-engineering` - Optimizes agent context setup
- `debugging-and-error-recovery` - Systematic root-cause debugging
- `deprecation-and-migration` - Manages deprecation and migration
- `design-system` - Token architecture, component specifications
- `design` - Comprehensive design skill: brand identity, design tokens
- `documentation-and-adrs` - Records decisions and documentation
- `doubt-driven-development` - Subjects every non-trivial decision to adversarial review
- `frontend-ui-engineering` - Builds production-quality, accessible user interfaces
- `git-workflow-and-versioning` - Structures git workflow practices
- `idea-refine` - Refines raw ideas into sharp, actionable concepts
- `incremental-implementation` - Delivers changes incrementally in thin, verifiable slices
- `interview-me` - Extracts what the user actually wants
- `observability-and-instrumentation` - Instruments code for production visibility
- `performance-optimization` - Optimizes application performance
- `planning-and-task-breakdown` - Breaks work into ordered tasks
- `ponytail-audit` - Whole-repo audit for over-engineering
- `ponytail-debt` - Harvests ponytail deferrals into a debt ledger
- `ponytail-gain` - Show ponytail's measured impact as a compact scoreboard
- `ponytail-help` - Quick-reference card for all ponytail modes, skills, and commands
- `ponytail-review` - Code review focused exclusively on over-engineering
- `security-and-hardening` - Hardens code against vulnerabilities
- `shipping-and-launch` - Prepares production launches
- `slides` - Create strategic HTML presentations with Chart.js
- `source-driven-development` - Grounds implementation decisions in official documentation
- `spec-driven-development` - Creates specs before coding
- `test-driven-development` - Drives development with tests using the red-green-refactor loop
- `ui-styling` - Create beautiful, accessible user interfaces with shadcn/ui

### Skill Usage Guide

To use a specific skill when you want to work on a particular type of task:

1. **Identify your development phase**: Use `using-agent-skills` to discover which skill is appropriate for your task
2. **Specify the skill in your prompt**: Start your prompt with `Skill: <skill-name>`
3. **The skill will be automatically loaded and applied**: Based on your task requirements

For example:
- `Skill: frontend-ui-engineering Build a responsive dashboard component`
- `Skill: test-driven-development Write tests for the user authentication flow`
- `Skill: api-and-interface-design Design a RESTful API for the orders endpoint`