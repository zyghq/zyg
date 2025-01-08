import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { WorkspaceSidebar } from "@/components/workspace/sidebar";
import { WorkspaceStoreState } from "@/db/store";
import { useAccountStore, useWorkspaceStore } from "@/providers";
import { createFileRoute, Outlet } from "@tanstack/react-router";
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
    <SidebarProvider>
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
