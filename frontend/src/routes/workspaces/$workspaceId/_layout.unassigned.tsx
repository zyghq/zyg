import { createFileRoute } from "@tanstack/react-router";
import { useStore } from "zustand";
import { WorkspaceStoreStateType } from "@/db/store";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { DoubleArrowUpIcon } from "@radix-ui/react-icons";
import { CheckCircle, CircleIcon, EclipseIcon } from "lucide-react";

import { ThreadList } from "@/components/workspace/threads";

export const Route = createFileRoute(
  "/workspaces/$workspaceId/_layout/unassigned"
)({
  component: () => <UnassignedThreads />,
});

function UnassignedThreads() {
  const { workspaceStore } = Route.useRouteContext();
  const workspaceId = useStore(
    workspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getWorkspaceId(state)
  );
  const threads = useStore(
    workspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.viewUnassignedThreads(state)
  );
  return (
    <main className="col-span-3 lg:col-span-4">
      <div className="container">
        <div className="mb-4 mt-4 text-xl">Unassigned Threads</div>
        <Tabs defaultValue="todo">
          <div className="mb-4 sm:flex sm:justify-between">
            <TabsList className="grid grid-cols-3">
              <TabsTrigger value="todo">
                <div className="flex items-center">
                  <CircleIcon className="mr-1 h-4 w-4 text-indigo-500" />
                  Todo
                </div>
              </TabsTrigger>
              <TabsTrigger value="snoozed">
                <div className="flex items-center">
                  <EclipseIcon className="mr-1 h-4 w-4 text-fuchsia-500" />
                  Snoozed
                </div>
              </TabsTrigger>
              <TabsTrigger value="done">
                <div className="flex items-center">
                  <CheckCircle className="mr-1 h-4 w-4 text-green-500" />
                  Done
                </div>
              </TabsTrigger>
            </TabsList>
            <div className="mt-4 flex gap-1 sm:my-auto">
              ...
              {/* <ThreadFilterDropDownMenu /> */}
              <Button variant="outline" size="sm" className="border-dashed">
                <DoubleArrowUpIcon className="mr-1 h-3 w-3" />
                Sort
              </Button>
            </div>
          </div>
          <TabsContent value="todo" className="m-0">
            <ThreadList
              workspaceId={workspaceId}
              threads={threads}
              className="h-[calc(100dvh-14rem)]"
            />
          </TabsContent>
          <TabsContent value="snoozed" className="m-0">
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
      </div>
    </main>
  );
}
