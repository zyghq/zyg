"use client";
import Image from "next/image";
import Link from "next/link";
import { MagnifyingGlassIcon } from "@radix-ui/react-icons";
import CloseButton from "@/components/close-btn";
import SendMessageCTA from "@/components/send-message-cta";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ScrollArea } from "@/components/ui/scroll-area";
import Threads from "@/components/threads";

import { useCustomer } from "@/lib/customer";

interface HomeFeed {
  id: string;
  type: string;
  title: string;
  previewText?: string;
  href?: string;
  ctas?: { link: string; text: string }[];
}

const homeFeeds = [
  {
    id: "1",
    type: "link",
    title: "What is Zyg?",
    previewText: "Zyg is purpose-built customer support for your SaaS products",
    href: "https://www.zyg.ai/",
  },
  {
    id: "3",
    type: "link",
    title: "What is Zyg?",
    previewText: "Zyg is purpose-built customer support for your SaaS products",
    href: "https://www.zyg.ai/",
  },
  {
    id: "4",
    type: "link",
    title: "What is Zyg?",
    previewText: "Zyg is purpose-built customer support for your SaaS products",
    href: "https://www.zyg.ai/",
  },
  {
    id: "2",
    type: "card",
    title: "Let's set up time for a demo",
    ctas: [
      {
        link: "https://www.zyg.ai/",
        text: "Schedule a demo",
      },
    ],
  },
];

const threads = [
  {
    threadChatId: "th_cpv9sektiduej80e6lug",
    sequence: 15526674508,
    status: "todo",
    read: false,
    replied: false,
    priority: "normal",
    customer: {
      customerId: "c_cpv9cn4tidu9gavgfk5g",
      name: "Tom Allison",
    },
    assignee: null,
    createdAt: "2024-06-28T11:27:54Z",
    updatedAt: "2024-06-28T11:27:54Z",
    messages: [
      {
        threadChatId: "th_cpv9sektiduej80e6lug",
        threadChatMessageId: "thm_cq04lrctidufe383koeg",
        body: "are you checking this?",
        sequence: 15636421359,
        customer: {
          customerId: "c_cpv9cn4tidu9gavgfk5g",
          name: "Tom Allison",
        },
        createdAt: "2024-06-29T17:57:01Z",
        updatedAt: "2024-06-29T17:57:01Z",
      },
    ],
  },
  {
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
    ],
  },
];

export default function Home() {
  const { isLoading, hasError, customer } = useCustomer();

  const renderHomeFeed = (feed: HomeFeed) => {
    if (feed.type === "link" && feed.href) {
      return (
        <Link
          key={feed.id}
          href={feed.href}
          target="_blank"
          className="flex font-normal"
        >
          <Card className="p-0 w-full">
            <CardHeader className="p-4">
              <CardTitle className="font-normal">{feed.title}</CardTitle>
              {feed.previewText && (
                <CardDescription>{feed.previewText}</CardDescription>
              )}
            </CardHeader>
          </Card>
        </Link>
      );
    } else if (feed.type === "card") {
      return (
        <Card key={feed.id} className="p-0 w-full">
          <CardHeader className="p-4">
            <CardTitle className="font-normal text-muted-foreground">
              {feed.title}
            </CardTitle>
            <CardDescription>{feed.previewText}</CardDescription>
          </CardHeader>
          <CardFooter className="p-4 flex space-y-1">
            {feed.ctas?.map((cta) => (
              <Button key={cta.link} variant="secondary" asChild>
                <Link href={cta.link}>{cta.text}</Link>
              </Button>
            ))}
          </CardFooter>
        </Card>
      );
    } else {
      return null;
    }
  };

  if (hasError) {
    return (
      <div className="absolute bg-white bg-opacity-60 z-10 h-full w-full flex items-center justify-center">
        <div className="flex flex-col items-center justify-center">
          <span className="text-xl">{`We're sorry, something went wrong.`}</span>
          <span className="text-xl">Please try again later.</span>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="absolute bg-white bg-opacity-60 z-10 h-full w-full flex items-center justify-center">
        <div className="flex flex-col items-center justify-center">
          <svg
            className="animate-spin h-5 w-5 text-gray-600"
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

  return (
    <div className="flex min-h-screen flex-col font-sans">
      <div className="z-10 w-full justify-between">
        <div className="flex w-full justify-between items-center p-4 bg-white">
          <div className="text-xl">Hey! How can we help?</div>
          <CloseButton />
        </div>
        <div className="fixed bottom-0 left-0 flex w-full border-t flex-col bg-white">
          <div className="w-full px-4 py-4">
            <SendMessageCTA />
          </div>
          <div className="w-full border-t flex justify-center items-center py-2">
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
        <ScrollArea className="h-[calc(100dvh-12rem)]">
          <div className="flex w-full flex-col mt-1 px-4 gap-2">
            <Button
              variant="secondary"
              className="w-full text-muted-foreground flex font-normal items-center"
              asChild
            >
              <Link href="/search">
                <MagnifyingGlassIcon className="h-4 w-4 mr-1" />
                Search for articles
              </Link>
            </Button>

            <Tabs defaultValue="home" className="">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="home">Home</TabsTrigger>
                <TabsTrigger value="threads">Threads</TabsTrigger>
              </TabsList>
              <TabsContent value="home" className="w-full flex flex-col gap-2">
                {homeFeeds.map((feed) => renderHomeFeed(feed))}
              </TabsContent>
              <TabsContent value="threads">
                <Threads threads={threads} />
              </TabsContent>
            </Tabs>
          </div>
        </ScrollArea>
      </main>
    </div>
  );
}
