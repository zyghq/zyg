"use client";

import * as React from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  ArrowLeftIcon,
  ChatBubbleIcon,
  ChevronRightIcon,
} from "@radix-ui/react-icons";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import {
  CaretSortIcon,
  GearIcon,
  ReaderIcon,
  OpenInNewWindowIcon,
  ExitIcon,
  WidthIcon,
} from "@radix-ui/react-icons";
import {
  LifeBuoyIcon,
  BugIcon,
  UsersIcon,
  WebhookIcon,
  TagsIcon,
} from "lucide-react";
import { Icons } from "@/components/icons";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import SettingsMenuScroll from "@/components/settings-menu-scroll";

// rename this to DashboardNav
export function MobileNav() {
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
          <div className="flex flex-col space-y-3">
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
                <Badge className="my-auto bg-indigo-500 font-mono hover:bg-indigo-700">
                  18
                </Badge>
              </Link>
            </Button>
            <Button
              variant="ghost"
              asChild
              className="flex w-full justify-between"
            >
              <MobileLink href={`/`} onOpenChange={setOpen}>
                <div className="flex">
                  <WebhookIcon className="my-auto mr-2 h-4 w-4" />
                  <div className="my-auto">Events</div>
                </div>
                <Badge className="my-auto bg-zinc-400 font-mono">37</Badge>
              </MobileLink>
            </Button>
          </div>
          <div className="mb-3 mt-4 text-xs text-zinc-500">Browse</div>
          <div className="flex flex-col space-y-3">
            <Button
              variant="ghost"
              asChild
              className="flex w-full justify-between bg-indigo-100 hover:bg-indigo-200"
            >
              <MobileLink href={`/`} onOpenChange={setOpen}>
                <div className="flex">
                  <TagsIcon className="my-auto mr-2 h-4 w-4" />
                  <div className="my-auto">Labels</div>
                </div>
                <ChevronRightIcon className="my-auto h-4 w-4" />
              </MobileLink>
            </Button>
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

export function SettingsMobileNav() {
  const [open, setOpen] = React.useState(false);
  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          variant="ghost"
          className="my-auto px-0 text-base hover:bg-transparent focus-visible:bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0 md:hidden"
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
      <SheetContent side="left" className="p-0">
        <div className="flex h-14 w-full items-center border-b">
          <div className="mx-4">
            <Button size="sm" variant="outline">
              <ArrowLeftIcon className="mr-2 h-4 w-4" />
              Settings
            </Button>
          </div>
        </div>
        <SettingsMenuScroll className="h-[calc(100dvh-7rem)]" />
        <div className="bottom-0 flex h-14 w-full border-t">
          <div className="flex w-full items-center">
            <div className="mx-4">
              <Button size="sm" variant="outline">
                <ArrowLeftIcon className="mr-2 h-4 w-4" />
                Support
              </Button>
            </div>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
