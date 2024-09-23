import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import SideNavLinks from "@/components/workspace/sidenav-links";
import { WorkspaceMetrics } from "@/db/models";
import { HamburgerMenuIcon } from "@radix-ui/react-icons";
import React from "react";

export default function SideNavMobileLinks({
  email,
  memberId,
  metrics,
  workspaceId,
  workspaceName,
}: {
  email: string;
  memberId: string;
  metrics: WorkspaceMetrics;
  workspaceId: string;
  workspaceName: string;
}) {
  const [open, setOpen] = React.useState(false);
  return (
    <Sheet onOpenChange={setOpen} open={open}>
      <SheetTrigger asChild>
        <Button
          className="flex md:hidden my-auto mr-4"
          size="icon"
          variant="ghost"
        >
          <HamburgerMenuIcon className="h-4 w-4" />
        </Button>
      </SheetTrigger>
      <SheetContent className="p-0" side="left">
        {/* adding this to stop aria warnings. */}
        <SheetHeader className="hidden">
          <SheetTitle>Open Menu</SheetTitle>
          <SheetDescription>
            Select menu items from the left sidebar to navigate to different
            pages.
          </SheetDescription>
        </SheetHeader>
        <SideNavLinks
          email={email}
          maxHeight="h-[calc(100dvh-4rem)]"
          memberId={memberId}
          metrics={metrics}
          openClose={setOpen}
          workspaceId={workspaceId}
          workspaceName={workspaceName}
        />
      </SheetContent>
    </Sheet>
  );
}
