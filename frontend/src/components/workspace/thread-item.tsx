import { Link } from "@tanstack/react-router";
import { Badge } from "@/components/ui/badge";
import { ChatBubbleIcon, PersonIcon } from "@radix-ui/react-icons";
import { LocateIcon } from "lucide-react";
import { Thread } from "@/db/models";
import { useStore } from "zustand";
import { formatDistanceToNow } from "date-fns";
import { useWorkspaceStore } from "@/providers";

export function ThreadLinkItem({
  workspaceId,
  thread,
}: {
  workspaceId: string;
  thread: Thread;
}) {
  const workspaceStore = useWorkspaceStore();
  const customerName = useStore(workspaceStore, (state) =>
    state.viewCustomerName(state, thread.customerId)
  );
  return (
    <Link
      to={"/workspaces/$workspaceId/threads/$threadId"}
      params={{ workspaceId, threadId: thread.threadId }}
      className="grid grid-cols-custom-thread-list-default sm:grid-cols-custom-thread-list-sm grid-rows-custom-thread-list-default sm:grid-rows-custom-thread-list-sm px-4 py-4 border-b sm:px-8 gap-x-4 gap-y-2 hover:bg-zinc-50 dark:hover:bg-accent"
    >
      <div className="col-span-1 sm:col-span-1">
        <ChatBubbleIcon className="w-4 h-4 text-muted-foreground" />
      </div>
      <div className="col-span-1 sm:col-span-1">
        <div className="flex flex-col">
          <div className="text-xs font-medium sm:text-sm">{customerName}</div>
          <div className="text-xs text-muted-foreground"></div>
        </div>
      </div>
      <div className="col-span-1 sm:col-span-1 sm:order-last">
        <div className="flex justify-end gap-4 items-center">
          <div className="text-xs">
            {formatDistanceToNow(new Date(thread.updatedAt), {
              addSuffix: true,
            })}
          </div>
          <PersonIcon className="w-4 h-4 text-muted-foreground" />
        </div>
      </div>
      <div className="col-span-3 sm:col-span-1 sm:order-3">
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
            variant="outline"
            className="p-1 bg-indigo-100 font-normal border-indigo-200 dark:bg-indigo-700 dark:border-indigo-600"
          >
            <span className="mr-1">
              <LocateIcon className="w-4 h-4 text-indigo-500 dark:text-accent-foreground" />
            </span>
            {"Need First Response"}
          </Badge>
        </div>
      </div>
    </Link>
  );
}
