import { NotFound } from "@/components/notfound";
import { ThreadActionsSidebar } from "@/components/thread/thread-actions-sidebar";
import { ThreadContent } from "@/components/thread/thread-content";
import { ThreadQueueSidebar } from "@/components/thread/thread-queue-sidebar";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { WorkspaceStoreState } from "@/db/store.ts";
import { useWorkspaceStore } from "@/providers.tsx";
import { createFileRoute } from "@tanstack/react-router";
import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/threads/$threadId",
)({
  component: ThreadDetailLayout,
});

function ThreadDetailLayout() {
  const { threadId, workspaceId } = Route.useParams();
  const { token } = Route.useRouteContext();

  const workspaceStore = useWorkspaceStore();
  const activeThread = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getThreadItem(state, threadId),
  );

  if (!activeThread) return <NotFound />;

  return (
    <SidebarProvider
      style={
        {
          "--sidebar-width": "450px",
          "--sidebar-width-mobile": "24rem",
        } as React.CSSProperties
      }
    >
      <ThreadQueueSidebar
        activeThread={activeThread}
        workspaceId={workspaceId}
      />
      <SidebarInset>
        <header className="sticky top-0 flex h-14 shrink-0 items-center gap-2 border-b bg-background p-4">
          <SidebarTrigger aria-label="Toggle Sidebar" className="-ml-1" />
          <Separator className="mr-2 h-4" orientation="vertical" />
          <div className="flex truncate font-serif font-medium">
            {activeThread?.title || ""}
          </div>
        </header>
        <ThreadContent />
      </SidebarInset>
      <ThreadActionsSidebar activeThread={activeThread} token={token} workspaceId={workspaceId} />
    </SidebarProvider>
  );
}
