#!/usr/bin/env node

// MCP Build Tools Server for supertokens-golang
// Provides test, lint, and task management via MCP (SSE) and stateless HTTP API.

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { SSEServerTransport } from "@modelcontextprotocol/sdk/server/sse.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import { createServer } from "http";
import { spawn } from "child_process";
import { randomBytes, createHash } from "crypto";
import {
  readFileSync,
  writeFileSync,
  mkdirSync,
  existsSync,
  readdirSync,
  rmSync,
  statSync,
} from "fs";
import { join, basename } from "path";
import { zodToJsonSchema } from "zod-to-json-schema";

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------
const WORKSPACE = process.env.GO_MCP_WORKSPACE || "/workspace";
const TEST_TIMEOUT_MS =
  (parseInt(process.env.GO_MCP_TEST_TIMEOUT_SECS, 10) || 3600) * 1000;
const LINT_TIMEOUT_MS =
  (parseInt(process.env.GO_MCP_LINT_TIMEOUT_SECS, 10) || 600) * 1000;
const MAX_OUTPUT_CHARS =
  parseInt(process.env.GO_MCP_MAX_OUTPUT_CHARS, 10) || 60000;
const MAX_TEST_RUNS =
  parseInt(process.env.GO_MCP_MAX_TEST_RUNS, 10) || 10;
const CROSS_SDK_TIMEOUT_MS =
  (parseInt(process.env.GO_MCP_CROSS_SDK_TIMEOUT_SECS, 10) || 1800) * 1000;
const CROSS_SDK_DIR = "/cross-sdk-tests";
const TRANSPORT = process.env.MCP_TRANSPORT || "sse";
const PORT = parseInt(process.env.MCP_PORT, 10) || 3000;

const TEST_RESULTS_DIR = join(WORKSPACE, "test-results");
const TEST_OUTPUT_DIR = join(WORKSPACE, "test-output");
const TEST_RUNS_DIR = join(WORKSPACE, ".test-runs");

// ---------------------------------------------------------------------------
// Utility helpers
// ---------------------------------------------------------------------------
function generateId(prefix = "") {
  return prefix + randomBytes(6).toString("hex");
}

function truncateOutput(text, max = MAX_OUTPUT_CHARS) {
  if (!text || text.length <= max) return text || "";
  const half = Math.floor(max / 2) - 50;
  return (
    text.slice(0, half) +
    `\n\n... [TRUNCATED ${text.length - max} chars] ...\n\n` +
    text.slice(-half)
  );
}

function ensureDir(dir) {
  if (!existsSync(dir)) mkdirSync(dir, { recursive: true });
}

function nowMs() {
  return Date.now();
}

// ---------------------------------------------------------------------------
// Task registry
// ---------------------------------------------------------------------------
const tasks = new Map();

function pruneOldTasks() {
  const cutoff = nowMs() - 3600_000; // 1 hour
  for (const [id, t] of tasks) {
    if (
      (t.status === "completed" || t.status === "failed" || t.status === "cancelled") &&
      t.endTime < cutoff
    ) {
      tasks.delete(id);
    }
  }
}

// ---------------------------------------------------------------------------
// Test-run archival
// ---------------------------------------------------------------------------
function generateRunId() {
  const d = new Date();
  const pad = (n) => String(n).padStart(2, "0");
  return `${d.getFullYear()}${pad(d.getMonth() + 1)}${pad(d.getDate())}_${pad(d.getHours())}${pad(d.getMinutes())}${pad(d.getSeconds())}_${randomBytes(3).toString("hex")}`;
}

function archiveTestRun(task) {
  ensureDir(TEST_RUNS_DIR);
  const runId = generateRunId();
  const runDir = join(TEST_RUNS_DIR, runId);
  mkdirSync(runDir, { recursive: true });

  // Copy per-test output JSON files
  if (existsSync(TEST_OUTPUT_DIR)) {
    const outputDir = join(runDir, "test-output");
    mkdirSync(outputDir, { recursive: true });
    for (const f of readdirSync(TEST_OUTPUT_DIR)) {
      if (f.endsWith(".json")) {
        writeFileSync(
          join(outputDir, f),
          readFileSync(join(TEST_OUTPUT_DIR, f))
        );
      }
    }
  }

  // Save summary
  const summary = {
    runId,
    taskId: task.id,
    taskType: task.type,
    startTime: task.startTime,
    endTime: task.endTime,
    status: task.status,
    testSummary: task.testSummary || null,
    command: task.command,
  };
  writeFileSync(join(runDir, "summary.json"), JSON.stringify(summary, null, 2));

  pruneTestRuns();
  return runId;
}

function listTestRuns() {
  if (!existsSync(TEST_RUNS_DIR)) return [];
  return readdirSync(TEST_RUNS_DIR)
    .filter((d) => {
      try {
        return statSync(join(TEST_RUNS_DIR, d)).isDirectory();
      } catch {
        return false;
      }
    })
    .sort()
    .reverse()
    .map((d) => {
      try {
        const s = JSON.parse(
          readFileSync(join(TEST_RUNS_DIR, d, "summary.json"), "utf-8")
        );
        return s;
      } catch {
        return { runId: d };
      }
    });
}

function getRunDir(runId) {
  if (!runId) {
    const runs = listTestRuns();
    if (runs.length === 0) return null;
    runId = runs[0].runId;
  }
  const dir = join(TEST_RUNS_DIR, runId);
  return existsSync(dir) ? dir : null;
}

function getArchivedTestOutputs(runId) {
  const dir = getRunDir(runId);
  if (!dir) return [];
  const outputDir = join(dir, "test-output");
  if (!existsSync(outputDir)) return [];
  const results = [];
  for (const f of readdirSync(outputDir)) {
    if (f.endsWith(".json")) {
      try {
        results.push(JSON.parse(readFileSync(join(outputDir, f), "utf-8")));
      } catch {
        // skip corrupt files
      }
    }
  }
  return results;
}

function pruneTestRuns() {
  if (!existsSync(TEST_RUNS_DIR)) return;
  const dirs = readdirSync(TEST_RUNS_DIR)
    .filter((d) => {
      try {
        return statSync(join(TEST_RUNS_DIR, d)).isDirectory();
      } catch {
        return false;
      }
    })
    .sort()
    .reverse();
  for (let i = MAX_TEST_RUNS; i < dirs.length; i++) {
    try {
      rmSync(join(TEST_RUNS_DIR, dirs[i]), { recursive: true, force: true });
    } catch {
      // ignore
    }
  }
}

