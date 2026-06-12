#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
README="$ROOT_DIR/README.md"
MAKEFILE="$ROOT_DIR/Makefile"
GITIGNORE="$ROOT_DIR/.gitignore"
DOCS_PLANS="$ROOT_DIR/docs/plans"

require_file() {
  path=$1
  if [ ! -f "$ROOT_DIR/$path" ]; then
    printf '%s\n' "Required file is missing: $path" >&2
    exit 1
  fi
}

for path in \
  ".gitignore" \
  ".github/workflows/check.yml" \
  "AGENTS.md" \
  "CHANGES.md" \
  "Makefile" \
  "README.md" \
  "SECURITY.md" \
  "VISION.md" \
  "go.mod" \
  "go.sum" \
  "main.go" \
  "main_test.go" \
  "docs/plans/2026-06-08-twilio-sms-baseline.md" \
  "docs/plans/2026-06-09-make-gate-targets.md" \
  "docs/plans/2026-06-09-scripted-baseline-check.md" \
  "docs/plans/2026-06-10-message-body-utf8-validation.md" \
  "docs/plans/2026-06-10-explicit-twilio-timeout.md" \
  "docs/plans/2026-06-10-hosted-go-validation.md" \
  "docs/plans/2026-06-12-redacted-twilio-send-errors.md" \
  "docs/plans/2026-06-12-patched-go-toolchain.md" \
  "docs/plans/2026-06-12-hosted-govulncheck.md" \
  "scripts/check-baseline.sh"; do
  require_file "$path"
done

if ! grep -Fq "utf8.ValidString(config.MessageBody)" "$ROOT_DIR/main.go"; then
  printf '%s\n' "main.go must reject invalid UTF-8 MESSAGE_BODY values." >&2
  exit 1
fi

if ! grep -Fq "TestLoadConfigRejectsInvalidUTF8MessageBody" "$ROOT_DIR/main_test.go"; then
  printf '%s\n' "main_test.go must cover invalid UTF-8 MESSAGE_BODY values." >&2
  exit 1
fi

if ! grep -Fq "const twilioRequestTimeout = 10 * time.Second" "$ROOT_DIR/main.go" || \
   ! grep -Fq "client.SetTimeout(twilioRequestTimeout)" "$ROOT_DIR/main.go" || \
   ! grep -Fq "configureTwilioClient(client)" "$ROOT_DIR/main.go"; then
  printf '%s\n' "main.go must apply the explicit 10-second Twilio request timeout." >&2
  exit 1
fi

if ! grep -Fq "TestConfigureTwilioClientSetsRequestTimeout" "$ROOT_DIR/main_test.go"; then
  printf '%s\n' "main_test.go must cover the configured Twilio request timeout." >&2
  exit 1
fi

if ! grep -Fq 'return "send SMS: request failed"' "$ROOT_DIR/main.go" || \
   ! grep -Fq "func (err smsSendError) Unwrap() error" "$ROOT_DIR/main.go" || \
   ! grep -Fq "return smsSendError{cause: err}" "$ROOT_DIR/main.go"; then
  printf '%s\n' "main.go must redact Twilio send errors while preserving unwrapping." >&2
  exit 1
fi

if ! grep -Fq "TestRunRedactsAndWrapsSenderError" "$ROOT_DIR/main_test.go" || \
   ! grep -Fq "errors.Is(err, cause)" "$ROOT_DIR/main_test.go"; then
  printf '%s\n' "main_test.go must cover redacted, unwrap-capable Twilio send errors." >&2
  exit 1
fi

if ! grep -Fq "scripts/check-baseline.sh" "$MAKEFILE"; then
  printf '%s\n' "Makefile must run scripts/check-baseline.sh from make check." >&2
  exit 1
fi

for make_target in \
  'ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))' \
  'check: lint test build vuln' \
  'lint:' \
  'test:' \
  'build:' \
  'vuln:' \
  'fmt:'; do
  if ! grep -Fxq "$make_target" "$MAKEFILE"; then
    printf '%s\n' "Makefile must preserve exact target $make_target" >&2
    exit 1
  fi
done

