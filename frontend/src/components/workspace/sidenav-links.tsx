import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ScrollArea } from "@/components/ui/scroll-area";
import { WorkspaceMetrics } from "@/db/models";
import { LabelMetrics } from "@/db/models";
import { defaultSortOp } from "@/lib/search-params";
import { cn } from "@/lib/utils";
import {
  BarChartIcon,
  CaretSortIcon,
  ChatBubbleIcon,
  CircleIcon,
  ExitIcon,
  GearIcon,
  MagnifyingGlassIcon,
  OpenInNewWindowIcon,
  PersonIcon,
  ReaderIcon,
  WidthIcon,
} from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import {
  Bug as BugIcon,
  Building2Icon,
  LifeBuoy as LifeBuoyIcon,
  TagIcon,
  Users as UsersIcon,
  LocateIcon,
  TagsIcon,
  ReplyIcon,
  ClockIcon,
  PauseIcon,
  CheckCircleIcon,
} from "lucide-react";
import React from "react";

function SideNavLabelLinks({
  labels,
  workspaceId,
}: {
  labels: LabelMetrics[];
  workspaceId: string;
}) {
  const [isOpen, setIsOpen] = React.useState(false);
  return (
    <Collapsible
      className="w-full space-y-2"
      onOpenChange={setIsOpen}
      open={isOpen}
    >
      <div className="flex items-center justify-between space-x-1">
        <Button className="w-full pl-3" variant="ghost">
          <div className="mr-auto flex">
            <TagsIcon className="my-auto mr-2 h-4 w-4" />
            <div className="font-normal">Labels</div>
          </div>
        </Button>
        <CollapsibleTrigger asChild>
          <Button size="icon" variant="ghost">
            <CaretSortIcon className="h-4 w-4" />
            <span className="sr-only">Toggle</span>
          </Button>
        </CollapsibleTrigger>
      </div>
      <CollapsibleContent className="space-y-1">
        {labels.map((label) => (
          <Link
            activeOptions={{ exact: true }}
            activeProps={{
              className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
            }}
            className={cn(
              buttonVariants({ variant: "ghost" }),
              "flex w-full justify-between px-3 dark:text-accent-foreground"
            )}
            key={label.labelId}
            params={{ labelId: label.labelId, workspaceId }}
            search={{ sort: defaultSortOp }}
            to="/workspaces/$workspaceId/threads/labels/$labelId"
          >
            {({ isActive }) => (
              <>
                {isActive ? (
                  <>
                    <div className="flex gap-x-2 ml-2">
                      <TagIcon className="my-auto h-3 w-3 text-muted-foreground" />
                      <div className="font-normal capitalize text-foreground">
                        {label.name}
                      </div>
                    </div>
                    <div className="font-mono font-light text-muted-foreground">
                      {label.count}
                    </div>
                  </>
                ) : (
                  <>
                    <div className="flex gap-x-2 ml-2">
                      <TagIcon className="my-auto h-3 w-3 text-muted-foreground" />
                      <div className="font-normal capitalize text-foreground">
                        {label.name}
                      </div>
                    </div>
                    <div className="font-mono font-light text-muted-foreground">
                      {label.count}
                    </div>
                  </>
                )}
              </>
            )}
          </Link>
        ))}
      </CollapsibleContent>
    </Collapsible>
  );
}

