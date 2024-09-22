import { createFileRoute, Outlet } from "@tanstack/react-router";
import { threadSearchSchema } from "@/lib/search-params";
import { zodSearchValidator } from "@tanstack/router-zod-adapter";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads"
)({
  validateSearch: zodSearchValidator(threadSearchSchema),
  component: Outlet,
});
