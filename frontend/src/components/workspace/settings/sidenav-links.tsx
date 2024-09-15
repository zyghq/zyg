import { Icons } from "@/components/icons";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";
import {
  ChatBubbleIcon,
  CodeIcon,
  OpenInNewWindowIcon,
  ReaderIcon,
} from "@radix-ui/react-icons";
import { Link, useParams } from "@tanstack/react-router";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { BlocksIcon } from "lucide-react";
import {
  Bug as BugIcon,
  LifeBuoy as LifeBuoyIcon,
  Users as UsersIcon,
} from "lucide-react";
import React from "react";

export default function SideNavLinks({
  accountId,
  accountName,
  maxHeight,
}: {
  accountId: string;
  accountName: string;
  maxHeight?: string;
}) {
  const { workspaceId } = useParams({
    from: "/_account/workspaces/$workspaceId/settings",
  });
  return (
    <React.Fragment>
      <ScrollArea className={maxHeight}>
        <div className="p-4">
          {/* G1 */}
          <div className="mb-4 flex items-center gap-x-2">
            <Avatar className="h-5 w-5">
              <AvatarImage
                alt={accountId}
                src={`https://avatar.vercel.sh/${accountId}`}
              />
              <AvatarFallback>CN</AvatarFallback>
            </Avatar>
            <div>
              <div className="text-xs font-medium">{accountName || "User"}</div>
              <div className="text-xs text-muted-foreground">Account</div>
            </div>
          </div>
          {/* G1 Items */}
          <div className="flex flex-col gap-1">
            <Button
              asChild
              className="flex w-full justify-between"
              variant="ghost"
            >
              <a href={`/`}>
                <div className="flex">
                  <div className="my-auto">Appearance</div>
                </div>
              </a>
            </Button>
            <Button
              asChild
              className="flex w-full justify-between"
              variant="ghost"
            >
              <a href={`/`}>
                <div className="flex">
                  <div className="my-auto">Personal Notifications</div>
                </div>
              </a>
            </Button>
          </div>
          {/* G2 */}
          <div className="my-4 flex items-center gap-1">
            <Icons.logo className="mx-1 h-5 w-5" />
            <div>
              <div className="text-xs font-medium">ZygHQ</div>
              <div className="text-xs text-foreground">Workspace</div>
            </div>
          </div>
          {/* G2 Items */}
          <div className="flex flex-col gap-1">
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings"
            >
              <div className="flex">General</div>
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
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/members"
            >
              <div className="flex">Members</div>
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
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/labels"
            >
              <div className="flex">Labels</div>
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
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/ai"
            >
              <div className="flex">AI</div>
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
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/billing"
            >
              <div className="flex">Billing</div>
            </Link>
          </div>
          {/* G3 */}
          <div className="my-4 flex items-center gap-1">
            <ChatBubbleIcon className="mx-1 h-4 w-4 text-muted-foreground" />
            <div className="text-xs font-medium text-muted-foreground">
              Channels
            </div>
          </div>
          {/* G3 Items */}
          <div className="flex flex-col gap-1">
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/slack"
            >
              <div className="flex">Slack</div>
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
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/email"
            >
              <div className="flex">Email</div>
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
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/chat"
            >
              <div className="flex">Chat</div>
            </Link>
          </div>
          {/* G4 */}
          <div className="my-4 flex items-center gap-1">
            <BlocksIcon className="mx-1 h-4 w-4 text-muted-foreground" />
            <div className="text-xs font-medium text-muted-foreground">
              Integrations
            </div>
          </div>
          {/* G4 Items */}
          <div className="flex flex-col gap-1">
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/github"
            >
              <div className="flex">Github</div>
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
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/linear"
            >
              <div className="flex">Linear</div>
            </Link>
          </div>
          {/* G5 */}
          <div className="my-4 flex items-center gap-1">
            <CodeIcon className="mx-1 h-4 w-4 text-muted-foreground" />
            <div className="text-xs font-medium text-muted-foreground">
              Build
            </div>
          </div>
          {/* G5 Items */}
          <div className="flex flex-col gap-1">
            <Link
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/webhooks"
            >
              <div className="flex">Webhooks</div>
            </Link>
            <Link
              activeOptions={{ exact: false, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
              }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/pats"
            >
              <div className="flex">Personal Access Tokens</div>
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
              params={{ workspaceId }}
              to="/workspaces/$workspaceId/settings/events"
            >
              <div className="flex">Events</div>
            </Link>
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
