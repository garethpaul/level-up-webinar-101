#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
CHECKER="$ROOT_DIR/scripts/check-go-version.sh"

expect_allowed() {
  version=$1
  if ! "$CHECKER" "$version" >/dev/null 2>&1; then
    printf '%s\n' "expected $version to be allowed" >&2
    exit 1
  fi
}

expect_rejected() {
  version=$1
  if "$CHECKER" "$version" >/dev/null 2>&1; then
    printf '%s\n' "expected $version to be rejected" >&2
    exit 1
  fi
}

expect_allowed go1.26.4
expect_allowed go1.27.0
expect_rejected go1.26.3
expect_rejected go1.25.11
expect_rejected devel
expect_rejected 1.26.4

printf '%s\n' "Go toolchain policy cases passed"
