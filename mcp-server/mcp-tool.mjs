#!/usr/bin/env node

// Stdio MCP server that proxies to the Docker container's stateless HTTP API.
// Configure in Claude Code's MCP settings to use the build tools directly.
//
// Usage in .mcp.json:
// {
//   "mcpServers": {
//     "go-build-tools": {
//       "command": "node",
//       "args": ["mcp-server/mcp-tool.mjs"],
//       "env": { "MCP_PORT": "3001" }
//     }
//   }
// }

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import http from "node:http";
import { readFileSync } from "node:fs";
import { resolve, dirname } from "node:path";
import { fileURLToPath } from "node:url";

function getDefaultPort() {
  try {
    const __dirname = dirname(fileURLToPath(import.meta.url));
    const content = readFileSync(resolve(__dirname, "../mcp.env"), "utf8");
    const match = content.match(/^MCP_PORT=(\d+)/m);
    if (match) return match[1];
  } catch {}
  return "3001";
}

const API_BASE = process.env.MCP_URL
  || `http://localhost:${process.env.MCP_PORT || getDefaultPort()}`;

function httpPost(path, body) {
  return new Promise((resolve, reject) => {
    const url = new URL(path, API_BASE);
    const payload = JSON.stringify(body);

    const req = http.request(
      {
        hostname: url.hostname,
        port: url.port || 80,
        path: url.pathname,
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Content-Length": Buffer.byteLength(payload),
        },
      },
      (res) => {
        let data = "";
        res.on("data", (chunk) => { data += chunk; });
        res.on("end", () => {
          if (res.statusCode >= 200 && res.statusCode < 300) {
            try { resolve(JSON.parse(data)); } catch { resolve(data); }
          } else {
            reject(new Error(`HTTP ${res.statusCode}: ${data}`));
          }
        });
      }
    );
    req.on("error", reject);
    req.setTimeout(0);
    req.write(payload);
    req.end();
  });
}

async function callRemoteTool(name, args) {
  try {
    const result = await httpPost("/api/call", { tool: name, arguments: args });
    if (result?.content) return result;
    return { content: [{ type: "text", text: JSON.stringify(result, null, 2) }] };
  } catch (err) {
    return {
      content: [{ type: "text", text: `Error calling ${name}: ${err.message}` }],
      isError: true,
    };
  }
}

const server = new McpServer({
  name: "go-build-tools",
  version: "1.0.0",
});

// --- test ---
server.tool(
  "test",
  "Run Go tests in the Docker container. Returns a task ID — use task_status to poll for results.",
  {
    filter: z.string().optional().describe("go test -run regex to filter tests by name"),
    path: z.string().optional().describe("Package path (default: ./...)"),
    maxFail: z.number().optional().describe("Stop after N failures (uses -failfast if 1)"),
  },
  async (params) => callRemoteTool("test", params)
);

// --- lint ---
server.tool(
  "lint",
  "Run linting tools (golangci-lint and/or go vet) in the Docker container. Returns a task ID.",
  {
    tool: z.enum(["all", "golangci-lint", "vet"]).optional().default("all")
      .describe("Which linter to run: all, golangci-lint, or vet"),
    fix: z.boolean().optional().default(false)
      .describe("Auto-fix lint issues (golangci-lint only — applies --fix)"),
  },
  async (params) => callRemoteTool("lint", params)
);

// --- cross_sdk_test ---
server.tool(
  "cross_sdk_test",
  "Run cross-SDK Mocha tests against the Go test-server in Docker. Returns a task ID — use task_status to poll for results.",
  {
    grep: z.string().optional().describe("Mocha --grep filter expression"),
    timeout: z.number().optional().describe("Mocha per-test timeout in ms"),
  },
  async (params) => callRemoteTool("cross_sdk_test", params)
);

// --- task_status ---
server.tool(
  "task_status",
  "Check the status of a running or completed task. Returns results when done.",
  {
    taskId: z.string().describe("The task ID returned by test or lint"),
  },
  async (params) => callRemoteTool("task_status", params)
);

// --- task_cancel ---
server.tool(
  "task_cancel",
  "Cancel a running task.",
  {
    taskId: z.string().describe("The task ID to cancel"),
  },
  async (params) => callRemoteTool("task_cancel", params)
);

// --- task_list ---
server.tool(
  "task_list",
  "List all tasks (running and recent completed).",
  {},
  async () => callRemoteTool("task_list", {})
);

// --- test_results ---
server.tool(
  "test_results",
  "Browse archived test results with optional filters. Use test_output to get detailed output for specific tests.",
  {
    runId: z.string().optional().describe("Test run ID (default: latest run). Use test_runs to list available runs."),
    filter: z.enum(["all", "failed", "passed", "skipped"]).optional().default("all")
      .describe("Filter results by status"),
    className: z.string().optional().describe("Filter by test package/name (substring match)"),
  },
  async (params) => callRemoteTool("test_results", params)
);

// --- test_output ---
server.tool(
  "test_output",
  "Get detailed stdout/stderr for a specific test. Use test_results first to find the testId.",
  {
    testId: z.string().describe("Full test ID (e.g. TestCreateSession)"),
    runId: z.string().optional().describe("Test run ID (default: latest run)"),
  },
  async (params) => callRemoteTool("test_output", params)
);

// --- test_runs ---
server.tool(
  "test_runs",
  "List archived test runs. Use the runId with test_results or test_output to access specific runs.",
  {},
  async () => callRemoteTool("test_runs", {})
);

// --- Start ---
const transport = new StdioServerTransport();
await server.connect(transport);
