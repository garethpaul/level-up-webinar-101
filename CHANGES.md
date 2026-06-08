# Changes

## 2026-06-08

- Added Go module metadata and checksum locking for the Twilio SDK dependency.
- Added environment validation for recipient, sender, account SID, and auth token configuration.
- Added E.164-style phone number validation without echoing configured phone number values in errors.
- Added `DRY_RUN` support for webinar demos that should validate configuration without sending SMS.
- Added an injectable run path so send success and Twilio failures can be tested without contacting Twilio.
- Added Go unit tests and a `make check` verification target that covers formatting, module verification, tests, and builds.
- Updated generated project context for the Go module, tests, and dry-run baseline.
