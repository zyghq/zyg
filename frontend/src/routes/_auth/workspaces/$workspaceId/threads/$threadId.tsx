import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/threads/$threadId"
)({
  component: () => <div>Hello /workspaces/threads/$threadId!</div>,
});
