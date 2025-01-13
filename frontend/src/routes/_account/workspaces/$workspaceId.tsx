import { bootstrapWorkspace } from "@/db/api";
import {
  memberRowToShape,
  membersToMap,
  workspaceRowToShape,
  WorkspaceShape,
} from "@/db/shapes";
import {
  IWorkspaceEntities,
  IWorkspaceValueObjects,
  MemberShapeMap,
} from "@/db/store";
import {
  MemberRow,
  syncWorkspaceMemberShape,
  syncWorkspaceShape,
  WorkspaceRow,
} from "@/db/sync";
import { WorkspaceStoreProvider } from "@/providers";
import { preloadShape } from "@electric-sql/react";
import { queryOptions, useQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

const bootstrapWorkspaceQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    queryFn: async () => {
      const data = await bootstrapWorkspace(token, workspaceId);
      const { error } = data;
      if (error) throw new Error("failed to fetch workspace details.");
      return data;
    },
    queryKey: ["workspaceStore", token, workspaceId],
  });

async function bootstrapWorkspaceShape(
  token: string,
  workspaceId: string,
): Promise<WorkspaceShape> {
  const shape = await preloadShape(syncWorkspaceShape({ token, workspaceId }));
  const rows = (await shape.rows) as unknown as WorkspaceRow[];
  const workspaces = rows.map((row) => workspaceRowToShape(row));
  if (workspaces.length === 0) throw new Error("the workspace does not exist.");
  return workspaces[0];
}

async function bootstrapMemberShape(
  token: string,
  workspaceId: string,
): Promise<MemberShapeMap> {
  const shape = await preloadShape(
    syncWorkspaceMemberShape({ token, workspaceId }),
  );
  const rows = (await shape.rows) as unknown as MemberRow[];
  const members = rows.map((row) => memberRowToShape(row));
  return membersToMap(members);
}

export const Route = createFileRoute("/_account/workspaces/$workspaceId")({
  component: Workspace,
  // check if we need this, add some kind of stale timer.
  // https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#using-staletime-to-control-how-long-data-is-considered-fresh
  loader: async ({
    context: { queryClient, supabaseClient },
    params: { workspaceId },
  }) => {
    const { data, error } = await supabaseClient.auth.getSession();
    if (error || !data?.session) throw redirect({ to: "/signin" });
    const token = data.session.access_token as string;

    const workspace = await bootstrapWorkspaceShape(token, workspaceId);
    const members = await bootstrapMemberShape(token, workspaceId);
    const bootstrapped = queryClient.ensureQueryData(
      bootstrapWorkspaceQueryOptions(token, workspaceId),
    );
    return {
      ...bootstrapped,
      members,
      workspace,
    } as unknown as IWorkspaceEntities & IWorkspaceValueObjects;
  },
});

function Workspace() {
  const { token } = Route.useRouteContext();
  const { workspaceId } = Route.useParams();
  const initialData = Route.useLoaderData();
  const { members, workspace, ...rest } = initialData;

  const response = useQuery({
    enabled: !!token && !!workspaceId,
    initialData: rest,
    queryFn: async () => {
      const data = await bootstrapWorkspace(token, workspaceId);
      const { error } = data;
      if (error) throw new Error("failed to fetch workspace details.");
      return data;
    },
    queryKey: ["workspaceStore", token, workspaceId],
    staleTime: 1000 * 60 * 3,
  });
  const { data } = response;

  return (
    <WorkspaceStoreProvider initialValue={{ ...data, members, workspace }}>
      <Outlet />
    </WorkspaceStoreProvider>
  );
}
