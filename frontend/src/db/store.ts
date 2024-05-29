import { createStore } from "zustand/vanilla";
import { z } from "zod";
import { workspaceResponseSchema, membershipResponseSchema } from "./schema";
import { AccountResponseType } from "./api";

// inferred from schema
export type WorkspaceStoreType = z.infer<typeof workspaceResponseSchema>;

// inferred from schema
export type MembershipStoreType = z.infer<typeof membershipResponseSchema>;

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

// add more like: Workspace | User | etc.
type AllowedEntities =
  | WorkspaceStoreType
  | ThreadChatStoreType
  | WorkspaceCustomerStoreType;

export type Dictionary<K extends string | number, V extends AllowedEntities> = {
  [key in K]: V;
};

export type ThreadChatMapStoreType = Dictionary<string, ThreadChatStoreType>;

export type WorkspaceCustomerMapStoreType = Dictionary<
  string,
  WorkspaceCustomerStoreType
>;

export interface IWorkspaceEntities {
  hasData: boolean;
  isPending: boolean;
  error: Error | null;
  workspace: WorkspaceStoreType | null;
  member: MembershipStoreType | null;
  metrics: WorkspaceMetricsStoreType;
  threadChats: ThreadChatMapStoreType | null;
  customers: WorkspaceCustomerMapStoreType | null;
}

type ReplyStatus = "replied" | "unreplied";

export type reasonsFiltersType = ReplyStatus | ReplyStatus[] | undefined;

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
    reasons: reasonsFiltersType
  ): ThreadChatStoreType[];
  viewMyTodoThreads(
    state: WorkspaceStoreStateType,
    memberId: string
  ): ThreadChatStoreType[];
  viewUnassignedThreads(state: WorkspaceStoreStateType): ThreadChatStoreType[];
  viewCustomerName(state: WorkspaceStoreStateType, customerId: string): string;
}

export type WorkspaceStoreStateType = IWorkspaceEntities &
  IWorkspaceStoreActions;

// (sanchitrk) for reference on using zustand, check this great article:
// https://tkdodo.eu/blog/working-with-zustand
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
    viewAllTodoThreads: (state: WorkspaceStoreStateType, reasons) => {
      const threads = state.threadChats ? Object.values(state.threadChats) : [];
      const todoThreads = threads.filter((t) => t.status === "todo");
      if (reasons && Array.isArray(reasons)) {
        const uniqueReasons = [...new Set(reasons)];
        const filtered = [];
        for (const reason of uniqueReasons) {
          if (reason === "replied") {
            filtered.push(...todoThreads.filter((t) => t.replied));
          } else if (reason === "unreplied") {
            filtered.push(...todoThreads.filter((t) => !t.replied));
          }
        }
        return filtered;
      }
      if (reasons && typeof reasons === "string") {
        if (reasons === "replied") {
          return todoThreads.filter((t) => t.replied);
        }
        if (reasons === "unreplied") {
          return todoThreads.filter((t) => !t.replied);
        }
      }
      return todoThreads;
    },
    viewMyTodoThreads: (state: WorkspaceStoreStateType, memberId: string) => {
      const threads = state.threadChats ? Object.values(state.threadChats) : [];
      return threads.filter(
        (t) => t.status === "todo" && t.assigneeId === memberId
      );
    },
    viewUnassignedThreads: (state: WorkspaceStoreStateType) => {
      const threads = state.threadChats ? Object.values(state.threadChats) : [];
      return threads.filter((t) => t.status === "todo" && !t.assigneeId);
    },
    viewCustomerName: (state: WorkspaceStoreStateType, customerId: string) => {
      const customer = state.customers?.[customerId];
      return customer ? customer.name : "";
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
