import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useStore } from "zustand";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { Filters } from "@/components/workspace/filters";
import { Sorts } from "@/components/workspace/sorts";
import { ThreadListV3 } from "@/components/workspace/thread-list";

import {
  ReasonsFiltersType,
  AssigneesFiltersType,
  PrioritiesFiltersType,
} from "@/db/store";
import { CircleIcon } from "lucide-react";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/todo"
)({
  component: Threads,
});

function Threads() {
  const workspaceStore = useWorkspaceStore();
  const { reasons, sort, assignees, priorities } = Route.useSearch();

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

  const assignedMembers = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewAssignees(state)
  );

  React.useEffect(() => {
    workspaceStore
      .getState()
      .applyThreadFilters(
        "todo",
        assignees as AssigneesFiltersType,
        reasons as ReasonsFiltersType,
        priorities as PrioritiesFiltersType,
        sort,
        null,
        null
      );
  }, [workspaceStore, assignees, reasons, priorities, sort]);

  return (
    <React.Fragment>
      <div className="px-4 sm:px-8 flex justify-between my-4">
        <div className="text-lg sm:text-xl font-medium items-center">
          <div className="flex items-center gap-x-2">
            <CircleIcon className="h-4 w-4 text-indigo-500 mt-1" />
            <span>Todo</span>
          </div>
        </div>
        <div className="flex gap-1 my-auto">
          <Filters assignedMembers={assignedMembers} />
          <Sorts />
        </div>
      </div>
      <ThreadListV3 workspaceId={workspaceId} threads={todoThreads} />
    </React.Fragment>
  );
}
