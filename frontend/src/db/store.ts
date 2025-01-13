import { MemberShape, WorkspaceShape } from "@/db/shapes";
import _ from "lodash";
import { immer } from "zustand/middleware/immer";
import { createStore } from "zustand/vanilla";

import {
  HOLD,
  NEEDS_FIRST_RESPONSE,
  NEEDS_NEXT_RESPONSE,
  RESOLVED,
  SPAM,
  WAITING_ON_CUSTOMER,
} from "./constants";
import { getFromLocalStorage, setInLocalStorage } from "./helpers";
import {
  Account,
  AuthMember,
  Customer,
  Label,
  Pat,
  Thread,
  WorkspaceMetrics,
} from "./models";

export type LabelMap = Dictionary<string, Label>;

export type MemberShapeMap = Dictionary<string, MemberShape>;

export type PatMap = Dictionary<string, Pat>;

export type CustomerMap = Dictionary<string, Customer>;

export type Dictionary<
  K extends number | string,
  V extends WorkspaceEntities,
> = {
  [key in K]: V;
};

export interface IWorkspaceEntities {
  customers: CustomerMap | null;
  labels: LabelMap | null;
  member: AuthMember | null;
  members: MemberShapeMap | null;
  metrics: WorkspaceMetrics;
  pats: null | PatMap;
  threads: null | ThreadMap;
  workspace: WorkspaceShape;
}

// Represents API bootstrapped entities.
export interface IWorkspaceEntitiesBootstrap {
  customers: CustomerMap | null;
  labels: LabelMap | null;
  member: AuthMember | null;
  metrics: WorkspaceMetrics;
  pats: null | PatMap;
  threads: null | ThreadMap;
}

export type SortBy =
  | "created-asc"
  | "created-dsc"
  | "inbound-message-dsc"
  | "outbound-message-dsc"
  | "priority-asc"
  | "priority-dsc"
  | "status-changed-asc"
  | "status-changed-dsc";

export type ThreadMap = Dictionary<string, Thread>;

type Priority = "high" | "low" | "normal" | "urgent";
type StageType =
  | "hold"
  | "needs_first_response"
  | "needs_next_response"
  | "resolved"
  | "spam"
  | "waiting_on_customer";

type StatusType = "done" | "todo";

// add more entities as supported by store
// e.g: Workspace | User | etc.
type WorkspaceEntities =
  | Customer
  | Label
  | MemberShape
  | Pat
  | Thread
  | WorkspaceShape;

export const defaultSortKey = "created-asc";

export type Assignee = {
  assigneeId: string;
  name: string;
};
export type AssigneesFiltersType = string | string[] | undefined;

export interface IWorkspaceValueObjects {
  error: Error | null;
  hasData: boolean;
  isPending: boolean;
  threadAppliedFilters: null | ThreadAppliedFilters;
  threadSortKey: null | SortBy;
}

export type PrioritiesFiltersType = Priority | Priority[] | undefined;

// export type sortByType = "last-message-dsc" | "created-asc" | "created-dsc";
// export const defaultSortKey = "last-message-dsc";

export type StagesFiltersType = StageType | StageType[] | undefined;

export type ThreadAppliedFilters = {
  assignees: AssigneesFiltersType;
  isUnassigned?: boolean | null;
  memberId?: null | string;
  priorities: PrioritiesFiltersType;
  sortBy: SortBy;
  stages: StagesFiltersType;
  status: StatusType;
};

export type WorkspaceStoreState = IWorkspaceEntities &
  IWorkspaceStoreActions &
  IWorkspaceValueObjects;

interface IWorkspaceStoreActions {
  addLabel(label: Label): void;

  addPat(pat: Pat): void;

  applyThreadFilters(
    status: StatusType,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: SortBy,
    memberId: null | string,
    isUnassigned: boolean | null,
  ): void;

  deletePat(patId: string): void;

  getMemberId(state: WorkspaceStoreState): string;

  getMemberName(state: WorkspaceStoreState): string;

  getMetrics(state: WorkspaceStoreState): WorkspaceMetrics;

  getThreadItem(state: WorkspaceStoreState, threadId: string): null | Thread;

  getWorkspaceId(state: WorkspaceStoreState): string;

  getWorkspaceName(state: WorkspaceStoreState): string;

  setThreadSortKey(sortKey: SortBy): void;

  updateLabel(labelId: string, label: Label): void;

  updateThread(thread: Thread): void;

  updateWorkspaceName(name: string): void;

  viewAssignees(state: WorkspaceStoreState): Assignee[];

