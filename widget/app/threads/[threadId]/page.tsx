"use client";
import * as React from "react";
import Avatar from "boring-avatars";
import HomeButton from "@/components/home-btn";
import CloseButton from "@/components/close-btn";
import { ScrollArea } from "@/components/ui/scroll-area";
import MessageThreadForm from "@/components/message-thread-form";

// const threadMessages = {
//   threadChatId: "th_cpv9sektiduej80e6lug",
//   sequence: 15526674508,
//   status: "todo",
//   read: false,
//   replied: false,
//   priority: "normal",
//   customer: {
//     customerId: "c_cpv9cn4tidu9gavgfk5g",
//     name: "Tom Allison",
//   },
//   assignee: null,
//   createdAt: "2024-06-28T11:27:54Z",
//   updatedAt: "2024-06-28T11:27:54Z",
//   messages: [
//     {
//       threadChatId: "th_cpv9sektiduej80e6lug",
//       threadChatMessageId: "thm_cq04lrctidufe383koeg",
//       body: "are you checking this?",
//       sequence: 15636421359,
//       customer: {
//         customerId: "c_cpv9cn4tidu9gavgfk5g",
//         name: "Tom Allison",
//       },
//       createdAt: "2024-06-29T17:57:01Z",
//       updatedAt: "2024-06-29T17:57:01Z",
//     },
//     {
//       threadChatId: "th_cpv9sektiduej80e6lug",
//       threadChatMessageId: "thm_cpv9sektiduej80e6lv0",
//       body: "I am having trouble when creating token, the personal access token. am I missing something here?\n\nlove if I can get some help\n\nThanks",
//       sequence: 15526674527,
//       customer: {
//         customerId: "c_cpv9cn4tidu9gavgfk5g",
//         name: "Tom Allison",
//       },
//       createdAt: "2024-06-28T11:27:54Z",
//       updatedAt: "2024-06-28T11:27:54Z",
//     },
//   ],
// };

const threadMessages = {
  threadChatId: "th_cpv9uj4tiduej80e6lvg",
  sequence: 15526948204,
  status: "todo",
  read: false,
  replied: true,
  priority: "normal",
  customer: {
    customerId: "c_cpv9cn4tidu9gavgfk5g",
    name: "Tom Allison",
  },
  assignee: {
    memberId: "m_co60epktidu7sod96la0",
    name: "Manmohini",
  },
  createdAt: "2024-06-28T11:32:28Z",
  updatedAt: "2024-06-29T17:54:28Z",
  messages: [
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfec4tidu6vcj5an10",
      body: "I am not sure that many messages are required!! But thanks!",
      sequence: 15549449004,
      customer: {
        customerId: "c_cpv9cn4tidu9gavgfk5g",
        name: "Tom Allison",
      },
      createdAt: "2024-06-28T17:47:28Z",
      updatedAt: "2024-06-28T17:47:28Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfe1ctidu6hi227dr0",
      body: "done!",
      sequence: 15549405549,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T17:46:45Z",
      updatedAt: "2024-06-28T17:46:45Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfdo4tidu6hi227dqg",
      body: "done!",
      sequence: 15549368157,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T17:46:08Z",
      updatedAt: "2024-06-28T17:46:08Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfclctidu6hi227dp0",
      body: "wow",
      sequence: 15549229943,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T17:43:49Z",
      updatedAt: "2024-06-28T17:43:49Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfce4tidu6hi227dog",
      body: "wow",
      sequence: 15549200308,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T17:43:20Z",
      updatedAt: "2024-06-28T17:43:20Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfc7ktidu6hi227do0",
      body: "wow",
      sequence: 15549174765,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T17:42:54Z",
      updatedAt: "2024-06-28T17:42:54Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfblktidu6hi227dn0",
      body: "wow",
      sequence: 15549102482,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T17:41:42Z",
      updatedAt: "2024-06-28T17:41:42Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfbl4tidu6hi227dmg",
      body: "wow",
      sequence: 15549100363,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T17:41:40Z",
      updatedAt: "2024-06-28T17:41:40Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvfat4tidu6hi227dl0",
      body: "Okay lets implement this!!",
      sequence: 15549004442,
      member: {
        memberId: "m_cpv7ts4tiduau9n48ufg",
        name: "Sanchit",
      },
      createdAt: "2024-06-28T17:40:04Z",
      updatedAt: "2024-06-28T17:40:04Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpve1tstidu6hi227dig",
      body: "I think you are right",
      sequence: 15543759340,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T16:12:39Z",
      updatedAt: "2024-06-28T16:12:39Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpve1g4tidu6hi227dhg",
      body: "I think you are right",
      sequence: 15543704895,
      member: {
        memberId: "m_co60epktidu7sod96la0",
        name: "Manmohini",
      },
      createdAt: "2024-06-28T16:11:44Z",
      updatedAt: "2024-06-28T16:11:44Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvce9stidu0nnp475c0",
      body: "Yeah, thanks we can think about it!",
      sequence: 15537151422,
      member: {
        memberId: "m_cpv7ts4tiduau9n48ufg",
        name: "Sanchit",
      },
      createdAt: "2024-06-28T14:22:31Z",
      updatedAt: "2024-06-28T14:22:31Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpvab6ktiduasjl7kqig",
      body: "I think it will be too much work",
      sequence: 15528562194,
      customer: {
        customerId: "c_cpv9cn4tidu9gavgfk5g",
        name: "Tom Allison",
      },
      createdAt: "2024-06-28T11:59:22Z",
      updatedAt: "2024-06-28T11:59:22Z",
    },
    {
      threadChatId: "th_cpv9uj4tiduej80e6lvg",
      threadChatMessageId: "thm_cpv9uj4tiduej80e6m00",
      body: "ah! now I am able to send messages, there seems to be issue with threadId - moving to typescript should help I guess for these type of issues.",
      sequence: 15526948205,
      customer: {
        customerId: "c_cpv9cn4tidu9gavgfk5g",
        name: "Tom Allison",
      },
      createdAt: "2024-06-28T11:32:28Z",
      updatedAt: "2024-06-28T11:32:28Z",
    },
  ],
};

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
  console.log(threadId);
  const bottomRef = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, []);

  const { messages } = threadMessages;
  const messagesReversed = messages.slice().reverse();

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
            <MessageThreadForm />
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
            <div ref={bottomRef}></div>
          </div>
        </ScrollArea>
      </main>
    </div>
  );
}
