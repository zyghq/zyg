import { createFileRoute, Outlet } from "@tanstack/react-router";
import { threadSearchSchema } from "@/lib/search-params";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads"
)({
  validateSearch: (search) => threadSearchSchema.parse(search),
  component: Outlet,
});
