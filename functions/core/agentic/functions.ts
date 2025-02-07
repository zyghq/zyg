import { google } from "@ai-sdk/google";
import { generateObject } from "ai";
import { ThreadSummary, threadSummarySchema } from "./schemas.ts";
import { ZYG_SRV_BASE_URL } from "./config.ts";

const model = google("gemini-2.0-flash-lite-preview-02-05");

export async function threadSummarizeFn(
  systemPrompt: string,
  userPrompt: string,
): Promise<ThreadSummary> {
  const { object } = await generateObject<ThreadSummary>({
    model,
    system: systemPrompt,
    prompt: userPrompt,
    schema: threadSummarySchema,
  });

  return object;
}

export async function upsertThreadSummaryRPC() {
  const payload = {
    jsonrpc: "2.0",
    method: "upsertThreadSummary",
    params: {},
    id: Date.now(),
  };

  const response = await fetch(`${ZYG_SRV_BASE_URL}/v1/rpc/threads/`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  const result = await response.json();
  if (result.error) {
    throw new Error(result.error.message);
  }

  return result.result;
}
