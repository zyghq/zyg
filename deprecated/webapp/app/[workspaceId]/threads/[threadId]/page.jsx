import { titleCase } from "@/lib/utils";
import { getSession, isAuthenticated } from "@/utils/supabase/helpers";
import { createClient } from "@/utils/supabase/server";
import Avatar from "boring-avatars";

import Link from "next/link";

import { Button } from "@/components/ui/button";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";

import { GoBack } from "@/components/commons/buttons";
import { OopsDefault } from "@/components/errors";
import { SidePanelThreadList } from "@/components/thread/sidepanel-thread-list";
import ThreadList from "@/components/thread/thread-list";

import {
  ArrowDownIcon,
  ArrowUpIcon,
  ChatBubbleIcon,
  DotsHorizontalIcon,
  HomeIcon,
  ResetIcon,
} from "@radix-ui/react-icons";

import { CircleIcon } from "lucide-react";

/**
 * Fetches the list of thread chats for a given workspace.
 *
 * @param {string} workspaceId - The ID of the workspace.
 * @param {string} [authToken=""] - The authentication token (optional).
 * @returns {Promise<{ data: any, error: Error | null }>} - The response object containing the data and error (if any).
 */
async function getThreadChatListAPI(workspaceId, authToken = "") {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}/threads/chat/`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `Error fetching thread chats: ${status} ${statusText}`
        ),
      };
    }

    const data = await response.json();
    return { data, error: null };
  } catch (err) {
    console.error("error fetching workspace thread chats", err);
    return { data: null, error: err };
  }
}

function getPrevNextFromCurrent(threads, threadId) {
  const currentIndex = threads.findIndex(
    (thread) => thread.threadChatId === threadId
  );
  const currentItem = threads[currentIndex] || null;

  const prevItem = threads[currentIndex - 1] || null;
  const nextItem = threads[currentIndex + 1] || null;

  return { currentItem, prevItem, nextItem };
}

export default async function ThreadItemPage({ params }) {
  const { workspaceId, threadId } = params;
  const supabase = createClient();

  if (!(await isAuthenticated(supabase))) {
    return redirect("/login/");
  }

  const { token, error: sessionErr } = await getSession(supabase);
  if (sessionErr) {
    return (
      <div className="container mt-12">
        <h1 className="mb-1 text-3xl font-bold">Error</h1>
        <p className="mb-4 text-red-500">
          There was an error fetching your thread. Please try again later.
        </p>
      </div>
    );
  }

  const threads = [];

  const { error, data } = await getThreadChatListAPI(workspaceId, token);
  if (error) {
    return <OopsDefault />;
  } else {
    threads.push(...data);
  }

  const { currentItem, prevItem, nextItem } = getPrevNextFromCurrent(
    threads,
    threadId
  );

  console.log("currentItem", currentItem);

  const customerName =
    currentItem?.customer?.name ||
    "Customer " + currentItem?.customer?.customerId?.slice(-4) ||
    "N/A";

  const status = currentItem?.status ? titleCase(currentItem.status) : "N/A";

  return (
    <div className="flex flex-1">
      <div className="flex flex-col items-center px-2 lg:border-r">
        <div className="mt-4 flex flex-col gap-4">
          <GoBack />
          <Button variant="outline" size="icon" asChild>
            <Link href={`/${workspaceId}/`}>
              <HomeIcon className="h-4 w-4" />
            </Link>
          </Button>
          <SidePanelThreadList
            workspaceId={workspaceId}
            threads={threads}
            activeThreadId={threadId}
          />
          {prevItem ? (
            <Button variant="outline" size="icon" asChild>
              <Link href={`/${workspaceId}/threads/${prevItem?.threadChatId}/`}>
                <ArrowUpIcon className="h-4 w-4" />
              </Link>
            </Button>
          ) : null}
          {nextItem ? (
            <Button variant="outline" size="icon" asChild>
              <Link href={`/${workspaceId}/threads/${nextItem?.threadChatId}/`}>
                <ArrowDownIcon className="h-4 w-4" />
              </Link>
            </Button>
          ) : null}
        </div>
      </div>
      <div className="flex flex-col">
        <ResizablePanelGroup direction="horizontal">
          <ResizablePanel
            defaultSize={25}
            minSize={20}
            maxSize={30}
            className="hidden sm:block"
          >
            <div className="flex h-full flex-col">
              <div className="flex h-14 flex-col justify-center border-b px-4">
                <div className="text-`md font-semibold">Threads</div>
              </div>
              <ThreadList
                workspaceId={workspaceId}
                items={threads}
                className="h-[calc(100dvh-8rem)] p-1"
                activeThreadId={threadId}
                variant="compress"
              />
            </div>
          </ResizablePanel>
          <ResizableHandle withHandle={false} />
          <ResizablePanel defaultSize={50} className="flex flex-col">
            <ResizablePanelGroup direction="vertical">
              <ResizablePanel defaultSize={75}>
                <div className="flex h-full flex-col">
                  <div className="flex h-14 min-h-14 flex-col justify-center border-b px-4">
                    <div className="flex">
                      <div className="text-sm font-semibold">
                        {customerName}
                      </div>
                    </div>
                    <div className="flex items-center">
                      <CircleIcon className="mr-1 h-3 w-3 text-indigo-500" />
                      <span className="items-center text-xs">{status}</span>
                      <Separator orientation="vertical" className="mx-2" />
                      <ChatBubbleIcon className="h-3 w-3" />
                      {/* disabled for now, enable for something else perhaps? */}
                      {/* <Separator orientation="vertical" className="mx-2" />
                      <span className="font-mono text-xs">12/44</span> */}
                    </div>
                  </div>
                  <ScrollArea className="flex h-full flex-auto flex-col px-2 pb-4">
                    <div className="flex flex-col gap-1">
                      <div className="m-4">
                        <div className="flex items-center font-mono text-sm font-medium">
                          <span className="mr-1 flex h-1 w-1 rounded-full bg-fuchsia-500" />
                          Monday, 14 February 2024
                        </div>
                      </div>
                      {/* message */}
                      <div className="flex flex-col gap-2 rounded-lg border bg-background p-4">
                        <div className="flex w-full flex-col gap-1">
                          <div className="flex items-center">
                            <div className="flex items-center gap-2">
                              <Avatar size={28} name="name" variant="beam" />
                              {/* <Avatar className="h-8 w-8">
                                <AvatarImage src="https://github.com/shadcn.png" />
                                <AvatarFallback>CN</AvatarFallback>
                              </Avatar> */}
                              <div className="font-medium">Emily Davis</div>
                              <span className="flex h-1 w-1 rounded-full bg-blue-600" />
                              <span className="text-xs">3d ago.</span>
                            </div>
                            <div className="ml-auto">
                              <Button variant="ghost" size="icon">
                                <DotsHorizontalIcon className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>
                          <div className="font-medium">
                            {"Welcome To Plain."}
                          </div>
                        </div>
                        <div className="rounded-lg p-4 text-left text-muted-foreground hover:bg-accent">
                          {`Hi, welcome to Plain! We're so happy you're here! This
                      message is automated but if you need to talk to us you can
                      press the support button in the bottom left at anytime. In
                      the meantime let's use this thread to show you how Plain
                      works. 🌱 When customers reach out to you, they will show
                      up in your Plain workspace in a thread just like this. 🏷️
                      Each thread has a priority, an assignee, labels, and
                      Linear issues. Use the right-hand panel to set and change
                      those. ✉️ Reply by clicking into the composer below or
                      pressing R. Try sending a reply now. 🚀`}
                        </div>
                      </div>
                    </div>
                    <div className="flex flex-col gap-1">
                      <div className="m-4">
                        <div className="flex items-center font-mono text-sm font-medium">
                          <span className="mr-1 flex h-1 w-1 rounded-full bg-fuchsia-500" />
                          Monday, 14 February 2024
                        </div>
                      </div>
                      {/* message */}
                      <div className="flex flex-col gap-2 rounded-lg border bg-background p-4">
                        <div className="flex w-full flex-col gap-1">
                          <div className="flex items-center">
                            <div className="flex items-center gap-2">
                              <Avatar size={28} name="name" variant="beam" />
                              <div className="font-medium">Emily Davis</div>
                              <span className="flex h-1 w-1 rounded-full bg-blue-600" />
                              <span className="text-xs">3d ago.</span>
                            </div>
                            <div className="ml-auto">
                              <Button variant="ghost" size="icon">
                                <DotsHorizontalIcon className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>
                          <div className="flex items-center text-muted-foreground">
                            <ResetIcon className="mr-1 h-3 w-3" />
                            <div className="text-xs">Welcom to Plain...</div>
                          </div>
                        </div>
                        <div className="rounded-lg p-4 text-left text-muted-foreground hover:bg-accent">
                          {`Nice! 👀 You can see how these messages appear in Slack by pressing O. For more details you can check out our docs.✅ When you are done just hit "Mark as done" on the bottom right.
                      ⌨️ If you want to do anything in Plain use ⌘ + K or CTRL + K on Windows.
                      (edited)`}
                        </div>
                      </div>
                    </div>
                    <div className="flex flex-col gap-1">
                      <div className="m-4">
                        <div className="flex items-center font-mono text-sm font-medium">
                          <span className="mr-1 flex h-1 w-1 rounded-full bg-fuchsia-500" />
                          Monday, 14 February 2024
                        </div>
                      </div>
                      {/* message */}
                      <div className="flex flex-col gap-2 rounded-lg border bg-background p-4">
                        <div className="flex w-full flex-col gap-1">
                          <div className="flex items-center">
                            <div className="flex items-center gap-2">
                              <Avatar size={28} name="name" variant="beam" />
                              <div className="font-medium">Emily Davis</div>
                              <span className="flex h-1 w-1 rounded-full bg-blue-600" />
                              <span className="text-xs">3d ago.</span>
                            </div>
                            <div className="ml-auto">
                              <Button variant="ghost" size="icon">
                                <DotsHorizontalIcon className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>
                          <div className="flex items-center text-muted-foreground">
                            <ResetIcon className="mr-1 h-3 w-3" />
                            <div className="text-xs">Welcom to Plain...</div>
                          </div>
                        </div>
                        <div className="rounded-lg p-4 text-left text-muted-foreground hover:bg-accent">
                          {`Nice! 👀 You can see how these messages appear in Slack by pressing O. For more details you can check out our docs.✅ When you are done just hit "Mark as done" on the bottom right.
                      ⌨️ If you want to do anything in Plain use ⌘ + K or CTRL + K on Windows.
                      (edited)`}
                        </div>
                      </div>
                    </div>
                  </ScrollArea>
                </div>
              </ResizablePanel>
              <ResizableHandle withHandle />
              <ResizablePanel defaultSize={25} maxSize={50} minSize={20}>
                <div className="flex h-full items-center justify-center p-6">
                  <span className="font-semibold">Editor</span>
                </div>
              </ResizablePanel>
            </ResizablePanelGroup>
          </ResizablePanel>
          <ResizableHandle withHandle={false} />
          <ResizablePanel
            defaultSize={25}
            minSize={20}
            maxSize={30}
            className="hidden sm:block"
          >
            <div className="flex h-full items-center justify-center p-6">
              <span className="font-semibold">Sidebar</span>
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
      {/* <div className="flex flex-col justify-center text-center flex-1 p-14 border">
        box 3
      </div>
      <div className="flex flex-col justify-center text-center p-14 border">
        box 3
      </div> */}
    </div>
  );
}
