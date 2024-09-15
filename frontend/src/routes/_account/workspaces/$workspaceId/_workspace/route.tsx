import * as React from "react";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useStore } from "zustand";
import { WorkspaceStoreState } from "@/db/store";
import { Header } from "@/components/workspace/header";
import SideNavLinks from "@/components/workspace/sidenav-links";
import { useAccountStore, useWorkspaceStore } from "@/providers";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace"
)({
  component: WorkspaceLayout,
});

function WorkspaceLayout() {
  const accountStore = useAccountStore();
  const workspaceStore = useWorkspaceStore();

  const email = useStore(accountStore, (state) => state.getEmail(state));

  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state)
  );
  const workspaceName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceName(state)
  );

  const memberId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getMemberId(state)
  );

  const metrics = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getMetrics(state)
  );

  return (
    <React.Fragment>
      <Header
        email={email}
        workspaceId={workspaceId}
        workspaceName={workspaceName}
        metrics={metrics}
        memberId={memberId}
      />
      <div className="flex min-h-screen">
        <aside className="hidden sticky top-14 h-[calc(100vh-theme(spacing.14))] overflow-y-auto md:block md:border-r bg-zinc-50 dark:bg-inherit min-w-80">
          <SideNavLinks
            maxHeight="h-[calc(100dvh-8rem)]"
            email={email}
            workspaceId={workspaceId}
            workspaceName={workspaceName}
            metrics={metrics}
            memberId={memberId}
          />
        </aside>
        <main className="flex-1 mt-14 pb-4">
          <Outlet />
        </main>
      </div>
    </React.Fragment>
  );
}
