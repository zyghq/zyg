import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import Avatar from "boring-avatars";
import { cn } from "@/lib/utils";

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
    name: string | null;
  } | null;
  member?: {
    memberId: string;
    name: string;
  } | null;
  createdAt: string;
  updatedAt: string;
}

function ThreadChatItem({
  customerId,
  thread,
}: {
  customerId: string;
  thread: Thread;
}) {
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
      key={thread.threadChatId}
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

export default function Threads({ threads }: { threads: Thread[] }) {
  return (
    <div className="flex flex-col gap-2">
      {threads.map((thread) => (
        <ThreadChatItem
          key={thread.threadChatId}
          customerId={thread.customer.customerId}
          thread={thread}
        />
      ))}
    </div>
  );
}
