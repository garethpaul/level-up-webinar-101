# Hosted Go Validation

status: completed

## Context

The repository has a locked Go module, focused configuration and sender tests,
and a canonical local gate, but no hosted validation. The lint target currently
checks formatting without running Go's static analyzer.

## Priorities

1. Add `go vet ./...` to the canonical lint gate.
2. Run `make check` on pinned Go 1.24 tooling in hosted Linux CI.
3. Enforce a pinned, read-only, bounded workflow from the baseline script.
4. Keep Twilio credentials, real phone numbers, and outbound SMS requests out
   of hosted validation.

## Implementation Units

### Gates And Workflow

Files:

- `Makefile`
- `.github/workflows/check.yml`
- `scripts/check-baseline.sh`

Add Go vet to lint. Add push, pull-request, and manual workflow triggers;
read-only permissions; concurrency cancellation; a bounded `ubuntu-24.04` job;
commit-pinned checkout and Go setup; Go `1.24.x`; and `make check`. Require that
contract from the baseline script.

### Documentation

Files:

- `README.md`
- `VISION.md`
- `SECURITY.md`
- `CHANGES.md`
- `docs/plans/2026-06-10-hosted-go-validation.md`

Document hosted tests and static analysis without implying a live Twilio
integration test.

## Verification

- `make lint`
- `make test`
- `make build`
- `make check`
- workflow YAML parse
- `git diff --check`
- successful hosted Linux `Check` workflow for the pushed commit

## Boundaries

- Do not configure Twilio credentials or real phone numbers in CI.
- Do not invoke the production sender or contact the Twilio API.
