# Plan: MCP & Local Testing Setup — supertokens-golang

## Goal & Context

`supertokens-python` (branch `dev`) has a Docker-based local development setup that exposes
the full test environment to Claude Code via an MCP server. This lets Claude run tests and
linters directly from the IDE without leaving the editor. Replicate this here.

The canonical reference for all files is `supertokens-python`. The Go adaptation differs
only in runtime (Go instead of Python) and tooling (`go test`/`golangci-lint` instead of
pytest/ruff/pyright).

---

## Architecture Overview

```
Claude Code (host)
  └─ mcp-tool.mjs          ← stdio MCP proxy (runs on host)
       └─ HTTP POST /api/call
            └─ server.mjs  ← MCP server (runs inside Docker container)
                 ├─ test          → go test ./...
                 ├─ lint          → golangci-lint / go vet
                 ├─ cross_sdk_test → mocha cross-SDK tests
                 └─ task_status / task_cancel / task_list
                    test_results / test_output / test_runs
```

**Why split-process?** Tests need to reach `supertokens-core` on a Docker network. Keeping
the MCP server inside Docker avoids port-forwarding complexity for the test runner itself.
The host-side proxy is a thin stdio↔HTTP bridge — no logic lives there.

---

## Files to Create

### `compose.yaml` — core test services

Minimal stack for running tests locally (no MCP):

```yaml
services:
  core:
    image: supertokens/supertokens-dev-postgresql:${SUPERTOKENS_CORE_VERSION:-master}
    entrypoint: [
      "/usr/lib/supertokens/jre/bin/java",
      "-classpath", "/usr/lib/supertokens/core/*:/usr/lib/supertokens/plugin-interface/*:/usr/lib/supertokens/ee/*",
      "io.supertokens.Main", "/usr/lib/supertokens/", "DEV", "test_mode"
    ]
    ports:
      - ${SUPERTOKENS_CORE_PORT:-3567}:3567
    platform: linux/amd64
    depends_on: [oauth]
    environment:
      OAUTH_PROVIDER_PUBLIC_SERVICE_URL: http://oauth:4444
      OAUTH_PROVIDER_ADMIN_SERVICE_URL: http://oauth:4445
      OAUTH_PROVIDER_CONSENT_LOGIN_BASE_URL: http://localhost:3001/auth
      OAUTH_CLIENT_SECRET_ENCRYPTION_KEY: asdfasdfasdfasdfasdf
      INFO_LOG_PATH: "null"
      ERROR_LOG_PATH: "null"
    healthcheck:
      test: bash -c 'curl -s "http://127.0.0.1:3567/hello" | grep "Hello"'
      interval: 10s
      timeout: 5s
      retries: 5

  oauth:
    image: supertokens/oauth2-test:latest
    platform: linux/amd64
```

### `compose.mcp.yml` — full MCP stack

Extends `compose.yaml` with PostgreSQL (for license-key tests), full Hydra, and the MCP
container. Used when starting the MCP server for Claude Code:

```yaml
services:
  core:
    depends_on:
      oauth:
        condition: service_started
      pg:
        condition: service_healthy
    environment:
      POSTGRESQL_HOST: "pg"
      POSTGRESQL_PORT: "5432"
      POSTGRESQL_USER: "root"
      POSTGRESQL_PASSWORD: "root"
      POSTGRESQL_DATABASE_NAME: "postgres"
      SUPERTOKENS_LICENSE_KEY: ${SUPERTOKENS_LICENSE_KEY:-}

  pg:
    image: percona/percona-distribution-postgresql:13
    environment:
      POSTGRES_USER: root
      POSTGRES_PASSWORD: root
      POSTGRES_DB: postgres
    tmpfs:
      - /var/lib/postgresql/data:size=4g
    mem_limit: 4g
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U root -d postgres"]
      interval: 5s
      timeout: 3s
      retries: 10
    command: >
      postgres
        -c max_connections=1000
        -c fsync=off
        -c synchronous_commit=off
        -c full_page_writes=off

  oauth:
    image: oryd/hydra:v2.2.0
    environment:
      DSN: memory
      URLS_SELF_ISSUER: http://oauth:4444
      URLS_LOGIN: http://localhost:3001/auth/oauth/login
      URLS_CONSENT: http://localhost:3001/auth/oauth/consent
      URLS_LOGOUT: http://localhost:3001/auth/oauth/logout
      SECRETS_SYSTEM: thisIsATestSecretThatIsAtLeast32Chars
      SECRETS_COOKIE: thisIsATestCookieSecretAtLeast32
      SERVE_COOKIES_SAME_SITE_MODE: Lax
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:4445/health/alive"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s
    depends_on:
      pg:
        condition: service_healthy
    entrypoint: ["sh", "-c"]
    command: ["hydra migrate sql -e --yes && hydra serve all --dev"]

  mcp:
    build:
      context: .
      dockerfile: Dockerfile.mcp
    ports:
      - "127.0.0.1:${MCP_PORT:-3001}:3000"
    environment:
      MCP_TRANSPORT: sse
      MCP_PORT: "3000"
      GO_MCP_WORKSPACE: /workspace
      SUPERTOKENS_CORE_HOST: core
      SUPERTOKENS_CORE_PORT: "3567"
      SUPERTOKENS_LICENSE_KEY: ${SUPERTOKENS_LICENSE_KEY:-}
    volumes:
      - .:/workspace
      - go-cache:/root/go/pkg/mod
      - ${BACKEND_SDK_TESTING_PATH:-../backend-sdk-testing}:/cross-sdk-tests
      - cross-sdk-node-modules:/cross-sdk-tests/node_modules
    depends_on:
      core:
        condition: service_healthy
    mem_limit: 8g
    cpus: 8

volumes:
  go-cache:
  cross-sdk-node-modules:
```

