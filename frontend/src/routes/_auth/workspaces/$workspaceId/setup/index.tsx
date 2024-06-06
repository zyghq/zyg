import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/workspaces/$workspaceId/setup/")({
  component: WorkspaceSetup,
});

function WorkspaceSetup() {
  return <div className="flex flex-col justify-center p-4">...</div>;
}
