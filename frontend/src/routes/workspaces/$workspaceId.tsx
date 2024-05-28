import { createFileRoute, Outlet } from "@tanstack/react-router";

import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";

import { bootstrapWorkspace } from "@/db/api";
import { WorkspaceStore } from "@/providers";

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

export const Route = createFileRoute("/workspaces/$workspaceId")({
  // attach the workspaceStore to the context
  // make children happy.
  beforeLoad: () => {
    return {
      // TODO: rename this
      workspaceStore: WorkspaceStore,
    };
  },
  // check if we need this, add some kind of stale timer.
  // https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#using-staletime-to-control-how-long-data-is-considered-fresh
  loader: async ({
    context: { queryClient, auth },
    params: { workspaceId },
  }) => {
    const { client } = auth;
    const { data } = await client.auth.getSession();
    const { session } = data;
    const token = session?.access_token || "";
    return queryClient.ensureQueryData(
      bootstrapWorkspaceQueryOptions(token, workspaceId)
    );
  },
  component: () => <WorkspaceContainer />,
});

function WorkspaceContainer() {
  const { auth } = Route.useRouteContext();
  const useAuth = auth.useContext();
  const session = useAuth?.session;
  const token = session?.access_token || "";
  const workspaceId = Route.useParams().workspaceId;
  const { data, isRefetching } = useSuspenseQuery(
    bootstrapWorkspaceQueryOptions(token, workspaceId)
  );

  // these state should be automatically handled by the router.
  // if (isPending) {
  //   console.log("loading...");
  //   return <div>Loading...</div>;
  // }

  // if (error) {
  //   console.log("error!!!!");
  //   return <div>Error: {error.message}</div>;
  // }

  if (isRefetching) console.log("handle refetching?");

  return (
    <WorkspaceStore.Provider initialValue={{ ...data }}>
      <Outlet />
    </WorkspaceStore.Provider>
  );
}
