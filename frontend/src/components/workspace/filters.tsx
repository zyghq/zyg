import React from "react";
import { getRouteApi, useNavigate } from "@tanstack/react-router";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Separator } from "@/components/ui/separator";

import { MixerHorizontalIcon } from "@radix-ui/react-icons";

const routeApi = getRouteApi("/workspaces/$workspaceId/_layout");

export function Filters() {
  const navigate = useNavigate();
  const [selectedReasons, setSelectedReasons] = React.useState<
    string[] | string | null
  >(null);
  const routeSearch = routeApi.useSearch();
  const { reasons } = routeSearch;

  React.useEffect(() => {
    if (reasons && Array.isArray(reasons)) {
      setSelectedReasons([...reasons]);
    } else if (reasons && typeof reasons === "string") {
      setSelectedReasons([reasons]);
    } else {
      setSelectedReasons([]);
    }
  }, [reasons]);

  function onChecked(reason: string) {
    return navigate({
      search: (prev: { reasons: string[] | string | null }) => {
        const { reasons } = prev;

        if (!reasons) {
          return { reasons: reason };
        }
        if (typeof reasons === "string") {
          return { reasons: [reasons, reason] };
        }
        if (Array.isArray(reasons)) {
          return { reasons: [...reasons, reason] };
        }
      },
    });
  }

  function onUnchecked(reason: string) {
    return navigate({
      search: (prev: { reasons: string[] | string | null }) => {
        const { reasons } = prev;

        if (!reasons) {
          return { reasons: null };
        }
        if (typeof reasons === "string" && reasons === reason) {
          return { reasons: null };
        }
        if (Array.isArray(reasons)) {
          const filtered = reasons.filter((r) => r !== reason);
          if (filtered.length === 0) {
            return { reasons: null };
          }
          if (filtered.length === 1) {
            return { reasons: filtered[0] };
          }
          return { reasons: filtered };
        }
      },
    });
  }

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
              {selectedReasons && selectedReasons.length > 0 && (
                <React.Fragment>
                  <Separator orientation="vertical" className="mx-1 h-3" />
                  <Badge
                    variant="secondary"
                    className="px-1 text-xs font-normal"
                  >
                    {selectedReasons.length} selected
                  </Badge>
                </React.Fragment>
              )}
            </DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent className="mx-2 w-44 sm:w-56">
                <DropdownMenuCheckboxItem
                  onSelect={(e) => e.preventDefault()}
                  checked={
                    selectedReasons
                      ? selectedReasons.includes("unreplied")
                      : false
                  }
                  onCheckedChange={(checked) => {
                    checked ? onChecked("unreplied") : onUnchecked("unreplied");
                  }}
                >
                  Awaiting Reply
                </DropdownMenuCheckboxItem>
                <DropdownMenuCheckboxItem
                  onSelect={(e) => e.preventDefault()}
                  checked={
                    selectedReasons
                      ? selectedReasons.includes("replied")
                      : false
                  }
                  onCheckedChange={(checked) => {
                    checked ? onChecked("replied") : onUnchecked("replied");
                  }}
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
