import * as React from "react";
import { createFileRoute, Outlet, Link } from "@tanstack/react-router";
import { ArrowLeftIcon } from "@radix-ui/react-icons";
import { SideNavLinks } from "@/components/workspace/settings/sidenav-links";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { SideNavMobile } from "@/components/workspace/settings/sidenav-mobile";
import { useStore } from "zustand";
import { useAccountStore } from "@/providers";

export const Route = createFileRoute("/_auth/workspaces/$workspaceId/settings")(
  {
    component: SettingsLayout,
  }
);

function SettingsLayout() {
  const { workspaceId } = Route.useParams();
  const accountStore = useAccountStore();
  const accountId = useStore(accountStore, (state) =>
    state.getAccountId(state)
  );

  return (
    <React.Fragment>
      <header className="fixed top-0 left-0 right-0 z-50 flex h-14 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 bottom-[calc(100vh-theme(spacing.14))]">
        <div className="mx-4 flex w-full items-center">
          <div className="hidden md:flex">
            <Link
              to="/workspaces/$workspaceId"
              params={{ workspaceId }}
              className={cn(buttonVariants({ variant: "outline", size: "sm" }))}
            >
              <ArrowLeftIcon className="mr-2 h-4 w-4" />
              Settings
            </Link>
          </div>
          <SideNavMobile accountId={accountId} workspaceId={workspaceId} />
        </div>
      </header>
      <div className="flex min-h-screen">
        <aside className="sticky top-14 h-[calc(100vh-theme(spacing.14))] w-80 overflow-y-auto">
          <SideNavLinks
            accountId={accountId}
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
