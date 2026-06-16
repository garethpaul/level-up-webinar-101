#!/usr/bin/env sh
set -eu

ROOT_DIR=${1:?repository root is required}
GO_MOD="$ROOT_DIR/go.mod"
JWT_PLAN="$ROOT_DIR/docs/plans/2026-06-16-jwt-transitive-update.md"

jwt_selection_count=$(grep -Ec '^[[:space:]]*github\.com/golang-jwt/jwt/v5[[:space:]]+v5\.3\.1[[:space:]]+// indirect[[:space:]]*$' "$GO_MOD" || true)
jwt_module_count=$(grep -Ec '^[[:space:]]*github\.com/golang-jwt/jwt/v5[[:space:]]+' "$GO_MOD" || true)
if [ "$jwt_selection_count" -ne 1 ] || [ "$jwt_module_count" -ne 1 ]; then
  printf '%s\n' "go.mod must select exactly one JWT v5.3.1 dependency." >&2
  exit 1
fi

for plan_contract in \
  "status: completed" \
  "github.com/golang-jwt/jwt/v5 v5.3.1" \
  "go mod why -m github.com/golang-jwt/jwt/v5" \
  "go list -deps ./..." \
  "root and external-directory" \
  "five isolated hostile mutations"; do
  if ! grep -Fq "$plan_contract" "$JWT_PLAN"; then
    printf '%s\n' "JWT update plan must record completed verification: $plan_contract" >&2
    exit 1
  fi
done

for doc in README.md SECURITY.md VISION.md CHANGES.md; do
  if ! grep -Fq "JWT v5.3.1" "$ROOT_DIR/$doc"; then
    printf '%s\n' "$doc must document the selected JWT v5.3.1 dependency." >&2
    exit 1
  fi
done
