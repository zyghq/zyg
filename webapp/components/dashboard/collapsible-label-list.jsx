"use client";

import * as React from "react";

import Link from "next/link";

import { Button } from "@/components/ui/button";
import { buttonVariants } from "@/components/ui/button";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";

import { CaretSortIcon } from "@radix-ui/react-icons";

import { TagIcon, TagsIcon } from "lucide-react";

export default function CollapsibleLabelList({ workspaceId, labels }) {
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
          <Button
            key={label.labelId}
            variant="ghost"
            asChild
            className="flex justify-between pr-3"
          >
            <Link href={`/${workspaceId}/threads/labels/${label.labelId}/`}>
              <div className="flex">
                <TagIcon className="my-auto mr-1 h-3 w-3 text-muted-foreground" />
                <div className="font-normal capitalize text-foreground">
                  {label.name}
                </div>
              </div>
              <div className="font-mono font-light text-muted-foreground">
                {label.count}
              </div>
            </Link>
          </Button>
        ))}
      </CollapsibleContent>
    </Collapsible>
  );
}