// ---------------------------------------------------------------------------
// Go test -json output parser
// ---------------------------------------------------------------------------
// `go test -json` emits one JSON object per line:
// { "Action": "run",    "Test": "TestFoo", "Package": "..." }
// { "Action": "output", "Test": "TestFoo", "Output": "...\n" }
// { "Action": "pass",   "Test": "TestFoo", "Elapsed": 0.001 }
// { "Action": "fail",   "Test": "TestFoo", "Elapsed": 0.123 }
// { "Action": "skip",   "Test": "TestFoo", "Elapsed": 0.0 }
// Events with Test: undefined are package-level.

function parseGoTestJson(stdout) {
  const testMap = new Map(); // testId -> { output, status, elapsed, pkg }
  const lines = stdout.split("\n");

  for (const line of lines) {
    if (!line.trim()) continue;
    let ev;
    try {
      ev = JSON.parse(line);
    } catch {
      continue; // non-JSON line (e.g. build errors)
    }

    if (!ev.Test) continue; // package-level event

    const testId = ev.Package ? `${ev.Package}/${ev.Test}` : ev.Test;

    if (!testMap.has(testId)) {
      testMap.set(testId, {
        testId,
        testName: ev.Test,
        className: ev.Package || "",
        stdout: "",
        stderr: "",
        resultType: "RUNNING",
        duration: 0,
        failureMessage: "",
        failureTrace: "",
      });
    }

    const t = testMap.get(testId);

    switch (ev.Action) {
      case "output":
        t.stdout += ev.Output || "";
        break;
      case "pass":
        t.resultType = "SUCCESS";
        t.duration = Math.round((ev.Elapsed || 0) * 1000);
        break;
      case "fail":
        t.resultType = "FAILURE";
        t.duration = Math.round((ev.Elapsed || 0) * 1000);
        // Extract failure info from captured output
        t.failureMessage = extractFailureMessage(t.stdout);
        t.failureTrace = t.stdout;
        break;
      case "skip":
        t.resultType = "SKIPPED";
        t.duration = Math.round((ev.Elapsed || 0) * 1000);
        break;
    }
  }

  return testMap;
}

function extractFailureMessage(output) {
  // Look for common Go test failure patterns
  const lines = output.split("\n");
  for (const line of lines) {
    if (line.includes("FAIL") || line.includes("Error") || line.includes("panic")) {
      return line.trim().substring(0, 200);
    }
  }
  return "";
}

function buildGoTestSummary(testMap) {
  let total = 0, passed = 0, failed = 0, skipped = 0;
  const failures = [];

  for (const t of testMap.values()) {
    if (t.resultType === "RUNNING") continue; // incomplete
    total++;
    if (t.resultType === "SUCCESS") passed++;
    else if (t.resultType === "FAILURE") {
      failed++;
      failures.push({
        className: t.className,
        testName: t.testName,
        message: t.failureMessage,
        trace: t.failureTrace,
      });
    } else if (t.resultType === "SKIPPED") skipped++;
  }

  return { total, passed, failed, errors: 0, skipped, failures };
}

function writeGoTestOutputs(testMap, outputDir) {
  ensureDir(outputDir);
  for (const t of testMap.values()) {
    if (t.resultType === "RUNNING") continue;
    const hash = createHash("md5").update(t.testId).digest("hex").slice(0, 8);
    const safeName = t.testId.replace(/[^a-zA-Z0-9._-]/g, "_").slice(0, 120);
    const filename = `${safeName}_${hash}.json`;
    writeFileSync(join(outputDir, filename), JSON.stringify(t, null, 2));
  }
}

// ---------------------------------------------------------------------------
// Mocha JSON parser (cross-SDK tests)
// ---------------------------------------------------------------------------
function parseMochaJson(jsonPath) {
  if (!existsSync(jsonPath)) return null;
  try {
    const data = JSON.parse(readFileSync(jsonPath, "utf-8"));
    const stats = data.stats || {};
    const total = stats.tests || 0;
    const passed = stats.passes || 0;
    const failed = stats.failures || 0;
    const skipped = stats.pending || 0;

    const failures = [];
    for (const t of data.failures || []) {
      failures.push({
        className: t.fullTitle || "",
        testName: t.title || "",
        message: t.err?.message || "",
        trace: t.err?.stack || "",
      });
    }

    return { total, passed, failed, errors: 0, skipped, failures };
  } catch {
    return null;
  }
}

function convertMochaToTestOutputs(jsonPath, outputDir) {
  if (!existsSync(jsonPath)) return;
  let data;
  try {
    data = JSON.parse(readFileSync(jsonPath, "utf-8"));
  } catch {
    return;
  }

  ensureDir(outputDir);

  const allTests = [
    ...(data.passes || []).map((t) => ({ ...t, _result: "SUCCESS" })),
    ...(data.failures || []).map((t) => ({ ...t, _result: "FAILURE" })),
    ...(data.pending || []).map((t) => ({ ...t, _result: "SKIPPED" })),
  ];

  for (const t of allTests) {
    const testId = t.fullTitle || t.title || "unknown";
    const className = t.fullTitle
      ? t.fullTitle.replace(/ /g, ".").replace(/::/g, ".")
      : "cross-sdk";
    const testName = t.title || "unknown";

    const output = {
      testId,
      className,
      testName,
      stdout: "",
      stderr: "",
      resultType: t._result,
      duration: t.duration || 0,
      failureMessage: t.err?.message || "",
      failureTrace: t.err?.stack || "",
    };

    const hash = createHash("md5").update(testId).digest("hex").slice(0, 8);
    const safeName = testId.replace(/[^a-zA-Z0-9._-]/g, "_").slice(0, 120);
    const filename = `${safeName}_${hash}.json`;
    writeFileSync(join(outputDir, filename), JSON.stringify(output, null, 2));
  }
}

