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

type Priority = "high" | "low" | "normal" | "urgent";
type StageType =
  | "hold"
  | "needs_first_response"
  | "needs_next_response"
  | "resolved"
  | "spam"
  | "waiting_on_customer";

type StatusType = "done" | "todo";

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

  getThreadItem(
    state: WorkspaceStoreState,
    threadId: string,
  ): null | ThreadShape;

  getWorkspaceId(state: WorkspaceStoreState): string;

  getWorkspaceName(state: WorkspaceStoreState): string;

  isInSync(state: WorkspaceStoreState): boolean;

  setCustomersShapeHandle(handle: null | string): void;

  setCustomersShapeOffset(offset: string): void;

  setInSync(f: boolean): void;

  setMembersShapeHandle(handle: null | string): void;

  setMembersShapeOffset(offset: string): void;

  setThreadSortKey(sortKey: SortBy): void;

  setThreadsShapeHandle(handle: null | string): void;

  setThreadsShapeOffset(offset: string): void;

  updateCustomer(member: CustomerShapeUpdates): void;

  updateLabel(labelId: string, label: Label): void;

  updateMember(member: MemberShapeUpdates): void;

  updateThread(thread: ThreadShapeUpdates): void;

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

// function sortThreads(threads: ThreadShape[], sortBy: SortBy): ThreadShape[] {
//   const priorityMap: { [key: string]: number } = {
//     high: 1,
//     low: 3,
//     normal: 2,
//     urgent: 0,
//   };
//   switch (sortBy) {
//     case "created-asc":
//       return _.sortBy(threads, "createdAt");
//     case "created-dsc":
//       return _.sortBy(threads, "createdAt").reverse();
//     case "inbound-message-dsc":
//       return _.sortBy(threads, "inboundLastSeqId").reverse();
//     case "outbound-message-dsc":
//       return _.sortBy(threads, "outboundLastSeqId").reverse();
//     case "priority-asc":
//       return _.sortBy(threads, (thread) => priorityMap[thread.priority]);
//     case "priority-dsc":
//       return _.sortBy(threads, (thread) => -priorityMap[thread.priority]);
//     case "status-changed-asc":
//       return _.sortBy(threads, "statusChangedAt");
//     case "status-changed-dsc":
//       return _.sortBy(threads, "statusChangedAt").reverse();
//     default:
//       return threads;
//   }
// }
//

