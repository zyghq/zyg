export function threadStageHumanized(key: string): string {
  switch (key) {
    case "spam":
      return "Spam";
    case "needs_first_response":
      return "Needs First Response";
    case "waiting_on_customer":
      return "Waiting on Customer";
    case "hold":
      return "Hold";
    case "needs_next_response":
      return "Needs Next Response";
    case "resolved":
      return "Resolved";
    default:
      return key;
  }
}

export function ThreadSortKeyHumanized(key: string): string {
  switch (key) {
    case "created-dsc": // when the thread is created
      return "Created, newest first";
    case "created-asc": // when the thread is created
      return "Created, oldest first";
    case "status-changed-dsc": // when the thread status is changed
      return "Status changed at, newest first";
    case "status-changed-asc": // when the thread status is changed
      return "Status changed at, oldest first";
    case "inbound-message-dsc": // when the inbound message was received
      return "Most recent message";
    case "outbound-message-dsc": // when the outbound message was sent
      return "Most recent reply";
    case "priority-asc": // when the thread priority is changed
      return "Priority, highest first";
    case "priority-dsc": // when the thread priority is changed
      return "Priority, lowest first";
    default:
      return key;
  }
}

export const sortKeys = [
  "created-dsc",
  "created-asc",
  "status-changed-dsc",
  "status-changed-asc",
  "inbound-message-dsc",
  "outbound-message-dsc",
  "priority-asc",
  "priority-dsc",
] as const;

export const todoThreadStages = [
  "needs_first_response",
  "waiting_on_customer",
  "hold",
  "needs_next_response",
] as const;
