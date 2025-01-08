import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { CodeIcon } from "@radix-ui/react-icons";
import { createFileRoute } from "@tanstack/react-router";
import { MessageCircle } from "lucide-react";
import * as React from "react";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/threads/done",
)({
  component: () => (
    <React.Fragment>
      <header className="flex h-14 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger className="-ml-1" />
        <Separator className="mr-2 h-4" orientation="vertical" />
        <div className="flex-1"></div>
      </header>
      <div className="container mt-4 max-w-md sm:mt-24">
        <div className="rounded-xl border p-4">
          <div className="flex">
            <div className="flex items-center gap-2">
              <CodeIcon className="h-5 w-5" />
              <div className="font-mono">In Development</div>
            </div>
          </div>
          <div className="font-mono text-sm text-muted-foreground">
            Soon, any threads marked as done will be shown here.
          </div>
          <div className="mt-2 flex gap-x-2">
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
    </React.Fragment>
  ),
});
