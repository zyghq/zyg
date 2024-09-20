import { createStore } from "zustand/vanilla";
import { immer } from "zustand/middleware/immer";
import _ from "lodash";

import {
  Account,
  Workspace,
  AuthMember,
  Pat,
  WorkspaceMetrics,
  Thread,
  Label,
  Member,
  Customer,
} from "./models";

// add more entities as supported by store
// e.g: Workspace | User | etc.
type WorkspaceEntities = Workspace | Thread | Customer | Label | Member | Pat;

export type Dictionary<
  K extends string | number,
  V extends WorkspaceEntities,
> = {
  [key in K]: V;
};

export type ThreadMap = Dictionary<string, Thread>;

export type CustomerMap = Dictionary<string, Customer>;

export type LabelMap = Dictionary<string, Label>;

export type MemberMap = Dictionary<string, Member>;

export type PatMap = Dictionary<string, Pat>;

export interface IWorkspaceEntities {
  workspace: Workspace | null;
  member: AuthMember | null;
  metrics: WorkspaceMetrics;
  threads: ThreadMap | null;
  customers: CustomerMap | null;
  labels: LabelMap | null;
  members: MemberMap | null;
  pats: PatMap | null;
}

type StatusType = "todo" | "done";
type Priority = "urgent" | "high" | "normal" | "low";

type StageType =
  | "needs_first_response"
  | "waiting_on_customer"
  | "hold"
  | "needs_next_response"
  | "resolved"
  | "spam";

export type StagesFiltersType = StageType | StageType[] | undefined;
export type AssigneesFiltersType = string | string[] | undefined;
export type PrioritiesFiltersType = Priority | Priority[] | undefined;

export type ThreadAppliedFilters = {
  status: StatusType;
  assignees: AssigneesFiltersType;
  stages: StagesFiltersType;
  priorities: PrioritiesFiltersType;
  sortBy: sortByType;
  memberId?: string | null;
  isUnassigned?: boolean | null;
};

export type sortByType = "last-message-dsc" | "created-asc" | "created-dsc";
export const defaultSortKey = "last-message-dsc";

export type Assignee = {
  assigneeId: string;
  name: string;
};

interface IWorkspaceStoreActions {
  getWorkspaceName(state: WorkspaceStoreState): string;
  getWorkspaceId(state: WorkspaceStoreState): string;
  getMemberId(state: WorkspaceStoreState): string;
  getMemberName(state: WorkspaceStoreState): string;
  getMetrics(state: WorkspaceStoreState): WorkspaceMetrics;
  getThreadItem(state: WorkspaceStoreState, threadId: string): Thread | null;
  applyThreadFilters(
    status: StatusType,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: sortByType,
    memberId: string | null,
    isUnassigned: boolean | null
  ): void;
  viewThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: sortByType
  ): Thread[];
  viewMyThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    memberId: string,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: sortByType
  ): Thread[];
  viewUnassignedThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: sortByType
  ): Thread[];
  viewCurrentThreadQueue(state: WorkspaceStoreState): Thread[] | null;
  viewCustomerName(state: WorkspaceStoreState, customerId: string): string;
  viewAssignees(state: WorkspaceStoreState): Assignee[];
  updateWorkspaceName(name: string): void;
  viewLabels(state: WorkspaceStoreState): Label[];
  viewMembers(state: WorkspaceStoreState): Member[];
  viewPats(state: WorkspaceStoreState): Pat[];
  addPat(pat: Pat): void;
  deletePat(patId: string): void;
  addLabel(label: Label): void;
  updateLabel(labelId: string, label: Label): void;
  viewMemberName(state: WorkspaceStoreState, memberId: string): string;
  updateThread(thread: Thread): void;
  viewThreadAssigneeId(
    state: WorkspaceStoreState,
    threadId: string
  ): string | null;
  setFromPath(path: string): void;
  viewFromPath(state: WorkspaceStoreState): string | null;
}

export interface IWorkspaceValueObjects {
  hasData: boolean;
  isPending: boolean;
  error: Error | null;
  threadAppliedFilters: ThreadAppliedFilters | null;
  fromPath?: string;
}

export type WorkspaceStoreState = IWorkspaceEntities &
  IWorkspaceValueObjects &
  IWorkspaceStoreActions;

