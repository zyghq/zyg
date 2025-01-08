import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
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
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar";
import { WorkspaceMetrics } from "@/db/models";
import { SortBy } from "@/db/store";
import {
  CircleIcon,
  ExitIcon,
  GearIcon,
  OpenInNewWindowIcon,
  PersonIcon,
  ReaderIcon,
  WidthIcon,
} from "@radix-ui/react-icons";
import { Link } from "@tanstack/react-router";
import {
  Building2Icon,
  ChartColumnIncreasing,
  ChevronsUpDown,
  Search,
} from "lucide-react";
import * as React from "react";

type WorkspaceSidebarProps = React.ComponentProps<typeof Sidebar> & {
  email: string;
  memberId: string;
  metrics: WorkspaceMetrics;
  sort: SortBy;
  workspaceId: string;
  workspaceName: string;
};

export function WorkspaceSidebar({
  email,
  memberId,
  metrics,
  sort,
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
        <SidebarGroup>
          <SidebarGroupLabel>Workspace</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link
                    activeOptions={{ exact: true, includeSearch: false }}
                    activeProps={{
                      className:
                        "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
                    }}
                    params={{ workspaceId }}
                    search={{ sort }}
                    to="/workspaces/$workspaceId/search"
                  >
                    {({ isActive }: { isActive: boolean }) => (
                      <>
                        {isActive ? (
                          <>
                            <Search className="h-5 w-5" />
                            <span className="font-semibold">Search</span>
                          </>
                        ) : (
                          <>
                            <Search className="h-5 w-5 text-muted-foreground" />
                            <span className="font-normal">Search</span>
                          </>
                        )}
                      </>
                    )}
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link
                    activeOptions={{ exact: true, includeSearch: false }}
                    activeProps={{
                      className:
                        "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
                    }}
                    params={{ workspaceId }}
                    search={{ sort }}
                    to="/workspaces/$workspaceId/insights"
                  >
                    {({ isActive }: { isActive: boolean }) => (
                      <>
                        {isActive ? (
                          <>
                            <ChartColumnIncreasing className="h-5 w-5" />
                            <span className="font-semibold">Insights</span>
                          </>
                        ) : (
                          <>
                            <ChartColumnIncreasing className="h-5 w-5 text-muted-foreground" />
                            <span className="font-normal">Insights</span>
                          </>
                        )}
                      </>
                    )}
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>Threads</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link
                    activeOptions={{ exact: true, includeSearch: false }}
                    activeProps={{
                      className:
                        "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
                    }}
                    params={{ workspaceId }}
                    search={{ sort }}
                    to="/workspaces/$workspaceId/threads/me"
                  >
                    {({ isActive }: { isActive: boolean }) => (
                      <>
                        {isActive ? (
                          <>
                            <Avatar className="h-5 w-5">
                              <AvatarImage
                                alt={memberId}
                                src={`https://avatar.vercel.sh/${memberId}`}
                              />
                              <AvatarFallback>U</AvatarFallback>
                            </Avatar>
                            <span className="font-semibold">Your Threads</span>
                          </>
                        ) : (
                          <>
                            <Avatar className="h-5 w-5">
                              <AvatarImage
                                alt={memberId}
                                src={`https://avatar.vercel.sh/${memberId}`}
                              />
                              <AvatarFallback>CN</AvatarFallback>
                            </Avatar>
                            <span className="font-normal">Your Threads</span>
                          </>
                        )}
                      </>
                    )}
                  </Link>
                </SidebarMenuButton>
                <SidebarMenuBadge>{metrics.assignedToMe}</SidebarMenuBadge>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link
                    activeOptions={{ exact: true, includeSearch: false }}
                    activeProps={{
                      className:
                        "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
                    }}
                    params={{ workspaceId }}
                    search={{ sort }}
                    to="/workspaces/$workspaceId/threads/unassigned"
                  >
                    {({ isActive }: { isActive: boolean }) => (
                      <>
                        {isActive ? (
                          <>
                            <PersonIcon className="h-5 w-5" />
                            <span className="font-semibold">
                              Unassigned Threads
                            </span>
                          </>
                        ) : (
                          <>
                            <PersonIcon className="h-5 w-5 text-muted-foreground" />
                            <span className="font-normal">
                              Unassigned Threads
                            </span>
                          </>
                        )}
                      </>
                    )}
                  </Link>
                </SidebarMenuButton>
                <SidebarMenuBadge>{metrics.unassigned}</SidebarMenuBadge>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link
                    activeOptions={{ exact: true, includeSearch: false }}
                    activeProps={{
                      className:
                        "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
                    }}
                    params={{ workspaceId }}
                    search={{ sort }}
                    to="/workspaces/$workspaceId/threads/todo"
                  >
                    {({ isActive }: { isActive: boolean }) => (
                      <>
                        {isActive ? (
                          <>
                            <CircleIcon className="h-5 w-5 text-indigo-500" />
                            <span className="font-semibold">Todo</span>
                          </>
                        ) : (
                          <>
                            <CircleIcon className="h-5 w-5 text-indigo-500" />
                            <span className="font-normal">Todo</span>
                          </>
                        )}
                      </>
                    )}
                  </Link>
                </SidebarMenuButton>
                <SidebarMenuBadge>{metrics.active}</SidebarMenuBadge>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  );
}

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
