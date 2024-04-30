"use client";
import * as React from "react";
// import Link from "next/link";
// import { Icons } from "@/components/icons";
// import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
// import { ScrollArea } from "@/components/ui/scroll-area";
import { PanelLeftIcon } from "lucide-react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";

export function SidePanelThreadList() {
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
          <div className="flex h-8 flex-col border-b">
            <div className="text-sm font-semibold">Threads</div>
          </div>
          {/* <ThreadList
                items={threads}
                className="h-[calc(100dvh-8rem)] pr-0"
                variant="compress"
              /> */}
          ... load some threads
        </div>
      </SheetContent>
    </Sheet>
  );
}
