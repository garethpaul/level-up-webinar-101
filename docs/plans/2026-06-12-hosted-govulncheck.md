# Hosted govulncheck enforcement

status: completed

## Goal

Make the official Go vulnerability scanner part of the canonical local and
hosted `make check` gate so the repository continuously verifies that its
Twilio request path does not call known vulnerable standard-library or module
symbols.

## Context

The repository already requires Go 1.25.11 because an earlier scan found
reachable standard-library vulnerabilities with Go 1.25.3. The patched
toolchain and current Twilio dependency pass a local scan, but GitHub Actions
currently runs formatting, vet, tests, build, and static policy only. Official
Go module metadata identifies `golang.org/x/vuln` v1.3.0 as the current tagged
scanner release.

## Implementation Units

1. Add a `vuln` Make target that runs the pinned
   `golang.org/x/vuln/cmd/govulncheck@v1.3.0` source scan and include it in the
   canonical `check` dependency chain.
2. Extend the static baseline to require the exact pinned scanner command,
   reject unversioned or bypassed scans, and register this completed plan.
3. Update README, security guidance, vision, and changes to describe hosted
   reachable-vulnerability enforcement and the scanner's vulnerability
   database network/privacy boundary.

## Verification

- Run `make lint`, `make test`, `make build`, `make vuln`, and `make check` with
  Go 1.25.11.
- Run `go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...` and require a
  zero-finding exit status.
- Run the full Make gate from an external working directory.
- Execute hostile mutations for scanner removal, floating `@latest` selection,
  version drift, check-chain bypass, and incomplete plan status.
- Require exact-head push and `pull_request` GitHub Actions runs to succeed.

## Boundaries

- Do not add the scanner to `go.mod` or `go.sum`; use the version-qualified Go
  tool invocation so application dependencies remain unchanged.
- Do not suppress findings or use JSON output, which exits successfully even
  when vulnerabilities are present.
- Do not add Twilio credentials or make live Twilio API calls.
- The scanner queries the public Go vulnerability database and sends module
  paths, not repository source code; offline execution is outside this unit.

## Assumptions

- Hosted validation retains outbound access to the official Go module proxy
  and vulnerability database.
- The pinned scanner remains compatible with the repository's exact Go 1.25.11
  floor.

## Completion Evidence

- `make lint`, `make test`, and `make build` passed with Go 1.25.11.
- `make vuln` and the exact pinned `go run` command completed with `No
  vulnerabilities found.`
- Repository-local and external-working-directory `make check` are the
  authoritative combined gate for the completed implementation.
- A clean temporary module and build cache completed the full gate, including
  Go 1.25.11 toolchain resolution, application dependencies, govulncheck
  v1.3.0 installation, and a zero-finding scan, in under the hosted timeout.
- Seven hostile mutations covering check-chain bypass, floating or drifted
  scanner versions, success-masking JSON output, incomplete plan status, and a
  nonportable Make command were rejected by the baseline checker, including a
  commented approved-recipe spoof paired with an executable bypass.
