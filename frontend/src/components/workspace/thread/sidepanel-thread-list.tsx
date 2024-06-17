import * as React from "react";

import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";

import { PanelLeftIcon } from "lucide-react";

export function SidePanelThreadList({ title }: { title: string }) {
  const [open, setOpen] = React.useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button variant="outline" size="icon" className="md:hidden">
          <PanelLeftIcon className="h-4 w-4" />
          <span className="sr-only">Toggle Thread Panel</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="w-full px-4">
        <div className="flex h-full flex-col">
          <div className="flex h-8 flex-col">
            <div className="text-sm font-semibold">{title}</div>
          </div>
          ...
        </div>
      </SheetContent>
    </Sheet>
  );
}
