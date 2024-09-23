import { Icons } from "@/components/icons";
import { NotFound } from "@/components/notfound";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { MessageForm } from "@/components/workspace/thread/message-form";
import { PropertiesForm } from "@/components/workspace/thread/properties-form";
import { SidePanelThreadList } from "@/components/workspace/thread/sidepanel-thread-list";
import { ThreadList } from "@/components/workspace/thread/threads";
import {
  deleteThreadLabel,
  getThreadLabels,
  getWorkspaceThreadChatMessages,
  putThreadLabel,
} from "@/db/api";
import { updateThread } from "@/db/api";
import { Thread, threadTransformer } from "@/db/models";
import {
  ThreadChatResponse,
  ThreadLabelResponse,
  ThreadResponse,
} from "@/db/schema";
import { WorkspaceStoreState } from "@/db/store";
import { defaultSortKey } from "@/db/store";
import { cn } from "@/lib/utils";
import { useWorkspaceStore } from "@/providers";
import {
  ArrowDownIcon,
  ArrowLeftIcon,
  ArrowUpIcon,
  ChatBubbleIcon,
  DotsHorizontalIcon,
  PlusIcon,
  ResetIcon,
} from "@radix-ui/react-icons";
import { CheckIcon } from "@radix-ui/react-icons";
import { useQuery } from "@tanstack/react-query";
import { useMutation } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { formatDistanceToNow } from "date-fns";
import { CheckCircleIcon, CircleIcon, EclipseIcon } from "lucide-react";
import React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/threads/$threadId"
)({
  component: ThreadDetail,
});

function getPrevNextFromCurrent(threads: null | Thread[], threadId: string) {
  if (!threads) return { nextItem: null, prevItem: null };

  const currentIndex = threads.findIndex(
    (thread) => thread.threadId === threadId
  );

  const prevItem = threads[currentIndex - 1] || null;
  const nextItem = threads[currentIndex + 1] || null;

  return { nextItem, prevItem };
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
      <Avatar className="h-7 w-7">
        <AvatarImage src={`https://avatar.vercel.sh/${customerOrMemberId}`} />
        <AvatarFallback>{isMe ? "M" : "U"}</AvatarFallback>
      </Avatar>
      <div className="flex flex-col flex-1">
        <div className="flex justify-between">
          <div className="flex items-center">
            <div className="text-md font-semibold">
              {isMe ? `You` : customerOrMemberName}
            </div>
            <Separator className="mx-2 h-3" orientation="vertical" />
            <div className="text-muted-foreground text-xs">
              {`${when} via chat`}
            </div>
          </div>
          <Button size="sm" variant="ghost">
            <DotsHorizontalIcon className="h-4 w-4" />
          </Button>
        </div>
        <div className="text-xs text-muted-foreground"></div>
        <Separator
          className="mt-3 mb-3 dark:bg-zinc-700"
          orientation="horizontal"
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
        fill="none"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
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
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          fill="currentColor"
        ></path>
      </svg>
    </div>
  );
}

