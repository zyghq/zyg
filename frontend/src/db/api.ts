import { z } from "zod";

import {
  workspaceResponseSchema,
  workspaceMetricsResponseSchema,
  threadChatResponseSchema,
  threadChatMessagePreviewSchema,
  accountSchema,
} from "./schema";
import {
  IWorkspaceEntities,
  ThreadChatStoreType,
  ThreadChatMapStoreType,
} from "./store";

// TODO: fix this.
// const TOKEN = `eyJhbGciOiJIUzI1NiIsImtpZCI6InhIaUQ3b1NRUSt5NWJ6ZUIiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJhdXRoZW50aWNhdGVkIiwiZXhwIjoxNzE2ODE2MjkwLCJpYXQiOjE3MTY3ODc0OTAsImlzcyI6Imh0dHBzOi8vdHdqbGxqdmltb21nb29jcm5pZWwuc3VwYWJhc2UuY28vYXV0aC92MSIsInN1YiI6IjYyMmZkNDFjLWE0MzctNGNiMi1iMTNkLThjOTdjZmRjNTYyZiIsImVtYWlsIjoic2FuY2hpdHJya0BnbWFpbC5jb20iLCJwaG9uZSI6IiIsImFwcF9tZXRhZGF0YSI6eyJwcm92aWRlciI6ImVtYWlsIiwicHJvdmlkZXJzIjpbImVtYWlsIl19LCJ1c2VyX21ldGFkYXRhIjp7fSwicm9sZSI6ImF1dGhlbnRpY2F0ZWQiLCJhYWwiOiJhYWwxIiwiYW1yIjpbeyJtZXRob2QiOiJwYXNzd29yZCIsInRpbWVzdGFtcCI6MTcxNTcwNTM0N31dLCJzZXNzaW9uX2lkIjoiM2U0Yjk2ZmUtNjk1MS00ZTI4LThiZWUtYzUyODczYmZjMmRjIiwiaXNfYW5vbnltb3VzIjpmYWxzZX0.zf9UuXjc2Reo4FrDJB-2ZUXz6QDBKZ4WKl7ic4llAdU`;

export type WorkspaceResponseType = z.infer<typeof workspaceResponseSchema>;

export type WorkspaceMetricsResponseType = z.infer<
  typeof workspaceMetricsResponseSchema
>;

export type ThreadChatResponseType = z.infer<typeof threadChatResponseSchema>;

export type ThreadChatMessagePreviewType = z.infer<
  typeof threadChatMessagePreviewSchema
>;

export type AccountType = z.infer<typeof accountSchema>;

function initialWorkspaceData(): IWorkspaceEntities {
  return {
    hasData: false,
    isPending: true,
    error: null,
    workspace: null,
    metrics: {
      active: 0,
      done: 0,
      snoozed: 0,
      assignedToMe: 0,
      unassigned: 0,
      otherAssigned: 0,
      labels: [],
    },
    threadChats: null,
  };
}

// API call to fetch workspace details
async function getWorkspace(
  token: string,
  workspaceId: string
): Promise<{ data: WorkspaceResponseType | null; error: Error | null }> {
  try {
    // make a request
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/`,
      {
        method: "GET",
        headers: {
          // "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `error fetching workspace details: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    const data = await response.json();
    // parse into schema
    try {
      const workspace = workspaceResponseSchema.parse({ ...data });
      return { error: null, data: workspace };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      }
      console.error(err);
      return {
        error: new Error("error parsing workspace schema"),
        data: null,
      };
    }
  } catch (error) {
    console.error(error);
    return {
      error: new Error(
        "error fetching workspace details - something went wrong"
      ),
      data: null,
    };
  }
}