  viewCurrentThreadQueue(state: WorkspaceStoreState): null | Thread[];

  viewCustomerEmail(
    state: WorkspaceStoreState,
    customerId: string,
  ): null | string;

  viewCustomerExternalId(
    state: WorkspaceStoreState,
    customerId: string,
  ): null | string;

  viewCustomerName(state: WorkspaceStoreState, customerId: string): string;

  viewCustomerPhone(
    state: WorkspaceStoreState,
    customerId: string,
  ): null | string;

  viewCustomerRole(state: WorkspaceStoreState, customerId: string): string;

  viewLabels(state: WorkspaceStoreState): Label[];

  viewMemberName(state: WorkspaceStoreState, memberId: string): string;

  viewMembers(state: WorkspaceStoreState): MemberShape[];

  viewMyThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    memberId: string,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: SortBy,
  ): Thread[];

  viewPats(state: WorkspaceStoreState): Pat[];

  viewThreadAssigneeId(
    state: WorkspaceStoreState,
    threadId: string,
  ): null | string;

  viewThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: SortBy,
  ): Thread[];

  viewThreadSortKey(state: WorkspaceStoreState): SortBy;

  viewUnassignedThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: SortBy,
  ): Thread[];
}

function filterByAssignees(threads: Thread[], assignees: AssigneesFiltersType) {
  if (assignees && Array.isArray(assignees)) {
    const uniqueAssignees = [...new Set(assignees)];
    const filtered = [];
    for (const assignee of uniqueAssignees) {
      filtered.push(...threads.filter((t) => t.assigneeId === assignee));
    }
    return filtered;
  }
  if (assignees) {
    return threads.filter((t) => t.assigneeId === assignees);
  }
  // no change
  return threads;
}

function filterByPriorities(
  threads: Thread[],
  priorities: PrioritiesFiltersType,
) {
  if (priorities && Array.isArray(priorities)) {
    const uniquePriorities = [...new Set(priorities)];
    const filtered = [];
    for (const priority of uniquePriorities) {
      filtered.push(...threads.filter((t) => t.priority === priority));
    }
    return filtered;
  }

  if (priorities) {
    return threads.filter((t) => t.priority === priorities);
  }
  // no change
  return threads;
}

function filterByStages(threads: Thread[], stages: StagesFiltersType) {
  const stageMap: Record<string, string> = {
    hold: HOLD,
    needs_first_response: NEEDS_FIRST_RESPONSE,
    needs_next_response: NEEDS_NEXT_RESPONSE,
    resolved: RESOLVED,
    spam: SPAM,
    waiting_on_customer: WAITING_ON_CUSTOMER,
  };

  if (stages) {
    if (Array.isArray(stages)) {
      const uniqueStages = [...new Set(stages)];
      return threads.filter((t) =>
        uniqueStages.some((stage) => t.stage === stageMap[stage]),
      );
    }
    if (stageMap[stages]) {
      return threads.filter((t) => t.stage === stageMap[stages]);
    }
  }

  // no change
  return threads;
}

function sortThreads(threads: Thread[], sortBy: SortBy): Thread[] {
  const priorityMap: { [key: string]: number } = {
    high: 1,
    low: 3,
    normal: 2,
    urgent: 0,
  };
  switch (sortBy) {
    case "created-asc":
      return _.sortBy(threads, "createdAt");
    case "created-dsc":
      return _.sortBy(threads, "createdAt").reverse();
    case "inbound-message-dsc":
      return _.sortBy(threads, "inboundLastSeqId").reverse();
    case "outbound-message-dsc":
      return _.sortBy(threads, "outboundLastSeqId").reverse();
    case "priority-asc":
      return _.sortBy(threads, (thread) => priorityMap[thread.priority]);
    case "priority-dsc":
      return _.sortBy(threads, (thread) => -priorityMap[thread.priority]);
    case "status-changed-asc":
      return _.sortBy(threads, "statusChangedAt");
    case "status-changed-dsc":
      return _.sortBy(threads, "statusChangedAt").reverse();
    default:
      return threads;
  }
}