// ---------------------------------------------------------------------------
// Async task runner
// ---------------------------------------------------------------------------
function startProcessTask(taskType, command, args, timeoutMs) {
  pruneOldTasks();

  const taskId = generateId("task_");
  const task = {
    id: taskId,
    type: taskType,
    status: "running",
    command: `${command} ${args.join(" ")}`,
    startTime: nowMs(),
    endTime: null,
    stdout: "",
    stderr: "",
    exitCode: null,
    testSummary: null,
    runId: null,
    _process: null,
    _timeout: null,
  };
  tasks.set(taskId, task);

  const proc = spawn(command, args, {
    cwd: WORKSPACE,
    env: { ...process.env },
    stdio: ["ignore", "pipe", "pipe"],
  });

  task._process = proc;

  proc.stdout.on("data", (chunk) => {
    const text = chunk.toString();
    task.stdout += text;
    process.stderr.write(text); // stream to container logs
  });

  proc.stderr.on("data", (chunk) => {
    const text = chunk.toString();
    task.stderr += text;
    process.stderr.write(text);
  });

  const timer = setTimeout(() => {
    if (task.status === "running") {
      task.status = "failed";
      task.endTime = nowMs();
      task.stderr += `\n[TIMEOUT after ${timeoutMs / 1000}s]`;
      try {
        proc.kill("SIGTERM");
        setTimeout(() => {
          try {
            proc.kill("SIGKILL");
          } catch {}
        }, 5000);
      } catch {}
    }
  }, timeoutMs);
  task._timeout = timer;

  proc.on("close", (code) => {
    clearTimeout(timer);
    if (task.status === "running") {
      task.exitCode = code;
      task.status = code === 0 ? "completed" : "failed";
      task.endTime = nowMs();

      // Post-process test results
      if (taskType === "test") {
        const testMap = parseGoTestJson(task.stdout);
        task.testSummary = buildGoTestSummary(testMap);

        // Clean previous output
        if (existsSync(TEST_OUTPUT_DIR)) {
          rmSync(TEST_OUTPUT_DIR, { recursive: true, force: true });
        }
        writeGoTestOutputs(testMap, TEST_OUTPUT_DIR);

        task.runId = archiveTestRun(task);
      } else if (taskType === "cross_sdk_test") {
        const mochaJsonPath = join(CROSS_SDK_DIR, "test-results.json");

        // Clean previous data so archive only contains cross-SDK results
        if (existsSync(TEST_OUTPUT_DIR)) {
          rmSync(TEST_OUTPUT_DIR, { recursive: true, force: true });
        }
        ensureDir(TEST_OUTPUT_DIR);

        // Convert Mocha JSON to per-test output files
        convertMochaToTestOutputs(mochaJsonPath, TEST_OUTPUT_DIR);

        // Parse summary
        task.testSummary = parseMochaJson(mochaJsonPath);

        // Archive
        task.runId = archiveTestRun(task);
      }
    }
  });

  proc.on("error", (err) => {
    clearTimeout(timer);
    if (task.status === "running") {
      task.status = "failed";
      task.endTime = nowMs();
      task.stderr += `\nProcess error: ${err.message}`;
    }
  });

  return taskId;
}

// ---------------------------------------------------------------------------
// Result formatters
// ---------------------------------------------------------------------------
function formatTestResult(task) {
  const lines = [];

  if (task.testSummary) {
    const s = task.testSummary;
    const header = task.status === "completed" ? "ALL TESTS PASSED" : "TESTS FAILED";
    lines.push(header, "");
    lines.push("--- Summary ---");
    lines.push(
      `Total: ${s.total}  |  Passed: ${s.passed}  |  ` +
      `Failed: ${s.failed}  |  Errors: ${s.errors}  |  Skipped: ${s.skipped}`
    );

    if (s.failures && s.failures.length > 0) {
      lines.push("", "--- Failure details ---");
      const cap = 25;
      for (const f of s.failures.slice(0, cap)) {
        lines.push(`\n> ${f.className}::${f.testName}`);
        if (f.message) lines.push(`  Message: ${f.message}`);
        if (f.trace) lines.push(`  ${truncateOutput(f.trace, 2000)}`);
      }
      if (s.failures.length > cap) {
        lines.push(`\n... and ${s.failures.length - cap} more`);
      }
    }
  } else {
    lines.push(task.status === "completed"
      ? "TESTS PASSED (no structured results found)"
      : "TESTS FAILED");
  }

  if (task.runId) {
    lines.push(
      "",
      `Archived as: ${task.runId}`,
      `   Use test_results({ runId: "${task.runId}" }) to browse.`,
      `   Use test_runs() to see all archived runs.`,
    );
  }

  return lines.join("\n");
}

function formatLintResult(task) {
  const lines = [
    task.status === "completed" ? "LINT CHECKS PASSED" : "LINT CHECKS FAILED",
  ];

  const combined = [task.stdout, task.stderr].filter(Boolean).join("\n");
  if (combined.trim()) {
    lines.push("", "--- Output ---", truncateOutput(combined.trim()));
  }
  return lines.join("\n");
}

function formatCrossSdkTestResult(task) {
  const lines = [];

  if (task.testSummary) {
    const s = task.testSummary;
    const header = task.status === "completed"
      ? "CROSS-SDK TESTS PASSED"
      : "CROSS-SDK TESTS FAILED";
    lines.push(header, "");
    lines.push("--- Summary ---");
    lines.push(
      `Total: ${s.total}  |  Passed: ${s.passed}  |  ` +
      `Failed: ${s.failed}  |  Skipped: ${s.skipped}`
    );

    if (s.failures && s.failures.length > 0) {
      lines.push("", "--- Failure details ---");
      const cap = 25;
      for (const f of s.failures.slice(0, cap)) {
        lines.push(`\n> ${f.className}`);
        if (f.message) lines.push(`  Message: ${f.message}`);
        if (f.trace) lines.push(`  ${truncateOutput(f.trace, 2000)}`);
      }
      if (s.failures.length > cap) {
        lines.push(`\n... and ${s.failures.length - cap} more`);
      }
    }
  } else {
    lines.push(
      task.status === "completed"
        ? "CROSS-SDK TESTS PASSED (no structured results found)"
        : "CROSS-SDK TESTS FAILED"
    );
    const combined = [task.stdout, task.stderr].filter(Boolean).join("\n");
    if (combined.trim()) {
      lines.push("", "--- Output ---", truncateOutput(combined.trim()));
    }
  }

  if (task.runId) {
    lines.push(
      "",
      `Archived as: ${task.runId}`,
      `   Use test_results({ runId: "${task.runId}" }) to browse.`,
      `   Use test_runs() to see all archived runs.`,
    );
  }

  return lines.join("\n");
}

function formatTaskStatus(task) {
  if (!task) return "Task not found.";

  if (task.status === "running") {
    const elapsed = Math.round((nowMs() - task.startTime) / 1000);
    const recentOutput = (task.stdout + task.stderr).slice(-2000);
    return [
      `Task ${task.id} is still running (${elapsed}s elapsed).`,
      "",
      `Type: ${task.type}`,
      "",
      "--- Recent output ---",
      recentOutput || "(no output yet)",
    ].join("\n");
  }

  const elapsed = Math.round((task.endTime - task.startTime) / 1000);
  let resultText;
  if (task.type === "test") {
    resultText = formatTestResult(task);
  } else if (task.type === "lint") {
    resultText = formatLintResult(task);
  } else if (task.type === "cross_sdk_test") {
    resultText = formatCrossSdkTestResult(task);
  } else {
    resultText = `Status: ${task.status}\n\n${(task.stdout + "\n" + task.stderr).trim()}`;
  }

  return `Task ${task.id} completed in ${elapsed}s.\n\n${resultText}`;
}

