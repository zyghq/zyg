import { createFileRoute, Link } from "@tanstack/react-router";
import { cn } from "@/lib/utils";
import React from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import {
  ArrowDownIcon,
  ArrowUpIcon,
  ChatBubbleIcon,
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
import { CurrentThreadQueueType, ThreadChatMessageType } from "@/db/store";
import {
  getWorkspaceThreadChatMessages,
  ThreadChatMessagesResponseType,
} from "@/db/api";
import { NotFound } from "@/components/notfound";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/threads/$threadId/"
)({
  component: ThreadDetail,
});

function getPrevNextFromCurrent(
  currentQueue: CurrentThreadQueueType | null,
  threadId: string
) {
  if (!currentQueue) return { prevItem: null, nextItem: null };

  const { threads } = currentQueue;
  const currentIndex = threads.findIndex(
    (thread) => thread.threadChatId === threadId
  );

  const prevItem = threads[currentIndex - 1] || null;
  const nextItem = threads[currentIndex + 1] || null;

  return { prevItem, nextItem };
}

function Message({
  message,
  memberId,
  memberName,
}: {
  message: ThreadChatMessageType;
  memberId: string;
  memberName: string;
}) {
  const { createdAt } = message;
  const date = new Date(createdAt);
  const time = date.toLocaleString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

  const isCustomer = message.customer ? true : false;
  const isMember = message.member ? true : false;

  const customerId = message.customer?.customerId || "C";

  const isMe = message.member?.memberId === memberId;

  return (
    <div className="flex">
      <div className={`flex ${isMe ? "ml-auto" : "mr-auto"}`}>
        <div className="flex space-x-2">
          {isCustomer && (
            <Avatar name={customerId} size={32} variant="marble" />
          )}
          <div className="p-2 rounded-lg bg-gray-100 dark:bg-accent">
            <div className="text-muted-foreground">{`${isMe ? "You" : memberName}`}</div>
            <p className="text-sm">{message.body}</p>
            <div className="flex text-xs justify-end text-muted-foreground mt-1">
              {time}
            </div>
          </div>
          {isMember && <Avatar name={memberId} size={32} variant="marble" />}
        </div>
      </div>
    </div>
  );
}

function ThreadDetail() {
  const { token } = Route.useRouteContext();
  const { workspaceId, threadId } = Route.useParams();
  const bottomRef = React.useRef<null | HTMLDivElement>(null);

  const workspaceStore = useWorkspaceStore();

  const currentQueue = useStore(
    workspaceStore,
    (state: WorkspaceStoreStateType) => state.viewCurrentThreadQueue(state)
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

  const memberId = useStore(workspaceStore, (state: WorkspaceStoreStateType) =>
    state.getMemberId(state)
  );
  const memberName = useStore(
    workspaceStore,
    (state: WorkspaceStoreStateType) => state.getMemberName(state)
  );

  const threadStatus = activeThread?.status || "";

  const { prevItem, nextItem } = getPrevNextFromCurrent(currentQueue, threadId);

  const { isPending, error, data } = useQuery({
    queryKey: ["messages", threadId, workspaceId, token],
    queryFn: async () => {
      const { error, data } = await getWorkspaceThreadChatMessages(
        token,
        workspaceId,
        threadId
      );
      if (error) throw new Error("failed to fetch thread messages");
      return data as ThreadChatMessagesResponseType;
    },
    enabled: !!activeThread,
  });

  React.useEffect(() => {
    if (bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [data]);

  function renderMessages(
    isPending: boolean,
    data?: ThreadChatMessagesResponseType
  ) {
    if (isPending) {
      return (
        <div className="flex justify-center mt-12">
          <svg
            className="animate-spin h-5 w-5 text-indigo-500"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            ></circle>
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
        </div>
      );
    }
    if (data && data.messages.length > 0) {
      const { messages } = data;
      const messagesReversed = Array.from(messages).reverse();
      return (
        <div className="p-4 space-y-2">
          {messagesReversed.map((message) => (
            <Message
              key={message.threadChatMessageId}
              message={message}
              memberId={memberId}
              memberName={memberName}
            />
          ))}
          <div ref={bottomRef}></div>
        </div>
      );
    }
    return (
      <div className="flex justify-center mt-12 text-muted-foreground">
        No results
      </div>
    );
  }

  if (!activeThread) {
    return <NotFound />;
  }

  if (error) {
    return <div>errror</div>;
  }

  return (
    <React.Fragment>
      <div className="flex min-h-screen">
        <aside
          className={cn(
            "sticky overflow-y-auto",
            currentQueue ? "border-r" : ""
          )}
        >
          <div className="flex">
            <div className="flex flex-col gap-4 px-2 py-4">
              <Button variant="outline" size="icon" asChild>
                <Link to={"/workspaces/$workspaceId"} params={{ workspaceId }}>
                  <ArrowLeftIcon className="h-4 w-4" />
                </Link>
              </Button>
              {currentQueue && (
                <SidePanelThreadList
                  threads={currentQueue.threads}
                  title="All Threads"
                  workspaceId={workspaceId}
                />
              )}
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
              className={cn("hidden", currentQueue ? "sm:block" : "")}
            >
              <div className="flex h-14 flex-col justify-center border-b px-4">
                <div className="font-semibold">{currentQueue?.from}</div>
              </div>
              <ScrollArea className="h-[calc(100dvh-4rem)]">
                {currentQueue && (
                  <ThreadList
                    workspaceId={workspaceId}
                    threads={currentQueue.threads}
                    variant="compress"
                  />
                )}
              </ScrollArea>
            </ResizablePanel>
            <ResizableHandle withHandle={true} />
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
                    <ScrollArea className="flex h-[calc(100dvh-4rem)] flex-col p-1">
                      {renderMessages(isPending, data)}
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
            <ResizableHandle withHandle={true} />
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
