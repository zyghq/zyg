import { cn } from "@/lib/utils";
import { Link } from "@tanstack/react-router";
import { Badge } from "@/components/ui/badge";
import { formatDistanceToNow } from "date-fns";
import { Thread } from "@/db/models";
import { ChatBubbleIcon, ResetIcon } from "@radix-ui/react-icons";
import Avatar from "boring-avatars";
import { useStore } from "zustand";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";

function ThreadItem({
  workspaceId,
  item,
  variant = "default",
}: {
  workspaceId: string;
  item: Thread;
  variant?: string;
}) {
  // const WorkspaceStore = useRouteContext({
  //   from: "/_auth/workspaces/$workspaceId/_workspace",
  //   select: (context) => context.WorkspaceStore,
  // });
  const workspaceStore = useWorkspaceStore();
  const customerName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewCustomerName(state, item.customerId)
  );

  // const bottomRef = React.useRef<null | HTMLDivElement>(null);

  // React.useEffect(() => {
  //   if (bottomRef.current) {
  //     bottomRef.current.scrollIntoView({ behavior: "smooth" });
  //   }
  // }, []);

  return (
    <Link
      to={"/workspaces/$workspaceId/threads/$threadId"}
      params={{ workspaceId, threadId: item.threadId }}
      className={cn(
        "flex flex-col items-start gap-2 rounded-lg px-3 py-3 text-left text-sm transition-all hover:bg-accent",
        variant === "compress" && "gap-0 rounded-none py-5 border-b"
      )}
      activeOptions={{ exact: true }}
      activeProps={{
        className: "border-l-2 border-l-indigo-500 bg-indigo-50 dark:bg-accent",
      }}
    >
      <div className="flex w-full flex-col gap-1">
        <div className="flex items-center">
          <div className="flex items-center gap-2">
            <ChatBubbleIcon />
            <div className="font-semibold">{customerName}</div>
            {!item.read && (
              <span className="flex h-2 w-2 rounded-full bg-blue-600" />
            )}
          </div>
          <div
            className={cn(
              "ml-auto mr-2 text-xs",
              !item.replied ? "text-foreground" : "text-muted-foreground"
            )}
          >
            {formatDistanceToNow(new Date(item.updatedAt), {
              addSuffix: true,
            })}
          </div>
          {item.assigneeId && (
            <Avatar size={28} name={item.assigneeId} variant="marble" />
          )}
        </div>
        {item.replied ? (
          <div className="flex">
            <Badge variant="outline" className="font-normal">
              <div className="flex items-center gap-1">
                <ResetIcon className="h-3 w-3" />
                replied to
              </div>
            </Badge>
          </div>
        ) : (
          <div className="flex">
            <Badge
              variant="outline"
              className="bg-indigo-100 font-normal dark:bg-indigo-500"
            >
              <div className="flex items-center gap-1">
                <ResetIcon className="h-3 w-3" />
                awaiting reply
              </div>
            </Badge>
          </div>
        )}
        {/* {item?.title ? <div className="font-medium">{item?.title}</div> : null} */}
      </div>
      <div className="line-clamp-2 text-muted-foreground">
        {item.previewText}
      </div>
    </Link>
  );
}

export function ThreadList({
  workspaceId,
  threads,
  variant = "default",
}: {
  workspaceId: string;
  threads: Thread[];
  variant?: string;
}) {
  return (
    <div
      className={cn("flex flex-col gap-2", variant === "compress" && "gap-0")}
    >
      {threads.map((item: Thread) => (
        <ThreadItem
          key={item.threadId}
          workspaceId={workspaceId}
          item={item}
          variant={variant}
        />
      ))}
    </div>
  );
}
