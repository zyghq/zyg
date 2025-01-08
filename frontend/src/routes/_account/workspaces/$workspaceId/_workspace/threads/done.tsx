import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { createFileRoute } from "@tanstack/react-router";
import { CheckCircleIcon } from "lucide-react";
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
      <div className="p-4">
        <div className="flex items-center space-x-2">
          <CheckCircleIcon className="h-5 w-5 text-green-600" />
          <span className="font-serif text-lg font-medium sm:text-xl">
            {"Done"}
          </span>
        </div>
      </div>
    </React.Fragment>
  ),
});
