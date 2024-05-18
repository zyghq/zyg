import { cn } from "@/lib/utils";

import Link from "next/link";

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

// import CollapsibleLabelList from "@/components/dashboard/collapsible-label-list";
import SidebarLinks from "@/components/dashboard/sidebar-links";

import {
  AvatarIcon,
  CaretSortIcon,
  ChatBubbleIcon,
  ExitIcon,
  GearIcon,
  OpenInNewWindowIcon,
  ReaderIcon,
  WidthIcon,
} from "@radix-ui/react-icons";

import {
  Bug as BugIcon,
  Building2Icon,
  LifeBuoy as LifeBuoyIcon,
  TagsIcon,
  Users as UsersIcon,
  Webhook as WebhookIcon,
} from "lucide-react";

export function Sidebar({ workspaceId, workspaceName, metrics }) {
  const { count } = metrics;

  return (
    <div className={cn("p-4", "hidden lg:block lg:border-r")}>
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
            <Link className="flex" target="_blank" href="https://zyg.ai/docs/">
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
      <ScrollArea className="my-4 h-[calc(100dvh-14rem)] pb-4">
        <SidebarLinks workspaceId={workspaceId} count={count} />
      </ScrollArea>
      <div>
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
    </div>
  );
}