export default function SideNavLinks({
  email,
  maxHeight,
  memberId,
  metrics,
  openClose = () => {},
  workspaceId,
  workspaceName,
}: {
  email: string;
  maxHeight: string;
  memberId: string;
  metrics: WorkspaceMetrics;
  openClose?: (isOpen: boolean) => undefined | void;
  workspaceId: string;
  workspaceName: string;
}) {
  return (
    <React.Fragment>
      <ScrollArea className={maxHeight}>
        <div className="p-2 sm:p-4">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button className="flex justify-between" variant="outline">
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
                  params={{ workspaceId }}
                  to="/workspaces/$workspaceId/settings"
                >
                  <GearIcon className="my-auto mr-2 h-4 w-4" />
                  <div className="my-auto">Settings</div>
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <a href="https://zyg.ai/docs/" target="_blank">
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
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/search"
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
                        className="font-mono text-muted-foreground bg-white dark:bg-accent"
                        variant="outline"
                      >
                        {`/`}
                      </Badge>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/insights"
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
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ sort: defaultSortOp }}
              to="/workspaces/$workspaceId/threads/me"
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex gap-x-2">
                        <Avatar className="h-5 w-5">
                          <AvatarImage
                            alt={memberId}
                            src={`https://avatar.vercel.sh/${memberId}`}
                          />
                          <AvatarFallback>CN</AvatarFallback>
                        </Avatar>
                        <div className="font-semibold">Your Threads</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {metrics.assignedToMe}
                      </span>
                    </>
                  ) : (
                    <>
                      <div className="flex gap-x-2">
                        <Avatar className="h-5 w-5">
                          <AvatarImage
                            alt={memberId}
                            src={`https://avatar.vercel.sh/${memberId}`}
                          />
                          <AvatarFallback>M</AvatarFallback>
                        </Avatar>
                        <div className="font-normal">Your Threads</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {metrics.assignedToMe}
                      </span>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ sort: defaultSortOp }}
              to="/workspaces/$workspaceId/threads/unassigned"
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex">
                        <PersonIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                        <div className="font-semibold">Unassigned Threads</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {metrics.unassigned}
                      </span>
                    </>
                  ) : (
                    <>
                      <div className="flex">
                        <PersonIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                        <div className="font-normal">Unassigned Threads</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {metrics.unassigned}
                      </span>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ sort: defaultSortOp }}
              to="/workspaces/$workspaceId/threads/todo"
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex gap-x-2">
                        <CircleIcon className="my-auto h-4 w-4 text-indigo-500" />
                        <div className="font-semibold">Todo</div>
                      </div>
                      <Badge className="bg-indigo-500 font-mono text-white hover:bg-indigo-600">
                        {metrics.active}
                      </Badge>
                    </>
                  ) : (
                    <>
                      <div className="flex gap-x-2">
                        <CircleIcon className="my-auto h-4 w-4 text-indigo-500" />
                        <div className="font-normal">Todo</div>
                      </div>
                      <Badge className="font-mono font-light" variant="outline">
                        {metrics.active}
                      </Badge>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ sort: defaultSortOp }}
              to="/workspaces/$workspaceId/threads/needs-first-response"
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex gap-x-2 ml-4">
                        <LocateIcon className="my-auto w-4 h-4 text-indigo-500" />
                        <div className="font-medium">Needs First Response</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {"00"}
                      </span>
                    </>
                  ) : (
                    <>
                      <div className="flex gap-x-2 ml-4">
                        <LocateIcon className="my-auto w-4 h-4 text-indigo-500" />
                        <div className="font-normal">Needs First Response</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {"00"}
                      </span>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ sort: defaultSortOp }}
              to="/workspaces/$workspaceId/threads/needs-next-response"
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex gap-x-2 ml-4">
                        <ReplyIcon className="my-auto w-4 h-4 text-indigo-500" />
                        <div className="font-medium">Needs Next Response</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {"00"}
                      </span>
                    </>
                  ) : (
                    <>
                      <div className="flex gap-x-2 ml-4">
                        <ReplyIcon className="my-auto w-4 h-4 text-indigo-500" />
                        <div className="font-normal">Needs Next Response</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {"00"}
                      </span>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ sort: defaultSortOp }}
              to="/workspaces/$workspaceId/threads/hold"
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex gap-x-2 ml-4">
                        <PauseIcon className="my-auto w-4 h-4 text-indigo-500" />
                        <div className="font-medium">Hold</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {"00"}
                      </span>
                    </>
                  ) : (
                    <>
                      <div className="flex gap-x-2 ml-4">
                        <PauseIcon className="my-auto w-4 h-4 text-indigo-500" />
                        <div className="font-normal">Hold</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {"00"}
                      </span>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ sort: defaultSortOp }}
              to="/workspaces/$workspaceId/threads/waiting-on-customer"
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex gap-x-2 ml-4">
                        <ClockIcon className="my-auto w-4 h-4 text-indigo-500" />
                        <div className="font-medium">Waiting on Customer</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {"00"}
                      </span>
                    </>
                  ) : (
                    <>
                      <div className="flex gap-x-2 ml-4">
                        <ClockIcon className="my-auto w-4 h-4 text-indigo-500" />
                        <div className="font-normal">Waiting on Customer</div>
                      </div>
                      <span className="font-mono text-muted-foreground pr-2">
                        {"00"}
                      </span>
                    </>
                  )}
                </>
              )}
            </Link>
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-green-100 hover:bg-green-200 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              onClick={() => openClose(false)}
              params={{ workspaceId }}
              search={{ sort: defaultSortOp }}
              to="/workspaces/$workspaceId/threads/done"
            >
              {({ isActive }: { isActive: boolean }) => (
                <>
                  {isActive ? (
                    <>
                      <div className="flex gap-x-2">
                        <CheckCircleIcon className="my-auto h-4 w-4 text-green-600" />
                        <div className="font-semibold">Done</div>
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="flex gap-x-2">
                        <CheckCircleIcon className="my-auto h-4 w-4 text-green-600" />
                        <div className="font-normal">Done</div>
                      </div>
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
                labels={metrics.labels}
                workspaceId={workspaceId}
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
                <Button size="sm" variant="outline">
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
                    href="https://zyg.ai/docs/"
                    target="_blank"
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
