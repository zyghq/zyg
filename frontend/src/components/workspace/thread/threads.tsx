import { channelIcon, stageIcon } from "@/components/icons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { getInitials } from "@/db/helpers";
import { Thread } from "@/db/models";
import { WorkspaceStoreState } from "@/db/store";
import { cn } from "@/lib/utils";
import { useWorkspaceStore } from "@/providers";
import { PersonIcon } from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import { formatDistanceToNowStrict } from "date-fns";
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
        variant === "compress" && "gap-0 rounded-none py-2 border-b"
      )}
      params={{ threadId: item.threadId, workspaceId }}
      to={"/workspaces/$workspaceId/threads/$threadId"}
    >
      <div className="flex w-full flex-col gap-1">
        <div className="flex justify-between">
          <div className="flex items-center gap-2">
            <Avatar className="h-7 w-7">
              <AvatarImage
                src={`https://avatar.vercel.sh/${item.customerId}.svg?text=${getInitials(
                  customerName
                )}`}
              />
              <AvatarFallback>CS</AvatarFallback>
            </Avatar>
            <div className="font-medium">{customerName}</div>
          </div>
          {stageIcon(item.stage, {
            className: "w-4 h-4 text-indigo-500 my-auto",
          })}
        </div>
        <div className="font-semibold">{item.title}</div>
        <div className="line-clamp-2 text-muted-foreground">
          {item.previewText}
        </div>
      </div>
      <div className="flex flex-col w-full gap-2">
        <div className="flex justify-end w-full">
          <div className="flex gap-2">
            {channelIcon(item.channel, {
              className: "h-4 w-4 text-muted-foreground",
            })}
            <div className="text-xs">
              {formatDistanceToNowStrict(new Date(item.updatedAt), {
                addSuffix: true,
              })}
            </div>
            {item.assigneeId ? (
              <Avatar className="h-5 w-5">
                <AvatarImage
                  src={`https://avatar.vercel.sh/${item.assigneeId}`}
                />
                <AvatarFallback>M</AvatarFallback>
              </Avatar>
            ) : (
              <PersonIcon className="w-5 h-5 text-muted-foreground" />
            )}
          </div>
        </div>
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
