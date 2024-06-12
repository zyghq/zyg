import { createFileRoute } from "@tanstack/react-router";
import { Separator } from "@/components/ui/separator";
import { CookingPotIcon } from "lucide-react";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/settings/chat"
)({
  component: ComingSoonSettings,
});

function ComingSoonSettings() {
  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="my-12">
            <header className="text-xl font-semibold">Chat</header>
          </div>
          <Separator />
        </div>
        <div className="flex text-muted-foreground">
          <div className="mr-1">Will be adding chat settings soon.</div>
          <CookingPotIcon className="h-4 w-4 my-auto" />
        </div>
      </div>
    </div>
  );
}
