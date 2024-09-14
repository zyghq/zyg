import React from "react";
import { getRouteApi, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { DoubleArrowUpIcon } from "@radix-ui/react-icons";

const routeApi = getRouteApi("/_account/workspaces/$workspaceId/_workspace");

export function Sorts() {
  const navigate = useNavigate();
  const [selectedSort, setSelectedSort] = React.useState<string>("");
  const routeSearch = routeApi.useSearch();
  const { sort } = routeSearch;

  React.useEffect(() => {
    setSelectedSort(sort || "");
  }, [sort]);

  const sortDescription = (v: string) => {
    if (status === "todo" && v === "created-dsc")
      return "In Todo, Latest First";
    if (status === "todo" && v === "created-asc")
      return "In Todo, Oldest First";
    if (status === "done" && v === "created-dsc")
      return "In Done, Latest First";
    if (status === "done" && v === "created-asc")
      return "In Done, Oldest First";
    if (status === "snoozed" && v === "created-dsc")
      return "In Snoozed, Latest First";
    if (status === "snoozed" && v === "created-asc")
      return "In Snoozed, Oldest First";
    return "";
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="border-dashed">
          <DoubleArrowUpIcon className="mr-1 h-3 w-3" />
          Sort
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-58 mx-1" align="end">
        <DropdownMenuRadioGroup
          onSelect={(e) => e.preventDefault()}
          value={selectedSort}
          onValueChange={(value) =>
            navigate({
              search: (prev: any) => {
                return { ...prev, sort: value };
              },
            })
          }
        >
          <DropdownMenuRadioItem
            onSelect={(e) => e.preventDefault()}
            value="last-message-dsc"
          >
            Most Recent Message
          </DropdownMenuRadioItem>
          <DropdownMenuRadioItem
            onSelect={(e) => e.preventDefault()}
            value="created-dsc"
          >
            {sortDescription("created-dsc")}
          </DropdownMenuRadioItem>
          <DropdownMenuRadioItem
            onSelect={(e) => e.preventDefault()}
            value="created-asc"
          >
            {sortDescription("created-asc")}
          </DropdownMenuRadioItem>
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
