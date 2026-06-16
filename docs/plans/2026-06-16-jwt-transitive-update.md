# Update Transitive JWT Dependency

status: planned

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

## Planned Verification

- Prove the runtime dependency path with `go mod why -m` and `go list -deps`.
- Run formatting, vet, module verification, tests, race tests, coverage, build,
  pinned `govulncheck`, and the static baseline.
- Run `make check` from the repository root and through the absolute Makefile
  path from an external directory.
- Reject isolated mutations for JWT rollback, missing static enforcement,
  incomplete plan status, erased dependency-path evidence, and unrelated
  Twilio or Go version drift.
