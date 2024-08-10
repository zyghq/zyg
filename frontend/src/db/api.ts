import { z } from "zod";
import {
  workspaceResponseSchema,
  workspaceMetricsResponseSchema,
  threadResponseSchema,
  accountResponseSchema,
  authMemberResponseSchema,
  customerResponseSchema,
  labelResponseSchema,
  memberResponseSchema,
  patResponseSchema,
  threadChatResponseSchema,
} from "./schema";
import {
  IWorkspaceEntities,
  IWorkspaceValueObjects,
  ThreadMap,
  CustomerMap,
  LabelMap,
  MemberMap,
  PatMap,
} from "./store";

import { Account, Workspace, AuthMember, Pat, Customer } from "./entities";

export type ThreadResponse = z.infer<typeof threadResponseSchema>;

export type ThreadChatResponse = z.infer<typeof threadChatResponseSchema>;

export type LabelResponse = z.infer<typeof labelResponseSchema>;

export type MemberResponse = z.infer<typeof memberResponseSchema>;

export type WorkspaceMetricsResponse = z.infer<
  typeof workspaceMetricsResponseSchema
>;

function initialWorkspaceData(): IWorkspaceEntities & IWorkspaceValueObjects {
  return {
    hasData: false,
    isPending: true,
    error: null,
    threadAppliedFilters: null,
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
    threads: null,
    customers: null,
    labels: null,
    members: null,
    pats: null,
  };
}

