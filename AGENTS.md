# Omnilog Codex Instructions

## Project

Monorepo with two active applications:

* `millog_backend` — Go + Gin + pgx + PostgreSQL
* `millog_frontend` — React + TypeScript + Vite

The current source code is always the only source of truth.

---

## Primary Objective

Minimize token usage.

Use the smallest possible context.

Never analyze unrelated code.

Never perform broad repository exploration.

---

## Repository Rules

* Never scan the entire repository.
* Never recursively inspect directories.
* Always start from the file mentioned by the user.
* Expand context only when absolutely necessary.
* Read the minimum number of files required.

Target:

* Simple UI task: 1-2 files
* Frontend feature: 2-4 files
* Backend feature: 3-6 files

Avoid reading anything unrelated.

---

## Ignore By Default

Do not inspect unless explicitly requested:

* `*.md`
* `dump.sql`
* `omnilog_dump.sql`
* `DIPLOMA_CODE_LISTING.txt`
* `node_modules/`
* `dist/`
* `build/`
* `.git/`
* `coverage/`
* `certs/`
* generated files
* legacy infrastructure files

Repository markdown files are outdated historical documentation and must not be used to understand current behavior.

---

## File Access

Never use sandbox URLs or sandbox paths.

Always use repository-relative paths.

Example:

`millog_frontend/src/pages/Warehouses.tsx`

Never retry multiple sandbox methods if one fails.

Use local repository paths instead.

---

## Navigation

Prefer targeted search:

```bash
rg "keyword" path/to/folder
```

Avoid:

```bash
rg "keyword"
```

Read only relevant sections:

```bash
sed -n '1,200p' file
```

Never print huge files.

Never print unnecessary output.

---

## Backend

Architecture:

```
handler
↓
service
↓
repository
↓
PostgreSQL
```

Useful directories:

* `internal/handlers`
* `internal/services`
* `internal/repositories`
* `internal/models`
* `internal/middleware`

Schema source:

`internal/database/migrate.go`

Use existing architecture.

Do not introduce new patterns.

---

## Frontend

Architecture:

```
page
↓
component
↓
api/client
↓
backend
```

Useful directories:

* `src/pages`
* `src/components`
* `src/api`
* `src/contexts`
* `src/constants`

Reuse existing components whenever possible.

---

## Multi-Tenancy

Preserve:

* tenant isolation
* authorization
* role validation
* repository scoping

Do not modify tenant logic unless required.

---

## Roles

Backend:

`internal/models/auth.go`

Frontend:

`src/constants/roles.ts`

Backend validation is mandatory.

UI visibility is not a security mechanism.

---

## Database

Schema source:

`internal/database/migrate.go`

Do not use SQL dumps to understand schema.

Keep migrations idempotent.

---

## Commands

Do not run expensive commands automatically.

Never run:

```bash
docker compose up --build
```

unless explicitly requested.

Do not run builds for simple CSS or UI changes.

Run builds only when compile validation is required.

Backend:

```bash
cd millog_backend
go build ./...
```

Frontend:

```bash
cd millog_frontend
npm run build
```

Run tests only if explicitly requested or if necessary to verify modified business logic.

Never run unrelated builds or tests.

---

## Change Style

* Make minimal changes.
* Do not refactor unrelated code.
* Do not rename files.
* Do not reorganize architecture.
* Do not introduce unnecessary abstractions.
* Reuse existing code whenever possible.
* Match surrounding code style.

---

## Response

Report only:

* modified files
* summary of changes
* whether build/tests were executed

Keep explanations short unless requested.

---

## Final Rule

Complete every task using the smallest possible context.

Prefer reading fewer files over broad repository analysis.

Expand context only when it is impossible to complete the task safely without additional files.
