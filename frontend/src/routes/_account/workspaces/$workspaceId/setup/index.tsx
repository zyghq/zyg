import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_account/workspaces/$workspaceId/setup/")({
  component: WorkspaceSetup,
});

function WorkspaceSetup() {
  return <div className="flex flex-col justify-center p-4">...</div>;
}
