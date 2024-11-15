import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/",
)({
  component: () => null,
  loader: ({ params }) => {
    throw redirect({
      params,
      to: "/workspaces/$workspaceId/threads/todo",
    });
  },
});
