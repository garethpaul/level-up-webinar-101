# Security Policy

## Supported Versions

The supported security scope for `level-up-webinar-101` is the current default branch, `main`. Older commits, tags, branches, forks, demos, and generated artifacts are not actively supported unless the repository explicitly marks them as maintained.

Project summary: Level-Up-Webinar 101

## Reporting a Vulnerability

Please report suspected vulnerabilities through GitHub's private vulnerability reporting or by opening a draft GitHub Security Advisory for `garethpaul/level-up-webinar-101` when that option is available. If GitHub does not show a private reporting option for this repository, contact the repository owner through GitHub and avoid posting exploit details publicly until the issue can be assessed.

Do not open a public issue that includes exploit code, secrets, personal data, or detailed reproduction steps for an unpatched vulnerability.

## What to Include

Helpful reports include:

- the affected file, endpoint, permission, dependency, or workflow
- a concise impact statement explaining what an attacker could do
- reproduction steps using test data and accounts you control
- the branch, commit SHA, platform version, device, runtime, or dependency versions used
- logs, screenshots, or proof-of-concept snippets that demonstrate impact without exposing private data

## Project Security Posture

- This repository appears to be a public sample, documentation, or utility project. The active security scope is the code and documentation on the default branch.
- Review found external API integrations or credential-adjacent configuration; changes in those areas should receive security-focused review before merge.
- The sample now uses `go.mod` and `go.sum` for Twilio SDK dependency metadata. Run `make check` after Go, dependency, or documentation changes.
- The pinned Linux workflow uses read-only permissions, disables checkout
  credential persistence, selects patched Go 1.25.11, and runs formatting,
  `go vet`, module verification, injected sender tests, and builds without
  Twilio credentials, real phone numbers, outbound SMS requests, or live API
  calls.
- Required real-send values are `TO_PHONE_NUMBER`, `TWILIO_PHONE_NUMBER`, `TWILIO_ACCOUNT_SID`, and `TWILIO_AUTH_TOKEN`; reports should note whether failures expose these values.
- `DRY_RUN=1` should validate non-secret E.164-style phone-number configuration without sending SMS or printing phone numbers, account SIDs, or auth tokens.
- Ambiguous `DRY_RUN` values should fail closed by naming `DRY_RUN` without echoing the configured value.
- Phone number validation errors should name environment variables rather than echoing configured phone number values.
- Matching sender/recipient phone number errors should name `TO_PHONE_NUMBER` and `TWILIO_PHONE_NUMBER` without echoing the shared value.
- Keep the explicit 10-second Twilio request timeout applied before live sends
  so dependency defaults cannot introduce an unbounded request.
- Real-send Account SID validation errors should name `TWILIO_ACCOUNT_SID` rather than echoing configured values.
- Real-send Auth Token validation errors should name `TWILIO_AUTH_TOKEN` rather than echoing configured values.
- All-zero Twilio Account SID and Auth Token placeholder-shaped credentials should be rejected by name rather than echoing configured values.
- Message body validation errors should name `MESSAGE_BODY` rather than echoing oversized or invalid UTF-8 content.
- Keep `.env` files, local shell exports, real phone numbers, account SIDs, auth tokens, API keys, and webhook secrets out of git.


## Dependency and Supply Chain Security

The canonical `make check` gate runs
`golang.org/x/vuln/cmd/govulncheck@v1.3.0` against all source packages on the
pinned Go 1.25.11 toolchain. Hosted validation must fail on reachable known
vulnerabilities rather than suppressing or converting findings to a
success-only output format. The scanner queries the public Go vulnerability
database with module paths. It does not upload repository source code.

Dependency updates should come from trusted package managers and should keep `go.mod` and `go.sum` in sync. Do not commit credentials, private keys, tokens, generated secrets, or machine-local configuration. If a vulnerability depends on a compromised package, typosquatting risk, insecure transitive dependency, or unsafe build step, include the package name, affected version, and the path through which it is used.

## Safe Research Guidelines

Good-faith research is welcome when it stays within these boundaries:

- use only accounts, devices, data, and infrastructure that you own or have explicit permission to test
- avoid destructive actions, persistence, spam, phishing, social engineering, or denial-of-service testing
- minimize access to personal data and stop testing immediately if private data is exposed
- do not exfiltrate secrets or third-party data; report the minimum evidence needed to verify impact
- keep vulnerability details confidential until the maintainer has assessed the report

## Maintainer Response

The maintainer will review complete reports as availability allows, prioritize issues by exploitability and impact, and coordinate a fix or mitigation when the affected code is still maintained. For sample, archived, or educational repositories, the likely remediation may be documentation, dependency updates, or clearly marking unsupported code rather than a production-style patch release.
