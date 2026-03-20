# Plan: CI/CD Modernisation — supertokens-golang

## Goal & Context

SuperTokens SDKs each declare which CDI and FDI interface versions they support. The current CI
workflows test against every declared interface version independently, creating redundant runs
because many versions map to the same upstream release branch. The fix is to deduplicate at the
**branch level**: for each interface version, find the upstream branch where that version is the
*maximum* — the release designed for that era. Versions subsumed by a newer branch fall back to
the nearest branch.

Additionally, workflows likely lack the consistent PR labelling + skip pattern that avoids
re-running tests on identical code.

The Python SDK (`supertokens-python`, branch `dev`) and Node SDK (`supertokens-node`) have the
full implementation. This plan describes what to replicate here. The Golang SDK differs only in
the runtime dimension (`go-version` instead of `py-version`/`node-version`) and likely has no
auth-react or website tests (backends only: unit tests + backend-sdk tests).

**Start by reading the existing workflows in this repo** before making any changes, since the
current state may differ from what this plan assumes.

---

## New GitHub Action: `generate-test-matrix`

**Repo:** `supertokens/actions` — already implemented and built at
`generate-test-matrix/dist/index.js`. Commit the `dist/` alongside the source when tagging.

**Replaces:** `supertokens/get-supported-versions-action` + `supertokens/actions/get-versions-from-repo`
+ inline shell matrix logic.

### Inputs

| Input | Description |
|---|---|
| `github-token` | Required |
| `include-cdi` | `'true'` to vary CDI/core axis |
| `include-fdi` | `'true'` to vary FDI axis |
| `upstream-fdi-repos` | JSON array. First entry drives row deduplication; others are companion lookups. Default: `["supertokens-node"]` |
| `strategy` | `boundary` (default) or `primary-full` |
| `extra-axes` | JSON object, e.g. `{"go-version":["1.21","1.22","1.23"]}` |
| `latest-extra` | Override "latest" per axis. Default: last element |
| `working-directory` | Path to repo root with `*InterfaceSupported.json` |

### Outputs

| Output | Description |
|---|---|
| `testMatrix` | `{"include":[...]}` ready for `strategy.matrix` |
| `coreCdiVersionMap` | `{"5.4":"master"}` — core branch per CDI representative |
| `fdiVersionMap` | `{"supertokens-node":{"4.2":"master"},...}` |
| `fdiVersions` | `["1.19","2.0","4.2"]` — deduplicated FDI rep list |
| `extraAxes` | The `extra-axes` input echoed back as JSON |

### Strategies

**`boundary`** — all FDI reps × latest CDI + other CDI reps × latest FDI. Extra axes
boundary-reduced **only on the latest interface cell**; non-latest cells get only the anchor.

**`primary-full`** — full CDI × FDI cross-product, same extra-axis rule.

### Local validation prototype

`supertokens-python/scripts/generate-test-matrix.mjs` — run against local checkouts:

```bash
# unit-test scenario
node /path/to/supertokens-python/scripts/generate-test-matrix.mjs \
  --local-repos-dir ~/repos/supertokens \
  --include-cdi \
  --extra-axes '{"go-version":["1.21","1.22","1.23"]}' \
  --latest-extra '{"go-version":"1.23"}'

# backend-sdk-test scenario (only if FDI file exists)
node /path/to/supertokens-python/scripts/generate-test-matrix.mjs \
  --local-repos-dir ~/repos/supertokens \
  --include-cdi --include-fdi --strategy primary-full \
  --extra-axes '{"go-version":["1.21","1.22","1.23"]}' \
  --latest-extra '{"go-version":"1.23"}'
```

---

## Shared PR Pattern: skip-check + test-gate

**Every test workflow** must gain two additional jobs.

### skip-check (new)

```yaml
skip-check:
  runs-on: ubuntu-latest
  outputs:
    should-skip: ${{ steps.skip.outputs.should_skip }}
  steps:
    - uses: fkirc/skip-duplicate-actions@v5
      id: skip
      with:
        paths: '["**/*.go", "go.mod", "go.sum"]'
        skip_after_successful_duplicate: 'true'
```

Adjust `paths` to the Go source directories. Do NOT include the version constant file if it
contains the SDK version — version bumps should not invalidate the hash.

All `define-versions` and `test` jobs gain:
```yaml
if: |
  needs.skip-check.outputs.should-skip != 'true' &&
  (github.event_name != 'pull_request' || contains(github.event.pull_request.labels.*.name, 'run-tests'))
needs: [skip-check, ...]
```

### test-gate (new)

