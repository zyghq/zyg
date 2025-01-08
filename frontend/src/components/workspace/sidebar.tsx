import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar";
import {
  ExitIcon,
  GearIcon,
  OpenInNewWindowIcon,
  ReaderIcon,
  WidthIcon,
} from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import { Building2Icon, ChevronsUpDown } from "lucide-react";
import * as React from "react";

function WorkspaceMenu({
  email,
  workspaceId,
  workspaceName,
}: {
  email: string;
  workspaceId: string;
  workspaceName: string;
}) {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              size="lg"
            >
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <Building2Icon className="size-4" />
              </div>
              <div className="flex flex-col gap-0.5 leading-none">
                <span className="font-semibold">{workspaceName}</span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="start"
            className="w-[--radix-dropdown-menu-trigger-width]"
          >
            <DropdownMenuLabel className="text-muted-foreground">
              {email}
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem asChild>
              <Link
                params={{ workspaceId }}
                to="/workspaces/$workspaceId/settings"
              >
                <GearIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Settings</div>
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <a href="https://zyg.ai/docs/" rel="noreferrer" target="_blank">
                <ReaderIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Documentation</div>
                <OpenInNewWindowIcon className="my-auto ml-2 h-4 w-4" />
              </a>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem asChild>
              <Link to="/workspaces">
                <WidthIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Switch Workspace</div>
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link to="/signout">
                <ExitIcon className="my-auto mr-2 h-4 w-4" />
                <div className="my-auto">Sign Out</div>
              </Link>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}

// This is sample data.
const data = {
  navMain: [
    {
      items: [
        {
          title: "Installation",
          url: "#",
        },
        {
          title: "Project Structure",
          url: "#",
        },
      ],
      title: "Getting Started",
      url: "#",
    },
    {
      items: [
        {
          title: "Routing",
          url: "#",
        },
        {
          isActive: true,
          title: "Data Fetching",
          url: "#",
        },
        {
          title: "Rendering",
          url: "#",
        },
        {
          title: "Caching",
          url: "#",
        },
        {
          title: "Styling",
          url: "#",
        },
        {
          title: "Optimizing",
          url: "#",
        },
        {
          title: "Configuring",
          url: "#",
        },
        {
          title: "Testing",
          url: "#",
        },
        {
          title: "Authentication",
          url: "#",
        },
        {
          title: "Deploying",
          url: "#",
        },
        {
          title: "Upgrading",
          url: "#",
        },
        {
          title: "Examples",
          url: "#",
        },
      ],
      title: "Building Your Application",
      url: "#",
    },
    {
      items: [
        {
          title: "Components",
          url: "#",
        },
        {
          title: "File Conventions",
          url: "#",
        },
        {
          title: "Functions",
          url: "#",
        },
        {
          title: "next.config.js Options",
          url: "#",
        },
        {
          title: "CLI",
          url: "#",
        },
        {
          title: "Edge Runtime",
          url: "#",
        },
      ],
      title: "API Reference",
      url: "#",
    },
    {
      items: [
        {
          title: "Accessibility",
          url: "#",
        },
        {
          title: "Fast Refresh",
          url: "#",
        },
        {
          title: "Next.js Compiler",
          url: "#",
        },
        {
          title: "Supported Browsers",
          url: "#",
        },
        {
          title: "Turbopack",
          url: "#",
        },
      ],
      title: "Architecture",
      url: "#",
    },
  ],
  versions: ["1.0.1", "1.1.0-alpha", "2.0.0-beta1"],
};

type WorkspaceSidebarProps = React.ComponentProps<typeof Sidebar> & {
  email: string;
  workspaceId: string;
  workspaceName: string;
};

export function WorkspaceSidebar({
  email,
  workspaceId,
  workspaceName,
  ...props
}: WorkspaceSidebarProps) {
  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <WorkspaceMenu
          email={email}
          workspaceId={workspaceId}
          workspaceName={workspaceName}
        />
      </SidebarHeader>
      <SidebarContent>
        {/* We create a SidebarGroup for each parent. */}
        {data.navMain.map((item) => (
          <SidebarGroup key={item.title}>
            <SidebarGroupLabel>{item.title}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {item.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild isActive={item.isActive}>
                      <a href={item.url}>{item.title}</a>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  );
}
