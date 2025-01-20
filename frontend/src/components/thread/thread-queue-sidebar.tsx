import { FooterMenu } from "@/components/sidebarcommans";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { ThreadQueue } from "@/components/workspace/thread/thread-queue";
import { getFromLocalStorage, threadStatusVerboseName } from "@/db/helpers";
import { ThreadShape } from "@/db/shapes";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { Link } from "@tanstack/react-router";
import { ArrowDownIcon, ArrowLeftIcon, ArrowUpIcon } from "lucide-react";
import * as React from "react";
import { useStore } from "zustand";

type ThreadSidebarProps = React.ComponentProps<typeof Sidebar> & {
  activeThread: ThreadShape;
  workspaceId: string;
};

function getPrevNextFromCurrent(
  threads: null | ThreadShape[],
  threadId: string,
) {
  if (!threads) return { nextItem: null, prevItem: null };

  const currentIndex = threads.findIndex(
    (thread) => thread.threadId === threadId,
  );

  const prevItem = threads[currentIndex - 1] || null;
  const nextItem = threads[currentIndex + 1] || null;

  return { nextItem, prevItem };
}

export function ThreadQueueSidebar({
  activeThread,
  workspaceId,
  ...props
}: ThreadSidebarProps) {
  const workspaceStore = useWorkspaceStore();

  const sort = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewThreadSortKey(state),
  );

  const currentQueue = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewCurrentThreadQueue(state),
  );

  const { nextItem, prevItem } = getPrevNextFromCurrent(
    currentQueue,
    activeThread?.threadId,
  );

  const stage = activeThread.stage;

  const threadsPath =
    getFromLocalStorage("zyg:threadsQueuePath") ||
    "/workspaces/$workspaceId/threads/todo";

  function renderQ(inQueue: null | ThreadShape[], active: null | ThreadShape) {
    if (inQueue && inQueue.length)
      return <ThreadQueue threads={inQueue} workspaceId={workspaceId} />;
    else if (active) {
      const items = [];
      items.push(active);
      return <ThreadQueue threads={items} workspaceId={workspaceId} />;
    }
    return null;
  }

  return (
    <Sidebar
      className="overflow-hidden [&>[data-sidebar=sidebar]]:flex-row"
      collapsible="icon"
      {...props}
    >
      {/* This is the first sidebar */}
      {/* We disable collapsible and adjust width to icon. */}
      {/* This will make the sidebar appear as icons. */}
      <Sidebar
        className="hidden !w-[calc(var(--sidebar-width-icon)_+_1px)] border-r md:flex"
        collapsible="none"
      >
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton asChild size="default" variant="outline">
                <Link
                  params={{ workspaceId }}
                  search={{ sort }}
                  to={threadsPath as string}
                >
                  <ArrowLeftIcon className="h-4 w-4" />
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupContent className="px-1.5 md:px-0">
              <SidebarMenu>
                {prevItem ? (
                  <SidebarMenuItem>
                    <SidebarMenuButton
                      asChild
                      className="px-2.5 md:px-2"
                      tooltip={{
                        children: "Previous In Queue",
                        hidden: false,
                      }}
                      variant="outline"
                    >
                      <Link
                        params={{ threadId: prevItem.threadId, workspaceId }}
                        to="/workspaces/$workspaceId/threads/$threadId"
                      >
                        <ArrowUpIcon className="h-4 w-4" />
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ) : null}
                {nextItem ? (
                  <SidebarMenuItem>
                    <SidebarMenuButton
                      asChild
                      className="px-2.5 md:px-2"
                      tooltip={{
                        children: "Next In Queue",
                        hidden: false,
                      }}
                      variant="outline"
                    >
                      <Link
                        params={{ threadId: nextItem.threadId, workspaceId }}
                        to="/workspaces/$workspaceId/threads/$threadId"
                      >
                        <ArrowDownIcon className="h-4 w-4" />
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ) : null}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <FooterMenu />
        </SidebarFooter>
      </Sidebar>

      {/* This is the second sidebar */}
      {/* We disable collapsible and let it fill remaining space */}
      <Sidebar className="flex-1 md:flex" collapsible="none">
        <SidebarHeader className="h-14 border-b px-4">
          <div className="flex w-full items-center justify-between">
            <div>
              <div className="font-serif text-sm font-medium">Threads</div>
              <div className="text-xs text-muted-foreground">
                {threadStatusVerboseName(stage)}
              </div>
            </div>
            <div></div>
          </div>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup className="p-0">
            <SidebarGroupContent>
              {renderQ(currentQueue, activeThread)}
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
      </Sidebar>
    </Sidebar>
  );
}
