import { priorityIcon } from "@/components/icons";
import { stageIcon } from "@/components/icons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { threadPriorityHumanized } from "@/db/helpers";
import { getInitials } from "@/db/helpers";
import { ThreadShape } from "@/db/shapes";
import { ThreadLabelShape, ThreadLabelShapeMap } from "@/db/shapes";
import { cn } from "@/lib/utils";
import { useWorkspaceStore } from "@/providers";
import { PersonIcon } from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import { formatDistanceToNow } from "date-fns";
import { useStore } from "zustand";

export function QueueItem({
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
  const customerEmail = useStore(workspaceStore, (state) =>
    state.viewCustomerEmail(state, thread.customerId),
  );

  const renderLabels = (labels: ThreadLabelShape[]) => {
    return labels.map((label) => (
      <Badge
        className="rounded-full bg-accent capitalize text-muted-foreground"
        key={label.labelId}
        variant="outline"
      >
        {label.name}
      </Badge>
    ));
  };

  return (
    <Link
      activeOptions={{ exact: true }}
      activeProps={{
        className: "bg-blue-50 dark:bg-accent",
      }}
      className={cn(
        "relative flex flex-col gap-2 border-b py-2 transition-colors hover:bg-muted/50",
      )}
      params={{ threadId: thread.threadId, workspaceId }}
      to={"/workspaces/$workspaceId/threads/$threadId"}
    >
      {({ isActive }: { isActive: boolean }) => (
        <>
          {isActive && (
            <div className="absolute bottom-0 left-0 top-0 w-[2px] bg-indigo-500"></div>
          )}
          <div className="flex flex-col items-start gap-2 px-3 sm:flex-row">
            <Avatar className="mt-1 h-7 w-7">
              <AvatarFallback>
                {getInitials(customerName) || "U"}
              </AvatarFallback>
              <AvatarImage
                alt={thread.customerId}
                src={`https://avatar.vercel.sh/${thread.customerId}`}
              />
            </Avatar>
            <div className="min-w-0 flex-1 space-y-2">
              <div className="flex flex-col justify-between sm:flex-row sm:items-center">
                <div>
                  <h3 className="truncate text-sm font-medium">
                    {customerName}
                  </h3>
                  <p className="text-xs text-muted-foreground">
                    {customerEmail}
                  </p>
                </div>
                {stageIcon(thread.stage, {
                  className:
                    "w-4 h-4 text-indigo-500 dark:text-accent-foreground",
                })}
              </div>
            </div>
          </div>
          <div className="px-3">
            <h4 className="truncate text-sm font-medium">{thread.title}</h4>
            <p className="truncate text-sm text-muted-foreground">
              {thread.previewText}
            </p>
          </div>
          <div className="px-3">
            <div className="flex flex-wrap items-center gap-2">
              {thread.labels &&
                renderLabels(
                  Object.values(thread.labels as ThreadLabelShapeMap),
                )}
            </div>
            <div className="flex flex-col justify-between gap-2 pt-2 sm:flex-row sm:items-center">
              <div className="flex items-center gap-1">
                {priorityIcon(thread.priority, { className: "h-5 w-5" })}
                <span>{threadPriorityHumanized(thread.priority)}</span>
              </div>
              <div className="flex items-center space-x-1">
                <div className="font-mono text-sm">
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
                  <PersonIcon className="h-5 w-5 text-muted-foreground" />
                )}
              </div>
            </div>
          </div>
        </>
      )}
    </Link>
  );
}
