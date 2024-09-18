import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/waiting-on-customer"
)({
  component: () => (
    <div>
      Hello
      /_account/workspaces/$workspaceId/_workspace/threads/waiting-on-customer!
    </div>
  ),
});
