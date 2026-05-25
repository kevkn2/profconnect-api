# Backend Architecture Guidelines (Go + Fiber)

## Tech Stack

* Language: **Go (Golang)**
* Framework: **Fiber**
* Architecture Style: **Clean / Hexagonal / Usecase-Oriented**
* Database Migrations: **golang-migrate**
* Development: **Air** (live reload)

---

## Project Structure

Follow this directory structure strictly:

```text
.
├── cmd
│   ├── rest
│   ├── worker        # Optional background jobs
│   └── seeder        # Database seeding (optional)
│
├── internal
│   ├── adapter       # Presentation layer (HTTP)
│   │   ├── handler   # Request handling & usecase wiring
│   │   └── router    # Route definitions
│   │
│   ├── config
│   │   └── config.go # Environment & config loader
│   │
│   ├── database
│   │   ├── migrations     # SQL migration files (golang-migrate)
│   │   └── sqlc
│   │       ├── queries    # SQL query definitions (.sql files)
│   │       └── generated  # Auto-generated Go code (sqlc)
│   │
│   ├── domain
│   │   ├── entities  # Domain models
│   │   ├── constants # Business-specific constants
│   │   └── port      # Interfaces (ports)
│   │
│   ├── infrastructure
│   │   ├── adapter   # Port implementations
│   │   └── pkg
│   │       ├── shared # Global constants
│   │       └── utils  # Utility functions
│   │
│   └── usecase       # Application business logic
```

---

## Layer Responsibilities

---

## 1. `cmd/` — Application Entrypoints

Contains executable entry points.

### `cmd/rest`

* Start Fiber server
* Load config
* Initialize dependencies
* Register routes
* Start HTTP service

### `cmd/worker`

* Background jobs
* Async processing
* Cron / queue consumers

### `cmd/seeder`

* Seed database
* Insert initial data
* Development utilities

No business logic allowed here.

---

## 2. `internal/adapter/` — Presentation Layer

Handles external communication (HTTP).

### `handler/`

Responsibilities:

* Parse request input
* Validate input
* Instantiate usecases
* Call usecase methods
* Format responses

Rules:

* No business rules
* No database access
* No infrastructure imports

Example:

```go
type UserHandler struct {
    usecase usecase.UserUsecase
}
```

---

### `router/`

Responsibilities:

* Register endpoints
* Bind routes to handlers
* Register middleware

Rules:

* No logic
* Routing only

---

## 3. `internal/config/` — Configuration

### `config.go`

Responsibilities:

* Load env variables
* Define config struct
* Centralize app settings

Rules:

* Single source of truth
* No business logic

---

## 4. `internal/database/` — Database Layer

### `migrations/`

* SQL migration files
* Versioned schema changes
* Managed by migration tools

Rules:

* No Go logic here
* Only schema-related files

### `sqlc/queries/`

* SQL query definitions
* One file per domain (`user.sql`, `post.sql`, etc.)
* Auto-discovered by SQLC

Rules:

* Write standard SQL with SQLC syntax
* Name queries clearly: `-- name: GetUserByID :one`
* Use parameterized queries: `$1, $2, $3`

### `sqlc/generated/`

* Auto-generated Go code
* DO NOT edit manually
* Regenerate with `sqlc generate` after query changes

Contains:

* Type-safe query functions
* Struct definitions matching database schema
* Database interface

---

## 5. `internal/domain/` — Core Domain

Contains pure business concepts.

Must not depend on:

* Fiber
* Database drivers
* Infrastructure
* External APIs

---

### `entities/`

Domain models.

Rules:

* One entity per file
* Plain Go structs
* Domain rules only

Example:

```text
auth.go → User, Session
order.go → Order, Payment
```

---

### `constants/`

Business-context constants.

Responsibilities:

* Store constants related to one domain
* Define typed constants

Rules:

* One context per file
* `common.go` for shared domain types

Example:

```go
type UserStatus string
```

---

### `port/`

Domain interfaces.

Responsibilities:

* Define contracts
* Describe required behaviors

Rules:

* One interface per file
* No implementation

Example:

```go
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*entities.User, error)
}
```

---

## 6. `internal/usecase/` — Application Logic

Main business logic layer.

Responsibilities:

* Implement workflows
* Orchestrate domain + ports
* Apply business rules
* Validate domain behavior

Rules:

* Depends on domain + ports
* No infrastructure imports
* No HTTP logic

Example:

```go
type UserUsecase struct {
    repo port.UserRepository
}
```

Usecases are the **heart of the system**.

---

## 7. `internal/infrastructure/` — Technical Implementation

Implements ports.

---

### `adapter/`

Responsibilities:

* Database repositories
* External API clients
* Cache adapters
* Message brokers

