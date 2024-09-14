import React from "react";
import { Button } from "@/components/ui/button";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { TagsIcon } from "lucide-react";
import { CaretSortIcon } from "@radix-ui/react-icons";
import { LabelMetrics } from "@/db/models";
import { LabelThreadsLink } from "@/components/workspace/sidenav-thread-links";

export function SideNavLabelLinks({
  workspaceId,
  labels,
}: {
  workspaceId: string;
  labels: LabelMetrics[];
}) {
  const [isOpen, setIsOpen] = React.useState(false);
  return (
    <Collapsible
      open={isOpen}
      onOpenChange={setIsOpen}
      className="w-full space-y-2"
    >
      <div className="flex items-center justify-between space-x-1">
        <Button variant="ghost" className="w-full pl-3">
          <div className="mr-auto flex">
            <TagsIcon className="my-auto mr-2 h-4 w-4" />
            <div className="font-normal">Labels</div>
          </div>
        </Button>
        <CollapsibleTrigger asChild>
          <Button variant="ghost" size="icon">
            <CaretSortIcon className="h-4 w-4" />
            <span className="sr-only">Toggle</span>
          </Button>
        </CollapsibleTrigger>
      </div>
      <CollapsibleContent className="space-y-1">
        {labels.map((label) => (
          <LabelThreadsLink
            key={label.labelId}
            workspaceId={workspaceId}
            label={label}
          />
        ))}
      </CollapsibleContent>
    </Collapsible>
  );
}
