# Message Body UTF-8 Validation

status: completed

## Context

`MESSAGE_BODY` is supplied through the environment and becomes the body of the
Twilio SMS request. The sample already trims and bounds the message body, but
invalid UTF-8 input should fail local validation before reaching the Twilio
client.

## Completed Scope

- Rejected invalid UTF-8 `MESSAGE_BODY` values before real sends.
- Kept message body errors limited to the `MESSAGE_BODY` environment variable
  name without echoing configured content.
- Added unit coverage for invalid UTF-8 body values.
- Extended the scripted baseline and docs so the encoding guard stays visible.

## Verification

- `make lint`
- `make test`
- `make build`
- `make check`
- `go test ./...`
- `git diff --check`
