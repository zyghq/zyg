import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/workspaces/add")({
  component: () => <div>Hello /workspaces/add!</div>,
});
