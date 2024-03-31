import formatDistanceToNow from "date-fns/formatDistanceToNow";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { ReloadIcon } from "@radix-ui/react-icons";

export default function Threads({ items, className, variant = "default" }) {
  const mail = {
    selected: false,
  };

  return (
    <ScrollArea className={cn("pr-1", className)}>
      <div
        className={cn("flex flex-col gap-2", variant === "compress" && "gap-0")}
      >
        {items.map((item) => (
          <button
            key={item.id}
            className={cn(
              "flex flex-col items-start gap-2 rounded-lg border px-3 py-3 text-left text-sm transition-all hover:bg-accent",
              mail.selected === item.id && "bg-muted",
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
                  <Avatar
                    className={cn(
                      "h-8 w-8",
                      variant === "compress" && "h-6 w-6",
                    )}
                  >
                    <AvatarImage src="https://github.com/shadcn.png" />
                    <AvatarFallback>CN</AvatarFallback>
                  </Avatar>
                  <div className="font-semibold">{item.name}</div>
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
                  {formatDistanceToNow(new Date(item.date), {
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
              <div className="font-medium">{item.subject}</div>
            </div>
            <div className="line-clamp-2 text-muted-foreground">
              {item.text.substring(0, 300)}
            </div>
            {item.labels.length && variant === "default" ? (
              <div className="flex items-center gap-2">
                {item.labels.map((label) => (
                  <Badge key={label} variant={getBadgeVariantFromLabel(label)}>
                    {label}
                  </Badge>
                ))}
              </div>
            ) : null}
          </button>
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

function getBadgeVariantFromLabel(label) {
  if (["work"].includes(label.toLowerCase())) {
    return "default";
  }

  if (["personal"].includes(label.toLowerCase())) {
    return "outline";
  }

  return "secondary";
}
