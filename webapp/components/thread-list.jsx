import formatDistanceToNow from "date-fns/formatDistanceToNow";
import { cn } from "@/lib/utils";
import Link from "next/link";
// import { Badge } from "@/components/ui/badge";
import { ChatBubbleIcon } from "@radix-ui/react-icons";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { ReloadIcon } from "@radix-ui/react-icons";

function ThreadItem({ workspaceId, item, variant = "default" }) {
  const message = item.messages[0];
  const name = item?.customer?.name || "Customer";
  const body = () => {
    return (
      message.body.substring(0, 300) + (message.body.length > 300 ? "..." : "")
    );
  };

  // TODO: remove
  const mail = {
    selected: false,
  };

  return (
    <Link
      key={item.threadId}
      href={`/${workspaceId}/threads/${item.threadId}/`}
      className={cn(
        "flex flex-col items-start gap-2 rounded-lg border px-3 py-3 text-left text-sm transition-all hover:bg-accent",
        mail.selected === item.threadId && "bg-muted",
        variant === "compress" && "gap-0 rounded-none py-5",
      )}
      // onClick={() =>
      //   setMail({
      //     ...mail,
      //     selected: item.id,
      //   })
      // }
    >
      <div className="flex w-full flex-col gap-1">
        <div className="flex items-center">
          <div className="flex items-center gap-2">
            <ChatBubbleIcon />
            {/* <Avatar
                    className={cn(
                      "h-8 w-8",
                      variant === "compress" && "h-6 w-6",
                    )}
                  >
                    <AvatarImage src="https://github.com/shadcn.png" />
                    <AvatarFallback>CN</AvatarFallback>
                  </Avatar> */}
            <div className="font-semibold">{name}</div>
            {!item.read && (
              <span className="flex h-2 w-2 rounded-full bg-blue-600" />
            )}
          </div>
          <div
            className={cn(
              "ml-auto mr-2 text-xs",
              mail.selected === item.id
                ? "text-foreground"
                : "text-muted-foreground",
            )}
          >
            {formatDistanceToNow(new Date(message.updatedAt), {
              addSuffix: true,
            })}
          </div>
          <Avatar
            className={cn("h-6 w-6", variant === "compress" && "h-5 w-5")}
          >
            <AvatarImage src="https://github.com/shadcn.png" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
        </div>
        {item?.title ? <div className="font-medium">{item?.title}</div> : null}
      </div>
      <div className="line-clamp-2 text-muted-foreground">
        {body(item.messages[0])}
      </div>
      {/* {item.labels.length && variant === "default" ? (
              <div className="flex items-center gap-2">
                {item.labels.map((label) => (
                  <Badge key={label} variant={getBadgeVariantFromLabel(label)}>
                    {label}
                  </Badge>
                ))}
              </div>
            ) : null} */}
    </Link>
  );
}

export default function ThreadList({
  workspaceId,
  items,
  className,
  variant = "default",
}) {
  console.log(items[0]);

  return (
    <ScrollArea className={cn("pr-1", className)}>
      <div
        className={cn("flex flex-col gap-2", variant === "compress" && "gap-0")}
      >
        {items.map((item) => (
          <ThreadItem
            key={item.threadId}
            workspaceId={workspaceId}
            item={item}
          />
        ))}
        <div
          className={cn(
            "flex justify-start",
            variant === "compress" && "m-1 justify-center",
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

// function getBadgeVariantFromLabel(label) {
//   if (["work"].includes(label.toLowerCase())) {
//     return "default";
//   }

//   if (["personal"].includes(label.toLowerCase())) {
//     return "outline";
//   }

//   return "secondary";
// }
