import {
  HOLD,
  NEEDS_FIRST_RESPONSE,
  NEEDS_NEXT_RESPONSE,
  RESOLVED,
  SPAM,
  WAITING_ON_CUSTOMER,
} from "@/db/constants";
import { getFromLocalStorage, setInLocalStorage } from "@/db/helpers";
import {
  Account,
  AuthMember,
  Label,
  Pat,
  Workspace,
  WorkspaceMetrics,
} from "@/db/models";
import {
  CustomerShape,
  CustomerShapeUpdates,
  MemberShape,
  MemberShapeUpdates,
  ThreadLabelShape,
  ThreadShape,
  ThreadShapeUpdates,
} from "@/db/shapes";
import { enableMapSet } from "immer";
import _ from "lodash";
import { immer } from "zustand/middleware/immer";
import { createStore } from "zustand/vanilla";

enableMapSet();

// Represents entities for KV dictionary values.
// add more entities as supported by KV store
// e.g: Workspace | User | etc.
type EntitiesKV =
  | CustomerShape
  | Label
  | MemberShape
  | Pat
  | ThreadShape
  | Workspace;

export type Dictionary<K extends number | string, V extends EntitiesKV> = {
  [key in K]: V;
};

export type LabelMap = Dictionary<string, Label>;

export type MemberShapeMap = Dictionary<string, MemberShape>;

export type PatMap = Dictionary<string, Pat>;

export type CustomerShapeMap = Dictionary<string, CustomerShape>;

export type ThreadShapeMap = Map<string, ThreadShape>;

// Represents the store entities.
export interface IWorkspaceEntities {
  customers: CustomerShapeMap | null;
  labels: LabelMap | null;
  member: AuthMember;
  members: MemberShapeMap | null;
  metrics: WorkspaceMetrics;
  pats: null | PatMap;
  threads: null | ThreadShapeMap;
  workspace: null | Workspace;
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

export type Priority = "high" | "low" | "normal" | "urgent";
export type StageType =
  | "hold"
  | "needs_first_response"
  | "needs_next_response"
  | "resolved"
  | "spam"
  | "waiting_on_customer";

export type StatusType = "done" | "todo";

export const defaultSortKey = "created-asc";

export type Assignee = {
  assigneeId: string;
  name: string;
};
export type AssigneesFiltersType = string | string[] | undefined;

// Represents workspace value objects.
// Unlike entities which represent backend system data, value objects are application level.
export interface IWorkspaceValueObjects {
  customersShapeHandle: null | string;
  customersShapeOffset: string;
  error: Error | null;
  inSync: boolean;
  membersShapeHandle: null | string;
  membersShapeOffset: string;
  threadAppliedFilters: null | ThreadAppliedFilters;
  threadSortKey: null | SortBy;
  threadsShapeHandle: null | string;
  threadsShapeOffset: string;
}

export type PrioritiesFiltersType = Priority | Priority[] | undefined;

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

// Represents the full store state with entities and application state.
export type WorkspaceStoreState = IWorkspaceEntities &
  IWorkspaceStoreActions &
  IWorkspaceValueObjects;

interface IWorkspaceStoreActions {
  // Currently being used in settings page
  // Adds label for the workspace.
  addLabel(label: Label): void;

  addPat(pat: Pat): void;

  addThreadLabel(threadId: string, label: ThreadLabelShape): void;
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

  getThreadItem(
    state: WorkspaceStoreState,
    threadId: string,
  ): null | ThreadShape;

  getWorkspaceId(state: WorkspaceStoreState): string;

  getWorkspaceName(state: WorkspaceStoreState): string;

  isInSync(state: WorkspaceStoreState): boolean;

  removeThreadLabel(threadId: string, labelId: string): void;

  setCustomersShapeHandle(handle: null | string): void;

  setCustomersShapeOffset(offset: string): void;

  setInSync(f: boolean): void;

  setMembersShapeHandle(handle: null | string): void;

  setMembersShapeOffset(offset: string): void;

  setThreadSortKey(sortKey: SortBy): void;

  setThreadsShapeHandle(handle: null | string): void;

  setThreadsShapeOffset(offset: string): void;

