import React from "react";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { Button, buttonVariants } from "@/components/ui/button";
import { Link } from "@tanstack/react-router";
import { WorkspaceMetrics } from "@/db/models";
import { SideNavLabelLinks } from "@/components/workspace/sidenav-label-links";
import {
  MyThreadsLink,
  UnassignedThreadsLink,
  TodoThreadsLink,
} from "@/components/workspace/sidenav-thread-links";
import {
  OpenInNewWindowIcon,
  ReaderIcon,
  ChatBubbleIcon,
  CaretSortIcon,
  GearIcon,
  WidthIcon,
  ExitIcon,
  MagnifyingGlassIcon,
  BarChartIcon,
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
                <a target="_blank" href="https://zyg.ai/docs/">
                  <ReaderIcon className="my-auto mr-2 h-4 w-4" />
                  <div className="my-auto">Documentation</div>
                  <OpenInNewWindowIcon className="my-auto ml-2 h-4 w-4" />
                </a>
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
              to="/workspaces/$workspaceId/search"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex">
                        <MagnifyingGlassIcon className="my-auto mr-2 h-5 w-5 text-muted-foreground" />
                        <div className="font-semibold">Search</div>
                      </div>
                      <Badge className="bg-white font-mono text-muted-foreground dark:bg-black">{`/`}</Badge>
                    </>
                  ) : (
                    <>
                      <div className="flex">
                        <MagnifyingGlassIcon className="my-auto mr-2 h-5 w-5 text-muted-foreground" />
                        <div className="font-normal">Search</div>
                      </div>
                      <Badge
                        variant="outline"
                        className="font-mono text-muted-foreground bg-white dark:bg-accent"
                      >
                        {`/`}
                      </Badge>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              onClick={() => openClose(false)}
              to="/workspaces/$workspaceId/insights"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex">
                        <BarChartIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                        <div className="font-semibold">Insights</div>
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="flex">
                        <BarChartIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                        <div className="font-normal">Insights</div>
                      </div>
                    </>
                  )}
                </>
              )}
            </Link>
            <MyThreadsLink
              workspaceId={workspaceId}
              memberId={memberId}
              assignedCount={metrics.assignedToMe}
              openClose={openClose}
            />
            <UnassignedThreadsLink
              workspaceId={workspaceId}
              unassignedCount={metrics.unassigned}
              openClose={openClose}
            />
            <TodoThreadsLink
              workspaceId={workspaceId}
              activeCount={metrics.active}
              openClose={openClose}
            />
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
                  <a
                    className="flex"
                    target="_blank"
                    href="https://zyg.ai/docs/"
                  >
                    <ReaderIcon className="my-auto mr-2 h-4 w-4" />
                    <div className="my-auto">Documentation</div>
                    <OpenInNewWindowIcon className="my-auto ml-2 h-4 w-4" />
                  </a>
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
