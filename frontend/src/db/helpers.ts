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

export const todoThreadStages = [
  "needs_first_response",
  "waiting_on_customer",
  "hold",
  "needs_next_response",
];
