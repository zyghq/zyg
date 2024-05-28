import { cn } from "@/lib/utils";
import { Link } from "@tanstack/react-router";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";

import { formatDistanceToNow } from "date-fns";
import { ThreadChatStoreType } from "@/db/store";

import { ChatBubbleIcon, ResetIcon } from "@radix-ui/react-icons";
import { ReloadIcon } from "@radix-ui/react-icons";
import Avatar from "boring-avatars";

function ThreadItem({
  workspaceId,
  item,
  variant = "default",
}: {
  workspaceId: string;
  item: ThreadChatStoreType;
  variant?: string;
}) {
  //   const message = item.messages[0];
  //   const name = item?.customer?.name || "Customer";
  //   const { assignee } = item;

  //   const renderLabels = () => {
  //     if (result.isSuccess && result.data && result.data.length) {
  //       return (
  //         <div className="flex gap-1">
  //           {result.data.map((label) => (
  //             <Badge
  //               key={label.labelId}
  //               variant="outline"
  //               className="font-normal"
  //             >
  //               {label.name}
  //             </Badge>
  //           ))}
  //         </div>
  //       );
  //     }
  //     return <div className="min-h-5"></div>;
  //   };

  return (
    <Link
      to={"/workspaces/$workspaceId/threads/$threadId"}
      params={{ workspaceId, threadId: item.threadChatId }}
      className={cn(
        "flex flex-col items-start gap-2 rounded-lg border px-3 py-3 text-left text-sm transition-all hover:bg-accent",
        variant === "compress" && "gap-0 rounded-none py-5"
      )}
    >
      <div className="flex w-full flex-col gap-1">
        <div className="flex items-center">
          <div className="flex items-center gap-2">
            <ChatBubbleIcon />
            <div className="font-semibold">{item.customerId}</div>
            {!item.read && (
              <span className="flex h-2 w-2 rounded-full bg-blue-600" />
            )}
          </div>

          {item.recentMessage && (
            <div
              className={cn(
                "ml-auto mr-2 text-xs",
                !item.replied ? "text-foreground" : "text-muted-foreground"
              )}
            >
              {formatDistanceToNow(new Date(item.recentMessage.updatedAt), {
                addSuffix: true,
              })}
            </div>
          )}
          {item.assigneeId && (
            <Avatar
              size={28}
              name={item.assigneeId}
              variant="beam"
              colors={["#92A1C6", "#146A7C", "#F0AB3D", "#C271B4", "#C20D90"]}
            />
          )}
        </div>
        {item.replied ? (
          <div className="flex">
            <Badge variant="outline" className="font-normal">
              <div className="flex items-center gap-1">
                <ResetIcon className="h-3 w-3" /> replied to
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
        {item.recentMessage?.body}
      </div>
    </Link>
  );
}

export function ThreadList({
  workspaceId,
  threads,
  className,
  variant = "default",
}: {
  workspaceId: string;
  threads: ThreadChatStoreType[];
  className: string;
  variant?: string;
}) {
  return (
    <ScrollArea className={cn("pr-1", className)}>
      <div
        className={cn("flex flex-col gap-2", variant === "compress" && "gap-0")}
      >
        {threads.map((item: ThreadChatStoreType) => (
          <ThreadItem
            key={item.threadChatId}
            workspaceId={workspaceId}
            item={item}
          />
        ))}
        <div
          className={cn(
            "flex justify-start",
            variant === "compress" && "m-1 justify-center"
          )}
        >
          <Button variant="outline" size="sm">
            <ReloadIcon className="mr-1 h-3 w-3" />
            Load more
          </Button>
        </div>
      </div>
    </ScrollArea>
  );
}
