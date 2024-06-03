import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/recover")({
  component: () => <div>Hello /recover!</div>,
});
