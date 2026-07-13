---
name: appcrons-backend
description: Use when working in this repo (the appcrons/"keep-active" Go backend) — adding or changing API endpoints, models, middleware, scheduling/event logic, auth, or tests. Covers architecture, routing, DB, auth, config, and build/test conventions actually used here.
---

# Appcrons backend (Go)

REST API + background scheduler that keeps free-tier backend hosts (e.g. Render) awake by
periodically pinging user-registered app URLs on a schedule, so they don't spin down from
inactivity. Also exposes an SSE stream for live request progress and a Prometheus `/metrics`
endpoint.

- Module path: `github.com/Tibz-Dankan/keep-active` — **note this differs from the repo/product
  name "appcrons"**. All internal imports use this module path, not `appcrons`.
- Go version: `go.mod` declares `go 1.22`, README says `1.23`, Dockerfile builds with
  `golang:1.23.0-alpine`. Slight mismatch across the repo — don't be surprised by it, just build
  with whatever's already pinned in the Dockerfile/CI.
- Server listens on port `8080`.

## Directory map

```
cmd/main.go                 entrypoint: wires router, CORS, schedulers, event bus, http.ListenAndServe
internal/
  apperror/                 sentinel errors (currently just ErrNoUser; lightly used)
  constants/                global consts (query limits, error strings)
  events/                   custom in-process pub/sub event bus
    eventBus.go              topic -> channel-of-subscribers
    publishers/               e.g. PublishRequestEvent -> "makeRequest" topic
    subscribers/               e.g. subscribes "makeRequest" -> MakeAppRequest; also permissions
  middlewares/
    auth.go                  JWT bearer auth, injects userId into request context
    permissions.go            HasPermissions: path-regex-based RBAC/ACL check
    logger.go                 colorized request logger
    rateLimiter.go            in-memory per-IP limiter (20 req/60s), skipped in testing/staging
    requestDuration.go        Prometheus histogram middleware
  models/                    GORM models AND all data-access methods (repository layer lives here)
    db.go                     Postgres/SQLite connection (env-driven) + AutoMigrate
    schema.go                 struct defs: User, App, Request, RequestTime, Feedback, OTP, BugReport, RequestCount
    app.go, user.go, request.go, requestTime.go, feedback.go, bugReport.go, opt.go, requestCount.go, ...
    cache.go                  Redis client init
    permissions.go            UserPermissions model (Redis-cached dev/prod, in-RAM map testing/staging)
  routes/
    router.go                 AppRouter(): mounts all subrouters + middleware chains
    admin/ app/ auth/ bugReport/ feedback/ monitor/ request/   one file per endpoint
  schedulers/
    schedulers.go              InitSchedulers(): kicks off background goroutines
    request.go                  5-min-aligned tick -> PublishRequestEvent
    UserAppMemory.go            periodic cleanup of in-memory per-user app request state
  services/                  flat grab-bag of stateless helpers — NOT domain services
    jwt.go, email.go, upload.go, firebaseAdmin.go, httpRequest.go, client.go, date.go, error.go, ...
  templates/email/            HTML email templates (reset-password.html, otp.html)
tests/                       Go integration tests (see Testing below) — nothing under internal/
testsprite_tests/            separate TestSprite-MCP-generated Python test suite, NOT part of `make test` or CI
```

## Architecture: "fat model / thin handler"

There is no separate service/repository layer. Two layers only:

1. **Handlers** in `internal/routes/<domain>/<verb><Noun>.go`. Each file has a handler function
   plus an exported `<Verb><Noun>Route(router *mux.Router)` that registers it:

   ```go
   // internal/routes/app/postApp.go
   func PostApp(w http.ResponseWriter, r *http.Request) { ... }
   func PostAppRoute(router *mux.Router) {
       router.HandleFunc("/post", PostApp).Methods("POST")
   }
   ```

   Handlers decode JSON directly (`json.NewDecoder(r.Body).Decode(&app)`), validate inline, call
   model methods, and write JSON responses by hand (`w.WriteHeader(...)`,
   `json.NewEncoder(w).Encode(...)`). No handler-level DTO/validation library.

