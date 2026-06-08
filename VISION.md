## Level Up Webinar 101 Vision

Level Up Webinar 101 is a Go sample that sends an SMS through Twilio using
environment-provided phone numbers and Twilio credentials.

The repository is useful as a compact webinar/training sample for the Twilio Go
SDK and environment-based configuration. Project context lives in
[`README.md`](README.md).

The goal is to keep the sample small, safe, and easy to run for demos.

The current focus is:

Priority:

- Preserve the simple Twilio SMS send flow
- Keep phone numbers and Twilio credentials in environment variables
- Avoid committing account SIDs, auth tokens, API keys, or phone numbers
- Maintain security policy for the sample

Next priorities:

- Add README setup, required environment variables, and run commands
- Add input validation for missing phone numbers and credentials
- Add tests or a dry-run path that does not send real SMS
- Document expected Twilio account prerequisites

Contribution rules:

- One PR = one focused Twilio, Go, test, or documentation change.
- Do not commit secrets or real phone numbers.
- Keep the webinar sample easy to understand.
- Prefer dry-run/testing improvements before adding broader messaging features.

## Security And Privacy

SMS samples can expose phone numbers and Twilio credentials. Keep all secrets
and real recipient/sender numbers out of git and avoid logging sensitive values.

## What We Will Not Merge (For Now)

- Committed Twilio credentials or real phone numbers
- Bulk messaging features
- Logging of sensitive messaging configuration
- Broad app scaffolding that obscures the demo

This list is a roadmap guardrail, not a permanent rule.
Strong user demand and strong technical rationale can change it.
