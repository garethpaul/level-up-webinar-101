# Message Body Length Validation

status: completed

## Context

`MESSAGE_BODY` lets demo operators override the default Twilio SMS body. The
sample trimmed the value but did not bound its size before the send path, so a
misconfigured environment value could produce an oversized request while still
passing local validation.

## Objectives

- Add a sample-level maximum message body length before sending SMS.
- Keep default body behavior and exact-limit custom bodies working.
- Reject oversized `MESSAGE_BODY` values by environment variable name without
  echoing the configured body.
- Document the guard in README, vision, security, and changelog surfaces.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
