import * as React from "react";

import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ThreadList } from "@/components/workspace/thread/threads";
import { PanelLeftIcon } from "lucide-react";
import { ThreadChatWithRecentMessage } from "@/db/entities";

export function SidePanelThreadList({
  threads,
  title,
  workspaceId,
}: {
  threads: ThreadChatWithRecentMessage[];
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
      <SheetContent side="left" className="w-full px-0">
        <div className="flex h-14 flex-col justify-center border-b px-4">
          <div className="font-semibold">{title}</div>
        </div>
        <ScrollArea className="h-[calc(100dvh-8rem)]">
          <ThreadList
            workspaceId={workspaceId}
            threads={threads}
            variant="compress"
          />
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}
