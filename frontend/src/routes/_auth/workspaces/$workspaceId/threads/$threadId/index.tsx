import { createFileRoute, Link } from "@tanstack/react-router";
import { cn } from "@/lib/utils";
import React from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  ArrowDownIcon,
  ArrowUpIcon,
  ChatBubbleIcon,
  ArrowLeftIcon,
  ResetIcon,
  DotsHorizontalIcon,
  PlusIcon,
} from "@radix-ui/react-icons";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { CheckIcon } from "@radix-ui/react-icons";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import Avatar from "boring-avatars";
import { SidePanelThreadList } from "@/components/workspace/thread/sidepanel-thread-list";
import { useStore } from "zustand";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { ThreadList } from "@/components/workspace/thread/threads";
import { formatDistanceToNow } from "date-fns";
import {
  getWorkspaceThreadChatMessages,
  putThreadLabel,
  getThreadLabels,
} from "@/db/api";
import { NotFound } from "@/components/notfound";
import { PropertiesForm } from "@/components/workspace/thread/properties-form";
import { CheckCircleIcon, EclipseIcon, CircleIcon } from "lucide-react";
import { Icons } from "@/components/icons";
import { updateThread } from "@/db/api";
import { useMutation } from "@tanstack/react-query";
import { Thread, threadTransformer } from "@/db/models";
import { MessageForm } from "@/components/workspace/thread/message-form";
import {
  ThreadChatResponse,
  ThreadResponse,
  ThreadLabelResponse,
} from "@/db/schema";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/threads/$threadId/"
)({
  component: ThreadDetail,
});

function getPrevNextFromCurrent(threads: Thread[] | null, threadId: string) {
  if (!threads) return { prevItem: null, nextItem: null };

  const currentIndex = threads.findIndex(
    (thread) => thread.threadId === threadId
  );

  const prevItem = threads[currentIndex - 1] || null;
  const nextItem = threads[currentIndex + 1] || null;

  return { prevItem, nextItem };
}

function Chat({
  chat,
  memberId,
}: {
  chat: ThreadChatResponse;
  memberId: string;
}) {
  const { createdAt } = chat;
  const when = formatDistanceToNow(new Date(createdAt), {
    addSuffix: true,
  });

  const customerOrMemberId = chat.customer?.customerId || chat.member?.memberId;
  const customerOrMemberName = chat.customer?.name || chat.member?.name;

  const isMe = chat.member?.memberId === memberId;

  return (
    <div className="flex rounded-lg px-3 py-4 space-x-2 bg-white dark:bg-accent">
      <Avatar name={customerOrMemberId} size={32} variant="marble" />
      <div className="flex flex-col flex-1">
        <div className="flex justify-between">
          <div className="flex items-center">
            <div className="text-md font-semibold">
              {isMe ? `You` : customerOrMemberName}
            </div>
            <Separator orientation="vertical" className="mx-2 h-3" />
            <div className="text-muted-foreground text-xs">
              {`${when} via chat`}
            </div>
          </div>
          <Button variant="ghost" size="sm">
            <DotsHorizontalIcon className="h-4 w-4" />
          </Button>
        </div>
        <div className="text-xs text-muted-foreground"></div>
        <Separator
          orientation="horizontal"
          className="mt-3 mb-3 dark:bg-zinc-700"
        />
        <div>{chat.body}</div>
      </div>
    </div>
  );
}

function ThreadPreview({
  activeThread,
}: {
  activeThread: Thread;
}): JSX.Element {
  return (
    <div className="flex flex-col px-4 py-2">
      {activeThread.title && (
        <div className="font-semibold text-md">{activeThread.title}</div>
      )}
      {activeThread.previewText && (
        <div className="text-md text-muted-foreground">
          {activeThread.previewText}
        </div>
      )}
    </div>
  );
}

