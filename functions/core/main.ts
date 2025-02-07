import * as restate from "@restatedev/restate-sdk/fetch";
import { threadSummarizeFn } from "./agentic/functions.ts";
import { threadSummarizeSystemPrompt } from "./agentic/prompts.ts";
import { ThreadSummary } from "./agentic/schemas.ts";
import { LLM_RETRY_CONFIG } from "./agentic/config.ts";

interface ThreadSummaryRequest {
  system?: string;
  prompt: string;
}

const thread = restate.service({
  name: "thread",
  handlers: {
    summarize: async (ctx: restate.Context, req: ThreadSummaryRequest) => {
      const systemPrompt = req.system || threadSummarizeSystemPrompt();
      const { summary } = await ctx.run<ThreadSummary>(
        "llm",
        () => threadSummarizeFn(systemPrompt, req.prompt),
        LLM_RETRY_CONFIG,
      );
      return summary;
    },
  },
});

const handler = restate
  .endpoint()
  .bind(thread)
  .bidirectional()
  .handler();

Deno.serve({ port: 9080 }, handler.fetch);
