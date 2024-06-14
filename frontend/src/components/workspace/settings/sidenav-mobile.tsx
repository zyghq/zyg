import React from "react";
import { Link } from "@tanstack/react-router";
import { ArrowLeftIcon, HamburgerMenuIcon } from "@radix-ui/react-icons";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { SideNavLinks } from "@/components/workspace/settings/sidenav-links";

export function SideNavMobile({
  workspaceId,
  accountId,
}: {
  workspaceId: string;
  accountId: string;
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
        <div className="flex h-14 w-full items-center border-b">
          <div className="mx-2">
            <Link
              to="/workspaces/$workspaceId"
              params={{ workspaceId }}
              className={cn(buttonVariants({ variant: "outline", size: "sm" }))}
            >
              <ArrowLeftIcon className="mr-2 h-4 w-4" />
              Settings
            </Link>
          </div>
        </div>
        <SideNavLinks accountId={accountId} maxHeight="h-[calc(100dvh-8rem)]" />
      </SheetContent>
    </Sheet>
  );
}
