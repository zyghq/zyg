import { Separator } from "@/components/ui/separator";
import { CodeIcon } from "@radix-ui/react-icons";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/notifications"
)({
  component: NotificationSettings,
});

function NotificationSettings() {
  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="my-12">
            <header className="text-xl font-semibold">
              Your Notifications
            </header>
          </div>
          <Separator />
        </div>
        <div className="flex text-muted-foreground gap-x-1">
          <CodeIcon className="h-5 w-5 my-auto" />
          <div className="mr-1">
            Notifications is not yet supported. Will be adding soon.
          </div>
        </div>
      </div>
    </div>
  );
}
