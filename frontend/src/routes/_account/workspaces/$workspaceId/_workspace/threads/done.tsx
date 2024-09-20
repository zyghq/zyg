import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/done"
)({
  component: () => (
    <div>Hello /_account/workspaces/$workspaceId/_workspace/threads/done!</div>
  ),
});