import { atom, createStore } from "jotai";

export const store = createStore();

export const workspace = atom({
  hasError: false,
  isFetching: false,
  hasData: false,
  workspaceId: "",
  accountId: "",
  name: "",
  createdAt: "",
  updatedAt: "",
});

export const threadChatMetrics = atom({
  hasError: false,
  isFetching: false,
  hasData: false,
  active: 0,
  done: 0,
  todo: 0,
  snoozed: 0,
  assignedToMe: 0,
  unassigned: 0,
  otherAssigned: 0,
  labels: [],
});
