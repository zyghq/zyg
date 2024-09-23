import { Button } from "@/components/ui/button";
import { CodeIcon } from "@radix-ui/react-icons";
import { createFileRoute } from "@tanstack/react-router";
import { MessageCircle } from "lucide-react";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/labels/$labelId"
)({
  component: () => (
    <div className="container mt-4 sm:mt-24 max-w-md">
      <div className="border p-4 rounded-xl">
        <div className="flex">
          <div className="flex items-center gap-2">
            <CodeIcon className="h-5 w-5" />
            <div className="font-mono">In Development</div>
          </div>
        </div>
        <div className="text-muted-foreground text-sm font-mono">
          Soon you'll be able to list threads by assigned labels.
        </div>
        <div className="flex gap-x-2 mt-2">
          <Button className="w-full sm:w-auto">Talk to Us</Button>
          <a
            className="w-full sm:w-auto"
            href="https://github.com/zyghq/zyg/discussions"
            rel="noopener noreferrer"
            target="_blank"
          >
            <Button className="w-full" variant="outline">
              <MessageCircle className="mr-2 h-4 w-4" />
              Start Discussion in GitHub
            </Button>
          </a>
        </div>
      </div>
    </div>
  ),
});
