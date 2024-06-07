import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

import { queryOptions } from "@tanstack/react-query";
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

export const Route = createFileRoute("/_auth/workspaces/$workspaceId")({
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
  component: () => <WorkspaceContainer />,
});

function WorkspaceContainer() {
  const data = Route.useLoaderData();
  // const { token } = Route.useRouteContext();
  // const accountStore = useAccountStore();
  // const workspaceId = Route.useParams().workspaceId;

  // const token = useStore(accountStore, (state) => state.getToken(state));

  // const { data, isRefetching } = useSuspenseQuery(
  //   bootstrapWorkspaceQueryOptions(token, workspaceId)
  // );

  // these state should be automatically handled by the router.
  // if (isPending) {
  //   console.log("loading...");
  //   return <div>Loading...</div>;
  // }

  // if (error) {
  //   console.log("error!!!!");
  //   return <div>Error: {error.message}</div>;
  // }

  // if (isRefetching) console.log("handle refetching?");

  return (
    <WorkspaceStoreProvider initialValue={{ ...data }}>
      <Outlet />
    </WorkspaceStoreProvider>
  );
}