function sortThreads(threads: ThreadShape[], sortBy: SortBy): ThreadShape[] {
  switch (sortBy) {
    case "created-asc":
      return _.orderBy(threads, "createdAt", "asc");
    case "created-dsc":
      return _.orderBy(threads, "createdAt", "desc");
    case "inbound-message-dsc":
      return _.orderBy(threads, "inboundSeqId", "desc");
    case "outbound-message-dsc":
      return _.orderBy(threads, "outboundSeqId", "desc");
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
        state.threads?.get(threadId) ?? null,
      getWorkspaceId: (state: WorkspaceStoreState) =>
        state.workspace?.workspaceId || "",
      getWorkspaceName: (state: WorkspaceStoreState) =>
        state.workspace?.name || "",
      isInSync: (state: WorkspaceStoreState) => state.inSync,
      setCustomersShapeHandle: (handle: null | string) => {
        set((state) => {
          state.customersShapeHandle = handle;
          return state;
        });
      },
      setCustomersShapeOffset: (offset: string) => {
        set((state) => {
          state.customersShapeOffset = offset;
          return state;
        });
      },
      setInSync: (f: boolean) => {
        set((state) => {
          state.inSync = f;
          return state;
        });
      },

      setMembersShapeHandle: (handle: null | string) => {
        set((state) => {
          state.membersShapeHandle = handle;
          return state;
        });
      },
      setMembersShapeOffset: (offset: string) => {
        set((state) => {
          state.membersShapeOffset = offset;
          return state;
        });
      },
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
      setThreadsShapeHandle: (handle: null | string) => {
        set((state) => {
          state.threadsShapeHandle = handle;
          return state;
        });
      },
      setThreadsShapeOffset: (offset: string) => {
        set((state) => {
          state.threadsShapeOffset = offset;
          return state;
        });
      },
      updateCustomer: (customer: CustomerShapeUpdates) => {
        // Guard against invalid input and assert memberId is string
        if (!customer?.customerId) return;
        const id: string = customer.customerId;

        set((state) => {
          // If member doesn't exist in state, return unchanged state
          if (!state.customers?.[id]) {
            return state;
          }

          // Return new state object with updated member
          return {
            ...state,
            customers: {
              ...state.customers,
              [id]: {
                ...state.customers[id],
                ...customer,
              },
            },
          };
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
      updateMember: (member: MemberShapeUpdates) => {
        // Guard against invalid input and assert memberId is string
        if (!member?.memberId) return;
        const id: string = member.memberId;

        set((state) => {
          // If member doesn't exist in state, return unchanged state
          if (!state.members?.[id]) {
            return state;
          }

          // Return new state object with updated member
          return {
            ...state,
            members: {
              ...state.members,
              [id]: {
                ...state.members[id],
                ...member,
              },
            },
          };
        });
      },
      updateThread: (thread: ThreadShapeUpdates) => {
        // Guard against invalid input
        if (!thread?.threadId) return;
        const id: string = thread.threadId;

        set((state) => {
          // If threads Map doesn't exist or thread doesn't exist, return unchanged state
          if (!state.threads?.has(id)) {
            return state;
          }

          // Get existing thread
          const existingThread = state.threads.get(id)!;

          // Create new Map with updated thread
          const newThreads = new Map(state.threads);
          newThreads.set(id, {
            ...existingThread,
            ...thread,
          });

          // Return new state with updated threads Map
          return {
            ...state,
            threads: newThreads,
          };
        });
      },
      // updateThread: (thread) => {
      //   set((state) => {
      //     if (!state.threads || !state.threads.has(thread.threadId))
      //       return state;
      //
      //     // Get existing thread if it exists
      //     const existingThread = state.threads.get(thread.threadId);
      //
      //     // Update the thread by merging existing with new data
      //     state.threads.set(thread.threadId, {
      //       ...existingThread, // spread existing data if any
      //       ...thread, // spread new updates
      //     });
      //     return state;
      //   });
      // },
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
      viewCurrentThreadQueue: (
        state: WorkspaceStoreState,
      ): null | ThreadShape[] => {
        if (!state.threadAppliedFilters) return null;
        if (!state.threads) return [];

        const {
          assignees,
          isUnassigned,
          memberId,
          priorities,
          sortBy,
          stages,
          status,
        } = state.threadAppliedFilters;

        // Convert Map to array once
        let results = [...state.threads.values()];

        // Apply filters in order of most restrictive first
        // 1. Status filter (required)
        results = results.filter((t) => t.status === status);

        // 2. Assignment filters (mutually exclusive)
        if (memberId) {
          // Single assignee filter
          results = results.filter((t) => t.assigneeId === memberId);
        } else {
          // Handle other assignment cases
          if (isUnassigned) {
            results = results.filter((t) => !t.assigneeId);
          } else if (assignees) {
            // Only apply assignees filter if not handling memberId or isUnassigned
            results = filterByAssignees(results, assignees);
          }
        }

        // 3. Apply remaining filters
        if (stages) {
          results = filterByStages(results, stages);
        }

        if (priorities) {
          results = filterByPriorities(results, priorities);
        }

        // 4. Finally sort the filtered results
        return sortThreads(results, sortBy);
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
      viewCustomersShapeHandle: (state: WorkspaceStoreState) => {
        return state.customersShapeHandle;
      },
      viewCustomersShapeOffset: (state: WorkspaceStoreState) => {
        return state.customersShapeOffset;
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
      viewMembersShapeHandle: (state: WorkspaceStoreState) => {
        return state.membersShapeHandle;
      },
      viewMembersShapeOffset: (state: WorkspaceStoreState) => {
        return state.membersShapeOffset;
      },
      viewMyThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        memberId: string,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: SortBy = defaultSortKey,
      ): ThreadShape[] => {
        if (!state.threads) return [];

        // Convert Map to array once and apply base filters
        let results = [...state.threads.values()].filter(
          (t) => t.status === status && t.assigneeId === memberId,
        );

        // Apply remaining filters conditionally
        if (assignees) {
          results = filterByAssignees(results, assignees);
        }

        if (stages) {
          results = filterByStages(results, stages);
        }

        if (priorities) {
          results = filterByPriorities(results, priorities);
        }

        // Apply sorting
        return sortThreads(results, sortBy);
      },
      viewPats: (state: WorkspaceStoreState) => {
        const pats = state.pats ? Object.values(state.pats) : [];
        return _.sortBy(pats, "patId").reverse();
      },
      viewThreadAssigneeId: (
        state: WorkspaceStoreState,
        threadId: string,
      ): null | string => {
        return state.threads?.get(threadId)?.assigneeId ?? null;
      },
      viewThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: SortBy = defaultSortKey,
      ): ThreadShape[] => {
        if (!state.threads) return [];

        // Convert Map to array once and apply base status filter
        let results = [...state.threads.values()].filter(
          (t) => t.status === status,
        );

        // Apply remaining filters conditionally
        if (assignees) {
          results = filterByAssignees(results, assignees);
        }

        if (stages) {
          results = filterByStages(results, stages);
        }

        if (priorities) {
          results = filterByPriorities(results, priorities);
        }

        // Apply sorting
        return sortThreads(results, sortBy);
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
      viewThreadsShapeHandle(state: WorkspaceStoreState): null | string {
        return state.threadsShapeHandle;
      },
      viewThreadsShapeOffset(state: WorkspaceStoreState): string {
        return state.threadsShapeOffset;
      },
      viewUnassignedThreads: (
        state: WorkspaceStoreState,
        status: StatusType,
        assignees: AssigneesFiltersType,
        stages: StagesFiltersType,
        priorities: PrioritiesFiltersType,
        sortBy: SortBy = defaultSortKey,
      ): ThreadShape[] => {
        if (!state.threads) return [];

        // Convert Map to array once and apply base filters
        let results = [...state.threads.values()].filter(
          (t) => t.status === status && !t.assigneeId,
        );

        // Apply remaining filters conditionally
        if (assignees) {
          results = filterByAssignees(results, assignees);
        }

        if (stages) {
          results = filterByStages(results, stages);
        }

        if (priorities) {
          results = filterByPriorities(results, priorities);
        }

        // Apply sorting
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
