import { z } from "zod";

const threadSummarySchema = z.object({
  summary: z
    .object({
      bulletPoints: z
        .array(z.string())
        .describe(
          "4-8 bullet points that summarize the key points of the conversation"
        ),
      oneLineSummary: z
        .string()
        .describe("a concise 1-line summary of the conversation thread"),
    })
    .describe("summary of the conversation thread"),
});

const githubIssueSchema = z
  .object({
    title: z
      .string()
      .describe(
        `A concise title summarizing the customer's problem (e.g., "Login fails after password reset")`
      ),
    body: z
      .string()
      .describe(
        `Concise GitHub Markdown description: ` +
          `problem summary, goal, steps, expected/actual results, error messages in code blocks. Omit irrelevant details.`
      ),
  })
  .describe(
    `Schema for representing GitHub issues created from customer support conversations. ` +
      `Captures the key information needed for engineering teams to understand and resolve reported problems efficiently.`
  );

export { threadSummarySchema, githubIssueSchema };

export type ThreadSummary = z.infer<typeof threadSummarySchema>;
export type GitHubIssue = z.infer<typeof githubIssueSchema>;
