import { z } from "zod";

import {
  Account,
  AuthMember,
  Customer,
  customerTransformer,
  labelTransformer,
  memberTransformer,
  Pat,
  patTransformer,
  threadTransformer,
  Workspace,
} from "./models";
import {
  accountResponseSchema,
  authMemberResponseSchema,
  CustomerEventResponse,
  customerEventSchema,
  customerResponseSchema,
  LabelResponse,
  labelResponseSchema,
  MemberResponse,
  memberResponseSchema,
  MessageAttachmentResponse,
  messageAttachmentResponseSchema,
  patResponseSchema,
  PostmarkMailServerSetting,
  postmarkMailServerSettingSchema,
  ThreadLabelResponse,
  threadLabelResponseSchema,
  ThreadMessageResponse,
  threadMessageResponseSchema,
  ThreadResponse,
  threadResponseSchema,
  WorkspaceMetricsResponse,
  workspaceMetricsResponseSchema,
  workspaceResponseSchema

} from "./schema";
import {
  CustomerMap,
  IWorkspaceEntities,
  IWorkspaceValueObjects,
  LabelMap,
  MemberMap,
  PatMap,
  ThreadMap,
} from "./store";

// Returns the default state of the workspace store.
function initialWorkspaceData(): IWorkspaceEntities & IWorkspaceValueObjects {
  return {
    customers: null,
    error: null,
    hasData: false,
    isPending: true,
    labels: null,
    member: null,
    members: null,
    metrics: {
      active: 0,
      assignedToMe: 0,
      hold: 0,
      labels: [],
      needsFirstResponse: 0,
      needsNextResponse: 0,
      otherAssigned: 0,
      snoozed: 0,
      unassigned: 0,
      waitingOnCustomer: 0,
    },
    pats: null,
    threadAppliedFilters: null,
    threads: null,
    threadSortKey: null,
    workspace: null,
  };
}

export async function getWorkspace(
  token: string,
  workspaceId: string,
): Promise<{ data: null | Workspace; error: Error | null }> {
  try {
    // make a request
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace details: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const workspace = workspaceResponseSchema.parse({ ...data });
      return { data: workspace, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace schema"),
      };
    }
  } catch (error) {
    console.error(error);
    return {
      data: null,
      error: new Error(
        "error fetching workspace details - something went wrong",
      ),
    };
  }
}

export async function getWorkspaceMember(
  token: string,
  workspaceId: string,
): Promise<{ data: AuthMember | null; error: Error | null }> {
  try {
    // make a request
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/members/me/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace member details: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const member = authMemberResponseSchema.parse({ ...data });
      return { data: member, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace membership schema"),
      };
    }
  } catch (error) {
    console.error(error);
    return {
      data: null,
      error: new Error(
        "error fetching workspace member details - something went wrong",
      ),
    };
  }
}

export async function getWorkspaceMetrics(
  token: string,
  workspaceId: string,
): Promise<{ data: null | WorkspaceMetricsResponse; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/metrics/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace metrics: ${status} ${statusText}`,
        ),
      };
    }
    const data = await response.json();
    try {
      const metrics = workspaceMetricsResponseSchema.parse({ ...data });
      return { data: metrics, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace metrics schema"),
      };
    }
  } catch (error) {
    console.error(error);
    return {
      data: null,
      error: new Error(
        "error fetching workspace metrics - something went wrong",
      ),
    };
  }
}

export async function getWorkspaceThreads(
  token: string,
  workspaceId: string,
): Promise<{ data: null | ThreadResponse[]; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace threads: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const threads = data.map((item: any) => {
        return threadResponseSchema.parse({ ...item });
      });
      return { data: threads, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace threads schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error(
        "error fetching workspace threads - something went wrong",
      ),
    };
  }
}

export async function getWorkspaceLabels(
  token: string,
  workspaceId: string,
): Promise<{
  data: LabelResponse[] | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/`,
      {
        headers: {
          // "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace labels: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      // schema validate for each item
      const labels = data.map((item: any) => {
        return labelResponseSchema.parse({ ...item });
      });
      return { data: labels, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace labels schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error(
        "error fetching workspace labels - something went wrong",
      ),
    };
  }
}

export async function getPats(token: string): Promise<{
  data: null | Pat[];
  error: Error | null;
}> {
  try {
    const response = await fetch(`${import.meta.env.VITE_ZYG_URL}/pats/`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      method: "GET",
    });

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching account pats: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const pats = data.map((item: any) => {
        return patResponseSchema.parse({ ...item });
      });
      return { data: pats, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing account pat schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error fetching account pats - something went wrong"),
    };
  }
}

