import { NotFound } from "@/components/notfound";
import { ThreadContent } from "@/components/thread/thread-content";
import { ThreadDetailsSheet } from "@/components/thread/thread-details-sheet";
import { ThreadDetailsSidebar } from "@/components/thread/thread-details-sidebar.tsx";
import { ThreadQueueSidebar } from "@/components/thread/thread-queue-sidebar";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { WorkspaceStoreState } from "@/db/store.ts";
import { useWorkspaceStore } from "@/providers.tsx";
import { createFileRoute } from "@tanstack/react-router";
import { PanelRight } from "lucide-react";
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

  const [detailSidebarHidden, setDetailSidebarHidden] = React.useState(false);

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
          <div className="flex w-full items-center space-x-2">
            <SidebarTrigger aria-label="Toggle Sidebar" className="shrink-0" />
            <Separator className="h-4 shrink-0" orientation="vertical" />
            <div className="min-w-0 flex-1">
              <div className="truncate font-serif text-sm font-medium sm:text-base">
                {activeThread?.title || ""}
              </div>
            </div>
            <Button
              className="hidden h-7 w-7 shrink-0 md:flex"
              onClick={() => setDetailSidebarHidden(!detailSidebarHidden)}
              size="icon"
              variant="ghost"
            >
              <PanelRight />
            </Button>
            <ThreadDetailsSheet
              activeThread={activeThread}
              token={token}
              workspaceId={workspaceId}
            />
          </div>
        </header>
        <ThreadContent />
      </SidebarInset>
      <ThreadDetailsSidebar
        activeThread={activeThread}
        className="hidden lg:block"
        hide={detailSidebarHidden}
        token={token}
        workspaceId={workspaceId}
      />
    </SidebarProvider>
  );
}
