## Level Up Webinar 101 Vision

This document explains the current state and direction of the project.
Project overview and developer docs: [`README.md`](README.md)

Level Up Webinar 101 is a Go sample that sends an SMS through Twilio using
environment-provided phone numbers and Twilio credentials.

The repository is useful as a compact webinar/training sample for the Twilio Go
SDK and environment-based configuration. Project context lives in
[`README.md`](README.md).

The goal is to keep the sample small, safe, and easy to run for demos.

Current baseline: the sample has Go module metadata, configuration validation,
`DRY_RUN` support for no-send demos, local Make gates for formatting, module
verification, unit tests, a local build through `make lint`, `make test`,
`make build`, and `make check`, plus `scripts/check-baseline.sh` for
repository metadata and hygiene checks.

The current focus is:

Priority:

- Preserve the simple Twilio SMS send flow
- Keep live Twilio requests bounded by an explicit 10-second timeout
- Keep phone numbers and Twilio credentials in environment variables
- Validate phone number shape before dry runs or real sends
- Reject matching sender and recipient phone numbers
- Validate Account SID shape before real sends
- Validate Auth Token shape before real sends
- Reject all-zero Twilio credential placeholders before real sends
- Validate message body encoding and length before demo sends, including invalid UTF-8 rejection
- Keep `DRY_RUN` available for webinar setup checks that should not send SMS
- Reject ambiguous `DRY_RUN` values instead of treating typos as real sends
- Keep `make lint`, `make test`, `make build`, and `make check` green before
  pushing changes
- Keep `go vet ./...` and pinned, credential-free, read-only Go 1.25.11 hosted
  validation in the canonical gate without contacting Twilio
- Keep `scripts/check-baseline.sh` green as repository metadata and local
  hygiene evolve
- Avoid committing account SIDs, auth tokens, API keys, or phone numbers
- Maintain security policy for the sample

Next priorities:

- Document expected Twilio account prerequisites
- Add a fake HTTP client if Twilio request payloads need deeper unit coverage
- Consider examples for Twilio trial-account verification limits

Contribution rules:

- One PR = one focused Twilio, Go, test, or documentation change.
- Do not commit secrets or real phone numbers.
- Run `make lint`, `make test`, `make build`, and `make check` before pushing.
- Update `scripts/check-baseline.sh` when required files or verification docs
  intentionally change.
- Keep validation errors limited to environment variable names.
- Keep the webinar sample easy to understand.
- Prefer dry-run/testing improvements before adding broader messaging features.

## Security And Privacy

Canonical security policy and reporting:

- [`SECURITY.md`](SECURITY.md)

SMS samples can expose phone numbers and Twilio credentials. Keep all secrets
and real recipient/sender numbers out of git and avoid logging sensitive values.

## What We Will Not Merge (For Now)

- Committed Twilio credentials or real phone numbers
- Bulk messaging features
- Logging of sensitive messaging configuration
- Broad app scaffolding that obscures the demo

This list is a roadmap guardrail, not a permanent rule.
Strong user demand and strong technical rationale can change it.
