import React from "react";
import { Link, useParams } from "@tanstack/react-router";
import { cn } from "@/lib/utils";
import { Button, buttonVariants } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ChatBubbleIcon, CodeIcon } from "@radix-ui/react-icons";
import { BlocksIcon } from "lucide-react";
import Avatar from "boring-avatars";
import { Icons } from "@/components/icons";
import { ArrowLeftIcon } from "@radix-ui/react-icons";

export function SideNavLinks({
  accountId,
  maxHeight,
}: {
  accountId: string;
  maxHeight?: string;
}) {
  const { workspaceId } = useParams({
    from: "/_auth/workspaces/$workspaceId/settings",
  });
  return (
    <React.Fragment>
      <ScrollArea className={cn("px-2", maxHeight)}>
        <div className="px-2 pb-8">
          {/* G1 */}
          <div className="my-4 flex items-center gap-1">
            <Avatar name={accountId} size={32} variant="marble" />
            <div>
              <div className="text-xs font-medium">Sanchit Rk</div>
              <div className="text-xs text-foreground">Account</div>
            </div>
          </div>
          {/* G1 Items */}
          <div className="flex flex-col gap-1">
            <Button
              variant="ghost"
              asChild
              className="flex w-full justify-between"
            >
              <Link href={`/`}>
                <div className="flex">
                  <div className="my-auto">Appearance</div>
                </div>
              </Link>
            </Button>
            <Button
              variant="ghost"
              asChild
              className="flex w-full justify-between"
            >
              <Link href={`/`}>
                <div className="flex">
                  <div className="my-auto">Personal Notifications</div>
                </div>
              </Link>
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
              to="/workspaces/$workspaceId/settings"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">General</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/members"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">Members</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/labels"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">Labels</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/ai"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">AI</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/billing"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
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
              to="/workspaces/$workspaceId/settings/slack"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">Slack</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/email"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">Email</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/chat"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
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
              to="/workspaces/$workspaceId/settings/github"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">Github</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/linear"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
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
              to="/workspaces/$workspaceId/settings/webhooks"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">Webhooks</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/pats"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: false, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">Personal Access Tokens</div>
            </Link>
            <Link
              to="/workspaces/$workspaceId/settings/events"
              params={{ workspaceId }}
              className={cn(
                buttonVariants({ variant: "ghost" }),
                "flex w-full justify-between px-3 dark:text-accent-foreground"
              )}
              activeOptions={{ exact: true, includeSearch: false }}
              activeProps={{
                className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
              }}
            >
              <div className="flex">Events</div>
            </Link>
          </div>
        </div>
      </ScrollArea>
      <div className="sticky bottom-0 flex h-14 border-t">
        <div className="flex w-full items-center">
          <div className="mx-4">
            <Button size="sm" variant="outline">
              <ArrowLeftIcon className="mr-2 h-4 w-4" />
              Support
            </Button>
          </div>
        </div>
      </div>
    </React.Fragment>
  );
}