export async function createPat(
  token: string,
  body: { description: string; name: string },
): Promise<{
  data: null | Pat;
  error: Error | null;
}> {
  try {
    const response = await fetch(`${import.meta.env.VITE_ZYG_URL}/pats/`, {
      body: JSON.stringify({ ...body }),
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      method: "POST",
    });
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error creating account pat with status: ${status} and statusText: ${statusText}`,
      );
      return { data: null, error };
    }

    try {
      const data = await response.json();
      console.log(data);
      const pat = patResponseSchema.parse({ ...data });
      return { data: pat, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing account pat schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error creating account pat - something went wrong"),
    };
  }
}

function makeThreadsStoreable(threads: ThreadResponse[]): ThreadMap {
  const mapped: ThreadMap = {};
  const transfomer = threadTransformer();
  for (const thread of threads) {
    const [threadId, normalized] = transfomer.normalize(thread);
    mapped[threadId] = normalized;
  }
  return mapped;
}

function makeLabelsStoreable(labels: LabelResponse[]): LabelMap {
  const mapped: LabelMap = {};
  const transfomer = labelTransformer();
  for (const label of labels) {
    const [labelId, normalized] = transfomer.normalize(label);
    mapped[labelId] = normalized;
  }
  return mapped;
}

function makeCustomersStoreable(customers: Customer[]): CustomerMap {
  const mapped: CustomerMap = {};
  const transformer = customerTransformer();
  for (const customer of customers) {
    const [customerId, normalized] = transformer.normalize(customer);
    mapped[customerId] = normalized;
  }
  return mapped;
}

function makeMembersStoreable(members: MemberResponse[]): MemberMap {
  const mapped: MemberMap = {};
  const transfomer = memberTransformer();
  for (const member of members) {
    const [memberId, normalized] = transfomer.normalize(member);
    mapped[memberId] = normalized;
  }
  return mapped;
}

function makePatsStoreable(pats: Pat[]): PatMap {
  const mapped: PatMap = {};
  const transfomer = patTransformer();
  for (const pat of pats) {
    const [patId, normalized] = transfomer.normalize(pat);
    mapped[patId] = normalized;
  }
  return mapped;
}

export async function getOrCreateZygAccount(
  token: string,
  body?: { name: string },
): Promise<{
  data: Account | null;
  error: Error | null;
}> {
  try {
    const reqBody = body ? JSON.stringify({ ...body }) : JSON.stringify({});
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/accounts/auth/`,
      {
        body: reqBody,
        headers: {
          // "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        method: "POST",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error creating Zyg auth account with with status: ${status} and statusText: ${statusText}`,
      );
      return { data: null, error };
    }

    try {
      const data = await response.json();
      console.log(data);
      const account = accountResponseSchema.parse({ ...data });
      return { data: account, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing account schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error(
        "error fetching auth account details - something went wrong",
      ),
    };
  }
}

export async function getWorkspaceCustomers(
  token: string,
  workspaceId: string,
): Promise<{
  data: Customer[] | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/customers/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace customers: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const customers = data.map((item: any) => {
        return customerResponseSchema.parse({ ...item });
      });
      return { data: customers, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace customers schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error(
        "error fetching workspace customers - something went wrong",
      ),
    };
  }
}

export async function createWorkspace(
  token: string,
  body: { name: string },
): Promise<{
  data: { workspaceId: string; workspaceName: string } | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/`,
      {
        body: JSON.stringify({ ...body }),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error creating workspace with status: ${status} and statusText: ${statusText}`,
      );
      return { data: null, error };
    }
    const data = await response.json();
    console.log(data);
    const { name, workspaceId } = data;
    return {
      data: { workspaceId, workspaceName: name },
      error: null,
    };
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error creating workspace - something went wrong"),
    };
  }
}

