import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { DoubleArrowUpIcon, MixerHorizontalIcon } from "@radix-ui/react-icons";
import { EclipseIcon, CircleIcon, CheckCircle } from "lucide-react";
import ThreadList from "@/components/thread-list";
// import { ScrollArea } from "@/components/ui/scroll-area";

// mock data for now
import { threads } from "@/data/threads";

export default async function ThreadListTab({ threads }) {
  return (
    <Tabs defaultValue="todo">
      <div className="mb-4 sm:flex sm:justify-between">
        <TabsList className="grid grid-cols-3">
          <TabsTrigger value="todo">
            <div className="flex items-center">
              <CircleIcon className="mr-1 h-3 w-3 text-indigo-500" />
              Todo
            </div>
          </TabsTrigger>
          <TabsTrigger value="inprogress">
            <div className="flex items-center">
              <EclipseIcon className="mr-1 h-3 w-3 text-fuchsia-500" />
              In Progress
            </div>
          </TabsTrigger>
          <TabsTrigger value="done">
            <div className="flex items-center">
              <CheckCircle className="mr-1 h-3 w-3 text-green-500" />
              Done
            </div>
          </TabsTrigger>
        </TabsList>
        <div className="mt-4 sm:my-auto">
          <Button variant="ghost" size="sm">
            <MixerHorizontalIcon className="mr-1 h-3 w-3" />
            Filters
          </Button>
          <Button variant="ghost" size="sm">
            <DoubleArrowUpIcon className="mr-1 h-3 w-3" />
            Sort
          </Button>
        </div>
      </div>
      <TabsContent value="todo" className="m-0">
        <ThreadList items={threads} className="h-[calc(100dvh-14rem)]" />
      </TabsContent>
      <TabsContent value="inprogress" className="m-0">
        {/* <ThreadList
          items={threads.filter((item) => !item.read)}
          className="h-[calc(100dvh-14rem)]"
        /> */}
      </TabsContent>
      <TabsContent value="done" className="m-0">
        {/* <Threads items={threads} /> */}
        {/* <ScrollArea className="h-[calc(100vh-14rem)] pr-1">
          <div className="flex flex-col gap-2">...</div>
        </ScrollArea> */}
      </TabsContent>
    </Tabs>
  );
}
