import { stageIcon } from "@/components/icons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { threadStatusVerboseName } from "@/db/helpers";
import { Thread } from "@/db/models";
import { useWorkspaceStore } from "@/providers";
import { ChatBubbleIcon, PersonIcon } from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import { formatDistanceToNow } from "date-fns";
import { useStore } from "zustand";

export function ThreadLinkItem({
  thread,
  workspaceId,
}: {
  thread: Thread;
  workspaceId: string;
}) {
  const workspaceStore = useWorkspaceStore();
  const customerName = useStore(workspaceStore, (state) =>
    state.viewCustomerName(state, thread.customerId),
  );
  return (
    <Link
      className="grid grid-cols-custom-thread-list-default grid-rows-custom-thread-list-default gap-x-4 gap-y-2 border-b px-4 py-4 hover:bg-accent dark:hover:bg-accent xl:grid-cols-custom-thread-list-xl xl:grid-rows-custom-thread-list-xl xl:px-8"
      params={{ threadId: thread.threadId, workspaceId }}
      to={"/workspaces/$workspaceId/threads/$threadId"}
    >
      <div className="col-span-1 xl:col-span-1">
        <ChatBubbleIcon className="h-4 w-4 text-muted-foreground" />
      </div>
      <div className="col-span-1 xl:col-span-1">
        <div className="flex flex-col">
          <div className="text-xs font-medium sm:text-sm">{customerName}</div>
          <div className="text-xs text-muted-foreground"></div>
        </div>
      </div>
      <div className="col-span-1 xl:order-last xl:col-span-1">
        <div className="flex items-center justify-end gap-4">
          <div className="font-mono text-xs">
            {formatDistanceToNow(new Date(thread.createdAt), {
              addSuffix: true,
            })}
          </div>
          {thread.assigneeId ? (
            <Avatar className="h-5 w-5">
              <AvatarImage
                alt={thread.assigneeId}
                src={`https://avatar.vercel.sh/${thread.assigneeId}`}
              />
              <AvatarFallback>M</AvatarFallback>
            </Avatar>
          ) : (
            <PersonIcon className="h-4 w-4 text-muted-foreground" />
          )}
        </div>
      </div>
      <div className="col-span-3 xl:order-3 xl:col-span-1">
        <span className="flex overflow-hidden text-ellipsis whitespace-nowrap">
          <span className="break-words text-sm font-medium">
            {thread.title}
          </span>
          <span className="ml-2 max-w-xl truncate text-sm text-muted-foreground">
            {thread.previewText}
          </span>
        </span>
        <div className="mt-1 flex flex-wrap justify-start gap-1">
          <Badge
            className="border-indigo-200 bg-indigo-100 p-1 font-normal dark:border-indigo-600 dark:bg-indigo-700"
            variant="outline"
          >
            <span className="mr-1">
              {stageIcon(thread.stage, {
                className:
                  "w-4 h-4 text-indigo-500 dark:text-accent-foreground",
              })}
            </span>
            {threadStatusVerboseName(thread.stage)}
          </Badge>
        </div>
      </div>
    </Link>
  );
}
