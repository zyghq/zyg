import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useStore } from "zustand";
import { WorkspaceStoreState } from "@/db/store";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CheckCircle, CircleIcon, EclipseIcon } from "lucide-react";

import { ThreadList } from "@/components/workspace/threads";
import { ThreadListV2 } from "@/components/workspace/threads-v2";
import { Filters } from "@/components/workspace/filters";
import { Sorts } from "@/components/workspace/sorts";
import {
  ReasonsFiltersType,
  AssigneesFiltersType,
  PrioritiesFiltersType,
} from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import * as React from "react";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/_workspace/unassigned"
)({
  component: () => <UnassignedThreads />,
});

function UnassignedThreads() {
  const workspaceStore = useWorkspaceStore();
  const { status, reasons, sort, assignees, priorities } = Route.useSearch();
  const navigate = useNavigate();

  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state)
  );
  const todoThreads = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewUnassignedThreads(
      state,
      "todo",
      assignees as AssigneesFiltersType,
      reasons as ReasonsFiltersType,
      priorities as PrioritiesFiltersType,
      sort
    )
  );

  const doneThreads = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewUnassignedThreads(
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
        true
      );
  }, [workspaceStore, status, assignees, reasons, priorities, sort]);

  return (
    <div>
      <div className="px-4 sm:px-8 flex justify-between my-4">
        <div className="text-lg sm:text-xl font-medium my-auto">
          Unassigned Threads
        </div>
        <div className="flex gap-1 my-auto">
          <Filters assignedMembers={assignedMembers} />
          <Sorts />
        </div>
      </div>
      <ThreadListV2 workspaceId={workspaceId} threads={todoThreads} />
    </div>
  );
}
