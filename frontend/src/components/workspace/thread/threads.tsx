import { Badge } from "@/components/ui/badge";
import { Thread } from "@/db/models";
import { WorkspaceStoreState } from "@/db/store";
import { cn } from "@/lib/utils";
import { useWorkspaceStore } from "@/providers";
import { ChatBubbleIcon, ResetIcon } from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { formatDistanceToNow } from "date-fns";
import { useStore } from "zustand";

function ThreadItem({
  item,
  variant = "default",
  workspaceId,
}: {
  item: Thread;
  variant?: string;
  workspaceId: string;
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
      activeOptions={{ exact: true }}
      activeProps={{
        className: "border-l-2 border-l-indigo-500 bg-indigo-50 dark:bg-accent",
      }}
      className={cn(
        "flex flex-col items-start gap-2 rounded-lg px-3 py-3 text-left text-sm transition-all hover:bg-accent",
        variant === "compress" && "gap-0 rounded-none py-5 border-b"
      )}
      params={{ threadId: item.threadId, workspaceId }}
      to={"/workspaces/$workspaceId/threads/$threadId"}
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
            <Avatar className="h-5 w-5">
              <AvatarImage
                src={`https://avatar.vercel.sh/${item.assigneeId}`}
              />
              <AvatarFallback>M</AvatarFallback>
            </Avatar>
          )}
        </div>
        {item.replied ? (
          <div className="flex">
            <Badge className="font-normal" variant="outline">
              <div className="flex items-center gap-1">
                <ResetIcon className="h-3 w-3" />
                replied to
              </div>
            </Badge>
          </div>
        ) : (
          <div className="flex">
            <Badge
              className="bg-indigo-100 font-normal dark:bg-indigo-500"
              variant="outline"
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
  threads,
  variant = "default",
  workspaceId,
}: {
  threads: Thread[];
  variant?: string;
  workspaceId: string;
}) {
  return (
    <div
      className={cn("flex flex-col gap-2", variant === "compress" && "gap-0")}
    >
      {threads.map((item: Thread) => (
        <ThreadItem
          item={item}
          key={item.threadId}
          variant={variant}
          workspaceId={workspaceId}
        />
      ))}
    </div>
  );
}