// API call to fetch workspace metrics
async function getWorkspaceMetrics(
  token: string,
  workspaceId: string
): Promise<{ data: WorkspaceMetricsResponseType | null; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/chat/metrics/`,
      {
        method: "GET",
        headers: {
          // "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `error fetching workspace metrics: ${status} ${statusText}`
        ),
        data: null,
      };
    }
    const data = await response.json();
    try {
      const metrics = workspaceMetricsResponseSchema.parse({ ...data });
      return { error: null, data: metrics };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      }
      console.error(err);
      return {
        error: new Error("error parsing workspace metrics schema"),
        data: null,
      };
    }
  } catch (error) {
    console.error(error);
    return {
      error: new Error(
        "error fetching workspace metrics - something went wrong"
      ),
      data: null,
    };
  }
}

// API call to fetch thread chats
async function getWorkspaceThreads(
  token: string,
  workspaceId: string
): Promise<{ data: ThreadChatResponseType[] | null; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/chat/`,
      {
        method: "GET",
        headers: {
          // "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `error fetching workspace threads: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    try {
      const data = await response.json();
      // schema validate for each item
      const threads = data.map((item: any) => {
        return threadChatResponseSchema.parse({ ...item });
      });
      return { error: null, data: threads };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      }
      console.error(err);
      return {
        error: new Error("error parsing workspace threads schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error(
        "error fetching workspace threads - something went wrong"
      ),
      data: null,
    };
  }
}

function makeThreadsStoreable(
  threads: ThreadChatResponseType[]
): ThreadChatMapStoreType {
  const mapped: ThreadChatMapStoreType = {};
  for (const thread of threads) {
    const {
      customer,
      assignee,
      messages,
      threadChatId,
      sequence,
      status,
      read,
      replied,
      createdAt,
      updatedAt,
    } = thread;
    const newMap: ThreadChatStoreType = {
      threadChatId,
      sequence,
      status,
      read,
      replied,
      customerId: customer.customerId,
      assigneeId: assignee ? assignee.memberId : null,
      createdAt,
      updatedAt,
      recentMessage: null,
    };
    const message = messages[0] || null;
    if (message) {
      const { customer, member, ...rest } = message;
      if (customer) {
        newMap.recentMessage = {
          ...rest,
          customerId: customer.customerId,
          memberId: null,
        };
      } else if (member) {
        newMap.recentMessage = {
          ...rest,
          customerId: null,
          memberId: member.memberId,
        };
      }
    }
    mapped[threadChatId] = newMap;
  }
  return mapped;
}

export async function bootstrapWorkspace(
  token: string,
  workspaceId: string
): Promise<IWorkspaceEntities> {
  const data = initialWorkspaceData();

  const getWorkspacePromise = getWorkspace(token, workspaceId);
  const getWorkspaceMetricsPromise = getWorkspaceMetrics(token, workspaceId);
  const getWorkspaceThreadsP = getWorkspaceThreads(token, workspaceId);

  const [workspaceData, metricsData, threadsData] = await Promise.all([
    getWorkspacePromise,
    getWorkspaceMetricsPromise,
    getWorkspaceThreadsP,
  ]);

  const { error: errWorkspace, data: workspace } = workspaceData;
  const { error: errMetrics, data: metrics } = metricsData;
  const { error: errThreads, data: threads } = threadsData;

  const hasErr = errWorkspace || errMetrics || errThreads;

  if (hasErr) {
    data.error = new Error("error bootsrapping workspace store information");
    data.isPending = false;
    return data;
  }

  if (workspace) {
    data.workspace = workspace;
    data.hasData = true;
    data.isPending = false;
  }

  if (metrics) {
    const { count } = metrics;
    data.metrics = count;
  }

  if (threads) {
    const threadsMap = makeThreadsStoreable(threads);
    data.threadChats = threadsMap;
  }

  return data;
}

export async function getOrCreateZygAccount(token: string): Promise<{
  data: AccountType | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/accounts/auth/`,
      {
        method: "POST",
        headers: {
          // "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({}),
      }
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error creating Zyg auth account with with status: ${status} and statusText: ${statusText}`
      );
      return { error, data: null };
    }

    try {
      const data = await response.json();
      const account = accountSchema.parse({ ...data });
      return { error: null, data: account };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      }
      console.error(err);
      return {
        error: new Error("error parsing account schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error(
        "error fetching auth account details - something went wrong"
      ),
      data: null,
    };
  }
}
