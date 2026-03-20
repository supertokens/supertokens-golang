# Plan: DevEx & CI Improvements — supertokens-golang

## Reference

`supertokens-python` (branch `dev`) is the canonical implementation. These improvements have
already been applied there — use its workflows as the reference when implementing.

Key differences from Python:
- Version lives in `supertokens/constants.go` (not `constants.py` + `setup.py`)
- Linting is `go vet` + `golangci-lint` (not pre-commit + pyright)
- Test result surfacing uses `gotestsum` for JUnit XML generation
- No `skip-duplicate-actions` in use yet — needs to be added as part of CI modernization
- No release pipeline exists yet — needs to be created

---

## Items

### 1. CI lint workflow

**Status: Not started** — `pre-commit-hook-run.yml` exists but runs the git hook in a
simulated environment. A proper lint workflow should run `go vet` and `golangci-lint`
directly.

Create `.github/workflows/lint.yml`:

```yaml
name: "Lint"

on:
  pull_request:
    types: [opened, reopened, synchronize]
  push:
    branches:
      - '[0-9]+.[0-9]+'

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

permissions:
  contents: read

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Run go vet
        run: go vet ./...

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
```

If there is no `.golangci.yml` in the repo, create a minimal one enabling at least `go vet`,
`staticcheck`, and `errcheck`. Read the existing config first to avoid overwriting anything.

Add `lint` as a required status check alongside `test-gate`.

---

### 2. Changie fragment enforcement

**Status: Not started** — currently uses `dangoslen/changelog-enforcer@v2` which checks for
changes to `CHANGELOG.md` directly. Needs to be replaced with changie fragment enforcement
once changie is set up.

In Python, this is implemented as the `lint-changelog` job inside `lint-pr.yml` (not a
separate workflow). The same pattern should be used here — combine PR title linting and
changelog enforcement into a single `lint-pr.yml` workflow.

Source path filter for this repo: `grep -E '\.go$|^go\.mod$|^go\.sum$'`

---

### 3. PR title linting

**Status: Partially done** — `lint-pr-title.yml` exists with `amannn/action-semantic-pull-request@v3`.

Updates needed:
- Bump `@v3` → **`@v6`** (v6 uses Node 24, resolving the Node.js 20 deprecation warning)
- Add explicit `types` list and `ignoreLabels` (currently uses `validateSingleCommit` only)
- Add `permissions: contents: read`
- Consider merging into a combined `lint-pr.yml` with changelog enforcement (matching the
  Python pattern)

---

### 4. `paths-ignore` on test workflows

**Status: Not started** — no test workflow has `paths-ignore`.

Once test workflows are modernized (per `cicd-update-plan.md`), add to the `push` trigger:

```yaml
push:
  branches:
    - '[0-9]+.[0-9]+'
  paths-ignore:
    - '**.md'
    - 'docs/**'
```

This only affects push-triggered runs on release branches. PR-triggered runs are unaffected.

---

### 5. Test result surfacing

**Status: Not started** — no test result annotation in PRs.

`go test -v -json` produces structured output that can be converted to JUnit XML with
`gotestsum`. Once the test workflow is modernized:

```yaml
      - name: Install gotestsum
        run: go install gotest.tools/gotestsum@latest

      - name: Run tests
        run: gotestsum --junitfile test-results/junit.xml -- ./... -p 1 -v -count=1

      - uses: dorny/test-reporter@v1
        if: always()
        with:
          name: "Test Results"
          path: test-results/junit.xml
          reporter: java-junit
```

Replace `go test ./...` with `gotestsum` in the test workflow.

---

### 6. `skip-duplicate-actions` with version file exclusion

**Status: Not started** — no duplicate action skipping in use.

When adding `skip-duplicate-actions` to test workflows, exclude the version constant file to
prevent release PRs (which only bump version strings) from re-running the full test matrix:

```yaml
      - uses: fkirc/skip-duplicate-actions@v5
        id: skip
        with:
          paths: '["supertokens/**", "recipe/**", "ingredients/**", "test/**", "go.mod", "go.sum"]'
          paths_ignore: '["supertokens/constants.go"]'
          skip_after_successful_duplicate: 'true'
```

**Lesson from Python:** Without `paths_ignore` for the version file, the dev-sync release
PR triggers a full test re-run for what is just a version string change.

Adjust the `paths` list to match the actual source directories in this repo.

---

### 7. Explicit `permissions:` on test workflows

**Status: Not started** — no workflow has `permissions:` defined.

Add to the top level of every test/lint workflow:

```yaml
permissions:
  contents: read
```

Applies to all current workflows: `tests.yml`, `tests-pass-check-pr.yml`,
`lint-pr-title.yml`, `pre-commit-hook-run.yml`, `enforce-go-mod-tidy.yml`,
`github-actions-changelog.yml`.

---

### 8. Release pipeline

**Status: Not started** — no release pipeline exists.

This repo needs a `pipeline-release-tag.yml` and `dev-sync.yml` modeled after the Python
implementation. Key design points from the Python implementation:

**Auto-trigger pattern:** The release pipeline should fire automatically when a dev-sync
release PR is merged, with manual dispatch as a fallback for overrides:

```yaml
on:
  pull_request:
    types: [closed]
    branches:
      - '[0-9]+.[0-9]+'
  workflow_dispatch:
    inputs:
      branch: ...
      skip-test-checks: ...
      skip-other-version-checks: ...
```

**Guard condition:** Only run for merged release PRs or manual dispatch:

```yaml
  setup:
    if: >-
      github.event_name == 'workflow_dispatch' ||
      (github.event.pull_request.merged == true &&
       startsWith(github.event.pull_request.head.ref, 'release/v'))
```

**Branch resolution:** Compute from event type, output for downstream jobs:

```yaml
      - name: Determine release branch
        id: params
        run: |
          if [[ "${{ github.event_name }}" == "workflow_dispatch" ]]; then
            echo "branch=${{ inputs.branch }}" >> $GITHUB_OUTPUT
          else
            echo "branch=${{ github.event.pull_request.base.ref }}" >> $GITHUB_OUTPUT
          fi
```

**Skip flag condition:** Use `!= true` (not `== 'false'`) so the check runs correctly when
inputs are empty (auto-trigger case):

```yaml
      - if: inputs.skip-other-version-checks != true
        name: Check if core and frontend released
```

The `publish` environment approval gates the actual release. Merging the PR just starts the
pipeline without requiring a separate manual dispatch step.

See the Python `pipeline-release-tag.yml` and `dev-sync.yml` on branch `dev` for the
complete implementation.