tab=$(printf '\t')
for make_recipe in \
  'cd "$(ROOT)" && test -z "$$(gofmt -l *.go)"' \
  'cd "$(ROOT)" && go vet ./...' \
  'cd "$(ROOT)" && go mod verify' \
  'cd "$(ROOT)" && go test ./...' \
  'cd "$(ROOT)" && go build ./...' \
  'cd "$(ROOT)" && gofmt -w *.go' \
  'cd "$(ROOT)" && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...' \
  'cd "$(ROOT)" && ./scripts/check-baseline.sh'; do
  if ! grep -Fxq "${tab}${make_recipe}" "$MAKEFILE"; then
    printf '%s\n' "Makefile must preserve executable recipe $make_recipe" >&2
    exit 1
  fi
done

if grep -Eq 'govulncheck@(latest|master|main)|govulncheck[[:space:]]+\./\.\.\.' "$MAKEFILE" || \
   [ "$(grep -Fc 'golang.org/x/vuln/cmd/govulncheck@v1.3.0' "$MAKEFILE")" -ne 1 ]; then
  printf '%s\n' "Makefile must run exactly one pinned govulncheck v1.3.0 source scan." >&2
  exit 1
fi

WORKFLOW="$ROOT_DIR/.github/workflows/check.yml"
for workflow_value in \
  "permissions:" \
  "contents: read" \
  "cancel-in-progress: true" \
  "runs-on: ubuntu-24.04" \
  "timeout-minutes: 10" \
  "actions/checkout@df4cb1c069e1874edd31b4311f1884172cec0e10" \
  "actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c" \
  "go-version: \"1.25.11\"" \
  "run: make check"; do
  if ! grep -Fq "$workflow_value" "$WORKFLOW"; then
    printf '%s\n' "Check workflow must keep $workflow_value" >&2
    exit 1
  fi
done

if [ "$(grep -Ec '^[[:space:]]*permissions:' "$WORKFLOW")" -ne 1 ] || \
   grep -Eq '^[[:space:]]+[A-Za-z-]+:[[:space:]]+write[[:space:]]*$' "$WORKFLOW" || \
   [ "$(grep -Ec '^[[:space:]]*persist-credentials:' "$WORKFLOW")" -ne 1 ] || \
   ! grep -Fq 'persist-credentials: false' "$WORKFLOW"; then
  printf '%s\n' "Check workflow must keep one read-only permission block and credential-free checkout." >&2
  exit 1
fi

workflow_actions=$(sed -n 's/^[[:space:]]*-\{0,1\}[[:space:]]*uses:[[:space:]]*\([^[:space:]#]*\).*$/\1/p' "$WORKFLOW")
expected_actions=$(printf '%s\n' \
  'actions/checkout@df4cb1c069e1874edd31b4311f1884172cec0e10' \
  'actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c')
if [ "$workflow_actions" != "$expected_actions" ]; then
  printf '%s\n' "Check workflow must use only the expected pinned checkout and setup-go actions." >&2
  exit 1
fi

AGENTS="$ROOT_DIR/AGENTS.md"
for guidance in \
  "go test ./..." \
  "go vet ./..." \
  "go build ./..." \
  "TO_PHONE_NUMBER" \
  "TWILIO_PHONE_NUMBER" \
  "TWILIO_ACCOUNT_SID" \
  "TWILIO_AUTH_TOKEN"; do
  if ! grep -Fq "$guidance" "$AGENTS"; then
    printf '%s\n' "AGENTS.md must preserve $guidance guidance." >&2
    exit 1
  fi
done

for documented in \
  "TO_PHONE_NUMBER" \
  "TWILIO_PHONE_NUMBER" \
  "TWILIO_ACCOUNT_SID" \
  "TWILIO_AUTH_TOKEN" \
  "DRY_RUN=1" \
  "MESSAGE_BODY" \
  "10-second Twilio request timeout" \
  "invalid UTF-8" \
  "redacted Twilio send errors" \
  "go test ./..." \
  "go vet ./..." \
  "make vuln" \
  "govulncheck" \
  "make check" \
  "scripts/check-baseline.sh"; do
  if ! grep -Fq "$documented" "$README"; then
    printf '%s\n' "README must document $documented." >&2
    exit 1
  fi
done

for doc in "SECURITY.md" "VISION.md" "CHANGES.md"; do
  if ! grep -Fq "invalid UTF-8" "$ROOT_DIR/$doc"; then
    printf '%s\n' "$doc must document invalid UTF-8 MESSAGE_BODY validation." >&2
    exit 1
  fi
  if ! grep -Fq "govulncheck" "$ROOT_DIR/$doc"; then
    printf '%s\n' "$doc must document canonical govulncheck enforcement." >&2
    exit 1
  fi
