# Go Implementation Audit

The implementation has a solid Go foundation, but it is not production-ready. The primary blockers are unresolved domain policies, data-model inconsistencies, unsafe request handling, and almost entirely untested domain logic.

## 1. Architecture Decisions

1. **Define domain ownership and cardinality.** Every authenticated user can create or modify global projects (`internal/project/routes.go:8-13`). The API permits multiple groups per user, and starting a project can create unlimited duplicate progress records. Decide:
   - Are projects admin-managed or globally editable?
   - Is there one group per user or many?
   - Can a project be started more than once?

   Enforce those decisions with authorization and unique indexes.

2. **Make server-owned state explicit.** Project progress is calculated by one endpoint but can be directly overwritten through the generic update endpoint (`internal/project/model.go:96-102`). Decide which fields are client-editable and which are derived by the server, then enforce that boundary in request types rather than relying on handler behavior.

3. **Keep the vertical-slice architecture, but introduce test seams.** The handler/service/repository separation is sensible. However, services depend directly on concrete Mongo repositories, and some services accept BSON filters (`internal/user/service.go:36-38`). Use small consumer-owned interfaces around repositories and the Together client where tests require them. Define those interfaces beside the consuming service and keep test fakes in `*_test.go`; this does not require a new production package. Avoid a full framework or generic "clean architecture" rewrite.

4. **Treat persistence shape as an explicit application contract.** Collection names, indexes, BSON field names, and migration ownership should be defined centrally. The current startup only creates user and refresh-token indexes (`internal/db/mongo.go:47-79`), leaving feature-level uniqueness and query performance implicit.

   Suggested structure if index management is split from connection management:

   ```text
   internal/db/
   |-- indexes.go  # All collection index definitions
   `-- mongo.go    # Connection, ping, and shutdown lifecycle
   ```

5. **Reconsider mandatory AI startup coupling.** `run()` refuses to start any route without a Together key (`cmd/server/main.go:68-102`). If AI is optional, only disable AI routes or expose them as unavailable. Also separate liveness from readiness instead of always returning healthy.

6. **Harden the cookie architecture.** Production cookies use `SameSite=None` (`internal/auth/cookie.go:9-34`) without CSRF protection. CORS does not prevent cross-site requests. If cross-site cookies are required, validate `Origin` on unsafe methods or add CSRF tokens. Otherwise prefer `SameSite=Lax`. The suggested middleware file placement is included in the `internal/api` structure below.

## 2. File Structure

### What Works

- `cmd/server` is a clear composition root.
- `internal/` correctly prevents accidental external imports.
- Feature-oriented packages are easier to navigate than global handler/service/model directories.
- Manual dependency wiring is readable at the current scale (`cmd/server/main.go:78-116`).

### What Should Improve

- `internal/project/model.go` combines persisted models and HTTP update payloads. The current tags are corrected, but keeping transport and persistence concerns together makes future schema changes riskier. Split request payloads from persisted models where their representations differ and add explicit JSON and BSON tags to nested persisted types.

  Suggested minimal structure:

  ```text
  internal/project/
  |-- errors.go
  |-- handler.go
  |-- model.go       # Project, ProjectProgress, and nested persisted types
  |-- payload.go     # Create and update request DTOs
  |-- repository.go
  |-- routes.go
  `-- service.go
  ```

- Centralize HTTP error translation in `internal/api`. Handlers currently duplicate response logic and disagree on plain text versus JSON (`internal/auth/handler.go:21-120`, `internal/project/handler.go:17-155`).

  Suggested structure:

  ```text
  internal/api/
  |-- cors.go
  |-- errors.go      # Domain error to HTTP status/public response mapping
  |-- handlers.go
  |-- origin.go      # Unsafe-method Origin or CSRF enforcement
  `-- response.go
  ```

- Normalize naming: `route.go` versus `routes.go`, `ListGroup` versus `ListGroups`, singular Mongo collection `group`, and exported project collection fields.

- Remove incomplete or unused pieces such as `ChatURLv2`, `verifyRefreshToken`, unused AI message persistence, and unused model constants unless implementation is imminent.

- Add `README.md`, an environment template, and `.dockerignore`. Credential-like values exist in ignored local files, but `Dockerfile:10` sends them into the Docker build context because Git ignore rules do not apply there. Rotate them if they are real.

  Suggested root structure:

  ```text
  implementations/go/
  |-- .dockerignore
  |-- .env.example   # Variable names and safe placeholders only
  |-- README.md      # Setup, configuration, architecture, and commands
  |-- Dockerfile
  |-- cmd/
  |-- go.mod
  |-- go.sum
  `-- internal/
  ```

