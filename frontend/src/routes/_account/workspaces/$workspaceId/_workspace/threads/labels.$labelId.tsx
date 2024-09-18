import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/labels/$labelId"
)({
  component: () => <div>hmm, work in progress</div>,
});
