import { stageIcon } from "@/components/icons";
import { channelIcon } from "@/components/icons";
import { priorityIcon } from "@/components/icons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { threadPriorityHumanized } from "@/db/helpers.ts";
import { ThreadShape } from "@/db/shapes";
import { useWorkspaceStore } from "@/providers";
import { PersonIcon } from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import { formatDistanceToNow } from "date-fns";
import { useStore } from "zustand";

export function ThreadLinkItem({
  thread,
  workspaceId,
}: {
  thread: ThreadShape;
  workspaceId: string;
}) {
  const workspaceStore = useWorkspaceStore();
  const customerName = useStore(workspaceStore, (state) =>
    state.viewCustomerName(state, thread.customerId),
  );
  return (
    <Link
      className="grid grid-cols-1 gap-2 border-b p-4 transition-colors duration-200 hover:bg-accent dark:hover:bg-accent sm:grid-cols-2 lg:grid-cols-5"
      params={{ threadId: thread.threadId, workspaceId }}
      to={"/workspaces/$workspaceId/threads/$threadId"}
    >
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between lg:col-span-1">
        <div className="flex items-center space-x-3">
          <div className="flex-shrink-0">
            {stageIcon(thread.stage, {
              className: "w-4 h-4 text-indigo-500 dark:text-accent-foreground",
            })}
          </div>
          <div className="min-w-0 flex-1">
            <p className="truncate text-sm font-medium">{customerName}</p>
            <p className="truncate text-xs text-muted-foreground lg:hidden">
              {thread.title}
            </p>
          </div>
        </div>
        <div className="mt-2 flex items-center sm:mt-0 lg:hidden">
          <div className="flex items-center gap-1">
            {priorityIcon(thread.priority, { className: "h-5 w-5" })}
            <span className="text-xs text-muted-foreground">
              {threadPriorityHumanized(thread.priority)}
            </span>
          </div>
        </div>
      </div>

      <div className="hidden lg:col-span-2 lg:block">
        <p className="truncate text-sm">
          <span className="font-medium">{thread.title}</span>
          <span className="ml-2 text-muted-foreground">
            {thread.previewText}
          </span>
        </p>
      </div>

      <div className="hidden lg:col-span-1 lg:flex lg:items-center lg:justify-center">
        <div className="flex-1"></div>
        <div className="flex flex-1 items-center gap-1">
          {priorityIcon(thread.priority, { className: "h-5 w-5" })}
          <div className="text-sm text-muted-foreground">
            {threadPriorityHumanized(thread.priority)}
          </div>
        </div>
      </div>

      <div className="mt-2 flex items-center justify-between space-x-3 sm:mt-0 sm:justify-end lg:col-span-1">
        <div className="flex items-center space-x-2 text-xs text-muted-foreground">
          {channelIcon(thread.channel, {
            className: "h-4 w-4",
          })}
          <span>
            {formatDistanceToNow(new Date(thread.createdAt), {
              addSuffix: true,
            })}
          </span>
        </div>
        {thread.assigneeId ? (
          <Avatar className="h-6 w-6">
            <AvatarImage
              alt={thread.assigneeId}
              src={`https://avatar.vercel.sh/${thread.assigneeId}`}
            />
            <AvatarFallback>
              {thread.assigneeId.charAt(0).toUpperCase()}
            </AvatarFallback>
          </Avatar>
        ) : (
          <PersonIcon className="h-5 w-5 text-muted-foreground" />
        )}
      </div>
    </Link>
  );
}
