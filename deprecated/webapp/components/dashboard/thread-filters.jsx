"use client";

import * as React from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Separator } from "@/components/ui/separator";

import { useReasonFilter } from "@/components/dashboard/hooks";

import { MixerHorizontalIcon } from "@radix-ui/react-icons";

export function ThreadFilterDropDownMenu() {
  const reasonFilter = useReasonFilter();
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="border-dashed">
          <MixerHorizontalIcon className="mr-1 h-3 w-3" />
          Filters
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="sm:58 w-48" align="start">
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>
              Reason
              {reasonFilter.countSelectedReasons() > 0 && (
                <React.Fragment>
                  <Separator orientation="vertical" className="mx-1 h-3" />
                  <Badge
                    variant="secondary"
                    className="px-1 text-xs font-normal"
                  >
                    {reasonFilter.countSelectedReasons()} selected
                  </Badge>
                </React.Fragment>
              )}
            </DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent className="mx-2 w-44 sm:w-56">
                <DropdownMenuCheckboxItem
                  onSelect={(e) => e.preventDefault()}
                  checked={reasonFilter.isChecked("unreplied")}
                  onCheckedChange={(checked) =>
                    checked
                      ? reasonFilter.setReason("unreplied")
                      : reasonFilter.clearReason("unreplied")
                  }
                >
                  Awaiting Reply
                </DropdownMenuCheckboxItem>
                <DropdownMenuCheckboxItem
                  onSelect={(e) => e.preventDefault()}
                  checked={reasonFilter.isChecked("replied")}
                  onCheckedChange={(checked) =>
                    checked
                      ? reasonFilter.setReason("replied")
                      : reasonFilter.clearReason("replied")
                  }
                >
                  Replied
                </DropdownMenuCheckboxItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>Priority</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent className="mx-2 w-44 sm:w-56">
                <DropdownMenuItem>Urgent</DropdownMenuItem>
                <DropdownMenuItem>High</DropdownMenuItem>
                <DropdownMenuItem>Normal</DropdownMenuItem>
                <DropdownMenuItem>Low</DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>Assignee</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent className="mx-2">
                <DropdownMenuItem>Email</DropdownMenuItem>
                <DropdownMenuItem>Message</DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem>More...</DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
