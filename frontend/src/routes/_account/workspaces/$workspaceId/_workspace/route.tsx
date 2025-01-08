import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { Header } from "@/components/workspace/header";
import { WorkspaceSidebar} from "@/components/workspace/sidebar";
import SideNavLinks from "@/components/workspace/sidenav-links";
import { WorkspaceStoreState } from "@/db/store";
import { useAccountStore, useWorkspaceStore } from "@/providers";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace",
)({
  component: WorkspaceLayoutV2,
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
        sort={sort}
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

function WorkspaceLayoutV2() {
  const accountStore = useAccountStore();
  const workspaceStore = useWorkspaceStore();

  const email = useStore(accountStore, (state) => state.getEmail(state));
  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state),
  );
  const workspaceName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceName(state),
  );
  return (
    <SidebarProvider>
      <WorkspaceSidebar email={email} workspaceId={workspaceId} workspaceName={workspaceName}  />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="-ml-1" />
          <Separator className="mr-2 h-4" orientation="vertical" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem className="hidden md:block">
                <BreadcrumbLink href="#">
                  Building Your Application
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator className="hidden md:block" />
              <BreadcrumbItem>
                <BreadcrumbPage>Data Fetching</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </header>
        <Outlet />
        <div className="flex flex-1 flex-col gap-4 p-4">
          <div className="grid auto-rows-min gap-4 md:grid-cols-3">
            <div className="aspect-video rounded-xl bg-muted/50" />
            <div className="aspect-video rounded-xl bg-muted/50" />
            <div className="aspect-video rounded-xl bg-muted/50" />
          </div>
          <div className="min-h-[100vh] flex-1 rounded-xl bg-muted/50 md:min-h-min" />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