Each folder implements domain ports.

Example:

```text
adapter/
├── postgres
├── redis
└── smtp
```

---

### `pkg/`

Shared infrastructure utilities.

#### `shared/`

* Global constants
* Cross-domain values

Example:

```go
const AppName = "MyService"
```

#### `utils/`

Utility helpers.

Rules:

* One context per file
* Stateless helpers only

Example:

```text
string.go
time.go
crypto.go
```

---

## Dependency Rules (Strict)

Allowed dependency flow:

```text
Router → Handler → Usecase → Port → Infrastructure
                     ↓
                  Domain
```

### Allowed

✅

* Handler → Usecase
* Usecase → Port
* Infrastructure → Port
* Usecase → Domain

### Forbidden

❌

* Domain → Infrastructure
* Domain → Usecase
* Port → Implementation
* Usecase → HTTP/DB

---

## Application Flow

Typical request:

```text
HTTP Request
   ↓
Router
   ↓
Handler
   ↓
Usecase
   ↓
Port Interface
   ↓
Infrastructure Adapter
   ↓
Database / External Service
```

---

## Dependency Injection

All wiring happens in:

```text
cmd/rest
cmd/worker
cmd/seeder
```

Responsibilities:

* Initialize DB
* Initialize repositories
* Initialize usecases
* Initialize handlers

No `new()` inside domain.

---

## Naming Conventions

| Item      | Convention    |
| --------- | ------------- |
| Package   | snake_case    |
| File      | snake_case.go |
| Struct    | PascalCase    |
| Interface | PascalCase    |
| Method    | PascalCase    |

Example:

```text
user_handler.go
user_usecase.go
postgres_user_repo.go
```

---

## Error Handling

Rules:

* Return `error`
* Wrap with context
* No panic in production
* Centralized error handler

Example:

```go
fmt.Errorf("create user: %w", err)
```

---

## Testing Strategy

| Layer          | Type        |
| -------------- | ----------- |
| Domain         | Unit        |
| Usecase        | Unit        |
| Port           | Mock        |
| Infrastructure | Integration |
| Handler        | HTTP        |

Domain & usecase tests must not use DB.

### Tooling

* Use the Go standard library `testing` package — no external assertion or
  mocking framework.
* Repository / port doubles live in [`internal/tests/mocks`](internal/tests/mocks).
  Each mock exposes one function field per interface method, so individual
  tests can override only the methods they exercise. Un-stubbed methods return
  an "unimplemented" error to make accidental usage visible.
* For asserting on returned `*domain.AppError` values, use
  `utils.AssertAppErr(t, err, want)` from
  [`internal/tests/utils`](internal/tests/utils) — it checks the error
  type via `domain.IsAppError`.

### Layout

* Tests live under [`internal/tests`](internal/tests), grouped by domain into
  one package per subfolder (`auth`, `professor`, `project`, `student`). Each
  test file imports the usecase package it exercises (e.g.
  `professor_usecase "profconnect-api/internal/usecase/professor"`) and drives
  it through its exported constructor. Shared port mocks live at
  [`internal/tests/mocks`](internal/tests/mocks); shared assertion helpers at
  [`internal/tests/utils`](internal/tests/utils).
  ```text
  internal/tests/
  ├── mocks/                       # hand-rolled port doubles
  ├── utils/                       # AssertAppErr and other shared helpers
  ├── auth/
  │   ├── auth_login_test.go
  │   ├── auth_refresh_test.go
  │   ├── auth_register_core_test.go
  │   ├── auth_register_admin_test.go
  │   ├── auth_register_professor_test.go
  │   └── auth_register_student_test.go
  ├── professor/
  │   ├── create_project_test.go
  │   ├── list_applications_test.go
  │   ├── list_invitations_test.go
  │   ├── list_students_test.go
  │   ├── professor_profile_test.go
  │   ├── send_invitation_test.go
  │   ├── cancel_invitation_test.go
  │   ├── review_application_test.go
  │   └── remove_member_test.go
  ├── project/
  │   ├── get_project_test.go
  │   ├── list_projects_test.go
  │   └── list_members_test.go
  └── student/
      ├── apply_project_test.go
      ├── check_application_test.go
      ├── my_applications_test.go
      ├── list_my_invitations_test.go
      ├── list_my_projects_test.go
      ├── student_profile_test.go
      ├── respond_invitation_test.go
      ├── leave_project_test.go
      └── withdraw_application_test.go
  ```
* Every usecase under [`internal/usecase/<domain>`](internal/usecase) has a
  matching `<name>_test.go` file in `internal/tests/<domain>/`. When you add a
  new usecase, drop a sibling test file into the matching domain folder.
* Use `utils.AssertAppErr(t, err, want)` to verify a usecase mapped a
  precondition to the expected `domain.ErrorType`.