export async function getWorkspace(
  token: string,
  workspaceId: string
): Promise<{ data: Workspace | null; error: Error | null }> {
  try {
    // make a request
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
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
    console.log(data);
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

export async function getWorkspaceMember(
  token: string,
  workspaceId: string
): Promise<{ data: AuthMember | null; error: Error | null }> {
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
    console.log(data);
    // parse into schema
    try {
      const member = authMemberResponseSchema.parse({ ...data });
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

export async function getWorkspaceMetrics(
  token: string,
  workspaceId: string
): Promise<{ data: WorkspaceMetricsResponse | null; error: Error | null }> {
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

export async function getWorkspaceThreads(
  token: string,
  workspaceId: string
): Promise<{ data: ThreadResponse[] | null; error: Error | null }> {
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
      console.log(data);
      // schema validate for each item
      const threads = data.map((item: any) => {
        return threadResponseSchema.parse({ ...item });
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

export async function getWorkspaceLabels(
  token: string,
  workspaceId: string
): Promise<{
  data: LabelResponse[] | null;
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
      console.log(data);
      // schema validate for each item
      const labels = data.map((item: any) => {
        return labelResponseSchema.parse({ ...item });
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

export async function getPats(token: string): Promise<{
  data: Pat[] | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(`${import.meta.env.VITE_ZYG_URL}/pats/`, {
      method: "GET",
      headers: {
        // "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `error fetching account pats: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      // schema validate for each item
      const pats = data.map((item: any) => {
        return patResponseSchema.parse({ ...item });
      });
      return { error: null, data: pats };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing account pat schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error fetching account pats - something went wrong"),
      data: null,
    };
  }
}

export async function createPat(
  token: string,
  body: { name: string; description: string }
): Promise<{
  data: Pat | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(`${import.meta.env.VITE_ZYG_URL}/pats/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ ...body }),
    });
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error creating account pat with status: ${status} and statusText: ${statusText}`
      );
      return { error, data: null };
    }
    try {
      const data = await response.json();

      console.log(data);
      const pat = patResponseSchema.parse({ ...data });
      return { error: null, data: pat };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing account pat schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error creating account pat - something went wrong"),
      data: null,
    };
  }
}

function makeThreadsStoreable(threads: ThreadResponse[]): ThreadMap {
  const mapped: ThreadMap = {};
  for (const thread of threads) {
    const {
      threadId,
      customer,
      assignee,
      egressFirstSeq,
      egressLastSeq,
      ingressFirstSeq,
      ingressLastSeq,
      ingressCustomer,
      egressMember,
      ...rest
    } = thread;
    const customerId = customer.customerId;
    const assigneeId = assignee?.memberId || null;
    const inboundFirstSeq = ingressFirstSeq || null;
    const inboundLastSeq = ingressLastSeq || null;
    const inboundCustomerId = ingressCustomer?.customerId || null;
    const outboundFirstSeq = egressFirstSeq || null;
    const outboundLastSeq = egressLastSeq || null;
    const outboundMemberId = egressMember?.memberId || null;
    mapped[threadId] = {
      threadId,
      ...rest,
      customerId: customerId,
      assigneeId: assigneeId,
      inboundFirstSeq,
      inboundLastSeq,
      inboundCustomerId,
      outboundFirstSeq,
      outboundLastSeq,
      outboundMemberId,
    };
  }
  return mapped;
}

function makeLabelsStoreable(labels: LabelResponse[]): LabelMap {
  const mapped: LabelMap = {};
  for (const label of labels) {
    const { labelId, name, icon, createdAt, updatedAt } = label;
    mapped[labelId] = {
      labelId,
      name,
      icon,
      createdAt,
      updatedAt,
    };
  }
  return mapped;
}

function makeCustomersStoreable(customers: Customer[]): CustomerMap {
  const mapped: CustomerMap = {};
  for (const customer of customers) {
    const { customerId, ...rest } = customer;
    mapped[customerId] = { customerId, ...rest };
  }
  return mapped;
}

function makeMembersStoreable(members: MemberResponse[]): MemberMap {
  const mapped: MemberMap = {};
  for (const member of members) {
    const { memberId, ...rest } = member;
    mapped[memberId] = { memberId, ...rest };
  }
  return mapped;
}

function makePatsStoreable(pats: Pat[]): PatMap {
  const mapped: PatMap = {};
  for (const pat of pats) {
    const { patId, ...rest } = pat;
    mapped[patId] = { patId, ...rest };
  }
  return mapped;
}

export async function getOrCreateZygAccount(
  token: string,
  body?: { name: string }
): Promise<{
  data: Account | null;
  error: Error | null;
}> {
  try {
    const reqBody = body ? JSON.stringify({ ...body }) : JSON.stringify({});
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/accounts/auth/`,
      {
        method: "POST",
        headers: {
          // "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: reqBody,
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

      console.log(data);
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
  data: Customer[] | null;
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

      console.log(data);
      // schema validate for each item
      const customers = data.map((item: any) => {
        return customerResponseSchema.parse({ ...item });
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

    console.log(data);
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
        `error updating workspace with status: ${status} and statusText: ${statusText}`
      );
      return { error, data: null };
    }
    const data = await response.json();

    console.log(data);
    const { workspaceId: id, name } = data;
    return {
      error: null,
      data: { workspaceId: id, workspaceName: name },
    };
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error updating workspace - something went wrong"),
      data: null,
    };
  }
}

export async function getWorkspaceMembers(
  token: string,
  workspaceId: string
): Promise<{
  data: MemberResponse[] | null;
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

      console.log(data);
      // schema validate for each item
      const members = data.map((item: any) => {
        return memberResponseSchema.parse({ ...item });
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

export async function deletePat(
  token: string,
  patId: string
): Promise<{
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/pats/${patId}/`,
      {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `error deleting pat with status: ${status} and statusText: ${statusText}`
        ),
      };
    }
    return { error: null };
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error deleting pat - something went wrong"),
    };
  }
}

export async function bootstrapWorkspace(
  token: string,
  workspaceId: string
): Promise<IWorkspaceEntities & IWorkspaceValueObjects> {
  const data = initialWorkspaceData();

  const getWorkspaceP = getWorkspace(token, workspaceId);
  const getWorkspaceMemberP = getWorkspaceMember(token, workspaceId);
  const getWorkspaceCustomersP = getWorkspaceCustomers(token, workspaceId);
  const getWorkspaceMetricsP = getWorkspaceMetrics(token, workspaceId);
  const getWorkspaceThreadsP = getWorkspaceThreads(token, workspaceId);
  const getWorkspaceLabelsP = getWorkspaceLabels(token, workspaceId);
  const getWorkspaceMembersP = getWorkspaceMembers(token, workspaceId);
  const getAccountPatsP = getPats(token);

  const [
    workspaceData,
    memberData,
    customerData,
    metricsData,
    threadsData,
    labelsData,
    membersData,
    patsData,
  ] = await Promise.all([
    getWorkspaceP,
    getWorkspaceMemberP,
    getWorkspaceCustomersP,
    getWorkspaceMetricsP,
    getWorkspaceThreadsP,
    getWorkspaceLabelsP,
    getWorkspaceMembersP,
    getAccountPatsP,
  ]);

  const { error: errWorkspace, data: workspace } = workspaceData;
  const { error: errMember, data: member } = memberData;
  const { error: errCustomer, data: customers } = customerData;
  const { error: errMetrics, data: metrics } = metricsData;
  const { error: errThreads, data: threads } = threadsData;
  const { error: errLabels, data: labels } = labelsData;
  const { error: errMembers, data: members } = membersData;
  const { error: errPats, data: pats } = patsData;

  const hasErr =
    errWorkspace ||
    errMember ||
    errCustomer ||
    errMetrics ||
    errThreads ||
    errLabels ||
    errMembers ||
    errPats;

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
    data.threads = threadsMap;
  }

  if (labels && labels.length > 0) {
    const labelsMap = makeLabelsStoreable(labels);
    data.labels = labelsMap;
  }

  if (pats && pats.length > 0) {
    const patsMap = makePatsStoreable(pats);
    data.pats = patsMap;
  }

  return data;
}

export async function getWorkspaceThreadChatMessages(
  token: string,
  workspaceId: string,
  threadId: string
): Promise<{
  data: ThreadChatResponse[] | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/chat/${threadId}/messages/`,
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
          `error fetching workspace thread chat messages: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const messages = data.map((item: any) => {
        return threadChatResponseSchema.parse({ ...item });
      });
      return { error: null, data: messages };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing workspace thread chat messages schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error(
        "error fetching workspace thread chat messages - something went wrong"
      ),
      data: null,
    };
  }
}

export async function createWorkspaceLabel(
  token: string,
  workspaceId: string,
  body: { name: string }
): Promise<{ data: LabelResponse | null; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/`,
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
        `error creating label with status: ${status} and statusText: ${statusText}`
      );
      return { error, data: null };
    }
    try {
      const data = await response.json();

      console.log(data);
      const parsed = labelResponseSchema.parse({ ...data });
      return { error: null, data: parsed };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing workspace label schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error creating label - something went wrong"),
      data: null,
    };
  }
}

export async function updateWorkspaceLabel(
  token: string,
  workspaceId: string,
  labelId: string,
  body: { name: string }
): Promise<{ data: LabelResponse | null; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/${labelId}/`,
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
        `error updating label with status: ${status} and statusText: ${statusText}`
      );
      return { error, data: null };
    }
    try {
      const data = await response.json();

      console.log(data);
      const parsed = labelResponseSchema.parse({ ...data });
      return { error: null, data: parsed };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing workspace label schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error workspace label - something went wrong"),
      data: null,
    };
  }
}

export async function updateThread(
  token: string,
  workspaceId: string,
  threadId: string,
  body: object
): Promise<{ data: ThreadResponse | null; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/chat/${threadId}/`,
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
        `error updating thread with status: ${status} and statusText: ${statusText}`
      );
      return { error, data: null };
    }
    try {
      const data = await response.json();

      console.log(data);
      const parsed = threadResponseSchema.parse({ ...data });
      return { error: null, data: parsed };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing thread chat schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error("error thread - something went wrong"),
      data: null,
    };
  }
}

export async function sendThreadChatMessage(
  token: string,
  workspaceId: string,
  threadChatId: string,
  body: { message: string }
): Promise<{
  data: ThreadChatResponse | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/chat/${threadChatId}/messages/`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ ...body }),
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `error fetching workspace thread chat messages: ${status} ${statusText}`
        ),
        data: null,
      };
    }

    try {
      const data = await response.json();

      console.log(data);
      const parsed = threadChatResponseSchema.parse({ ...data });
      return { error: null, data: parsed };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        error: new Error("error parsing workspace thread chat messages schema"),
        data: null,
      };
    }
  } catch (err) {
    console.error(err);
    return {
      error: new Error(
        "error sending workspace thread chat messages - something went wrong"
      ),
      data: null,
    };
  }
}
