import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { ChatBubbleIcon, OpenInNewWindowIcon } from "@radix-ui/react-icons";
import {
  ChevronsUpDown,
  GitGraphIcon,
  LifeBuoyIcon,
  UsersIcon,
} from "lucide-react";

export function FooterMenu() {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              size="lg"
            >
              <LifeBuoyIcon />
              Support
              <ChevronsUpDown className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="start"
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            sideOffset={4}
          >
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <UsersIcon className="mr-2 h-4 w-4" />
                Join Community
                <OpenInNewWindowIcon className="ml-2 h-4 w-4" />
              </DropdownMenuItem>
              <DropdownMenuItem>
                <ChatBubbleIcon className="mr-2 h-4 w-4" />
                Chat with Us
              </DropdownMenuItem>
              <DropdownMenuItem>
                <GitGraphIcon className="mr-2 h-4 w-4" />
                Changelog
                <OpenInNewWindowIcon className="ml-2 h-4 w-4" />
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
