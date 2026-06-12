# Explicit Twilio Timeout

status: completed

## Problem

The live send path relies on the Twilio SDK's current default HTTP timeout.
That upstream default is not an explicit application contract and can change
across dependency upgrades.

## Scope

- Define a 10-second Twilio request timeout in application code.
- Apply it to the REST client before creating a message.
- Test the configured underlying HTTP transport without making a network call.
- Extend the scripted baseline to guard the timeout constant and application.
- Preserve dry-run behavior and existing sender injection tests.

## Verification

- `go test ./...`
- `go vet ./...`
- `make check`
- mutation check that removes timeout application
- `git diff --check`
