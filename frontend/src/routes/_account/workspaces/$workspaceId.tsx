import {
  bootstrapWorkspace,
  getPats,
  getWorkspace, getWorkspaceCustomers,
  getWorkspaceLabels,
  getWorkspaceMember,
  getWorkspaceMetrics,
  getWorkspaceThreads
} from "@/db/api";
import { AuthMember, customersToMap, labelsToMap, patsToMap, threadsToMap, Workspace, WorkspaceMetrics } from "@/db/models";
import { CustomerMap, IWorkspaceEntities, IWorkspaceValueObjects, LabelMap, PatMap, ThreadMap } from "@/db/store";
import { WorkspaceStoreProvider } from "@/providers";
import { queryOptions, useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";


/**
 * Retrieves workspace data based on the provided token and workspace ID.
 *
 * @param {string} token - The authentication token used to access the workspace data.
 * @param {string} workspaceId - The unique identifier of the workspace to retrieve.
 * @return {Promise<Workspace>} A promise that resolves to the workspace data.
 * @throws {Error} If there is an error retrieving the workspace data or the workspace does not exist.
 */
async function getWorkspaceData(token: string, workspaceId: string): Promise<Workspace> {
  const { data, error } = await getWorkspace(token, workspaceId);
  if (error) throw new Error(error.message)
  if (!data) throw new Error("the workspace does not exist.")
  return data// doesn't require any transformation
}


/**
 * Fetches authenticated member data for a specified workspace.
 *
 * @param {string} token - The authentication token used for verifying the request.
 * @param {string} workspaceId - The unique identifier of the workspace.
 * @return {Promise<AuthMember>} A promise that resolves to the authenticated member data.
 * @throws {Error} If an error occurs during the data retrieval process or if the member does not exist.
 */
async function getAuthMemberData(token: string, workspaceId: string): Promise<AuthMember> {
  const { data, error } = await getWorkspaceMember(token, workspaceId);
  if (error) throw new Error(error.message)
  if (!data) throw new Error("member does not exist")
  return data // doesn't require any transformation
}


/**
 * Fetches and transforms customer data for a given workspace.
 *
 * @param {string} token - The authorization token used for API requests.
 * @param {string} workspaceId - The unique identifier for the workspace whose customer data is being retrieved.
 * @return {Promise<CustomerMap | null>} A promise that resolves to a map of transformed customer data or null if no data is available.
 * @throws {Error} If there is an error while fetching the customer data.
 */
async function getCustomersData(token: string, workspaceId: string): Promise<CustomerMap | null> {
  const { data, error } = await getWorkspaceCustomers(token, workspaceId);
  if (error) throw new Error(error.message)
  if (data && data.length > 0) {
   return customersToMap(data) 
  }
  return null
}


async function getMetricsData(token: string, workspaceId: string): Promise<WorkspaceMetrics> {
  const { data, error } = await getWorkspaceMetrics(token, workspaceId);
  if (error) throw new Error(error.message)
  if (!data) throw new Error("workspace metrics does not exist")
  const { count } = data;
  return count
}

async function getThreadsData(token: string, workspaceId: string): Promise<null | ThreadMap> {
 const { data, error } = await getWorkspaceThreads(token, workspaceId); 
 if (error) throw new Error(error.message)
  if (data && data.length > 0) {
    return threadsToMap(data);
  }
  return null
}

async function getLabelsData(token: string, workspaceId: string): Promise<LabelMap | null> {
 const {data, error} = await getWorkspaceLabels(token, workspaceId); 
 if (error) throw new Error(error.message)
  if (data && data.length > 0) {
    return labelsToMap(data);
  }
  return null
}

async function getAccountPatsData(token: string): Promise<null | PatMap> {
  const { data, error} = await getPats(token);
  if (error) throw new Error(error.message)
  if (data && data.length > 0) {
    return patsToMap(data);
  }
  return null
}


// Tanstack query options
const workspaceQueryOptions = (token: string, workspaceId: string) => queryOptions({
  queryFn: async () => getWorkspaceData(token, workspaceId),
  queryKey: ['workspace']
});

const memberQueryOptions = (token: string, workspaceId: string) => queryOptions({
  queryFn: async () => getAuthMemberData(token, workspaceId),
  queryKey: ['member']
});

const customersQueryOptions = (token: string, workspaceId: string) => queryOptions({
  queryFn: async () => getCustomersData(token, workspaceId),
  queryKey: ['customers']
});

const metricsQueryOptions = (token: string, workspaceId: string) => queryOptions({
  queryFn: async () => getMetricsData(token, workspaceId),
  queryKey: ['metrics']
});

const threadsQueryOptions = (token: string, workspaceId: string) => queryOptions({
  queryFn: async () => getThreadsData(token, workspaceId),
  queryKey: ['threads']
});

const labelsQueryOptions = (token: string, workspaceId: string) => queryOptions({
  queryFn: async () => getLabelsData(token, workspaceId),
  queryKey: ['labels']
});

const patsQueryOptions = (token: string) => queryOptions({
  queryFn: async () => getAccountPatsData(token),
  queryKey: ['pats']
});





// const bootstrapWorkspaceQueryOptions = (token: string, workspaceId: string) =>
//   queryOptions({
//     queryFn: async () => {
//       const data = await bootstrapWorkspace(token, workspaceId);
//       const { error } = data;
//       if (error) throw new Error("failed to fetch workspace details.");
//       return data;
//     },
//     queryKey: ["workspaceStore", token, workspaceId],
//   });

// async function bootstrapWorkspaceShape(
//   token: string,
//   workspaceId: string,
// ): Promise<WorkspaceShape> {
//   const shape = await preloadShape(syncWorkspaceShape({ token, workspaceId }));
//   const rows = (await shape.rows) as unknown as WorkspaceRow[];
//   const workspaces = rows.map((row) => workspaceRowToShape(row));
//   if (workspaces.length === 0) throw new Error("the workspace does not exist.");
//   return workspaces[0];
// }

// async function bootstrapMemberShape(
//   token: string,
//   workspaceId: string,
// ): Promise<MemberShapeMap> {
//   const shape = await preloadShape(
//     syncWorkspaceMemberShape({ token, workspaceId }),
//   );
//   const rows = (await shape.rows) as unknown as MemberRow[];
//   const members = rows.map((row) => memberRowToShape(row));
//   return membersToMap(members);
// }

export const Route = createFileRoute("/_account/workspaces/$workspaceId")({
  component: WorkspaceComponent,
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
      queryClient.ensureQueryData(customersQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(metricsQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(threadsQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(labelsQueryOptions(token, workspaceId)),
      queryClient.ensureQueryData(patsQueryOptions(token)),
    ])
  },
});

function WorkspaceComponent() {
  const { token } = Route.useRouteContext();
  const { workspaceId } = Route.useParams();
  const { data: workspace } = useSuspenseQuery(workspaceQueryOptions(token, workspaceId))
  const { data: member } = useSuspenseQuery(memberQueryOptions(token, workspaceId))
  const {data: customers } = useSuspenseQuery(customersQueryOptions(token, workspaceId))
  const {data: metrics} = useSuspenseQuery(metricsQueryOptions(token, workspaceId))
  const { data: threads} = useSuspenseQuery(threadsQueryOptions(token, workspaceId))
  const { data: labels } = useSuspenseQuery(labelsQueryOptions(token, workspaceId))
  const { data: pats } = useSuspenseQuery(patsQueryOptions(token))










  // const initialData = Route.useLoaderData();
  // const { members, workspace, ...rest } = initialData;

  // console.log("*************** initialData **************", initialData);

  // const response = useQuery({
  //   enabled: !!token && !!workspaceId,
  //   initialData: rest,
  //   queryFn: async () => {
  //     const data = await bootstrapWorkspace(token, workspaceId);
  //     const { error } = data;
  //     if (error) throw new Error("failed to fetch workspace details.");
  //     return data;
  //   },
  //   queryKey: ["workspaceStore", token, workspaceId],
  //   staleTime: 1000 * 60 * 3,
  // });
  // const { data } = response;

  return <div>....</div>

  // return (
  //   <WorkspaceStoreProvider initialValue={{ ...data, members, workspace }}>
  //     <Outlet />
  //   </WorkspaceStoreProvider>
  // );
}
