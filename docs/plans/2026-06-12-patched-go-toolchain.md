# Patched Go Toolchain

status: completed

## Context

`govulncheck` found reachable standard-library vulnerabilities in the Twilio
HTTP and TLS request path when validation used Go 1.25.3. The official Go
release feed identifies Go 1.25.11 as the patched release for the complete
reported set.

## Changes

- Raised the module Go version to 1.25.11.
- Pinned hosted validation to exact Go 1.25.11.
- Made the baseline reject a different selected Go toolchain.
- Kept the latest Twilio Go dependency and existing request timeout.

## Verification

- `go test -race ./...`
- `go vet ./...`
- `go build ./...`
- `govulncheck ./...`
- `make check`
- hosted GitHub Actions at the final commit
