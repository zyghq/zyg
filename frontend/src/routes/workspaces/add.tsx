import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/workspaces/add")({
  component: () => <div>Hello /workspaces/add!</div>,
});
