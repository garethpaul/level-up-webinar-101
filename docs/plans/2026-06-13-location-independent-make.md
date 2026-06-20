# Location-Independent Webinar Verification

status: completed

## Context

Absolute Makefile invocations run the baseline script and Go commands from the
caller's directory instead of the checkout, so documented gates fail outside
the repository.

## Scope

1. Derive the checkout root from the loaded Makefile.
2. Run formatting checks, vet, module verification, tests, builds, and the
   baseline script from that root.
3. Add exact Makefile, completed-plan, external-run, and guidance contracts.
4. Preserve Twilio request behavior, dependency locks, workflow policy, and all
   existing stacked-branch artifacts unchanged.

## Verification Plan

- Run `make check` and the other non-mutating Make gates from the checkout and
  through an absolute Makefile path from a temporary directory.
- Run shell syntax, Go formatting/vet/test/build/module verification, checker
  compilation where applicable, and diff checks.
- Reject root derivation, Go working directory, baseline-script location,
  plan status, plan evidence, and documentation mutations independently.
- Inspect intended paths, secret patterns, conflict markers, generated
  artifacts, and Go/runtime/workflow changes before commit.

## Risk And Rollback

This changes verification path resolution only. Rollback restores the relative
recipes and removes their checker, plan, and documentation contracts.

## Verification

- `make lint`, `make test`, `make build`, and `make check` passed in root and external-directory
  runs through an absolute Makefile path.
- Shell syntax, Go formatting, vet, module verification, unit tests, builds,
  baseline checks, and `git diff --check` passed.
- Verification rejected six isolated hostile mutations by their intended
  contracts: root derivation, Go working directory, baseline-script location,
  plan status, plan evidence, and README guidance.
- The intended five-file diff passed secret-pattern, conflict-marker,
  generated-artifact, and Go/runtime/workflow change audits.
