#!/usr/bin/env sh
set -eu

version=${1:?Go version is required}
case "$version" in
  go[0-9]*.[0-9]*.[0-9]*) ;;
  *)
    printf '%s\n' "unsupported Go version format: $version" >&2
    exit 1
    ;;
esac

numeric=${version#go}
if ! printf '%s\n' "$numeric" | awk -F. '
  NF == 3 && $1 ~ /^[0-9]+$/ && $2 ~ /^[0-9]+$/ && $3 ~ /^[0-9]+$/ {
    valid = ($1 > 1) || ($1 == 1 && ($2 > 26 || ($2 == 26 && $3 >= 4)))
  }
  END { exit valid ? 0 : 1 }
'; then
  printf '%s\n' "verification requires patched Go 1.26.4 or newer" >&2
  exit 1
fi
