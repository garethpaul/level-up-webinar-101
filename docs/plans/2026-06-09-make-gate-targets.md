# Make Gate Targets

status: completed

## Context

The repository had one `make check` target that handled formatting, module
verification, tests, and build output. The shared maintenance workflow expects
`make lint`, `make test`, `make build`, and `make check` to be available before
a change is pushed.

## Completed Scope

- Added `make lint` for the `gofmt` check.
- Added `make test` for module verification and unit tests.
- Added `make build` for local compilation.
- Kept `make check` as the full gate by chaining lint, test, and build.
- Updated README, VISION, and CHANGES so the gate contract is visible.

## Verification

- `make lint`
- `make test`
- `make build`
- `make check`
- `git diff --check`