function formatTestOutputResult(test) {
  const lines = [
    `--- Test Output: ${test.testId} ---`,
    "",
    `Status: ${test.resultType === "SUCCESS" ? "PASSED" : test.resultType === "FAILURE" ? "FAILED" : "SKIPPED"}`,
    `Duration: ${test.duration || 0}ms`,
    "",
  ];

  if (test.failureMessage) {
    lines.push("--- Failure Message ---");
    lines.push(test.failureMessage);
    lines.push("");
  }

  if (test.failureTrace) {
    lines.push("--- Stack Trace ---");
    lines.push(truncateOutput(test.failureTrace, 5000));
    lines.push("");
  }

  if (test.stdout) {
    lines.push("--- stdout ---");
    lines.push(truncateOutput(test.stdout, 5000));
    lines.push("");
  } else {
    lines.push("--- stdout ---");
    lines.push("(no output)");
    lines.push("");
  }

  if (test.stderr) {
    lines.push("--- stderr ---");
    lines.push(truncateOutput(test.stderr, 5000));
  } else {
    lines.push("--- stderr ---");
    lines.push("(no output)");
  }

  return lines.join("\n");
}

// ---------------------------------------------------------------------------
// Workspace validation
// ---------------------------------------------------------------------------
function validateWorkspace() {
  const hasGoMod = existsSync(join(WORKSPACE, "go.mod"));
  if (!hasGoMod) {
    throw new Error(
      `Workspace ${WORKSPACE} does not contain go.mod. Is the volume mounted correctly?`
    );
  }
}

// ---------------------------------------------------------------------------
// MCP Server setup
// ---------------------------------------------------------------------------
const server = new McpServer({
  name: "go-build-tools",
  version: "1.0.0",
});

// --- Tool: test ---
const testSchema = {
  filter: z
    .string()
    .optional()
    .describe("go test -run regex to filter tests by name"),
  path: z
    .string()
    .optional()
    .describe("Package path (default: ./...)"),
  maxFail: z
    .number()
    .optional()
    .describe("Stop after N failures (uses -failfast if 1)"),
};

server.tool(
  "test",
  "Run Go tests asynchronously. Returns a task ID to poll for results.",
  testSchema,
  async (params) => {
    validateWorkspace();

    ensureDir(TEST_RESULTS_DIR);
    if (existsSync(TEST_OUTPUT_DIR)) {
      rmSync(TEST_OUTPUT_DIR, { recursive: true, force: true });
    }
    ensureDir(TEST_OUTPUT_DIR);

    const args = ["test", "-v", "-json", "-count=1", "-p", "1"];
    const testPath = params.path || "./...";

    if (params.filter) {
      args.push("-run", params.filter);
    }
    if (params.maxFail === 1) {
      args.push("-failfast");
    }

    args.push(testPath);

    const taskId = startProcessTask("test", "go", args, TEST_TIMEOUT_MS);

    return {
      content: [
        {
          type: "text",
          text: `Test task started.\n**Task ID:** ${taskId}\n**Command:** go ${args.join(" ")}\n\nUse \`task_status\` to check progress.`,
        },
      ],
    };
  }
);

// --- Tool: lint ---
const lintSchema = {
  tool: z
    .enum(["all", "golangci-lint", "vet"])
    .optional()
    .default("all")
    .describe("Which linter to run: all, golangci-lint, or vet"),
  fix: z
    .boolean()
    .optional()
    .default(false)
    .describe("Auto-fix lint issues (golangci-lint only — applies --fix)"),
};

server.tool(
  "lint",
  "Run linting tools (golangci-lint and/or go vet). Returns a task ID.",
  lintSchema,
  async (params) => {
    validateWorkspace();

    const lintTool = params.tool || "all";
    const fix = params.fix || false;
    let command, args;

    if (lintTool === "vet") {
      command = "go";
      args = ["vet", "./..."];
    } else if (lintTool === "golangci-lint") {
      command = "golangci-lint";
      args = fix ? ["run", "--fix", "./..."] : ["run", "./..."];
    } else {
      // all
      command = "bash";
      const golangciCmd = fix
        ? "golangci-lint run --fix ./..."
        : "golangci-lint run ./...";
      args = [
        "-c",
        `echo '=== go vet ===' && go vet ./... && echo '=== golangci-lint ===' && ${golangciCmd}`,
      ];
    }

    const taskId = startProcessTask("lint", command, args, LINT_TIMEOUT_MS);

    return {
      content: [
        {
          type: "text",
          text: `Lint task started (${lintTool}).\n**Task ID:** ${taskId}\n\nUse \`task_status\` to check progress.`,
        },
      ],
    };
  }
);

// --- Tool: cross_sdk_test ---
const DEFAULT_PARALLEL_JOBS = 4;
const crossSdkTestSchema = {
  grep: z
    .string()
    .optional()
    .describe("Mocha --grep filter expression"),
  timeout: z
    .number()
    .optional()
    .describe("Mocha per-test timeout in ms"),
  parallel: z
    .boolean()
    .optional()
    .describe("Run tests in parallel (Mocha --parallel). Defaults to true"),
  jobs: z
    .number()
    .optional()
    .describe(
      `Number of parallel worker processes (Mocha --jobs). Each worker gets its own test-server instance. Defaults to ${DEFAULT_PARALLEL_JOBS}`
    ),
};

