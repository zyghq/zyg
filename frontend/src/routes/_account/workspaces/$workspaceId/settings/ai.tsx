import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import { MagicWandIcon } from "@radix-ui/react-icons";
import { createFileRoute } from "@tanstack/react-router";
import { FlaskConicalIcon } from "lucide-react";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/ai"
)({
  component: AISettings,
});

function AISettings() {
  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="my-12">
            <header className="text-xl font-semibold">AI</header>
          </div>
          <Separator />
        </div>
        <div className="flex items-center">
          <MagicWandIcon className="mr-2 h-4 w-4" />
          <div className="text-lg">AI Workflows</div>
          <div className="ml-2 flex items-center rounded-lg border bg-muted p-1">
            <FlaskConicalIcon className="h-4 w-4" />
            <div className="text-xs">Experimental</div>
          </div>
        </div>
        <div className="pt-4 flex flex-col gap-1">
          <div className="flex flex-col items-start rounded-lg border p-3 text-left">
            <div className="flex w-full flex-col">
              <div className="flex items-center">
                <div className="mr-2 flex max-w-lg flex-col">
                  <div className="font-normal">Auto Labelling</div>
                  <div className="text-xs text-muted-foreground">
                    {`New threads will automatically be tagged with labels.
                    Threads with labels added via the API won't be affected`}
                  </div>
                </div>
                <div className="ml-auto">
                  <div className="flex">
                    <Switch id="ai-auto-label" />
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="flex flex-col items-start rounded-lg border p-3 text-left">
            <div className="flex w-full flex-col">
              <div className="flex items-center">
                <div className="mr-2 flex max-w-lg flex-col">
                  <div className="font-normal">Thread Summarisation</div>
                  <div className="text-xs text-muted-foreground">
                    {`We'll summarise your conversations`}
                  </div>
                </div>
                <div className="ml-auto">
                  <div className="flex">
                    <Switch id="ai-auto-summarize" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
