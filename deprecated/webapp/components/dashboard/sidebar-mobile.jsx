"use client";

import { cn } from "@/lib/utils";
import Avatar from "boring-avatars";
import * as React from "react";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { usePathname } from "next/navigation";

import { Badge } from "@/components/ui/badge";
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

import CollapsibleLabelListMobile from "@/components/dashboard/collapsible-label-list-mobile";
import { Icons } from "@/components/icons";

import {
  ArrowLeftIcon,
  ChatBubbleIcon,
  ChevronRightIcon,
} from "@radix-ui/react-icons";
import {
  AvatarIcon,
  CaretSortIcon,
  ExitIcon,
  GearIcon,
  OpenInNewWindowIcon,
  ReaderIcon,
  WidthIcon,
} from "@radix-ui/react-icons";

import {
  BugIcon,
  LifeBuoyIcon,
  TagsIcon,
  UsersIcon,
  WebhookIcon,
} from "lucide-react";

export function SidebarMobile({ workspaceId, workspaceName, metrics }) {
  const [open, setOpen] = React.useState(false);

  const { count } = metrics;

  const pathname = usePathname();
  const isActive = (path) => pathname === path;
  const activeClasses = "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent";

  const applyClasses = (path) => {
    return cn(
      "flex w-full justify-between px-3 dark:text-accent-foreground",
      isActive(path) ? activeClasses : ""
    );
  };

  const applyBadge = (path, count) => {
    if (isActive(path)) {
      return (
        <Badge className="bg-indigo-500 font-mono text-white hover:bg-indigo-600">
          {count}
        </Badge>
      );
    }
    return (
      <Badge variant="outline" className="font-mono font-light">
        {count}
      </Badge>
    );
  };

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
          <DropdownMenuContent className="mx-1">
            <DropdownMenuLabel className="text-zinc-500">
              sanchitrrk@gmail.com
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <GearIcon className="my-auto mr-2 h-4 w-4" />
              <div className="my-auto">Settings</div>
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
              <Link href="/workspaces/">
                <WidthIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Switch Workspace</div>
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <ExitIcon className="my-auto mr-2 h-4 w-4" />
              <div className="my-auto">Logout</div>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <ScrollArea className="my-4 h-full pb-4">
          <div className="flex flex-col space-y-2">
            <Button
              variant="ghost"
              asChild
              className={applyClasses(`/${workspaceId}/`)}
            >
              <MobileLink href={`/${workspaceId}/`} onOpenChange={setOpen}>
                <div className="flex">
                  <ChatBubbleIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                  <div className="my-auto">All Threads</div>
                </div>
                {applyBadge(`/${workspaceId}/`, count.active)}
              </MobileLink>
            </Button>
            <Button
              variant="ghost"
              asChild
              className={applyClasses(`/${workspaceId}/threads/me/`)}
            >
              <MobileLink
                href={`/${workspaceId}/threads/me/`}
                onOpenChange={setOpen}
              >
                <div className="flex">
                  <div className="mr-2">
                    <Avatar size={18} name="name" variant="beam" />
                  </div>
                  <div className="my-auto">My Threads</div>
                </div>
                {applyBadge(`/${workspaceId}/threads/me/`, count.assignedToMe)}
              </MobileLink>
            </Button>
            <Button
              variant="ghost"
              asChild
              className={applyClasses(`/${workspaceId}/threads/unassigned/`)}
            >
              <MobileLink
                href={`/${workspaceId}/threads/unassigned/`}
                onOpenChange={setOpen}
              >
                <div className="flex">
                  <AvatarIcon className="my-auto mr-2 h-5 w-5 text-muted-foreground" />
                  <div className="my-auto">Unassigned Threads</div>
                </div>
                {applyBadge(
                  `/${workspaceId}/threads/unassigned/`,
                  count.unassigned
                )}
              </MobileLink>
            </Button>
          </div>
          <div className="mb-3 mt-4 text-xs text-zinc-500">Browse</div>
          <div className="flex flex-col space-y-2">
            <div className="flex">
              <CollapsibleLabelListMobile
                workspaceId={workspaceId}
                labels={count.labels}
                sheetOpen={setOpen}
              />
            </div>
          </div>
        </ScrollArea>
        <div className="fixed bottom-4">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button className="flex" variant="outline">
                <LifeBuoyIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Support</div>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="mx-1">
              <DropdownMenuLabel className="text-zinc-500">
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

function MobileLink({ href, onOpenChange, className, children, ...props }) {
  const router = useRouter();
  return (
    <Link
      href={href}
      onClick={() => {
        router.push(href.toString());
        onOpenChange?.(false);
      }}
      className={cn(className)}
      {...props}
    >
      {children}
    </Link>
  );
}