### `Dockerfile.mcp`

```dockerfile
FROM golang:1.23-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl ca-certificates socat nodejs npm \
    && rm -rf /var/lib/apt/lists/*

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b /usr/local/bin latest

# Install MCP server Node dependencies
WORKDIR /opt/mcp-server
COPY mcp-server/package.json ./
RUN npm install --omit=dev && npm cache clean --force

COPY mcp-server/server.mjs ./

# Pre-warm Go module cache
COPY go.mod go.sum /tmp/sdk-deps/
RUN cd /tmp/sdk-deps && go mod download && rm -rf /tmp/sdk-deps

RUN mkdir -p /workspace /workspace/test-results /workspace/test-output

COPY <<'ENTRY' /opt/mcp-server/entrypoint.sh
#!/bin/bash
set -e
cd /workspace
go mod download 2>/dev/null || echo "[mcp] Warning: go mod download failed"

# Forward localhost:<core-port> → <core-host>:<core-port>
CORE_HOST="${SUPERTOKENS_CORE_HOST:-localhost}"
CORE_PORT="${SUPERTOKENS_CORE_PORT:-3567}"
if [ "$CORE_HOST" != "localhost" ] && [ "$CORE_HOST" != "127.0.0.1" ]; then
    echo "[mcp] Forwarding localhost:${CORE_PORT} → ${CORE_HOST}:${CORE_PORT}"
    socat TCP-LISTEN:${CORE_PORT},fork,reuseaddr TCP:${CORE_HOST}:${CORE_PORT} &
    export SUPERTOKENS_CORE_HOST=localhost
fi

exec node /opt/mcp-server/server.mjs
ENTRY
RUN chmod +x /opt/mcp-server/entrypoint.sh

EXPOSE 3000
ENTRYPOINT ["/opt/mcp-server/entrypoint.sh"]
```

### `mcp-server/server.mjs` — MCP server (inside Docker)

Node.js MCP server (Node is only for the MCP plumbing; the tests themselves run via `go test`).
Exposes both SSE transport and a stateless HTTP API (`POST /api/call`). Tools:

| Tool | Description |
|---|---|
| `test` | Run `go test ./...`. Accepts `filter` (`-run` regex), `path`, `maxFail` (`-failfast`). |
| `lint` | Run `golangci-lint run` and/or `go vet`. Accepts `tool: "all"\|"golangci-lint"\|"vet"`. |
| `cross_sdk_test` | Run Mocha cross-SDK tests from `/cross-sdk-tests`. Accepts `grep`. |
| `task_status` | Poll a running task by `taskId`. Returns progress + results when done. |
| `task_cancel` | Cancel a running task. |
| `task_list` | List all running and recently completed tasks. |
| `test_results` | Browse archived test runs. Accepts `runId`, `filter`, `testName`. |
| `test_output` | Get stdout/stderr for a specific test by `testId`. |
| `test_runs` | List archived test run IDs. |

Implementation pattern (copy from `supertokens-python/mcp-server/server.mjs`, then adapt):
- Replace `pytest` invocation with `go test -v -json -run <filter> ./...`
- Replace `ruff`/`pyright` with `golangci-lint run` and/or `go vet ./...`
- Replace JUnit XML parsing with Go's `-json` output (`go test -json` emits structured events)
  — parse `Action: "pass"/"fail"/"output"` lines to build per-test results
