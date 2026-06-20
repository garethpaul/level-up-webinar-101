#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
TMP_ROOT=$(mktemp -d)
trap 'if [ -d "$TMP_ROOT" ]; then rm -rf -- "$TMP_ROOT"; fi' EXIT HUP INT TERM

write_valid_fixture() {
  rm -rf -- "$TMP_ROOT/fixture"
  mkdir -p "$TMP_ROOT/fixture/docs/plans"
  cat >"$TMP_ROOT/fixture/go.mod" <<'EOF'
module example.com/fixture

go 1.26.4

require github.com/twilio/twilio-go v1.30.9

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
)
EOF
  cat >"$TMP_ROOT/fixture/docs/plans/2026-06-16-jwt-transitive-update.md" <<'EOF'
status: completed
github.com/golang-jwt/jwt/v5 v5.3.1
go mod why -m github.com/golang-jwt/jwt/v5
go list -deps ./...
root and external-directory
five isolated hostile mutations
EOF
  for doc in README.md SECURITY.md VISION.md CHANGES.md; do
    printf '%s\n' "Selected dependency: JWT v5.3.1" >"$TMP_ROOT/fixture/$doc"
  done
}

expect_rejected() {
  name=$1
  if "$ROOT_DIR/scripts/check-jwt-dependency.sh" "$TMP_ROOT/fixture" >/dev/null 2>&1; then
    printf '%s\n' "JWT contract accepted hostile mutation: $name" >&2
    exit 1
  fi
}

replace_in_file() {
  expression=$1
  path=$2
  sed "$expression" "$path" >"$path.tmp"
  mv "$path.tmp" "$path"
}

write_valid_fixture
"$ROOT_DIR/scripts/check-jwt-dependency.sh" "$TMP_ROOT/fixture"

write_valid_fixture
replace_in_file 's/v5.3.1/v5.2.2/' "$TMP_ROOT/fixture/go.mod"
expect_rejected "version rollback"

write_valid_fixture
printf '%s\n' "github.com/golang-jwt/jwt/v5 v5.3.1 // indirect" >>"$TMP_ROOT/fixture/go.mod"
expect_rejected "duplicate selection"

write_valid_fixture
replace_in_file 's/status: completed/status: planned/' "$TMP_ROOT/fixture/docs/plans/2026-06-16-jwt-transitive-update.md"
expect_rejected "incomplete plan"

write_valid_fixture
replace_in_file '/go mod why -m/d' "$TMP_ROOT/fixture/docs/plans/2026-06-16-jwt-transitive-update.md"
expect_rejected "missing graph evidence"

write_valid_fixture
printf '%s\n' "Dependency guidance removed" >"$TMP_ROOT/fixture/SECURITY.md"
expect_rejected "missing security guidance"

printf '%s\n' "JWT dependency contract mutations rejected"
