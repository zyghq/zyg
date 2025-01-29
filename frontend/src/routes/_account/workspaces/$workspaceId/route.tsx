import {
  getPats,
  getWorkspace,
  getWorkspaceLabels,
  getWorkspaceMember,
  getWorkspaceMetrics,
} from "@/db/api";
import {
  AuthMember,
  labelsToMap,
  patsToMap,
  Workspace,
  WorkspaceMetrics,
} from "@/db/models";
import {
  customerRowToShape,
  customersToMap,
  memberRowToShape,
  membersToMap,
  takeCustomerUpdates,
  takeMemberUpdates,
  takeThreadUpdates,
  threadRowToShape,
  threadsToMap,
} from "@/db/shapes";
import { LabelMap, PatMap } from "@/db/store";
import {
  CustomerRow,
  CustomerRowUpdates,
  MemberRow,
  MemberRowUpdates,
  syncCustomersShape,
  syncMembersShape,
  syncThreadsShape,
  ThreadRow,
  ThreadRowUpdates,
} from "@/db/sync";
import { useWorkspaceStore, WorkspaceStoreProvider } from "@/providers";
import {
  isChangeMessage,
  isControlMessage,
  Message,
  Offset,
  ShapeStream,
} from "@electric-sql/client";
import { preloadShape } from "@electric-sql/react";
import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import React, { useCallback, useEffect, useMemo, useTransition } from "react";
import { useStore } from "zustand";

/**
 * Retrieves workspace data based on the provided token and workspace ID.
 *
 * @param {string} token - The authentication token used to access the workspace data.
 * @param {string} workspaceId - The unique identifier of the workspace to retrieve.
 * @return {Promise<Workspace>} A promise that resolves to the workspace data.
 * @throws {Error} If there is an error retrieving the workspace data or the workspace does not exist.
 */
async function getWorkspaceData(
  token: string,
  workspaceId: string,
): Promise<Workspace> {
  const { data, error } = await getWorkspace(token, workspaceId);
  if (error) throw new Error(error.message);
  if (!data) throw new Error("the workspace does not exist.");
  return data; // doesn't require any transformation
}

/**
 * Fetches authenticated member data for a specified workspace.
 *
 * @param {string} token - The authentication token used for verifying the request.
 * @param {string} workspaceId - The unique identifier of the workspace.
 * @return {Promise<AuthMember>} A promise that resolves to the authenticated member data.
 * @throws {Error} If an error occurs during the data retrieval process or if the member does not exist.
 */
async function getAuthMemberData(
  token: string,
  workspaceId: string,
): Promise<AuthMember> {
  const { data, error } = await getWorkspaceMember(token, workspaceId);
  if (error) throw new Error(error.message);
  if (!data) throw new Error("member does not exist");
  return data; // doesn't require any transformation
}

// /**
//  * Fetches and transforms customer data for a given workspace.
//  *
//  * @param {string} token - The authorization token used for API requests.
//  * @param {string} workspaceId - The unique identifier for the workspace whose customer data is being retrieved.
//  * @return {Promise<CustomerMap | null>} A promise that resolves to a map of transformed customer data or null if no data is available.
//  * @throws {Error} If there is an error while fetching the customer data.
//  */
// async function getCustomersData(
//   token: string,
//   workspaceId: string,
// ): Promise<CustomerMap | null> {
//   const { data, error } = await getWorkspaceCustomers(token, workspaceId);
//   if (error) throw new Error(error.message);
//   if (data && data.length > 0) {
//     return customersToMap(data);
//   }
//   return null;
// }

/**
 * Retrieves metrics data for a specific workspace.
 *
 * @param {string} token - The authentication token used to access the workspace data.
 * @param {string} workspaceId - The unique identifier of the workspace for which metrics are being retrieved.
 * @return {Promise<WorkspaceMetrics>} A promise that resolves to the metrics data of the specified workspace.
 * @throws {Error} Throws an error if the retrieval fails or the metrics data does not exist.
 */
async function getMetricsData(
  token: string,
  workspaceId: string,
): Promise<WorkspaceMetrics> {
  const { data, error } = await getWorkspaceMetrics(token, workspaceId);
  if (error) throw new Error(error.message);
  if (!data) throw new Error("workspace metrics does not exist");
  const { count } = data;
  return count;
}

