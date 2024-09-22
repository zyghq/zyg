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
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { PriorityIcons, stageIcon } from "@/components/icons";
import {
  PrioritiesFiltersType,
  AssigneesFiltersType,
  StagesFiltersType,
} from "@/db/store";
import { todoThreadStages, threadStageHumanized } from "@/db/helpers";
import React from "react";

function StagesSubMenu({
  stages,
  onChecked = () => {},
  onUnchecked = () => {},
}: {
  stages: string | string[] | undefined;
  onChecked?: (stage: string) => void;
  onUnchecked?: (stage: string) => void;
}) {
  const [selectedStages, setSelectedStages] = React.useState<
    string | string[] | undefined
  >("");

  React.useEffect(() => {
    // check if multiple stages are selected
    if (stages && Array.isArray(stages)) {
      setSelectedStages([...stages]);
      // check if only 1 stage(s) is selected
    } else if (stages && typeof stages === "string") {
      setSelectedStages([stages]);
      // if no stages are selected
    } else {
      setSelectedStages(undefined);
    }
  }, [stages, setSelectedStages]);

  return (
    <DropdownMenuSub>
      <DropdownMenuSubTrigger>
        Status
        {selectedStages && selectedStages.length > 0 && (
          <React.Fragment>
            <Separator className="mx-1 h-3" orientation="vertical" />
            <Badge className="px-1 text-xs font-normal p-0" variant="secondary">
              {selectedStages.length} selected
            </Badge>
          </React.Fragment>
        )}
      </DropdownMenuSubTrigger>
      <DropdownMenuPortal>
        <DropdownMenuSubContent className="mx-2 w-56">
          {todoThreadStages.map((stage) => (
            <DropdownMenuCheckboxItem
              key={stage}
              checked={selectedStages ? selectedStages.includes(stage) : false}
              onCheckedChange={(checked) => {
                checked ? onChecked(stage) : onUnchecked(stage);
              }}
              onSelect={(e) => e.preventDefault()}
            >
              <div className="flex items-center gap-x-2">
                {stageIcon(stage, {
                  className: "w-4 h-4 text-indigo-500",
                })}
                <span>{threadStageHumanized(stage)}</span>
              </div>
            </DropdownMenuCheckboxItem>
          ))}
        </DropdownMenuSubContent>
      </DropdownMenuPortal>
    </DropdownMenuSub>
  );
}

function PrioritiesSubMenu({
  priorities,
  onChecked = () => {},
  onUnchecked = () => {},
}: {
  priorities: PrioritiesFiltersType;
  onChecked?: (priority: string) => void;
  onUnchecked?: (priority: string) => void;
}) {
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
        <DropdownMenuSubContent className="mx-2 w-48" alignOffset={4}>
          <DropdownMenuCheckboxItem
            checked={
              selectedPriorities ? selectedPriorities.includes("urgent") : false
            }
            onCheckedChange={(checked) => {
              checked ? onChecked("urgent") : onUnchecked("urgent");
            }}
            onSelect={(e) => e.preventDefault()}
          >
            <div className="flex items-center gap-1">
              <PriorityIcons.urgent className="h-5 w-5" />
              <span>Urgent</span>
            </div>
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
            <div className="flex items-center gap-1">
              <PriorityIcons.high className="h-5 w-5" />
              <span>High</span>
            </div>
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
            <div className="flex items-center gap-1">
              <PriorityIcons.normal className="h-5 w-5" />
              <span>Normal</span>
            </div>
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
            <div className="flex items-center gap-1">
              <PriorityIcons.low className="h-5 w-5" />
              <span>Low</span>
            </div>
          </DropdownMenuCheckboxItem>
        </DropdownMenuSubContent>
      </DropdownMenuPortal>
    </DropdownMenuSub>
  );
}

function AssigneeSubMenu({
  assignedMembers,
  assignees,
  onChecked = () => {},
  onUnchecked = () => {},
}: {
  assignedMembers: Assignee[];
  assignees: AssigneesFiltersType;
  onChecked?: (member: string) => void;
  onUnchecked?: (member: string) => void;
}) {
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
        <DropdownMenuSubContent className="mx-2" alignOffset={4}>
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
                    <div className="flex gap-x-2">
                      <Avatar className="h-5 w-5">
                        <AvatarImage
                          src={`https://avatar.vercel.sh/${m.assigneeId}`}
                        />
                        <AvatarFallback>M</AvatarFallback>
                      </Avatar>
                      <span>{m.name}</span>
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
  stages,
  priorities,
  assignedMembers,
  assignees,
  statusOnChecked = () => {},
  statusOnUnchecked = () => {},
  priorityOnChecked = () => {},
  priorityOnUnchecked = () => {},
  assigneeOnChecked = () => {},
  assigneeOnUnchecked = () => {},
  disableAssigneeFilter = false,
}: {
  stages: StagesFiltersType;
  priorities: PrioritiesFiltersType;
  assignedMembers: Assignee[];
  assignees: AssigneesFiltersType;
  statusOnChecked?: (stage: string) => void;
  statusOnUnchecked?: (stage: string) => void;
  priorityOnChecked?: (priority: string) => void;
  priorityOnUnchecked?: (priority: string) => void;
  assigneeOnChecked?: (member: string) => void;
  assigneeOnUnchecked?: (member: string) => void;
  disableAssigneeFilter?: boolean;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button className="border-dashed" size="sm" variant="outline">
          <MixerHorizontalIcon className="mr-1 h-3 w-3" />
          Filters
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-48 mx-1">
        <DropdownMenuGroup>
          <StagesSubMenu
            stages={stages}
            onChecked={statusOnChecked}
            onUnchecked={statusOnUnchecked}
          />
          <PrioritiesSubMenu
            priorities={priorities}
            onChecked={priorityOnChecked}
            onUnchecked={priorityOnUnchecked}
          />
        </DropdownMenuGroup>
        {!disableAssigneeFilter && (
          <React.Fragment>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <AssigneeSubMenu
                assignedMembers={assignedMembers}
                assignees={assignees}
                onChecked={assigneeOnChecked}
                onUnchecked={assigneeOnUnchecked}
              />
            </DropdownMenuGroup>
          </React.Fragment>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