  updateCustomer(member: CustomerShapeUpdates): void;

  // Currently being used in settings page, adding workspace labels
  updateLabel(labelId: string, label: Label): void;

  updateMember(member: MemberShapeUpdates): void;

  updateThread(thread: ThreadShapeUpdates): void;

  updateThreadAssignee(threadId: string, memberId: null | string): void;

  updateThreadPriority(threadId: string, priority: Priority): void;

  updateThreadStage(threadId: string, stage: StageType): void;

  updateWorkspaceName(name: string): void;

  viewAssignees(state: WorkspaceStoreState): Assignee[];

  viewCurrentThreadQueue(state: WorkspaceStoreState): null | ThreadShape[];

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

  viewCustomersShapeHandle(state: WorkspaceStoreState): null | string;

  viewCustomersShapeOffset(state: WorkspaceStoreState): string;

  // List workspace labels
  viewLabels(state: WorkspaceStoreState): Label[];

  viewMemberName(state: WorkspaceStoreState, memberId: string): string;

  viewMembers(state: WorkspaceStoreState): MemberShape[];

  viewMembersShapeHandle(state: WorkspaceStoreState): null | string;

  viewMembersShapeOffset(state: WorkspaceStoreState): string;

  viewMyThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    memberId: string,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: SortBy,
  ): ThreadShape[];

  viewPats(state: WorkspaceStoreState): Pat[];

  viewThreadAssigneeId(
    state: WorkspaceStoreState,
    threadId: string,
  ): null | string;

  viewThreadLabels(
    state: WorkspaceStoreState,
    threadId: string,
  ): ThreadLabelShape[];

  viewThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: SortBy,
  ): ThreadShape[];

  viewThreadSortKey(state: WorkspaceStoreState): SortBy;

  viewThreadsShapeHandle(state: WorkspaceStoreState): null | string;

  viewThreadsShapeOffset(state: WorkspaceStoreState): string;

  viewUnassignedThreads(
    state: WorkspaceStoreState,
    status: StatusType,
    assignees: AssigneesFiltersType,
    stages: StagesFiltersType,
    priorities: PrioritiesFiltersType,
    sortBy: SortBy,
  ): ThreadShape[];
}

// function filterByAssignees(
//   threads: ThreadShape[],
//   assignees: AssigneesFiltersType,
// ) {
//   if (assignees && Array.isArray(assignees)) {
//     const uniqueAssignees = [...new Set(assignees)];
//     const filtered = [];
//     for (const assignee of uniqueAssignees) {
//       filtered.push(...threads.filter((t) => t.assigneeId === assignee));
//     }
//     return filtered;
//   }
//   if (assignees) {
//     return threads.filter((t) => t.assigneeId === assignees);
//   }
//   // no change
//   return threads;
// }

function filterByAssignees(
  threads: ThreadShape[],
  assignees: AssigneesFiltersType,
): ThreadShape[] {
  if (!assignees) return threads;

  // Convert assignees to Set for O(1) lookup
  const assigneeSet = Array.isArray(assignees)
    ? new Set(assignees)
    : new Set([assignees]);

  return threads.filter((t) => t.assigneeId && assigneeSet.has(t.assigneeId));
}

// function filterByPriorities(
//   threads: ThreadShape[],
//   priorities: PrioritiesFiltersType,
// ) {
//   if (priorities && Array.isArray(priorities)) {
//     const uniquePriorities = [...new Set(priorities)];
//     const filtered = [];
//     for (const priority of uniquePriorities) {
//       filtered.push(...threads.filter((t) => t.priority === priority));
//     }
//     return filtered;
//   }
//
//   if (priorities) {
//     return threads.filter((t) => t.priority === priorities);
//   }
//   // no change
//   return threads;
// }

function filterByPriorities(
  threads: ThreadShape[],
  priorities: PrioritiesFiltersType,
): ThreadShape[] {
  if (!priorities) return threads;

  // Convert priorities to Set for O(1) lookup
  const prioritySet = Array.isArray(priorities)
    ? new Set(priorities)
    : new Set([priorities]);

  return threads.filter((t) => prioritySet.has(t.priority as Priority));
}

