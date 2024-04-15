"use client";
import * as React from "react";
import { cn } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import { useAuth } from "@/lib/auth";
import formatDistanceToNow from "date-fns/formatDistanceToNow";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Skeleton } from "@/components/ui/skeleton";
import { Icons } from "@/components/icons";
import Link from "next/link";

// watch this video for scroll bottom
// https://www.youtube.com/watch?v=yaIytT_Y0DA
function ThreadChatItem({ customerId, thread }) {
  const { threadId, messages } = thread;
  const message = messages[0];

  const isMe = message?.customer?.customerId === customerId;
  const body = message?.body || "...";
  const updatedAt = message?.updatedAt || new Date();

  const memberName = (member) => {
    return member?.name || "Member";
  };

  return (
    <Link
      href={`/threads/${threadId}/`}
      key={thread.threadId}
      className={cn(
        "w-full",
        "flex flex-col items-start gap-2 rounded-lg border px-3 py-3 text-left text-sm transition-all hover:bg-accent"
      )}
    >
      <div className="flex w-full flex-col gap-2">
        <div className="flex items-center">
          <div className="flex items-center gap-1">
            <Avatar className="h-8 w-8 rounded-sm">
              <AvatarImage src="https://github.com/shadcn.png" />
              <AvatarFallback>CN</AvatarFallback>
            </Avatar>
            <div className="font-semibold">{`Zyg Team`}</div>
            <span className="flex h-2 w-2 rounded-full bg-blue-600" />
          </div>
          <div className={cn("ml-auto text-xs", "text-foreground")}>
            {formatDistanceToNow(new Date(updatedAt), {
              addSuffix: true,
            })}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Avatar className="h-5 w-5">
            <AvatarImage src="https://github.com/shadcn.png" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
          <div className="flex flex-col">
            <div className="text-xs font-medium">{`${isMe ? "You" : memberName(message?.member)}`}</div>
            <div className="text-xs">
              {body.substring(0, 200) + (body.length > 220 ? " ..." : "")}
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}

export default function ThreadList({ threads }) {
  const auth = useAuth();
  const { authUser, isAuthLoading } = auth;
  const result = useQuery({
    queryKey: ["thchats", authUser?.authToken?.value],
    queryFn: async () => {
      const token = authUser?.authToken?.value || "";
      const response = await fetch(`http://localhost:8080/-/threads/chat/`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });
      if (!response.ok) {
        throw new Error("Failed to fetch threads");
      }
      return response.json();
    },
    refetchOnWindowFocus: false,
    enabled: authUser?.authToken?.value ? true : false,
    initialData: threads,
    refetchInterval: 10000,
  });

  if (result.isPending || isAuthLoading || !authUser) {
    return (
      <div className="flex flex-col gap-2">
        <div className="flex items-center space-x-2">
          <Skeleton className="h-8 w-8 rounded-sm" />
          <Skeleton className="h-4 w-[250px]" />
        </div>
        <div className="flex items-center gap-2">
          <Skeleton className="h-5 w-5 rounded-full" />
          <Skeleton className="h-4 w-[250px]" />
        </div>
      </div>
    );
  }

  if (result.isError) {
    return (
      <div className="flex flex-col justify-center items-center mt-4 px-2 space-y-1">
        <Icons.oops className="h-10 w-10" />
        <div className="text-xs">something went wrong.</div>
      </div>
    );
  }

  const { data } = result;

  return (
    <ScrollArea className="h-[calc(100dvh-18rem)]">
      <div className="space-y-2">
        {data.map((thread) => (
          <ThreadChatItem
            key={thread.threadId}
            customerId={authUser?.customerId}
            thread={thread}
          />
        ))}
      </div>
    </ScrollArea>
  );
}
