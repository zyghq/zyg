import { createStore } from "zustand/vanilla";
import { z } from "zod";
import _ from "lodash";
import {
  workspaceResponseSchema,
  membershipResponseSchema,
  workspaceLabelResponseSchema,
  workspaceMemberResponseSchema,
  accountPatSchema,
} from "./schema";
import { AccountResponseType } from "./api";

// inferred from schema
export type WorkspaceStoreType = z.infer<typeof workspaceResponseSchema>;

// inferred from schema
export type MembershipStoreType = z.infer<typeof membershipResponseSchema>;

// inferred from schema
export type WorkspaceLabelStoreType = z.infer<
  typeof workspaceLabelResponseSchema
>;

// inferred from schema
export type WorkspaceMemberStoreType = z.infer<
  typeof workspaceMemberResponseSchema
>;

// inferred from schema
export type AccountPatStoreType = z.infer<typeof accountPatSchema>;

export type LabelMetricsStoreType = {
  labelId: string;
  name: string;
  icon: string;
  count: number;
};

export type WorkspaceMetricsStoreType = {
  active: number;
  done: number;
  snoozed: number;
  assignedToMe: number;
  unassigned: number;
  otherAssigned: number;
  labels: LabelMetricsStoreType[] | [];
};

export type ThreadChatStoreType = {
  threadChatId: string;
  sequence: number;
  status: string;
  read: boolean;
  replied: boolean;
  customerId: string;
  assigneeId: string | null;
  createdAt: string;
  updatedAt: string;
  recentMessage: {
    threadChatId: string;
    threadChatMessageId: string;
    body: string;
    sequence: number;
    customerId: string | null;
    memberId: string | null;
    createdAt: string;
    updatedAt: string;
  } | null;
};

export type WorkspaceCustomerStoreType = {
  workspaceId: string;
  customerId: string;
  externalId: string | null;
  email: string | null;
  phone: string | null;
  name: string;
  createdAt: string;
  updatedAt: string;
};

// add more entitites as supported by store
// e.g: Workspace | User | etc.
type AllowedEntities =
  | WorkspaceStoreType
  | ThreadChatStoreType
  | WorkspaceCustomerStoreType
  | WorkspaceLabelStoreType
  | WorkspaceMemberStoreType
  | AccountPatStoreType;

export type Dictionary<K extends string | number, V extends AllowedEntities> = {
  [key in K]: V;
};

export type ThreadChatMapStoreType = Dictionary<string, ThreadChatStoreType>;

export type WorkspaceCustomerMapStoreType = Dictionary<
  string,
  WorkspaceCustomerStoreType
>;

export type WorkspaceLabelMapStoreType = Dictionary<
  string,
  WorkspaceLabelStoreType
>;

export type WorkspaceMemberMapStoreType = Dictionary<
  string,
  WorkspaceMemberStoreType
>;

export type AccountPatMapStoreType = Dictionary<string, AccountPatStoreType>;

export interface IWorkspaceEntities {
  hasData: boolean;
  isPending: boolean;
  error: Error | null;
  workspace: WorkspaceStoreType | null;
  member: MembershipStoreType | null;
  metrics: WorkspaceMetricsStoreType;
  threadChats: ThreadChatMapStoreType | null;
  customers: WorkspaceCustomerMapStoreType | null;
  labels: WorkspaceLabelMapStoreType | null;
  members: WorkspaceMemberMapStoreType | null;
  pats: AccountPatMapStoreType | null;
}

type ReplyStatus = "replied" | "unreplied";

export type reasonsFiltersType = ReplyStatus | ReplyStatus[] | undefined;

export type sortByType = "last-message-dsc" | "created-asc" | "created-dsc";

interface IWorkspaceStoreActions {
  updateWorkspaceStore(): void;
  getWorkspaceName(state: WorkspaceStoreStateType): string;
  getWorkspaceId(state: WorkspaceStoreStateType): string;
  getMemberId(state: WorkspaceStoreStateType): string;
  getMemberName(state: WorkspaceStoreStateType): string;
  getMemberRole(state: WorkspaceStoreStateType): string;
  getMetrics(state: WorkspaceStoreStateType): WorkspaceMetricsStoreType;
  getThreadChatItem(
    state: WorkspaceStoreStateType,
    threadChatId: string
  ): ThreadChatStoreType | null;
  viewAllTodoThreads(
    state: WorkspaceStoreStateType,
    reasons: reasonsFiltersType,
    sortBy: sortByType
  ): ThreadChatStoreType[];
  viewMyTodoThreads(
    state: WorkspaceStoreStateType,
    memberId: string,
    reasons: reasonsFiltersType,
    sortBy: sortByType
  ): ThreadChatStoreType[];
  viewUnassignedThreads(
    state: WorkspaceStoreStateType,
    reasons: reasonsFiltersType,
    sortBy: sortByType
  ): ThreadChatStoreType[];
  viewCustomerName(state: WorkspaceStoreStateType, customerId: string): string;
  updateWorkspaceName(name: string): void;
  viewLabels(state: WorkspaceStoreStateType): WorkspaceLabelStoreType[];
  viewMembers(state: WorkspaceStoreStateType): WorkspaceMemberStoreType[];
  viewPats(state: WorkspaceStoreStateType): AccountPatStoreType[];
  addPat(pat: AccountPatStoreType): void;
}

