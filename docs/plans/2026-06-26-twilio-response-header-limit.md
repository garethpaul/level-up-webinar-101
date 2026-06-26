# Twilio Response Header Limit

status: completed

## Problem

Twilio response bodies were bounded to 256 KiB, but real sends still used Go's
shared default transport with its much larger response-header allowance. A
provider or intermediary could consume substantially more memory in headers
before the body guard was reached.

## Requirements

1. Bound real-send response headers independently of response bodies.
2. Clone the default transport instead of mutating process-global HTTP state.
3. Preserve the 10-second request timeout, redirect refusal, body cap, official
   Twilio host, error redaction, and no-retry behavior.
4. Add focused regression, baseline contracts, hostile mutations, and full Go
   verification without live Twilio credentials.

## Work Completed

- Added a 64 KiB `MaxResponseHeaderBytes` limit on a cloned default transport.
- Routed real Twilio sends through the isolated transport while preserving
  injected round trippers for provider tests.
- Added a regression proving the limit and global transport isolation.
- Added five isolated hostile mutations covering the limit, assignment, real
  send routing, regression test, and completed plan.

## Verification Completed

- All 28 focused tests and five isolated hostile mutations passed.
- `gofmt`, `go vet ./...`, `go mod verify`, `go test ./...`, and `go build ./...`
  passed with the module-pinned Go 1.26.4 toolchain.
- `make check` passed from the repository root and an external working directory.
- Pinned `govulncheck@v1.3.0`, `git diff --check`, strict Git validation,
  generated-artifact checks, and secret/conflict scans passed without Twilio
  credentials or a live SMS request.