2. **Models** in `internal/models/<entity>.go` hold the struct (in `schema.go`) AND all
   data-access as receiver methods (`Create`, `FindOne`, `FindByUser`, `Update`, `Delete`, ...).
   This _is_ the repository layer — don't create a separate `repository`/`store` package for new
   entities, follow this pattern instead. GORM `BeforeCreate` hooks generate UUID PKs and set
   defaults (e.g. `App.BeforeCreate` sets `IsDisabled = true`).

`internal/services/` is a flat bucket of stateless helpers (JWT signing, email, Firebase upload,
outbound HTTP client, date/URL utilities) — not domain services tied to a specific entity.

**End-to-end example — `POST /api/v1/apps/post`:**

1. `middlewares.Auth` parses the bearer JWT, loads the user via `models.User.FindOne`, injects
   `userId` into context.
2. `middlewares.HasPermissions` loads `models.Permissions.Get(userId)` (Redis in dev/prod, in-RAM
   map in testing/staging) and checks the required permission for this path via regex matching.
3. `app.PostApp` reads `userId` from context, decodes the body into `models.App`, validates
   required fields inline, checks uniqueness (`FindByName`/`FindByURL`), calls `app.Create(app)`.
4. Permissions are refreshed synchronously in testing/staging, or asynchronously via
   `events.EB.Publish("permissions", user)` in dev/prod.
5. Response is written directly as JSON with `status`/`message`/`app`/`data` keys.

## Routing

- Router: **gorilla/mux** (not net/http ServeMux, chi, gin, echo, or fiber).
- `AppRouter()` in `internal/routes/router.go` creates one `mux.NewRouter()`, applies global
  middleware via `router.Use(...)` (`RequestDuration`, `Logger`, `RateLimit`), then mounts
  path-prefixed subrouters per domain, each with its own `.Use(middlewares.Auth,
middlewares.HasPermissions)` where auth is required:
  - `/api/v1/apps` — Auth + HasPermissions
  - `/api/v1/requests` — Auth + HasPermissions
  - `/api/v1/auth` — signup/signin/forgot/reset are unauthenticated; a second
    `authorizedAuthRouter` on the same prefix adds Auth + HasPermissions for
    update-details/change-password/get-user
  - `/api/v1/admin` — Auth + HasPermissions
  - `/api/v1/feedback` — Auth + HasPermissions
  - `/api/v1/bugreport` — unauthenticated (bug reports can be filed without login)
  - plus `GetActiveRoute` (health check), `monitor.GetMetrics` (Prometheus), `NotFoundRoute`
- CORS is wrapped around the mux router in `cmd/main.go` via `github.com/rs/cors`
  (`AllowedOrigins: ["*"]`, GET/POST/PATCH/DELETE, `AllowCredentials: true`).

## Scheduling & events (the core product feature)

Custom in-process pub/sub (`internal/events`), not a message queue:

1. `schedulers.InitSchedulers()` starts a goroutine that ticks on 5-minute-aligned boundaries.
2. `publishers.PublishRequestEvent()` fetches all apps (`models.App.FindAll`) and publishes each
   to the `"makeRequest"` topic on the event bus.
3. A subscriber (`internal/events/subscribers/request.go`) receives each app and calls
   `request.MakeAppRequest(app)`, which validates the app's interval/time window
   (`validateApp`), performs the HTTP GET (`services.MakeHTTPRequest`/`MakeExternalHTTPRequest`),
   saves a `models.Request` row, and publishes progress on `"appRequestProgress"`.
4. The SSE handler (`getLiveRequests.go`) relays that progress topic to connected clients as
   `text/event-stream`.

When adding a new async/decoupled flow, follow this publish/subscribe pattern rather than calling
across packages directly.

## Database

