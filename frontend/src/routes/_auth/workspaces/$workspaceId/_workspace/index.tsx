import * as React from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useStore } from "zustand";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CheckCircleIcon, CircleIcon, EclipseIcon } from "lucide-react";

import { Filters } from "@/components/workspace/filters";
import { Sorts } from "@/components/workspace/sorts";
import { ThreadList } from "@/components/workspace/threads";

import {
  ReasonsFiltersType,
  AssigneesFiltersType,
  PrioritiesFiltersType,
} from "@/db/store";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/_workspace/"
)({
  component: () => <AllThreads />,
});

function AllThreads() {
  const workspaceStore = useWorkspaceStore();
  const navigate = useNavigate();
  const { status, reasons, sort, assignees, priorities } = Route.useSearch();

  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state)
  );
  const todoThreads = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewThreads(
      state,
      "todo",
      assignees as AssigneesFiltersType,
      reasons as ReasonsFiltersType,
      priorities as PrioritiesFiltersType,
      sort
    )
  );
  const doneThreads = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewThreads(
      state,
      "done",
      assignees as AssigneesFiltersType,
      reasons as ReasonsFiltersType,
      priorities as PrioritiesFiltersType,
      sort
    )
  );

  const assignedMembers = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewAssignees(state)
  );

  React.useEffect(() => {
    workspaceStore
      .getState()
      .applyThreadFilters(
        status,
        assignees as AssigneesFiltersType,
        reasons as ReasonsFiltersType,
        priorities as PrioritiesFiltersType,
        sort,
        null,
        null
      );
  }, [workspaceStore, status, assignees, reasons, priorities, sort]);

  return (
    <div className="container">
      <div className="mx-1 my-2 text-xl sm:my-4">All Threads</div>
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
                <CheckCircleIcon className="mr-1 h-4 w-4 text-green-500" />
                Done
              </div>
            </TabsTrigger>
          </TabsList>
          <div className="mt-4 flex gap-1 sm:my-auto">
            <Filters assignedMembers={assignedMembers} />
            <Sorts />
          </div>
        </div>
        <TabsContent value="todo" className="m-0">
          <ThreadList workspaceId={workspaceId} threads={todoThreads} />
        </TabsContent>
        <TabsContent value="snoozed" className="m-0">
          {/* <ThreadList
          items={threads.filter((item) => !item.read)}
          className="h-[calc(100dvh-14rem)]"
        /> */}
        </TabsContent>
        <TabsContent value="done" className="m-0">
          <ThreadList workspaceId={workspaceId} threads={doneThreads} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
