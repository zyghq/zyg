import { HTTPError } from "@/db/errors";
import {
  Account,
  AuthMember,
  Customer,
  CustomerEventResponse,
  customerTransformer,
  LabelResponse,
  labelTransformer,
  MessageAttachmentResponse,
  Pat,
  patTransformer,
  PostmarkMailServerSetting,
  ThreadLabelResponse,
  ThreadMessageResponse,
  ThreadResponse,
  threadTransformer,
  Workspace,
  WorkspaceMetricsResponse,
} from "@/db/models";
import {
  accountSchema,
  authMemberSchema,
  customerEventSchema,
  customerSchema,
  labelSchema,
  messageAttachmentSchema,
  patSchema,
  postmarkMailServerSettingSchema,
  threadLabelSchema,
  threadMessageSchema,
  threadSchema,
  workspaceMetricsSchema,
  workspaceSchema,
} from "@/db/schema";
import {
  CustomerMap,
  IWorkspaceEntitiesBootstrap,
  IWorkspaceValueObjects,
  LabelMap,
  PatMap,
  ThreadMap,
} from "@/db/store";
import { z } from "zod";

interface ApiResponse<T> {
  data: null | T;
  error: Error | null;
}

// Returns the default state of the workspace store.
function initialWorkspaceData(): IWorkspaceEntitiesBootstrap &
  IWorkspaceValueObjects {
  return {
    customers: null,
    error: null,
    hasData: false,
    isPending: true,
    labels: null,
    member: null,
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

/**
 * Retrieves the details of a specified workspace.
 *
 * @param {string} token - The authorization token required for accessing the workspace.
 * @param {string} workspaceId - The unique identifier of the workspace to retrieve.
 * @return {Promise<ApiResponse<Workspace>>} A promise that resolves with the workspace data or an error.
 */
export async function getWorkspace(
  token: string,
  workspaceId: string,
): Promise<ApiResponse<Workspace>> {
  try {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch workspace details",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const workspace = workspaceSchema.parse(data);
    return {
      data: workspace,
      error: null,
    };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid workspace schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch workspace details";

    console.error("[getWorkspace]", {
      error,
      timestamp: new Date().toISOString(),
      workspaceId,
    });

    return {
      data: null,
      error: new Error(errorMessage),
    };
  }
}

/**
 * Retrieves details of a workspace member.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @return {Promise<ApiResponse<AuthMember>>} A promise resolving with member data or an error.
 */
export async function getWorkspaceMember(
  token: string,
  workspaceId: string,
): Promise<ApiResponse<AuthMember>> {
  try {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message:
          errorData?.message || "Failed to fetch workspace member details",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const member = authMemberSchema.parse(data);
    return { data: member, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid member schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch workspace member details";

    console.error("[getWorkspaceMember]", {
      error,
      timestamp: new Date().toISOString(),
      workspaceId,
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Retrieves workspace metrics.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @return {Promise<ApiResponse<WorkspaceMetricsResponse>>} A promise resolving with metrics data or an error.
 */
export async function getWorkspaceMetrics(
  token: string,
  workspaceId: string,
): Promise<ApiResponse<WorkspaceMetricsResponse>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch workspace metrics",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const metrics = workspaceMetricsSchema.parse(data);
    return {
      data: metrics,
      error: null,
    };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid workspace metrics schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch workspace metrics";

    console.error("[getWorkspaceMetrics]", {
      error,
      timestamp: new Date().toISOString(),
      workspaceId,
    });

    return {
      data: null,
      error: new Error(errorMessage),
    };
  }
}

/**
 * Retrieves workspace threads.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @return {Promise<ApiResponse<ThreadResponse[]>>} A promise resolving with threads data or an error.
 */
export async function getWorkspaceThreads(
  token: string,
  workspaceId: string,
): Promise<ApiResponse<ThreadResponse[]>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch workspace threads",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const threads = data.map((item: any) => threadSchema.parse(item));
    return { data: threads, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid workspace threads schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch workspace threads";

    console.error("[getWorkspaceThreads]", {
      error,
      timestamp: new Date().toISOString(),
      workspaceId,
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Retrieves workspace labels.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @return {Promise<ApiResponse<LabelResponse[]>>} A promise resolving with label data or an error.
 */
export async function getWorkspaceLabels(
  token: string,
  workspaceId: string,
): Promise<ApiResponse<LabelResponse[]>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "GET",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch workspace labels",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const labels = data.map((item: any) => labelSchema.parse(item));
    return { data: labels, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid workspace labels schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch workspace labels";

    console.error("[getWorkspaceLabels]", {
      error,
      timestamp: new Date().toISOString(),
      workspaceId,
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Retrieves personal access tokens (PATs).
 *
 * @param {string} token - The authorization token.
 * @return {Promise<ApiResponse<Pat[]>>} A promise resolving with PAT data or an error.
 */
export async function getPats(token: string): Promise<ApiResponse<Pat[]>> {
  try {
    const response = await fetch(`${import.meta.env.VITE_ZYG_URL}/pats/`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      method: "GET",
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch account PATs",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const pats = data.map((item: any) => patSchema.parse(item));
    return { data: pats, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid account PAT schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch account PATs";

    console.error("[getPats]", { error, timestamp: new Date().toISOString() });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Creates a new personal access token (PAT).
 *
 * @param {string} token - The authorization token.
 * @param {object} body - Object containing PAT name and description.
 * @return {Promise<ApiResponse<Pat>>} A promise resolving with newly created PAT data or an error.
 */
export async function createPat(
  token: string,
  body: { description: string; name: string },
): Promise<ApiResponse<Pat>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to create account PAT",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const pat = patSchema.parse(data);
    return { data: pat, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid account PAT schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to create account PAT";

    console.error("[createPat]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}








/**
 * Gets an existing ZYG account or creates a new one.
 *
 * @param {string} token - The authorization token.
 * @param {object} body - Object containing possible account name.
 * @return {Promise<ApiResponse<Account>>} A promise resolving with account data or an error.
 */
export async function getOrCreateZygAccount(
  token: string,
  body?: { name: string },
): Promise<ApiResponse<Account>> {
  try {
    const reqBody = body ? JSON.stringify(body) : "{}";
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/accounts/auth/`,
      {
        body: reqBody,
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "POST",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message:
          errorData?.message || "Failed to create or get ZYG auth account",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const account = accountSchema.parse(data);
    return { data: account, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid account schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to create or get ZYG auth account";

    console.error("[getOrCreateZygAccount]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Retrieves workspace customers.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @return {Promise<ApiResponse<Customer[]>>} A promise resolving with customer data or an error.
 */
export async function getWorkspaceCustomers(
  token: string,
  workspaceId: string,
): Promise<ApiResponse<Customer[]>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch workspace customers",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const customers = data.map((item: any) => customerSchema.parse(item));
    return { data: customers, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid workspace customers schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch workspace customers";

    console.error("[getWorkspaceCustomers]", {
      error,
      timestamp: new Date().toISOString(),
      workspaceId,
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Creates a new workspace.
 *
 * @param {string} token - The authorization token.
 * @param {object} body - Containing workspace name.
 * @return {Promise<ApiResponse<{ workspaceId: string; workspaceName: string }>>} A promise resolving with workspace data or an error.
 */
export async function createWorkspace(
  token: string,
  body: { name: string },
): Promise<ApiResponse<{ workspaceId: string; workspaceName: string }>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to create workspace",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const { name, workspaceId } = await response.json();
    return {
      data: { workspaceId, workspaceName: name },
      error: null,
    };
  } catch (error) {
    const errorMessage =
      error instanceof HTTPError ? error.message : "Failed to create workspace";

    console.error("[createWorkspace]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Updates an existing workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID to update.
 * @param {object} body - The new workspace parameters.
 * @return {Promise<ApiResponse<{ workspaceId: string; workspaceName: string }>>} A promise resolving with the updated workspace data or an error.
 */
export async function updateWorkspace(
  token: string,
  workspaceId: string,
  body: { name: string },
): Promise<ApiResponse<{ workspaceId: string; workspaceName: string }>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PATCH",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to update workspace",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const { name, workspaceId: id } = await response.json();
    return {
      data: { workspaceId: id, workspaceName: name },
      error: null,
    };
  } catch (error) {
    const errorMessage =
      error instanceof HTTPError ? error.message : "Failed to update workspace";

    console.error("[updateWorkspace]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

// export async function getWorkspaceMembers(
//   token: string,
//   workspaceId: string,
// ): Promise<{
//   data: MemberResponse[] | null;
//   error: Error | null;
// }> {
//   try {
//     const response = await fetch(
//       `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/members/`,
//       {
//         headers: {
//           Authorization: `Bearer ${token}`,
//         },
//         method: "GET",
//       },
//     );
//
//     if (!response.ok) {
//       const { status, statusText } = response;
//       return {
//         data: null,
//         error: new Error(
//           `error fetching workspace members: ${status} ${statusText}`,
//         ),
//       };
//     }
//
//     try {
//       const data = await response.json();
//       console.log(data);
//       const members = data.map((item: any) => {
//         return memberResponseSchema.parse({ ...item });
//       });
//       return { data: members, error: null };
//     } catch (err) {
//       if (err instanceof z.ZodError) {
//         console.error(err.message);
//       } else console.error(err);
//       return {
//         data: null,
//         error: new Error("error parsing workspace members schema"),
//       };
//     }
//   } catch (err) {
//     console.error(err);
//     return {
//       data: null,
//       error: new Error(
//         "error fetching workspace members - something went wrong",
//       ),
//     };
//   }
// }

/**
 * Deletes a personal access token (PAT).
 *
 * @param {string} token - The authorization token.
 * @param {string} patId - The PAT ID.
 * @return {Promise<ApiResponse<null>>} A promise resolving with null or an error.
 */
export async function deletePat(
  token: string,
  patId: string,
): Promise<ApiResponse<null>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to delete PAT",
        status: response.status,
        statusText: response.statusText,
      });
    }

    return { data: null, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof HTTPError ? error.message : "Failed to delete PAT";

    console.error("[deletePat]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

// @sanchitrk:
// improve usage with react query with zustand store.
// can we do this better?
// export async function bootstrapWorkspace(
//   token: string,
//   workspaceId: string,
// ): Promise<IWorkspaceEntitiesBootstrap & IWorkspaceValueObjects> {
//   const data = initialWorkspaceData();
//
//   const getWorkspaceP = getWorkspace(token, workspaceId);
//   const getWorkspaceMemberP = getWorkspaceMember(token, workspaceId);
//   const getWorkspaceCustomersP = getWorkspaceCustomers(token, workspaceId);
//   const getWorkspaceMetricsP = getWorkspaceMetrics(token, workspaceId);
//   const getWorkspaceThreadsP = getWorkspaceThreads(token, workspaceId);
//   const getWorkspaceLabelsP = getWorkspaceLabels(token, workspaceId);
//   // const getWorkspaceMembersP = getWorkspaceMembers(token, workspaceId);
//   const getAccountPatsP = getPats(token);
//
//   const [
//     workspaceData,
//     memberData,
//     customerData,
//     metricsData,
//     threadsData,
//     labelsData,
//     // membersData,
//     patsData,
//   ] = await Promise.all([
//     getWorkspaceP,
//     getWorkspaceMemberP,
//     getWorkspaceCustomersP,
//     getWorkspaceMetricsP,
//     getWorkspaceThreadsP,
//     getWorkspaceLabelsP,
//     // getWorkspaceMembersP,
//     getAccountPatsP,
//   ]);
//
//   const { data: workspace, error: errWorkspace } = workspaceData;
//   const { data: member, error: errMember } = memberData;
//   const { data: customers, error: errCustomer } = customerData;
//   const { data: metrics, error: errMetrics } = metricsData;
//   const { data: threads, error: errThreads } = threadsData;
//   const { data: labels, error: errLabels } = labelsData;
//   const { data: pats, error: errPats } = patsData;
//
//   const hasErr =
//     errWorkspace ||
//     errMember ||
//     errCustomer ||
//     errMetrics ||
//     errThreads ||
//     errLabels ||
//     errPats;
//
//   if (hasErr) {
//     data.error = new Error("error bootstrapping workspace store information");
//     data.isPending = false;
//     return data;
//   }
//
//   if (workspace) {
//     data.workspace = workspace;
//     data.hasData = true;
//     data.isPending = false;
//   }
//
//   if (member) {
//     data.member = member;
//   }
//
//   if (customers && customers.length > 0) {
//     const customersMap = makeCustomersStoreable(customers);
//     data.customers = customersMap;
//   }
//
//   if (metrics) {
//     const { count } = metrics;
//     data.metrics = count;
//   }
//
//   // if (members && members.length > 0) {
//   //   const membersMap = makeMembersStoreable(members);
//   //   data.members = membersMap;
//   // }
//
//   if (threads && threads.length > 0) {
//     const threadsMap = makeThreadsStoreable(threads);
//     data.threads = threadsMap;
//   }
//
//   if (labels && labels.length > 0) {
//     const labelsMap = makeLabelsStoreable(labels);
//     data.labels = labelsMap;
//   }
//
//   if (pats && pats.length > 0) {
//     const patsMap = makePatsStoreable(pats);
//     data.pats = patsMap;
//   }
//
//   return data;
// }

/**
 * Fetches all messages for a specific thread within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} threadId - The thread ID.
 * @return {Promise<ApiResponse<ThreadMessageResponse[]>>} A promise resolving with thread messages data or an error.
 */
export async function getWorkspaceThreadMessages(
  token: string,
  workspaceId: string,
  threadId: string,
): Promise<ApiResponse<ThreadMessageResponse[]>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message:
          errorData?.message || "Failed to fetch workspace thread messages",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const messages = data.map((item: any) => threadMessageSchema.parse(item));
    return { data: messages, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid workspace thread messages schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch workspace thread messages";

    console.error("[getWorkspaceThreadMessages]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Creates a new label for a workspace.
 *
 * @param {string} token - Authorization token.
 * @param {string} workspaceId - ID of the workspace to create the label for.
 * @param {object} body - The new label parameters.
 * @return {Promise<ApiResponse<LabelResponse>>} A promise resolving with the newly created label data or an error.
 */
export async function createWorkspaceLabel(
  token: string,
  workspaceId: string,
  body: { name: string },
): Promise<ApiResponse<LabelResponse>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to create label",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const parsed = labelSchema.parse(data);
    return { data: parsed, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid workspace label schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to create label";

    console.error("[createWorkspaceLabel]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Updates a label for a specific workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} labelId - The label ID.
 * @param {object} body - The updated label parameters.
 * @return {Promise<ApiResponse<LabelResponse>>} A promise resolving with the updated label data or an error.
 */
export async function updateWorkspaceLabel(
  token: string,
  workspaceId: string,
  labelId: string,
  body: { name: string },
): Promise<ApiResponse<LabelResponse>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/labels/${labelId}/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PATCH",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to update label",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const parsed = labelSchema.parse(data);
    return { data: parsed, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid workspace label schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to update label";

    console.error("[updateWorkspaceLabel]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Updates a specific thread within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} threadId - The thread ID.
 * @param {object} body - The updated thread parameters.
 * @return {Promise<ApiResponse<ThreadResponse>>} A promise resolving with the updated thread data or an error.
 */
export async function updateThread(
  token: string,
  workspaceId: string,
  threadId: string,
  body: object,
): Promise<ApiResponse<ThreadResponse>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/${threadId}/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PATCH",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to update thread",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const parsed = threadSchema.parse(data);
    return { data: parsed, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid thread schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to update thread";

    console.error("[updateThread]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Fetches all labels for a specific thread within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} threadId - The thread ID.
 * @return {Promise<ApiResponse<ThreadLabelResponse[]>>} A promise resolving with thread labels data or an error.
 */
export async function getThreadLabels(
  token: string,
  workspaceId: string,
  threadId: string,
): Promise<ApiResponse<ThreadLabelResponse[]>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch thread labels",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const labels = data.map((item: any) => threadLabelSchema.parse(item));
    return { data: labels, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid thread labels schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch thread labels";

    console.error("[getThreadLabels]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Updates or adds a label to a specific thread within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} threadId - The thread ID.
 * @param {object} body - The label parameters.
 * @return {Promise<ApiResponse<ThreadLabelResponse>>} A promise resolving with the updated/added label data or an error.
 */
export async function putThreadLabel(
  token: string,
  workspaceId: string,
  threadId: string,
  body: { icon: string; name: string },
): Promise<ApiResponse<ThreadLabelResponse>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/${threadId}/labels/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PUT",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to set thread label",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const label = threadLabelSchema.parse(data);
    return { data: label, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid thread label schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to set thread label";

    console.error("[putThreadLabel]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Deletes a label from a specific thread within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} threadId - The thread ID.
 * @param {string} labelId - The ID of the label to be deleted.
 * @return {Promise<ApiResponse<boolean>>} A promise resolving with deletion status or an error.
 */
export async function deleteThreadLabel(
  token: string,
  workspaceId: string,
  threadId: string,
  labelId: string,
): Promise<ApiResponse<boolean>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to delete thread label",
        status: response.status,
        statusText: response.statusText,
      });
    }

    return { data: true, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof HTTPError
        ? error.message
        : "Failed to delete thread label";

    console.error("[deleteThreadLabel]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Sends a chat message to a specific thread within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} threadId - The thread ID.
 * @param {object} body - The chat message parameters.
 * @return {Promise<ApiResponse<ThreadMessageResponse>>} A promise resolving with the thread message data or an error.
 */
export async function sendThreadChatMessage(
  token: string,
  workspaceId: string,
  threadId: string,
  body: { message: string },
): Promise<ApiResponse<ThreadMessageResponse>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/chat/${threadId}/messages/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to send thread chat messages",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const parsed = threadMessageSchema.parse(data);
    return { data: parsed, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid thread chat messages schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to send thread chat messages";

    console.error("[sendThreadChatMessage]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Sends a mail message to a specific thread within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} threadId - The thread ID.
 * @param {object} body - The mail message parameters.
 * @return {Promise<ApiResponse<ThreadMessageResponse>>} A promise resolving with the thread message data or an error.
 */
export async function sendThreadMailMessage(
  token: string,
  workspaceId: string,
  threadId: string,
  body: { htmlBody: string },
): Promise<ApiResponse<ThreadMessageResponse>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/threads/email/${threadId}/messages/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to send thread mail messages",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const parsed = threadMessageSchema.parse(data);
    return { data: parsed, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid thread mail messages schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to send thread mail messages";

    console.error("[sendThreadMailMessage]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Fetches customer events for a specific customer within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} customerId - The customer ID.
 * @return {Promise<ApiResponse<CustomerEventResponse[]>>} A promise resolving with customer event data or an error.
 */
export async function getCustomerEvents(
  token: string,
  workspaceId: string,
  customerId: string,
): Promise<ApiResponse<CustomerEventResponse[]>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch customer events",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const events = data.map((item: any) => customerEventSchema.parse(item));
    return { data: events, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid customer event schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch customer events";

    console.error("[getCustomerEvents]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Fetches a specific attachment for a message within a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {string} messageId - The message ID.
 * @param {string} attachmentId - The attachment ID.
 * @return {Promise<ApiResponse<MessageAttachmentResponse>>} A promise resolving with message attachment data or an error.
 */
export async function getMessageAttachment(
  token: string,
  workspaceId: string,
  messageId: string,
  attachmentId: string,
): Promise<ApiResponse<MessageAttachmentResponse>> {
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
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch message attachment",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const attachment = messageAttachmentSchema.parse(data);
    return { data: attachment, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid message attachment schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch message attachment";

    console.error("[getMessageAttachment]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Creates an email setting for a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {object} body - The email setting parameters.
 * @return {Promise<ApiResponse<PostmarkMailServerSetting>>} A promise resolving with the created email setting data or an error.
 */
export async function createEmailSetting(
  token: string,
  workspaceId: string,
  body: { email: string },
): Promise<ApiResponse<PostmarkMailServerSetting>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/postmark/servers/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to create email setting",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const setting = postmarkMailServerSettingSchema.parse(data);
    return { data: setting, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid email setting schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to create email setting";

    console.error("[createEmailSetting]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Fetches the email setting for a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @return {Promise<ApiResponse<PostmarkMailServerSetting>>} A promise resolving with the email setting data or an error.
 */
export async function getEmailSetting(
  token: string,
  workspaceId: string,
): Promise<ApiResponse<PostmarkMailServerSetting>> {
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
        return { data: null, error: null };
      }

      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to fetch email setting",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const setting = postmarkMailServerSettingSchema.parse(data);
    return { data: setting, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid email setting schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to fetch email setting";

    console.error("[getEmailSetting]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Updates the email setting for a workspace.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {object} body - The email setting update parameters.
 * @return {Promise<ApiResponse<PostmarkMailServerSetting>>} A promise resolving with the updated email setting data or an error.
 */
export async function updateEmailSetting(
  token: string,
  workspaceId: string,
  body: { enabled?: boolean; hasForwardingEnabled?: boolean },
): Promise<ApiResponse<PostmarkMailServerSetting>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/postmark/servers/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PATCH",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to update email setting",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const setting = postmarkMailServerSettingSchema.parse(data);
    return { data: setting, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid email setting schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to update email setting";

    console.error("[updateEmailSetting]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Add an email domain to a workspace's email settings.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @param {object} body - The email domain to add.
 * @return {Promise<ApiResponse<PostmarkMailServerSetting>>} A promise resolving with the updated email setting data or an error.
 */
export async function addEmailDomain(
  token: string,
  workspaceId: string,
  body: { domain: string },
): Promise<ApiResponse<PostmarkMailServerSetting>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/postmark/servers/parts/dns/add/`,
      {
        body: JSON.stringify(body),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to add email domain",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const setting = postmarkMailServerSettingSchema.parse(data);
    return { data: setting, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid email setting schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to add email domain";

    console.error("[addEmailDomain]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}

/**
 * Verify DNS setting for a workspace's email settings.
 *
 * @param {string} token - The authorization token.
 * @param {string} workspaceId - The workspace ID.
 * @return {Promise<ApiResponse<PostmarkMailServerSetting>>} A promise resolving with the updated email setting data or an error.
 */
export async function verifyDNS(
  token: string,
  workspaceId: string,
): Promise<ApiResponse<PostmarkMailServerSetting>> {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/${workspaceId}/postmark/servers/parts/dns/verify/`,
      {
        body: JSON.stringify({}),
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "PUT",
      },
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new HTTPError({
        message: errorData?.message || "Failed to verify DNS",
        status: response.status,
        statusText: response.statusText,
      });
    }

    const data = await response.json();
    const setting = postmarkMailServerSettingSchema.parse(data);
    return { data: setting, error: null };
  } catch (error) {
    const errorMessage =
      error instanceof z.ZodError
        ? "Invalid email setting schema"
        : error instanceof HTTPError
          ? error.message
          : "Failed to verify DNS";

    console.error("[verifyDNS]", {
      error,
      timestamp: new Date().toISOString(),
    });
    return { data: null, error: new Error(errorMessage) };
  }
}
