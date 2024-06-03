import { createFileRoute, Outlet } from "@tanstack/react-router";

import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";
import { useStore } from "zustand";
import { bootstrapWorkspace } from "@/db/api";
import { WorkspaceStoreProvider, useAccountStore } from "@/providers";

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
  // attach the workspaceStore to the context
  // make children happy.
  // beforeLoad: () => {
  //   return {
  //     WorkspaceStore: WorkspaceStore,
  //   };
  // },
  // check if we need this, add some kind of stale timer.
  // https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#using-staletime-to-control-how-long-data-is-considered-fresh
  loader: async ({
    context: { queryClient, token },
    params: { workspaceId },
  }) => {
    return queryClient.ensureQueryData(
      bootstrapWorkspaceQueryOptions(token, workspaceId)
    );
  },
  component: () => <WorkspaceContainer />,
});

function WorkspaceContainer() {
  const accountStore = useAccountStore();
  const workspaceId = Route.useParams().workspaceId;

  const token = useStore(accountStore, (state) => state.getToken(state));

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
    <WorkspaceStoreProvider initialValue={{ ...data }}>
      <Outlet />
    </WorkspaceStoreProvider>
  );
}
