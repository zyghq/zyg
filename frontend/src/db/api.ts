import { z } from "zod";

import {
  workspaceResponseSchema,
  workspaceMetricsResponseSchema,
  threadChatResponseSchema,
  threadChatMessagePreviewSchema,
  accountResponseSchema,
  membershipResponseSchema,
} from "./schema";
import {
  IWorkspaceEntities,
  ThreadChatStoreType,
  ThreadChatMapStoreType,
} from "./store";

export type WorkspaceResponseType = z.infer<typeof workspaceResponseSchema>;

export type WorkspaceMetricsResponseType = z.infer<
  typeof workspaceMetricsResponseSchema
>;

export type ThreadChatResponseType = z.infer<typeof threadChatResponseSchema>;

export type ThreadChatMessagePreviewType = z.infer<
  typeof threadChatMessagePreviewSchema
>;

export type AccountResponseType = z.infer<typeof accountResponseSchema>;

export type MembershipResponseType = z.infer<typeof membershipResponseSchema>;

function initialWorkspaceData(): IWorkspaceEntities {
  return {
    hasData: false,
    isPending: true,
    error: null,
    workspace: null,
    member: null,
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

async function getWorkspaceMember(
  token: string,
  workspaceId: string
): Promise<{ data: MembershipResponseType | null; error: Error | null }> {
  try {
    // make a request
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/members/me/`,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `error fetching workspace member details: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    const data = await response.json();
    // parse into schema
    try {
      const member = membershipResponseSchema.parse({ ...data });
      return { error: null, data: member };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      }
      console.error(err);
      return {
        error: new Error("error parsing workspace member schema"),
        data: null,
      };
    }
  } catch (error) {
    console.error(error);
    return {
      error: new Error(
        "error fetching workspace member details - something went wrong"
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
  const getWorkspaceMemberPromise = getWorkspaceMember(token, workspaceId);
  const getWorkspaceMetricsPromise = getWorkspaceMetrics(token, workspaceId);
  const getWorkspaceThreadsP = getWorkspaceThreads(token, workspaceId);

  const [workspaceData, memberData, metricsData, threadsData] =
    await Promise.all([
      getWorkspacePromise,
      getWorkspaceMemberPromise,
      getWorkspaceMetricsPromise,
      getWorkspaceThreadsP,
    ]);

  const { error: errWorkspace, data: workspace } = workspaceData;
  const { error: errMember, data: member } = memberData;
  const { error: errMetrics, data: metrics } = metricsData;
  const { error: errThreads, data: threads } = threadsData;

  const hasErr = errWorkspace || errMember || errMetrics || errThreads;

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

  if (member) {
    data.member = member;
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
  data: AccountResponseType | null;
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
      const account = accountResponseSchema.parse({ ...data });
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
