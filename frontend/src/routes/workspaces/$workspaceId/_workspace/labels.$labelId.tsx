import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/workspaces/$workspaceId/_workspace/labels/$labelId"
)({
  component: () => <div>hmm, work in progress</div>,
});