function buildCrossSdkScript(grep, timeout, parallel, jobs) {
  const basePort = 3030;
  const useParallel = parallel !== false;
  const numJobs = useParallel ? (jobs || DEFAULT_PARALLEL_JOBS) : 1;
  const mochaArgs = [
    grep ? `--grep ${JSON.stringify(grep)}` : "",
    timeout ? `--timeout ${timeout}` : "",
    useParallel ? "--parallel" : "",
    useParallel ? `--jobs ${numJobs}` : "",
    useParallel ? "--reporter json" : "",
  ]
    .filter(Boolean)
    .join(" ");

  const ports = Array.from({ length: numJobs }, (_, i) => basePort + i);
  const portsCSV = ports.join(",");

  return `
set -euo pipefail
export TEST_MODE=testing

export SUPERTOKENS_CORE_HOST="\${SUPERTOKENS_CORE_HOST:-localhost}"
export SUPERTOKENS_CORE_PORT="\${SUPERTOKENS_CORE_PORT:-3567}"

SERVER_PIDS=""
cleanup() {
  echo "[cross-sdk] Cleaning up background processes..."
  for pid in $SERVER_PIDS; do
    kill "$pid" 2>/dev/null || true
  done
  wait 2>/dev/null || true
}
trap cleanup EXIT

# Build the Go test server
echo "[cross-sdk] Building Go test server..."
cd /workspace/test/test-server
go build -o /tmp/test-server .

# Start ${numJobs} test-server(s) on ports ${portsCSV}
${ports.map(p => `echo "[cross-sdk] Starting test-server on port ${p}..."
API_PORT=${p} /tmp/test-server &
SERVER_PIDS="$SERVER_PIDS $!"
`).join("")}

# Wait for all test-servers to be ready
echo "[cross-sdk] Waiting for test-server readiness..."
for port in ${portsCSV.replace(/,/g, " ")}; do
  ready=false
  for i in $(seq 1 30); do
    if curl -sf "http://localhost:$port/test/ping" >/dev/null 2>&1; then
      echo "[cross-sdk] Test-server on port $port is ready."
      ready=true
      break
    fi
    sleep 1
  done
  if [ "$ready" != "true" ]; then
    echo "[cross-sdk] ERROR: Test-server on port $port did not become ready in 30s." >&2
    exit 1
  fi
done

# Install and build the cross-SDK test suite
echo "[cross-sdk] Installing cross-SDK test dependencies..."
cd /cross-sdk-tests
npm install 2>&1 | tail -3
npm run build 2>&1 | tail -3

# Run Mocha tests
echo "[cross-sdk] Running Mocha tests (${numJobs} workers, ports: ${portsCSV})..."
export API_PORTS="${portsCSV}"
${useParallel
    ? `set +e
npx mocha ${mochaArgs} > test-results-raw.txt 2>&1
MOCHA_EXIT=$?
set -e

node -e '
  const fs = require("fs");
  const raw = fs.readFileSync("test-results-raw.txt", "utf-8");
  const lines = raw.split("\\n");
  const nonJson = [];
  const jsonLines = [];
  let inJson = false;
  for (const line of lines) {
    if (!inJson && line.startsWith("{")) inJson = true;
    if (inJson) { jsonLines.push(line); }
    else { nonJson.push(line); }
  }
  for (const l of nonJson) { if (l.trim()) console.log(l); }

  const jsonStr = jsonLines.join("\\n");
  fs.writeFileSync("test-results.json", jsonStr);

  try {
    const d = JSON.parse(jsonStr);
    const s = d.stats || {};
    console.log("[cross-sdk] Results: " +
      s.tests + " tests, " +
      s.passes + " passing, " +
      s.failures + " failing, " +
      s.pending + " pending (" +
      (s.duration / 1000).toFixed(1) + "s)");
    for (const f of (d.failures || [])) {
      console.log("  FAIL: " + f.fullTitle);
      if (f.err && f.err.message) console.log("        " + f.err.message.split("\\n")[0]);
    }
  } catch(e) { console.error("[cross-sdk] Could not parse results:", e.message); }
'
exit $MOCHA_EXIT`
    : `npx mocha ${mochaArgs}`}
`;
}

server.tool(
  "cross_sdk_test",
  "Run cross-SDK Mocha tests against the Go test-server. Returns a task ID to poll for results.",
  crossSdkTestSchema,
  async (params) => {
    if (!existsSync(join(CROSS_SDK_DIR, "package.json"))) {
      return {
        content: [
          {
            type: "text",
            text: `Cross-SDK test directory not found. Ensure the backend-sdk-testing repo is mounted at ${CROSS_SDK_DIR}.\n\nSet BACKEND_SDK_TESTING_PATH in your environment or place the repo at ../backend-sdk-testing relative to this project.`,
          },
        ],
        isError: true,
      };
    }

    const script = buildCrossSdkScript(params.grep, params.timeout, params.parallel, params.jobs);
    const taskId = startProcessTask(
      "cross_sdk_test",
      "bash",
      ["-c", script],
      CROSS_SDK_TIMEOUT_MS
    );

    const useParallel = params.parallel !== false;
    const numJobs = useParallel ? (params.jobs || DEFAULT_PARALLEL_JOBS) : 1;
    return {
      content: [
        {
          type: "text",
          text: `Cross-SDK test task started.\n**Task ID:** ${taskId}\n${params.grep ? `**Grep:** ${params.grep}\n` : ""}${params.timeout ? `**Timeout:** ${params.timeout}ms\n` : ""}**Parallel:** ${useParallel} (${numJobs} workers)\n\nUse \`task_status\` to check progress.`,
        },
      ],
    };
  }
);

// --- Tool: task_status ---
const taskStatusSchema = {
  taskId: z.string().describe("The task ID to check"),
};

server.tool(
  "task_status",
  "Check the status of a running or completed task. Returns results when done.",
  taskStatusSchema,
  async (params) => {
    const task = tasks.get(params.taskId);
    if (!task) {
      return {
        content: [{ type: "text", text: `Task not found: ${params.taskId}` }],
        isError: true,
      };
    }
    return {
      content: [{ type: "text", text: formatTaskStatus(task) }],
      isError: task.status === "failed",
    };
  }
);

// --- Tool: task_cancel ---
const taskCancelSchema = {
  taskId: z.string().describe("The task ID to cancel"),
};

server.tool(
  "task_cancel",
  "Cancel a running task.",
  taskCancelSchema,
  async (params) => {
    const task = tasks.get(params.taskId);
    if (!task) {
      return {
        content: [{ type: "text", text: `Task not found: ${params.taskId}` }],
        isError: true,
      };
    }
    if (task.status !== "running") {
      return {
        content: [{ type: "text", text: `Task ${task.id} is not running (status: ${task.status})` }],
      };
    }

    task.status = "cancelled";
    task.endTime = nowMs();
    if (task._timeout) clearTimeout(task._timeout);
    if (task._process) {
      try {
        task._process.kill("SIGTERM");
        setTimeout(() => {
          try {
            task._process.kill("SIGKILL");
          } catch {}
        }, 5000);
      } catch {}
    }

    return {
      content: [{ type: "text", text: `Cancellation signal sent to task ${task.id}` }],
    };
  }
);