done

for doc in "SECURITY.md" "VISION.md" "CHANGES.md"; do
  if ! grep -Fq "redacted Twilio send errors" "$ROOT_DIR/$doc"; then
    printf '%s\n' "$doc must document redacted Twilio send errors." >&2
    exit 1
  fi
done

if ! grep -Fq "public Go vulnerability database" "$README" || \
   ! grep -Fq "does not upload repository source code" "$README" || \
   ! grep -Fq "does not upload repository source code" "$ROOT_DIR/SECURITY.md"; then
  printf '%s\n' "README and SECURITY must document the govulncheck database privacy boundary." >&2
  exit 1
fi

if ! grep -Fq "credential-free" "$README" || \
   ! grep -Fq "credential-free" "$ROOT_DIR/VISION.md" || \
   ! grep -Fq "credential-free" "$ROOT_DIR/CHANGES.md" || \
   ! grep -Fq "disables checkout" "$ROOT_DIR/SECURITY.md"; then
  printf '%s\n' "Docs must preserve the credential-free hosted validation boundary." >&2
  exit 1
fi

for module_line in \
  "module github.com/garethpaul/level-up-webinar-101" \
  "go 1.25.11" \
  "github.com/twilio/twilio-go v1.30.9"; do
  if ! grep -Fq "$module_line" "$ROOT_DIR/go.mod"; then
    printf '%s\n' "go.mod must keep module baseline: $module_line" >&2
    exit 1
  fi
done

selected_go_version=$(go env GOVERSION | sed 's/^go//')
if ! printf '%s\n' "$selected_go_version" | awk -F. '
  $1 > 1 || ($1 == 1 && ($2 > 25 || ($2 == 25 && $3 >= 11))) { valid = 1 }
  END { exit valid ? 0 : 1 }
'; then
  printf '%s\n' "Verification requires patched Go 1.25.11 or newer." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$ROOT_DIR/docs/plans/2026-06-12-patched-go-toolchain.md" || \
   ! grep -Fq "govulncheck" "$ROOT_DIR/docs/plans/2026-06-12-patched-go-toolchain.md"; then
  printf '%s\n' "Patched Go toolchain plan must record completed vulnerability validation." >&2
  exit 1
fi

GOVULNCHECK_PLAN="$ROOT_DIR/docs/plans/2026-06-12-hosted-govulncheck.md"
for plan_contract in \
  "status: completed" \
  "govulncheck@v1.3.0" \
  "make vuln" \
  "make check" \
  "zero-finding"; do
  if ! grep -Fq "$plan_contract" "$GOVULNCHECK_PLAN"; then
    printf '%s\n' "Hosted govulncheck plan must record $plan_contract." >&2
    exit 1
  fi
done

for ignored in ".env" ".env.*" ".vscode/" ".idea/" "*.iml" "*.log" "coverage.out" "level-up-webinar-101"; do
  if ! grep -Fq "$ignored" "$GITIGNORE"; then
    printf '%s\n' ".gitignore must include $ignored" >&2
    exit 1
  fi
done

tracked_local=$(git -C "$ROOT_DIR" ls-files '.env' '.env.*' '.idea' '.vscode' '*.iml' || true)
if [ -n "$tracked_local" ]; then
  printf '%s\n%s\n' "Local secrets or editor metadata must not be tracked:" "$tracked_local" >&2
  exit 1
fi

found_plan=0
for plan in "$DOCS_PLANS"/*.md; do
  [ -e "$plan" ] || continue
  found_plan=1
  if ! grep -iq "status" "$plan" || ! grep -iq "completed" "$plan"; then
    printf '%s\n' "$plan must record completed status." >&2
    exit 1
  fi
  if ! grep -iq "verification" "$plan"; then
    printf '%s\n' "$plan must document verification." >&2
    exit 1
  fi
  if ! grep -Fq "make check" "$plan"; then
    printf '%s\n' "$plan must document make check verification." >&2
    exit 1
  fi
done

if [ "$found_plan" -eq 0 ]; then
  printf '%s\n' "docs/plans must contain completed markdown plans." >&2
  exit 1
fi
