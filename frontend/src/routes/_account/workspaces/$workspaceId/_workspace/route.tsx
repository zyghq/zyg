import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { WorkspaceSidebar } from "@/components/workspace/sidebar";
// import { memberRowToShape, membersToMap } from "@/db/shapes";
import { WorkspaceStoreState } from "@/db/store";
// import { MemberRow, syncMembersShape } from "@/db/sync";
import { useAccountStore, useWorkspaceStore } from "@/providers";
// import { useShape } from "@electric-sql/react";
import { createFileRoute, Outlet } from "@tanstack/react-router";
// import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace",
)({
  component: WorkspaceLayout,
});

function WorkspaceLayout() {
  // const { token } = Route.useRouteContext();
  // const { workspaceId: paramWorkspaceId } = Route.useParams();

  // const {
  //   data: memberRows,
  //   isError: isMembersError,
  //   isLoading: isMembersLoading,
  //   ...rest
  // } = useShape(syncMembersShape({ token, workspaceId: paramWorkspaceId })) as unknown as { data: MemberRow[]; isError: boolean; isLoading: boolean };
  //
  // const members = membersToMap(memberRows.map(memberRowToShape));

  // console.log('**** rest ***', rest)

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

  // const updateMembers = useStore(
  //   workspaceStore,
  //   (state) => state.updateMembers,
  // );

  // React.useEffect(() => {
  //   if (isMembersError || isMembersLoading) return;
  //   console.log("*************** useEffect: memberRows", memberRows);
  // }, [memberRows, updateMembers]);

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
