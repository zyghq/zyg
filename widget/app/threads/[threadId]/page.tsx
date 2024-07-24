"use client";
import * as React from "react";
import Avatar from "boring-avatars";
import HomeButton from "@/components/home-btn";
import CloseButton from "@/components/close-btn";
import { ScrollArea } from "@/components/ui/scroll-area";
import MessageThreadForm from "@/components/message-thread-form";
import { Icons } from "@/components/icons";
import { useCustomer } from "@/lib/customer";
import { useQuery } from "@tanstack/react-query";
import AskEmailForm from "@/components/ask-email-form";

interface Thread {
  threadChatId: string;
  sequence: number;
  status: string;
  read: boolean;
  replied: boolean;
  priority: string;
  customer: {
    customerId: string;
    name: string;
  };
  assignee: {
    memberId: string;
    name: string;
  } | null;
  createdAt: string;
  updatedAt: string;
  messages: Message[];
}

interface Message {
  threadChatId: string;
  threadChatMessageId: string;
  body: string;
  sequence: number;
  customer?: {
    customerId: string;
    name: string;
  } | null;
  member?: {
    memberId: string;
    name: string;
  } | null;
  createdAt: string;
  updatedAt: string;
}

function Message({ message }: { message: Message }) {
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

  const memberId = message?.member?.memberId || "";
  const memberName = message?.member?.name || "";

  return (
    <div className="flex">
      <div className={`flex max-w-sm ${isMe ? "ml-auto" : "mr-auto"}`}>
        <div className="flex space-x-2">
          {isMe ? null : <Avatar size={22} name={memberId} variant="marble" />}
          <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800">
            <div className="text-xs text-muted-foreground">{`${
              isMe ? "Me" : memberName
            }`}</div>
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

export default function ThreadMessages({
  params,
}: {
  params: { threadId: string };
}) {
  const { threadId } = params;
  const { isLoading, hasError, customer } = useCustomer();

  const hasIdentity =
    customer?.customerEmail ||
    customer?.customerPhone ||
    customer?.customerExternalId;

  const {
    data: thread,
    isLoading: isLoadingThread,
    error: errorThread,
    refetch,
  } = useQuery({
    queryKey: ["messages", threadId],
    queryFn: async () => {
      const jwt = customer?.jwt;
      if (!jwt) {
        console.error("No JWT found");
        return null;
      }
      const { widgetId } = customer;
      const response = await fetch(
        `/api/widgets/${widgetId}/threads/${threadId}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            jwt,
          }),
        }
      );

      if (!response.ok) {
        throw new Error("Not Found");
      }
      const thread = await response.json();
      return thread as Thread;
    },
    enabled: !!customer,
  });

  const bottomRef = React.useRef<HTMLDivElement>(null);
  React.useEffect(() => {
    if (bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, []);

  if (hasError || errorThread) {
    return (
      <div className="absolute z-10 h-full w-full flex items-center justify-center">
        <div className="flex flex-col items-center justify-center text-muted-foreground">
          <span className="text-lg">{`We're sorry, something went wrong.`}</span>
          <span className="text-lg">Please try again later.</span>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="absolute z-10 h-full w-full flex items-center justify-center">
        <div className="flex flex-col items-center justify-center">
          <svg
            className="animate-spin h-5 w-5 text-muted-foreground"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            ></circle>
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
        </div>
      </div>
    );
  }

  if (isLoadingThread) return null;

  if (!thread || !thread?.messages?.length) {
    return (
      <div className="flex flex-col items-center justify-center space-y-4 mt-24">
        <Icons.nothing className="w-40" />
        <p className="text-center text-muted-foreground">
          Nothing to see here yet.
        </p>
      </div>
    );
  }

  const { messages } = thread;

  const hasSentMessageWithoutIdentity =
    messages.length > 0 && !hasIdentity && customer;

  const messagesReversed = hasSentMessageWithoutIdentity
    ? messages.reverse().slice(messages.length - 1)
    : messages.reverse();

  return (
    <div className="flex min-h-screen flex-col font-sans">
      <div className="z-10 w-full justify-between">
        <div className="flex items-center justify-start py-4 border-b px-4 gap-1">
          <HomeButton />
          <div>
            <div className="flex flex-col">
              <div className="font-semibold">Zyg Team</div>
              <div className="text-xs text-muted-foreground">
                Ask us anything, or share your feedback.
              </div>
            </div>
          </div>
          <div className="ml-auto">
            <CloseButton />
          </div>
        </div>
        <div className="fixed bottom-0 left-0 flex w-full border-t flex-col bg-white">
          <div className="flex flex-col px-4 pt-4">
            {customer && (
              <MessageThreadForm
                disabled={!!hasSentMessageWithoutIdentity}
                widgetId={customer.widgetId}
                threadId={threadId}
                jwt={customer.jwt}
                refetch={refetch}
              />
            )}
          </div>
          <div className="w-full flex justify-center items-center py-2">
            <a
              href="https://www.zyg.ai/"
              className="text-xs font-semibold text-muted-foreground"
              target="_blank"
            >
              Powered by Zyg
            </a>
          </div>
        </div>
      </div>
      <main>
        <ScrollArea className="p-4 h-[calc(100dvh-12rem)]">
          <div className="space-y-2">
            {messagesReversed.map((message) => (
              <Message key={message.threadChatMessageId} message={message} />
            ))}
            {hasSentMessageWithoutIdentity && (
              <div className="flex flex-col px-2">
                <div className="text-sm max-w-xs font-semibold mb-1">
                  Please provide your email address so we can reach you.
                </div>
                <AskEmailForm
                  widgetId={customer.widgetId}
                  threadId={threadId}
                  jwt={customer.jwt}
                />
              </div>
            )}
            <div ref={bottomRef}></div>
          </div>
        </ScrollArea>
      </main>
    </div>
  );
}
