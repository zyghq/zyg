"use client";

import { cn } from "@/lib/utils";
import Avatar from "boring-avatars";
import * as React from "react";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

import CollapsibleLabelList from "@/components/dashboard/collapsible-label-list";

import { AvatarIcon, ChatBubbleIcon } from "@radix-ui/react-icons";

export default function SidebarLinks({ workspaceId, count }) {
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
    <React.Fragment>
      <div className="flex flex-col space-y-2">
        <Button
          variant="ghost"
          asChild
          className={applyClasses(`/${workspaceId}/`)}
        >
          <Link href={`/${workspaceId}/`}>
            <div className="flex">
              <ChatBubbleIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
              <div className="my-auto">All Threads</div>
            </div>
            {applyBadge(`/${workspaceId}/`, count.active)}
          </Link>
        </Button>
        <Button
          variant="ghost"
          asChild
          className={applyClasses(`/${workspaceId}/threads/me/`)}
        >
          <Link href={`/${workspaceId}/threads/me/`}>
            <div className="flex">
              <div className="mr-2">
                <Avatar size={18} name="name" variant="beam" />
              </div>
              <div className="my-auto">My Threads</div>
            </div>
            {applyBadge(`/${workspaceId}/threads/me/`, count.assignedToMe)}
          </Link>
        </Button>
        <Button
          variant="ghost"
          asChild
          className={applyClasses(`/${workspaceId}/threads/unassigned/`)}
        >
          <Link href={`/${workspaceId}/threads/unassigned/`}>
            <div className="flex">
              <AvatarIcon className="my-auto mr-2 h-5 w-5 text-muted-foreground" />
              <div className="my-auto">Unassigned Threads</div>
            </div>
            {applyBadge(
              `/${workspaceId}/threads/unassigned/`,
              count.unassigned
            )}
          </Link>
        </Button>
      </div>
      <div className="mb-3 mt-4 text-xs text-zinc-500">Browse</div>
      <div className="flex flex-col space-y-2">
        <div className="flex">
          <CollapsibleLabelList
            workspaceId={workspaceId}
            labels={count.labels}
          />
        </div>
      </div>
    </React.Fragment>
  );
}
