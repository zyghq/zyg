"use client";
import * as React from "react";

import { Button } from "@/components/ui/button";

import { PanelLeftIcon } from "lucide-react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import ThreadList from "@/components/thread/thread-list";

export function SidePanelThreadList({ workspaceId, threads }) {
  const [open, setOpen] = React.useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button variant="outline" size="icon" className="md:hidden">
          <PanelLeftIcon className="h-4 w-4" />
          <span className="sr-only">Toggle Menu</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="w-full px-4">
        <div className="flex h-full flex-col">
          <div className="flex h-8 flex-col">
            <div className="text-sm font-semibold">Threads</div>
          </div>
          <ThreadList
            workspaceId={workspaceId}
            items={threads}
            className="h-[calc(100dvh-8rem)] py-1"
          />
        </div>
      </SheetContent>
    </Sheet>
  );
}
