import React from "react";
import { cn } from "@/lib/utils";
import { getRouteApi, useNavigate } from "@tanstack/react-router";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";

import { Separator } from "@/components/ui/separator";
import { MixerHorizontalIcon, CheckIcon } from "@radix-ui/react-icons";
import Avatar from "boring-avatars";
import { AssigneeType } from "@/db/store";

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
            <Badge variant="secondary" className="px-1 text-xs font-normal p-0">
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
            <Badge variant="secondary" className="px-1 text-xs font-normal p-0">
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

function AssigneeSubMenu({
  assignees,
  assignedMembers,
}: {
  assignees: string[] | string | undefined;
  assignedMembers: AssigneeType[];
}) {
  const navigate = useNavigate();
  const [selectedMembers, setSelectedMembers] = React.useState<
    string | string[]
  >("");

  React.useEffect(() => {
    // check if multiple members are selected
    if (assignees && Array.isArray(assignees)) {
      setSelectedMembers([...assignees]);
      // check if only 1 member(s) is selected
    } else if (assignees && typeof assignees === "string") {
      setSelectedMembers([assignees]);
      // if no members are selected
    } else {
      setSelectedMembers("");
    }
  }, [assignees]);

  function onChecked(member: string) {
    return navigate({
      search: (prev: { assignees: string[] | string }) => {
        // search params
        const { assignees, ...others } = prev;

        // no existing members - add new member
        if (!assignees || assignees === "") {
          return { assignees: member, ...others };
        }

        // found a member - merge with existing
        if (typeof assignees === "string") {
          return { assignees: [assignees, member], ...others };
        }
        // multiple members selected add more to existing
        if (Array.isArray(assignees)) {
          const uniques = [...new Set([...assignees, member])];
          return { assignees: uniques, ...others };
        }
      },
    });
  }

  function onUnchecked(member: string) {
    return navigate({
      search: (prev: { assignees: string[] | string | null }) => {
        const { assignees, ...others } = prev;

        // no existing members - nothing to do
        if (!assignees || assignees === "") {
          return { ...others };
        }

        // found a member - remove it
        if (typeof assignees === "string" && assignees === member) {
          return { ...others };
        }

        // multiple members selected - remove the member
        if (Array.isArray(assignees)) {
          const filtered = assignees.filter((r) => r !== member);
          if (filtered.length === 0) {
            return { ...others };
          }
          if (filtered.length === 1) {
            return { assignees: filtered[0], ...others };
          }
          return { assignees: filtered, ...others };
        }
      },
    });
  }

  function onSelect(member: string) {
    const isChecked =
      member === selectedMembers ||
      (Array.isArray(selectedMembers) && selectedMembers.includes(member));

    if (isChecked) {
      onUnchecked(member);
    } else {
      onChecked(member);
    }
  }

  const isChecked = (member: string) => {
    const t =
      member === selectedMembers ||
      (Array.isArray(selectedMembers) && selectedMembers.includes(member));
    return t;
  };

  return (
    <DropdownMenuSub>
      <DropdownMenuSubTrigger>
        Assignee
        {selectedMembers && selectedMembers.length > 0 && (
          <React.Fragment>
            <Separator orientation="vertical" className="mx-1 h-3" />
            <Badge variant="secondary" className="px-1 text-xs font-normal p-0">
              {selectedMembers.length} selected
            </Badge>
          </React.Fragment>
        )}
      </DropdownMenuSubTrigger>
      <DropdownMenuPortal>
        <DropdownMenuSubContent className="mx-2">
          <Command>
            <CommandList>
              <CommandInput placeholder="Filter" />
              <CommandEmpty>No results</CommandEmpty>
              <CommandGroup>
                {assignedMembers.map((m) => (
                  <CommandItem
                    key={m.assigneeId}
                    onSelect={() => onSelect(m.assigneeId)}
                    className="text-sm"
                  >
                    <div className="flex gap-2">
                      <Avatar name={m.assigneeId} size={20} />
                      {m.name}
                    </div>
                    <CheckIcon
                      className={cn(
                        "ml-auto h-4 w-4",
                        isChecked(m.assigneeId) ? "opacity-100" : "opacity-0"
                      )}
                    />
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </DropdownMenuSubContent>
      </DropdownMenuPortal>
    </DropdownMenuSub>
  );
}

export function Filters({
  assignedMembers,
  disableAssigneeFilter = false,
}: {
  assignedMembers: AssigneeType[];
  disableAssigneeFilter?: boolean;
}) {
  const routeSearch = routeApi.useSearch();
  const { reasons, priorities, assignees } = routeSearch;
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="border-dashed">
          <MixerHorizontalIcon className="mr-1 h-3 w-3" />
          Filters
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="sm:58 w-48 mx-1" align="end">
        <DropdownMenuGroup>
          <ReasonsSubMenu reasons={reasons} />
          <PrioritiesSubMenu priorities={priorities} />
        </DropdownMenuGroup>
        {!disableAssigneeFilter && (
          <React.Fragment>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <AssigneeSubMenu
                assignees={assignees}
                assignedMembers={assignedMembers}
              />
            </DropdownMenuGroup>
          </React.Fragment>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
