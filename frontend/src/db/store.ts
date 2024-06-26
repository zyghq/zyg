import { createStore } from "zustand/vanilla";
import { immer } from "zustand/middleware/immer";
import _ from "lodash";

import {
  Account,
  Workspace,
  UserMember,
  AccountPat,
  WorkspaceMetrics,
  ThreadChatWithRecentMessage,
  ThreadChat,
  Label,
  Member,
  Customer,
} from "./entities";

// add more entitites as supported by store
// e.g: Workspace | User | etc.
type AllowedEntities =
  | Workspace
  | ThreadChatWithRecentMessage
  | Customer
  | Label
  | Member
  | AccountPat;

export type Dictionary<K extends string | number, V extends AllowedEntities> = {
  [key in K]: V;
};

export type ThreadChatMap = Dictionary<string, ThreadChatWithRecentMessage>;

export type WorkspaceCustomerMap = Dictionary<string, Customer>;

export type WorkspaceLabelMap = Dictionary<string, Label>;

export type WorkspaceMemberMap = Dictionary<string, Member>;

export type AccountPatMap = Dictionary<string, AccountPat>;

// export type CurrentThreadQueueType = {
//   status: string;
//   from: string;
//   threads: ThreadChatWithRecentMessage[];
// };

export interface IWorkspaceEntities {
  workspace: Workspace | null;
  member: UserMember | null;
  metrics: WorkspaceMetrics;
  threadChats: ThreadChatMap | null;
  customers: WorkspaceCustomerMap | null;
  labels: WorkspaceLabelMap | null;
  members: WorkspaceMemberMap | null;
  pats: AccountPatMap | null;
}

type ReplyStatus = "replied" | "unreplied";
type Priority = "urgent" | "high" | "normal" | "low";

export type StatusType = "todo" | "done" | "snoozed";
export type ReasonsFiltersType = ReplyStatus | ReplyStatus[] | undefined;
export type AssigneesFiltersType = string | string[] | undefined;
export type PrioritiesFiltersType = Priority | Priority[] | undefined;

export type ThreadAppliedFilters = {
  status: StatusType;
  assignees: AssigneesFiltersType;
  reasons: ReasonsFiltersType;
  priorities: PrioritiesFiltersType;
  sortBy: sortByType;
  memberId?: string | null;
  isUnassigned?: boolean | null;
};

export type sortByType = "last-message-dsc" | "created-asc" | "created-dsc";

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
  getThreadChatItem(
    state: WorkspaceStoreState,
    threadChatId: string
  ): ThreadChatWithRecentMessage | null;
  applyThreadFilters(
    status: StatusType,
    assignees: AssigneesFiltersType,
    reasons: ReasonsFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: sortByType,
    memberId: string | null,
    isUnassigned: boolean | null
  ): void;
  viewThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    assignees: AssigneesFiltersType,
    reasons: ReasonsFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: sortByType
  ): ThreadChatWithRecentMessage[];
  viewMyThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    memberId: string,
    assignees: AssigneesFiltersType,
    reasons: ReasonsFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: sortByType
  ): ThreadChatWithRecentMessage[];
  viewUnassignedThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    assignees: AssigneesFiltersType,
    reasons: ReasonsFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: sortByType
  ): ThreadChatWithRecentMessage[];
  viewCurrentThreadQueue(
    state: WorkspaceStoreState
  ): ThreadChatWithRecentMessage[] | null;
  viewCustomerName(state: WorkspaceStoreState, customerId: string): string;
  viewAssignees(state: WorkspaceStoreState): Assignee[];
  updateWorkspaceName(name: string): void;
  viewLabels(state: WorkspaceStoreState): Label[];
  viewMembers(state: WorkspaceStoreState): Member[];
  viewPats(state: WorkspaceStoreState): AccountPat[];
  addPat(pat: AccountPat): void;
  deletePat(patId: string): void;
  addLabel(label: Label): void;
  updateLabel(labelId: string, label: Label): void;
  viewMemberName(state: WorkspaceStoreState, memberId: string): string;
  updateMainThreadChat(threadChat: ThreadChat): void;
  viewThreadChatAssigneeId(
    state: WorkspaceStoreState,
    threadChatId: string
  ): string | null;
}

export interface IWorkspaceValueObjects {
  hasData: boolean;
  isPending: boolean;
  error: Error | null;
  threadAppliedFilters: ThreadAppliedFilters | null;
}

export type WorkspaceStoreState = IWorkspaceEntities &
  IWorkspaceValueObjects &
  IWorkspaceStoreActions;

