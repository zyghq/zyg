import { cn } from "@/lib/utils";
import { Link } from "@tanstack/react-router";
import { buttonVariants } from "@/components/ui/button";
import { Icons } from "@/components/icons";
import { ThemeToggler } from "@/components/theme-toggler";
import { ArrowLeftRightIcon } from "lucide-react";
import { WorkspaceMetrics } from "@/db/entities";
import SideNavMobile from "@/components/workspace/sidenav-mobile";

export function Header({
  email,
  workspaceId,
  workspaceName,
  metrics,
  memberId,
}: {
  email: string;
  workspaceId: string;
  workspaceName: string;
  metrics: WorkspaceMetrics;
  memberId: string;
}) {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex h-14 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 bottom-[calc(100vh-theme(spacing.14))]">
      <div className="mx-4 flex w-full items-center">
        <div className="hidden md:flex">
          <Link to="/" className="flex items-center space-x-2">
            <Icons.logo className="h-5 w-5" />
            <span className="hidden font-semibold sm:inline-block">Zyg.</span>
          </Link>
        </div>
        <SideNavMobile
          email={email}
          workspaceId={workspaceId}
          workspaceName={workspaceName}
          metrics={metrics}
          memberId={memberId}
        />
        <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
          <div className="w-full flex-1 md:w-auto md:flex-none">...</div>
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
