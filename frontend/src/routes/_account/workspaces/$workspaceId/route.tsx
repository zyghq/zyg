import {
  getPats,
  getWorkspace,
  getWorkspaceLabels,
  getWorkspaceMember,
  getWorkspaceMetrics,
  getWorkspaceThreads,
} from "@/db/api";
import {
  AuthMember,
  labelsToMap,
  patsToMap,
  threadsToMap,
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
} from "@/db/shapes";
import { LabelMap, PatMap, ThreadShapeMap, WorkspaceStoreState } from "@/db/store";
import {
  CustomerRow,
  CustomerRowUpdates,
  MemberRow,
  MemberRowUpdates,
  syncCustomersShape,
  syncMembersShape,
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
import React from "react";
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

/**
 * Fetches and processes thread data for a given workspace.
 *
 * @param {string} token - Authentication token used to access the workspace data.
 * @param {string} workspaceId - Identifier of the workspace to retrieve thread data for.
 * @return {Promise<null | ThreadShapeMap>} A promise that resolves to a map of threads if data is available, or null if no data is found.
 * @throws {Error} If an error occurs during the data retrieval process.
 */
async function getThreadsData(
  token: string,
  workspaceId: string,
): Promise<null | ThreadShapeMap> {
  const { data, error } = await getWorkspaceThreads(token, workspaceId);
  if (error) throw new Error(error.message);
  if (data && data.length > 0) {
    return threadsToMap(data);
  }
  return null;
}

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

const threadsQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    queryFn: async () => getThreadsData(token, workspaceId),
    queryKey: ["threads"],
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
  const workspaceStore = useWorkspaceStore();

  // queries
  const membersShapeOffset = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewMembersShapeOffset(state),
  ) as Offset;
  const membersShapeHandle = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewMembersShapeHandle(state),
  );
  const customersShapeOffset = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewCustomersShapeOffset(state),
  ) as Offset;
  const customersShapeHandle = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewCustomersShapeHandle(state),
  );

  // mutations
  const setMembersShapeOffset = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.setMembersShapeOffset,
  );
  const setMembersShapeHandle = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.setMembersShapeHandle,
  );

  const setCustomersShapeOffset = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.setCustomersShapeOffset,
  );
  const setCustomersShapeHandle = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.setCustomersShapeHandle,
  );

  const updateMember = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.updateMember,
  );
  const updateCustomer = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.updateCustomer,
  );

  const setInSync = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.setInSync,
  );

  const memberShape = React.useMemo(() => {
    const shapeDefaultOpts = syncMembersShape({ token, workspaceId });
    const shapeOpts = {
      ...shapeDefaultOpts,
      handle: membersShapeHandle || undefined,
      offset: membersShapeOffset,
    };
    return new ShapeStream<MemberRow>(shapeOpts);
  }, [token, workspaceId, membersShapeHandle, membersShapeOffset]);

  const customerShape = React.useMemo(() => {
    const shapeDefaultOpts = syncCustomersShape({ token, workspaceId });
    const shapeOpts = {
      ...shapeDefaultOpts,
      handle: customersShapeHandle || undefined,
      offset: customersShapeOffset,
    };
    return new ShapeStream<CustomerRow>(shapeOpts);
  }, [token, workspaceId, customersShapeHandle, customersShapeOffset]);

  const handleMemberSyncMessages = React.useCallback(
    (messages: Message<MemberRow>[]) => {
      messages.forEach((message) => {
        if (isChangeMessage(message) && message.value.member_id) {
          setInSync(false);
          if (message.headers.operation === "update") {
            const { value } = message;
            const updates = takeMemberUpdates(value as MemberRowUpdates);
            updateMember(updates);
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
    },
    [],
  );

  const handleCustomerSyncMessages = React.useCallback(
    (messages: Message<CustomerRow>[]) => {
      messages.forEach((message) => {
        if (isChangeMessage(message) && message.value.customer_id) {
          setInSync(false);
          if (message.headers.operation === "update") {
            const { value } = message;
            const updates = takeCustomerUpdates(value as CustomerRowUpdates);
            updateCustomer(updates);
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
    },
    [],
  );

  React.useEffect(() => {
    const unsubscribe = memberShape.subscribe(handleMemberSyncMessages);
    return () => {
      unsubscribe();
    };
  }, [memberShape]);

  React.useEffect(() => {
    const unsubscribe = customerShape.subscribe(handleCustomerSyncMessages);
    return () => {
      unsubscribe();
    };
  }, [customerShape]);

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
      queryClient.ensureQueryData(threadsQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(labelsQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(patsQueryOptions(token)),
      queryClient.ensureQueryData(syncMembersQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(
        syncCustomersQueryOptions(token, workspaceId),
      ),
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
  const { data: customersData } = useSuspenseQuery(
    syncCustomersQueryOptions(token, workspaceId),
  );
  const { data: metrics } = useSuspenseQuery(
    metricsQueryOptions(token, workspaceId),
  );
  const { data: threads } = useSuspenseQuery(
    threadsQueryOptions(token, workspaceId),
  );
  const { data: labels } = useSuspenseQuery(
    labelsQueryOptions(token, workspaceId),
  );
  const { data: pats } = useSuspenseQuery(patsQueryOptions(token));

  const { data: membersData } = useSuspenseQuery(
    syncMembersQueryOptions(token, workspaceId),
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
