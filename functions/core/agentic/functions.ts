import { google } from "@ai-sdk/google";
import { generateObject } from "ai";
import {
  ThreadSummary,
  threadSummarySchema,
  GitHubIssue,
  githubIssueSchema,
} from "./schemas.ts";

const model = google("gemini-2.0-flash-lite-preview-02-05");

export async function threadSummarizeFn(
  systemPrompt: string,
  userPrompt: string
): Promise<ThreadSummary> {
  const { object } = await generateObject<ThreadSummary>({
    model,
    system: systemPrompt,
    prompt: userPrompt,
    schema: threadSummarySchema,
  });

  return object;
}

// experimental; working on it.
export async function expThreadCreateGithubIssueFn(
  sysmtemPrompt: string,
  userPrompt: string
) {
  const { object } = await generateObject<GitHubIssue>({
    model,
    system: sysmtemPrompt,
    prompt: userPrompt,
    schema: githubIssueSchema,
  });

  return object;
}
