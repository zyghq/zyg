import { NotFound } from "@/components/notfound";
import { DetailsSheet } from "@/components/thread/details-sheet.tsx";
import { DetailsSidebar } from "@/components/thread/details-sidebar.tsx";
import { QueueSidebar } from "@/components/thread/queue-sidebar.tsx";
import { SummarySnippet } from "@/components/thread/summary-snippet";
import { ThreadContent } from "@/components/thread/thread-content";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
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
      <QueueSidebar activeThread={activeThread} workspaceId={workspaceId} />
      <SidebarInset>
        <header className="sticky top-0 flex shrink-0 flex-col gap-2 bg-background p-4 shadow-sm">
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
            <DetailsSheet
              activeThread={activeThread}
              token={token}
              workspaceId={workspaceId}
            />
          </div>
          <SummarySnippet />
        </header>
        <ThreadContent />
      </SidebarInset>
      <DetailsSidebar
        activeThread={activeThread}
        className="hidden lg:block"
        hide={detailSidebarHidden}
        token={token}
        workspaceId={workspaceId}
      />
    </SidebarProvider>
  );
}