- GORM (`gorm.io/gorm`) with `gorm.io/driver/postgres` and `gorm.io/driver/sqlite`.
- DB selection is driven by `GO_ENV` in `internal/models/db.go`:
  - `development` → Postgres via `APPCRONS_DEV_DSN`
  - `production` → Postgres via `APPCRONS_PROD_DSN`
  - `testing` / `staging` → **opens a local SQLite file** (`./../../appcrons_test.db` relative to
    `internal/models`) instead of actually using `APPCRONS_TEST_DSN`/`APPCRONS_STAG_DSN`. Those
    DSN env vars are read but not used for the connection in these environments — don't expect
    setting them to change which DB test runs hit.
- **No migration tool** — schema is managed entirely by `gormDB.AutoMigrate(&User{}, &App{},
&Request{}, &RequestTime{}, &Feedback{}, &OTP{}, &BugReport{}, &RequestCount{})` in `Db()`. Add
  new models to this call.
- String UUID primary keys generated in `BeforeCreate` hooks (`google/uuid`), soft deletes via
  `gorm.DeletedAt` on most models (not `Feedback`/`OTP`/`BugReport`/`RequestCount`).
- **Columns are camelCase, not GORM's default snake_case** (`gorm:"column:userId"`). Any new raw
  query must quote camelCase column names, e.g. `db.Where("\"userId\" = ?", userId)`.
- Redis (`internal/models/cache.go`) backs the permissions cache in dev/prod only; testing/staging
  use an in-process RAM map instead (`RedisClient()` has no testing/staging case — it would
  `log.Fatal` if reached, which is why permissions bypass Redis in those envs).

## Auth & permissions

- JWT bearer auth, HS256, `github.com/golang-jwt/jwt` **v3** (`+incompatible`), secret from
  `JWT_SECRET`. Claims: `userId`, `exp` (9h from issuance), `iat`. Issued via
  `services.SignJWTToken` (`internal/services/jwt.go`).
- `middlewares.Auth` parses `Authorization: Bearer <token>`, verifies HMAC signing method, loads
  the user via `models.User.FindOne`, rejects if the user no longer exists, injects `userId` into
  context under `middlewares.UserIDKey`.
- `middlewares.HasPermissions` is a custom, path-regex-driven RBAC layer (not a library):
  - Role is `user` or `sys_admin` on `models.UserPermissions`. `sys_admin` bypasses all checks.
  - Ordinary `user` gets `["READ","WRITE","EDIT","DELETE"]`; default `sys_admin` permission is
    `READ` only.
  - `getPermissionID(r)` inspects the URL path (and sometimes decodes the body, e.g. for
    `post-request-time`) to determine the required `{ID, Type, Permission}` tuple, then checks the
    resource belongs to the authenticated user and their permission list allows it.
- Passwords: bcrypt cost 12 (`models.User.HashPassword`/`PasswordMatches`). Reset tokens: random
  UUID, SHA-256 hashed at rest, 20-minute expiry.
- Separate admin auth endpoints: `POST /api/v1/auth/signup-admin`, `POST
/api/v1/auth/signin-admin`.

## Config

`GO_ENV` (`development` | `testing` | `staging` | `production`) is the master switch controlling
DSN/DB driver, Redis var, rate-limiting (skipped outside dev/prod), sync-vs-async permission
writes, and whether real emails are sent. Unrecognized `GO_ENV` values call `log.Fatal`.

`.env` / `appcrons.env` are loaded via `godotenv.Load()` **only when `GO_ENV=development`**; other
environments expect vars to already be present in the process environment (CI/CD or the host).
Config access throughout is plain `os.Getenv(...)` — no Viper or similar.

Env var categories present (names only):

