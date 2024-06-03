import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/_auth/")({
  component: Index,
});

function Index() {
  return (
    <div className="p-2">
      <h3>Auth Index.</h3>
    </div>
  );
}