function filterByReasons(
  threads: ThreadChatWithRecentMessage[],
  reasons: ReasonsFiltersType
) {
  if (reasons && Array.isArray(reasons)) {
    const uniqueReasons = [...new Set(reasons)];
    const filtered = [];
    for (const reason of uniqueReasons) {
      if (reason === "replied") {
        filtered.push(...threads.filter((t) => t.replied));
      } else if (reason === "unreplied") {
        filtered.push(...threads.filter((t) => !t.replied));
      }
    }
    return filtered;
  }
  if (reasons && typeof reasons === "string") {
    if (reasons === "replied") {
      return threads.filter((t) => t.replied);
    }
    if (reasons === "unreplied") {
      return threads.filter((t) => !t.replied);
    }
  }
  // no change
  return threads;
}

function filterByPriorities(
  threads: ThreadChatWithRecentMessage[],
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

function filterByAssignees(
  threads: ThreadChatWithRecentMessage[],
  assignees: AssigneesFiltersType
) {
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

function sortThreads(
  threads: ThreadChatWithRecentMessage[],
  sortBy: sortByType
) {
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
      getThreadChatItem: (state: WorkspaceStoreState, threadChatId: string) =>
        state.threadChats?.[threadChatId] || null,
      viewThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        assignees: AssigneesFiltersType,
        reasons: ReasonsFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: sortByType = "last-message-dsc"
      ) => {
        const threads = state.threadChats
          ? Object.values(state.threadChats)
          : [];
        const statusFiltered = threads.filter((t) => t.status === status);
        const assigneesFiltered = filterByAssignees(statusFiltered, assignees);
        const reasonsFiltered = filterByReasons(assigneesFiltered, reasons);
        const prioritiesFiltered = filterByPriorities(
          reasonsFiltered,
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
        reasons: ReasonsFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: sortByType = "last-message-dsc"
      ) => {
        const threads = state.threadChats
          ? Object.values(state.threadChats)
          : [];
        const myThreads = threads.filter(
          (t) => t.status === status && t.assigneeId === memberId
        );
        const assigneesFiltered = filterByAssignees(myThreads, assignees);
        const reasonsFiltered = filterByReasons(assigneesFiltered, reasons);
        const prioritiesFiltered = filterByPriorities(
          reasonsFiltered,
          priorities
        );
        const sortedThreads = sortThreads(prioritiesFiltered, sortBy);
        // state.currentThreadQueue = {
        //   status: "todo",
        //   from: "My Threads",
        //   threads: [...sortedThreads],
        // };
        return sortedThreads;
      },
      viewUnassignedThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        assignees: AssigneesFiltersType,
        reasons: ReasonsFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy = "last-message-dsc"
      ) => {
        const threads = state.threadChats
          ? Object.values(state.threadChats)
          : [];
        const unassignedThreads = threads.filter(
          (t) => t.status === status && !t.assigneeId
        );
        const assigneesFiltered = filterByAssignees(
          unassignedThreads,
          assignees
        );
        const reasonsFiltered = filterByReasons(assigneesFiltered, reasons);
        const prioritiesFiltered = filterByPriorities(
          reasonsFiltered,
          priorities
        );
        const sortedThreads = sortThreads(prioritiesFiltered, sortBy);
        return sortedThreads;
      },
      applyThreadFilters: (
        status: StatusType,
        assignees: AssigneesFiltersType,
        reasons: ReasonsFiltersType,
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
              reasons,
              priorities,
              sortBy,
              memberId,
              isUnassigned,
            };
          } else {
            state.threadAppliedFilters = {
              status,
              assignees,
              reasons,
              priorities,
              sortBy,
              memberId,
              isUnassigned,
            };
          }
        });
      },
      viewCurrentThreadQueue: (
        state: WorkspaceStoreState
      ): ThreadChatWithRecentMessage[] | null => {
        if (state.threadAppliedFilters) {
          const {
            status,
            assignees,
            reasons,
            priorities,
            sortBy,
            memberId,
            isUnassigned,
          } = state.threadAppliedFilters;

          const threads = state.threadChats
            ? Object.values(state.threadChats)
            : [];

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
          results = filterByReasons(results, reasons);
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
        const threads = state.threadChats
          ? Object.values(state.threadChats)
          : [];

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
      addPat: (pat: AccountPat) => {
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
      updateMainThreadChat: (threadChat) => {
        set((state) => {
          if (state.threadChats) {
            if (state.threadChats?.[threadChat.threadChatId]) {
              // existing record
              const { recentMessage, ...rest } =
                state.threadChats[threadChat.threadChatId];
              // updates
              const { assignee, customer, ...updates } = threadChat;
              state.threadChats[threadChat.threadChatId] = {
                recentMessage,
                ...rest,
                ...updates,
                customerId: customer.customerId,
                assigneeId: assignee?.memberId || null,
              };
            }
          }
          return state;
        });
      },
      viewThreadChatAssigneeId: (
        state: WorkspaceStoreState,
        threadChatId: string
      ) => state.threadChats?.[threadChatId]?.assigneeId || null,
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