// function filterByStages(threads: ThreadShape[], stages: StagesFiltersType) {
//   const stageMap: Record<string, string> = {
//     hold: HOLD,
//     needs_first_response: NEEDS_FIRST_RESPONSE,
//     needs_next_response: NEEDS_NEXT_RESPONSE,
//     resolved: RESOLVED,
//     spam: SPAM,
//     waiting_on_customer: WAITING_ON_CUSTOMER,
//   };
//
//   if (stages) {
//     if (Array.isArray(stages)) {
//       const uniqueStages = [...new Set(stages)];
//       return threads.filter((t) =>
//         uniqueStages.some((stage) => t.stage === stageMap[stage]),
//       );
//     }
//     if (stageMap[stages]) {
//       return threads.filter((t) => t.stage === stageMap[stages]);
//     }
//   }
//
//   // no change
//   return threads;
// }

function filterByStages(
  threads: ThreadShape[],
  stages: StagesFiltersType,
): ThreadShape[] {
  if (!stages) return threads;

  const stageMap: Record<string, string> = {
    hold: HOLD,
    needs_first_response: NEEDS_FIRST_RESPONSE,
    needs_next_response: NEEDS_NEXT_RESPONSE,
    resolved: RESOLVED,
    spam: SPAM,
    waiting_on_customer: WAITING_ON_CUSTOMER,
  };

  // Convert stages to mapped Set for O(1) lookup
  const stageSet = Array.isArray(stages)
    ? new Set(stages.map((stage) => stageMap[stage]))
    : new Set([stageMap[stages]]);

  return threads.filter((t) => stageSet.has(t.stage));
}

const PRIORITY_MAP: Record<string, number> = {
  high: 1,
  low: 3,
  normal: 2,
  urgent: 0,
} as const;

