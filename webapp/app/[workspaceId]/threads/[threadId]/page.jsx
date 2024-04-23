// import { ThreadSidebar } from "@/components/threadsidebar";
// import Title from "@/components/title";
// import ThreadTabs from "@/components/thread-tabs";
// import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  ChevronLeftIcon,
  ArrowUpIcon,
  ArrowDownIcon,
  ChatBubbleIcon,
  DotsHorizontalIcon,
  ResetIcon,
} from "@radix-ui/react-icons";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { CircleIcon } from "lucide-react";
import ThreadList from "@/components/thread-list";
import { threads } from "@/data/threads";

export default function ThreadItemPage() {
  return (
    <div className="flex flex-1">
      <div className="flex flex-col items-center border px-2">
        <div className="mt-4 flex flex-col gap-4">
          <Button variant="outline" size="icon">
            <ChevronLeftIcon className="h-4 w-4" />
          </Button>
          <Button variant="outline" size="icon">
            <ArrowUpIcon className="h-4 w-4" />
          </Button>
          <Button variant="outline" size="icon">
            <ArrowDownIcon className="h-4 w-4" />
          </Button>
        </div>
      </div>
      <div className="flex flex-1 flex-col">
        <ResizablePanelGroup direction="horizontal">
          <ResizablePanel defaultSize={25} minSize={20} maxSize={30}>
            <div className="flex h-full flex-col">
              <div className="flex h-14 flex-col justify-center border-b px-4">
                <div className="text-md font-semibold">Threads</div>
              </div>
              {/* <ThreadList
                items={threads}
                className="h-[calc(100dvh-8rem)] pr-0"
                variant="compress"
              /> */}
              ...
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
                        Emily Davis via Chat
                      </div>
                    </div>
                    <div className="flex items-center">
                      <CircleIcon className="mr-1 h-3 w-3 text-indigo-500" />
                      <span className="items-center text-xs">Todo</span>
                      <Separator orientation="vertical" className="mx-2" />
                      <ChatBubbleIcon className="h-3 w-3" />
                      <Separator orientation="vertical" className="mx-2" />
                      <span className="font-mono text-xs">12/44</span>
                    </div>
                  </div>
                  <ScrollArea className="flex h-full flex-auto flex-col px-12 pb-4">
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
                              <Avatar className="h-8 w-8">
                                <AvatarImage src="https://github.com/shadcn.png" />
                                <AvatarFallback>CN</AvatarFallback>
                              </Avatar>
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
                              <Avatar className="h-8 w-8">
                                <AvatarImage src="https://github.com/shadcn.png" />
                                <AvatarFallback>CN</AvatarFallback>
                              </Avatar>
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
                              <Avatar className="h-8 w-8">
                                <AvatarImage src="https://github.com/shadcn.png" />
                                <AvatarFallback>CN</AvatarFallback>
                              </Avatar>
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
          <ResizableHandle withHandle />
          <ResizablePanel defaultSize={25} minSize={20} maxSize={30}>
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
