import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
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
import { Separator } from "@/components/ui/separator";
import { Assignee } from "@/db/store";
import { cn } from "@/lib/utils";
import { CheckIcon, MixerHorizontalIcon } from "@radix-ui/react-icons";
import { getRouteApi, useNavigate } from "@tanstack/react-router";
import Avatar from "boring-avatars";
import React from "react";

const routeApi = getRouteApi(
  "/_account/workspaces/$workspaceId/_workspace/threads"
);

function ReasonsSubMenu({
  reasons,
}: {
  reasons: string | string[] | undefined;
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
      search: (prev: { reasons: string | string[] }) => {
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
      search: (prev: { reasons: null | string | string[] }) => {
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
            <Separator className="mx-1 h-3" orientation="vertical" />
            <Badge className="px-1 text-xs font-normal p-0" variant="secondary">
              {selectedReasons.length} selected
            </Badge>
          </React.Fragment>
        )}
      </DropdownMenuSubTrigger>
      <DropdownMenuPortal>
        <DropdownMenuSubContent className="mx-2 w-44 sm:w-56">
          <DropdownMenuCheckboxItem
            checked={
              selectedReasons ? selectedReasons.includes("unreplied") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("unreplied") : onUnchecked("unreplied");
            }}
            onSelect={(e) => e.preventDefault()}
          >
            Awaiting Reply
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            checked={
              selectedReasons ? selectedReasons.includes("replied") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("replied") : onUnchecked("replied");
            }}
            onSelect={(e) => e.preventDefault()}
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
  priorities: string | string[] | undefined;
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
      search: (prev: { priorities: string | string[] }) => {
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
      search: (prev: { priorities: null | string | string[] }) => {
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
            <Separator className="mx-1 h-3" orientation="vertical" />
            <Badge className="px-1 text-xs font-normal p-0" variant="secondary">
              {selectedPriorities.length} selected
            </Badge>
          </React.Fragment>
        )}
      </DropdownMenuSubTrigger>
      <DropdownMenuPortal>
        <DropdownMenuSubContent className="mx-2 w-44 sm:w-56">
          <DropdownMenuCheckboxItem
            checked={
              selectedPriorities ? selectedPriorities.includes("urgent") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("urgent") : onUnchecked("urgent");
            }}
            onSelect={(e) => e.preventDefault()}
          >
            Urgent
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            checked={
              selectedPriorities ? selectedPriorities.includes("high") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("high") : onUnchecked("high");
            }}
            onSelect={(e) => e.preventDefault()}
          >
            High
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            checked={
              selectedPriorities ? selectedPriorities.includes("normal") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("normal") : onUnchecked("normal");
            }}
            onSelect={(e) => e.preventDefault()}
          >
            Normal
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            checked={
              selectedPriorities ? selectedPriorities.includes("low") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("low") : onUnchecked("low");
            }}
            onSelect={(e) => e.preventDefault()}
          >
            low
          </DropdownMenuCheckboxItem>
        </DropdownMenuSubContent>
      </DropdownMenuPortal>
    </DropdownMenuSub>
  );
}

function AssigneeSubMenu({
  assignedMembers,
  assignees,
}: {
  assignedMembers: Assignee[];
  assignees: string | string[] | undefined;
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
      search: (prev: { assignees: string | string[] }) => {
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
          const uniques = [...new Set([member, ...assignees])];
          return { assignees: uniques, ...others };
        }
      },
    });
  }

  function onUnchecked(member: string) {
    return navigate({
      search: (prev: { assignees: null | string | string[] }) => {
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
            <Separator className="mx-1 h-3" orientation="vertical" />
            <Badge className="px-1 text-xs font-normal p-0" variant="secondary">
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
                    className="text-sm"
                    key={m.assigneeId}
                    onSelect={() => onSelect(m.assigneeId)}
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
  assignedMembers: Assignee[];
  disableAssigneeFilter?: boolean;
}) {
  const routeSearch = routeApi.useSearch();
  const { assignees, priorities, reasons } = routeSearch;
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button className="border-dashed" size="sm" variant="outline">
          <MixerHorizontalIcon className="mr-1 h-3 w-3" />
          Filters
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="sm:58 w-48 mx-1">
        <DropdownMenuGroup>
          <ReasonsSubMenu reasons={reasons} />
          <PrioritiesSubMenu priorities={priorities} />
        </DropdownMenuGroup>
        {!disableAssigneeFilter && (
          <React.Fragment>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <AssigneeSubMenu
                assignedMembers={assignedMembers}
                assignees={assignees}
              />
            </DropdownMenuGroup>
          </React.Fragment>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
