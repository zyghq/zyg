"use client";
// TODO: add authentication for this page...
import * as React from "react";
import { ThreadHeader } from "@/components/headers";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AvatarImage, AvatarFallback, Avatar } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
// import Room from "@/components/room";
import MessageThreadForm from "@/components/message-thread-form";
import { useQuery } from "@tanstack/react-query";
import Cookies from "js-cookie";
import { Icons } from "@/components/icons";

function Message({ message }) {
  const { createdAt } = message;
  const date = new Date(createdAt);
  const time = date.toLocaleString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

  const isMe = typeof message.customer === "object";
  const memName = () => {
    return message?.member?.name || "C";
  };
  return (
    <div className="flex">
      <div className={`flex max-w-sm ${isMe ? "ml-auto" : "mr-auto"}`}>
        <div className="flex space-x-2">
          {isMe ? null : (
            <Avatar className="h-6 w-6">
              <AvatarImage alt="User" src="/images/profile.jpg" />
              <AvatarFallback>U</AvatarFallback>
            </Avatar>
          )}
          <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
            <div className="text-xs text-muted-foreground">{`${isMe ? "Me" : memName()}`}</div>
            <p className="text-sm">{message.body}</p>
            <div className="flex text-xs justify-end text-muted-foreground mt-1">
              {time}
            </div>
          </div>
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
      const token = Cookies.get("__zygtoken") || "";
      const response = await fetch(
        `http://localhost:8080/-/threads/chat/${threadId}/messages/`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (!response.ok) {
        console.log("response", response);
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
        <div className="flex flex-col items-center mt-24">
          <Icons.oops className="h-8 w-8" />
          <div className="text-sm text-red-500">
            Something went wrong. Please try again later.
          </div>
        </div>
      );
    }

    const { data } = result;
    const { messages } = data;
    const messagesReversed = Array.from(messages).reverse();
    return (
      <div className="space-y-2">
        {messagesReversed.map((message) => (
          <Message key={message.threadChatMessageId} message={message} />
        ))}
      </div>
    );
  };

  return (
    <React.Fragment>
      <ThreadHeader />
      <ScrollArea type="always" className="p-4 h-[calc(100dvh-12rem)]">
        {renderContent()}
      </ScrollArea>
      <div className="pt-2 px-2 mt-auto border-t">
        <MessageThreadForm threadId={threadId} refetch={result.refetch} />
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
