import { createFileRoute } from "@tanstack/react-router";
import { QueueSize } from "@/components/workspace/insights/overview";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/insights"
)({
  component: InsightsComponent,
});

function InsightsComponent() {
  return (
    <div className="px-2 sm:px-4">
      <div className="flex mt-2 sm:mt-4">
        <QueueSize />
      </div>
    </div>
  );
}
