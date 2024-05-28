import { createStore } from "zustand/vanilla";
import { z } from "zod";
import { workspaceResponseSchema } from "./schema";
import { AccountType } from "./api";

// inferred from schema.ts - no change.
export type WorkspaceStoreType = z.infer<typeof workspaceResponseSchema>;

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

// add more like: Workspace | User | etc.
type AllowedEntities = WorkspaceStoreType | ThreadChatStoreType;

export type Dictionary<K extends string | number, V extends AllowedEntities> = {
  [key in K]: V;
};

export type ThreadChatMapStoreType = Dictionary<string, ThreadChatStoreType>;

export interface IWorkspaceEntities {
  hasData: boolean;
  isPending: boolean;
  error: Error | null;
  workspace: WorkspaceStoreType | null;
  metrics: WorkspaceMetricsStoreType;
  threadChats: ThreadChatMapStoreType | null;
}

interface IWorkspaceStoreActions {
  updateWorkspaceStore(): void;
  getWorkspaceName(state: WorkspaceStoreStateType): string;
  getWorkspaceId(state: WorkspaceStoreStateType): string;
  getMetrics(state: WorkspaceStoreStateType): WorkspaceMetricsStoreType;
  getThreadChatItem(
    state: WorkspaceStoreStateType,
    threadChatId: string
  ): ThreadChatStoreType | null;
  viewAllTodoThreads(state: WorkspaceStoreStateType): ThreadChatStoreType[];
  viewMyTodoThreads(
    state: WorkspaceStoreStateType,
    memberId: string
  ): ThreadChatStoreType[];
  viewUnassignedThreads(state: WorkspaceStoreStateType): ThreadChatStoreType[];
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
    getMetrics: (state: WorkspaceStoreStateType) => state.metrics,
    getThreadChatItem: (state: WorkspaceStoreStateType, threadChatId: string) =>
      state.threadChats?.[threadChatId] || null,
    viewAllTodoThreads: (state: WorkspaceStoreStateType) => {
      const threads = state.threadChats ? Object.values(state.threadChats) : [];
      return threads.filter((t) => t.status === "todo");
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
  }));
};

export interface IAccount {
  hasData: boolean;
  error: Error | null;
  account: AccountType | null;
}

interface IAccountStoreActions {
  updateStore(): void;
  getAccount(state: AccountStoreStateType): AccountType | null;
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