export type WorkspaceStoreStateType = IWorkspaceEntities &
  IWorkspaceStoreActions;

function filterByReasons(
  threads: ThreadChatStoreType[],
  reasons: reasonsFiltersType
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

function sortThreads(threads: ThreadChatStoreType[], sortBy: sortByType) {
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
export const buildStore = (initialState: IWorkspaceEntities) => {
  return createStore<WorkspaceStoreStateType>()((set) => ({
    ...initialState,
    updateWorkspaceStore: () => set((state) => ({ ...state })),
    getWorkspaceName: (state: WorkspaceStoreStateType) =>
      state.workspace?.name || "",
    getWorkspaceId: (state: WorkspaceStoreStateType) =>
      state.workspace?.workspaceId || "",
    getMemberId: (state: WorkspaceStoreStateType) =>
      state.member?.memberId || "",
    getMemberName: (state: WorkspaceStoreStateType) => state.member?.name || "",
    getMemberRole: (state: WorkspaceStoreStateType) => state.member?.role || "",
    getMetrics: (state: WorkspaceStoreStateType) => state.metrics,
    getThreadChatItem: (state: WorkspaceStoreStateType, threadChatId: string) =>
      state.threadChats?.[threadChatId] || null,
    viewAllTodoThreads: (
      state: WorkspaceStoreStateType,
      reasons: reasonsFiltersType,
      sortBy: sortByType = "last-message-dsc"
    ) => {
      const threads = state.threadChats ? Object.values(state.threadChats) : [];
      const todoThreads = threads.filter((t) => t.status === "todo");
      const reasonsFiltered = filterByReasons(todoThreads, reasons);
      const sortedThreads = sortThreads(reasonsFiltered, sortBy);
      return sortedThreads;
    },
    viewMyTodoThreads: (
      state: WorkspaceStoreStateType,
      memberId: string,
      reasons: reasonsFiltersType,
      sortBy: sortByType = "last-message-dsc"
    ) => {
      const threads = state.threadChats ? Object.values(state.threadChats) : [];
      const myThreads = threads.filter(
        (t) => t.status === "todo" && t.assigneeId === memberId
      );
      const reasonsFiltered = filterByReasons(myThreads, reasons);
      const sortedThreads = sortThreads(reasonsFiltered, sortBy);
      return sortedThreads;
    },
    viewUnassignedThreads: (
      state: WorkspaceStoreStateType,
      reasons,
      sortBy = "last-message-dsc"
    ) => {
      const threads = state.threadChats ? Object.values(state.threadChats) : [];
      const unassignedThreads = threads.filter(
        (t) => t.status === "todo" && !t.assigneeId
      );
      const reasonsFiltered = filterByReasons(unassignedThreads, reasons);
      const sortedThreads = sortThreads(reasonsFiltered, sortBy);
      return sortedThreads;
    },
    viewCustomerName: (state: WorkspaceStoreStateType, customerId: string) => {
      const customer = state.customers?.[customerId];
      return customer ? customer.name : "";
    },
    viewLabels: (state: WorkspaceStoreStateType) => {
      const labels = state.labels ? Object.values(state.labels) : [];
      return labels;
    },
    viewMembers: (state: WorkspaceStoreStateType) => {
      const members = state.members ? Object.values(state.members) : [];
      return members;
    },
    viewPats: (state: WorkspaceStoreStateType) => {
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
    addPat: (pat: AccountPatStoreType) => {
      const { patId } = pat;
      set((state) => {
        if (state.pats) {
          state.pats[patId] = { ...pat };
          return state;
        } else {
          return state;
        }
      });
    },
  }));
};

export interface IAccount {
  hasData: boolean;
  error: Error | null;
  account: AccountResponseType | null;
}

interface IAccountStoreActions {
  updateStore(): void;
  getAccount(state: AccountStoreStateType): AccountResponseType | null;
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