// --- Tool: task_list ---
server.tool(
  "task_list",
  "List all tasks (running and recent completed).",
  {},
  async () => {
    pruneOldTasks();
    if (tasks.size === 0) {
      return {
        content: [{ type: "text", text: "No tasks." }],
      };
    }

    const lines = ["--- Tasks ---", ""];
    const sorted = [...tasks.values()].sort((a, b) => b.startTime - a.startTime);

    for (const t of sorted) {
      const elapsed = t.endTime
        ? Math.round((t.endTime - t.startTime) / 1000)
        : Math.round((nowMs() - t.startTime) / 1000);

      const status = t.status === "running" ? "running" :
                     t.status === "completed" ? "completed" :
                     t.status === "cancelled" ? "cancelled" : "failed";

      lines.push(`${t.id}  ${status}  ${t.type}  (${elapsed}s)`);
    }

    return {
      content: [{ type: "text", text: lines.join("\n") }],
    };
  }
);

// --- Tool: test_results ---
const testResultsSchema = {
  runId: z
    .string()
    .optional()
    .describe("Test run ID (default: latest run)"),
  filter: z
    .enum(["all", "failed", "passed", "skipped"])
    .optional()
    .default("all")
    .describe("Filter results by status"),
  className: z
    .string()
    .optional()
    .describe("Filter by test package/name (substring match)"),
};

server.tool(
  "test_results",
  "Browse archived test results with optional filters.",
  testResultsSchema,
  async (params) => {
    const outputs = getArchivedTestOutputs(params.runId);
    if (outputs.length === 0) {
      return {
        content: [{ type: "text", text: "No test output files found. Run tests first using the test tool." }],
      };
    }

    let filtered = outputs;
    const statusFilter = params.filter || "all";
    if (statusFilter === "failed") {
      filtered = filtered.filter((o) => o.resultType === "FAILURE");
    } else if (statusFilter === "passed") {
      filtered = filtered.filter((o) => o.resultType === "SUCCESS");
    } else if (statusFilter === "skipped") {
      filtered = filtered.filter((o) => o.resultType === "SKIPPED");
    }

    if (params.className) {
      const cn = params.className.toLowerCase();
      filtered = filtered.filter(
        (o) => o.className && o.className.toLowerCase().includes(cn)
      );
    }

    const runDir = getRunDir(params.runId);
    const effectiveRunId = runDir ? basename(runDir) : null;

    const lines = [
      `--- Test Results (${filtered.length} tests)${effectiveRunId ? ` [${effectiveRunId}]` : ""} ---`,
      "",
    ];

    const byStatus = { SUCCESS: 0, FAILURE: 0, SKIPPED: 0 };
    filtered.forEach((r) => { byStatus[r.resultType] = (byStatus[r.resultType] || 0) + 1; });
    lines.push(`Passed: ${byStatus.SUCCESS}  |  Failed: ${byStatus.FAILURE}  |  Skipped: ${byStatus.SKIPPED}`);
    lines.push("");

    for (const t of filtered) {
      const icon = t.resultType === "SUCCESS" ? "PASS" : t.resultType === "FAILURE" ? "FAIL" : "SKIP";
      const dur = t.duration ? `${t.duration}ms` : "?";
      const hasOut = (t.stdout || t.stderr) ? " [has output]" : "";
      const msg = t.failureMessage ? ` - ${t.failureMessage.substring(0, 80)}` : "";
      lines.push(`${icon}  ${t.testId} (${dur})${hasOut}${msg}`);
    }

    if (filtered.some((r) => r.resultType === "FAILURE")) {
      lines.push("");
      lines.push('Use test_output({ testId: "<testId>" }) to see stdout/stderr for a specific test.');
    }

    return {
      content: [{ type: "text", text: lines.join("\n") }],
    };
  }
);

// --- Tool: test_output ---
const testOutputSchema = {
  testId: z
    .string()
    .describe("Full test ID (e.g. TestCreateSession)"),
  runId: z
    .string()
    .optional()
    .describe("Test run ID (default: latest run)"),
};

server.tool(
  "test_output",
  "Get detailed stdout/stderr for a specific test.",
  testOutputSchema,
  async (params) => {
    const outputs = getArchivedTestOutputs(params.runId);

    let test = outputs.find((o) => o.testId === params.testId);
    if (!test) {
      test = outputs.find(
        (o) =>
          o.testId?.endsWith(params.testId) ||
          params.testId.endsWith(o.testId) ||
          o.testId?.includes(params.testId)
      );
    }

    if (!test) {
      const candidates = outputs.filter(
        (o) => o.testId?.includes(params.testId) || o.testName === params.testId
      );
      if (candidates.length > 1) {
        return {
          content: [{
            type: "text",
            text: `Multiple tests match "${params.testId}":\n\n${candidates.map(m => `  ${m.testId}`).join("\n")}\n\nPlease specify the full testId.`,
          }],
        };
      }
      return {
        content: [{
          type: "text",
          text: `Test not found: ${params.testId}\n\nUse test_results to see available tests.`,
        }],
        isError: true,
      };
    }

    return {
      content: [{ type: "text", text: formatTestOutputResult(test) }],
    };
  }
);

// --- Tool: test_runs ---
server.tool(
  "test_runs",
  "List archived test runs.",
  {},
  async () => {
    const runs = listTestRuns();
    if (runs.length === 0) {
      return {
        content: [{ type: "text", text: "No archived test runs found. Run tests first using the test tool." }],
      };
    }

    const lines = [
      `--- Archived Test Runs (${runs.length}) ---`,
      "",
    ];

    for (const r of runs) {
      const s = r.testSummary;
      const date = r.startTime
        ? new Date(r.startTime).toISOString().replace("T", " ").replace(/\.\d+Z$/, "")
        : "?";
      const duration = r.startTime && r.endTime
        ? Math.round((r.endTime - r.startTime) / 1000)
        : "?";
      const icon = (s?.failed > 0 || s?.errors > 0) ? "FAIL" : "PASS";
      const label = r.taskType === "cross_sdk_test" ? " [cross-sdk]" : r.taskType === "test" ? " [go-test]" : "";

      lines.push(`${icon}  ${r.runId}${label}`);
      lines.push(`   ${date}  |  ${duration}s  |  ${s?.total || "?"} tests: ${s?.passed || 0} passed, ${s?.failed || 0} failed, ${s?.errors || 0} errors, ${s?.skipped || 0} skipped`);
      lines.push("");
    }

    lines.push('Use test_results({ runId: "<runId>" }) to see results for a specific run.');

    return {
      content: [{ type: "text", text: lines.join("\n") }],
    };
  }
);

