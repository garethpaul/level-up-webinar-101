# Update Transitive JWT Dependency

status: completed

## Context

Twilio Go v1.30.9 selects `github.com/golang-jwt/jwt/v5` v5.2.2 through its
runtime dependency graph. Upstream v5.3.1 is the current compatible minor
release and includes parser and signature-handling corrections, so retaining
the older selected version is unnecessary dependency risk.

## Requirements

- Select `github.com/golang-jwt/jwt/v5` v5.3.1 without changing Twilio Go
  v1.30.9 or the Go 1.25.11 module floor.
- Preserve application behavior, public interfaces, environment variables,
  validation, redacted errors, dry-run isolation, and the bounded Twilio
  client.
- Add a static contract that rejects JWT version rollback and incomplete plan
  evidence.
- Verify the selected module graph, checksums, tests, race behavior, coverage,
  build, vet, vulnerability scan, and both supported Make invocation locations.

## Scope Boundaries

- Do not update unrelated application modules or hosted action pins.
- Do not add live Twilio credentials or send a real SMS.
- Do not suppress vulnerability findings or weaken an existing gate.
- Do not merge or close any existing pull request.

## Work Completed

- Updated the selected runtime dependency from
  `github.com/golang-jwt/jwt/v5 v5.2.2` to
  `github.com/golang-jwt/jwt/v5 v5.3.1` while retaining Twilio Go v1.30.9 and
  Go 1.25.11.
- Added a focused POSIX dependency contract and fixture-based mutation suite to
  the canonical static baseline.
- Documented the selected JWT version and its Twilio runtime path in the
  repository's maintenance and security guidance.

## Verification Completed

- `go mod why -m github.com/golang-jwt/jwt/v5` proved the application reaches
  JWT through `github.com/twilio/twilio-go`; `go list -deps ./...` included the
  JWT package, and `go list -m all` selected JWT v5.3.1 with Twilio Go v1.30.9.
- `go mod verify`, `go test ./...`, `go test -race ./...`, `go vet ./...`, and
  `go build ./...` passed. Coverage remained 86.2% of statements.
- Pinned `govulncheck@v1.3.0` reported `No vulnerabilities found.`
- The authoritative `make check` gates cover both root and external-directory
  invocation through the absolute Makefile path.
- The focused fixture suite rejected five isolated hostile mutations: JWT
  rollback, duplicate selection, incomplete plan status, erased dependency-
  path evidence, and missing security guidance.