// /**
//  * Fetches and processes thread data for a given workspace.
//  *
//  * @param {string} token - Authentication token used to access the workspace data.
//  * @param {string} workspaceId - Identifier of the workspace to retrieve thread data for.
//  * @return {Promise<null | ThreadShapeMap>} A promise that resolves to a map of threads if data is available, or null if no data is found.
//  * @throws {Error} If an error occurs during the data retrieval process.
//  */
// async function getThreadsData(
//   token: string,
//   workspaceId: string,
// ): Promise<null | ThreadShapeMap> {
//   const { data, error } = await getWorkspaceThreads(token, workspaceId);
//   if (error) throw new Error(error.message);
//   if (data && data.length > 0) {
//     return threadsToMap(data);
//   }
//   return null;
// }

/**
 * Retrieves label data for a specific workspace and converts it into a map structure.
 *
 * @param {string} token - The authentication token required for the API request.
 * @param {string} workspaceId - The unique identifier for the targeted workspace.
 * @return {Promise<LabelMap | null>} A promise that resolves to a mapped representation of labels if data is available, or null if no data is found.
 * @throws {Error} Throws an error if the API request fails or returns an error.
 */
async function getLabelsData(
  token: string,
  workspaceId: string,
): Promise<LabelMap | null> {
  const { data, error } = await getWorkspaceLabels(token, workspaceId);
  if (error) throw new Error(error.message);
  if (data && data.length > 0) {
    return labelsToMap(data);
  }
  return null;
}

/**
 * Fetches and processes account PATs (Personal Access Tokens) data.
 *
 * @param {string} token - The authentication token used to fetch PATs data.
 * @return {Promise<null | PatMap>} Resolves to a map of PATs if data is available, otherwise resolves to null.
 * @throws {Error} Throws an error if fetching PATs data fails.
 */
async function getAccountPatsData(token: string): Promise<null | PatMap> {
  const { data, error } = await getPats(token);
  if (error) throw new Error(error.message);
  if (data && data.length > 0) {
    return patsToMap(data);
  }
  return null;
}

// Tanstack query options
const workspaceQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    queryFn: async () => getWorkspaceData(token, workspaceId),
    queryKey: ["workspace"],
  });

const memberQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    queryFn: async () => getAuthMemberData(token, workspaceId),
    queryKey: ["member"],
  });

const metricsQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    queryFn: async () => getMetricsData(token, workspaceId),
    queryKey: ["metrics"],
  });

const labelsQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    queryFn: async () => getLabelsData(token, workspaceId),
    queryKey: ["labels"],
  });

const patsQueryOptions = (token: string) =>
  queryOptions({
    queryFn: async () => getAccountPatsData(token),
    queryKey: ["pats"],
  });

// This is fetched from the sync engine.
const syncMembersQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    gcTime: Infinity,
    queryFn: async () => {
      const shape = await preloadShape(
        syncMembersShape({ token, workspaceId }),
      );
      const { handle, lastOffset } = shape;
      const rows = (await shape.rows) as unknown as MemberRow[];
      const members = rows.map((row) => memberRowToShape(row));
      const membersMap = membersToMap(members);
      return {
        handle: handle,
        members: membersMap,
        offset: lastOffset,
      };
    },
    queryKey: ["members"],
    staleTime: Infinity,
  });

const syncCustomersQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    gcTime: Infinity,
    queryFn: async () => {
      const shape = await preloadShape(
        syncCustomersShape({ token, workspaceId }),
      );
      const { handle, lastOffset } = shape;
      const rows = (await shape.rows) as unknown as CustomerRow[];
      const customers = rows.map((row) => customerRowToShape(row));
      const customersMap = customersToMap(customers);
      return {
        customers: customersMap,
        handle: handle,
        offset: lastOffset,
      };
    },
    queryKey: ["customers"],
    staleTime: Infinity,
  });

const syncThreadsQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    gcTime: Infinity,
    queryFn: async () => {
      const shape = await preloadShape(
        syncThreadsShape({ token, workspaceId }),
      );
      const { handle, lastOffset } = shape;
      const rows = (await shape.rows) as unknown as ThreadRow[];
      const threads = rows.map((row) => threadRowToShape(row));
      const threadsMap = threadsToMap(threads);
      return {
        handle: handle,
        offset: lastOffset,
        threads: threadsMap,
      };
    },
    queryKey: ["threads"],
    staleTime: Infinity,
  });

type ElectricSyncWrapperProps = {
  children: React.ReactNode;
  token: string;
  workspaceId: string;
};

