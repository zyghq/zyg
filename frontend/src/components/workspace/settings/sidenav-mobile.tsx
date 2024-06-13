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
        {/* <Button
          variant="ghost"
          className="my-auto px-0 text-base hover:bg-transparent focus-visible:bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0 md:hidden"
        >
          <svg
            strokeWidth="1.5"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
          >
            <path
              d="M3 5H11"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            ></path>
            <path
              d="M3 12H16"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            ></path>
            <path
              d="M3 19H21"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            ></path>
          </svg>
          <span className="sr-only">Toggle Menu</span>
        </Button> */}
      </SheetTrigger>
      <SheetContent side="left" className="p-0">
        <div className="flex h-14 w-full items-center border-b">
          <div className="mx-4">
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
        <SideNavLinks accountId={accountId} maxHeight="h-[calc(100dvh-7rem)]" />
      </SheetContent>
    </Sheet>
  );
}