- Do not create a broad shared utility package merely for the duplicated `getStringID`, `currentUserID`, and `UpdatePayload` helpers. Either keep the small duplication or centralize only BSON-specific helpers in a narrowly named persistence package.

  Suggested structure only if BSON helpers are centralized:

  ```text
  internal/db/
  |-- bson.go        # Small ObjectID conversion helpers only
  |-- indexes.go
  `-- mongo.go
  ```

## 3. Problem-Solving Approach

### Immediate Correctness and Security Fixes

- `ListProjects` and AI chat continue after writing errors, producing contradictory responses and potentially invoking Together with invalid input (`internal/project/handler.go:17-24`, `internal/ai/handler.go:17-33`).

- `UpdateCompletedItems` dereferences a nullable pointer (`internal/project/service.go:160-164`). A missing or `null` `completedItems` field can panic.

- Progress can exceed 100 because arbitrary keys count as tasks, and zero-task projects cause division by zero/NaN (`internal/project/service.go:165-187`). Build a valid task-ID set, reject unknown IDs, handle zero tasks, and store an `int`.

- The progress update path lets clients directly set `progress` and replace `completedItems`, bypassing recalculation (`internal/project/model.go:96-102`). Make computed fields server-owned.

- AI handlers accept unbounded uploads and trust the client MIME header (`internal/ai/handler.go:36-76`, `internal/ai/service.go:84-97`). Use `http.MaxBytesReader`, limit multipart parsing, sniff the MIME type, and enforce an explicit size and MIME policy, such as PNG/JPEG/WebP up to 10 MB.

- Together responses index `Choices[0]` without checking length (`internal/ai/service.go:75-81,149-158,281-287`). Chat also ignores non-2xx status codes. Close response bodies immediately after successful `Do`, bound provider error bodies, and validate response shape.

- Receipt persistence errors are silently ignored (`internal/ai/service.go:291-301`).

- Group creation discards `payload.People` (`internal/group/service.go:19-32`). Empty people lists also cannot clear existing values because of `omitempty` (`internal/group/model.go:28-30`).

### HTTP and Resilience Improvements

- Malformed JSON, duplicate users, invalid credentials, invalid IDs, missing records, and authorization failures frequently become HTTP 500. Map domain errors to 400, 401, 403, 404, and 409, and return one consistent JSON error shape.

- Raw Mongo and Together errors are exposed to clients. Log internal details and return stable public messages.

- The Together client allows 60 seconds while the server write timeout is 10 seconds (`internal/ai/service.go:34-41`, `cmd/server/main.go:104-110`). Align these deadlines and add `ReadHeaderTimeout` and `IdleTimeout`.

- Registration and refresh rotation are not failure-atomic (`internal/auth/service.go:27-66,137-154`). A token-storage failure can leave a created account or revoke the only usable refresh token. Use a transaction or an explicitly recoverable sequence.

- Add indexes for `(userId, projectId)`, group ownership, and `(userId, createdAt)` extraction queries. Extraction history should use pagination or a bounded result limit instead of loading every record into memory (`internal/ai/repository.go:47-58`).

### Testing

Only six tests exist, covering root, health, CORS, and partial configuration. Auth, projects, groups, AI, persistence, authorization, and error semantics have no tests.

Highest-value additions:

1. HTTP integration tests for routes, status codes, and response shapes.
2. JSON/BSON round-trip tests for project models.
3. Progress tests for nil payloads, zero tasks, unknown IDs, values over 100, ownership, and concurrent updates.
4. AI tests for non-2xx responses, empty choices, oversized uploads, invalid MIME types, timeouts, and persistence failures.
5. Auth tests for duplicate registration, incorrect credentials, refresh rotation, CSRF/origin enforcement, and storage failures.

Suggested test structure:

```text
cmd/server/
`-- routes_test.go             # Whole-router HTTP integration tests
internal/ai/
|-- handler_test.go            # Upload limits and request validation
`-- service_test.go            # Provider errors and persistence failures
internal/auth/
|-- middleware_test.go
`-- service_test.go            # Registration, login, and refresh behavior
internal/group/
`-- service_test.go            # Ownership and people persistence
internal/project/
|-- model_test.go              # JSON/BSON round trips
`-- service_test.go            # Progress calculation and ownership
```

`go test ./...`, `go test -race ./...`, and `go vet ./...` all pass, but the affected feature packages have no tests, so those results provide limited confidence.

## Recommended Order

1. Define project ownership, group cardinality, progress ownership, and uniqueness rules.
2. Fix panic, data-loss, upload-limit, CSRF, and ignored-error paths.
3. Separate project request payloads from persistence models before the schema grows further.
4. Establish shared validation and HTTP error handling.
5. Add domain authorization, uniqueness constraints, and query indexes.
6. Add service and HTTP integration tests before further feature work.
