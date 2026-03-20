# Plan: PR & Release Flow — supertokens-golang

## Goal & Context

This plan describes the development workflow, branching model, changelog tooling, and release
pipeline used in `supertokens-python` (branch `dev`). Replicate this here.

The key idea is that **`dev` is the single integration branch** for all work in progress.
Version branches (`X.Y`) are stable release lines. Releases are prepared automatically by a
`dev-sync` workflow that fires on every push to `dev`.

---

## Branching Model

```
dev  ──────────────────────────────────────────────────────────►  (integration)
         │              │
         └─► release/v0.31.2 ──► 0.31  (version branch, merged for release)
                        │
                        └─► release/v0.32.0 ──► 0.32  (new version branch if minor bump)
```

- **`dev`** — all feature/fix PRs target `dev`. Never commit directly.
- **`X.Y`** (e.g. `0.31`) — version branches. Each minor version has its own branch.
  These are the branches that CI runs on (`push: branches: '[0-9]+.[0-9]+'`).
- **`release/vX.Y.Z`** — ephemeral release prep branch, created automatically by `dev-sync`.
  Always force-updated; never manually edited.

---

## Changelog Tooling: changie

`changie` manages changelog fragments. Each PR adds a fragment file under `.changes/unreleased/`.

### `.changie.yaml` (create at repo root)

```yaml
changesDir: .changes
unreleasedDir: unreleased
headerPath: header.md
changelogPath: CHANGELOG.md
versionExt: md
versionFormat: "## [{{.VersionNoPrefix}}] - {{.Time.Format \"2006-01-02\"}}"
kindFormat: "### {{.Kind}}"
changeFormat: "- {{.Body}}"

kinds:
  - label: Added
    auto: minor
  - label: Changed
    auto: patch
  - label: Fixed
    auto: patch
  - label: Breaking Changes
    auto: minor   # treat as minor while major version is 0
  - label: Infrastructure
    auto: none    # tooling changes don't affect the version
  - label: Deprecated
    auto: patch
  - label: Removed
    auto: minor
  - label: Security
    auto: patch

newlines:
  afterChangelogHeader: 2
  afterVersion: 1

replacements:
  - path: supertokens/constants.go   # adjust to the file that holds the version string
    find: 'SDKVersion = "[0-9]+\.[0-9]+\.[0-9]+"'
    replace: 'SDKVersion = "{{.VersionNoPrefix}}"'
    flags: ""
```

**Key points:**
- `auto: none` on `Infrastructure` means tooling-only PRs don't bump the version.
- `changie next auto` reads unreleased fragments and computes the next semver bump.
  It returns e.g. `v0.31.2`; strip the leading `v` with `version="${raw#v}"`.
- `changie batch "v$version"` collects fragments into `.changes/v{version}.md` and
  applies `replacements` (version bumps in source files).
- `changie merge` regenerates `CHANGELOG.md` from all version files.
- **Read the actual version constant file** before writing the `replacements` regex — the exact
  variable name and file path will differ from the example above.

### Adding a changelog entry (contributor flow)

```bash
changie new          # interactive prompt; writes .changes/unreleased/<timestamp>-<slug>.yaml
git add .changes/unreleased/
git commit -m "chore: add changelog fragment"
```

---

## `dev-sync.yml` — Automated Release Prep

Create `.github/workflows/dev-sync.yml`:

