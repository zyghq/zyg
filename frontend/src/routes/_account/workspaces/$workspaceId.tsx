import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { queryOptions, useQuery } from "@tanstack/react-query";
import { bootstrapWorkspace } from "@/db/api";
import { WorkspaceStoreProvider } from "@/providers";

const bootstrapWorkspaceQueryOptions = (token: string, workspaceId: string) =>
  queryOptions({
    queryKey: ["workspaceStore", token, workspaceId],
    queryFn: async () => {
      const data = await bootstrapWorkspace(token, workspaceId);
      const { error } = data;
      if (error) throw new Error("failed to fetch workspace details.");
      return data;
    },
  });

export const Route = createFileRoute("/_account/workspaces/$workspaceId")({
  // check if we need this, add some kind of stale timer.
  // https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#using-staletime-to-control-how-long-data-is-considered-fresh
  loader: async ({
    context: { queryClient, supabaseClient },
    params: { workspaceId },
  }) => {
    const { error, data } = await supabaseClient.auth.getSession();
    if (error || !data?.session) throw redirect({ to: "/signin" });
    const token = data.session.access_token as string;
    return queryClient.ensureQueryData(
      bootstrapWorkspaceQueryOptions(token, workspaceId)
    );
  },
  component: Workspace,
});

function Workspace() {
  const { token } = Route.useRouteContext();
  const { workspaceId } = Route.useParams();
  const initialData = Route.useLoaderData();

  const response = useQuery({
    queryKey: ["workspaceStore", token, workspaceId],
    queryFn: async () => {
      const data = await bootstrapWorkspace(token, workspaceId);
      const { error } = data;
      if (error) throw new Error("failed to fetch workspace details.");
      return data;
    },
    initialData: initialData,
    enabled: !!token && !!workspaceId,
    staleTime: 1000 * 60 * 3,
  });

  const { data } = response;

  return (
    <WorkspaceStoreProvider initialValue={{ ...data }}>
      <Outlet />
    </WorkspaceStoreProvider>
  );
}
