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
}

async function getThreadMessages(threadId: string): Promise<ThreadMessage[]> {
  console.log("fetching messages for threadId", threadId);
  // Simulate network delay
  await new Promise((resolve) => setTimeout(resolve, 100));
  const mockMessages: ThreadMessage[] = [
    {
      messageId: "msg_001",
      content:
        "Hello! Thanks for reaching out to Acme Support. My name is Sarah. How may I assist you today?",
    },
    {
      messageId: "msg_002",
      content:
        "Hi Sarah, I'm trying to deploy my application to production but keep getting a 503 error after the build succeeds. I've checked the basic configuration but can't figure out what's wrong.",
    },
    {
      messageId: "msg_003",
      content:
        "I understand how frustrating deployment issues can be. Let me help you investigate this. Could you please provide your deployment ID or the last 4 digits of your project ID?",
    },
    {
      messageId: "msg_004",
      content:
        "Sure, my project ID ends in 8472. The deployment was attempted about 20 minutes ago.",
    },
    {
      messageId: "msg_005",
      content:
        "Thank you for providing that information. I'm looking at your logs now. I can see that your application is failing the health check because it's unable to connect to the database. Are you sure the DATABASE_URL environment variable is properly configured in your production environment?",
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
        **Customer Impact:** Briefly describe the impact on the customer (e.g., "Unable to access paid features", "Workaround required", "Significant inconvenience"). Quantify the impact if possible.
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

function sendHumanApproval(
  requestId: string,
  result: { title: string; body: string }
) {
  console.clear();
  console.log("=== HUMAN APPROVAL REVIEW ===");
  console.log(`Request ID: ${requestId}`);
  console.log("GitHub Title:", result.title);
  console.log("GitHub Body:", result.body);
  console.log("Approve or Deny?");
  console.log("=============================");
}

function approvalApproved(requestId: string) {
  console.clear();
  console.log("=== HUMAN APPROVAL APPROVED ===");
  console.log(`Request ID: ${requestId}`);
  console.log("=============================");
}

function approvalDenied(requestId: string) {
  console.clear();
  console.log("=== HUMAN APPROVAL DENIED ===");
  console.log(`Request ID: ${requestId}`);
  console.log("=============================");
}

function approvalCancelled(requestId: string) {
  console.clear();
  console.log("=== HUMAN APPROVAL CANCELLED ===");
  console.log(`Request ID: ${requestId}`);
  console.log("=============================");
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

      // generate request ID required to identify the approval
      const requestId = ctx.rand.uuidv4();
      await ctx.run(() => sendHumanApproval(requestId, result));

      // Wait until the user approves or denies
      // Promise gets resolved or rejected
      const approval = await ctx.promise<string>("approval");
      if (approval === "approved") {
        await ctx.run(() => approvalApproved(requestId));
      } else if (approval === "denied") {
        await ctx.run(() => approvalDenied(requestId));
      } else if (approval === "cancelled") {
        await ctx.run(() => approvalCancelled(requestId));
      } else {
        throw new restate.TerminalError(
          `Worflow terminated for unknown approval for session ID: ${sessionId}`
        );
      }
      return {
        result,
        requestId,
        approval,
      };
    },
    approval: async (
      ctx: restate.WorkflowSharedContext,
      request: { approval: string; requestId: string }
    ) => {
      // Send data to workflow via durable promise
      ctx.console.log(`approval for request ID: ${request.requestId}`);
      await ctx.promise<string>("approval").resolve(request.approval);
    },
  },
});