* Because tests are out-of-package, they can only touch the exported API of
  each usecase. That's a feature — the same surface the rest of the app
  consumes.
* For usecases with several preconditions, use a table-driven `t.Run` for the
  validation/rejection branches, and dedicated test functions for the
  multi-step success paths.

### Writing a usecase test

1. Build a small fixtures helper (e.g. `sendInvitationFixtures()`) returning a
   stock project, professor, and student — tests mutate fields as needed.
2. Construct the mocks for the ports the usecase depends on.
3. Stub only the methods the test cares about (the rest will panic / return
   `unimplemented`, surfacing accidental coupling).
4. Call `Execute(ctx, &input)` and assert on the output / error type.

Example (abbreviated):

```go
package professor

import (
    "context"
    "testing"

    "profconnect-api/internal/domain/constants"
    "profconnect-api/internal/domain/entities"
    inputoutput "profconnect-api/internal/domain/input_output"
    "profconnect-api/internal/tests/mocks"
    professor_usecase "profconnect-api/internal/usecase/professor"
)

func TestSendInvitation_Success(t *testing.T) {
    project, professor, student := sendInvitationFixtures()
    projectRepo := &mocks.ProjectRepository{}
    invRepo := &mocks.ProjectInvitationRepository{}
    // ... other repos

    projectRepo.GetByIDFn = func(_ context.Context, _ string) (*entities.Project, error) {
        return project, nil
    }
    invRepo.CreateFn = func(_ context.Context, inv *entities.ProjectInvitation) (*entities.ProjectInvitation, error) {
        return &entities.ProjectInvitation{ID: "inv-1", Status: inv.Status}, nil
    }
    // ... stub the rest

    uc := professor_usecase.NewSendInvitationUsecase(projectRepo, invRepo, /* ... */)
    out, err := uc.Execute(context.Background(), &inputoutput.SendInvitationInput{
        ProfessorUserID: professor.User.ID,
        ProjectID:       project.ID,
        StudentID:       student.ID,
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if out.ID != "inv-1" {
        t.Fatalf("unexpected output: %+v", out)
    }
}
```

### Running tests

```bash
# Run everything
make test

# Run every test package
go test ./internal/tests/...

# Run only one domain's tests
go test ./internal/tests/professor/

# Run a single test
go test -run TestSendInvitation_Success ./internal/tests/professor/

# Verbose output
go test -v ./internal/tests/...
```

### What to cover

For each new usecase, prioritise:

* The happy path (writes happen, status transitions correctly).
* Any side-effect that depends on the slot accounting (member created on
  approve/accept; project closes when full; project reopens when a slot frees).
* Each precondition branch that returns a `domain.AppError` (validation /
  forbidden / not-found / conflict).

Pure list/get usecases that are mostly straight-through don't need exhaustive
coverage.

---

## General Principles

* Keep domain pure
* Keep usecases thin
* Prefer composition
* No circular deps
* Small focused packages
* Explicit wiring
* Interface-driven design

---

## Database Migrations (golang-migrate)

### Overview

Database migrations are managed using `golang-migrate`. Migration files are in `migrations/` directory.

### Migration File Naming

Follow the strict convention:

```
YYYYMMDDHHMMSS_description.up.sql
YYYYMMDDHHMMSS_description.down.sql
```

**Or sequential:**

```
000001_description.up.sql
000001_description.down.sql
```

### Migration Files Structure

**Up migrations** (`*.up.sql`): Apply schema changes
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);
```

**Down migrations** (`*.down.sql`): Revert schema changes
```sql
DROP TABLE IF EXISTS users;
```

### Making Migrations

```bash
# Create a new migration pair
make migrate-create NAME=add_user_table
# Creates: migrations/000002_add_user_table.up.sql
#          migrations/000002_add_user_table.down.sql
```

### Running Migrations

```bash
# Apply all pending migrations
make migrate-up

# Revert last migration
make migrate-down

# Apply to running container
docker-compose exec -T postgres migrate -path /migrations -database "postgres://..." up
```

### Setup Database

```bash
# Start PostgreSQL and apply migrations automatically
make setup-db
```

---

## SQLC (SQL Compiler)

### Overview

SQLC generates type-safe Go code from SQL queries. Instead of writing raw SQL strings in Go, you:

1. Write SQL queries in `.sql` files
2. Run `sqlc generate` to create Go code
3. Use generated functions in repositories

**Benefits:**
- ✅ Compile-time query validation
- ✅ Type-safe parameter binding
- ✅ Auto-generated repository code
- ✅ IDE autocomplete for queries
- ✅ Easy to add complex queries

### Directory Structure

```text
internal/database/sqlc/
├── queries/
│   ├── user.sql
│   ├── post.sql
│   └── comment.sql
└── generated/
    ├── db.go           # Auto-generated
    ├── queries.go      # Auto-generated
    └── models.go       # Auto-generated