// (sanchitrk) for reference on using zustand, check this great article:
// https://tkdodo.eu/blog/working-with-zustand
export const buildWorkspaceStore = (
  initialState: IWorkspaceEntities & IWorkspaceValueObjects,
) => {
  return createStore<WorkspaceStoreState>()(
    immer((set) => ({
      ...initialState,
      addLabel: (label: Label) => {
        const { labelId } = label;
        set((state) => {
          if (state.labels) {
            state.labels[labelId] = { ...label };
            return state;
          } else {
            state.labels = { [labelId]: { ...label } };
            return state;
          }
        });
      },
      addPat: (pat: Pat) => {
        const { patId } = pat;
        set((state) => {
          if (state.pats) {
            state.pats[patId] = { ...pat };
            return state;
          } else {
            state.pats = { [patId]: { ...pat } };
            return state;
          }
        });
      },
      applyThreadFilters: (
        status: StatusType,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: SortBy,
        memberId: null | string,
        isUnassigned: boolean | null,
      ) => {
        set((state) => {
          if (state.threadAppliedFilters) {
            state.threadAppliedFilters = {
              ...state.threadAppliedFilters,
              assignees,
              isUnassigned,
              memberId,
              priorities,
              sortBy,
              stages,
              status,
            };
          } else {
            state.threadAppliedFilters = {
              assignees,
              isUnassigned,
              memberId,
              priorities,
              sortBy,
              stages,
              status,
            };
          }
        });
      },
      deletePat: (patId: string) => {
        set((state) => {
          if (state.pats) {
            delete state.pats[patId];
            return state;
          } else {
            return state;
          }
        });
      },
      getMemberId: (state: WorkspaceStoreState) => state.member?.memberId || "",
      getMemberName: (state: WorkspaceStoreState) => state.member?.name || "",
      getMetrics: (state: WorkspaceStoreState) => state.metrics,
      getThreadItem: (state: WorkspaceStoreState, threadId: string) =>
        state.threads?.[threadId] || null,
      getWorkspaceId: (state: WorkspaceStoreState) =>
        state.workspace?.workspaceId || "",
      getWorkspaceName: (state: WorkspaceStoreState) =>
        state.workspace?.name || "",
      setThreadSortKey: (sortKey: SortBy) => {
        set((state) => {
          state.threadSortKey = sortKey;
          const key = `zyg:${state.workspace?.workspaceId}:sortKey`;
          setTimeout(() => {
            setInLocalStorage(key, sortKey);
          }, 0);
          return state;
        });
      },
      updateLabel: (labelId: string, label: Label) => {
        set((state) => {
          if (state.labels) {
            state.labels[labelId] = { ...label };
            return state;
          } else {
            return state;
          }
        });
      },
      updateThread: (thread) => {
        set((state) => {
          if (state.threads) {
            if (state.threads?.[thread.threadId]) {
              // existing record
              const { ...rest } = state.threads[thread.threadId];
              // updates
              const { ...updates } = thread;
              state.threads[thread.threadId] = {
                ...rest,
                ...updates,
              };
            }
          }
          return state;
        });
      },
      updateWorkspaceName: (name: string) => {
        set((state) => {
          if (state.workspace) {
            state.workspace.name = name;
            return state;
          } else {
            return state;
          }
        });
      },
      viewAssignees: (state: WorkspaceStoreState) => {
        // Get all threads
        const threads = state.threads ? Object.values(state.threads) : [];

        // Extract unique, valid assignee IDs
        const assigneeIds = _.uniq(
          threads
            .map((t) => t.assigneeId)
            .filter((a): a is string => a !== undefined),
        );

        // Map assignee IDs to members
        return assigneeIds
          .map((a) => {
            const member = state.members?.[a];
            if (member) {
              return {
                assigneeId: member.memberId,
                name: member.name || "n/a",
              } as Assignee;
            }
          })
          .filter((m): m is Assignee => m !== undefined);
      },
      viewCurrentThreadQueue: (state: WorkspaceStoreState): null | Thread[] => {
        if (state.threadAppliedFilters) {
          const {
            assignees,
            isUnassigned,
            memberId,
            priorities,
            sortBy,
            stages,
            status,
          } = state.threadAppliedFilters;

          const threads = state.threads ? Object.values(state.threads) : [];

          let results = [];
          results = threads.filter((t) => t.status === status);

          if (memberId) {
            results = threads.filter(
              (t) => t.status === status && t.assigneeId === memberId,
            );
          } else if (isUnassigned) {
            results = threads.filter(
              (t) => t.status === status && !t.assigneeId,
            );
          }

          results = filterByAssignees(results, assignees);
          results = filterByStages(results, stages);
          results = filterByPriorities(results, priorities);
          results = sortThreads(results, sortBy);
          return results;
        }
        return null;
      },
      viewCustomerEmail: (state: WorkspaceStoreState, customerId: string) => {
        const customer = state.customers?.[customerId];
        return customer ? customer.email : null;
      },
      viewCustomerExternalId: (
        state: WorkspaceStoreState,
        customerId: string,
      ) => {
        const customer = state.customers?.[customerId];
        return customer ? customer.externalId : null;
      },
      viewCustomerName: (state: WorkspaceStoreState, customerId: string) => {
        const customer = state.customers?.[customerId];
        return customer ? customer.name : "";
      },
      viewCustomerPhone: (state: WorkspaceStoreState, customerId: string) => {
        const customer = state.customers?.[customerId];
        return customer ? customer.phone : null;
      },
      viewCustomerRole: (state: WorkspaceStoreState, customerId: string) => {
        const customer = state.customers?.[customerId];
        return customer ? customer.role : "";
      },
      viewLabels: (state: WorkspaceStoreState) => {
        const labels = state.labels ? Object.values(state.labels) : [];
        return _.sortBy(labels, "labelId").reverse();
      },
      viewMemberName: (state: WorkspaceStoreState, memberId: string) => {
        const member = state.members?.[memberId];
        return member ? member.name || "" : "";
      },
      viewMembers: (state: WorkspaceStoreState) => {
        return state.members ? Object.values(state.members) : [];
      },
      viewMyThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        memberId: string,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: SortBy = defaultSortKey,
      ) => {
        const threads = state.threads ? Object.values(state.threads) : [];
        const myThreads = threads.filter(
          (t) => t.status === status && t.assigneeId === memberId,
        );
        const assigneesFiltered = filterByAssignees(myThreads, assignees);
        const stagesFiltered = filterByStages(assigneesFiltered, stages);
        const prioritiesFiltered = filterByPriorities(
          stagesFiltered,
          priorities,
        );
        return sortThreads(prioritiesFiltered, sortBy);
      },
      viewPats: (state: WorkspaceStoreState) => {
        const pats = state.pats ? Object.values(state.pats) : [];
        return _.sortBy(pats, "patId").reverse();
      },
      viewThreadAssigneeId: (state: WorkspaceStoreState, threadId: string) =>
        state.threads?.[threadId]?.assigneeId || null,
      viewThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: SortBy = defaultSortKey,
      ) => {
        const threads = state.threads ? Object.values(state.threads) : [];
        const statusFiltered = threads.filter((t) => t.status === status);
        const assigneesFiltered = filterByAssignees(statusFiltered, assignees);
        const stagesFiltered = filterByStages(assigneesFiltered, stages);
        const prioritiesFiltered = filterByPriorities(
          stagesFiltered,
          priorities,
        );
        return sortThreads(prioritiesFiltered, sortBy);
      },
      viewThreadSortKey: (state: WorkspaceStoreState) => {
        if (state.threadSortKey) {
          return state.threadSortKey;
        }
        // if not set in store, check from local storage.
        const key = `zyg:${state.workspace?.workspaceId}:sortKey`;
        const sortKey = getFromLocalStorage(key);
        if (sortKey) {
          return sortKey as SortBy;
        }
        // use default otherwise.
        return defaultSortKey;
      },
      viewUnassignedThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy = defaultSortKey,
      ) => {
        const threads = state.threads ? Object.values(state.threads) : [];
        const unassignedThreads = threads.filter(
          (t) => t.status === status && !t.assigneeId,
        );
        const assigneesFiltered = filterByAssignees(
          unassignedThreads,
          assignees,
        );
        const stagesFiltered = filterByStages(assigneesFiltered, stages);
        const prioritiesFiltered = filterByPriorities(
          stagesFiltered,
          priorities,
        );
        return sortThreads(prioritiesFiltered, sortBy);
      },
    })),
  );
};

export type AccountStoreStateType = IAccount & IAccountStoreActions;

export interface IAccount {
  account: Account | null;
  error: Error | null;
  hasData: boolean;
}

interface IAccountStoreActions {
  getAccount(state: AccountStoreStateType): Account | null;

  getAccountId(state: AccountStoreStateType): string;

  getEmail(state: AccountStoreStateType): string;

  getName(state: AccountStoreStateType): string;

  updateStore(): void;
}

export const buildAccountStore = (initialState: IAccount) => {
  return createStore<AccountStoreStateType>()((set) => ({
    ...initialState,
    getAccount: (state: AccountStoreStateType) => state.account,
    getAccountId: (state: AccountStoreStateType) =>
      state.account?.accountId || "",
    getEmail: (state: AccountStoreStateType) => state.account?.email || "",
    getName: (state: AccountStoreStateType) => state.account?.name || "",
    updateStore: () => set((state) => ({ ...state })),
  }));
};
