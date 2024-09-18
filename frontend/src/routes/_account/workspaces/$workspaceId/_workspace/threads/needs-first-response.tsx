import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/needs-first-response"
)({
  component: () => (
    <div>
      Hello
      /_account/workspaces/$workspaceId/_workspace/threads/needs-first-response!
    </div>
  ),
});
