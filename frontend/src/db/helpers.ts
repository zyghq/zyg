import { SortBy } from "@/db/store";

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

export function ThreadSortKeyHumanized(key: SortBy): string {
  switch (key) {
    case "created-asc": // when the thread is created
      return "Created, oldest first";
    case "created-dsc": // when the thread is created
      return "Created, newest first";
    case "inbound-message-dsc": // when the inbound message was received
      return "Most recent message";
    case "outbound-message-dsc": // when the outbound message was sent
      return "Most recent reply";
    case "priority-asc": // when the thread priority is changed
      return "Priority, highest first";
    case "priority-dsc": // when the thread priority is changed
      return "Priority, lowest first";
    case "status-changed-asc": // when the thread status is changed
      return "Status changed at, oldest first";
    case "status-changed-dsc": // when the thread status is changed
      return "Status changed at, newest first";
    default:
      return key;
  }
}

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
    } catch (error) {
      return item; // If not JSON, return the raw string
    }
  } catch (error) {
    console.error("Error accessing localStorage:", error);
    return null; // Return null if there's an error
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
