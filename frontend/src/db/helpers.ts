import { SortBy } from "@/db/store";

export function customerRoleVerboseName(key: string): string {
  switch (key) {
    case "engaged":
      return "Engaged";
    case "lead":
      return "Lead";
    case "visitor":
      return "Visitor";
    default:
      return key;
  }
}

export function ThreadSortKeyHumanized(key: SortBy): string {
  switch (key) {
    case "created-asc": // when the thread is created
      return "Created / Oldest First";
    case "created-dsc": // when the thread is created
      return "Created / Newest First";
    case "inbound-message-dsc": // when the inbound message was received
      return "Recent Message";
    case "outbound-message-dsc": // when the outbound message was sent
      return "Recent Reply";
    case "priority-asc": // when the thread priority is changed
      return "Priority / Highest First";
    case "priority-dsc": // when the thread priority is changed
      return "Priority / Lowest First";
    case "status-changed-asc": // when the thread status is changed
      return "Status Changed / Oldest First";
    case "status-changed-dsc": // when the thread status is changed
      return "Status Changed / Newest First";
    default:
      return key;
  }
}

export function threadStatusVerboseName(key: string): string {
  switch (key) {
    case "hold":
      return "Hold";
    case "needs_first_response":
      return "Needs First Response";
    case "needs_next_response":
      return "Needs Next Response";
    case "resolved":
      return "Resolved";
    case "spam":
      return "Spam";
    case "waiting_on_customer":
      return "Waiting on Customer";
    default:
      return key;
  }
}

export const sortKeys = [
  "inbound-message-dsc",
  "outbound-message-dsc",
  "priority-asc",
  "priority-dsc",
  "status-changed-dsc",
  "status-changed-asc",
  "created-dsc",
  "created-asc",
] as const;

export const todoThreadStages = [
  "needs_first_response",
  "waiting_on_customer",
  "hold",
  "needs_next_response",
] as const;

export const priorityKeys = ["urgent", "high", "normal", "low"] as const;

export function getFromLocalStorage(key: string): any | null | string {
  try {
    // Get the item from localStorage
    const item = localStorage.getItem(key);

    // Check if it's already a string, return if so
    if (item === null || item === undefined) {
      return null;
    }

    // Try to parse as JSON, if it fails return as string
    try {
      return JSON.parse(item);
    } catch {
      return item; // If not JSON, return the raw string
    }
  } catch (err) {
    console.error("Error accessing localStorage:", err);
    return null; // Return null if there's an error
  }
}

export function getInitials(name: string): string {
  // Split the name by spaces
  const nameParts = name.trim().split(/\s+/);

  if (nameParts.length === 1) {
    // If there's only one name, return the first character
    return nameParts[0].charAt(0).toUpperCase();
  } else {
    // Otherwise, return the first letter of the first and last names
    const firstInitial = nameParts[0].charAt(0).toUpperCase();
    const lastInitial = nameParts[nameParts.length - 1].charAt(0).toUpperCase();
    return firstInitial + lastInitial;
  }
}

export function setInLocalStorage(key: string, value: any) {
  try {
    // Check if the value is an object, if so, stringify it before storing
    const item = typeof value === "object" ? JSON.stringify(value) : value;

    // Store the item in localStorage
    localStorage.setItem(key, item);
  } catch (error) {
    console.error("error setting in localStorage:", error);
  }
}