function ElectricSyncWrapper({
  children,
  token,
  workspaceId,
}: ElectricSyncWrapperProps) {
  const [isPending, startTransition] = useTransition();
  const workspaceStore = useWorkspaceStore();

  // Queries with explicit state passing
  const membersShapeOffset = useStore(workspaceStore, (state) =>
    state.viewMembersShapeOffset(state),
  ) as Offset;
  const membersShapeHandle = useStore(workspaceStore, (state) =>
    state.viewMembersShapeHandle(state),
  );
  const customersShapeOffset = useStore(workspaceStore, (state) =>
    state.viewCustomersShapeOffset(state),
  ) as Offset;
  const customersShapeHandle = useStore(workspaceStore, (state) =>
    state.viewCustomersShapeHandle(state),
  );
  const threadsShapeOffset = useStore(workspaceStore, (state) =>
    state.viewThreadsShapeOffset(state),
  ) as Offset;
  const threadsShapeHandle = useStore(workspaceStore, (state) =>
    state.viewThreadsShapeHandle(state),
  );

  // Mutations
  const { setMembersShapeOffset, setMembersShapeHandle, updateMember } =
    useStore(workspaceStore, (state) => ({
      setMembersShapeOffset: state.setMembersShapeOffset,
      setMembersShapeHandle: state.setMembersShapeHandle,
      updateMember: state.updateMember,
    }));

  const { setCustomersShapeOffset, setCustomersShapeHandle, updateCustomer } =
    useStore(workspaceStore, (state) => ({
      setCustomersShapeOffset: state.setCustomersShapeOffset,
      setCustomersShapeHandle: state.setCustomersShapeHandle,
      updateCustomer: state.updateCustomer,
    }));

  const {
    setThreadsShapeOffset,
    setThreadsShapeHandle,
    updateThread,
    setInSync,
  } = useStore(workspaceStore, (state) => ({
    setThreadsShapeOffset: state.setThreadsShapeOffset,
    setThreadsShapeHandle: state.setThreadsShapeHandle,
    updateThread: state.updateThread,
    setInSync: state.setInSync,
  }));

  // Memoized shape streams
  const memberShape = useMemo(
    () =>
      new ShapeStream<MemberRow>({
        ...syncMembersShape({ token, workspaceId }),
        handle: membersShapeHandle || undefined,
        offset: membersShapeOffset,
      }),
    [token, workspaceId, membersShapeHandle, membersShapeOffset],
  );

  const customerShape = useMemo(
    () =>
      new ShapeStream<CustomerRow>({
        ...syncCustomersShape({ token, workspaceId }),
        handle: customersShapeHandle || undefined,
        offset: customersShapeOffset,
      }),
    [token, workspaceId, customersShapeHandle, customersShapeOffset],
  );

  const threadShape = useMemo(
    () =>
      new ShapeStream<ThreadRow>({
        ...syncThreadsShape({ token, workspaceId }),
        handle: threadsShapeHandle || undefined,
        offset: threadsShapeOffset,
      }),
    [token, workspaceId, threadsShapeHandle, threadsShapeOffset],
  );

  // Optimized sync handlers using `useTransition`
  const handleMemberSyncMessages = useCallback(
    (messages: Message<MemberRow>[]) => {
      startTransition(() => {
        messages.forEach((message) => {
          if (isChangeMessage(message) && message.value.member_id) {
            setInSync(false);
            if (message.headers.operation === "update") {
              updateMember(
                takeMemberUpdates(message.value as MemberRowUpdates),
              );
            }
          } else if (
            isControlMessage(message) &&
            message.headers.control === "up-to-date"
          ) {
            setInSync(true);
            if (memberShape.lastOffset && memberShape.shapeHandle) {
              setMembersShapeOffset(memberShape.lastOffset);
              setMembersShapeHandle(memberShape.shapeHandle as string);
            }
          }
        });
      });
    },
    [
      setInSync,
      setMembersShapeOffset,
      setMembersShapeHandle,
      updateMember,
      memberShape,
    ],
  );

  const handleCustomerSyncMessages = useCallback(
    (messages: Message<CustomerRow>[]) => {
      startTransition(() => {
        messages.forEach((message) => {
          if (isChangeMessage(message) && message.value.customer_id) {
            setInSync(false);
            if (message.headers.operation === "update") {
              updateCustomer(
                takeCustomerUpdates(message.value as CustomerRowUpdates),
              );
            }
          } else if (
            isControlMessage(message) &&
            message.headers.control === "up-to-date"
          ) {
            setInSync(true);
            if (customerShape.lastOffset && customerShape.shapeHandle) {
              setCustomersShapeOffset(customerShape.lastOffset);
              setCustomersShapeHandle(customerShape.shapeHandle as string);
            }
          }
        });
      });
    },
    [
      setInSync,
      setCustomersShapeOffset,
      setCustomersShapeHandle,
      updateCustomer,
      customerShape,
    ],
  );

  const handleThreadSyncMessages = useCallback(
    (messages: Message<ThreadRow>[]) => {
      startTransition(() => {
        messages.forEach((message) => {
          if (isChangeMessage(message) && message.value.thread_id) {
            setInSync(false);
            if (message.headers.operation === "update") {
              updateThread(
                takeThreadUpdates(message.value as ThreadRowUpdates),
              );
            }
          } else if (
            isControlMessage(message) &&
            message.headers.control === "up-to-date"
          ) {
            setInSync(true);
            if (threadShape.lastOffset && threadShape.shapeHandle) {
              setThreadsShapeOffset(threadShape.lastOffset);
              setThreadsShapeHandle(threadShape.shapeHandle as string);
            }
          }
        });
      });
    },
    [
      setInSync,
      setThreadsShapeOffset,
      setThreadsShapeHandle,
      updateThread,
      threadShape,
    ],
  );

  // Subscribe to shape streams
  useEffect(() => {
    return memberShape.subscribe(handleMemberSyncMessages);
  }, [memberShape]);

  useEffect(() => {
    return customerShape.subscribe(handleCustomerSyncMessages);
  }, [customerShape]);

  useEffect(() => {
    return threadShape.subscribe(handleThreadSyncMessages);
  }, [threadShape]);

  if (isPending) console.log("syncing data...");

  return <>{children}</>;
}

