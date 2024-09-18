import { buttonVariants } from "@/components/ui/button";
import SideNavLinks from "@/components/workspace/settings/sidenav-links";
import SideNavMobileLinks from "@/components/workspace/settings/sidenav-mobile-links";
import { cn } from "@/lib/utils";
import { useAccountStore } from "@/providers";
import { ArrowLeftIcon } from "@radix-ui/react-icons";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings"
)({
  component: SettingsLayout,
});

function SettingsLayout() {
  const { workspaceId } = Route.useParams();
  const accountStore = useAccountStore();
  const accountId = useStore(accountStore, (state) =>
    state.getAccountId(state)
  );
  const accountName = useStore(accountStore, (state) => state.getName(state));

  return (
    <React.Fragment>
      <header className="fixed top-0 left-0 right-0 z-50 flex h-14 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 bottom-[calc(100vh-theme(spacing.14))]">
        <div className="mx-4 flex w-full items-center">
          <div className="hidden md:flex">
            <Link
              className={cn(buttonVariants({ size: "sm", variant: "outline" }))}
              params={{ workspaceId }}
              to="/workspaces/$workspaceId"
            >
              <ArrowLeftIcon className="mr-2 h-4 w-4" />
              Settings
            </Link>
          </div>
          <SideNavMobileLinks
            accountId={accountId}
            accountName={accountName}
            workspaceId={workspaceId}
          />
        </div>
      </header>
      <div className="flex min-h-screen">
        <aside className="hidden sticky top-14 h-[calc(100vh-theme(spacing.14))] w-80 overflow-y-auto md:block md:border-r">
          <SideNavLinks
            accountId={accountId}
            accountName={accountName}
            maxHeight="h-[calc(100dvh-8rem)]"
          />
        </aside>
        <main className="flex-1 mt-14 pb-4">
          <Outlet />
        </main>
      </div>
    </React.Fragment>
  );
}
