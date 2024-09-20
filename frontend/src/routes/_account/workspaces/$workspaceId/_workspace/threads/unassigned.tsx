import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useStore } from "zustand";
import { WorkspaceStoreState } from "@/db/store";

import { Filters } from "@/components/workspace/filters";
import { Sorts } from "@/components/workspace/sorts";
import { StagesFiltersType, PrioritiesFiltersType } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { ThreadListV3 } from "@/components/workspace/thread-list";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/unassigned"
)({
  component: UnassignedThreads,
});

function UnassignedThreads() {
  const workspaceStore = useWorkspaceStore();
  const { stages, sort, priorities } = Route.useSearch();
  const navigate = Route.useNavigate();

  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state)
  );
  const todoThreads = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewUnassignedThreads(
      state,
      "todo",
      undefined,
      stages as StagesFiltersType,
      priorities as PrioritiesFiltersType,
      sort
    )
  );

  React.useEffect(() => {
    workspaceStore
      .getState()
      .applyThreadFilters(
        "todo",
        undefined,
        stages as StagesFiltersType,
        priorities as PrioritiesFiltersType,
        sort,
        null,
        true
      );
  }, [workspaceStore, stages, priorities, sort]);

  function onStatusChecked(stage: string) {
    return navigate({
      search: (prev) => {
        const { stages, ...others } = prev;

        // no existing stages - add new stage
        if (!stages || stages === "") {
          return { stages: stage, ...others };
        }

        // found a stage - merge with existing
        if (typeof stages === "string") {
          return { stages: [stages, stage], ...others };
        }
        // multiple stages selected add more to existing
        if (Array.isArray(stages)) {
          return { stages: [...stages, stage], ...others };
        }
        // return without side effects
        return prev;
      },
    });
  }

  function onStatusUnchecked(stage: string) {
    return navigate({
      search: (prev) => {
        const { stages, ...others } = prev;

        // no existing stages - nothing to do
        if (!stages || stages === "") {
          return { ...others };
        }

        // found a stage - remove it
        if (typeof stages === "string" && stages === stage) {
          return { ...others };
        }

        // multiple stages selected - remove the stage
        if (Array.isArray(stages)) {
          const filtered = stages.filter((r) => r !== stage);
          if (filtered.length === 0) {
            return { ...others };
          }
          if (filtered.length === 1) {
            return { stages: filtered[0], ...others };
          }
          return { stages: filtered, ...others };
        }

        // return without side effects
        return prev;
      },
    });
  }

  function onPriorityChecked(priority: string) {
    return navigate({
      search: (prev) => {
        const { priorities, ...others } = prev;

        // no existing priorities - add new priority
        if (!priorities || priorities === "") {
          return { priorities: priority, ...others };
        }

        // found a priority - merge with existing
        if (typeof priorities === "string") {
          return { priorities: [priorities, priority], ...others };
        }
        // multiple priorities selected add more to existing
        if (Array.isArray(priorities)) {
          return { priorities: [...priorities, priority], ...others };
        }
        // return without side effects
        return prev;
      },
    });
  }

  function onPriorityUnchecked(priority: string) {
    return navigate({
      search: (prev) => {
        const { priorities, ...others } = prev;

        // no existing priorities - nothing to do
        if (!priorities || priorities === "") {
          return { ...others };
        }

        // found a priority - remove it
        if (typeof priorities === "string" && priorities === priority) {
          return { ...others };
        }

        // multiple priorities selected - remove the priority
        if (Array.isArray(priorities)) {
          const filtered = priorities.filter((r) => r !== priority);
          if (filtered.length === 0) {
            return { ...others };
          }
          if (filtered.length === 1) {
            return { priorities: filtered[0], ...others };
          }
          return { priorities: filtered, ...others };
        }
        // return without side effects
        return prev;
      },
    });
  }

  return (
    <React.Fragment>
      <div className="px-4 sm:px-8 flex justify-between my-4">
        <div className="text-lg sm:text-xl font-medium my-auto">
          Unassigned Threads
        </div>
        <div className="flex gap-1 my-auto">
          <Filters
            stages={stages as StagesFiltersType}
            priorities={priorities as PrioritiesFiltersType}
            statusOnChecked={onStatusChecked}
            statusOnUnchecked={onStatusUnchecked}
            priorityOnChecked={onPriorityChecked}
            priorityOnUnchecked={onPriorityUnchecked}
            assignedMembers={[]}
            assignees={undefined}
            disableAssigneeFilter={true}
          />
          <Sorts />
        </div>
      </div>
      <ThreadListV3 workspaceId={workspaceId} threads={todoThreads} />
    </React.Fragment>
  );
}
