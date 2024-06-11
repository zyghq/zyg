import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Pencil1Icon, DotsHorizontalIcon } from "@radix-ui/react-icons";
import { ShieldAlertIcon } from "lucide-react";

// this component will not only render readme label, based on actions
// like edit, archive it will also render the form to edit the label
export function LabelItem({ label }: { label: string }) {
  return (
    <div className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left">
      <div className="flex w-full flex-col">
        <div className="flex items-center">
          <div className="flex items-center gap-2">
            <ShieldAlertIcon className="h-4 w-4 text-muted-foreground" />
            <div className="font-normal">{label}</div>
          </div>
          <div className="ml-auto">
            <Button variant="ghost" size="sm">
              <Pencil1Icon className="h-4 w-4" />
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm">
                  <DotsHorizontalIcon className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="mx-1">
                <DropdownMenuItem>
                  <div className="mr-4">Copy Label ID</div>
                  <DropdownMenuShortcut>⌘C</DropdownMenuShortcut>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <div className="mr-4">Archive</div>
                  <DropdownMenuShortcut>⌘D</DropdownMenuShortcut>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>
    </div>
  );
}