export const Route = createFileRoute("/_account/workspaces/$workspaceId")({
  component: CurrentWorkspace,
  // check if we need this, add some kind of stale timer.
  // https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#using-staletime-to-control-how-long-data-is-considered-fresh
  loader: async ({
    context: { queryClient, supabaseClient },
    params: { workspaceId },
  }) => {
    const { data, error } = await supabaseClient.auth.getSession();
    if (error || !data?.session) throw redirect({ to: "/signin" });
    const token = data.session.access_token as string;
    await Promise.all([
      queryClient.ensureQueryData(workspaceQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(memberQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(metricsQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(labelsQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(patsQueryOptions(token)),
      queryClient.ensureQueryData(syncMembersQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(
        syncCustomersQueryOptions(token, workspaceId),
      ),
      queryClient.ensureQueryData(syncThreadsQueryOptions(token, workspaceId)),
    ]);
    return token;
  },
});

function CurrentWorkspace() {
  const { token } = Route.useRouteContext();
  const { workspaceId } = Route.useParams();
  const { data: workspace } = useSuspenseQuery(
    workspaceQueryOptions(token, workspaceId),
  );
  const { data: member } = useSuspenseQuery(
    memberQueryOptions(token, workspaceId),
  );
  const { data: metrics } = useSuspenseQuery(
    metricsQueryOptions(token, workspaceId),
  );
  const { data: labels } = useSuspenseQuery(
    labelsQueryOptions(token, workspaceId),
  );
  const { data: pats } = useSuspenseQuery(patsQueryOptions(token));

  const { data: membersData } = useSuspenseQuery(
    syncMembersQueryOptions(token, workspaceId),
  );

  const { data: customersData } = useSuspenseQuery(
    syncCustomersQueryOptions(token, workspaceId),
  );

  const { data: threadsData } = useSuspenseQuery(
    syncThreadsQueryOptions(token, workspaceId),
  );

  const {
    handle: membersShapeHandle,
    members,
    offset: membersShapeOffset,
  } = membersData;

  const {
    customers,
    handle: customersShapeHandle,
    offset: customersShapeOffset,
  } = customersData;

  const {
    handle: threadsShapeHandle,
    offset: threadsShapeOffset,
    threads,
  } = threadsData;

  const initialData = {
    customers: customers,
    customersShapeHandle: customersShapeHandle || null,
    customersShapeOffset: customersShapeOffset || "-1",
    error: null,
    inSync: false,
    labels: labels,
    member: member,
    members: members,
    membersShapeHandle: membersShapeHandle || null,
    membersShapeOffset: membersShapeOffset || "-1",
    metrics: metrics,
    pats: pats,
    threadAppliedFilters: null,
    threads: threads,
    threadSortKey: null,
    threadsShapeHandle: threadsShapeHandle || null,
    threadsShapeOffset: threadsShapeOffset || "-1",
    workspace: workspace,
  };

  return (
    <WorkspaceStoreProvider initialValue={initialData}>
      <ElectricSyncWrapper token={token} workspaceId={workspaceId}>
        <Outlet />
      </ElectricSyncWrapper>
    </WorkspaceStoreProvider>
  );
}