function SettingThreadLabel() {
  return (
    <div className="flex mr-1">
      <svg
        className="animate-spin h-3 w-3 text-indigo-500"
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

function ThreadLabels({
  token,
  workspaceId,
  threadId,
}: {
  token: string;
  workspaceId: string;
  threadId: string;
}) {
  const workspaceStore = useWorkspaceStore();
  const workspaceLabels = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewLabels(state)
  );

  const {
    isPending,
    data: threadLabels,
    error,
    refetch,
  } = useQuery({
    queryKey: ["threadLabels", workspaceId, threadId, token],
    queryFn: async () => {
      const { error, data } = await getThreadLabels(
        token,
        workspaceId,
        threadId
      );
      if (error) throw new Error("failed to fetch thread labels");
      return data as ThreadLabelResponse[];
    },
    enabled: !!threadId,
  });

  const threadLabelMutation = useMutation({
    mutationFn: async (values: { name: string; icon: string }) => {
      const { error, data } = await putThreadLabel(
        token,
        workspaceId,
        threadId,
        values
      );
      if (error) {
        throw new Error(error.message);
      }

      if (!data) {
        throw new Error("no data returned");
      }
      return data as ThreadLabelResponse;
    },
    onError: (error) => {
      console.error(error);
    },
    onSuccess: () => {
      refetch();
    },
  });

  const isChecked = (labelId: string) => {
    return threadLabels?.some((label) => label.labelId === labelId);
  };

  function onSelect(name: string, icon?: string) {
    threadLabelMutation.mutate({ name, icon: icon || "" });
  }

  const renderLabels = () => {
    if (isPending) {
      return null;
    }

    if (error) {
      return (
        <div className="flex items-center gap-1">
          <Icons.oops className="h-5 w-5" />
          <div className="text-xs text-red-500">Something went wrong</div>
        </div>
      );
    }

    return (
      <React.Fragment>
        {threadLabels?.map((label) => (
          <Badge key={label.labelId} variant="outline">
            <div className="flex items-center gap-1">
              <div>{label.icon}</div>
              <div className="text-muted-foreground capitalize">
                {label.name}
              </div>
            </div>
          </Badge>
        ))}
      </React.Fragment>
    );
  };

  return (
    <div className="flex flex-col px-4 py-2 gap-1">
      <div className="flex justify-between">
        <div className="text-muted-foreground font-semibold items-center">
          Labels
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm" className="border-dashed h-7">
              {threadLabelMutation.isPending ? (
                <SettingThreadLabel />
              ) : (
                <PlusIcon className="mr-1 h-3 w-3" />
              )}
              Add
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="sm:58 w-48" align="end">
            <Command>
              <CommandList>
                <CommandInput placeholder="Filter" />
                <CommandEmpty>No results</CommandEmpty>
                <CommandGroup>
                  {workspaceLabels.map((label) => (
                    <CommandItem
                      key={label.labelId}
                      onSelect={() => onSelect(label.name, label.icon)}
                      className="text-sm"
                    >
                      <div className="flex gap-2">
                        <div>{label.icon}</div>
                        <div className="capitalize">{label.name}</div>
                      </div>
                      <CheckIcon
                        className={cn(
                          "ml-auto h-4 w-4",
                          isChecked(label.labelId) ? "opacity-100" : "opacity-0"
                        )}
                      />
                    </CommandItem>
                  ))}
                </CommandGroup>
              </CommandList>
            </Command>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="flex gap-1 flex-wrap">{renderLabels()}</div>
      {threadLabelMutation.isError && (
        <div className="text-xs text-red-500">Something went wrong</div>
      )}
    </div>
  );
}

function ChatLoading() {
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

function ThreadDetail() {
  const { token } = Route.useRouteContext();
  const { workspaceId, threadId } = Route.useParams();

  const bottomRef = React.useRef<null | HTMLDivElement>(null);
  const workspaceStore = useWorkspaceStore();

  const currentQueue = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewCurrentThreadQueue(state)
  );

  const activeThread = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getThreadItem(state, threadId)
  );

  const customerName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewCustomerName(state, activeThread?.customerId || "")
  );

  const memberId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getMemberId(state)
  );

  const threadStatus = activeThread?.status || "";
  const isAwaitingReply = activeThread?.replied === false;

  const { prevItem, nextItem } = getPrevNextFromCurrent(currentQueue, threadId);

  const {
    isPending,
    error,
    data: chats,
    refetch,
  } = useQuery({
    queryKey: ["chats", threadId, workspaceId, token],
    queryFn: async () => {
      const { error, data } = await getWorkspaceThreadChatMessages(
        token,
        workspaceId,
        threadId
      );
      if (error) throw new Error("failed to fetch thread messages");
      return data as ThreadChatResponse[];
    },
    enabled: !!activeThread,
  });

  const assigneeId = activeThread?.assigneeId || "unassigned";
  const priority = activeThread?.priority || "normal";

  React.useEffect(() => {
    if (bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [chats]);

  const statusMutation = useMutation({
    mutationFn: async (values: { status: string }) => {
      const { error, data } = await updateThread(token, workspaceId, threadId, {
        ...values,
      });
      if (error) {
        throw new Error(error.message);
      }
      if (!data) {
        throw new Error("no data returned");
      }
      return data as ThreadResponse;
    },
    onError: (error) => {
      console.error(error);
    },
    onSuccess: (data) => {
      const transformer = threadTransformer();
      const [, thread] = transformer.normalize(data);
      workspaceStore.getState().updateThread(thread);
    },
  });

  const { isError: isStatusMutErr, isPending: isStatusMutPending } =
    statusMutation;

  function renderChats(chats?: ThreadChatResponse[]) {
    if (chats && chats.length > 0) {
      const chatsReversed = Array.from(chats).reverse();
      return (
        <div className="p-4 space-y-4">
          {chatsReversed.map((chat) => (
            <Chat key={chat.chatId} chat={chat} memberId={memberId} />
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
    return (
      <div className="flex flex-col container h-screen">
        <div className="my-auto mx-auto">
          <h1 className="mb-1 text-3xl font-bold">Error</h1>
          <p className="mb-4 text-red-500">
            There was an error fetching thread details. Try again later.
          </p>
        </div>
      </div>
    );
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
                  threads={currentQueue}
                  title="Threads"
                  workspaceId={workspaceId}
                />
              )}
              {prevItem ? (
                <Button variant="outline" size="icon" asChild>
                  <Link
                    to="/workspaces/$workspaceId/threads/$threadId"
                    params={{ workspaceId, threadId: prevItem.threadId }}
                  >
                    <ArrowUpIcon className="h-4 w-4" />
                  </Link>
                </Button>
              ) : null}
              {nextItem ? (
                <Button variant="outline" size="icon" asChild>
                  <Link
                    to="/workspaces/$workspaceId/threads/$threadId"
                    params={{ workspaceId, threadId: nextItem.threadId }}
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
                <div className="font-semibold">Threads</div>
              </div>
              <ScrollArea className="h-[calc(100dvh-4rem)]">
                {currentQueue && (
                  <ThreadList
                    workspaceId={workspaceId}
                    threads={currentQueue}
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
                        {isAwaitingReply && (
                          <React.Fragment>
                            <Separator
                              orientation="vertical"
                              className="mx-2"
                            />
                            <Badge
                              variant="outline"
                              className="bg-indigo-100 font-normal dark:bg-indigo-500"
                            >
                              <div className="flex items-center gap-1">
                                <ResetIcon className="h-2 w-2" />
                              </div>
                            </Badge>
                            <Separator
                              orientation="vertical"
                              className="mx-2"
                            />
                            <div className="text-xs">
                              {formatDistanceToNow(
                                new Date(activeThread.createdAt),
                                {
                                  addSuffix: true,
                                }
                              )}
                            </div>
                          </React.Fragment>
                        )}
                      </div>
                    </div>
                    <ScrollArea className="flex h-[calc(100dvh-4rem)] flex-col p-1 bg-gray-100 dark:bg-background">
                      {isPending ? <ChatLoading /> : renderChats(chats)}
                    </ScrollArea>
                  </div>
                </ResizablePanel>
                <ResizableHandle withHandle />
                <ResizablePanel defaultSize={25} maxSize={50} minSize={20}>
                  <div className="flex flex-col h-full p-2 overflow-auto gap-2">
                    <MessageForm
                      token={token}
                      workspaceId={workspaceId}
                      threadId={threadId}
                      customerName={customerName}
                      refetch={refetch}
                    />
                    <div className="flex flex-col mt-auto">
                      <div className="flex justify-end gap-2">
                        <Button size="sm" variant="outline">
                          <EclipseIcon className="mr-1 h-4 w-4 text-fuchsia-500" />
                          Snooze
                        </Button>
                        {threadStatus === "todo" && (
                          <Button
                            onClick={() => {
                              statusMutation.mutate({
                                status: "done",
                              });
                            }}
                            disabled={isStatusMutPending}
                            size="sm"
                            variant="outline"
                          >
                            <CheckCircleIcon className="mr-1 h-4 w-4 text-green-500" />
                            Mark as Done
                          </Button>
                        )}
                        {threadStatus === "done" && (
                          <Button
                            onClick={() => {
                              statusMutation.mutate({
                                status: "todo",
                              });
                            }}
                            disabled={isStatusMutPending}
                            size="sm"
                            variant="outline"
                          >
                            <CircleIcon className="mr-1 h-4 w-4 text-indigo-500" />
                            Mark as Todo
                          </Button>
                        )}
                      </div>
                      <div className="flex justify-end">
                        {isStatusMutErr && (
                          <div className="text-xs text-red-500">
                            Something went wrong.
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                </ResizablePanel>
              </ResizablePanelGroup>
            </ResizablePanel>
            <ResizableHandle withHandle={true} />
            <ResizablePanel
              defaultSize={25}
              minSize={20}
              maxSize={30}
              className="hidden sm:block bg-gray-100 dark:bg-background p-2"
            >
              <div className="flex flex-col gap-2 bg-white dark:bg-background rounded-lg">
                <ThreadPreview activeThread={activeThread} />
                <PropertiesForm
                  token={token}
                  workspaceId={workspaceId as string}
                  threadId={threadId as string}
                  priority={priority}
                  assigneeId={assigneeId}
                />
                <ThreadLabels
                  token={token}
                  workspaceId={workspaceId}
                  threadId={threadId}
                />
              </div>
            </ResizablePanel>
          </ResizablePanelGroup>
        </main>
      </div>
    </React.Fragment>
  );
}
