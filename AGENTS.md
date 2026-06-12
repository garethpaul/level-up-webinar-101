# AGENTS.md

## Repository purpose

`garethpaul/level-up-webinar-101` is a Go project. Level-Up-Webinar 101

## Project structure

- `Makefile` - repository verification targets
- `scripts` - baseline checks and helper scripts
- `docs` - plans, notes, and generated README assets
- `go.mod` - Go module definition

## Development commands

- Install dependencies: `go mod download`
- Full baseline: `make check`
- Lint/static checks: `make lint`
- Tests: `make test`
- Build: `make build`
- Go test all packages: `go test ./...`
- Go vet all packages: `go vet ./...`
- Go build all packages: `go build ./...`
- Use Go 1.25.11 or newer; hosted validation pins 1.25.11 while local checks
  reject vulnerable earlier Go releases.
- If a command above skips because a platform toolchain is missing, verify on a machine with that SDK before claiming platform behavior is tested.

## Coding conventions

- Language mix noted in the README: Go (1).
- Keep imports compatible with module path `github.com/garethpaul/level-up-webinar-101`.
- Run gofmt on changed Go files and keep table-driven tests close to the package under change.

## Testing guidance

- Test-related files detected: `main_test.go`
- Start with the narrowest relevant test or Make target, then run `make check` before handing off if the change is not documentation-only.
- Keep README verification notes in sync when commands, fixtures, or supported toolchains change.

## PR / change guidance

- Keep diffs focused on the requested repository and avoid unrelated modernization or formatting churn.
- Preserve public APIs, sample behavior, file formats, and documented environment variables unless the task explicitly changes them.
- Update tests, README notes, or docs/plans when behavior, security posture, or validation commands change.
- Call out skipped platform validation, legacy toolchain assumptions, and any risky files touched in the final summary.

## Safety and gotchas

- Required for real sends: `TO_PHONE_NUMBER`, `TWILIO_PHONE_NUMBER`, `TWILIO_ACCOUNT_SID`, and `TWILIO_AUTH_TOKEN`.
- Required for `DRY_RUN=1`: `TO_PHONE_NUMBER` and `TWILIO_PHONE_NUMBER`.
- Phone number values must be E.164-style strings beginning with `+` followed by digits.
- `TO_PHONE_NUMBER` and `TWILIO_PHONE_NUMBER` must not be the same value.
- Real sends validate that `TWILIO_ACCOUNT_SID` is an `AC`-prefixed Twilio Account SID.
- Real sends validate that `TWILIO_AUTH_TOKEN` is a 32-character hexadecimal Twilio Auth Token.

## Agent workflow

1. Inspect the README, Makefile, manifests, and the files directly related to the request.
2. Make the smallest source or docs change that satisfies the task; avoid generated, vendored, or local-environment files unless required.
3. Run the narrowest useful validation first, then `make check` or the documented package/platform gate when available.
4. If a required SDK, service credential, or external runtime is unavailable, record the skipped command and why.
5. Summarize changed files, commands run, and remaining risks or follow-up validation.
