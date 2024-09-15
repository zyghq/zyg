import React from "react";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { HamburgerMenuIcon } from "@radix-ui/react-icons";
import { WorkspaceMetrics } from "@/db/models";
import SideNavLinks from "@/components/workspace/sidenav-links";

export default function SideNavMobileLinks({
  email,
  workspaceId,
  workspaceName,
  metrics,
  memberId,
}: {
  email: string;
  workspaceId: string;
  workspaceName: string;
  metrics: WorkspaceMetrics;
  memberId: string;
}) {
  const [open, setOpen] = React.useState(false);
  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          size="icon"
          variant="ghost"
          className="flex md:hidden my-auto mr-4"
        >
          <HamburgerMenuIcon className="h-4 w-4" />
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="p-0">
        <SideNavLinks
          maxHeight="h-[calc(100dvh-4rem)]"
          workspaceId={workspaceId}
          workspaceName={workspaceName}
          metrics={metrics}
          memberId={memberId}
          email={email}
          openClose={setOpen}
        />
      </SheetContent>
    </Sheet>
  );
}
