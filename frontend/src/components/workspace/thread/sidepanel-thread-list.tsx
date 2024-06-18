import * as React from "react";

import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ThreadList } from "@/components/workspace/thread/threads";
import { PanelLeftIcon } from "lucide-react";
import { ThreadChatStoreType } from "@/db/store";

export function SidePanelThreadList({
  threads,
  title,
  workspaceId,
}: {
  threads: ThreadChatStoreType[];
  title: string;
  workspaceId: string;
}) {
  const [open, setOpen] = React.useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button variant="outline" size="icon" className="md:hidden">
          <PanelLeftIcon className="h-4 w-4" />
          <span className="sr-only">Toggle Thread Panel</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="w-full px-1">
        <ScrollArea className="h-full px-2">
          <div className="flex flex-col">
            <div className="flex h-8 flex-col">
              <div className="text-sm font-semibold">{title}</div>
            </div>
            <ThreadList
              workspaceId={workspaceId}
              threads={threads}
              variant="compress"
            />
          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}