// ---------------------------------------------------------------------------
// HTTP server (SSE transport + stateless API)
// ---------------------------------------------------------------------------
function buildToolSchemas() {
  const tools = [
    { name: "test", schema: testSchema, description: "Run Go tests" },
    { name: "lint", schema: lintSchema, description: "Run linting" },
    {
      name: "cross_sdk_test",
      schema: crossSdkTestSchema,
      description: "Run cross-SDK Mocha tests",
    },
    {
      name: "task_status",
      schema: taskStatusSchema,
      description: "Check task status",
    },
    {
      name: "task_cancel",
      schema: taskCancelSchema,
      description: "Cancel a task",
    },
    { name: "task_list", schema: {}, description: "List all tasks" },
    {
      name: "test_results",
      schema: testResultsSchema,
      description: "Browse test results",
    },
    {
      name: "test_output",
      schema: testOutputSchema,
      description: "Get per-test output",
    },
    { name: "test_runs", schema: {}, description: "List test runs" },
  ];

  return tools.map((t) => {
    const zObj = z.object(
      Object.fromEntries(
        Object.entries(t.schema).map(([k, v]) => [k, v])
      )
    );
    return {
      name: t.name,
      description: t.description,
      inputSchema: zodToJsonSchema(zObj),
    };
  });
}

// Stateless API handlers
const statelessHandlers = {
  test: async (a) => {
    validateWorkspace();
    ensureDir(TEST_RESULTS_DIR);
    if (existsSync(TEST_OUTPUT_DIR)) {
      rmSync(TEST_OUTPUT_DIR, { recursive: true, force: true });
    }
    ensureDir(TEST_OUTPUT_DIR);

    const args = ["test", "-v", "-json", "-count=1", "-p", "1"];
    if (a.filter) args.push("-run", a.filter);
    if (a.maxFail === 1) args.push("-failfast");
    args.push(a.path || "./...");

    const taskId = startProcessTask("test", "go", args, TEST_TIMEOUT_MS);
    return {
      content: [{ type: "text", text: `Test task started.\nTask ID: ${taskId}` }],
    };
  },
  lint: async (a) => {
    validateWorkspace();
    const lintTool = a.tool || "all";
    const fix = a.fix || false;
    let cmd, cargs;
    if (lintTool === "vet") {
      cmd = "go";
      cargs = ["vet", "./..."];
    } else if (lintTool === "golangci-lint") {
      cmd = "golangci-lint";
      cargs = fix ? ["run", "--fix", "./..."] : ["run", "./..."];
    } else {
      cmd = "bash";
      const golangciCmd = fix
        ? "golangci-lint run --fix ./..."
        : "golangci-lint run ./...";
      cargs = [
        "-c",
        `echo '=== go vet ===' && go vet ./... && echo '=== golangci-lint ===' && ${golangciCmd}`,
      ];
    }
    const taskId = startProcessTask("lint", cmd, cargs, LINT_TIMEOUT_MS);
    return {
      content: [{ type: "text", text: `Lint task started (${lintTool}${fix ? ", fix" : ""}).\nTask ID: ${taskId}` }],
    };
  },
  cross_sdk_test: async (a) => {
    if (!existsSync(join(CROSS_SDK_DIR, "package.json"))) {
      return {
        content: [{ type: "text", text: `Cross-SDK test directory not found at ${CROSS_SDK_DIR}.` }],
        isError: true,
      };
    }
    const script = buildCrossSdkScript(a.grep, a.timeout, a.parallel, a.jobs);
    const taskId = startProcessTask("cross_sdk_test", "bash", ["-c", script], CROSS_SDK_TIMEOUT_MS);
    return {
      content: [{ type: "text", text: `Cross-SDK test task started.\nTask ID: ${taskId}` }],
    };
  },
  task_status: async (a) => {
    const task = tasks.get(a.taskId);
    if (!task) {
      return {
        content: [{ type: "text", text: `Task not found: ${a.taskId}` }],
        isError: true,
      };
    }
    return {
      content: [{ type: "text", text: formatTaskStatus(task) }],
      isError: task.status === "failed",
    };
  },
  task_cancel: async (a) => {
    const task = tasks.get(a.taskId);
    if (!task) {
      return {
        content: [{ type: "text", text: `Task not found: ${a.taskId}` }],
        isError: true,
      };
    }
    if (task.status !== "running") {
      return {
        content: [{ type: "text", text: `Task ${task.id} is not running (status: ${task.status})` }],
      };
    }
    task.status = "cancelled";
    task.endTime = nowMs();
    if (task._timeout) clearTimeout(task._timeout);
    if (task._process) {
      try { task._process.kill("SIGTERM"); } catch {}
    }
    return {
      content: [{ type: "text", text: `Cancellation signal sent to task ${task.id}` }],
    };
  },
  task_list: async () => {
    pruneOldTasks();
    if (tasks.size === 0) {
      return { content: [{ type: "text", text: "No tasks." }] };
    }
    const lines = ["--- Tasks ---", ""];
    const sorted = [...tasks.values()].sort((a, b) => b.startTime - a.startTime);
    for (const t of sorted) {
      const elapsed = t.endTime
        ? Math.round((t.endTime - t.startTime) / 1000)
        : Math.round((nowMs() - t.startTime) / 1000);
      const status = t.status === "running" ? "running" :
                     t.status === "completed" ? "completed" :
                     t.status === "cancelled" ? "cancelled" : "failed";
      lines.push(`${t.id}  ${status}  ${t.type}  (${elapsed}s)`);
    }
    return { content: [{ type: "text", text: lines.join("\n") }] };
  },
  test_results: async (a) => {
    const outputs = getArchivedTestOutputs(a.runId);
    if (outputs.length === 0) {
      return { content: [{ type: "text", text: "No test output files found. Run tests first using the test tool." }] };
    }
    let filtered = outputs;
    if (a.filter === "failed") filtered = filtered.filter((o) => o.resultType === "FAILURE");
    else if (a.filter === "passed") filtered = filtered.filter((o) => o.resultType === "SUCCESS");
    else if (a.filter === "skipped") filtered = filtered.filter((o) => o.resultType === "SKIPPED");
    if (a.className) {
      const cn = a.className.toLowerCase();
      filtered = filtered.filter((o) => o.className && o.className.toLowerCase().includes(cn));
    }
    const runDir = getRunDir(a.runId);
    const effectiveRunId = runDir ? basename(runDir) : null;
    const lines = [
      `--- Test Results (${filtered.length} tests)${effectiveRunId ? ` [${effectiveRunId}]` : ""} ---`,
      "",
    ];
    const byStatus = { SUCCESS: 0, FAILURE: 0, SKIPPED: 0 };
    filtered.forEach((r) => { byStatus[r.resultType] = (byStatus[r.resultType] || 0) + 1; });
    lines.push(`Passed: ${byStatus.SUCCESS}  |  Failed: ${byStatus.FAILURE}  |  Skipped: ${byStatus.SKIPPED}`);
    lines.push("");
    for (const t of filtered) {
      const icon = t.resultType === "SUCCESS" ? "PASS" : t.resultType === "FAILURE" ? "FAIL" : "SKIP";
      const dur = t.duration ? `${t.duration}ms` : "?";
      const hasOut = (t.stdout || t.stderr) ? " [has output]" : "";
      const msg = t.failureMessage ? ` - ${t.failureMessage.substring(0, 80)}` : "";
      lines.push(`${icon}  ${t.testId} (${dur})${hasOut}${msg}`);
    }
    if (filtered.some((r) => r.resultType === "FAILURE")) {
      lines.push("");
      lines.push('Use test_output({ testId: "<testId>" }) to see stdout/stderr for a specific test.');
    }
    return { content: [{ type: "text", text: lines.join("\n") }] };
  },
  test_output: async (a) => {
    const outputs = getArchivedTestOutputs(a.runId);
    let test = outputs.find((o) => o.testId === a.testId);
    if (!test) {
      test = outputs.find((o) =>
        o.testId?.endsWith(a.testId) || a.testId.endsWith(o.testId) || o.testId?.includes(a.testId)
      );
    }
    if (!test) {
      return {
        content: [{ type: "text", text: `Test not found: ${a.testId}\n\nUse test_results to see available tests.` }],
        isError: true,
      };
    }
    return { content: [{ type: "text", text: formatTestOutputResult(test) }] };
  },
  test_runs: async () => {
    const runs = listTestRuns();
    if (runs.length === 0) {
      return { content: [{ type: "text", text: "No archived test runs found. Run tests first using the test tool." }] };
    }
    const lines = [`--- Archived Test Runs (${runs.length}) ---`, ""];
    for (const r of runs) {
      const s = r.testSummary;
      const date = r.startTime ? new Date(r.startTime).toISOString().replace("T", " ").replace(/\.\d+Z$/, "") : "?";
      const duration = r.startTime && r.endTime ? Math.round((r.endTime - r.startTime) / 1000) : "?";
      const icon = (s?.failed > 0 || s?.errors > 0) ? "FAIL" : "PASS";
      const label = r.taskType === "cross_sdk_test" ? " [cross-sdk]" : r.taskType === "test" ? " [go-test]" : "";
      lines.push(`${icon}  ${r.runId}${label}`);
      lines.push(`   ${date}  |  ${duration}s  |  ${s?.total || "?"} tests: ${s?.passed || 0} passed, ${s?.failed || 0} failed, ${s?.errors || 0} errors, ${s?.skipped || 0} skipped`);
      lines.push("");
    }
    lines.push('Use test_results({ runId: "<runId>" }) to see results for a specific run.');
    return { content: [{ type: "text", text: lines.join("\n") }] };
  },
};

