import { z } from "zod";

import {
  workspaceResponseSchema,
  workspaceMetricsResponseSchema,
  threadChatResponseSchema,
  threadChatMessagePreviewSchema,
  accountResponseSchema,
  membershipResponseSchema,
  workspaceCustomerResponseSchema,
  workspaceCustomersResponseSchema,
  workspaceLabelResponseSchema,
  workspaceMemberResponseSchema,
} from "./schema";
import {
  IWorkspaceEntities,
  ThreadChatStoreType,
  ThreadChatMapStoreType,
  WorkspaceCustomerMapStoreType,
  WorkspaceLabelMapStoreType,
  WorkspaceMemberMapStoreType,
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

export type WorkspaceCustomersResponseType = z.infer<
  typeof workspaceCustomersResponseSchema
>;

export type WorkspaceLabelResponseType = z.infer<
  typeof workspaceLabelResponseSchema
>;

export type WorkspaceMemberResponseType = z.infer<
  typeof workspaceMemberResponseSchema
>;

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
    customers: null,
    labels: null,
    members: null,
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
      } else console.error(err);
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
      } else console.error(err);
      return {
        error: new Error("error parsing workspace membership schema"),
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
      } else console.error(err);
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
      } else console.error(err);
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

// API call to fetch workspace labels
async function getWorkspaceLabels(
  token: string,
  workspaceId: string
): Promise<{
  data: WorkspaceLabelResponseType[] | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/`,
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
          `error fetching workspace labels: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    try {
      const data = await response.json();
      // schema validate for each item
      const labels = data.map((item: any) => {
        return workspaceLabelResponseSchema.parse({ ...item });
      });
      return { error: null, data: labels };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing workspace labels schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error(
        "error fetching workspace labels - something went wrong"
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

function makeLabelsStoreable(
  labels: WorkspaceLabelResponseType[]
): WorkspaceLabelMapStoreType {
  const mapped: WorkspaceLabelMapStoreType = {};
  for (const label of labels) {
    const { workspaceId, labelId, name, icon, createdAt, updatedAt } = label;
    mapped[labelId] = {
      workspaceId,
      labelId,
      name,
      icon,
      createdAt,
      updatedAt,
    };
  }
  return mapped;
}

function makeCustomersStoreable(
  customers: WorkspaceCustomersResponseType
): WorkspaceCustomerMapStoreType {
  const mapped: WorkspaceCustomerMapStoreType = {};
  for (const customer of customers) {
    const { customerId, ...rest } = customer;
    mapped[customerId] = { customerId, ...rest };
  }
  return mapped;
}

function makeMembersStoreable(
  members: WorkspaceMemberResponseType[]
): WorkspaceMemberMapStoreType {
  const mapped: WorkspaceMemberMapStoreType = {};
  for (const member of members) {
    const { memberId, ...rest } = member;
    mapped[memberId] = { memberId, ...rest };
  }
  return mapped;
}

export async function bootstrapWorkspace(
  token: string,
  workspaceId: string
): Promise<IWorkspaceEntities> {
  const data = initialWorkspaceData();

  const getWorkspaceP = getWorkspace(token, workspaceId);
  const getWorkspaceMemberP = getWorkspaceMember(token, workspaceId);
  const getWorkspaceCustomersP = getWorkspaceCustomers(token, workspaceId);
  const getWorkspaceMetricsP = getWorkspaceMetrics(token, workspaceId);
  const getWorkspaceThreadsP = getWorkspaceThreads(token, workspaceId);
  const getWorkspaceLabelsP = getWorkspaceLabels(token, workspaceId);
  const getWorkspaceMembersP = getWorkspaceMembers(token, workspaceId);

  const [
    workspaceData,
    memberData,
    customerData,
    metricsData,
    threadsData,
    labelsData,
    membersData,
  ] = await Promise.all([
    getWorkspaceP,
    getWorkspaceMemberP,
    getWorkspaceCustomersP,
    getWorkspaceMetricsP,
    getWorkspaceThreadsP,
    getWorkspaceLabelsP,
    getWorkspaceMembersP,
  ]);

  const { error: errWorkspace, data: workspace } = workspaceData;
  const { error: errMember, data: member } = memberData;
  const { error: errCustomer, data: customers } = customerData;
  const { error: errMetrics, data: metrics } = metricsData;
  const { error: errThreads, data: threads } = threadsData;
  const { error: errLabels, data: labels } = labelsData;
  const { error: errMembers, data: members } = membersData;

  const hasErr =
    errWorkspace ||
    errMember ||
    errCustomer ||
    errMetrics ||
    errThreads ||
    errLabels ||
    errMembers;

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

  if (customers && customers.length > 0) {
    const customersMap = makeCustomersStoreable(customers);
    data.customers = customersMap;
  }

  if (metrics) {
    const { count } = metrics;
    data.metrics = count;
  }

  if (members && members.length > 0) {
    const membersMap = makeMembersStoreable(members);
    data.members = membersMap;
  }

  if (threads && threads.length > 0) {
    const threadsMap = makeThreadsStoreable(threads);
    data.threadChats = threadsMap;
  }

  if (labels && labels.length > 0) {
    const labelsMap = makeLabelsStoreable(labels);
    data.labels = labelsMap;
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
      } else console.error(err);
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

export async function getWorkspaceCustomers(
  token: string,
  workspaceId: string
): Promise<{
  data: WorkspaceCustomersResponseType | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/customers/`,
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
          `error fetching workspace customers: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    try {
      const data = await response.json();
      // schema validate for each item
      const customers = data.map((item: any) => {
        return workspaceCustomerResponseSchema.parse({ ...item });
      });
      return { error: null, data: customers };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing workspace customers schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error(
        "error fetching workspace customers - something went wrong"
      ),
      data: null,
    };
  }
}

export async function createWorkspace(
  token: string,
  body: { name: string }
): Promise<{
  data: { workspaceId: string; workspaceName: string } | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ ...body }),
      }
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error creating workspace with status: ${status} and statusText: ${statusText}`
      );
      return { error, data: null };
    }
    const data = await response.json();
    const { workspaceId, name } = data;
    return {
      error: null,
      data: { workspaceId, workspaceName: name },
    };
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error creating workspace - something went wrong"),
      data: null,
    };
  }
}

export async function updateWorkspace(
  token: string,
  workspaceId: string,
  body: { name: string }
): Promise<{
  data: { workspaceId: string; workspaceName: string } | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/`,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ ...body }),
      }
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error creating workspace with status: ${status} and statusText: ${statusText}`
      );
      return { error, data: null };
    }
    const data = await response.json();
    const { workspaceId: id, name } = data;
    return {
      error: null,
      data: { workspaceId: id, workspaceName: name },
    };
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error creating workspace - something went wrong"),
      data: null,
    };
  }
}

export async function getWorkspaceMembers(
  token: string,
  workspaceId: string
): Promise<{
  data: WorkspaceMemberResponseType[] | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/members/`,
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
          `error fetching workspace members: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    try {
      const data = await response.json();
      // schema validate for each item
      const members = data.map((item: any) => {
        return workspaceMemberResponseSchema.parse({ ...item });
      });
      return { error: null, data: members };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing workspace members schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error(
        "error fetching workspace members - something went wrong"
      ),
      data: null,
    };
  }
}