export async function updateWorkspace(
  token: string,
  workspaceId: string,
  body: { name: string },
): Promise<{
  data: { workspaceId: string; workspaceName: string } | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/`,
      {
        body: JSON.stringify({ ...body }),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PATCH",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error updating workspace with status: ${status} and statusText: ${statusText}`,
      );
      return { data: null, error };
    }

    const data = await response.json();
    console.log(data);
    const { name, workspaceId: id } = data;
    return {
      data: { workspaceId: id, workspaceName: name },
      error: null,
    };
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error updating workspace - something went wrong"),
    };
  }
}

export async function getWorkspaceMembers(
  token: string,
  workspaceId: string,
): Promise<{
  data: MemberResponse[] | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/members/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace members: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const members = data.map((item: any) => {
        return memberResponseSchema.parse({ ...item });
      });
      return { data: members, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace members schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error(
        "error fetching workspace members - something went wrong",
      ),
    };
  }
}

export async function deletePat(
  token: string,
  patId: string,
): Promise<{
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/pats/${patId}/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "DELETE",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `error deleting pat with status: ${status} and statusText: ${statusText}`,
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

// @sanchitrk:
// improve usage with react query with zustand store.
// can we do this better?
export async function bootstrapWorkspace(
  token: string,
  workspaceId: string,
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

  const { data: workspace, error: errWorkspace } = workspaceData;
  const { data: member, error: errMember } = memberData;
  const { data: customers, error: errCustomer } = customerData;
  const { data: metrics, error: errMetrics } = metricsData;
  const { data: threads, error: errThreads } = threadsData;
  const { data: labels, error: errLabels } = labelsData;
  const { data: members, error: errMembers } = membersData;
  const { data: pats, error: errPats } = patsData;

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
    data.error = new Error("error bootstrapping workspace store information");
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

export async function getWorkspaceThreadMessages(
  token: string,
  workspaceId: string,
  threadId: string,
): Promise<{
  data: null | ThreadMessageResponse[];
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/${threadId}/messages/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace thread chat messages: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const messages = data.map((item: any) => {
        return threadMessageResponseSchema.parse({ ...item });
      });

      return { data: messages, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace thread chat messages schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error(
        "error fetching workspace thread chat messages - something went wrong",
      ),
    };
  }
}

export async function createWorkspaceLabel(
  token: string,
  workspaceId: string,
  body: { name: string },
): Promise<{ data: LabelResponse | null; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/`,
      {
        body: JSON.stringify({ ...body }),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error creating label with status: ${status} and statusText: ${statusText}`,
      );
      return { data: null, error };
    }

    try {
      const data = await response.json();
      console.log(data);
      const parsed = labelResponseSchema.parse({ ...data });
      return { data: parsed, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace label schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error creating label - something went wrong"),
    };
  }
}

export async function updateWorkspaceLabel(
  token: string,
  workspaceId: string,
  labelId: string,
  body: { name: string },
): Promise<{ data: LabelResponse | null; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/${labelId}/`,
      {
        body: JSON.stringify({ ...body }),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PATCH",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error updating label with status: ${status} and statusText: ${statusText}`,
      );
      return { data: null, error };
    }

    try {
      const data = await response.json();
      console.log(data);
      const parsed = labelResponseSchema.parse({ ...data });
      return { data: parsed, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace label schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error workspace label - something went wrong"),
    };
  }
}

