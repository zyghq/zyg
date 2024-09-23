import { Filters } from "@/components/workspace/filters";
import { Sorts } from "@/components/workspace/sorts";
import { ThreadListV3 } from "@/components/workspace/thread-list";
import { STATUS_TODO, WAITING_ON_CUSTOMER } from "@/db/constants";
import { setInLocalStorage } from "@/db/helpers";
import { WorkspaceStoreState } from "@/db/store";
import {
  AssigneesFiltersType,
  PrioritiesFiltersType,
  SortBy,
} from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { createFileRoute } from "@tanstack/react-router";
import { ClockIcon } from "lucide-react";
import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/waiting-on-customer"
)({
  component: Threads,
});

function Threads() {
  const workspaceStore = useWorkspaceStore();
  const { assignees, priorities, sort } = Route.useSearch();
  const navigate = Route.useNavigate();

  const workspaceId = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceId(state)
  );
  const threads = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewThreads(
      state,
      STATUS_TODO,
      assignees as AssigneesFiltersType,
      WAITING_ON_CUSTOMER,
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
        STATUS_TODO,
        assignees as AssigneesFiltersType,
        WAITING_ON_CUSTOMER,
        priorities as PrioritiesFiltersType,
        sort,
        null,
        null
      );
  }, [workspaceStore, assignees, priorities, sort]);

  React.useEffect(() => {
    setTimeout(() => {
      setInLocalStorage("zyg:threadsQueuePath", Route.fullPath);
    }, 0);
  }, []);

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

  function onAssigneeChecked(member: string) {
    return navigate({
      search: (prev) => {
        // search params
        const { assignees, ...others } = prev;

        // no existing members - add new member
        if (!assignees || assignees === "") {
          return { assignees: member, ...others };
        }

        // found a member - merge with existing
        if (typeof assignees === "string") {
          return { assignees: [assignees, member], ...others };
        }
        // multiple members selected add more to existing
        if (Array.isArray(assignees)) {
          const uniques = [...new Set([member, ...assignees])];
          return { assignees: uniques, ...others };
        }
        return prev;
      },
    });
  }

  function onAssigneeUnchecked(member: string) {
    return navigate({
      search: (prev) => {
        const { assignees, ...others } = prev;

        // no existing members - nothing to do
        if (!assignees || assignees === "") {
          return { ...others };
        }

        // found a member - remove it
        if (typeof assignees === "string" && assignees === member) {
          return { ...others };
        }

        // multiple members selected - remove the member
        if (Array.isArray(assignees)) {
          const filtered = assignees.filter((r) => r !== member);
          if (filtered.length === 0) {
            return { ...others };
          }
          if (filtered.length === 1) {
            return { assignees: filtered[0], ...others };
          }
          return { assignees: filtered, ...others };
        }
        return prev;
      },
    });
  }

  function onSortChecked(sort: string) {
    return navigate({
      search: (prev) => ({ ...prev, sort }),
    });
  }

  return (
    <React.Fragment>
      <div className="px-4 sm:px-8 flex justify-between my-4">
        <div className="text-lg sm:text-xl font-medium items-center">
          <div className="flex items-center gap-x-2">
            <ClockIcon className="my-auto w-5 h-5 text-indigo-500" />
            <span>Waiting on Customer</span>
          </div>
        </div>
        <div className="flex gap-1 my-auto">
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
      </div>
      {threads && threads.length > 0 ? (
        <ThreadListV3 threads={threads} workspaceId={workspaceId} />
      ) : (
        <div>
          <div className="container mt-4 sm:mt-24 max-w-md">
            <div className="border p-4 rounded-xl">
              <div className="font-medium text-sm">
                No threads are waiting on the customer.
              </div>
              <div className="text-muted-foreground text-sm">
                When you are waiting for a reply from a customer in order to
                resolve an issue, you can mark a thread as waiting on customer
                and it will show up here
              </div>
            </div>
          </div>
        </div>
      )}
    </React.Fragment>
  );
}
