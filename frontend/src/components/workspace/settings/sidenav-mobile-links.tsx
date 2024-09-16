import { Button, buttonVariants } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from "@/components/ui/sheet";

import SideNavLinks from "@/components/workspace/settings/sidenav-links";
import { cn } from "@/lib/utils";
import { ArrowLeftIcon, HamburgerMenuIcon } from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import React from "react";

export default function SideNavMobileLinks({
  accountId,
  accountName,
  workspaceId,
}: {
  accountId: string;
  accountName: string;
  workspaceId: string;
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
        <div className="flex h-14 w-full items-center border-b">
          <div className="mx-2">
            <Link
              className={cn(buttonVariants({ size: "sm", variant: "outline" }))}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId"
            >
              <ArrowLeftIcon className="mr-2 h-4 w-4" />
              Settings
            </Link>
          </div>
        </div>
        <SideNavLinks
          accountId={accountId}
          accountName={accountName}
          maxHeight="h-[calc(100dvh-8rem)]"
        />
      </SheetContent>
    </Sheet>
  );
}
