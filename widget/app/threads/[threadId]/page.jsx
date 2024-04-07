"use client";

import * as React from "react";
import { ThreadHeader } from "@/components/headers";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AvatarImage, AvatarFallback, Avatar } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
import Room from "@/components/room";
import MessageThreadForm from "@/components/message-thread-form";
import { useQuery } from "@tanstack/react-query";

function Message({ message, isMe = true }) {
  return (
    <div className={`max-w-xs ${isMe ? "ml-auto" : "mr-auto"}`}>
      <div className="flex space-x-2">
        <Avatar className="h-6 w-6">
          <AvatarImage alt="User" src="/images/profile.jpg" />
          <AvatarFallback>U</AvatarFallback>
        </Avatar>
        <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
          <div className="text-xs text-muted-foreground">{`${isMe ? "Me" : "C"}`}</div>
          <p className="text-sm">{message.body}</p>
        </div>
      </div>
    </div>
  );
}

export default function ThreadPage({ params }) {
  const { threadId } = params;
  const result = useQuery({
    queryKey: ["thchats", threadId],
    queryFn: async () => {
      const response = await fetch(
        `http://localhost:8080/-/threads/chat/${threadId}/messages/`,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization:
              "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ3b3Jrc3BhY2VJZCI6Indya2NvNjBlcGt0aWR1N3NvZDk2bDkwIiwiZXh0ZXJuYWxJZCI6Inh4eHgtMTExLXp6enoiLCJlbWFpbCI6InNhbmNoaXRycmtAZ21haWwuY29tIiwicGhvbmUiOiIrOTE3NzYwNjg2MDY4IiwiaXNzIjoiYXV0aC56eWcuYWkiLCJzdWIiOiJjX2NvNjFhYmt0aWR1MXQzaTNkbjYwIiwiYXVkIjpbImN1c3RvbWVyIl0sImV4cCI6MTc0Mzc1Nzg3MSwibmJmIjoxNzEyMjIxODcxLCJpYXQiOjE3MTIyMjE4NzEsImp0aSI6Indya2NvNjBlcGt0aWR1N3NvZDk2bDkwOmNfY282MWFia3RpZHUxdDNpM2RuNjAifQ.epCQ4aXvYPXIhVrX6TtfYrq0XxYXT18kIWsOae8HvUQ",
          },
        }
      );
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    },
  });

  const renderContent = () => {
    if (result.isPending) {
      return (
        <div className="space-y-4">
          <div className="max-w-xs ml-auto">
            <div className="flex items-center space-x-1">
              <Skeleton className="h-7 w-7 rounded-full" />
              <div className="space-y-2">
                <Skeleton className="h-4 w-[250px]" />
                <Skeleton className="h-4 w-[200px]" />
              </div>
            </div>
          </div>
          <div className="max-w-xs ml-auto">
            <div className="flex items-center space-x-1">
              <Skeleton className="h-7 w-7 rounded-full" />
              <div className="space-y-2">
                <Skeleton className="h-4 w-[250px]" />
                <Skeleton className="h-4 w-[200px]" />
              </div>
            </div>
          </div>
          <div className="flex items-center space-x-1">
            <Skeleton className="h-7 w-7 rounded-full" />
            <div className="space-y-2">
              <Skeleton className="h-4 w-[250px]" />
              <Skeleton className="h-4 w-[200px]" />
            </div>
          </div>

          <div className="max-w-xs ml-auto">
            <div className="flex items-center space-x-1">
              <Skeleton className="h-7 w-7 rounded-full" />
              <div className="space-y-2">
                <Skeleton className="h-4 w-[250px]" />
                <Skeleton className="h-4 w-[200px]" />
              </div>
            </div>
          </div>
        </div>
      );
    }

    if (result.isError) {
      return (
        <div className="flex justify-center">
          <div className="text-red-500">
            Something went wrong. Please try again later.
          </div>
        </div>
      );
    }

    const { data } = result;
    const { messages } = data;
    const reversedMessages = messages.reverse();
    return (
      <div className="space-y-2">
        {reversedMessages.map((message) => (
          <Message key={message.threadChatMessageId} message={message} />
        ))}
      </div>
    );
  };

  return (
    <React.Fragment>
      <ThreadHeader />
      <ScrollArea className="p-4 h-[calc(100dvh-12rem)]">
        {renderContent()}
      </ScrollArea>
      <div className="pt-2 px-2 mt-auto border-t">
        <MessageThreadForm threadId={threadId} />
        <footer className="flex flex-col justify-center items-center border-t w-full h-8 mt-2">
          <a
            href="https://www.zyg.ai/"
            className="text-xs font-semibold text-gray-500"
            target="_blank"
          >
            Powered by Zyg.
          </a>
        </footer>
      </div>
    </React.Fragment>
  );
}
