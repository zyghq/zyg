import React from "react";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { Button, buttonVariants } from "@/components/ui/button";

import Avatar from "boring-avatars";
import { Link, getRouteApi } from "@tanstack/react-router";

import { WorkspaceMetrics } from "@/db/entities";
import { SideNavLabelLinks } from "@/components/workspace/sidenav-label-links";
import {
  OpenInNewWindowIcon,
  ReaderIcon,
  ChatBubbleIcon,
  CaretSortIcon,
  GearIcon,
  WidthIcon,
  ExitIcon,
} from "@radix-ui/react-icons";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Building2Icon,
  Bug as BugIcon,
  LifeBuoy as LifeBuoyIcon,
  Users as UsersIcon,
} from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";

const routeApi = getRouteApi("/_auth/workspaces/$workspaceId/_workspace");

export default function SideNavLinks({
  maxHeight,
  workspaceId,
  workspaceName,
  metrics,
  memberId,
  email,
  openClose = () => {},
}: {
  maxHeight: string;
  workspaceId: string;
  workspaceName: string;
  metrics: WorkspaceMetrics;
  memberId: string;
  email: string;
  openClose?: (isOpen: boolean) => void | undefined;
}) {
  const routeSearch = routeApi.useSearch();
  const { status, sort } = routeSearch;
  return (
    <React.Fragment>
      <ScrollArea className={maxHeight}>
        <div className="p-2 sm:p-4">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" className="flex justify-between">
                <div className="flex justify-start">
                  <Building2Icon className="mr-2 h-5 w-5" />
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
                <Link target="_blank" href="https://zyg.ai/docs/">
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
          <div className="mt-4 flex flex-col space-y-1">
            <Link
              onClick={() => openClose(false)}
              to="/workspaces/$workspaceId"
              params={{ workspaceId }}
              search={{ status: status }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              {({ isActive }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex">
                        <ChatBubbleIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                        <div className="font-normal">All Threads</div>
                      </div>
                      <Badge className="bg-indigo-500 font-mono text-white hover:bg-indigo-600">
                        {metrics.active}
                      </Badge>
                    </>
                  ) : (
                    <>
                      <div className="flex">
                        <ChatBubbleIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                        <div className="font-normal">All Threads</div>
                      </div>
                      <Badge variant="outline" className="font-mono font-light">
                        {metrics.active}
                      </Badge>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              to="/workspaces/$workspaceId/me"
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ status, sort }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              {({ isActive }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex">
                        {/* <ChatBubbleIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" /> */}
                        <div className="flex my-auto mr-2">
                          <Avatar size={18} name={memberId} variant="marble" />
                        </div>
                        <div className="font-normal">My Threads</div>
                      </div>
                      <Badge className="bg-indigo-500 font-mono text-white hover:bg-indigo-600">
                        {metrics.assignedToMe}
                      </Badge>
                    </>
                  ) : (
                    <>
                      <div className="flex">
                        {/* <ChatBubbleIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" /> */}
                        <div className="flex my-auto mr-2">
                          <Avatar size={18} name={memberId} variant="marble" />
                        </div>
                        <div className="font-normal">My Threads</div>
                      </div>
                      <Badge variant="outline" className="font-mono font-light">
                        {metrics.assignedToMe}
                      </Badge>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              to="/workspaces/$workspaceId/unassigned"
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ status, sort }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              {({ isActive }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex">
                        <ChatBubbleIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                        <div className="font-normal">Unassigned Threads</div>
                      </div>
                      <Badge className="bg-indigo-500 font-mono text-white hover:bg-indigo-600">
                        {metrics.unassigned}
                      </Badge>
                    </>
                  ) : (
                    <>
                      <div className="flex">
                        <ChatBubbleIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                        <div className="font-normal">Unassigned Threads</div>
                      </div>
                      <Badge variant="outline" className="font-mono font-light">
                        {metrics.unassigned}
                      </Badge>
                    </>
                  )}
                </>
              )}
            </Link>
          </div>
          <div className="mb-3 mt-4 text-xs text-muted-foreground">Browse</div>
          <div className="flex flex-col space-y-2">
            <div className="flex">
              <SideNavLabelLinks
                workspaceId={workspaceId}
                labels={metrics.labels}
              />
            </div>
          </div>
        </div>
      </ScrollArea>
      <div className="sticky bottom-0 flex h-14 border-t">
        <div className="flex w-full items-center">
          <div className="mx-4">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <LifeBuoyIcon className="mr-2 h-4 w-4" />
                  Support
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
        </div>
      </div>
    </React.Fragment>
  );
}
