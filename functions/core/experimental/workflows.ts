import * as restate from "@restatedev/restate-sdk/fetch";
import Handlebars from "handlebars";
import { threadCreateGitHubIssueSystemPrompt } from "../agentic/prompts.ts";
import { GitHubIssue } from "../agentic/schemas.ts";
import { expThreadCreateGithubIssueFn } from "../agentic/functions.ts";
import { LLM_RETRY_CONFIG } from "../config.ts";

interface ThreadCreateGitHubIssueReq {
  workspaceId: string;
  threadId: string;
}

interface ThreadMessage {
  messageId: string;
  content: string;
  timestamp: number;
}

async function getThreadMessages(threadId: string): Promise<ThreadMessage[]> {
  console.log("fetching messages for threadId", threadId);
  // Simulate network delay
  await new Promise((resolve) => setTimeout(resolve, 100));

  // Mock conversation messages
  const mockMessages: ThreadMessage[] = [
    {
      messageId: "msg_" + Math.random().toString(36).substring(7),
      content: "Hey, how can I help you today?",
      timestamp: Date.now() - 300000,
    },
    {
      messageId: "msg_" + Math.random().toString(36).substring(7),
      content: "I'm having trouble with my deployment.",
      timestamp: Date.now() - 200000,
    },
    {
      messageId: "msg_" + Math.random().toString(36).substring(7),
      content: "Could you check the logs and help me debug?",
      timestamp: Date.now() - 100000,
    },
  ];
  return mockMessages;
}

function makeGitHubIssuePrompt(messages: ThreadMessage[]) {
  const template = `
        <INSTRUCTION>
        **Description**: Expand on the summary with details from the conversation. Use bullet points or numbered lists for clarity if multiple steps are involved.
        **Steps to Reproduce (If Possible)**:  If the customer provided steps to reproduce the issue, list them clearly and concisely. This is crucial for debugging.
        **Expected Result:** What should have happened.
        **Actual Result:** What actually happened.
        **Environment:** Include any relevant environment information provided by the customer (e.g., browser, operating system, device, app version).
        **Customer Impact:** Briefly describe the impact on the customer (e.g., "Unable to access paid features", "Workaround required", "Significant inconvenience").  Quantify the impact if possible.
        </INSTRUCTION>
        
        <CONVERSATION>
        {{#each messages}}
        <MESSAGE {{@index}}>
            {{this.content}}
        </MESSAGE>
        {{/each}}
        </CONVERSATION>
    `;

  const compiled = Handlebars.compile(template);
  return compiled({ messages });
}

export const threadCreateGitHubIssueWorflow = restate.workflow({
  name: "threadCreateGithubIssue",
  handlers: {
    run: async (
      ctx: restate.WorkflowContext,
      req: ThreadCreateGitHubIssueReq
    ) => {
      // workflow ID == session ID; workflows runs per session.
      // workflows are executed within session.
      const sessionId = ctx.key;
      ctx.console.log("invoking for sesssion ID" + sessionId);
      const messages = await ctx.run("read messages", async () => {
        return await getThreadMessages(req.threadId);
      });
      const systemPrompt = threadCreateGitHubIssueSystemPrompt();
      const userPrompt = makeGitHubIssuePrompt(messages);
      const result = await ctx.run<GitHubIssue>(
        "llm create github issue",
        () => expThreadCreateGithubIssueFn(systemPrompt, userPrompt),
        LLM_RETRY_CONFIG
      );
      console.log("************** result ********************");
      console.log(result);
      return result;
    },
  },
});
