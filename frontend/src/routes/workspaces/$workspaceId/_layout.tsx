import { z } from "zod";
import { createFileRoute, Outlet } from "@tanstack/react-router";

import { useStore } from "zustand";
import { WorkspaceStoreStateType } from "@/db/store";
import { Header } from "@/components/workspace/header";
import { SideNav } from "@/components/workspace/sidenav";

//
// for more: https://tanstack.com/router/latest/docs/framework/react/guide/search-params
// usage of `.catch` or `default` matters.

const threadFiltersSearchSchema = z.object({
  status: z.enum(["todo", "snoozed", "done"]).catch("todo"),
});

export const Route = createFileRoute("/workspaces/$workspaceId/_layout")({
  validateSearch: threadFiltersSearchSchema,
  component: () => <WorkspaceLayout />,
});

function WorkspaceLayout() {
  const { workspaceStore, AccountStore } = Route.useRouteContext();

  const email = useStore(AccountStore.useContext(), (state) =>
    state.getEmail(state)
  );

  const workspaceId = useStore(
    workspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getWorkspaceId(state)
  );
  const workspaceName = useStore(
    workspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getWorkspaceName(state)
  );

  const memberId = useStore(
    workspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getMemberId(state)
  );

  const metrics = useStore(
    workspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getMetrics(state)
  );

  return (
    <div vaul-drawer-wrapper="">
      <div className="flex flex-col">
        <Header
          workspaceId={workspaceId}
          workspaceName={workspaceName}
          metrics={metrics}
          memberId={memberId}
        />
        <div className="flex flex-col">
          <div className="grid lg:grid-cols-5">
            <SideNav
              email={email}
              workspaceId={workspaceId}
              workspaceName={workspaceName}
              metrics={metrics}
              memberId={memberId}
            />
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}
