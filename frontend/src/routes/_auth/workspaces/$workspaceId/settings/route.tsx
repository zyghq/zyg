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
    <div className="flex flex-col">
      <header className="sticky top-0 z-50 flex h-14 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
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
      <div className="flex flex-col">
        <div className="flex">
          <div className="hidden min-w-80 flex-col border-r lg:flex">
            <SideNavLinks
              accountId={accountId}
              maxHeight="h-[calc(100dvh-8rem)]"
            />
          </div>
          <main className="flex-1">
            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
}