async function handleToolCall(toolName, args) {
  const handler = statelessHandlers[toolName];
  if (!handler) {
    return {
      content: [{ type: "text", text: `Unknown tool: ${toolName}` }],
      isError: true,
    };
  }
  return handler(args || {});
}

// ---------------------------------------------------------------------------
// Start server
// ---------------------------------------------------------------------------
async function main() {
  console.error(`[go-build-tools] workspace=${WORKSPACE}`);
  console.error(`[go-build-tools] transport=${TRANSPORT}, port=${PORT}`);

  try {
    validateWorkspace();
    console.error("[go-build-tools] Workspace validated OK");
  } catch (e) {
    console.error(`[go-build-tools] WARNING: ${e.message}`);
  }

  if (TRANSPORT === "stdio") {
    const transport = new StdioServerTransport();
    await server.connect(transport);
    console.error("[go-build-tools] Connected via stdio");
    return;
  }

  // SSE transport with stateless API
  let sseTransport = null;

  const httpServer = createServer(async (req, res) => {
    const url = new URL(req.url, `http://localhost:${PORT}`);

    // CORS
    res.setHeader("Access-Control-Allow-Origin", "*");
    res.setHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS");
    res.setHeader("Access-Control-Allow-Headers", "Content-Type");

    if (req.method === "OPTIONS") {
      res.writeHead(204);
      res.end();
      return;
    }

    // Health check
    if (url.pathname === "/health" && req.method === "GET") {
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ status: "ok", transport: "sse" }));
      return;
    }

    // SSE endpoint
    if (url.pathname === "/sse" && req.method === "GET") {
      sseTransport = new SSEServerTransport("/messages", res);
      await server.connect(sseTransport);
      console.error("[go-build-tools] SSE client connected");
      return;
    }

    // MCP messages endpoint
    if (url.pathname === "/messages" && req.method === "POST") {
      if (!sseTransport) {
        res.writeHead(400, { "Content-Type": "application/json" });
        res.end(JSON.stringify({ error: "No SSE connection. Connect to /sse first." }));
        return;
      }
      await sseTransport.handlePostMessage(req, res);
      return;
    }

    // Stateless API: list tools
    if (url.pathname === "/api/tools" && req.method === "GET") {
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify(buildToolSchemas(), null, 2));
      return;
    }

    // Stateless API: call tool
    if (url.pathname === "/api/call" && req.method === "POST") {
      let body = "";
      for await (const chunk of req) body += chunk;
      try {
        const parsed = JSON.parse(body);
        const toolName = parsed.tool;
        const toolArgs = parsed.arguments || parsed.args;
        const result = await handleToolCall(toolName, toolArgs);
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(JSON.stringify(result, null, 2));
      } catch (e) {
        res.writeHead(400, { "Content-Type": "application/json" });
        res.end(JSON.stringify({ error: e.message }));
      }
      return;
    }

    // 404
    res.writeHead(404, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ error: "Not found" }));
  });

  httpServer.listen(PORT, "0.0.0.0", () => {
    console.error(`[go-build-tools] HTTP server listening on port ${PORT}`);
    console.error(`[go-build-tools] SSE endpoint: http://localhost:${PORT}/sse`);
    console.error(`[go-build-tools] API endpoint: http://localhost:${PORT}/api/tools`);
  });
}

main().catch((err) => {
  console.error("Fatal error:", err);
  process.exit(1);
});
