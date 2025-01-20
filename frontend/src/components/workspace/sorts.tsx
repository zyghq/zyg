import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { sortKeys, threadSortKeyHumanized } from "@/db/helpers";
import { DoubleArrowUpIcon } from "@radix-ui/react-icons";
import React from "react";

export function Sorts({
  onChecked = () => {},
  sort,
}: {
  onChecked?: (sort: string) => void;
  sort: string;
}) {
  const [selectedSort, setSelectedSort] = React.useState<string>("");
  React.useEffect(() => {
    setSelectedSort(sort);
  }, [sort]);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button className="border-dashed" size="sm" variant="outline">
          <DoubleArrowUpIcon className="mr-1 h-2 w-2" />
          Sort
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="mx-1 w-64">
        <DropdownMenuRadioGroup
          onSelect={(e) => e.preventDefault()}
          onValueChange={(value) => onChecked(value)}
          value={selectedSort}
        >
          {sortKeys.map((sortKey) => (
            <DropdownMenuRadioItem
              key={sortKey}
              onSelect={(e) => e.preventDefault()}
              value={sortKey}
            >
              {threadSortKeyHumanized(sortKey)}
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
