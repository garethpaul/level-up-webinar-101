# Patched Go Toolchain

status: completed

## Context

`govulncheck` found reachable standard-library vulnerabilities in the Twilio
HTTP and TLS request path when validation used older Go 1.25 patch releases.
The floor was raised to Go 1.25.11, then to Go 1.25.12 after GO-2026-5856
(`crypto/tls` Encrypted Client Hello privacy leak) appeared on the live send
path under Go 1.25.11. The official Go release feed identifies Go 1.25.12 as
the patched release for the current reported set.

## Changes

- Raised the module Go version to 1.25.12.
- Pinned hosted validation to exact Go 1.25.12.
- Made the baseline reject a different selected Go toolchain below 1.25.12.
- Kept the latest Twilio Go dependency and existing request timeout.

## Verification

- `go test -race ./...`
- `go vet ./...`
- `go build ./...`
- `govulncheck ./...`
- `make check`
- hosted GitHub Actions at the final commit
