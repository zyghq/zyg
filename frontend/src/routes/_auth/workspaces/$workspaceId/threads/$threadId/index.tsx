import { createFileRoute, Link } from "@tanstack/react-router";
import React from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import {
  ArrowDownIcon,
  ArrowUpIcon,
  ChatBubbleIcon,
  DotsHorizontalIcon,
  ResetIcon,
  ArrowLeftIcon,
} from "@radix-ui/react-icons";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import Avatar from "boring-avatars";
import { CircleIcon } from "lucide-react";
import { SidePanelThreadList } from "@/components/workspace/thread/sidepanel-thread-list";
import { useStore } from "zustand";
import { WorkspaceStoreStateType } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { ThreadList } from "@/components/workspace/thread/threads";
import { ThreadChatStoreType } from "@/db/store";
import { getWorkspaceThreadChatMessages } from "@/db/api";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/threads/$threadId/"
)({
  component: ThreadDetail,
});

function getPrevNextFromCurrent(
  threads: ThreadChatStoreType[],
  threadId: string
) {
  const currentIndex = threads.findIndex(
    (thread) => thread.threadChatId === threadId
  );

  const prevItem = threads[currentIndex - 1] || null;
  const nextItem = threads[currentIndex + 1] || null;

  return { prevItem, nextItem };
}