function filterByStages(threads: Thread[], stages: StagesFiltersType) {
  if (stages && Array.isArray(stages)) {
    const uniqueReasons = [...new Set(stages)];
    const filtered = [];
    for (const stage of uniqueReasons) {
      if (stage === "needs_first_response") {
        filtered.push(...threads.filter((t) => t.stage));
      } else if (stage === "needs_next_response") {
        filtered.push(...threads.filter((t) => !t.stage));
      } else if (stage === "waiting_on_customer") {
        filtered.push(...threads.filter((t) => !t.stage));
      } else if (stage === "hold") {
        filtered.push(...threads.filter((t) => !t.stage));
      } else if (stage === "resolved") {
        filtered.push(...threads.filter((t) => !t.stage));
      }
    }
    return filtered;
  }
  if (stages && typeof stages === "string") {
    if (stages === "needs_first_response") {
      return threads.filter((t) => t.stage);
    }
    if (stages === "needs_next_response") {
      return threads.filter((t) => !t.stage);
    }
    if (stages === "waiting_on_customer") {
      return threads.filter((t) => !t.stage);
    }
    if (stages === "hold") {
      return threads.filter((t) => !t.stage);
    }
    if (stages === "resolved") {
      return threads.filter((t) => !t.stage);
    }
    if (stages === "spam") {
      return threads.filter((t) => !t.stage);
    }
  }
  // no change
  return threads;
}

function filterByPriorities(
  threads: Thread[],
  priorities: PrioritiesFiltersType
) {
  if (priorities && Array.isArray(priorities)) {
    const uniquePriorities = [...new Set(priorities)];
    const filtered = [];
    for (const priority of uniquePriorities) {
      filtered.push(...threads.filter((t) => t.priority === priority));
    }
    return filtered;
  }

  if (priorities && typeof priorities === "string") {
    return threads.filter((t) => t.priority === priorities);
  }
  // no change
  return threads;
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
  if (assignees && typeof assignees === "string") {
    return threads.filter((t) => t.assigneeId === assignees);
  }
  // no change
  return threads;
}

function sortThreads(threads: Thread[], sortBy: sortByType) {
  if (sortBy === "created-dsc") {
    return _.sortBy(threads, "createdAt").reverse();
  } else if (sortBy === "created-asc") {
    return _.sortBy(threads, "createdAt");
  }
  // default sorted by last-message-dsc (from server)
  return threads;
}