- Per-test output capture: collect all `"output"` lines for each `Test:` into JSON files in
  `test-output/` (mirrors `pytest_capture.py` approach but reading from stdout stream)
- Test run archival, task registry, truncation logic — copy unchanged

Environment variables consumed by `server.mjs`:
- `GO_MCP_WORKSPACE` — path to the SDK repo (default `/workspace`)
- `SUPERTOKENS_CORE_HOST` / `SUPERTOKENS_CORE_PORT`
- `MCP_TRANSPORT` — `sse` (default) or `stdio`
- `MCP_PORT` — HTTP port (default `3000`)

### `mcp-server/mcp-tool.mjs` — host-side proxy

Thin stdio↔HTTP bridge. Copy `supertokens-python/mcp-server/mcp-tool.mjs` unchanged and
rename the server name from `"python-build-tools"` to `"go-build-tools"`. Update tool
registrations to match the tools in `server.mjs`.

### `mcp-server/package.json`

```json
{
  "name": "go-build-tools",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "dependencies": {
    "@modelcontextprotocol/sdk": "^1.12.0",
    "zod": "^3.24.0",
    "zod-to-json-schema": "^3.24.0"
  }
}
```

### `.mcp.json` — project MCP config

```json
{
  "mcpServers": {
    "go-build-tools": {
      "command": "node",
      "args": ["mcp-server/mcp-tool.mjs"]
    }
  }
}
```

### `mcp.env` (gitignored)

```
MCP_PORT=3001
```

Add `mcp.env` to `.gitignore`. This file is read by `mcp-tool.mjs` to determine which port
the Docker container is listening on.

### `.claude/settings.json` — permissions

```json
{
  "permissions": {
    "allow": [
      "mcp__github__get_me",
      "mcp__github__search_repositories"
    ]
  }
}
```

Extend with `Read(...)` and `Bash(...)` allows as needed during development sessions.

---

## Go `test -json` Output Parsing

`go test -json` emits one JSON object per line. The important fields:

```json
{ "Action": "run",    "Test": "TestFoo" }
{ "Action": "output", "Test": "TestFoo", "Output": "    --- PASS: TestFoo\n" }
{ "Action": "pass",   "Test": "TestFoo", "Elapsed": 0.001 }
{ "Action": "fail",   "Test": "TestFoo", "Elapsed": 0.123 }
{ "Action": "output", "Test": null, "Output": "ok  \tgithub.com/...\n" }
```

Collect all `"output"` lines per `Test` name; write to `test-output/<hashed-name>.json`.
Use `"pass"`/`"fail"` events to build the summary (total, passed, failed, skipped).
Events with `"Test": null` are package-level — add to a global log, not per-test files.

---

## Developer Workflow

### Start the full MCP stack

```bash
# First time or after Dockerfile changes:
docker compose -f compose.mcp.yml build mcp

# Start everything:
docker compose -f compose.mcp.yml up -d

# Watch logs:
docker compose -f compose.mcp.yml logs -f mcp
```

### Run tests without MCP (CI / quick local)

```bash
docker compose up --wait        # starts core + oauth
go test ./... -p 1              # run all tests
```

### Use from Claude Code

Claude Code reads `.mcp.json` automatically. Once the MCP stack is up, the `go-build-tools`
server appears in Claude Code's tool palette. Use the `test`, `lint`, and `task_status` tools
to run tests and check results without leaving the editor.

---

## Reference Implementation

All files have a direct counterpart in `supertokens-python` (branch `dev`):

| This repo | supertokens-python |
|---|---|
| `compose.yaml` | `compose.yaml` |
| `compose.mcp.yml` | `compose.mcp.yml` |
| `Dockerfile.mcp` | `Dockerfile.mcp` |
| `mcp-server/server.mjs` | `mcp-server/server.mjs` |
| `mcp-server/mcp-tool.mjs` | `mcp-server/mcp-tool.mjs` |
| `mcp-server/package.json` | `mcp-server/package.json` |
| `.mcp.json` | `.mcp.json` |
| `mcp.env` | `mcp.env` |

**Key adaptation points** from Python → Go:
- `pytest` → `go test -v -json`; parse JSON stream instead of JUnit XML
- `ruff`/`pyright` → `golangci-lint run`/`go vet`
- `pytest_capture.py` plugin → inline JSON-stream parser in `server.mjs`
- `python:3.12-slim` base → `golang:1.23-bookworm`
- pip cache volume → Go module cache volume (`~/go/pkg/mod`)
- Editable install step → `go mod download`
