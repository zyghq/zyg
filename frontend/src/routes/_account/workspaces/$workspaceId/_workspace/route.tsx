import { Header } from "@/components/workspace/header";
import SideNavLinks from "@/components/workspace/sidenav-links";
import { WorkspaceStoreState } from "@/db/store";
import { useAccountStore, useWorkspaceStore } from "@/providers";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace",
)({
  component: WorkspaceLayout,
});

function WorkspaceLayout() {
  const accountStore = useAccountStore();
  const workspaceStore = useWorkspaceStore();

  const email = useStore(accountStore, (state) => state.getEmail(state));
  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state),
  );
  const workspaceName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceName(state),
  );
  const memberId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getMemberId(state),
  );
  const metrics = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getMetrics(state),
  );

  const sort = useStore(workspaceStore, (state) =>
    state.viewThreadSortKey(state),
  );

  return (
    <React.Fragment>
      <Header
        email={email}
        memberId={memberId}
        metrics={metrics}
        workspaceId={workspaceId}
        workspaceName={workspaceName}
      />
      <div className="flex min-h-screen">
        <aside className="sticky top-14 hidden h-[calc(100vh-theme(spacing.14))] min-w-80 overflow-y-auto bg-neutral-50 dark:bg-inherit md:block md:border-r">
          <SideNavLinks
            email={email}
            maxHeight="h-[calc(100dvh-8rem)]"
            memberId={memberId}
            metrics={metrics}
            sort={sort}
            workspaceId={workspaceId}
            workspaceName={workspaceName}
          />
        </aside>
        <main className="mt-14 flex-1 pb-4">
          <Outlet />
        </main>
      </div>
    </React.Fragment>
  );
}
