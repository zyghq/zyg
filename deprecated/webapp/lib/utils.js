import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs) {
  return twMerge(clsx(inputs));
}

export function titleCase(str) {
  return str.replace(/\b\w/g, (l) => l.toUpperCase());
}
