import { z } from "zod";

const threadSummarySchema = z.object({
  summary: z.object({
    bulletPoints: z.array(z.string()).describe(
      "4-8 bullet points that summarize the key points of the conversation",
    ),
    oneLineSummary: z.string().describe(
      "a concise 1-line summary of the conversation thread",
    ),
  }).describe("summary of the conversation thread"),
});

export { threadSummarySchema };

export type ThreadSummary = z.infer<typeof threadSummarySchema>;
