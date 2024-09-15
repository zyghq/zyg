import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { ThreadList } from "@/components/workspace/thread/threads";
import { Thread } from "@/db/models";
import { PanelLeftIcon } from "lucide-react";
import * as React from "react";

export function SidePanelThreadList({
  threads,
  title,
  workspaceId,
}: {
  threads: Thread[];
  title: string;
  workspaceId: string;
}) {
  const [open, setOpen] = React.useState(false);

  return (
    <Sheet onOpenChange={setOpen} open={open}>
      <SheetTrigger asChild>
        <Button className="md:hidden" size="icon" variant="outline">
          <PanelLeftIcon className="h-4 w-4" />
          <span className="sr-only">Toggle Thread Panel</span>
        </Button>
      </SheetTrigger>
      <SheetContent className="w-full px-0" side="left">
        <div className="flex h-14 flex-col justify-center border-b px-4">
          <div className="font-semibold">{title}</div>
        </div>
        <ScrollArea className="h-[calc(100dvh-8rem)]">
          <ThreadList
            threads={threads}
            variant="compress"
            workspaceId={workspaceId}
          />
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}
