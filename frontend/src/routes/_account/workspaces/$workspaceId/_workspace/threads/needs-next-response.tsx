import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Filters } from "@/components/workspace/filters";
import { Sorts } from "@/components/workspace/sorts";
import { ThreadList } from "@/components/workspace/thread-list";
import { NEEDS_NEXT_RESPONSE, STATUS_TODO } from "@/db/constants";
import { setInLocalStorage } from "@/db/helpers";
import { WorkspaceStoreState } from "@/db/store";
import {
  AssigneesFiltersType,
  PrioritiesFiltersType,
  SortBy,
} from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { createFileRoute } from "@tanstack/react-router";
import { ReplyIcon } from "lucide-react";
import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/needs-next-response",
)({
  component: Threads,
});

type AssigneeSearchParams = {
  [key: string]: any; // To accommodate other search parameters
  assignees?: string | string[];
};

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

function Threads() {
  const workspaceStore = useWorkspaceStore();
  const { assignees, priorities, sort } = Route.useSearch();
  const navigate = Route.useNavigate();

  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state),
  );
  const threads = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewThreads(
      state,
      STATUS_TODO,
      assignees as AssigneesFiltersType,
      NEEDS_NEXT_RESPONSE,
      priorities as PrioritiesFiltersType,
      sort,
    ),
  );

  const assignedMembers = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewAssignees(state),
  );

  React.useEffect(() => {
    workspaceStore
      .getState()
      .applyThreadFilters(
        STATUS_TODO,
        assignees as AssigneesFiltersType,
        NEEDS_NEXT_RESPONSE,
        priorities as PrioritiesFiltersType,
        sort,
        null,
        null,
      );
  }, [workspaceStore, assignees, priorities, sort]);

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

  function onAssigneeChecked(member: string) {
    return navigate({
      search: (prev: AssigneeSearchParams) => {
        const { assignees, ...others } = prev;

        if (!member) {
          return prev; // Return unchanged if input is invalid
        }

        // No existing members - add the new member
        if (!assignees || assignees === "") {
          return { assignees: member, ...others };
        }

        // Single member as string - convert to array and add the new member
        if (typeof assignees === "string") {
          if (assignees === member) return prev; // Avoid duplicates
          return { assignees: [assignees, member], ...others };
        }

        // Multiple members in an array - add the new member if not already present
        if (Array.isArray(assignees)) {
          if (assignees.includes(member)) return prev; // Avoid duplicates
          return { assignees: [...assignees, member], ...others };
        }

        // Fallback: return unchanged if assignees type is unexpected
        return prev;
      },
    });
  }

  function onAssigneeUnchecked(member: string) {
    return navigate({
      search: (prev: AssigneeSearchParams) => {
        const { assignees, ...others } = prev;

        // Case 1: No existing assignees, nothing to do
        if (!assignees || assignees === "") {
          return { ...others }; // Return without modifying if there are no assignees
        }

        // Case 2: Single assignee as string - remove if it matches
        if (typeof assignees === "string") {
          if (assignees === member) {
            return { ...others }; // Remove assignee if it matches the member
          }
          return prev; // Return unchanged if member does not match
        }

        // Case 3: Multiple assignees in an array - filter out the member
        if (Array.isArray(assignees)) {
          const filtered = assignees.filter((r) => r !== member);

          // Case 3a: No members left, return without assignees
          if (filtered.length === 0) {
            return { ...others };
          }

          // Case 3b: Only one member left, return as a string
          if (filtered.length === 1) {
            return { assignees: filtered[0], ...others };
          }

          // Case 3c: Multiple members remain, return the filtered array
          return { assignees: filtered, ...others };
        }

        // Fallback: Return unchanged if assignees is of an unexpected type
        console.warn("Unexpected assignees type detected.");
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
            assignedMembers={assignedMembers}
            assigneeOnChecked={onAssigneeChecked}
            assigneeOnUnchecked={onAssigneeUnchecked}
            assignees={assignees as AssigneesFiltersType}
            disableStagesFilter={true}
            priorities={priorities as PrioritiesFiltersType}
            priorityOnChecked={onPriorityChecked}
            priorityOnUnchecked={onPriorityUnchecked}
            stages={undefined}
            statusOnChecked={onStatusChecked}
            statusOnUnchecked={onStatusUnchecked}
          />
          <Sorts onChecked={onSortChecked} sort={sort as SortBy} />
        </div>
      </header>
      <div className="p-4">
        <div className="flex items-center space-x-2">
          <ReplyIcon className="h-5 w-5 text-indigo-500" />
          <span className="font-serif text-lg font-medium sm:text-xl">
            {"Needs Next Response"}
          </span>
        </div>
      </div>
      <ThreadList threads={threads} workspaceId={workspaceId} />
    </React.Fragment>
  );
}