export async function updateThread(
  token: string,
  workspaceId: string,
  threadId: string,
  body: object,
): Promise<{ data: null | ThreadResponse; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/${threadId}/`,
      {
        body: JSON.stringify({ ...body }),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PATCH",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      const error = new Error(
        `error updating thread with status: ${status} and statusText: ${statusText}`,
      );
      return { data: null, error };
    }

    try {
      const data = await response.json();
      console.log(data);
      const parsed = threadResponseSchema.parse({ ...data });
      return { data: parsed, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing thread chat schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error thread - something went wrong"),
    };
  }
}

export async function getThreadLabels(
  token: string,
  workspaceId: string,
  threadId: string,
): Promise<{
  data: null | ThreadLabelResponse[];
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/${threadId}/labels/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "GET",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace thread labels: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      const labels = data.map((item: any) => {
        return threadLabelResponseSchema.parse({ ...item });
      });
      return { data: labels, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.log(err);
      return {
        data: null,
        error: new Error("error parsing workspace thread label schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error(
        "error fetching workspace thread labels - something went wrong",
      ),
    };
  }
}

export async function putThreadLabel(
  token: string,
  workspaceId: string,
  threadId: string,
  body: { icon: string; name: string },
): Promise<{
  data: null | ThreadLabelResponse;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/${threadId}/labels/`,
      {
        body: JSON.stringify({ ...body }),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PUT",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(`error setting thread label: ${status} ${statusText}`),
      };
    }

    try {
      const data = await response.json();
      console.log(data);
      const label = threadLabelResponseSchema.parse({ ...data });
      return { data: label, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.log(err);
      return {
        data: null,
        error: new Error("error parsing workspace thread label schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error setting thread label - something went wrong"),
    };
  }
}

export async function deleteThreadLabel(
  token: string,
  workspaceId: string,
  threadId: string,
  labelId: string,
): Promise<{
  data: boolean | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/${threadId}/labels/${labelId}/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "DELETE",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error deleting thread label: ${status} ${statusText}`,
        ),
      };
    }
    return {
      data: true,
      error: null,
    };
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("error deleting thread label - something went wrong"),
    };
  }
}

export async function sendThreadChatMessage(
  token: string,
  workspaceId: string,
  threadId: string,
  body: { message: string },
): Promise<{
  data: null | ThreadMessageResponse;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/chat/${threadId}/messages/`,
      {
        body: JSON.stringify({ ...body }),
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "POST",
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching workspace thread chat messages: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();

      console.log(data);
      const parsed = threadMessageResponseSchema.parse({ ...data });
      return { data: parsed, error: null };
    } catch (err) {
      if (err instanceof z.ZodError) {
        console.error(err.message);
      } else console.error(err);
      return {
        data: null,
        error: new Error("error parsing workspace thread chat messages schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error(
        "error sending workspace thread chat messages - something went wrong",
      ),
    };
  }
}

export async function getCustomerEvents(
  token: string,
  workspaceId: string,
  customerId: string,
): Promise<{
  data: CustomerEventResponse[] | null;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/customers/events/${customerId}/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "GET",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching customer events: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      const events = data.map((item: any) => {
        return customerEventSchema.parse({ ...item });
      });
      return { data: events, error: null };
    } catch (err) {
      console.error(err);
      return {
        data: null,
        error: new Error("error parsing customer event schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("something went wrong"),
    };
  }
}

export async function getMessageAttachment(
  token: string,
  workspaceId: string,
  messageId: string,
  attachmentId: string,
): Promise<{ data: MessageAttachmentResponse | null; error: Error | null }> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/messages/${messageId}/attachments/${attachmentId}/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "GET",
      },
    );
    if (!response.ok) {
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching customer events: ${status} ${statusText}`,
        ),
      };
    }

    try {
      const data = await response.json();
      const attachment = messageAttachmentResponseSchema.parse({ ...data });
      return { data: attachment, error: null };
    } catch (err) {
      console.error(err);
      return {
        data: null,
        error: new Error("error parsing message attachment schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("something went wrong"),
    };
  }
}


export async function getEmailSetting(
  token: string,
  workspaceId: string,
): Promise<{
  data: null | PostmarkMailServerSetting;
  error: Error | null;
}> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/postmark/servers/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "GET",
      },
    );
    if (!response.ok) {
      if (response.status === 404) {
        return {
          data: null,
          error: null,
        };
      }
      // something went wrong
      const { status, statusText } = response;
      return {
        data: null,
        error: new Error(
          `error fetching mail setting: ${status} ${statusText}`,
        ),
      };
    }
    try {
      const data = await response.json();
      const setting = postmarkMailServerSettingSchema.parse({ ...data });
      return { data: setting, error: null };
    } catch (err) {
      console.error(err);
      return {
        data: null,
        error: new Error("error parsing email setting schema"),
      };
    }
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: new Error("something went wrong"),
    };
  }
}