```

### Configuration

**sqlc.yaml** (root of project):

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/database/sqlc/queries"
    schema: "migrations"
    gen:
      go:
        package: "generated"
        out: "internal/database/sqlc/generated"
        emit_prepared_queries: false
        emit_interface: true
```

### Writing Queries

**internal/database/sqlc/queries/user.sql:**

```sql
-- name: GetUserByID :one
SELECT id, name, email, hashed_password FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, hashed_password FROM users WHERE email = $1;

-- name: ListAllUsers :many
SELECT id, name, email, hashed_password FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (name, email, hashed_password) 
VALUES ($1, $2, $3) 
RETURNING id, name, email, hashed_password;

-- name: UpdateUser :one
UPDATE users 
SET name = $1, email = $2, hashed_password = $3 
WHERE id = $4 
RETURNING id, name, email, hashed_password;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;
```

**Syntax:**
- `-- name: FunctionName :one` - Returns single row
- `-- name: FunctionName :many` - Returns multiple rows
- `-- name: FunctionName :exec` - Executes without return
- `$1, $2, $3` - Parameterized placeholders

### Generating Code

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate Go code from SQL queries
sqlc generate
```

### Using Generated Code in Repositories

**internal/infrastructure/adapter/postgres/user_repository.go:**

```go
package postgres

import (
    "context"
    "profconnect-api/internal/database/sqlc/generated"
    "profconnect-api/internal/domain/entities"
    "profconnect-api/internal/domain/port"
)

type UserRepository struct {
    queries *generated.Queries
}

func NewUserRepository(queries *generated.Queries) port.UserRepository {
    return &UserRepository{queries: queries}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
    user, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get user by id: %w", err)
    }
    
    return mapToEntity(user), nil
}

// Helper function to map from generated type to domain entity
func mapToEntity(u generated.User) *entities.User {
    return &entities.User{
        ID:             u.ID,
        Name:           u.Name,
        Email:          u.Email,
        HashedPassword: u.HashedPassword,
    }
}
```

### Dependency Injection

**cmd/rest/main.go:**

```go
func main() {
    db := initDatabase() // PostgreSQL connection
    
    // Initialize sqlc queries
    queries := generated.New(db)
    
    // Initialize repositories with generated queries
    userRepository := postgres.NewUserRepository(queries)
    
    // Initialize usecases
    registerUsecase := usecase.NewRegisterUsecase(userRepository)
    
    // ... rest of setup ...
}
```

### Advantages Over Raw SQL

| Aspect | Raw SQL | SQLC |
|--------|---------|------|
| Type Safety | ❌ Runtime errors | ✅ Compile-time |
| IDE Support | ❌ No | ✅ Full autocomplete |
| Query Validation | ❌ Runtime | ✅ At generation time |
| Boilerplate | ❌ Manual Scan() | ✅ Auto-generated |
| Maintainability | ❌ Scattered SQL | ✅ Centralized queries |
| Performance | ✅ Same | ✅ Same (zero overhead) |

---

## Development Workflow

### Live Reload with Air

Use `air` for automatic app restart on file changes.

```bash
# Start development server with live reload
make dev
```

Air configuration: `.air.toml`

**Features:**
- ✅ Watches Go files (except `*_test.go`)
- ✅ Rebuilds on changes (1s delay)
- ✅ Colored output
- ✅ Ignores vendor, tmp directories

### Initial Setup

```bash
# Install all tools
make setup

# This installs:
# - golang-migrate CLI
# - air (live reload)
# - sqlc (SQL compiler)
```

### Development Commands

```bash
# Setup database with migrations
make setup-db

# Start with live reload
make dev

# Or run directly
make run

# View docker logs
make docker-logs

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean
```

### Makefile Variables

Override default values:

```bash
# Custom database config
make migrate-up DB_HOST=db.example.com DB_USER=admin

# Or set environment variables
export DB_HOST=db.example.com
make migrate-up
```

**Default values:**
```makefile
DB_HOST = localhost
DB_PORT = 5432
DB_USER = postgres
DB_PASSWORD = postgres
DB_NAME = profconnect
DB_SSLMODE = disable
```

---

## Summary

This backend must:

* Use Go + Fiber
* Follow clean architecture
* Separate concerns
* Centralize business logic in usecases
* Use ports for abstraction
* Keep domain independent
* Wire everything in `cmd`
* Use golang-migrate for schema versioning
* Use SQLC for type-safe database queries
* Use Air for development live reload
* Follow Makefile conventions for common tasks