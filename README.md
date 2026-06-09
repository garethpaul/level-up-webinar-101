# level-up-webinar-101

<!-- README-OVERVIEW-IMAGE -->
![Project overview](docs/readme-overview.svg)

## Overview

`garethpaul/level-up-webinar-101` is a Go project. Level-Up-Webinar 101

This README is based on the checked-in source, manifests, scripts, and repository metadata on the `main` branch. The project language mix found during review was: Go (1).

## Repository Contents

- `.gitignore` - local secret, log, coverage, and binary ignores
- `CHANGES.md` - recent maintenance changes
- `Makefile` - local verification entry points
- `README.md` - project overview and local usage notes
- `go.mod` - Go module and Twilio SDK dependency metadata
- `go.sum` - Go dependency checksums
- `main.go` - Twilio SMS sample entry point
- `main_test.go` - configuration and dry-run unit tests
- `SECURITY.md` - security reporting and disclosure guidance
- `VISION.md` - project direction and maintenance guardrails

Additional scan context:

- Source directories: no top-level source directories detected
- Dependency and build manifests: go.mod, go.sum
- Entry points or build surfaces: `make lint`, `make test`, `make build`, `make check`, main.go
- Test-looking files: main_test.go

## Getting Started

### Prerequisites

- Git
- Go 1.24 or newer

### Setup

```bash
git clone https://github.com/garethpaul/level-up-webinar-101.git
cd level-up-webinar-101
make lint
make test
make build
make check
```

The setup commands above are derived from repository files. Legacy mobile, Python, or JavaScript samples may require older SDKs or package versions than a modern workstation uses by default.

## Running or Using the Project

Set the required environment variables before sending a real SMS. Populate
`TWILIO_ACCOUNT_SID` and `TWILIO_AUTH_TOKEN` from your local secret store or
shell history-safe prompt. Phone numbers should use E.164 format, such as `+15558675310`.
The recipient and Twilio sender phone numbers must be different.
`TWILIO_ACCOUNT_SID` should use Twilio's `AC`-prefixed Account SID format.
`TWILIO_AUTH_TOKEN` should use Twilio's 32-character hexadecimal Auth Token format.

```bash
export TO_PHONE_NUMBER="+15558675310"
export TWILIO_PHONE_NUMBER="+15558675309"
export TWILIO_ACCOUNT_SID
export TWILIO_AUTH_TOKEN
go run .
```

Optional environment variables:

- `MESSAGE_BODY` overrides the default `Hello from Golang!` body.
- `MESSAGE_BODY` is trimmed and must be no more than 1600 characters.
- `DRY_RUN=1` validates phone-number configuration without requiring Twilio credentials or sending SMS. Accepted true values are `1`, `true`, `t`, `yes`, `y`, and `on`; accepted false values are empty, `0`, `false`, `f`, `no`, `n`, and `off`. Any other `DRY_RUN` value is rejected instead of being treated as a real send.

For a no-send setup check:

```bash
TO_PHONE_NUMBER="+15558675310" TWILIO_PHONE_NUMBER="+15558675309" DRY_RUN=1 go run .
```

## Testing and Verification

- `make lint` verifies Go formatting with `gofmt`.
- `make test` verifies Go module checksums and runs the unit tests.
- `make build` compiles the local package.
- `make check` runs `make lint`, `make test`, and `make build`.
- `go test ./...` covers missing environment variables, strict dry-run value parsing, dry-run behavior, E.164-style phone number validation, matching sender/recipient rejection, Account SID validation, Auth Token validation, custom message body handling, message body length validation, whitespace trimming, sender success, and sender error wrapping without contacting Twilio.

When the required SDK or runtime is unavailable, use static checks and source review first, then verify on a machine that has the matching platform toolchain.

## Configuration and Secrets

- Required for real sends: `TO_PHONE_NUMBER`, `TWILIO_PHONE_NUMBER`, `TWILIO_ACCOUNT_SID`, and `TWILIO_AUTH_TOKEN`.
- Required for `DRY_RUN=1`: `TO_PHONE_NUMBER` and `TWILIO_PHONE_NUMBER`.
- Phone number values must be E.164-style strings beginning with `+` followed by digits.
- `TO_PHONE_NUMBER` and `TWILIO_PHONE_NUMBER` must not be the same value.
- Real sends validate that `TWILIO_ACCOUNT_SID` is an `AC`-prefixed Twilio Account SID.
- Real sends validate that `TWILIO_AUTH_TOKEN` is a 32-character hexadecimal Twilio Auth Token.
- All-zero Twilio Account SID and Auth Token placeholder-shaped credentials are rejected by name.
- `MESSAGE_BODY` values longer than 1600 characters are rejected by name without echoing the body.
- Ambiguous `DRY_RUN` values are rejected by name without echoing the configured value.
- Keep Twilio credentials and real phone numbers in local environment variables or secret stores only.

## Security and Privacy Notes

- Review changes touching external API calls or credential-adjacent configuration; examples from the scan include main.go.
- Do not log Twilio auth tokens, account SIDs, or real phone numbers. The dry-run path confirms configuration without printing sensitive values.

## Maintenance Notes

- Run `make lint`, `make test`, `make build`, and `make check` before pushing Go, dependency, README, or security-policy changes.
- See `SECURITY.md` for vulnerability reporting and safe research guidance.
- See `VISION.md` for project direction and contribution guardrails.
- See `docs/plans/2026-06-09-make-gate-targets.md` for the local gate target guardrail.

## Contributing

Keep changes small and tied to the project that is already present in this repository. For code changes, document the toolchain used, avoid committing generated dependency directories or local configuration, and update this README when setup or verification steps change.
