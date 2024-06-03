import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useStore } from "zustand";
import { WorkspaceStoreStateType } from "@/db/store";
import { useWorkspaceStore } from "@/providers";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CheckCircle, CircleIcon, EclipseIcon } from "lucide-react";

import { Filters } from "@/components/workspace/filters";
import { Sorts } from "@/components/workspace/sorts";
import { ThreadList } from "@/components/workspace/threads";

import { reasonsFiltersType } from "@/db/store";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/_workspace/"
)({
  component: () => <AllThreads />,
});

function AllThreads() {
  const workspaceStore = useWorkspaceStore();
  const navigate = useNavigate();
  const { status, reasons, sort } = Route.useSearch();

  const workspaceId = useStore(
    workspaceStore,
    (state: WorkspaceStoreStateType) => state.getWorkspaceId(state)
  );
  const threads = useStore(workspaceStore, (state: WorkspaceStoreStateType) =>
    state.viewAllTodoThreads(state, reasons as reasonsFiltersType, sort)
  );
  return (
    <main className="col-span-3 lg:col-span-4">
      <div className="container">
        <div className="mb-4 mt-4 text-xl">All Threads</div>
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
              <Sorts />
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
