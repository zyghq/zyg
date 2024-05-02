import Link from "next/link";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  ReaderIcon,
  OpenInNewWindowIcon,
  ChatBubbleIcon,
  CaretSortIcon,
  GearIcon,
  WidthIcon,
  ExitIcon,
  ChevronRightIcon,
} from "@radix-ui/react-icons";
import {
  LifeBuoy as LifeBuoyIcon,
  Bug as BugIcon,
  Users as UsersIcon,
  Webhook as WebhookIcon,
  TagsIcon,
  Building2Icon,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";

export function Sidebar({ className }) {
  return (
    <div className={cn("p-4", className)}>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" className="flex justify-between">
            <div className="flex justify-start">
              <Building2Icon className="mr-2 h-5 w-5" />
              <div className="my-auto">ZygHQ</div>
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
        <div className="flex flex-col space-y-2">
          <Button
            variant="ghost"
            asChild
            // className={`w-full flex justify-between ${
            //   isActive(`/${slug}/`, pathname)
            //     ? "bg-indigo-100 hover:bg-indigo-200"
            //     : ""
            // }`}
            className="flex w-full justify-between bg-indigo-100 hover:bg-indigo-200"
          >
            <Link href={`/`}>
              <div className="flex">
                <ChatBubbleIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Threads</div>
              </div>
              <Badge className="my-auto bg-indigo-500 font-mono">18</Badge>
            </Link>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <WebhookIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Events</div>
              </div>
              <Badge className="my-auto bg-zinc-400 font-mono">37</Badge>
            </Link>
          </Button>
        </div>
        <div className="mb-3 mt-4 text-xs text-zinc-500">Browse</div>
        <div className="flex flex-col space-y-2">
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between bg-indigo-100 hover:bg-indigo-200"
          >
            <Link href={`/`}>
              <div className="flex">
                <TagsIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Labels</div>
              </div>
              <ChevronRightIcon className="my-auto h-4 w-4" />
            </Link>
          </Button>
        </div>
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
