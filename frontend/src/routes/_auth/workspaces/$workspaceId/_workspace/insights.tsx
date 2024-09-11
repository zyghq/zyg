import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/_workspace/insights"
)({
  component: () => (
    <div>Hello /_auth/workspaces/$workspaceId/_workspace/insights!</div>
  ),
});