function ThreadDetail() {
  const { token } = Route.useRouteContext();
  const { workspaceId, threadId } = Route.useParams();
  const workspaceStore = useWorkspaceStore();

  const currentThreads = useStore(
    workspaceStore,
    (state: WorkspaceStoreStateType) => state.viewCurrentViewableThreads(state)
  );

  const activeThread = useStore(
    workspaceStore,
    (state: WorkspaceStoreStateType) => state.getThreadChatItem(state, threadId)
  );

  const customerName = useStore(
    workspaceStore,
    (state: WorkspaceStoreStateType) =>
      state.viewCustomerName(state, activeThread?.customerId || "")
  );

  const threadStatus = activeThread?.status || "n/a";

  const { prevItem, nextItem } = getPrevNextFromCurrent(
    currentThreads,
    threadId
  );

  const { isPending, error, data } = useQuery({
    queryKey: ["messages", threadId, workspaceId, token],
    queryFn: async () => {
      const { error, data } = await getWorkspaceThreadChatMessages(
        token,
        workspaceId,
        threadId
      );
      if (error) throw new Error("failed to fetch thread messages");
      return data;
    },
  });

  console.log("isPending", isPending);
  console.log("error", error);
  console.log("data", data);

  return (
    <React.Fragment>
      <div className="flex min-h-screen">
        <aside className="sticky overflow-y-auto md:border-r">
          <div className="flex">
            <div className="flex flex-col gap-4 px-2 py-4">
              <Button variant="outline" size="icon" asChild>
                <Link to={"/workspaces/$workspaceId"} params={{ workspaceId }}>
                  <ArrowLeftIcon className="h-4 w-4" />
                </Link>
              </Button>
              <SidePanelThreadList title="All Threads" />
              {prevItem ? (
                <Button variant="outline" size="icon" asChild>
                  <Link
                    to="/workspaces/$workspaceId/threads/$threadId"
                    params={{ workspaceId, threadId: prevItem.threadChatId }}
                  >
                    <ArrowUpIcon className="h-4 w-4" />
                  </Link>
                </Button>
              ) : null}
              {nextItem ? (
                <Button variant="outline" size="icon" asChild>
                  <Link
                    to="/workspaces/$workspaceId/threads/$threadId"
                    params={{ workspaceId, threadId: nextItem.threadChatId }}
                  >
                    <ArrowDownIcon className="h-4 w-4" />
                  </Link>
                </Button>
              ) : null}
            </div>
          </div>
        </aside>
        <main className="flex flex-col flex-1">
          <ResizablePanelGroup direction="horizontal">
            <ResizablePanel
              defaultSize={25}
              minSize={20}
              maxSize={30}
              className="hidden sm:block"
            >
              <div className="flex h-14 flex-col justify-center border-b px-4">
                <div className="font-semibold">All Threads</div>
              </div>
              <ScrollArea className="h-[calc(100dvh-4rem)]">
                <ThreadList
                  workspaceId={workspaceId}
                  threads={currentThreads}
                  variant="compress"
                  activeThread={activeThread}
                />
              </ScrollArea>
            </ResizablePanel>
            <ResizableHandle withHandle={false} />
            <ResizablePanel defaultSize={50} className="flex flex-col">
              <ResizablePanelGroup direction="vertical">
                <ResizablePanel defaultSize={75}>
                  <div className="flex h-full flex-col">
                    <div className="flex h-14 min-h-14 flex-col justify-center border-b px-4">
                      <div className="flex">
                        <div className="text-sm font-semibold">
                          {customerName}
                        </div>
                      </div>
                      <div className="flex items-center">
                        <CircleIcon className="mr-1 h-3 w-3 text-indigo-500" />
                        <span className="items-center text-xs capitalize">
                          {threadStatus}
                        </span>
                        <Separator orientation="vertical" className="mx-2" />
                        <ChatBubbleIcon className="h-3 w-3" />
                        {/* disabled for now, enable for something else perhaps? */}
                        {/* <Separator orientation="vertical" className="mx-2" />
                      <span className="font-mono text-xs">12/44</span> */}
                      </div>
                    </div>
                    <ScrollArea className="flex h-full flex-auto flex-col px-2 pb-4">
                      <div className="flex flex-col gap-1">
                        <div className="m-4">
                          <div className="flex items-center font-mono text-sm font-medium">
                            <span className="mr-1 flex h-1 w-1 rounded-full bg-fuchsia-500" />
                            Monday, 14 February 2024
                          </div>
                        </div>
                        {/* message */}
                        <div className="flex flex-col gap-2 rounded-lg border bg-background p-4">
                          <div className="flex w-full flex-col gap-1">
                            <div className="flex items-center">
                              <div className="flex items-center gap-2">
                                <Avatar
                                  size={28}
                                  name="name"
                                  variant="marble"
                                />
                                <div className="font-medium">Emily Davis</div>
                                <span className="flex h-1 w-1 rounded-full bg-blue-600" />
                                <span className="text-xs">3d ago.</span>
                              </div>
                              <div className="ml-auto">
                                <Button variant="ghost" size="icon">
                                  <DotsHorizontalIcon className="h-4 w-4" />
                                </Button>
                              </div>
                            </div>
                            <div className="font-medium">
                              {"Welcome To Plain."}
                            </div>
                          </div>
                          <div className="rounded-lg p-4 text-left text-muted-foreground hover:bg-accent">
                            {`Hi, welcome to Plain! We're so happy you're here! This
                      message is automated but if you need to talk to us you can
                      press the support button in the bottom left at anytime. In
                      the meantime let's use this thread to show you how Plain
                      works. 🌱 When customers reach out to you, they will show
                      up in your Plain workspace in a thread just like this. 🏷️
                      Each thread has a priority, an assignee, labels, and
                      Linear issues. Use the right-hand panel to set and change
                      those. ✉️ Reply by clicking into the composer below or
                      pressing R. Try sending a reply now. 🚀`}
                          </div>
                        </div>
                      </div>
                      <div className="flex flex-col gap-1">
                        <div className="m-4">
                          <div className="flex items-center font-mono text-sm font-medium">
                            <span className="mr-1 flex h-1 w-1 rounded-full bg-fuchsia-500" />
                            Monday, 14 February 2024
                          </div>
                        </div>
                        {/* message */}
                        <div className="flex flex-col gap-2 rounded-lg border bg-background p-4">
                          <div className="flex w-full flex-col gap-1">
                            <div className="flex items-center">
                              <div className="flex items-center gap-2">
                                <Avatar size={28} name="name" variant="beam" />
                                <div className="font-medium">Emily Davis</div>
                                <span className="flex h-1 w-1 rounded-full bg-blue-600" />
                                <span className="text-xs">3d ago.</span>
                              </div>
                              <div className="ml-auto">
                                <Button variant="ghost" size="icon">
                                  <DotsHorizontalIcon className="h-4 w-4" />
                                </Button>
                              </div>
                            </div>
                            <div className="flex items-center text-muted-foreground">
                              <ResetIcon className="mr-1 h-3 w-3" />
                              <div className="text-xs">Welcom to Plain...</div>
                            </div>
                          </div>
                          <div className="rounded-lg p-4 text-left text-muted-foreground hover:bg-accent">
                            {`Nice! 👀 You can see how these messages appear in Slack by pressing O. For more details you can check out our docs.✅ When you are done just hit "Mark as done" on the bottom right.
                      ⌨️ If you want to do anything in Plain use ⌘ + K or CTRL + K on Windows.
                      (edited)`}
                          </div>
                        </div>
                      </div>
                      <div className="flex flex-col gap-1">
                        <div className="m-4">
                          <div className="flex items-center font-mono text-sm font-medium">
                            <span className="mr-1 flex h-1 w-1 rounded-full bg-fuchsia-500" />
                            Monday, 14 February 2024
                          </div>
                        </div>
                        {/* message */}
                        <div className="flex flex-col gap-2 rounded-lg border bg-background p-4">
                          <div className="flex w-full flex-col gap-1">
                            <div className="flex items-center">
                              <div className="flex items-center gap-2">
                                <Avatar size={28} name="name" variant="beam" />
                                <div className="font-medium">Emily Davis</div>
                                <span className="flex h-1 w-1 rounded-full bg-blue-600" />
                                <span className="text-xs">3d ago.</span>
                              </div>
                              <div className="ml-auto">
                                <Button variant="ghost" size="icon">
                                  <DotsHorizontalIcon className="h-4 w-4" />
                                </Button>
                              </div>
                            </div>
                            <div className="flex items-center text-muted-foreground">
                              <ResetIcon className="mr-1 h-3 w-3" />
                              <div className="text-xs">Welcom to Plain...</div>
                            </div>
                          </div>
                          <div className="rounded-lg p-4 text-left text-muted-foreground hover:bg-accent">
                            {`Nice! 👀 You can see how these messages appear in Slack by pressing O. For more details you can check out our docs.✅ When you are done just hit "Mark as done" on the bottom right.
                      ⌨️ If you want to do anything in Plain use ⌘ + K or CTRL + K on Windows.
                      (edited)`}
                          </div>
                        </div>
                      </div>
                    </ScrollArea>
                  </div>
                </ResizablePanel>
                <ResizableHandle withHandle />
                <ResizablePanel defaultSize={25} maxSize={50} minSize={20}>
                  <div className="flex h-full items-center justify-center p-6">
                    <span className="font-semibold">Editor</span>
                  </div>
                </ResizablePanel>
              </ResizablePanelGroup>
            </ResizablePanel>
            <ResizableHandle withHandle={false} />
            <ResizablePanel
              defaultSize={25}
              minSize={20}
              maxSize={30}
              className="hidden sm:block"
            >
              <div className="flex h-full items-center justify-center p-6">
                <span className="font-semibold">Sidebar</span>
              </div>
            </ResizablePanel>
          </ResizablePanelGroup>
        </main>
      </div>
    </React.Fragment>
  );
}
