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
  "docs/plans/2026-06-12-checkout-credential-boundary.md" \
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

TIMEOUT_PLAN="$ROOT_DIR/docs/plans/2026-06-10-explicit-twilio-timeout.md"
timeout_status_count=$(grep -Ec '^status: completed$' "$TIMEOUT_PLAN" || true)
timeout_work=$(awk '
  /^## Work Completed$/ { found = 1; next }
  /^## / && found { exit }
  found { print }
' "$TIMEOUT_PLAN")
timeout_verification=$(awk '
  /^## Verification Completed$/ { found = 1; next }
  /^## / && found { exit }
  found { print }
' "$TIMEOUT_PLAN")

if [ "$timeout_status_count" -ne 1 ] || [ -z "$timeout_work" ]; then
  printf '%s\n' "Twilio timeout plan must record one completed status and completed work." >&2
  exit 1
fi

if [ -z "$timeout_verification" ] || printf '%s\n' "$timeout_verification" | grep -Eiq '\b(pending|todo|tbd|not run)\b'; then
  printf '%s\n' "Twilio timeout plan must record finished verification without pending markers." >&2
  exit 1
fi

for evidence in \
  "make check" \
  "gofmt" \
  "go vet ./..." \
  "go mod verify" \
  "go test ./..." \
  "go build ./..." \
  "scripts/check-baseline.sh" \
  "git diff --check" \
  "27287294148" \
  "27402325505" \
  "58653e513077cdfae819dee89fe70ab4cc48d8e9" \
  "const twilioRequestTimeout = 10 * time.Second" \
  "client.SetTimeout(twilioRequestTimeout)" \
  "TestConfigureTwilioClientSetsRequestTimeout"; do
  if ! printf '%s\n' "$timeout_verification" | grep -Fq "$evidence"; then
    printf '%s\n' "Twilio timeout plan must preserve verification evidence: $evidence" >&2
    exit 1
  fi
done

if ! grep -Fq "scripts/check-baseline.sh" "$MAKEFILE"; then
  printf '%s\n' "Makefile must run scripts/check-baseline.sh from make check." >&2
  exit 1
fi

if ! grep -Fq "go vet ./..." "$MAKEFILE"; then
  printf '%s\n' "Makefile must run go vet from make lint." >&2
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
  "go-version: \"1.24.x\"" \
  "run: make check"; do
  if ! grep -Fq "$workflow_value" "$WORKFLOW"; then
    printf '%s\n' "Check workflow must keep $workflow_value" >&2
    exit 1
  fi
done

workflow_files=$(find "$ROOT_DIR/.github/workflows" -mindepth 1 -maxdepth 1 -type f -print | sort)
if [ "$workflow_files" != "$WORKFLOW" ]; then
  printf '%s\n' "Workflow inventory must contain only .github/workflows/check.yml." >&2
  exit 1
fi

if [ "$(grep -Fc "actions/checkout@" "$WORKFLOW")" -ne 1 ] ||
   [ "$(grep -Fc "persist-credentials:" "$WORKFLOW")" -ne 1 ] ||
   ! grep -Fq "persist-credentials: false" "$WORKFLOW" ||
   grep -Fq "persist-credentials: true" "$WORKFLOW"; then
  printf '%s\n' "Check workflow must use one pinned credential-free checkout." >&2
  exit 1
fi

CHECKOUT_PLAN="$ROOT_DIR/docs/plans/2026-06-12-checkout-credential-boundary.md"
if ! grep -Fq "status: completed" "$CHECKOUT_PLAN" ||
   ! grep -Fq "persist-credentials: false" "$CHECKOUT_PLAN" ||
   ! grep -Fq "hostile mutations rejected" "$CHECKOUT_PLAN"; then
  printf '%s\n' "Checkout credential plan must record completed verification." >&2
  exit 1
fi

guidance=$(cat "$ROOT_DIR/README.md" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md" | tr '\n' ' ' | tr -s '[:space:]' ' ' | tr '[:upper:]' '[:lower:]')
case "$guidance" in
  *"checkout credentials are not persisted"*"credential-free checkout"*) ;;
  *)
    printf '%s\n' "Repository guidance must document the credential-free checkout boundary." >&2
    exit 1
    ;;
esac

for target in "lint:" "test:" "build:" "fmt:" "check:"; do
  if ! grep -Fq "$target" "$MAKEFILE"; then
    printf '%s\n' "Makefile must expose the $target gate." >&2
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
  "go test ./..." \
  "go vet ./..." \
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
done

for module_line in \
  "module github.com/garethpaul/level-up-webinar-101" \
  "go 1.24" \
  "github.com/twilio/twilio-go v1.30.9"; do
  if ! grep -Fq "$module_line" "$ROOT_DIR/go.mod"; then
    printf '%s\n' "go.mod must keep module baseline: $module_line" >&2
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
