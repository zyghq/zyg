/**
 * Entry point for the core server.
 * This script initializes and starts the server that handles Restate endpoints.
 * It:
 * 1. Imports required Restate SDK and service definitions
 * 2. Creates a bidirectional endpoint handler for the thread service
 * 3. Starts a Deno server on port 9080 to handle incoming requests
 */

import "dotenv/config"; // first import to load envs
import * as restate from "@restatedev/restate-sdk/fetch";
import { thread } from "./services/llm.ts";
import { inSync, sync } from "./services/db.ts";
import { threadCreateGitHubIssueWorflow } from "./experimental/workflows.ts";

const handler = restate
  .endpoint()
  .bind(inSync)
  .bind(sync)
  .bind(thread)
  .bind(threadCreateGitHubIssueWorflow)
  .bidirectional()
  .handler();

Deno.serve({ port: 9080 }, handler.fetch);
