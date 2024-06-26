"use client";
import * as React from "react";
import { cn } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import { useAuth } from "@/lib/auth";
import formatDistanceToNow from "date-fns/formatDistanceToNow";
import Avatar from "boring-avatars";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Skeleton } from "@/components/ui/skeleton";
import { Icons } from "@/components/icons";
import Link from "next/link";

function ThreadChatItem({ customerId, thread }) {
  const { threadChatId, messages } = thread;
  const message = messages[0];
  const isMe = message?.customer?.customerId === customerId;
  const body = message?.body || "...";
  const updatedAt = message?.updatedAt || new Date();
  const isRead = thread?.read || false;
  const isReplied = thread?.replied || false;

  const memberName = message?.member?.name || "Member";
  const memberId = message?.member?.memberId || "";

  return (
    <Link
      href={`/threads/${threadChatId}/`}
      key={thread.threadId}
      className={cn(
        "flex flex-col items-start gap-2 rounded-lg border px-3 py-3 text-left text-sm transition-all hover:bg-accent"
      )}
    >
      <div className="flex w-full flex-col gap-1">
        <div className="flex items-center">
          <div className="flex items-center gap-2">
            {isMe ? (
              <Avatar name={customerId} size={24} variant="marble" />
            ) : (
              <Avatar name={memberId} size={24} variant="marble" />
            )}
            <div className="font-medium">{`${isMe ? "You" : memberName}`}</div>
            {!isRead && (
              <span className="flex h-2 w-2 rounded-full bg-blue-600" />
            )}
          </div>
          <div
            className={cn(
              "ml-auto mr-2 text-xs",
              !isReplied ? "text-foreground" : "text-muted-foreground"
            )}
          >
            {formatDistanceToNow(new Date(updatedAt), {
              addSuffix: true,
            })}
          </div>
        </div>
      </div>
      <div className="line-clamp-2 text-muted-foreground text-xs">{body}</div>
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
      const response = await fetch(`http://localhost:8000/threads/chat/`, {
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
    refetchOnWindowFocus: true,
    enabled: authUser?.authToken?.value ? true : false,
    initialData: threads,
    refetchInterval: 1000,
    refetchOnMount: "always",
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
    <ScrollArea className="h-[calc(100dvh-18rem)] pr-1">
      <div className="space-y-2">
        {data.map((thread) => (
          <ThreadChatItem
            key={thread.threadChatId}
            customerId={authUser?.customerId}
            thread={thread}
          />
        ))}
      </div>
    </ScrollArea>
  );
}
