import { google } from "@ai-sdk/google";
import { generateObject } from "ai";
import { ThreadSummary, threadSummarySchema } from "./schemas.ts";

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
