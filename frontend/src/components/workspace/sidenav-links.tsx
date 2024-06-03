import React from "react";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { buttonVariants } from "@/components/ui/button";

import Avatar from "boring-avatars";
import { Link, getRouteApi } from "@tanstack/react-router";

import { ChatBubbleIcon } from "@radix-ui/react-icons";
import { WorkspaceMetricsStoreType } from "@/db/store";
import { SideNavLabelLinks } from "@/components/workspace/sidenav-label-links";

const routeApi = getRouteApi("/_auth/workspaces/$workspaceId/_workspace");

export default function SideNavLinks({
  workspaceId,
  metrics,
  memberId,
  openClose = () => {},
}: {
  workspaceId: string;
  metrics: WorkspaceMetricsStoreType;
  memberId: string;
  openClose?: (isOpen: boolean) => void | undefined;
}) {
  const routeSearch = routeApi.useSearch();
  const { status, sort } = routeSearch;
  return (
    <React.Fragment>
      <div className="flex flex-col space-y-2">
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
    </React.Fragment>
  );
}
