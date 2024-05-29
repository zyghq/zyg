import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useStore } from "zustand";
import { WorkspaceStoreStateType } from "@/db/store";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { DoubleArrowUpIcon } from "@radix-ui/react-icons";
import { CheckCircle, CircleIcon, EclipseIcon } from "lucide-react";

import { Filters } from "@/components/workspace/filters";
import { ThreadList } from "@/components/workspace/threads";

export const Route = createFileRoute("/workspaces/$workspaceId/_layout/me")({
  component: () => <MyThreads />,
});

function MyThreads() {
  const { WorkspaceStore } = Route.useRouteContext();

  const { status } = Route.useSearch();
  const navigate = useNavigate();

  const workspaceId = useStore(
    WorkspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getWorkspaceId(state)
  );

  const memberId = useStore(
    WorkspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getMemberId(state)
  );
  const threads = useStore(
    WorkspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.viewMyTodoThreads(state, memberId)
  );
  return (
    <main className="col-span-3 lg:col-span-4">
      <div className="container">
        <div className="mb-4 mt-4 text-xl">My Threads</div>
        <Tabs defaultValue={status}>
          <div className="mb-4 sm:flex sm:justify-between">
            <TabsList className="grid grid-cols-3">
              <TabsTrigger
                onClick={() => {
                  navigate({ search: () => ({ status: "todo" }) });
                }}
                value="todo"
              >
                <div className="flex items-center">
                  <CircleIcon className="mr-1 h-4 w-4 text-indigo-500" />
                  Todo
                </div>
              </TabsTrigger>
              <TabsTrigger
                onClick={() => {
                  navigate({ search: () => ({ status: "snoozed" }) });
                }}
                value="snoozed"
              >
                <div className="flex items-center">
                  <EclipseIcon className="mr-1 h-4 w-4 text-fuchsia-500" />
                  Snoozed
                </div>
              </TabsTrigger>
              <TabsTrigger
                onClick={() => {
                  navigate({ search: () => ({ status: "done" }) });
                }}
                value="done"
              >
                <div className="flex items-center">
                  <CheckCircle className="mr-1 h-4 w-4 text-green-500" />
                  Done
                </div>
              </TabsTrigger>
            </TabsList>
            <div className="mt-4 flex gap-1 sm:my-auto">
              <Filters />
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
