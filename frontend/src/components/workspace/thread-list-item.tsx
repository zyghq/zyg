import { Badge } from "@/components/ui/badge";
import { Thread } from "@/db/models";
import { useWorkspaceStore } from "@/providers";
import { ChatBubbleIcon, PersonIcon } from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import { formatDistanceToNow } from "date-fns";
import { useStore } from "zustand";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { threadStageHumanized } from "@/db/helpers";
import { stageIcon } from "@/components/icons";

export function ThreadLinkItem({
  thread,
  workspaceId,
}: {
  thread: Thread;
  workspaceId: string;
}) {
  const workspaceStore = useWorkspaceStore();
  const customerName = useStore(workspaceStore, (state) =>
    state.viewCustomerName(state, thread.customerId)
  );
  return (
    <Link
      className="grid grid-cols-custom-thread-list-default xl:grid-cols-custom-thread-list-xl grid-rows-custom-thread-list-default xl:grid-rows-custom-thread-list-xl px-4 py-4 border-b xl:px-8 gap-x-4 gap-y-2 hover:bg-zinc-50 dark:hover:bg-accent"
      params={{ threadId: thread.threadId, workspaceId }}
      to={"/workspaces/$workspaceId/threads/$threadId"}
    >
      <div className="col-span-1 xl:col-span-1">
        <ChatBubbleIcon className="w-4 h-4 text-muted-foreground" />
      </div>
      <div className="col-span-1 xl:col-span-1">
        <div className="flex flex-col">
          <div className="text-xs font-medium sm:text-sm">{customerName}</div>
          <div className="text-xs text-muted-foreground"></div>
        </div>
      </div>
      <div className="col-span-1 xl:col-span-1 xl:order-last">
        <div className="flex justify-end gap-4 items-center">
          <div className="text-xs font-mono">
            {formatDistanceToNow(new Date(thread.updatedAt), {
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
            <PersonIcon className="w-4 h-4 text-muted-foreground" />
          )}
        </div>
      </div>
      <div className="col-span-3 xl:col-span-1 xl:order-3">
        <span className="flex whitespace-nowrap overflow-hidden text-ellipsis">
          <span className="text-sm font-medium break-words">
            {thread.title}
          </span>
          <span className="text-sm ml-2 text-muted-foreground truncate">
            {thread.previewText}
          </span>
        </span>
        <div className="flex flex-wrap justify-start gap-1 mt-1">
          <Badge
            className="p-1 bg-indigo-100 font-normal border-indigo-200 dark:bg-indigo-700 dark:border-indigo-600"
            variant="outline"
          >
            <span className="mr-1">
              {stageIcon(thread.stage, {
                className:
                  "w-4 h-4 text-indigo-500 dark:text-accent-foreground",
              })}
            </span>
            {threadStageHumanized(thread.stage)}
          </Badge>
        </div>
      </div>
    </Link>
  );
}
