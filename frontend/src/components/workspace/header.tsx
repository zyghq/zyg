import { Icons } from "@/components/icons";
import { ThemeToggler } from "@/components/theme-toggler";
import { buttonVariants } from "@/components/ui/button";
import SideNavMobileLinks from "@/components/workspace/sidenav-mobile-links";
import { WorkspaceMetrics } from "@/db/models";
import { SortBy } from "@/db/store";
import { cn } from "@/lib/utils";
import { Link } from "@tanstack/react-router";
import { ArrowLeftRightIcon } from "lucide-react";

export function Header({
  email,
  memberId,
  metrics,
  sort,
  workspaceId,
  workspaceName,
}: {
  email: string;
  memberId: string;
  metrics: WorkspaceMetrics;
  sort: SortBy,
  workspaceId: string;
  workspaceName: string;
}) {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex h-14 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 bottom-[calc(100vh-theme(spacing.14))]">
      <div className="mx-4 flex w-full items-center">
        <div className="hidden md:flex">
          <Link className="flex items-center space-x-2" to="/">
            <Icons.logo className="h-5 w-5" />
            <span className="hidden font-semibold sm:inline-block">Zyg.</span>
          </Link>
        </div>
        <SideNavMobileLinks
          email={email}
          memberId={memberId}
          metrics={metrics}
          sort={sort}
          workspaceId={workspaceId}
          workspaceName={workspaceName}
        />
        <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
          <div className="w-full flex-1 md:w-auto md:flex-none"></div>
          <nav className="flex items-center">
            <Link to={"/workspaces"}>
              <div
                className={cn(
                  buttonVariants({
                    variant: "ghost",
                  }),
                  "w-9 px-0"
                )}
              >
                <ArrowLeftRightIcon className="h-4 w-4" />
                <span className="sr-only">Switch Workspace</span>
              </div>
            </Link>
            <ThemeToggler />
          </nav>
        </div>
      </div>
    </header>
  );
}