- DB DSNs: `APPCRONS_DEV_DSN`, `APPCRONS_TEST_DSN`, `APPCRONS_STAG_DSN`, `APPCRONS_PROD_DSN`
- Redis: `REDIS_DEV_URL`, `REDIS_PROD_URL`
- App: `APPCRONS_EXTERNAL_URL`, `ADMIN_EMAIL`, `ADMIN_EMAIL_APPCRONS`, `JWT_SECRET`
- Email: `GMAIL_SENDER_MAIL`, `GMAIL_APP_PASSWORD`, `MJ_APIKEY_PUBLIC`, `MJ_APIKEY_PRIVATE`,
  `MJ_SENDER_MAIL`
- Firebase (service account): `FIREBASE_TYPE`, `FIREBASE_PROJECT_ID`, `FIREBASE_PRIVATE_KEY_ID`,
  `FIREBASE_PRIVATE_KEY`, `FIREBASE_CLIENT_EMAIL`, `FIREBASE_CLIENT_ID`, `FIREBASE_AUTH_URI`,
  `FIREBASE_TOKEN_URI`, `FIREBASE_AUTH_PROVIDER_X509_CERT_URL`, `FIREBASE_CLIENT_X509_CERT_URL`,
  `FIREBASE_UNIVERSE_DOMAIN`, `FIREBASE_STORAGE_BUCKET`

## Testing

- Go tests live **only** under `tests/` — there are no `_test.go` files inside `internal/`.
- Stdlib `testing` + `httptest` only, no testify. Tests drive the full router
  (`routes.AppRouter()`) via `tests/setup/setup.go`'s `ExecuteRequest`, so they're effectively
  HTTP-level integration tests against a real (SQLite, in testing env) DB.
- Shared helpers: `tests/setup/setup.go` (`ExecuteRequest`, `ClearAllTables`, `CreateSignInUser`,
  `SignInUser`, `CheckResponseCode`), `tests/data/testData.go` (random unique test data
  generator).
- Feature test files live under `tests/app/` and `tests/auth/`.
- Run with `make test` (`GO_ENV=testing`) or `make stage` (`GO_ENV=staging`).
- `testsprite_tests/` is a **separate** Python E2E suite generated by TestSprite MCP — it is not
  Go, not run by `make test`, and not wired into CI. Don't confuse it with the real test suite.

## Build / run / deploy

```
make run       # GO_ENV=development go run ./cmd
make test      # GO_ENV=testing go test -v ./tests/...
make stage     # GO_ENV=staging go test -v ./tests/...
make install   # go mod tidy && go mod download
```

Dockerfile: `golang:1.23.0-alpine`, builds `./bin/appcrons` from `./cmd`, `ENV
GO_ENV=production`, exposes `8080`.

CI (`.github/workflows/deploy-to-render.yaml`, triggered on push to `main`) is **deploy-only** —
it checks out and deploys straight to Render via `JorgeLNJunior/render-deploy@v1.4.4`. A more
complete `test-deploy` job (checkout → Go 1.23 setup → `make install` → `make stage` → deploy) is
present in the file but **entirely commented out**, so pushes to `main` currently deploy to
production without running any tests in CI. Keep this in mind — passing `make test`/`make stage`
locally is currently the only gate before deploy.

## Code-style conventions

- No linter is configured (no `.golangci.yml`, no `.editorconfig`).
- Logging is stdlib `log` only (no zap/logrus). `internal/middlewares/logger.go` has a custom
  colorized request logger.
- HTTP error responses go through the shared `services.AppError(message string, statusCode int,
w http.ResponseWriter)` helper — use it for new handlers instead of writing errors ad hoc.
- File naming: route files are `<verb><Noun>.go` (e.g. `postApp.go`) with a matching handler and
  an exported `<Verb><Noun>Route(router)` registration function; model files are named after the
  entity, with `schema.go` centralizing struct + GORM tag definitions.
- For new env-specific behavior, follow the existing `GO_ENV`-based branching pattern rather than
  introducing a new config mechanism.
- Some model query methods retain large commented-out "previous implementation" blocks next to
  the current optimized version (e.g. `App.FindByUser`/`FindAll` in `internal/models/app.go`) —
  this is existing style, not something to necessarily clean up unless asked.
