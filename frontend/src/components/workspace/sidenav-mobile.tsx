import React from "react";
import { Link } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Icons } from "@/components/icons";
import { ChatBubbleIcon } from "@radix-ui/react-icons";
import {
  CaretSortIcon,
  ExitIcon,
  GearIcon,
  OpenInNewWindowIcon,
  ReaderIcon,
  WidthIcon,
} from "@radix-ui/react-icons";
import { BugIcon, LifeBuoyIcon, UsersIcon } from "lucide-react";
import { WorkspaceMetricsStoreType } from "@/db/store";

import SideNavLinks from "@/components/workspace/sidenav-links";

export default function SideNavMobile({
  email,
  workspaceId,
  workspaceName,
  metrics,
  memberId,
}: {
  email: string;
  workspaceId: string;
  workspaceName: string;
  metrics: WorkspaceMetricsStoreType;
  memberId: string;
}) {
  const [open, setOpen] = React.useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          variant="ghost"
          className="my-auto mr-4 px-0 text-base hover:bg-transparent focus-visible:bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0 md:hidden"
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
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="p-2">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="flex justify-between">
              <div className="flex justify-start">
                <Icons.logo className="mr-2 h-5 w-5" />
                <div className="my-auto">{workspaceName}</div>
              </div>
              <CaretSortIcon className="my-auto h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start">
            <DropdownMenuLabel className="text-muted-foreground">
              {email}
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem asChild>
              <Link
                to="/workspaces/$workspaceId/settings"
                params={{ workspaceId }}
              >
                <GearIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Settings</div>
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link
                className="flex"
                target="_blank"
                href="https://zyg.ai/docs/"
              >
                <ReaderIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Documentation</div>
                <OpenInNewWindowIcon className="my-auto ml-2 h-4 w-4" />
              </Link>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem asChild>
              <Link to="/workspaces">
                <WidthIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Switch Workspace</div>
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link to="/signout">
                <ExitIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Sign Out</div>
              </Link>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <ScrollArea className="my-4 h-full pb-2">
          <SideNavLinks
            workspaceId={workspaceId}
            metrics={metrics}
            memberId={memberId}
            openClose={setOpen}
          />
        </ScrollArea>
        <div className="fixed bottom-4">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button className="flex" variant="outline">
                <LifeBuoyIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Support</div>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start">
              <DropdownMenuLabel className="text-muted-foreground">
                How can we help?
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <ChatBubbleIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Get in touch</div>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link
                  className="flex"
                  target="_blank"
                  href="https://zyg.ai/docs/"
                >
                  <ReaderIcon className="my-auto mr-2 h-4 w-4" />
                  <div className="my-auto">Documentation</div>
                  <OpenInNewWindowIcon className="my-auto ml-2 h-4 w-4" />
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <UsersIcon className="my-auto mr-2 h-4 w-4" />
                Join Slack
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <BugIcon className="my-auto mr-2 h-4 w-4" />
                Bug Report
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </SheetContent>
    </Sheet>
  );
}