function sortThreads(threads: ThreadShape[], sortBy: SortBy): ThreadShape[] {
  switch (sortBy) {
    case "created-asc":
      return _.orderBy(threads, "createdAt", "asc");
    case "created-dsc":
      return _.orderBy(threads, "createdAt", "desc");
    case "inbound-message-dsc":
      return _.orderBy(threads, "lastInboundAt", "desc");
    case "outbound-message-dsc":
      return _.orderBy(threads, "lastOutboundAt", "desc");
    case "priority-asc":
      return _.orderBy(
        threads,
        (thread) => PRIORITY_MAP[thread.priority] ?? Infinity,
        "asc",
      );
    case "priority-dsc":
      return _.orderBy(
        threads,
        (thread) => PRIORITY_MAP[thread.priority] ?? -Infinity,
        "desc",
      );
    case "status-changed-asc":
      return _.orderBy(threads, "statusChangedAt", "asc");
    case "status-changed-dsc":
      return _.orderBy(threads, "statusChangedAt", "desc");
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
      // Label Management
      addLabel: (label: Label) => {
        const { labelId } = label;
        set((state) => {
          state.labels ??= {};
          state.labels[labelId] = { ...label };
        });
      },

      updateLabel: (labelId: string, label: Label) => {
        set((state) => {
          if (!state.labels) return;
          state.labels[labelId] = { ...label };
        });
      },
      viewMemberName: (state, memberId: string) =>
        state.members?.[memberId]?.name ?? "",
      viewThreadAssigneeId: (state, threadId: string) =>
        state.threads?.get(threadId)?.assigneeId ?? null,

      // Pat Management
      addPat: (pat: Pat) => {
        const { patId } = pat;
        set((state) => {
          state.pats ??= {};
          state.pats[patId] = { ...pat };
        });
      },

      deletePat: (patId: string) => {
        set((state) => {
          if (!state.pats) return;
          delete state.pats[patId];
        });
      },

      // Thread Management
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
        });
      },

      updateThread: (thread: ThreadShapeUpdates) => {
        if (!thread?.threadId) return;
        const id: string = thread.threadId;

        set((state) => {
          if (!state.threads?.has(id)) return;

          const existingThread = state.threads.get(id)!;
          state.threads.set(id, {
            ...existingThread,
            ...thread,
          });
        });
      },

      // Thread Property Updates
      updateThreadStage: (threadId: string, stage: StageType) => {
        set((state) => {
          const thread = state.threads?.get(threadId);
          if (thread) thread.stage = stage;
        });
      },

      updateThreadAssignee: (threadId: string, memberId: null | string) => {
        set((state) => {
          const thread = state.threads?.get(threadId);
          if (thread) thread.assigneeId = memberId;
        });
      },

      updateThreadPriority: (threadId: string, priority: Priority) => {
        set((state) => {
          const thread = state.threads?.get(threadId);
          if (thread) thread.priority = priority;
        });
      },

      // adds label to a specific thread.
      addThreadLabel: (threadId: string, label: ThreadLabelShape) => {
        set((state) => {
          const thread = state.threads?.get(threadId);
          if (thread) {
            thread.labels ??= {};
            thread.labels[label.labelId] = label;
          }
        });
      },

      removeThreadLabel: async (threadId: string, labelId: string) => {
        set((state) => {
          const thread = state.threads?.get(threadId);
          if (thread) {
            delete thread.labels?.[labelId];
          }
        });
      },

      // Customer Management
      updateCustomer: (customer: CustomerShapeUpdates) => {
        if (!customer?.customerId) return;
        const id: string = customer.customerId;

        set((state) => {
          if (!state.customers?.[id]) return;

          state.customers[id] = {
            ...state.customers[id],
            ...customer,
          };
        });
      },

      // Member Management
      updateMember: (member: MemberShapeUpdates) => {
        if (!member?.memberId) return;
        const id: string = member.memberId;

        set((state) => {
          if (!state.members?.[id]) return;

          state.members[id] = {
            ...state.members[id],
            ...member,
          };
        });
      },

      // Shape Management
      setCustomersShapeHandle: (handle: null | string) =>
        set((state) => {
          state.customersShapeHandle = handle;
        }),

      setCustomersShapeOffset: (offset: string) =>
        set((state) => {
          state.customersShapeOffset = offset;
        }),

      setMembersShapeHandle: (handle: null | string) =>
        set((state) => {
          state.membersShapeHandle = handle;
        }),

      setMembersShapeOffset: (offset: string) =>
        set((state) => {
          state.membersShapeOffset = offset;
        }),

      setThreadsShapeHandle: (handle: null | string) =>
        set((state) => {
          state.threadsShapeHandle = handle;
        }),

      setThreadsShapeOffset: (offset: string) =>
        set((state) => {
          state.threadsShapeOffset = offset;
        }),

      // Sync Management
      setInSync: (f: boolean) =>
        set((state) => {
          state.inSync = f;
        }),

      // Sort Management
      setThreadSortKey: (sortKey: SortBy) => {
        set((state) => {
          state.threadSortKey = sortKey;
          const key = `zyg:${state.workspace?.workspaceId}:sortKey`;
          setTimeout(() => setInLocalStorage(key, sortKey), 0);
        });
      },

      // Workspace Management
      updateWorkspaceName: (name: string) => {
        set((state) => {
          if (state.workspace) state.workspace.name = name;
        });
      },

      // Getters
      getMemberId: (state) => state.member?.memberId || "",
      getMemberName: (state) => state.member?.name || "",
      getMetrics: (state) => state.metrics,
      getThreadItem: (state, threadId) => state.threads?.get(threadId) ?? null,
      getWorkspaceId: (state) => state.workspace?.workspaceId || "",
      getWorkspaceName: (state) => state.workspace?.name || "",
      isInSync: (state) => state.inSync,

      // View Methods
      viewThreadLabels: (state, threadId) => {
        const thread = state.threads?.get(threadId);
        return thread?.labels ? Object.values(thread.labels) : [];
      },
      viewAssignees: (state: WorkspaceStoreState): Assignee[] => {
        const threads = state.threads ? [...state.threads.values()] : [];
        const assigneeIds = _.uniq(
          threads
            .map((t) => t.assigneeId)
            .filter((id): id is string => Boolean(id)),
        );

        return assigneeIds
          .map((id) => {
            const member = state.members?.[id];
            return member
              ? {
                  assigneeId: member.memberId,
                  name: member.name || "n/a",
                }
              : undefined;
          })
          .filter((assignee): assignee is Assignee => Boolean(assignee));
      },

      viewCurrentThreadQueue: (state) => {
        if (!state.threadAppliedFilters || !state.threads) {
          return state.threadAppliedFilters ? [] : null;
        }

        const {
          status,
          assignees,
          stages,
          priorities,
          sortBy,
          memberId,
          isUnassigned,
        } = state.threadAppliedFilters;

        let results = [...state.threads.values()].filter(
          (t) => t.status === status,
        );

        if (memberId) {
          results = results.filter((t) => t.assigneeId === memberId);
        } else if (isUnassigned) {
          results = results.filter((t) => !t.assigneeId);
        } else if (assignees) {
          results = filterByAssignees(results, assignees);
        }

        if (stages) results = filterByStages(results, stages);
        if (priorities) results = filterByPriorities(results, priorities);

        return sortThreads(results, sortBy);
      },

      viewThreadSortKey: (state) => {
        if (state.threadSortKey) return state.threadSortKey;

        const key = `zyg:${state.workspace?.workspaceId}:sortKey`;
        return (getFromLocalStorage(key) as SortBy) || defaultSortKey;
      },

      // Customer View Methods
      viewCustomerEmail: (state, customerId) =>
        state.customers?.[customerId]?.email ?? null,
      viewCustomerExternalId: (state, customerId) =>
        state.customers?.[customerId]?.externalId ?? null,
      viewCustomerName: (state, customerId) =>
        state.customers?.[customerId]?.name ?? "",
      viewCustomerPhone: (state, customerId) =>
        state.customers?.[customerId]?.phone ?? null,
      viewCustomerRole: (state, customerId) =>
        state.customers?.[customerId]?.role ?? "",

      // Shape View Methods
      viewCustomersShapeHandle: (state) => state.customersShapeHandle,
      viewCustomersShapeOffset: (state) => state.customersShapeOffset,
      viewMembersShapeHandle: (state) => state.membersShapeHandle,
      viewMembersShapeOffset: (state) => state.membersShapeOffset,
      viewThreadsShapeHandle: (state) => state.threadsShapeHandle,
      viewThreadsShapeOffset: (state) => state.threadsShapeOffset,

      // Collection View Methods
      viewLabels: (state) =>
        _.sortBy(Object.values(state.labels ?? {}), "labelId").reverse(),
      viewMembers: (state) => Object.values(state.members ?? {}),
      viewPats: (state) =>
        _.sortBy(Object.values(state.pats ?? {}), "patId").reverse(),

      // Thread View Methods
      viewThreads: (
        state,
        status,
        assignees,
        stages,
        priorities,
        sortBy = defaultSortKey,
      ) => {
        if (!state.threads) return [];

        let results = [...state.threads.values()].filter(
          (t) => t.status === status,
        );

        if (assignees) results = filterByAssignees(results, assignees);
        if (stages) results = filterByStages(results, stages);
        if (priorities) results = filterByPriorities(results, priorities);

        return sortThreads(results, sortBy);
      },

      viewMyThreads: (
        state,
        status,
        memberId,
        assignees,
        stages,
        priorities,
        sortBy = defaultSortKey,
      ) => {
        if (!state.threads) return [];

        let results = [...state.threads.values()].filter(
          (t) => t.status === status && t.assigneeId === memberId,
        );

        if (assignees) results = filterByAssignees(results, assignees);
        if (stages) results = filterByStages(results, stages);
        if (priorities) results = filterByPriorities(results, priorities);

        return sortThreads(results, sortBy);
      },

      viewUnassignedThreads: (
        state,
        status,
        assignees,
        stages,
        priorities,
        sortBy = defaultSortKey,
      ) => {
        if (!state.threads) return [];

        let results = [...state.threads.values()].filter(
          (t) => t.status === status && !t.assigneeId,
        );

        if (assignees) results = filterByAssignees(results, assignees);
        if (stages) results = filterByStages(results, stages);
        if (priorities) results = filterByPriorities(results, priorities);

        return sortThreads(results, sortBy);
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