function ThreadLabels({
  threadId,
  token,
  workspaceId,
}: {
  threadId: string;
  token: string;
  workspaceId: string;
}) {
  const workspaceStore = useWorkspaceStore();
  const workspaceLabels = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewLabels(state)
  );

  const {
    data: threadLabels,
    error,
    isPending,
    refetch,
  } = useQuery({
    enabled: !!threadId,
    queryFn: async () => {
      const { data, error } = await getThreadLabels(
        token,
        workspaceId,
        threadId
      );
      if (error) throw new Error("failed to fetch thread labels");
      return data as ThreadLabelResponse[];
    },
    queryKey: ["threadLabels", workspaceId, threadId, token],
  });

  const threadLabelMutation = useMutation({
    mutationFn: async (values: { icon: string; name: string }) => {
      const { data, error } = await putThreadLabel(
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

  const deleteThreadLabelMutation = useMutation({
    mutationFn: async (labelId: string) => {
      const { data, error } = await deleteThreadLabel(
        token,
        workspaceId,
        threadId,
        labelId
      );
      if (error) {
        throw new Error(error.message);
      }

      if (!data) {
        throw new Error("no data returned");
      }
      return data;
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

  function onSelect(labelId: string, name: string, icon?: string) {
    if (isChecked(labelId)) {
      deleteThreadLabelMutation.mutate(labelId);
    } else {
      threadLabelMutation.mutate({ icon: icon || "", name });
    }
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
            <Button className="border-dashed h-7" size="sm" variant="outline">
              {threadLabelMutation.isPending ? (
                <SettingThreadLabel />
              ) : (
                <PlusIcon className="mr-1 h-3 w-3" />
              )}
              Add
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="sm:58 w-48">
            <Command>
              <CommandList>
                <CommandInput placeholder="Filter" />
                <CommandEmpty>No results</CommandEmpty>
                <CommandGroup>
                  {workspaceLabels.map((label) => (
                    <CommandItem
                      className="text-sm"
                      key={label.labelId}
                      onSelect={() =>
                        onSelect(label.labelId, label.name, label.icon)
                      }
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
        fill="none"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
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
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          fill="currentColor"
        ></path>
      </svg>
    </div>
  );
}

function ThreadDetail() {
  const { token } = Route.useRouteContext();
  const { threadId, workspaceId } = Route.useParams();

  const bottomRef = React.useRef<HTMLDivElement | null>(null);
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

  const { nextItem, prevItem } = getPrevNextFromCurrent(currentQueue, threadId);

  const {
    data: chats,
    error,
    isPending,
    refetch,
  } = useQuery({
    enabled: !!activeThread,
    queryFn: async () => {
      const { data, error } = await getWorkspaceThreadChatMessages(
        token,
        workspaceId,
        threadId
      );
      if (error) throw new Error("failed to fetch thread messages");
      return data as ThreadChatResponse[];
    },
    queryKey: ["chats", threadId, workspaceId, token],
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
      const { data, error } = await updateThread(token, workspaceId, threadId, {
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
            <Chat chat={chat} key={chat.chatId} memberId={memberId} />
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
              <Button asChild size="icon" variant="outline">
                <Link
                  params={{ workspaceId }}
                  search={{ sort: defaultSortKey }}
                  to={"/workspaces/$workspaceId/threads/todo"}
                >
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
                <Button asChild size="icon" variant="outline">
                  <Link
                    params={{ threadId: prevItem.threadId, workspaceId }}
                    to="/workspaces/$workspaceId/threads/$threadId"
                  >
                    <ArrowUpIcon className="h-4 w-4" />
                  </Link>
                </Button>
              ) : null}
              {nextItem ? (
                <Button asChild size="icon" variant="outline">
                  <Link
                    params={{ threadId: nextItem.threadId, workspaceId }}
                    to="/workspaces/$workspaceId/threads/$threadId"
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
              className={cn("hidden", currentQueue ? "sm:block" : "")}
              defaultSize={25}
              maxSize={30}
              minSize={20}
            >
              <div className="flex h-14 flex-col justify-center border-b px-4">
                <div className="font-semibold">Threads</div>
              </div>
              <ScrollArea className="h-[calc(100dvh-4rem)]">
                {currentQueue && (
                  <ThreadList
                    threads={currentQueue}
                    variant="compress"
                    workspaceId={workspaceId}
                  />
                )}
              </ScrollArea>
            </ResizablePanel>
            <ResizableHandle withHandle={true} />
            <ResizablePanel className="flex flex-col" defaultSize={50}>
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
                        <Separator className="mx-2" orientation="vertical" />
                        <ChatBubbleIcon className="h-3 w-3" />
                        {isAwaitingReply && (
                          <React.Fragment>
                            <Separator
                              className="mx-2"
                              orientation="vertical"
                            />
                            <Badge
                              className="bg-indigo-100 font-normal dark:bg-indigo-500"
                              variant="outline"
                            >
                              <div className="flex items-center gap-1">
                                <ResetIcon className="h-2 w-2" />
                              </div>
                            </Badge>
                            <Separator
                              className="mx-2"
                              orientation="vertical"
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
                      customerName={customerName}
                      refetch={refetch}
                      threadId={threadId}
                      token={token}
                      workspaceId={workspaceId}
                    />
                    <div className="flex flex-col mt-auto">
                      <div className="flex justify-end gap-2">
                        <Button size="sm" variant="outline">
                          <EclipseIcon className="mr-1 h-4 w-4 text-fuchsia-500" />
                          Snooze
                        </Button>
                        {threadStatus === "todo" && (
                          <Button
                            disabled={isStatusMutPending}
                            onClick={() => {
                              statusMutation.mutate({
                                status: "done",
                              });
                            }}
                            size="sm"
                            variant="outline"
                          >
                            <CheckCircleIcon className="mr-1 h-4 w-4 text-green-500" />
                            Mark as Done
                          </Button>
                        )}
                        {threadStatus === "done" && (
                          <Button
                            disabled={isStatusMutPending}
                            onClick={() => {
                              statusMutation.mutate({
                                status: "todo",
                              });
                            }}
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
              className="hidden sm:block bg-gray-100 dark:bg-background p-2"
              defaultSize={25}
              maxSize={30}
              minSize={20}
            >
              <div className="flex flex-col gap-2 bg-white dark:bg-background rounded-lg">
                <ThreadPreview activeThread={activeThread} />
                <PropertiesForm
                  assigneeId={assigneeId}
                  priority={priority}
                  threadId={threadId as string}
                  token={token}
                  workspaceId={workspaceId as string}
                />
                <ThreadLabels
                  threadId={threadId}
                  token={token}
                  workspaceId={workspaceId}
                />
              </div>
            </ResizablePanel>
          </ResizablePanelGroup>
        </main>
      </div>
    </React.Fragment>
  );
}
