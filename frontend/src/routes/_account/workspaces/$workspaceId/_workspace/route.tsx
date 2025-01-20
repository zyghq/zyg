import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { WorkspaceSidebar } from "@/components/workspace/sidebar";
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


  // There seems to be issue with sidebar-width-mobile having no effect:
  // see: https://github.com/shadcn-ui/ui/issues/5509

  return (
    <SidebarProvider
      style={
        {
          "--sidebar-width": "20rem",
          "--sidebar-width-mobile": "20rem",
        } as React.CSSProperties
      }
    >
      <WorkspaceSidebar
        email={email}
        memberId={memberId}
        metrics={metrics}
        sort={sort}
        workspaceId={workspaceId}
        workspaceName={workspaceName}
      />
      <SidebarInset>
        <Outlet />
      </SidebarInset>
    </SidebarProvider>
  );
}
