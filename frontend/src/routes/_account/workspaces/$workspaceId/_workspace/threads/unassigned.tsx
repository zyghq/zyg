import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Filters } from "@/components/workspace/filters";
import { Sorts } from "@/components/workspace/sorts";
import { ThreadList } from "@/components/workspace/thread-list";
import { STATUS_TODO } from "@/db/constants";
import { setInLocalStorage } from "@/db/helpers";
import { WorkspaceStoreState } from "@/db/store";
import { PrioritiesFiltersType, SortBy, StagesFiltersType } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { PersonIcon } from "@radix-ui/react-icons";
import { createFileRoute } from "@tanstack/react-router";
import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/unassigned",
)({
  component: UnassignedThreads,
});

type PrioritySearchParams = {
  [key: string]: any; // To accommodate other search parameters
  priorities?: string | string[];
};

type SortSearchParams = {
  [key: string]: any; // To accommodate other search parameters
  sort?: string;
};

type StatusSearchParams = {
  [key: string]: any; // To accommodate other search parameters
  stages?: string | string[];
};

function UnassignedThreads() {
  const workspaceStore = useWorkspaceStore();
  const { priorities, sort, stages } = Route.useSearch();
  const navigate = Route.useNavigate();

  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state),
  );
  const todoThreads = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewUnassignedThreads(
      state,
      STATUS_TODO,
      undefined,
      stages as StagesFiltersType,
      priorities as PrioritiesFiltersType,
      sort,
    ),
  );

  React.useEffect(() => {
    workspaceStore
      .getState()
      .applyThreadFilters(
        STATUS_TODO,
        undefined,
        stages as StagesFiltersType,
        priorities as PrioritiesFiltersType,
        sort,
        null,
        true,
      );
  }, [workspaceStore, stages, priorities, sort]);

  React.useEffect(() => {
    setTimeout(() => {
      setInLocalStorage("zyg:threadsQueuePath", Route.fullPath);
    }, 0);
  }, []);

  function onStatusChecked(stage: string) {
    return navigate({
      search: (prev: StatusSearchParams) => {
        const { stages, ...others } = prev;

        if (!stage) {
          return prev; // Return unchanged if input is invalid
        }

        // No existing stages - add new stage
        if (!stages || stages === "") {
          return { stages: stage, ...others };
        }

        // Single stage as string - convert to array and add new stage
        if (typeof stages === "string") {
          if (stages === stage) return prev; // Avoid duplicates
          return { stages: [stages, stage], ...others };
        }

        // Multiple stages in array - add stage if not already present
        if (Array.isArray(stages)) {
          if (stages.includes(stage)) return prev; // Avoid duplicates
          return { stages: [...stages, stage], ...others };
        }

        // Fallback: return unchanged if stages type is unexpected
        return prev;
      },
    });
  }

  function onStatusUnchecked(stage: string) {
    return navigate({
      search: (prev: StatusSearchParams) => {
        const { stages, ...others } = prev;

        if (!stage) {
          return prev; // Return unchanged if input is invalid
        }

        // No existing stages - nothing to do
        if (!stages || stages === "") {
          return { ...others };
        }

        // Single stage as string - remove it if it matches
        if (typeof stages === "string") {
          if (stages === stage) {
            return { ...others };
          }
          return prev; // No change if the stage doesn't match
        }

        // Multiple stages in an array - filter out the stage
        if (Array.isArray(stages)) {
          const filtered = stages.filter((r) => r !== stage);

          // If no stages remain, return others without stages
          if (filtered.length === 0) {
            return { ...others };
          }

          // If only one stage remains, collapse to string
          if (filtered.length === 1) {
            return { stages: filtered[0], ...others };
          }

          // Return the filtered array of stages
          return { stages: filtered, ...others };
        }

        // Fallback: return unchanged if stages type is unexpected
        return prev;
      },
    });
  }

  function onPriorityChecked(priority: string) {
    return navigate({
      search: (prev: PrioritySearchParams) => {
        const { priorities, ...others } = prev;

        if (!priority) {
          return prev; // Return unchanged if input is invalid
        }

        // No existing priorities - add new priority
        if (!priorities || priorities === "") {
          return { priorities: priority, ...others };
        }

        // Single priority as string - convert to array and add new priority
        if (typeof priorities === "string") {
          if (priorities === priority) return prev; // Avoid duplicates
          return { priorities: [priorities, priority], ...others };
        }

        // Multiple priorities in array - add priority if not already present
        if (Array.isArray(priorities)) {
          if (priorities.includes(priority)) return prev; // Avoid duplicates
          return { priorities: [...priorities, priority], ...others };
        }

        // Fallback: return unchanged if priorities type is unexpected
        return prev;
      },
    });
  }

  function onPriorityUnchecked(priority: string) {
    return navigate({
      search: (prev: PrioritySearchParams) => {
        const { priorities, ...others } = prev;

        if (!priority) {
          return prev; // Return unchanged if input is invalid
        }

        // No existing priorities - nothing to do
        if (!priorities || priorities === "") {
          return { ...others };
        }

        // Single priority as string - remove it if it matches
        if (typeof priorities === "string") {
          if (priorities === priority) {
            return { ...others };
          }
          return prev; // No change if the priority doesn't match
        }

        // Multiple priorities in an array - filter out the priority
        if (Array.isArray(priorities)) {
          const filtered = priorities.filter((r) => r !== priority);

          // If no priorities remain, return others without priorities
          if (filtered.length === 0) {
            return { ...others };
          }

          // If only one priority remains, collapse to string
          if (filtered.length === 1) {
            return { priorities: filtered[0], ...others };
          }

          // Return the filtered array of priorities
          return { priorities: filtered, ...others };
        }

        // Fallback: return unchanged if priorities type is unexpected
        return prev;
      },
    });
  }

  function onSortChecked(sort: string) {
    return navigate({
      search: (prev: SortSearchParams) => {
        // If sort is different from the previous one, update it
        if (prev.sort !== sort) {
          return { ...prev, sort };
        }
        // Return unchanged if the sort is the same
        return prev;
      },
    });
  }

  return (
    <React.Fragment>
      <header className="flex h-14 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger className="-ml-1" />
        <Separator className="mr-2 h-4" orientation="vertical" />
        <div className="flex-1"></div>
        <div className="flex space-x-1">
          <Filters
            assignedMembers={[]}
            assignees={undefined}
            disableAssigneeFilter={true}
            priorities={priorities as PrioritiesFiltersType}
            priorityOnChecked={onPriorityChecked}
            priorityOnUnchecked={onPriorityUnchecked}
            stages={stages as StagesFiltersType}
            statusOnChecked={onStatusChecked}
            statusOnUnchecked={onStatusUnchecked}
          />
          <Sorts onChecked={onSortChecked} sort={sort as SortBy} />
        </div>
      </header>

      <div className="p-4">
        <div className="flex items-center space-x-2">
          <PersonIcon className="h-5 w-5" />
          <span className="font-serif text-lg font-medium">
            {"Unassigned Threads"}
          </span>
        </div>
      </div>

      <ThreadList threads={todoThreads} workspaceId={workspaceId} />
    </React.Fragment>
  );
}
