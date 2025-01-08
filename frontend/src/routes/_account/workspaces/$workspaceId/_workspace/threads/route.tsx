import { threadSearchSchema } from "@/lib/search-params";
import { useWorkspaceStore } from "@/providers";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { zodSearchValidator } from "@tanstack/router-zod-adapter";
import * as React from "react";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads",
)({
  component: ThreadRoute,
  validateSearch: zodSearchValidator(threadSearchSchema),
});

function ThreadRoute() {
  const { sort } = Route.useSearch();
  const workspaceStore = useWorkspaceStore();

  React.useEffect(() => {
    workspaceStore.getState().setThreadSortKey(sort);
  }, [sort, workspaceStore]);

  return <Outlet />;
}
