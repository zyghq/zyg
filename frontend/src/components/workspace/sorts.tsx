import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { sortKeys, ThreadSortKeyHumanized } from "@/db/helpers";
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
          <DoubleArrowUpIcon className="mr-1 h-3 w-3" />
          Sort
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-64 mx-1">
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
              {ThreadSortKeyHumanized(sortKey)}
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
