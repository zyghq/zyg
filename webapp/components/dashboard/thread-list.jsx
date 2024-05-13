"use client";

import { cn } from "@/lib/utils";
import { createClient } from "@/utils/supabase/client";
import { getSession } from "@/utils/supabase/helpers";
import { useQuery } from "@tanstack/react-query";
import Avatar from "boring-avatars";
import formatDistanceToNow from "date-fns/formatDistanceToNow";

import Link from "next/link";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";

import { OopsDefault } from "@/components/errors";

import { ChatBubbleIcon, ResetIcon } from "@radix-ui/react-icons";
import { ReloadIcon } from "@radix-ui/react-icons";

export default function ThreadList({
  workspaceId,
  threads,
  className,
  variant = "default",
  endpoint = "/threads/chat/",
}) {
  const supabase = createClient();
  const result = useQuery({
    queryKey: ["threads", workspaceId, endpoint, supabase],
    queryFn: async () => {
      const { token, error: sessErr } = await getSession(supabase);
      if (sessErr) throw new Error("session expired or not found");
      const url = `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}${endpoint}`;
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const { status, statusText } = response;
        return {
          error: new Error(`error fetching threads: ${status} ${statusText}`),
        };
      }

      const data = await response.json();
      return data;
    },
    refetchOnWindowFocus: true,
    initialData: threads,
    refetchInterval: 1000 * 60 * 1,
    refetchOnMount: "always",
  });

  const { data } = result;

  if (result.isError) {
    return <OopsDefault />;
  }

  // TODO: handle result.isPending

  return (
    <ScrollArea className={cn("pr-1", className)}>
      <div
        className={cn("flex flex-col gap-2", variant === "compress" && "gap-0")}
      >
        {data.map((item) => (
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

function ThreadItem({ workspaceId, item, variant = "default" }) {
  const supabase = createClient();
  const result = useQuery({
    queryKey: ["thlabels", workspaceId, supabase, item.threadChatId],
    queryFn: async () => {
      const { token, error: sessErr } = await getSession(supabase);
      if (sessErr) throw new Error("session expired or not found");
      const url = `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}/threads/chat/${item.threadChatId}/labels/`;
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const { status, statusText } = response;
        return {
          error: new Error(`error fetching threads: ${status} ${statusText}`),
        };
      }

      const data = await response.json();
      return data;
    },
    refetchOnWindowFocus: true,
    refetchInterval: 1000 * 60 * 3,
    refetchOnMount: true,
  });

  const message = item.messages[0];
  const name = item?.customer?.name || "Customer";
  const { assignee } = item;

  const renderLabels = () => {
    if (result.isSuccess && result.data && result.data.length) {
      return (
        <div className="flex gap-1">
          {result.data.map((label) => (
            <Badge
              key={label.labelId}
              variant="outline"
              className="font-normal"
            >
              {label.name}
            </Badge>
          ))}
        </div>
      );
    }
    return <div className="min-h-5"></div>;
  };

  return (
    <Link
      key={item.threadChatId}
      href={`/${workspaceId}/threads/${item.threadChatId}/`}
      className={cn(
        "flex flex-col items-start gap-2 rounded-lg border px-3 py-3 text-left text-sm transition-all hover:bg-accent",
        variant === "compress" && "gap-0 rounded-none py-5"
      )}
    >
      <div className="flex w-full flex-col gap-1">
        <div className="flex items-center">
          <div className="flex items-center gap-2">
            <ChatBubbleIcon />
            <div className="font-semibold">{name}</div>
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
            {formatDistanceToNow(new Date(message.updatedAt), {
              addSuffix: true,
            })}
          </div>
          {assignee && (
            <Avatar
              size={28}
              name={assignee.name}
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
            <Badge variant="outline" className="bg-indigo-100 font-normal">
              <div className="flex items-center gap-1">
                <ResetIcon className="h-3 w-3" />
                awaiting reply
              </div>
            </Badge>
          </div>
        )}
        {item?.title ? <div className="font-medium">{item?.title}</div> : null}
      </div>
      <div className="line-clamp-2 text-muted-foreground">{message.body}</div>
      {renderLabels()}
    </Link>
  );
}
