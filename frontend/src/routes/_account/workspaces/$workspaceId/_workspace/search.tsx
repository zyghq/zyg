import { Input } from "@/components/ui/input";
import { createFileRoute } from "@tanstack/react-router";
import { Search } from "lucide-react";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/search"
)({
  component: SearchComponent,
});

function SearchComponent() {
  return (
    <div className="mt-4">
      <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
          <Input
            className="w-full pl-10 pr-4 py-2 text-sm rounded-lg"
            placeholder="Search..."
            type="search"
          />
        </div>
      </div>
      <div>
        <div className="container mt-4 sm:mt-24 max-w-md">
          <div className="border p-4 rounded-xl">
            <div className="font-medium text-sm">
              Search for threads, customers, and more.
            </div>
            <div className="text-muted-foreground text-sm">
              Search by thread content, customer name, email or external ID.
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
