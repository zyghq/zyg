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

const routeApi = getRouteApi("/_auth/workspaces/$workspaceId/_workspace");

function ReasonsSubMenu({
  reasons,
}: {
  reasons: string[] | string | undefined;
}) {
  const navigate = useNavigate();
  const [selectedReasons, setSelectedReasons] = React.useState<
    string | string[]
  >("");

  React.useEffect(() => {
    // check if multiple reasons are selected
    if (reasons && Array.isArray(reasons)) {
      setSelectedReasons([...reasons]);
      // check if only 1 reason(s) is selected
    } else if (reasons && typeof reasons === "string") {
      setSelectedReasons([reasons]);
      // if no reasons are selected
    } else {
      setSelectedReasons("");
    }
  }, [reasons, setSelectedReasons]);

  function onChecked(reason: string) {
    return navigate({
      search: (prev: { reasons: string[] | string }) => {
        const { reasons, ...others } = prev;

        // no existing reasons - add new reason
        if (!reasons || reasons === "") {
          return { reasons: reason, ...others };
        }

        // found a reason - merge with existing
        if (typeof reasons === "string") {
          return { reasons: [reasons, reason], ...others };
        }
        // multiple reasons selected add more to existing
        if (Array.isArray(reasons)) {
          return { reasons: [...reasons, reason], ...others };
        }
      },
    });
  }

  function onUnchecked(reason: string) {
    return navigate({
      search: (prev: { reasons: string[] | string | null }) => {
        const { reasons, ...others } = prev;

        // no existing reasons - nothing to do
        if (!reasons || reasons === "") {
          return { ...others };
        }

        // found a reason - remove it
        if (typeof reasons === "string" && reasons === reason) {
          return { ...others };
        }

        // multiple reasons selected - remove the reason
        if (Array.isArray(reasons)) {
          const filtered = reasons.filter((r) => r !== reason);
          if (filtered.length === 0) {
            return { ...others };
          }
          if (filtered.length === 1) {
            return { reasons: filtered[0], ...others };
          }
          return { reasons: filtered, ...others };
        }
      },
    });
  }

  return (
    <DropdownMenuSub>
      <DropdownMenuSubTrigger>
        Reason
        {selectedReasons && selectedReasons.length > 0 && (
          <React.Fragment>
            <Separator orientation="vertical" className="mx-1 h-3" />
            <Badge variant="secondary" className="px-1 text-xs font-normal">
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
              selectedReasons ? selectedReasons.includes("unreplied") : false
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
              selectedReasons ? selectedReasons.includes("replied") : false
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
  );
}

function PrioritiesSubMenu({
  priorities,
}: {
  priorities: string[] | string | undefined;
}) {
  const navigate = useNavigate();
  const [selectedPriorities, setSelectedPriorities] = React.useState<
    string | string[]
  >("");

  React.useEffect(() => {
    // check if multiple priorities are selected
    if (priorities && Array.isArray(priorities)) {
      setSelectedPriorities([...priorities]);
      // check if only 1 priority(s) is selected
    } else if (priorities && typeof priorities === "string") {
      setSelectedPriorities([priorities]);
      // if no priorities are selected
    } else {
      setSelectedPriorities("");
    }
  }, [priorities]);

  function onChecked(priority: string) {
    return navigate({
      search: (prev: { priorities: string[] | string }) => {
        const { priorities, ...others } = prev;

        // no existing priorities - add new priority
        if (!priorities || priorities === "") {
          return { priorities: priority, ...others };
        }

        // found a priority - merge with existing
        if (typeof priorities === "string") {
          return { priorities: [priorities, priority], ...others };
        }
        // multiple priorities selected add more to existing
        if (Array.isArray(priorities)) {
          return { priorities: [...priorities, priority], ...others };
        }
      },
    });
  }

  function onUnchecked(priority: string) {
    return navigate({
      search: (prev: { priorities: string[] | string | null }) => {
        const { priorities, ...others } = prev;

        // no existing priorities - nothing to do
        if (!priorities || priorities === "") {
          return { ...others };
        }

        // found a priority - remove it
        if (typeof priorities === "string" && priorities === priority) {
          return { ...others };
        }

        // multiple priorities selected - remove the priority
        if (Array.isArray(priorities)) {
          const filtered = priorities.filter((r) => r !== priority);
          if (filtered.length === 0) {
            return { ...others };
          }
          if (filtered.length === 1) {
            return { priorities: filtered[0], ...others };
          }
          return { priorities: filtered, ...others };
        }
      },
    });
  }

  return (
    <DropdownMenuSub>
      <DropdownMenuSubTrigger>
        Priority
        {selectedPriorities && selectedPriorities.length > 0 && (
          <React.Fragment>
            <Separator orientation="vertical" className="mx-1 h-3" />
            <Badge variant="secondary" className="px-1 text-xs font-normal">
              {selectedPriorities.length} selected
            </Badge>
          </React.Fragment>
        )}
      </DropdownMenuSubTrigger>
      <DropdownMenuPortal>
        <DropdownMenuSubContent className="mx-2 w-44 sm:w-56">
          <DropdownMenuCheckboxItem
            onSelect={(e) => e.preventDefault()}
            checked={
              selectedPriorities ? selectedPriorities.includes("urgent") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("urgent") : onUnchecked("urgent");
            }}
          >
            Urgent
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            onSelect={(e) => e.preventDefault()}
            checked={
              selectedPriorities ? selectedPriorities.includes("high") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("high") : onUnchecked("high");
            }}
          >
            High
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            onSelect={(e) => e.preventDefault()}
            checked={
              selectedPriorities ? selectedPriorities.includes("normal") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("normal") : onUnchecked("normal");
            }}
          >
            Normal
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            onSelect={(e) => e.preventDefault()}
            checked={
              selectedPriorities ? selectedPriorities.includes("low") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("low") : onUnchecked("low");
            }}
          >
            low
          </DropdownMenuCheckboxItem>
        </DropdownMenuSubContent>
      </DropdownMenuPortal>
    </DropdownMenuSub>
  );
}

export function Filters() {
  const routeSearch = routeApi.useSearch();
  const { reasons, priorities } = routeSearch;
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="border-dashed">
          <MixerHorizontalIcon className="mr-1 h-3 w-3" />
          Filters
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="sm:58 w-48" align="end">
        <DropdownMenuGroup>
          <ReasonsSubMenu reasons={reasons} />
          <PrioritiesSubMenu priorities={priorities} />
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