```yaml
test-gate:
  runs-on: ubuntu-latest
  needs: [test]
  if: always() && github.event_name == 'pull_request'
  steps:
    - name: Evaluate test requirement
      env:
        LABELS: ${{ toJSON(github.event.pull_request.labels.*.name) }}
        TEST_RESULT: ${{ needs.test.result }}
      run: |
        has_run_tests=$(echo "$LABELS"  | jq 'contains(["run-tests"])')
        has_skip_tests=$(echo "$LABELS" | jq 'contains(["skip-tests"])')

        if [[ "$has_skip_tests" == "true" ]]; then
          echo "Tests explicitly skipped via 'skip-tests' label."
          exit 0
        elif [[ "$has_run_tests" == "true" ]]; then
          if [[ "$TEST_RESULT" == "success" || "$TEST_RESULT" == "skipped" ]]; then
            echo "Tests passed (result: $TEST_RESULT)."
            exit 0
          else
            echo "::error::Tests did not pass (result: $TEST_RESULT)."
            exit 1
          fi
        else
          echo "::error::No test label found. Add 'run-tests' to run tests, or 'skip-tests' to bypass."
          exit 1
        fi
```

Add `test-gate` as a required status check on version branches and PRs.

---

## Workflow-by-Workflow Changes

The current `tests.yml` is a `workflow_dispatch`-only workflow — it needs to be replaced with
proper push/PR-triggered workflows. The old `tests-pass-check-pr.yml` (which polls for a prior
manual run) should be retired in favour of the `test-gate` pattern below.

### Unit tests — new `unit-test.yml`

Add triggers:
```yaml
on:
  pull_request:
    types: [opened, reopened, synchronize, labeled, unlabeled]
  push:
    branches:
      - '[0-9]+.[0-9]+'
```

**Target `define-versions` job:**

```yaml
define-versions:
  outputs:
    coreCdiVersionMap: ${{ steps.matrix.outputs.coreCdiVersionMap }}
    testMatrix: ${{ steps.matrix.outputs.testMatrix }}
  steps:
    - uses: actions/checkout@v4
    - uses: supertokens/actions/generate-test-matrix@main
      id: matrix
      with:
        github-token: ${{ secrets.GITHUB_TOKEN }}
        include-cdi: 'true'
        extra-axes: '{"go-version":["1.21","1.22","1.23"]}'
        latest-extra: '{"go-version":"1.23"}'
```

Test step: look up core branch:
```bash
coreVersion=$(echo '${{ needs.define-versions.outputs.coreCdiVersionMap }}' \
  | jq -r '.["${{ matrix.cdi-version }}"]')
```

Use `matrix.go-version` in `actions/setup-go` (or equivalent). Adjust the `go-version` list to
the versions actually supported by this repo.

---

### Backend SDK tests — new `backend-sdk-test.yml`

Only needed if `frontendDriverInterfaceSupported.json` exists. If the repo has no FDI file,
skip this workflow.

**Target:**

```yaml
define-versions:
  outputs:
    coreCdiVersionMap: ${{ steps.matrix.outputs.coreCdiVersionMap }}
    testMatrix: ${{ steps.matrix.outputs.testMatrix }}
  steps:
    - uses: actions/checkout@v4
    - uses: supertokens/actions/generate-test-matrix@main
      id: matrix
      with:
        github-token: ${{ secrets.GITHUB_TOKEN }}
        include-cdi: 'true'
        include-fdi: 'true'
        strategy: primary-full
        extra-axes: '{"go-version":["1.21"]}'
        latest-extra: '{"go-version":"1.21"}'
```

Adjust the `go-version` list to match what the repo actually tests. Use a single version if only
one Go version is needed for backend-sdk tests (same pattern as Node using only `"20"`).

`backend-sdk-testing-action` receives `version: ${{ matrix.fdi-version }}` (raw FDI version
number — unchanged from current).

---

### If there are website / auth-react tests

Unlikely for a pure backend SDK, but if they exist, follow the same pattern as `supertokens-node`:

- website-test: `upstream-fdi-repos: '["supertokens-website","supertokens-node"]'`
- auth-react: `upstream-fdi-repos: '["supertokens-auth-react","supertokens-node"]'`

---

## Files to Read First

Before making changes, read these files in this repo:

- All files under `.github/workflows/`
- `coreDriverInterfaceSupported.json`
- `frontendDriverInterfaceSupported.json` (may not exist)
- Any changelog/release tooling config (`.changie.yaml`, `Makefile`, etc.)

And these canonical reference implementations in `supertokens-python`:

- `.github/workflows/unit-test.yml`
- `.github/workflows/backend-sdk-testing.yml`

And in `supertokens-node`:

- `.github/workflows/unit-test.yml`
- `.github/workflows/backend-sdk-test.yml`
- `.claude/cicd-update-plan.md` (same plan as this one, tailored for Node)