```yaml
name: "Sync Dev to Version Branch"

# On every push to dev, computes the next version from unreleased changelog
# fragments and creates or force-updates a release PR:
#
#   release/vX.Y.Z  →  X.Y  (version branch)
#
# Can also be triggered manually to override the bump type.

on:
  push:
    branches:
      - dev
  workflow_dispatch:
    inputs:
      bump:
        description: "Override bump type (default: auto — inferred from change kinds)"
        type: choice
        default: auto
        options:
          - auto
          - minor
          - patch

permissions:
  contents: write
  pull-requests: write

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          ref: dev
          fetch-depth: 0
          token: ${{ secrets.ALL_REPO_PAT }}

      - name: Setup git
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
          git fetch --all --tags

      - uses: miniscruff/changie-action@v2
        with:
          version: latest

      - name: Compute next version
        id: next
        run: |
          bump="${{ inputs.bump || 'auto' }}"

          if ! raw_version=$(changie next $bump 2>&1); then
            echo "changie next returned: $raw_version"
            echo "No bumpable fragments found. Skipping sync."
            echo "skip=true" >> $GITHUB_OUTPUT
            exit 0
          fi

          version="${raw_version#v}"
          version_branch=$(echo "$version" | grep -oE '^[0-9]+\.[0-9]+')

          echo "version=$version"               >> $GITHUB_OUTPUT
          echo "version_branch=$version_branch" >> $GITHUB_OUTPUT
          echo "pr_branch=release/v$version"    >> $GITHUB_OUTPUT

      - name: Create version branch if it doesn't exist
        if: steps.next.outputs.skip != 'true'
        run: |
          version_branch="${{ steps.next.outputs.version_branch }}"
          if ! git show-ref --verify --quiet "refs/remotes/origin/$version_branch"; then
            git checkout -b "$version_branch"
            git push origin "$version_branch"
            echo "Created new version branch: $version_branch"
          else
            echo "Version branch $version_branch already exists."
          fi

      - name: Set up Go (for any build/codegen steps)
        if: steps.next.outputs.skip != 'true'
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'   # adjust to latest supported

      - name: Build release branch
        if: steps.next.outputs.skip != 'true'
        run: |
          pr_branch="${{ steps.next.outputs.pr_branch }}"
          version_branch="${{ steps.next.outputs.version_branch }}"
          version="${{ steps.next.outputs.version }}"

          # Start from version branch, merge dev on top (dev wins on conflicts).
          if git show-ref --verify --quiet "refs/remotes/origin/$pr_branch"; then
            git checkout "$pr_branch"
            git reset --hard "origin/$version_branch"
          else
            git checkout -b "$pr_branch" "origin/$version_branch"
          fi

          git merge origin/dev --no-edit -X theirs

          changie batch "v$version"
          changie merge

          git add .changes/ CHANGELOG.md
          # Add other files touched by replacements, e.g.:
          # git add supertokens/constants.go
          git diff --staged --quiet \
            || git commit -m "chore: prepare release v$version"

          git push origin "$pr_branch" --force

      - name: Create or update pull request
        if: steps.next.outputs.skip != 'true'
        env:
          GH_TOKEN: ${{ secrets.ALL_REPO_PAT }}
        run: |
          version="${{ steps.next.outputs.version }}"
          pr_branch="${{ steps.next.outputs.pr_branch }}"
          target_branch="${{ steps.next.outputs.version_branch }}"

          pr_number=$(gh pr list \
            --base "$target_branch" \
            --head "$pr_branch" \
            --json number \
            --jq '.[0].number')

          notes=$(sed '1d' ".changes/v${version}.md" | sed '/./,$!d')

          body=$(cat <<EOF
          ## Release v${version}

          This PR is automatically kept in sync with \`dev\`. It will be force-updated on every push to \`dev\`.

          **Review checklist:**
          - [ ] Version \`${version}\` is correct (check version constant file)
          - [ ] Changelog entries below are accurate and complete

          **After merging:** trigger the [Release Pipeline](../actions/workflows/pipeline-release-tag.yml) with branch \`${target_branch}\`.

          ---

          ${notes}
          EOF
          )

          if [[ -n "$pr_number" ]]; then
            gh pr edit "$pr_number" \
              --title "chore: prepare release v${version}" \
              --body "$body" \
              --add-label "run-tests"
            echo "Updated PR #$pr_number"
          else
            gh pr create \
              --base "$target_branch" \
              --head "$pr_branch" \
              --title "chore: prepare release v${version}" \
              --body "$body" \
              --label "run-tests,Skip-Changelog"
            echo "Created new release PR"
          fi
```

### Required secret

`ALL_REPO_PAT` — a PAT with `contents: write` and `pull-requests: write` on this repo. Without
it, PRs created by `GITHUB_TOKEN` cannot trigger CI (GitHub blocks workflow runs on
bot-created PRs by default).

---

## PR Labels

| Label | Meaning |
|---|---|
| `run-tests` | Trigger the test matrix. Required by `test-gate`. |
| `skip-tests` | Explicitly bypass `test-gate` (e.g. docs-only PRs). |
| `Skip-Changelog` | Used on release prep PRs so they don't need a changelog fragment. |

`test-gate` (in each test workflow) enforces that every PR must have either `run-tests` or
`skip-tests` before it can be merged. See `cicd-update-plan.md` for the full `test-gate` spec.

---

## Release Pipeline

After merging the release PR (`release/vX.Y.Z → X.Y`), manually trigger
`.github/workflows/pipeline-release-tag.yml` with the target version branch. That workflow:

1. Tags the commit (e.g. `v0.31.2`)
2. Publishes to GitHub Releases / pkg.go.dev (via tag push)
3. Optionally updates `master` to point to the latest release

Read `supertokens-python/.github/workflows/pipeline-release-tag.yml` as the canonical reference,
then adapt for Go (no package registry publish step; Go modules are version-tagged automatically
by the tag push).

---

## Files to Create / Adapt

- `.changie.yaml` — changelog config (adjust `replacements` paths for this repo)
- `.changes/header.md` — short header for CHANGELOG.md (e.g. `# Changelog`)
- `.github/workflows/dev-sync.yml` — release prep automation (adapt build steps for Go)
- `.github/workflows/pipeline-release-tag.yml` — release/tag pipeline (port from Python)

Reference: `supertokens-python` branch `dev` has all of these in their final form.
