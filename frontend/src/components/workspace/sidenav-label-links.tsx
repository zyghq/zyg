import React from "react";
import { cn } from "@/lib/utils";
import { Link, getRouteApi } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { buttonVariants } from "@/components/ui/button";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";

import { TagIcon, TagsIcon } from "lucide-react";
import { CaretSortIcon } from "@radix-ui/react-icons";
import { LabelMetrics } from "@/db/models";

const routeApi = getRouteApi("/_auth/workspaces/$workspaceId/_workspace");

export function SideNavLabelLinks({
  workspaceId,
  labels,
}: {
  workspaceId: string;
  labels: LabelMetrics[];
}) {
  const routeSearch = routeApi.useSearch();
  const { status, sort } = routeSearch;
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
          <Link
            key={label.labelId}
            to="/workspaces/$workspaceId/labels/$labelId"
            params={{ workspaceId, labelId: label.labelId }}
            search={{ status: status, sort: sort }}
            className={cn(
              buttonVariants({ variant: "ghost" }),
              "flex w-full justify-between px-3 dark:text-accent-foreground"
            )}
            activeOptions={{ exact: true }}
            activeProps={{
              className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
            }}
          >
            {({ isActive }) => (
              <>
                {isActive ? (
                  <>
                    <div className="flex">
                      <TagIcon className="my-auto mr-1 h-3 w-3 text-muted-foreground" />
                      <div className="font-normal capitalize text-foreground">
                        {label.name}
                      </div>
                    </div>
                    <div className="font-mono font-light text-muted-foreground">
                      {label.count}
                    </div>
                  </>
                ) : (
                  <>
                    <div className="flex">
                      <TagIcon className="my-auto mr-1 h-3 w-3 text-muted-foreground" />
                      <div className="font-normal capitalize text-foreground">
                        {label.name}
                      </div>
                    </div>
                    <div className="font-mono font-light text-muted-foreground">
                      {label.count}
                    </div>
                  </>
                )}
              </>
            )}
          </Link>
        ))}
      </CollapsibleContent>
    </Collapsible>
  );
}