// (sanchitrk) for reference on using zustand, check this great article:
// https://tkdodo.eu/blog/working-with-zustand
//
// @sanchitrk: shall we rename it to `buildWorkspaceStore` ?
export const buildStore = (
  initialState: IWorkspaceEntities & IWorkspaceValueObjects
) => {
  return createStore<WorkspaceStoreState>()(
    immer((set) => ({
      ...initialState,
      getWorkspaceName: (state: WorkspaceStoreState) =>
        state.workspace?.name || "",
      getWorkspaceId: (state: WorkspaceStoreState) =>
        state.workspace?.workspaceId || "",
      getMemberId: (state: WorkspaceStoreState) => state.member?.memberId || "",
      getMemberName: (state: WorkspaceStoreState) => state.member?.name || "",
      getMetrics: (state: WorkspaceStoreState) => state.metrics,
      getThreadItem: (state: WorkspaceStoreState, threadId: string) =>
        state.threads?.[threadId] || null,
      viewThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: sortByType = defaultSortKey
      ) => {
        const threads = state.threads ? Object.values(state.threads) : [];
        const statusFiltered = threads.filter((t) => t.status === status);
        const assigneesFiltered = filterByAssignees(statusFiltered, assignees);
        const stagesFiltered = filterByStages(assigneesFiltered, stages);
        const prioritiesFiltered = filterByPriorities(
          stagesFiltered,
          priorities
        );
        const sortedThreads = sortThreads(prioritiesFiltered, sortBy);
        return sortedThreads;
      },
      viewMyThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        memberId: string,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: sortByType = defaultSortKey
      ) => {
        const threads = state.threads ? Object.values(state.threads) : [];
        const myThreads = threads.filter(
          (t) => t.status === status && t.assigneeId === memberId
        );
        const assigneesFiltered = filterByAssignees(myThreads, assignees);
        const stagesFiltered = filterByStages(assigneesFiltered, stages);
        const prioritiesFiltered = filterByPriorities(
          stagesFiltered,
          priorities
        );
        const sortedThreads = sortThreads(prioritiesFiltered, sortBy);
        return sortedThreads;
      },
      viewUnassignedThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy = defaultSortKey
      ) => {
        const threads = state.threads ? Object.values(state.threads) : [];
        const unassignedThreads = threads.filter(
          (t) => t.status === status && !t.assigneeId
        );
        const assigneesFiltered = filterByAssignees(
          unassignedThreads,
          assignees
        );
        const stagesFiltered = filterByStages(assigneesFiltered, stages);
        const prioritiesFiltered = filterByPriorities(
          stagesFiltered,
          priorities
        );
        const sortedThreads = sortThreads(prioritiesFiltered, sortBy);
        return sortedThreads;
      },
      applyThreadFilters: (
        status: StatusType,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: sortByType,
        memberId: string | null,
        isUnassigned: boolean | null
      ) => {
        set((state) => {
          if (state.threadAppliedFilters) {
            state.threadAppliedFilters = {
              ...state.threadAppliedFilters,
              status,
              assignees,
              stages,
              priorities,
              sortBy,
              memberId,
              isUnassigned,
            };
          } else {
            state.threadAppliedFilters = {
              status,
              assignees,
              stages,
              priorities,
              sortBy,
              memberId,
              isUnassigned,
            };
          }
        });
      },
      viewCurrentThreadQueue: (state: WorkspaceStoreState): Thread[] | null => {
        if (state.threadAppliedFilters) {
          const {
            status,
            assignees,
            stages,
            priorities,
            sortBy,
            memberId,
            isUnassigned,
          } = state.threadAppliedFilters;

          const threads = state.threads ? Object.values(state.threads) : [];

          let results = [];
          results = threads.filter((t) => t.status === status);

          if (memberId) {
            results = threads.filter(
              (t) => t.status === status && t.assigneeId === memberId
            );
          } else if (isUnassigned) {
            results = threads.filter(
              (t) => t.status === status && !t.assigneeId
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
      viewCustomerName: (state: WorkspaceStoreState, customerId: string) => {
        const customer = state.customers?.[customerId];
        return customer ? customer.name : "";
      },
      viewAssignees: (state: WorkspaceStoreState) => {
        // Get all threads
        const threads = state.threads ? Object.values(state.threads) : [];

        // Extract unique, valid assignee IDs
        const assigneeIds = _.uniq(
          threads
            .map((t) => t.assigneeId)
            .filter((a): a is string => a !== undefined)
        );

        // Map assignee IDs to members
        const assignees = assigneeIds
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

        return assignees;
      },
      viewLabels: (state: WorkspaceStoreState) => {
        const labels = state.labels ? Object.values(state.labels) : [];
        return _.sortBy(labels, "labelId").reverse();
      },
      viewMembers: (state: WorkspaceStoreState) => {
        const members = state.members ? Object.values(state.members) : [];
        return members;
      },
      viewPats: (state: WorkspaceStoreState) => {
        const pats = state.pats ? Object.values(state.pats) : [];
        return _.sortBy(pats, "patId").reverse();
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
      viewMemberName: (state: WorkspaceStoreState, memberId: string) => {
        const member = state.members?.[memberId];
        return member ? member.name || "" : "";
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
      setFromPath: (path: string) => {
        set((state) => {
          if (state.fromPath) {
            state.fromPath = path;
            return state;
          } else {
            state.fromPath = path;
            return state;
          }
        });
      },
      viewFromPath: (state: WorkspaceStoreState) => state.fromPath || null,
      viewThreadAssigneeId: (state: WorkspaceStoreState, threadId: string) =>
        state.threads?.[threadId]?.assigneeId || null,
    }))
  );
};

export interface IAccount {
  hasData: boolean;
  error: Error | null;
  account: Account | null;
}

interface IAccountStoreActions {
  updateStore(): void;
  getAccount(state: AccountStoreStateType): Account | null;
  getName(state: AccountStoreStateType): string;
  getAccountId(state: AccountStoreStateType): string;
  getEmail(state: AccountStoreStateType): string;
}

export type AccountStoreStateType = IAccount & IAccountStoreActions;

export const buildAccountStore = (initialState: IAccount) => {
  return createStore<AccountStoreStateType>()((set) => ({
    ...initialState,
    getAccount: (state: AccountStoreStateType) => state.account,
    updateStore: () => set((state) => ({ ...state })),
    getName: (state: AccountStoreStateType) => state.account?.name || "",
    getAccountId: (state: AccountStoreStateType) =>
      state.account?.accountId || "",
    getEmail: (state: AccountStoreStateType) => state.account?.email || "",
  }));
};
